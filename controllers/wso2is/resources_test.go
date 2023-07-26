package wso2is

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceRequirements", func() {
	It("should create the correct resource requirements", func() {
		resourceReqs := MakeResourceRequirements()

		// Verify Requests
		expectedRequests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1"),
			corev1.ResourceMemory: resource.MustParse("4096Mi"),
		}
		Expect(resourceReqs.Requests).To(Equal(expectedRequests))

		// Verify Limits
		expectedLimits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("2"),
			corev1.ResourceMemory: resource.MustParse("8000Mi"),
		}
		Expect(resourceReqs.Limits).To(Equal(expectedLimits))
	})
})
