package wso2is

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LabelsForWso2IS", func() {
	It("should create the correct labels for WSO2 IS", func() {
		depname := "wso2is"
		version := "v1.0"

		labels := labelsForWso2IS(depname, version)

		// Verify the length of the labels map
		expectedLength := 6
		Expect(len(labels)).To(Equal(expectedLength))

		// Verify the content of each label
		expectedLabels := map[string]string{
			"app":        "wso2is",
			"deployment": "wso2is",
			"instance":   "wso2is",
			"monitoring": "jmx",
			"pod":        "wso2is",
			"version":    "v1.0",
		}
		for key, expectedValue := range expectedLabels {
			value, ok := labels[key]
			Expect(ok).To(BeTrue(), "Label %s is missing", key)
			Expect(value).To(Equal(expectedValue), "Expected label %s value %s, got %s", key, expectedValue, value)
		}
	})
})
