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
	"fmt"

	"github.com/go-logr/logr"
	configv1 "github.com/kenmoini/certificate-sentinel-operator/apis/config/v1"
	defaults "github.com/kenmoini/certificate-sentinel-operator/controllers/defaults"
	helpers "github.com/kenmoini/certificate-sentinel-operator/controllers/helpers"
	corev1 "k8s.io/api/core/v1"

	"strconv"
	"strings"
	"text/template"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// processKeystoreReport processes reports for the KeystoreSentinel CRD
func processKeystoreReport(keystoreSentinel configv1.KeystoreSentinel, lggr logr.Logger, clnt client.Client) int64 {
	// Set up variables
	currentUnixTime := time.Now().Unix()
	lastReportSentTime := defaults.SetDefaultInt64(currentUnixTime, keystoreSentinel.Status.LastReportSent)
	effectiveReportInterval := defaults.SetDefaultString(defaults.ReportInterval, keystoreSentinel.Spec.Alert.AlertConfiguration.ReportInterval)

	// Add the times together
	addedTime := (lastReportSentTime + intervalToSeconds(effectiveReportInterval))

	// Next expected time to send is before the current time, overdue to send
	if (addedTime < time.Now().Unix()) || (lastReportSentTime == currentUnixTime) {
		LogWithLevel("Dispatching report for "+keystoreSentinel.Spec.Alert.AlertName, 2, lggr, setLogLevelK)

		// Send out alert based on alert type
		switch keystoreSentinel.Spec.Alert.AlertType {
		case "smtp":
			createKeystoreSMTPReport(keystoreSentinel, lggr, clnt)
		case "logger":
			loggr := createKeystoreLoggerReport(keystoreSentinel, lggr)
			lggr.Info(loggr)
		default:
			loggr := createKeystoreLoggerReport(keystoreSentinel, lggr)
			lggr.Info(loggr)
		}

		return time.Now().Unix()
	}
	return lastReportSentTime
}

// createKeystoreLoggerReport loops through KeystoreSentinel.Status and creates a stdout report
func createKeystoreLoggerReport(keystoreSentinel configv1.KeystoreSentinel, lggr logr.Logger) string {
	return createKeystoreTextTableReport(keystoreSentinel, lggr)
}

// createKeystoreTextTableReport creates a Text-based table of the report, used in logger reports and text-based SMTP reports
func createKeystoreTextTableReport(keystoreSentinel configv1.KeystoreSentinel, lggr logr.Logger) string {
	currentConfig, _ := config.GetConfig()
	clusterEndpoint := currentConfig.Host
	apiPath := currentConfig.APIPath
	expiredCertificateCount := 0

	// Set up init vars
	var reportLines string
	APIVersionLongest := "APIVersion"
	KindLongest := "Kind"
	NamespaceLongest := "Namespace"
	NameLongest := "Name"
	DataKeyLongest := "Data Key"
	KeystoreAliasLongest := "Keystore Alias"
	CertCNLongest := "Certificate CN"
	IsCALongest := "Is CA"
	CACNLongest := "Signing CA CN"
	ExpirationDateLongest := "Expiration Date"
	TriggeredDaysOutLongest := "Triggered Days Out"

	// Loop through the .status.DiscoveredCertificates
	for _, certInfo := range keystoreSentinel.Status.DiscoveredKeystoreCertificates {
		// If this is an expired certificate
		if len(certInfo.TriggeredDaysOut) > 0 {
			expiredCertificateCount++
			// Set up Logger Lines for length
			APIVersionLongest = helpers.ReturnLonger(APIVersionLongest, certInfo.APIVersion)
			KindLongest = helpers.ReturnLonger(KindLongest, certInfo.Kind)
			NamespaceLongest = helpers.ReturnLonger(NamespaceLongest, certInfo.Namespace)
			NameLongest = helpers.ReturnLonger(NameLongest, certInfo.Name)
			DataKeyLongest = helpers.ReturnLonger(DataKeyLongest, certInfo.DataKey)
			KeystoreAliasLongest = helpers.ReturnLonger(KeystoreAliasLongest, certInfo.KeystoreAlias)
			CertCNLongest = helpers.ReturnLonger(CertCNLongest, certInfo.CommonName)
			IsCALongest = helpers.ReturnLonger(IsCALongest, strconv.FormatBool(certInfo.IsCertificateAuthority))
			CACNLongest = helpers.ReturnLonger(CACNLongest, certInfo.CertificateAuthorityCommonName)
			ExpirationDateLongest = helpers.ReturnLonger(ExpirationDateLongest, certInfo.Expiration)
			TriggeredDaysOutLongest = helpers.ReturnLonger(TriggeredDaysOutLongest, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(certInfo.TriggeredDaysOut)), ", "), "[]"))
		}
	}

	// Set the longest length value
	APIVersionLength := len(APIVersionLongest)
	KindLength := len(KindLongest)
	NamespaceLength := len(NamespaceLongest)
	NameLength := len(NameLongest)
	DataKeyLength := len(DataKeyLongest)
	KeystoreAliasLength := len(KeystoreAliasLongest)
	CertCNLength := len(CertCNLongest)
	IsCALength := len(IsCALongest)
	CACNLength := len(CACNLongest)
	ExpirationDateLength := len(ExpirationDateLongest)
	TriggeredDaysOutLength := len(TriggeredDaysOutLongest)
	TotalLineLength := (APIVersionLength + KindLength + NamespaceLength + NameLength + DataKeyLength + KeystoreAliasLength + CertCNLength + IsCALength + CACNLength + ExpirationDateLength + TriggeredDaysOutLength + 34)
	LineBreak := helpers.StrPad("-", TotalLineLength, "-", "BOTH")

	// Loop through the .status.DiscoveredKeystoreCertificates
	for _, certInfo := range keystoreSentinel.Status.DiscoveredKeystoreCertificates {
		if len(certInfo.TriggeredDaysOut) > 0 {
			// Set up Logger Lines
			loggerReportLineStructure := LoggerKeystoreReportLineStructure{
				APIVersion:                     helpers.StrPad(certInfo.APIVersion, APIVersionLength, " ", "BOTH"),
				Kind:                           helpers.StrPad(certInfo.Kind, KindLength, " ", "BOTH"),
				Namespace:                      helpers.StrPad(certInfo.Namespace, NamespaceLength, " ", "BOTH"),
				Name:                           helpers.StrPad(certInfo.Name, NameLength, " ", "BOTH"),
				Key:                            helpers.StrPad(certInfo.DataKey, DataKeyLength, " ", "BOTH"),
				KeystoreAlias:                  helpers.StrPad(certInfo.KeystoreAlias, KeystoreAliasLength, " ", "BOTH"),
				CommonName:                     helpers.StrPad(certInfo.CommonName, CertCNLength, " ", "BOTH"),
				IsCA:                           helpers.StrPad(strconv.FormatBool(certInfo.IsCertificateAuthority), IsCALength, " ", "BOTH"),
				CertificateAuthorityCommonName: helpers.StrPad(certInfo.CertificateAuthorityCommonName, CACNLength, " ", "BOTH"),
				ExpirationDate:                 helpers.StrPad(certInfo.Expiration, ExpirationDateLength, " ", "BOTH"),
				TriggeredDaysOut:               helpers.StrPad(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(certInfo.TriggeredDaysOut)), ", "), "[]"), TriggeredDaysOutLength, " ", "BOTH"),
			}
			lineBuf := new(bytes.Buffer)
			loggerLineTemplate, err := template.New("loggerLine").Parse(LoggerKeystoreReportLine)
			if err != nil {
				lggr.Error(err, "Error parsing loggerLineTemplate template!")
			}
			err = loggerLineTemplate.Execute(lineBuf, loggerReportLineStructure)
			if err != nil {
				lggr.Error(err, "Error executing loggerLineTemplate template!")
			}
			// Append to total reportLines
			reportLines = (reportLines + lineBuf.String())
		}
	}

	// Setup Logger Headers
	loggerReportHeaderStructure := LoggerKeystoreReportHeaderStructure{
		APIVersion:                     helpers.StrPad("APIVersion", APIVersionLength, " ", "BOTH"),
		Kind:                           helpers.StrPad("Kind", KindLength, " ", "BOTH"),
		Namespace:                      helpers.StrPad("Namespace", NamespaceLength, " ", "BOTH"),
		Name:                           helpers.StrPad("Name", NameLength, " ", "BOTH"),
		Key:                            helpers.StrPad("Data Key", DataKeyLength, " ", "BOTH"),
		KeystoreAlias:                  helpers.StrPad("Keystore Alias", KeystoreAliasLength, " ", "BOTH"),
		CommonName:                     helpers.StrPad("Certificate CN", CertCNLength, " ", "BOTH"),
		IsCA:                           helpers.StrPad("Is CA", IsCALength, " ", "BOTH"),
		CertificateAuthorityCommonName: helpers.StrPad("Signing CA CN", CACNLength, " ", "BOTH"),
		ExpirationDate:                 helpers.StrPad("Expiration Date", ExpirationDateLength, " ", "BOTH"),
		TriggeredDaysOut:               helpers.StrPad("Triggered Days Out", TriggeredDaysOutLength, " ", "BOTH"),
	}
	headerBuf := new(bytes.Buffer)
	loggerHeaderTemplate, err := template.New("loggerHeader").Parse(LoggerKeystoreReportHeader)
	if err != nil {
		lggr.Error(err, "Error parsing loggerHeaderTemplate template!", 1, lggr)
	}
	err = loggerHeaderTemplate.Execute(headerBuf, loggerReportHeaderStructure)
	if err != nil {
		lggr.Error(err, "Error executing loggerHeaderTemplate template!")
	}

	// Set up Logger Report
	loggerReportStructure := LoggerKeystoreReportStructure{
		Namespace:          keystoreSentinel.Namespace,
		Name:               keystoreSentinel.Name,
		DateSent:           time.Now().UTC().String(),
		ClusterAPIEndpoint: clusterEndpoint + apiPath,
		TotalKeystores:     strconv.Itoa(keystoreSentinel.Status.TotalKeystoresFound),
		KeystoresAtRisk:    strconv.Itoa(keystoreSentinel.Status.KeystoresAtRisk),
		TotalCerts:         strconv.Itoa(len(keystoreSentinel.Status.DiscoveredKeystoreCertificates)),
		ExpiringCerts:      strconv.Itoa(keystoreSentinel.Status.ExpiringCertificates),
		ReportLines:        reportLines,
		Footer:             headerBuf.String(),
		Header:             headerBuf.String(),
		Divider:            LineBreak,
	}
	// Build template
	reportBuf := new(bytes.Buffer)
	loggerReportTemplate, err := template.New("loggerReport").Parse(LoggerKeystoreReport)
	if err != nil {
		lggr.Error(err, "Error parsing loggerReportTemplate template!")
	}
	// Execute our data against the template
	err = loggerReportTemplate.Execute(reportBuf, loggerReportStructure)
	if err != nil {
		lggr.Error(err, "Error executing loggerReportTemplate template!")
	}

	return reportBuf.String()
}

// createKeystoreSMTPReport loops through CertificateSentinel.Status and sends an email report
func createKeystoreSMTPReport(keystoreSentinel configv1.KeystoreSentinel, lggr logr.Logger, clnt client.Client) string {

	textEmailReport := tableTextReportToBasicHTMLReport(createKeystoreTextTableReport(keystoreSentinel, lggr))
	htmlEmailReport := createKeystoreSMTPHTMLReport(keystoreSentinel, lggr)

	// Loop through the alerts
	alert := keystoreSentinel.Spec.Alert
	if alert.AlertType == "smtp" {
		// Set up basic SMTP vars
		var username string
		var password string
		var identity string
		var cramSecret string

		// Set defaults
		useSSL := defaults.SetDefaultBool(&defaults.SMTPAuthUseSSL, alert.AlertConfiguration.SMTPAuthUseSSL)
		useSTARTTLS := defaults.SetDefaultBool(&defaults.SMTPAuthUseSTARTTLS, alert.AlertConfiguration.SMTPAuthUseSTARTTLS)

		// Get SMTP Authentication Secret if the AuthType is not `none`
		smtpAuthSecret := &corev1.Secret{}
		if alert.AlertConfiguration.SMTPAuthType != "none" {
			smtpAuthSecret, _ = GetSecret(alert.AlertConfiguration.SMTPAuthSecretName, keystoreSentinel.Namespace, clnt)
		}

		// Assign SMTP Auth vars where needed
		if len(smtpAuthSecret.Data) > 0 {
			username = string(smtpAuthSecret.Data["username"])
			password = string(smtpAuthSecret.Data["password"])
			identity = string(smtpAuthSecret.Data["identity"])
			cramSecret = string(smtpAuthSecret.Data["cram"])
		}

		// Send the message
		helpers.SendSMTPMail(alert.AlertConfiguration.SMTPAuthType,
			username,
			password,
			identity,
			cramSecret,
			useSSL,
			useSTARTTLS,
			alert.AlertConfiguration.SMTPDestinationEmailAddresses,
			alert.AlertConfiguration.SMTPSenderEmailAddress,
			alert.AlertConfiguration.SMTPEndpoint,
			"Certificate Sentinel Operator - KeystoreSentinel Report",
			textEmailReport,
			htmlEmailReport)
	}

	return textEmailReport
}

// createKeystoreSMTPHTMLReport creates a rich HTML-based table of the report, used in SMTP reports
func createKeystoreSMTPHTMLReport(keystoreSentinel configv1.KeystoreSentinel, lggr logr.Logger) string {

	// Set up init vars
	var reportLines string
	expiredCertificateCount := 0

	currentConfig, _ := config.GetConfig()
	clusterEndpoint := currentConfig.Host
	apiPath := currentConfig.APIPath

	// Loop through the .status.DiscoveredCertificates
	for iCI, certInfo := range keystoreSentinel.Status.DiscoveredKeystoreCertificates {
		if len(certInfo.TriggeredDaysOut) > 0 {
			expiredCertificateCount++
			// Set up styles
			var rowStyles string
			cellStyles := "padding:6px;text-align:left;"
			if iCI%2 == 0 {
				rowStyles = "background:#FFF;"
			} else {
				rowStyles = "background:#FAFAFA;"
			}

			// Parse time to format
			inputTimelayout := "2006-01-02 15:04:05 -0700 MST"
			t, err := time.Parse(inputTimelayout, certInfo.Expiration)
			if err != nil {
				fmt.Printf("Error parsing time! %+v\n", err)
			}

			// Set up HTML Lines
			htmlSMTPReportLine := HTMLKeystoreReportLineStructure{
				APIVersion:                     certInfo.APIVersion,
				Kind:                           certInfo.Kind,
				Namespace:                      certInfo.Namespace,
				Name:                           certInfo.Name,
				Key:                            certInfo.DataKey,
				KeystoreAlias:                  certInfo.KeystoreAlias,
				CommonName:                     certInfo.CommonName,
				IsCA:                           strconv.FormatBool(certInfo.IsCertificateAuthority),
				CertificateAuthorityCommonName: certInfo.CertificateAuthorityCommonName,
				ExpirationDate:                 string(t.Format(time.RFC822Z)),
				TriggeredDaysOut:               strings.Trim(strings.Join(strings.Fields(fmt.Sprint(certInfo.TriggeredDaysOut)), ", "), "[]"),
				RowStyles:                      rowStyles,
				CellStyles:                     cellStyles,
			}
			lineBuf := new(bytes.Buffer)
			htmlLineTemplate, err := template.New("tableLine").Parse(HTMLSMTPKeystoreReportLine)
			if err != nil {
				lggr.Error(err, "Error parsing htmlSMTPReportLine template!")
			}
			err = htmlLineTemplate.Execute(lineBuf, htmlSMTPReportLine)
			if err != nil {
				lggr.Error(err, "Error executing htmlSMTPReportLine template!")
			}
			// Append to total reportLines
			reportLines = (reportLines + lineBuf.String())
		}
	}

	// Set up styles for headers
	rowStyles := "background-color:#EEE;"
	cellStyles := "padding:6px;border-top: 1px solid #999;border-bottom: 1px solid #999;text-align:left;"

	// Setup HTML Headers
	htmlReportHeaderStructure := HTMLKeystoreReportHeaderStructure{
		APIVersion:                     "APIVersion",
		Kind:                           "Kind",
		Namespace:                      "Namespace",
		Name:                           "Name",
		Key:                            "Data Key",
		KeystoreAlias:                  "Keystore Alias",
		CommonName:                     "Certificate CN",
		IsCA:                           "Is CA",
		CertificateAuthorityCommonName: "Signing CA CN",
		ExpirationDate:                 "Expiration Date",
		TriggeredDaysOut:               "Triggered Days Out",
		RowStyles:                      rowStyles,
		CellStyles:                     cellStyles,
	}
	headerBuf := new(bytes.Buffer)
	tableHeaderTemplate, err := template.New("tableHeader").Parse(HTMLSMTPKeystoreReportHeader)
	if err != nil {
		lggr.Error(err, "Error parsing tableHeaderTemplate template!")
	}
	err = tableHeaderTemplate.Execute(headerBuf, htmlReportHeaderStructure)
	if err != nil {
		lggr.Error(err, "Error executing tableHeaderTemplate template!")
	}

	// Set up HTML Report
	htmlReportStructure := HTMLKeystoreReportStructure{
		Namespace:          keystoreSentinel.Namespace,
		Name:               keystoreSentinel.Name,
		DateSent:           string(time.Now().UTC().Format(time.RFC822Z)),
		ClusterAPIEndpoint: clusterEndpoint + apiPath,
		TotalKeystores:     strconv.Itoa(keystoreSentinel.Status.TotalKeystoresFound),
		KeystoresAtRisk:    strconv.Itoa(keystoreSentinel.Status.KeystoresAtRisk),
		TotalCerts:         strconv.Itoa(len(keystoreSentinel.Status.DiscoveredKeystoreCertificates)),
		ExpiringCerts:      strconv.Itoa(keystoreSentinel.Status.ExpiringCertificates),
		TableRows:          reportLines,
		THead:              headerBuf.String(),
		TFoot:              headerBuf.String(),
		BodyDivider:        HTMLSMTPReportBodyDivider,
	}
	reportBuf := new(bytes.Buffer)
	htmlReportTemplate, err := template.New("HTMLReport").Parse(HTMLSMTPKeystoreReportBody)
	if err != nil {
		lggr.Error(err, "Error parsing htmlReportTemplate template!")
	}
	err = htmlReportTemplate.Execute(reportBuf, htmlReportStructure)
	if err != nil {
		lggr.Error(err, "Error executing htmlReportTemplate template!")
	}

	return reportBuf.String()
}
