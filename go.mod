module github.com/kenmoini/certificate-sentinel-operator

go 1.16

require (
	github.com/go-logr/logr v1.2.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pavel-v-chernykh/keystore-go/v4 v4.1.0
	github.com/xhit/go-simple-mail/v2 v2.10.0
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	sigs.k8s.io/controller-runtime v0.11.2
)
