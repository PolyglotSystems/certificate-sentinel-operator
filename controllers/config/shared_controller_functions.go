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
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
	defaults "github.com/kenmoini/certificate-sentinel-operator/controllers/defaults"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LogWithLevel implements simple log levels
func LogWithLevel(s string, level int, l logr.Logger, setLogLevel int) {
	if setLogLevel >= level {
		l.Info(s)
	}
}

// SetupSingleLabelSelector takes the YAML definition of a LabelSelector and creates the actual object to use in filtering lists
func SetupSingleLabelSelector(targetLabels []configv1.LabelSelector) (labels.Selector, error) {
	targetLabelSelectorFN := labels.NewSelector()
	// Create the label selectors for the target label filter
	if len(targetLabels) > 0 {
		// Loop through label filters
		for _, label := range targetLabels {
			activeFilter := returnFilterType(label.Filter)
			req, err := labels.NewRequirement(label.Key, activeFilter, label.Values)
			if err != nil {
				// Failed to build labelSelector requirement for target!
				return labels.NewSelector(), err
			}
			targetLabelSelectorFN = targetLabelSelectorFN.Add(*req)
		}
	}
	return targetLabelSelectorFN, nil
}

// SetupLabelSelectors wraps some shared functions
func SetupLabelSelectors(targetLabels []configv1.LabelSelector, targetNamespaceLabels []configv1.LabelSelector, lgr logr.Logger) (labels.Selector, labels.Selector) {
	targetLabelSelectorFN := labels.NewSelector()
	targetNamespaceLabelSelectorFN := labels.NewSelector()
	// Create the label selectors for the target label filter
	if len(targetLabels) > 0 {
		// Loop through label filters
		for _, label := range targetLabels {
			activeFilter := returnFilterType(label.Filter)
			req, err := labels.NewRequirement(label.Key, activeFilter, label.Values)
			if err != nil {
				lgr.Error(err, "Failed to build labelSelector requirement for target!")
			}
			targetLabelSelectorFN = targetLabelSelectorFN.Add(*req)
		}
	}

	// Create the label selectors for the namespace label filter
	if len(targetNamespaceLabels) > 0 {
		// Loop through label filters
		for _, label := range targetNamespaceLabels {
			activeFilter := returnFilterType(label.Filter)
			req, err := labels.NewRequirement(label.Key, activeFilter, label.Values)
			if err != nil {
				lgr.Error(err, "Failed to build labelSelector requirement for namespace!")
			}
			targetNamespaceLabelSelectorFN = targetNamespaceLabelSelectorFN.Add(*req)
		}
	}
	return targetLabelSelectorFN, targetNamespaceLabelSelectorFN
}

// returnFilterType just returns whatever string type simple name of label selection operations to that actual type
func returnFilterType(labelFilter string) selection.Operator {
	switch labelFilter {
	case "in":
		return selection.In
	case "notIn":
		return selection.NotIn
	case "equals":
		return selection.Equals
	case "doubleEquals":
		return selection.DoubleEquals
	case "notEquals":
		return selection.NotEquals
	case "greaterThan":
		return selection.GreaterThan
	case "lessThan":
		return selection.LessThan
	case "exists":
		return selection.Exists
	case "doesNotExist":
		return selection.DoesNotExist
	default:
		return selection.Equals
	}
}

// SetupNamespaceSlice sets up the shared effectiveNamespaces from the provided YAML structures
func SetupNamespaceSlice(namespaces []string, cl client.Client, lgr logr.Logger, setLogLevel int, serviceAccount string, targetNamespaceLabelSelector labels.Selector, scanningInterval int) ([]string, error) {

	var effectiveNamespaces []string

	// Loop through target namespaces
	for _, elem := range namespaces {
		namespaceList := &corev1.NamespaceList{}
		ns := strings.TrimSpace(elem)
		activeNamespace := ns
		activeNamespaceDisplayName := ns

		if ns == "*" {
			activeNamespace = ""
			activeNamespaceDisplayName = "*"
		}

		LogWithLevel("Querying for namespace/"+activeNamespaceDisplayName+" with sa/"+serviceAccount, 3, lgr, setLogLevel)
		// Get Namespace with the cached context
		namespaceListOptions := &client.ListOptions{Namespace: activeNamespace, LabelSelector: targetNamespaceLabelSelector}
		err := cl.List(context.Background(), namespaceList, namespaceListOptions)
		if err != nil {
			lgr.Error(err, "Failed to list namespace in cluster!")
			lgr.Info("Running reconciler again in " + strconv.Itoa(scanningInterval) + "s")
			time.Sleep(time.Second * time.Duration(scanningInterval))
			return []string{}, err
		}
		// Loop through NamespaceList, create the effectiveNamespaces slice
		for _, el := range namespaceList.Items {
			if !defaults.ContainsString(effectiveNamespaces, el.Name) {
				LogWithLevel("Adding ns/"+el.Name+" to scope", 3, lgr, setLogLevel)
				effectiveNamespaces = append(effectiveNamespaces, el.Name)
			}
		}
	}

	return effectiveNamespaces, nil
}

// createUniqueCertificateChecksum takes a seedString and a certificate byte stream and creates a unique SHA1 hash to track
func createUniqueCertificateChecksum(seedString string, cert *x509.Certificate) string {
	// Hash the Certificate and add it to the string slice
	h := sha1.New()
	byteStream := []byte(seedString)
	// Append the raw certificate data
	h.Write(append(byteStream, cert.Raw...))

	return hex.EncodeToString(h.Sum(nil))
}

// zeroing allows sensitive data to be removed from memory as soon as the use of them is complete via defered function calls
func zeroing(s []byte) {
	for i := 0; i < len(s); i++ {
		s[i] = 0
	}
}

// DaysOutToTimeOut converts an int slice of the number of days out to trigger an expiration alert on into a []configv1.TimeSlice time.Time array of computed date values to compare against certificate expiration dates with time.After
func DaysOutToTimeOut(targetDaysOut []int) []configv1.TimeSlice {
	// Set Active DaysOut and time.Time formatted future dates
	daysOut := targetDaysOut
	if len(targetDaysOut) == 0 {
		daysOut = defaults.DaysOut
	}

	timeNow := metav1.Now()
	timeOut := []configv1.TimeSlice{}

	for _, tR := range daysOut {
		futureTime := time.Hour * 24 * time.Duration(tR)
		tSlice := configv1.TimeSlice{Time: metav1.NewTime(timeNow.Add(futureTime)), DaysOut: tR}
		timeOut = append(timeOut, tSlice)
	}
	return timeOut
}

// GetServiceAccount returns a single ServiceAccount by name in a given Namespace
func GetServiceAccount(serviceAccount string, namespace string, clnt client.Client) (*corev1.ServiceAccount, error) {
	targetServiceAccount := &corev1.ServiceAccount{}
	err := clnt.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      serviceAccount,
	}, targetServiceAccount)

	if err != nil {
		return targetServiceAccount, err
	}
	return targetServiceAccount, nil
}

// GetSecret returns a single Secret by name in a given Namespace
func GetSecret(name string, namespace string, clnt client.Client) (*corev1.Secret, error) {
	targetSecret := &corev1.Secret{}
	err := clnt.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, targetSecret)

	if err != nil {
		return targetSecret, err
	}
	return targetSecret, nil
}

// GetConfigMap returns a single ConfigMap by name in a given Namespace
func GetConfigMap(name string, namespace string, clnt client.Client) (*corev1.ConfigMap, error) {
	targetConfigMap := &corev1.ConfigMap{}
	err := clnt.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, targetConfigMap)

	if err != nil {
		return targetConfigMap, err
	}
	return targetConfigMap, nil
}

// SetupNewClient takes a cached client, serviceAccount name, the namespace it is in, the cluster endpoint and API Path and creates a new Kubernetes client to act against the API on
func SetupNewClient(lgr logr.Logger, setLogLevel int, r client.Client, serviceAccount string, serviceAccountNamespace string, clusterEndpoint string, apiPath string) (client.Client, error) {
	// Get ServiceAccount
	LogWithLevel("Using ServiceAccount: "+serviceAccount, 2, lgr, setLogLevel)
	targetServiceAccount, err := GetServiceAccount(serviceAccount, serviceAccountNamespace, r)
	if err != nil {
		return *new(client.Client), err
	}
	var serviceAccountSecretName string
	targetServiceAccountSecret := &corev1.Secret{}

	// Find the right secret
	for _, em := range targetServiceAccount.Secrets {
		secret, _ := GetSecret(em.Name, serviceAccountNamespace, r)
		if secret.Type == "kubernetes.io/service-account-token" {
			// Get Secret
			serviceAccountSecretName = em.Name
			LogWithLevel("Using Secret: "+serviceAccountSecretName, 2, lgr, setLogLevel)
			targetServiceAccountSecret, _ = GetSecret(serviceAccountSecretName, serviceAccountNamespace, r)
		}
	}

	// We didn't find a Secret to work against the API and thus can't create a new client
	if serviceAccountSecretName == "" {
		return *new(client.Client), errors.New("new Kubernetes client creation failed!  No ServiceAccount Secret found")
	}

	// Set up new client config
	newConfig := &rest.Config{
		BearerToken: string(targetServiceAccountSecret.Data[corev1.ServiceAccountTokenKey]),
		Host:        clusterEndpoint,
		APIPath:     apiPath,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: targetServiceAccountSecret.Data[corev1.ServiceAccountRootCAKey],
		},
	}

	// Set up new Client
	cl, err := client.New(newConfig, client.Options{})
	if err != nil {
		return *new(client.Client), err
	}

	return cl, nil
}
