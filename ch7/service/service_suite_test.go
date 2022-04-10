package service_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Service Suite")
}

var _ = ginkgo.BeforeSuite(func() {

})
