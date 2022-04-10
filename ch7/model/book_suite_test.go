package model_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Book Suite")
}
