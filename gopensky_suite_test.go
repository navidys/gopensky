package gopensky_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGopensky(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gopensky Suite")
}
