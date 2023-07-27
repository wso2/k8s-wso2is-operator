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
		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{},
		}
		resourceReqs := MakeResourceRequirements(instance)

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
