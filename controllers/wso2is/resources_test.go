package wso2is

import (
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceRequirements", func() {
	It("should create the correct resource requirements", func() {
		// Define your variables for CPU, memory, etc.
		cpuRequests := "1"
		memoryRequests := "4096Mi"
		cpuLimits := "2"
		memoryLimits := "8000Mi"

		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				Resources: wso2v1beta1.Resources{
					Requests: wso2v1beta1.ResourceRequests{
						Cpu:    cpuRequests,
						Memory: memoryRequests,
					},
					Limits: wso2v1beta1.ResourceLimits{
						Cpu:    cpuLimits,
						Memory: memoryLimits,
					},
				},
			},
		}
		resourceReqs := MakeResourceRequirements(instance)

		// Verify Requests
		expectedRequests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuRequests),
			corev1.ResourceMemory: resource.MustParse(memoryRequests),
		}
		Expect(resourceReqs.Requests).To(Equal(expectedRequests))

		// Verify Limits
		expectedLimits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimits),
			corev1.ResourceMemory: resource.MustParse(memoryLimits),
		}
		Expect(resourceReqs.Limits).To(Equal(expectedLimits))
	})
})
