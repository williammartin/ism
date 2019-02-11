package ui_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	. "github.com/pivotal-cf/ism/ui"
)

var _ = Describe("UI", func() {

	var testUI *UI

	BeforeEach(func() {
		testUI = &UI{
			Out: NewBuffer(),
			Err: NewBuffer(),
		}
	})

	Describe("DisplayText", func() {
		It("prints text with templated values to the out buffer", func() {
			testUI.DisplayText("This is a test for the {{.Struct}} struct", map[string]interface{}{"Struct": "UI"})
			Expect(testUI.Out).To(Say("This is a test for the UI struct\n"))
		})
	})
})
