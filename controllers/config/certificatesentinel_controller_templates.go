/*
Copyright 2021 Ken Moini.

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

const LoggerReport = `-------------------------------------------------------------------------------------------------------------------------
CertificateSentinel Report: {{ .Namespace }}/{{ .Name }} ({{ .DateSent }})
-------------------------------------------------------------------------------------------------------------------------
  Cluster: {{ .ClusterAPIEndpoint }}
  Total Certificates Found: {{ .TotalCerts }}
  Expiring Certificates Found: {{ .ExpiringCerts }}
-------------------------------------------------------------------------------------------------------------------------

-------------------------------------------------------------------------------------------------------------------------
| apiVersion | kind | namespace | name | key | isCA | signingCertificateAuthorityCN | expirationDate | triggeredDaysOut |
-------------------------------------------------------------------------------------------------------------------------
{{ .ReportLines }}
-------------------------------------------------------------------------------------------------------------------------
| apiVersion | kind | namespace | name | key | isCA | signingCertificateAuthorityCN | expirationDate | triggeredDaysOut |
-------------------------------------------------------------------------------------------------------------------------
`

const LoggerReportLine = `| {{ .APIVersion }} | {{ .Kind }} | {{ .Namespace }} | {{ .Name }} | {{ .Key }} | {{ .CertificateAuthorityCommonName }} | {{ .IsCA }} | {{ .ExpirationDate }} | {{ .TriggeredDaysOut }} |
`

// loggerReportStructure provides the overall structure to the loggerReport template
type LoggerReportStructure struct {
	Namespace          string
	Name               string
	DateSent           string
	ClusterAPIEndpoint string
	TotalCerts         string
	ExpiringCerts      string
	ReportLines        string
}

// loggerReportLineStructure provides the struct for the loggerReportLine template
type LoggerReportLineStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
}