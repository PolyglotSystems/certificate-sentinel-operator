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
	"bytes"
	"context"
	"crypto/x509"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
	keystore "github.com/pavel-v-chernykh/keystore-go/v4"
)

// ReadKeyStoreFromBytes takes in a byte slice and password and decodes the
func ReadKeyStoreFromBytes(byteData []byte, password []byte) (keystore.KeyStore, error) {
	f := bytes.NewReader(byteData)

	keyStore := keystore.New(keystore.WithCaseExactAliases(), keystore.WithOrderedAliases())
	if err := keyStore.Load(f, password); err != nil {
		return keyStore, err
	}

	return keyStore, nil
}

// ProcessKeystoreIntoCertificates takes a JKS object and turns it into a list of decoded certificates
func ProcessKeystoreIntoCertificates(keystoreObj keystore.KeyStore) (map[string][]x509.Certificate, error) {
	certificateMap := make(map[string][]x509.Certificate)

	for _, iV := range keystoreObj.Aliases() {
		// Check if the entry is a Certificate
		if keystoreObj.IsTrustedCertificateEntry(iV) {
			// Pull Certificate bytes from the keystore
			cert, err := keystoreObj.GetTrustedCertificateEntry(iV)
			if err != nil {
				return certificateMap, err
			}
			// Make sure this is an X.509 type certificate
			if cert.Certificate.Type == "X.509" {
				// Decode certificate bytes into proper x509.Certificate object
				certsDecode, err := x509.ParseCertificate(cert.Certificate.Content)
				if err != nil {
					return certificateMap, err
				}
				// Add to certificate list
				certificateMap[iV] = append(certificateMap[iV], *certsDecode)
			}
		}
	}
	return certificateMap, nil
}

func getPasswordBytesFromSpecTarget(keystorePasswordDef configv1.KeystorePassword, namespace string, clnt client.Client) ([]byte, error) {
	passwordBytes := []byte("changeit")
	defer zeroing(passwordBytes)

	switch keystorePasswordDef.Type {
	case "secret":
		scrt, err := GetSecret(keystorePasswordDef.Secret.Name, namespace, clnt)
		if err != nil {
			return []byte{}, err
		}
		passwordBytes = scrt.Data[keystorePasswordDef.Secret.Key]
	case "labels":
		labelSelector, err := SetupSingleLabelSelector(keystorePasswordDef.Labels.LabelSelectors)
		if err != nil {
			return []byte{}, err
		}

		// Build List Options
		targetListOptions := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}

		// Get secrets matching the label
		secretList := &corev1.SecretList{}
		err = clnt.List(context.Background(), secretList, targetListOptions)
		if err != nil {
			return []byte{}, err
		}

		// Loop through the list of secrets, find the matching key
		for _, sV := range secretList.Items {
			if sV.Data[keystorePasswordDef.Labels.Key] != nil {
				passwordBytes = sV.Data[keystorePasswordDef.Labels.Key]
			}
		}

	case "plaintext":
		passwordBytes = []byte(keystorePasswordDef.Plaintext)
	}

	return passwordBytes, nil
}
