/*
Copyright 2021 Polyglot Systems.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"strconv"
	"time"

	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
	defaults "github.com/kenmoini/certificate-sentinel-operator/controllers/defaults"
	helpers "github.com/kenmoini/certificate-sentinel-operator/controllers/helpers"
)

//===========================================================================================
// SPOT TYPES
//===========================================================================================

// KeystoreSentinelReconciler reconciles a KeystoreSentinel object
type KeystoreSentinelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//===========================================================================================
// INIT VARS
//===========================================================================================

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &KeystoreSentinelReconciler{}
var lggrK = log.Log.WithName("keystore-sentinel-controller")
var setLogLevelK int

//===========================================================================================
// INIT FUNC
//===========================================================================================

// init is fired when this controller is started
func init() {
	// Set logging
	//log.SetLogger(zap.New())
}

//===========================================================================================
// RBAC GENERATORS
//===========================================================================================

//+kubebuilder:rbac:groups=config.polyglot.systems,resources=keystoresentinels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.polyglot.systems,resources=keystoresentinels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.polyglot.systems,resources=keystoresentinels/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

//===========================================================================================
// RECONCILE
//===========================================================================================

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeystoreSentinel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KeystoreSentinelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentConfig, _ := config.GetConfig()
	clusterEndpoint := currentConfig.Host
	apiPath := currentConfig.APIPath
	statusLists := configv1.KeystoreSentinelStatus{
		DiscoveredKeystoreCertificates: []configv1.KeystoreInformation{},
	}

	LogWithLevel("Connecting to:"+clusterEndpoint, 3, lggrK, setLogLevelK)
	LogWithLevel("API Path:"+apiPath, 3, lggrK, setLogLevelK)

	// Fetch the KeystoreSentinel instance
	keystoreSentinel := &configv1.KeystoreSentinel{}
	err := r.Get(ctx, req.NamespacedName, keystoreSentinel)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			lggrK.Error(err, "KeystoreSentinel resource not found on the cluster.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		lggrK.Error(err, "Failed to get KeystoreSentinel")
		return ctrl.Result{}, err
	}

	// Log the used keystoreSentinel
	LogWithLevel("KeystoreSentinel loaded!  Found '"+keystoreSentinel.Name+"' in 'namespace/"+keystoreSentinel.Namespace+"'", 1, lggrK, setLogLevelK)
	setLogLevelK = defaults.SetDefaultInt(defaults.LogLevel, keystoreSentinel.Spec.LogLevel)

	// Set default vars
	var CertHashList []string
	keystoreCount := 0
	keystoreAtRisk := false
	var expiredKeystoreCount int
	var expiredKeystoreCertificatesCount int

	scanningInterval := defaults.SetDefaultInt(defaults.ScanningInterval, keystoreSentinel.Spec.ScanningInterval)
	targetName := keystoreSentinel.Spec.Target.TargetName

	serviceAccount := keystoreSentinel.Spec.Target.ServiceAccount
	targetKind := keystoreSentinel.Spec.Target.Kind
	targetAPIVersion := keystoreSentinel.Spec.Target.APIVersion
	targetDaysOut := keystoreSentinel.Spec.Target.DaysOut
	timeOut := DaysOutToTimeOut(targetDaysOut)

	targetNamespaceLabels := keystoreSentinel.Spec.Target.NamespaceLabels
	targetLabels := keystoreSentinel.Spec.Target.TargetLabels
	targetLabelSelector, targetNamespaceLabelSelector := SetupLabelSelectors(targetLabels, targetNamespaceLabels, lggrK)

	LogWithLevel("Processing KeystoreSentinel target: "+targetName, 2, lggrK, setLogLevelK)

	// Set up new Client
	cl, err := SetupNewClient(lggrK, setLogLevelK, r.Client, serviceAccount, keystoreSentinel.Namespace, clusterEndpoint, apiPath)
	if err != nil {
		lggrK.Error(err, "Failed to create client")
		lggrK.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
		time.Sleep(time.Second * time.Duration(scanningInterval))
		return ctrl.Result{}, err
	}

	// Figure out what namespaces we actually have to work with
	effectiveNamespaces, _ := SetupNamespaceSlice(keystoreSentinel.Spec.Target.Namespaces, cl, lggrK, setLogLevelK, serviceAccount, targetNamespaceLabelSelector, scanningInterval)

	// Loop through the namespaces in scope for this target
	for _, el := range effectiveNamespaces {

		// Check to make sure we can even access the Keystore password before scanning for them in this namespace
		passwordBytes, err := getPasswordBytesFromSpecTarget(keystoreSentinel.Spec.Target.KeystorePassword, el, cl)
		if err != nil {
			lggrK.Error(err, "Failed to process keystore password source!")
		}
		defer zeroing(passwordBytes)

		targetListOptions := &client.ListOptions{Namespace: el, LabelSelector: targetLabelSelector}

		switch targetKind {
		//=========================== SECRETS
		case "Secret":
			// Get a list of Secrets in this Namespace
			secretList := &corev1.SecretList{}
			err = cl.List(context.Background(), secretList, targetListOptions)
			if err != nil {
				lggrK.Error(err, "Failed to list secrets in namespace/"+el)
			}

			// Loop through Secrets
			for _, e := range secretList.Items {
				secretType := string(e.Type)
				if secretType == string(corev1.SecretTypeOpaque) || secretType == string(corev1.SecretTypeTLS) {
					LogWithLevel("Getting secret/"+e.Name+" in namespace/"+el+" (type="+secretType+")", 3, lggrK, setLogLevelK)

					secretItem, _ := GetSecret(string(e.Name), el, cl)

					// Get the actual secret data
					for k, s := range secretItem.Data {
						// Store the secret as a base64 decoded string from the byte slice
						//sDataStr := string(s)
						keystoreAtRisk = false

						keystoreObj, err := ReadKeyStoreFromBytes(s, passwordBytes)

						if err != nil {
							// No JKS object found in this secret
						} else {
							// Keystore is found
							LogWithLevel("KEYSTORE FOUND! - ns/"+el+" - secret/"+string(e.Name)+" - key:"+k, 3, lggrK, setLogLevelK)
							keystoreCount++

							certs, err := ProcessKeystoreIntoCertificates(keystoreObj)
							if err != nil {
								lggrK.Error(err, "Failed to process keystore into certificates!")
							}
							for keystoreAlias, certSlice := range certs {

								for _, cert := range certSlice {
									// Check to see if this has already been added
									shaStr := createUniqueCertificateChecksum(targetKind+"-"+el+"-"+e.Name+"-"+cert.Subject.CommonName+"-"+cert.Issuer.CommonName, &cert)

									if defaults.ContainsString(CertHashList, shaStr) {
										// Skipping Certificate
										LogWithLevel("Already found "+shaStr, 3, lggrK, setLogLevelK)
									} else {
										// Add + Process
										LogWithLevel("Adding "+shaStr, 3, lggrK, setLogLevelK)
										CertHashList = append(CertHashList, shaStr)

										discovered, messages := helpers.ParseKeystoreCertificateIntoObjects(&cert, timeOut, el, e.Name, k, targetKind, targetAPIVersion, keystoreAlias)
										// Loop through passed messages for log level 3
										for _, m := range messages {
											LogWithLevel(m, 3, lggrK, setLogLevelK)
										}

										// Add discovered keystore to DiscoveredKeystore
										if len(discovered) != 0 {
											for _, iv := range discovered {
												if len(iv.TriggeredDaysOut) > 0 {
													expiredKeystoreCertificatesCount++
													keystoreAtRisk = true
												}
												statusLists.DiscoveredKeystoreCertificates = append(statusLists.DiscoveredKeystoreCertificates, iv)
											}
										}
									}

								}
							}
						}

						if keystoreAtRisk {
							expiredKeystoreCount++
						}

					}
				}
			}
		//=========================== CONFIGMAPS
		case "ConfigMap":
			LogWithLevel("Checking for access to ConfigMap in ns/"+el, 3, lggrK, setLogLevelK)
			// Get the list of ConfigMaps in this namespace
			configMapList := &corev1.ConfigMapList{}
			err = cl.List(context.Background(), configMapList, targetListOptions)
			if err != nil {
				lggrK.Error(err, "Failed to list ConfigMaps in ns/"+el)
			}

			// Loop through ConfigMaps
			for _, e := range configMapList.Items {
				LogWithLevel("Getting configmap/"+e.Name+" in namespace/"+el, 3, lggrK, setLogLevelK)
				configMapItem, _ := GetConfigMap(string(e.Name), el, cl)

				// Loop through the actual ConfigMap data
				for k, cm := range configMapItem.Data {

					keystoreAtRisk = false

					keystoreObj, err := ReadKeyStoreFromBytes([]byte(cm), passwordBytes)

					if err != nil {
						// No JKS object found in this configmap
					} else {
						// Keystore is found
						LogWithLevel("KEYSTORE FOUND! - ns/"+el+" - configmap/"+string(e.Name)+" - key:"+k, 3, lggrK, setLogLevelK)
						keystoreCount++

						certs, err := ProcessKeystoreIntoCertificates(keystoreObj)
						if err != nil {
							lggrK.Error(err, "Failed to process keystore into certificates!")
						}

						for keystoreAlias, certSlice := range certs {
							for _, cert := range certSlice {
								// Check to see if this has already been added
								shaStr := createUniqueCertificateChecksum(targetKind+"-"+el+"-"+e.Name+"-"+cert.Subject.CommonName+"-"+cert.Issuer.CommonName, &cert)

								if defaults.ContainsString(CertHashList, shaStr) {
									// Skipping Certificate
									LogWithLevel("Already found "+shaStr, 3, lggrK, setLogLevelK)
								} else {
									// Add + Process
									LogWithLevel("Adding "+shaStr, 3, lggrK, setLogLevelK)
									CertHashList = append(CertHashList, shaStr)

									discovered, messages := helpers.ParseKeystoreCertificateIntoObjects(&cert, timeOut, el, e.Name, k, targetKind, targetAPIVersion, keystoreAlias)
									// Loop through passed messages for log level 3
									for _, m := range messages {
										LogWithLevel(m, 3, lggrK, setLogLevelK)
									}

									// Add discovered keystore to DiscoveredKeystore
									if len(discovered) != 0 {
										for _, iv := range discovered {
											if len(iv.TriggeredDaysOut) > 0 {
												expiredKeystoreCertificatesCount++
												keystoreAtRisk = true
											}
											//decodedCertificates = append(decodedCertificates, cert)
											statusLists.DiscoveredKeystoreCertificates = append(statusLists.DiscoveredKeystoreCertificates, iv)
										}
									}
								}
							}
						}
					}

					if keystoreAtRisk {
						expiredKeystoreCount++
					}
				}
			}
		//=========================== DEFAULT - INVALID KIND
		default:
			// Unsupported Object Kind
			lggrK.Info("Invalid Target Kind!")
			lggrK.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
			time.Sleep(time.Second * time.Duration(scanningInterval))
			return ctrl.Result{}, nil
		}
	}

	// Set updater check vars
	statusChanged := true
	oldDiscoveredKeystoreCertificates := keystoreSentinel.Status.DiscoveredKeystoreCertificates
	oldLastReportSent := keystoreSentinel.Status.LastReportSent
	oldExpiringCertificates := keystoreSentinel.Status.ExpiringCertificates
	oldKeystoresAtRisk := keystoreSentinel.Status.KeystoresAtRisk
	oldTotalKeystoresFound := keystoreSentinel.Status.TotalKeystoresFound

	// Merge the Certificates into the .status of the KeystoreSentinel object
	keystoreSentinel.Status.DiscoveredKeystoreCertificates = statusLists.DiscoveredKeystoreCertificates
	keystoreSentinel.Status.ExpiringCertificates = expiredKeystoreCertificatesCount
	keystoreSentinel.Status.KeystoresAtRisk = expiredKeystoreCount
	keystoreSentinel.Status.TotalKeystoresFound = keystoreCount

	// Process reports if needed, only if there are new certificates at risk
	if expiredKeystoreCertificatesCount > 0 {
		keystoreSentinel.Status.LastReportSent = processKeystoreReport(*keystoreSentinel, lggrK, r.Client)
	}

	// Check the difference in structs
	if reflect.DeepEqual(oldDiscoveredKeystoreCertificates, keystoreSentinel.Status.DiscoveredKeystoreCertificates) && oldLastReportSent == keystoreSentinel.Status.LastReportSent && oldExpiringCertificates == keystoreSentinel.Status.ExpiringCertificates && oldKeystoresAtRisk == keystoreSentinel.Status.KeystoresAtRisk && oldTotalKeystoresFound == keystoreSentinel.Status.TotalKeystoresFound {
		statusChanged = false
	}

	if statusChanged {
		err = r.Status().Update(ctx, keystoreSentinel)
		if err != nil {
			lggrK.Error(err, "Failed to update KeystoreSentinel status")
			return ctrl.Result{}, err
		}
	}
	LogWithLevel("Found "+strconv.Itoa(len(keystoreSentinel.Status.DiscoveredKeystoreCertificates))+" Certificates in "+strconv.Itoa(keystoreCount)+" Keystores, "+strconv.Itoa(expiredKeystoreCount)+" Keystores of which have "+strconv.Itoa(expiredKeystoreCertificatesCount)+" certificates that are at risk of expiring", 2, lggrK, setLogLevelK)

	// Reconcile for any reason other than an error after 5 seconds
	lggrK.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
	return ctrl.Result{RequeueAfter: time.Second * time.Duration(scanningInterval)}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *KeystoreSentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1.KeystoreSentinel{}).
		Complete(r)
}
