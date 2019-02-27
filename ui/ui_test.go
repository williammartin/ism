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

	Describe("DisplayTable", func() {
		It("prints a table with headers", func() {
			testUI.DisplayTable([][]string{
				{"header1", "header2", "header3"},
				{"data1", "mydata2", "data3"},
				{"data4", "data5", "data6"},
			})
			Expect(testUI.Out).To(Say("header1"))
			Expect(testUI.Out).To(Say("header2"))
			Expect(testUI.Out).To(Say("header3"))
			Expect(testUI.Out).To(Say(`data1\s+mydata2\s+data3`))
			Expect(testUI.Out).To(Say(`data4\s+data5\s+data6`))
		})
	})
})
