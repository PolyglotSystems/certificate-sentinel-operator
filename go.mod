module github.com/kenmoini/certificate-sentinel-operator

go 1.16

require (
	github.com/go-logr/logr v1.2.3
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.24.1
	github.com/pavel-v-chernykh/keystore-go/v4 v4.1.0
	github.com/xhit/go-simple-mail/v2 v2.10.0
	k8s.io/api v0.26.0
	k8s.io/apimachinery v0.26.0
	k8s.io/client-go v0.26.0
	sigs.k8s.io/controller-runtime v0.14.1
)
