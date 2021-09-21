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

//==================================================================================================
// Logger Reports - Plain text based for SMTP too
//==================================================================================================

// LoggerKeystoreReport is a template constant that provides the overall structure of the KeystoreSentinel STDOUT Logger report
const LoggerKeystoreReport = `{{ .Divider }}
KeystoreSentinel Report: {{ .Namespace }}/{{ .Name }} ({{ .DateSent }})
{{ .Divider }}
  Cluster: {{ .ClusterAPIEndpoint }}
  Total Keystores Found: {{ .TotalKeystores }}
  Keystores at Risk: {{ .KeystoresAtRisk }}
  Total Certificates Found: {{ .TotalCerts }}
  Expiring Certificates Found: {{ .ExpiringCerts }}
{{ .Divider }}

{{ .Divider }}
{{ .Header }}
{{ .Divider }}
{{ .ReportLines }}{{ .Divider }}
{{ .Footer }}
{{ .Divider }}
`

// LoggerKeystoreReportLine is a template constant provides the template of each individual tabulated and delimited line in a KeystoreSentinel STDOUT Logger report that represents an identified expiring certificate in keystores
const LoggerKeystoreReportLine = `| {{ .APIVersion }} | {{ .Kind }} | {{ .Namespace }} | {{ .Name }} | {{ .Key }} | {{ .KeystoreAlias }} | {{ .CommonName }} | {{ .IsCA }} | {{ .CertificateAuthorityCommonName }} | {{ .ExpirationDate }} | {{ .TriggeredDaysOut }} |
`

// LoggerKeystoreReportHeader is a template constant provides the template of the tabulated and delimited table header (and footer, technically) columns in a KeystoreSentinel STDOUT Logger report
const LoggerKeystoreReportHeader = `| {{ .APIVersion }} | {{ .Kind }} | {{ .Namespace }} | {{ .Name }} | {{ .Key }} | {{ .KeystoreAlias }} | {{ .CommonName }} | {{ .IsCA }} | {{ .CertificateAuthorityCommonName }} | {{ .ExpirationDate }} | {{ .TriggeredDaysOut }} |`

// LoggerKeystoreReportStructure provides the overall data structure to the KeystoreSentinel STDOUT LoggerKeystoreReport template
type LoggerKeystoreReportStructure struct {
	Namespace          string
	Name               string
	DateSent           string
	ClusterAPIEndpoint string
	TotalKeystores     string
	KeystoresAtRisk    string
	TotalCerts         string
	ExpiringCerts      string
	ReportLines        string
	Header             string
	Footer             string
	Divider            string
}

// LoggerKeystoreReportHeaderStructure provides the data structure for the KeystoreSentinel LoggerKeystoreReport header (and footer, technically)
type LoggerKeystoreReportHeaderStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	KeystoreAlias                  string
	CommonName                     string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
}

// LoggerKeystoreReportLineStructure provides the data structure for the KeystoreSentinel LoggerKeystoreReportLine template
type LoggerKeystoreReportLineStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	KeystoreAlias                  string
	CommonName                     string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
}

//==================================================================================================
// SMTP HTML KeystoreSentinel Report Objects
//==================================================================================================

// HTMLSMTPKeystoreReportBody is the constant that provides an HTML template for KeystoreSentinel SMTP reports
const HTMLSMTPKeystoreReportBody = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html><head></head><body><div width="100%" style="margin:0!important;padding:10px 0!important;background-color:#ffffff">
<center style="width:100%;background-color:#ffffff">
<div>
<div style="text-align:left;">
<h1>KeystoreSentinel Report</h1>
<h3>{{ .DateSent }}</h3>
</div>
{{ .BodyDivider }}
<div style="text-align:left;">
<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="max-width:640px;">
<tbody>
<tr><td style="text-align:left;padding:6px;"><strong>Cluster:</strong></td><td>{{ .ClusterAPIEndpoint }}</td></tr>
<tr style="background:#EEE;"><td style="text-align:left;padding:6px;"><strong>KeystoreSentinel Namespace:</strong></td><td>{{ .Namespace }}</td></tr>
<tr><td style="text-align:left;padding:6px;"><strong>KeystoreSentinel Name:</strong></td><td>{{ .Name }}</td></tr>
<tr style="background:#EEE;"><td style="text-align:left;padding:6px;"><strong>Total Keystores Found:</strong></td><td>{{ .TotalKeystores }}</td></tr>
<tr><td style="text-align:left;padding:6px;"><strong>Keystores At Risk:</strong></td><td>{{ .KeystoresAtRisk }}</td></tr>
<tr style="background:#EEE;"><td style="text-align:left;padding:6px;"><strong>Total Certificates Found:</strong></td><td>{{ .TotalCerts }}</td></tr>
<tr><td style="text-align:left;padding:6px;"><strong>Expiring Certificates Found:</strong></td><td>{{ .ExpiringCerts }}</td></tr>
</tbody>
</table>
</div>
<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin:auto">
<thead style="font-weight:bold;">{{ .THead }}</thead>
<tbody>
{{ .TableRows }}
</tbody>
<tfoot style="font-weight:bold;">{{ .TFoot }}</tfoot>
</table>
</div>
</center></div></div>
</body>
</html>
`

// HTMLSMTPKeystoreReportBodyDivider is the constant template for a full-width HTML divider, the <hr /> element is all really
const HTMLSMTPKeystoreReportBodyDivider = `<div style="width:100%;"><hr /></div>`

// HTMLSMTPKeystoreReportBodyTableDivider is the constant template for an empty single-cell table row divider for a rich HTML KeystoreSentinel report
const HTMLSMTPKeystoreReportBodyTableDivider = `<tr><td style="text-align:left">&nbsp;</td></tr>`

// HTMLSMTPKeystoreReportLine is the constant template for a styled table row representing a certificate at risk for a rich HTML KeystoreSentinel report
const HTMLSMTPKeystoreReportLine = `<tr style="{{ .RowStyles }}">
<td style="{{ .CellStyles }}">{{ .APIVersion }}</td>
<td style="{{ .CellStyles }}">{{ .Kind }}</td>
<td style="{{ .CellStyles }}">{{ .Namespace }}</td>
<td style="{{ .CellStyles }}">{{ .Name }}</td>
<td style="{{ .CellStyles }}">{{ .Key }}</td>
<td style="{{ .CellStyles }}">{{ .KeystoreAlias }}</td>
<td style="{{ .CellStyles }}">{{ .CommonName }}</td>
<td style="{{ .CellStyles }}">{{ .IsCA }}</td>
<td style="{{ .CellStyles }}">{{ .CertificateAuthorityCommonName }}</td>
<td style="{{ .CellStyles }}">{{ .ExpirationDate }}</td>
<td style="{{ .CellStyles }}">{{ .TriggeredDaysOut }}</td>
</tr>`

// HTMLSMTPKeystoreReportHeader is the constant template for a styled table header (and tfoot, technically...) for a rich HTML KeystoreSentinel report
const HTMLSMTPKeystoreReportHeader = `<tr style="background:#EEE;">
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .APIVersion }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Kind }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Namespace }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Name }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Key }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .KeystoreAlias }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .CommonName }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .IsCA }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .CertificateAuthorityCommonName }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .ExpirationDate }}</td>
<td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .TriggeredDaysOut }}</td>
</tr>`

// HTMLKeystoreReportStructure provides the overall data structure to the HTMLSMTPKeystoreReport template
type HTMLKeystoreReportStructure struct {
	Namespace          string
	Name               string
	DateSent           string
	ClusterAPIEndpoint string
	TotalKeystores     string
	KeystoresAtRisk    string
	TotalCerts         string
	ExpiringCerts      string
	TableRows          string
	THead              string
	TFoot              string
	BodyDivider        string
}

// HTMLKeystoreReportLineStructure provides the data structure for the HTMLKeystoreReportLine template
type HTMLKeystoreReportLineStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	KeystoreAlias                  string
	CommonName                     string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
	RowStyles                      string
	CellStyles                     string
}

// HTMLKeystoreReportHeaderStructure provides the data structure for the HTMLKeystoreReportHeader template
type HTMLKeystoreReportHeaderStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	KeystoreAlias                  string
	CommonName                     string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
	RowStyles                      string
	CellStyles                     string
}

//==================================================================================================
// General SMTP Structs and Functions
//==================================================================================================
