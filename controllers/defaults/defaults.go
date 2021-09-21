// Package defaults contains the default values for various configurable options
package defaults

import (
	"strconv"
	"strings"
)

const ()

var (
	// ScanningInterval is the number of seconds to wait before the controller starts again
	ScanningInterval = 60
	// LogLevel is the log level - can be 1-4, with 1 being least verbose and 4 being most
	LogLevel = 2
	// DaysOut is the default number of days out to gate certificate expiration at
	DaysOut = []int{30, 60, 90}
	// ReportInterval is how frequently a report should be submitted for triggered targeted alerts
	ReportInterval = "daily"
	// SMTPAuthUseSSL is a boolean for if the Golang SMTP Client will use TLS against the server
	SMTPAuthUseSSL = true
	// SMTPAuthUseSTARTTLS is a boolean for if the Golang SMTP Client will use STARTTLS against the server
	SMTPAuthUseSTARTTLS = true
	// SMTPMessageSubject is the default subject sent with emailed messages
	SMTPMessageSubject = "Certificate Sentinel Operator - Report"
)

// SetDefaultInt64 will return either the default int64 or an overridden value
func SetDefaultInt64(defaultVal int64, overrideVal int64) int64 {
	if overrideVal == 0 {
		return defaultVal
	}
	return overrideVal
}

// SetDefaultInt32 will return either the default int32 or an overridden value
func SetDefaultInt32(defaultVal int32, overrideVal int32) int32 {
	iString := strings.TrimSpace(I32ToString(overrideVal))
	if overrideVal == 0 {
		return defaultVal
	}
	if len(iString) > 0 {
		return overrideVal
	}
	return defaultVal
}

// SetDefaultInt will return either the default int or an overridden value
func SetDefaultInt(defaultVal int, overrideVal int) int {
	if overrideVal == 0 {
		return defaultVal
	}
	if len(strings.TrimSpace(strconv.Itoa(overrideVal))) > 0 {
		return overrideVal
	}
	return defaultVal
}

// SetDefaultString will return either the default string or an overridden value
func SetDefaultString(defaultVal string, overrideVal string) string {
	if len(strings.TrimSpace(overrideVal)) > 0 {
		return overrideVal
	}
	return defaultVal
}

// SetDefaultBool will return either the default bool or an overridden value
func SetDefaultBool(defaultVal *bool, overrideVal *bool) *bool {
	if overrideVal != nil {
		return overrideVal
	}
	return defaultVal
}

// I32ToString will convert an int32 to a string for length comparison
func I32ToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// ContainsString checks if a string is present in a slice
func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
