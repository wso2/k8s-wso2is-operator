package wso2is

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PodNames", func() {
	It("should return the correct pod names", func() {
		pods := []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod3",
				},
			},
		}

		podNames := getPodNames(pods)

		// Verify the length of the podNames slice
		expectedLength := len(pods)
		Expect(len(podNames)).To(Equal(expectedLength))

		// Verify the content of each pod name
		expectedPodNames := []string{"pod1", "pod2", "pod3"}
		for i, expectedPodName := range expectedPodNames {
			podName := podNames[i]
			Expect(podName).To(Equal(expectedPodName), "Expected pod name %s at index %d, got %s", expectedPodName, i, podName)
		}
	})
})
