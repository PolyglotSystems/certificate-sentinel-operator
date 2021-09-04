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

const LoggerReport = `{{ .Divider }}
CertificateSentinel Report: {{ .Namespace }}/{{ .Name }} ({{ .DateSent }})
{{ .Divider }}
  Cluster: {{ .ClusterAPIEndpoint }}
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

const LoggerReportLine = `| {{ .APIVersion }} | {{ .Kind }} | {{ .Namespace }} | {{ .Name }} | {{ .Key }} | {{ .CertificateAuthorityCommonName }} | {{ .IsCA }} | {{ .ExpirationDate }} | {{ .TriggeredDaysOut }} |
`

const LoggerReportHeader = `| {{ .APIVersion }} | {{ .Kind }} | {{ .Namespace }} | {{ .Name }} | {{ .Key }} | {{ .CertificateAuthorityCommonName }} | {{ .IsCA }} | {{ .ExpirationDate }} | {{ .TriggeredDaysOut }} |`

// loggerReportStructure provides the overall structure to the loggerReport template
type LoggerReportStructure struct {
	Namespace          string
	Name               string
	DateSent           string
	ClusterAPIEndpoint string
	TotalCerts         string
	ExpiringCerts      string
	ReportLines        string
	Header             string
	Footer             string
	Divider            string
}

// LoggerReportHeaderStructure provides the structure for the LoggerReport header
type LoggerReportHeaderStructure struct {
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

//==================================================================================================
// SMTP HTML Reports
//==================================================================================================

const TextSMTPReportDocument = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
<head></head>
<body style="font-family:monospace">
<pre>{{ .Content }}</pre>
</body>
</html>
`

const HTMLSMTPReportBody = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html><head></head><body><div width="100%" style="margin:0!important;padding:10px 0!important;background-color:#ffffff">
<center style="width:100%;background-color:#ffffff">
<div>
<div style="text-align:left;">
<h1>CertificateSentinel Report</h1>
<h3>{{ .DateSent }}</h3>
</div>
{{ .BodyDivider }}
<div style="text-align:left;">
<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="max-width:640px;">
<tbody>
<tr><td style="text-align:left;padding:6px;"><strong>Cluster:</strong></td><td>{{ .ClusterAPIEndpoint }}</td></tr>
<tr style="background:#EEE;"><td style="text-align:left;padding:6px;"><strong>CertificateSentinel Namespace:</strong></td><td>{{ .Namespace }}</td></tr>
<tr><td style="text-align:left;padding:6px;"><strong>CertificateSentinel Name:</strong></td><td>{{ .Name }}</td></tr>
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

const HTMLSMTPReportBodyDivider = `<div style="width:100%;"><hr /></div>`
const HTMLSMTPReportBodyTableDivider = `<tr><td style="text-align:left">&nbsp;</td></tr>`

const HTMLSMTPReportLine = `<tr style="{{ .RowStyles }}"><td style="{{ .CellStyles }}">{{ .APIVersion }}</td><td style="{{ .CellStyles }}">{{ .Kind }}</td><td style="{{ .CellStyles }}">{{ .Namespace }}</td><td style="{{ .CellStyles }}">{{ .Name }}</td><td style="{{ .CellStyles }}">{{ .Key }}</td><td style="{{ .CellStyles }}">{{ .CertificateAuthorityCommonName }}</td><td style="{{ .CellStyles }}">{{ .IsCA }}</td><td style="{{ .CellStyles }}">{{ .ExpirationDate }}</td><td style="{{ .CellStyles }}">{{ .TriggeredDaysOut }}</td></tr>`

const HTMLSMTPReportHeader = `<tr style="background:#EEE;"><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .APIVersion }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Kind }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Namespace }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Name }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .Key }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .CertificateAuthorityCommonName }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .IsCA }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .ExpirationDate }}</td><td style="{{ .CellStyles }}border-bottom:1px solid #999;border-top:1px solid #999;">{{ .TriggeredDaysOut }}</td></tr>`

// TextSMTPReportStructure is just a wrapper for the text-report in an HTML document
type TextSMTPReportStructure struct {
	Content string
}

// HTMLReportStructure provides the overall structure to the HTMLSMTPReport template
type HTMLReportStructure struct {
	Namespace          string
	Name               string
	DateSent           string
	ClusterAPIEndpoint string
	TotalCerts         string
	ExpiringCerts      string
	TableRows          string
	THead              string
	TFoot              string
	BodyDivider        string
}

// HTMLReportLineStructure provides the struct for the htmlReportLine template
type HTMLReportLineStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
	IsCA                           string
	CertificateAuthorityCommonName string
	ExpirationDate                 string
	TriggeredDaysOut               string
	RowStyles                      string
	CellStyles                     string
}

// HTMLReportHeaderStructure provides the struct for the htmlReportHeader template
type HTMLReportHeaderStructure struct {
	APIVersion                     string
	Kind                           string
	Namespace                      string
	Name                           string
	Key                            string
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
