package gopensky_test

import (
	"github.com/navidys/gopensky"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("floatToString", func() {
		It("converts float64 to string", func() {
			tests := []struct {
				have  float64
				wants string
			}{
				{have: 3.14, wants: "3.140000"},
				{have: 1, wants: "1.000000"},
				{have: 2.1, wants: "2.100000"},
			}

			for _, test := range tests {
				Expect(gopensky.FloatToString(test.have)).To(Equal(test.wants))
			}
		})
	})
})
