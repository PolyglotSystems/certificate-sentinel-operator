// Package helpers contains the supporting functions for the operator at large that are shared
package helpers

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/go-logr/logr"
	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
)

/*=====================================================================================
| x509 Certificate Helper Functions
=====================================================================================*/

// DecodeCertificateBytes decodes the byte slice from a Secret data item into a PEM and then into an x509 DER object
func DecodeCertificateBytes(s []byte, lggr logr.Logger) ([]*x509.Certificate, error) {
	block, rest := pem.Decode(s)

	// Check to see if this can be decoded into a PEM block
	if block == nil || block.Type != "CERTIFICATE" {
		lggr.Info("Failed to decode PEM block containing a Certificate: " + string(rest))
	}

	// Parse the Certificate
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		lggr.Error(err, "Failed to decode certificate!")
		return nil, err
	}
	return certs, nil
}

// ParseCertificateIntoObjects takes a slice of certificates and all the other supporting information to create a CertificateInformation return to add to .status/alert report
func ParseCertificateIntoObjects(cert *x509.Certificate, timeOut []configv1.TimeSlice, namespace string, name string, dataKey string, kind string, apiVersion string) (discovered []configv1.CertificateInformation, messages []string) {
	discoveredL := []configv1.CertificateInformation{}
	var messagesL []string

	var triggeredDaysOut []int
	expirationDate := cert.NotAfter

	// Loop through the timeOut slice, add triggered values to the certInfo config for priority ranking
	for _, d := range timeOut {
		utcTime := d.Time.UTC()
		if utcTime.After(expirationDate) {
			messagesL = append(messagesL, "Certificate will expire in less than "+fmt.Sprint(d.DaysOut)+" days! Date: "+expirationDate.String())
			triggeredDaysOut = append(triggeredDaysOut, d.DaysOut)
		}
	}
	// Create CertificateInformation object
	certInfo := configv1.CertificateInformation{Namespace: namespace, Name: name, DataKey: dataKey, Kind: kind, APIVersion: apiVersion, Expiration: expirationDate.String(), CommonName: cert.Subject.CommonName, CertificateAuthorityCommonName: cert.Issuer.CommonName, IsCertificateAuthority: cert.IsCA, TriggeredDaysOut: triggeredDaysOut}

	// This certificate is not expired
	discoveredL = append(discoveredL, certInfo)
	// If this certificate triggered any alerted daysOut, it is at risk of expiring

	return discoveredL, messagesL
}

// ParseKeystoreCertificateIntoObjects takes a slice of certificates and all the other supporting information to create a KeystoreInformation return to add to .status/alert report
func ParseKeystoreCertificateIntoObjects(cert *x509.Certificate, timeOut []configv1.TimeSlice, namespace string, name string, dataKey string, kind string, apiVersion string, keystoreAlias string) (discovered []configv1.KeystoreInformation, messages []string) {
	discoveredL := []configv1.KeystoreInformation{}
	var messagesL []string

	var triggeredDaysOut []int
	expirationDate := cert.NotAfter

	// Loop through the timeOut slice, add triggered values to the certInfo config for priority ranking
	for _, d := range timeOut {
		utcTime := d.Time.UTC()
		if utcTime.After(expirationDate) {
			messagesL = append(messagesL, "Certificate will expire in less than "+fmt.Sprint(d.DaysOut)+" days! Date: "+expirationDate.String())
			triggeredDaysOut = append(triggeredDaysOut, d.DaysOut)
		}
	}
	// Create KeystoreInformation object
	certInfo := configv1.KeystoreInformation{Namespace: namespace, Name: name, DataKey: dataKey, Kind: kind, APIVersion: apiVersion, KeystoreAlias: keystoreAlias, Expiration: expirationDate.String(), CommonName: cert.Subject.CommonName, CertificateAuthorityCommonName: cert.Issuer.CommonName, IsCertificateAuthority: cert.IsCA, TriggeredDaysOut: triggeredDaysOut}

	// This certificate is not expired
	discoveredL = append(discoveredL, certInfo)
	// If this certificate triggered any alerted daysOut, it is at risk of expiring

	return discoveredL, messagesL
}

// ParseCertificatesIntoLists [DEPRECIATED] takes a slice of certificates and all the other supporting information to create a CertificateInformation return to add to .status/alert report
func ParseCertificatesIntoLists(certs []*x509.Certificate, timeOut []configv1.TimeSlice, namespace string, name string, dataKey string, kind string, apiVersion string) (discovered []configv1.CertificateInformation, expired []configv1.CertificateInformation, messages []string) {
	discoveredL := []configv1.CertificateInformation{}
	expiredL := []configv1.CertificateInformation{}
	var messagesL []string
	var RHashL []string

	// Loop over parsed certificates
	for _, cert := range certs {
		var triggeredDaysOut []int
		expirationDate := cert.NotAfter

		// Loop through the timeOut slice, add triggered values to the certInfo config for priority ranking
		for _, d := range timeOut {
			utcTime := d.Time.UTC()
			if utcTime.After(expirationDate) {
				messagesL = append(messagesL, "Certificate will expire in less than "+fmt.Sprint(d.DaysOut)+" days! Date: "+expirationDate.String())
				triggeredDaysOut = append(triggeredDaysOut, d.DaysOut)
			}
		}
		// Create CertificateInformation object
		certInfo := configv1.CertificateInformation{Namespace: namespace, Name: name, DataKey: dataKey, Kind: kind, APIVersion: apiVersion, Expiration: expirationDate.String(), CommonName: cert.Subject.CommonName, CertificateAuthorityCommonName: cert.Issuer.CommonName, IsCertificateAuthority: cert.IsCA, TriggeredDaysOut: triggeredDaysOut}

		// Hash the Certificate and add it to the string slice
		h := sha1.New()
		h.Write(cert.Raw)
		//shaStr := fmt.Sprintf("%x", h.Sum(nil))
		shaStr := hex.EncodeToString(h.Sum(nil))

		elFound := false
		for _, v := range RHashL {
			if v == shaStr {
				elFound = true
			}
		}

		// Check to see if this certificate has been added yet
		//if containsString(RHashL, shaStr) {
		if elFound {
			// Skipping Certificate
		} else {
			RHashL = append(RHashL, shaStr)

			// This certificate is not expired
			discoveredL = append(discoveredL, certInfo)
			// If this certificate triggered any alerted daysOut, it is at risk of expiring
			if len(triggeredDaysOut) != 0 {
				expiredL = append(expiredL, certInfo)
			}
		}
	}
	return discoveredL, expiredL, messagesL
}
