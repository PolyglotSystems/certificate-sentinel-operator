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
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	//"k8s.io/api"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
	defaults "github.com/kenmoini/certificate-sentinel-operator/controllers/defaults"
	helpers "github.com/kenmoini/certificate-sentinel-operator/controllers/helpers"
)

//===========================================================================================
// SPOT TYPES
//===========================================================================================

// CertificateSentinelReconciler reconciles a CertificateSentinel object
type CertificateSentinelReconciler struct {
	// client can be used to retrieve objects from the APIServer with the cached response.
	client.Client
	Scheme *runtime.Scheme
}

//===========================================================================================
// INIT VARS
//===========================================================================================

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &CertificateSentinelReconciler{}
var lggr = log.Log.WithName("certificate-sentinel-controller")
var setLogLevel int

//===========================================================================================
// INIT FUNC
//===========================================================================================

// init is fired when this controller is started
func init() {
	// Set logging
	log.SetLogger(zap.New())
}

//===========================================================================================
// RBAC GENERATORS
//===========================================================================================

//+kubebuilder:rbac:groups=config.polyglot.systems,resources=certificatesentinels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.polyglot.systems,resources=certificatesentinels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.polyglot.systems,resources=certificatesentinels/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

//===========================================================================================
// RECONCILE
//===========================================================================================

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CertificateSentinel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *CertificateSentinelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentConfig, _ := config.GetConfig()
	clusterEndpoint := currentConfig.Host
	apiPath := currentConfig.APIPath
	statusLists := configv1.CertificateSentinelStatus{
		DiscoveredCertificates: []configv1.CertificateInformation{},
	}

	LogWithLevel("Connecting to:"+clusterEndpoint, 3, lggr, setLogLevel)
	LogWithLevel("API Path:"+apiPath, 3, lggr, setLogLevel)

	// Fetch the CertificateSentinel instance
	certificateSentinel := &configv1.CertificateSentinel{}
	err := r.Get(ctx, req.NamespacedName, certificateSentinel)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			lggr.Error(err, "CertificateSentinel resource not found on the cluster.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		lggr.Error(err, "Failed to get CertificateSentinel")
		return ctrl.Result{}, err
	}

	// Log the used certificateSentinel
	LogWithLevel("CertificateSentinel loaded!  Found '"+certificateSentinel.Name+"' in 'namespace/"+certificateSentinel.Namespace+"'", 1, lggr, setLogLevel)
	setLogLevel = defaults.SetDefaultInt(defaults.LogLevel, certificateSentinel.Spec.LogLevel)

	// Set default vars
	scanningInterval := defaults.SetDefaultInt(defaults.ScanningInterval, certificateSentinel.Spec.ScanningInterval)
	targetName := certificateSentinel.Spec.Target.TargetName

	serviceAccount := certificateSentinel.Spec.Target.ServiceAccount
	targetKind := certificateSentinel.Spec.Target.Kind
	targetAPIVersion := certificateSentinel.Spec.Target.APIVersion
	targetDaysOut := certificateSentinel.Spec.Target.DaysOut
	timeOut := DaysOutToTimeOut(targetDaysOut)

	targetLabels := certificateSentinel.Spec.Target.TargetLabels
	targetNamespaceLabels := certificateSentinel.Spec.Target.NamespaceLabels

	targetLabelSelector, targetNamespaceLabelSelector := SetupLabelSelectors(targetLabels, targetNamespaceLabels, lggr)

	CertHashList := []string{}
	expiredCertificateCount := 0

	LogWithLevel("Processing CertificateSentinel target: "+targetName, 2, lggr, setLogLevel)

	// Set up new Client
	cl, err := SetupNewClient(lggr, setLogLevel, r.Client, serviceAccount, certificateSentinel.Namespace, clusterEndpoint, apiPath)
	if err != nil {
		lggr.Error(err, "Failed to create client")
		lggr.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
		time.Sleep(time.Second * time.Duration(scanningInterval))
		return ctrl.Result{}, err
	}

	// Figure out what namespaces we actually have to work with
	effectiveNamespaces, _ := SetupNamespaceSlice(certificateSentinel.Spec.Target.Namespaces, cl, lggr, setLogLevel, serviceAccount, targetNamespaceLabelSelector, scanningInterval)

	// Loop through the namespaces in scope for this target
	for _, el := range effectiveNamespaces {
		targetListOptions := &client.ListOptions{Namespace: el, LabelSelector: targetLabelSelector}
		switch targetKind {
		//=========================== SECRETS
		case "Secret":
			// Get a list of Secrets in this Namespace
			secretList := &corev1.SecretList{}
			err = cl.List(context.Background(), secretList, targetListOptions)
			if err != nil {
				lggr.Error(err, "Failed to list secrets in namespace/"+el)
				return ctrl.Result{}, err
			}

			// Loop through Secrets
			for _, e := range secretList.Items {
				secretType := string(e.Type)
				if secretType == string(corev1.SecretTypeOpaque) || secretType == string(corev1.SecretTypeTLS) {
					LogWithLevel("Getting secret/"+e.Name+" in namespace/"+el+" (type="+secretType+")", 3, lggr, setLogLevel)

					secretItem, _ := GetSecret(string(e.Name), el, cl)

					// Get the actual secret data
					for k, s := range secretItem.Data {
						// Store the secret as a base64 decoded string from the byte slice
						sDataStr := string(s)

						// See if this contains text about a Certificate
						if strings.Contains(sDataStr, "-----BEGIN CERTIFICATE-----") {
							LogWithLevel("CERTIFICATE FOUND! - ns/"+el+" - secret/"+string(e.Name)+" - key:"+k, 3, lggr, setLogLevel)
							certs, _ := helpers.DecodeCertificateBytes(s, lggr)

							// Loop through the current collection of certificates
							for _, cert := range certs {
								// Check to see if this has already been added
								shaStr := createUniqueCertificateChecksum(targetKind+"-"+el+"-"+e.Name+"-"+cert.Subject.CommonName+"-"+cert.Issuer.CommonName, cert)

								if defaults.ContainsString(CertHashList, shaStr) {
									// Skipping Certificate
									LogWithLevel("Already found "+shaStr, 3, lggr, setLogLevel)
								} else {
									// Add + Process
									LogWithLevel("Adding "+shaStr, 3, lggr, setLogLevel)
									CertHashList = append(CertHashList, shaStr)

									discovered, messages := helpers.ParseCertificateIntoObjects(cert, timeOut, el, e.Name, k, targetKind, targetAPIVersion)
									// Loop through passed messages for log level 3
									for _, m := range messages {
										LogWithLevel(m, 3, lggr, setLogLevel)
									}

									// Add decoded certificate to DiscoveredCertificates
									if len(discovered) != 0 {
										for _, iv := range discovered {
											if len(iv.TriggeredDaysOut) > 0 {
												expiredCertificateCount++
											}
											statusLists.DiscoveredCertificates = append(statusLists.DiscoveredCertificates, iv)
										}
									}
								}
							}

						}

					}
				}
			}
		//=========================== CONFIGMAPS
		case "ConfigMap":
			LogWithLevel("Checking for access to ConfigMap in ns/"+el, 3, lggr, setLogLevel)
			// Get the list of ConfigMaps in this namespace
			configMapList := &corev1.ConfigMapList{}
			err = cl.List(context.Background(), configMapList, targetListOptions)
			if err != nil {
				lggr.Error(err, "Failed to list ConfigMaps in ns/"+el)
				return ctrl.Result{}, err
			}

			// Loop through ConfigMaps
			for _, e := range configMapList.Items {
				LogWithLevel("Getting configmap/"+e.Name+" in namespace/"+el, 3, lggr, setLogLevel)
				configMapItem, _ := GetConfigMap(string(e.Name), el, cl)

				// Loop through the actual ConfigMap data
				for k, cm := range configMapItem.Data {
					// See if this contains text about a Certificate
					if strings.Contains(string(cm), "-----BEGIN CERTIFICATE-----") {
						LogWithLevel("CERTIFICATE FOUND! - ns/"+el+" - configmap/"+string(e.Name)+" - key:"+k, 3, lggr, setLogLevel)
						certs, _ := helpers.DecodeCertificateBytes([]byte(cm), lggr)

						// Loop through the current collection of certificates
						for _, cert := range certs {
							// Check to see if this has already been added
							shaStr := createUniqueCertificateChecksum(targetKind+"-"+el+"-"+e.Name+"-"+cert.Subject.CommonName+"-"+cert.Issuer.CommonName, cert)

							if defaults.ContainsString(CertHashList, shaStr) {
								// Skipping Certificate
								LogWithLevel("Already found "+shaStr, 3, lggr, setLogLevel)
							} else {
								// Add + Process
								LogWithLevel("Adding "+shaStr, 3, lggr, setLogLevel)
								CertHashList = append(CertHashList, shaStr)

								discovered, messages := helpers.ParseCertificateIntoObjects(cert, timeOut, el, e.Name, k, targetKind, targetAPIVersion)
								// Loop through passed messages for log level 3
								for _, m := range messages {
									LogWithLevel(m, 3, lggr, setLogLevel)
								}

								// Add decoded certificate to DiscoveredCertificates
								if len(discovered) != 0 {
									for _, iv := range discovered {
										if len(iv.TriggeredDaysOut) > 0 {
											expiredCertificateCount++
										}
										statusLists.DiscoveredCertificates = append(statusLists.DiscoveredCertificates, iv)
									}
								}
							}
						}

					}
				}
			}
		//=========================== DEFAULT - INVALID KIND
		default:
			// Unsupported Object Kind
			lggr.Info("Invalid Target Kind!")
			lggr.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
			time.Sleep(time.Second * time.Duration(scanningInterval))
			return ctrl.Result{}, nil
		}
	}
	//}

	// Set updater check vars
	statusChanged := true
	oldDiscoveredCertificates := certificateSentinel.Status.DiscoveredCertificates
	oldLastReportSent := certificateSentinel.Status.LastReportSent
	oldExpiringCertificates := certificateSentinel.Status.ExpiringCertificates

	// Merge the Certificates into the .status of the CertificateSentinel object
	certificateSentinel.Status.DiscoveredCertificates = statusLists.DiscoveredCertificates
	certificateSentinel.Status.ExpiringCertificates = expiredCertificateCount

	// Process reports if needed, only if there are new certificates at risk
	if expiredCertificateCount > 0 {
		certificateSentinel.Status.LastReportSent = processReport(*certificateSentinel, lggr, r.Client)
	}

	// Check the difference in structs
	if reflect.DeepEqual(oldDiscoveredCertificates, certificateSentinel.Status.DiscoveredCertificates) && oldLastReportSent == certificateSentinel.Status.LastReportSent && oldExpiringCertificates == certificateSentinel.Status.ExpiringCertificates {
		statusChanged = false
	}

	if statusChanged {
		err = r.Status().Update(ctx, certificateSentinel)
		if err != nil {
			lggr.Error(err, "Failed to update CertificateSentinel status")
			return ctrl.Result{}, err
		}
	}
	LogWithLevel("Found "+strconv.Itoa(len(certificateSentinel.Status.DiscoveredCertificates))+" Certificates, "+strconv.Itoa(expiredCertificateCount)+" of which are at risk of expiring", 2, lggr, setLogLevel)

	// Reconcile for any reason other than an error after 5 seconds
	lggr.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
	return ctrl.Result{RequeueAfter: time.Second * time.Duration(scanningInterval)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CertificateSentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1.CertificateSentinel{}).
		Complete(r)
}
