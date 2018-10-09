package k8sutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8SUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8S Util Suite")
}
