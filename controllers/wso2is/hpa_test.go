package wso2is

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Hpa", func() {
	Context("Define HPA", func() {
		It("should define the correct HPA", func() {
			// Prepare scheme for the fake client
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = autoscalingv2.AddToScheme(scheme)

			// Create an instance of your reconciler with the fake client
			r := &Wso2IsReconciler{
				Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
				Scheme: scheme,
			}

			m := wso2v1beta1.Wso2Is{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "wso2is",
					Namespace: "wso2-iam-system",
				},
			}

			// Call the method under test
			hpa := r.defineHpa(m)

			Expect(hpa).NotTo(BeNil(), "HPA is nil")
			Expect(hpa.Name).To(Equal(m.Name+"-hpa"), "Expected HPA Name %s, got %s", m.Name+"-hpa", hpa.Name)
			Expect(hpa.Namespace).To(Equal(m.Namespace), "Expected HPA Namespace %s, got %s", m.Namespace, hpa.Namespace)
			Expect(*hpa.Spec.MinReplicas).To(BeEquivalentTo(1), "Expected MinReplicas %d, got %d", 1, *hpa.Spec.MinReplicas)
			Expect(hpa.Spec.MaxReplicas).To(BeEquivalentTo(3), "Expected MaxReplicas %d, got %d", 3, hpa.Spec.MaxReplicas)
			Expect(len(hpa.Spec.Metrics)).To(Equal(1), "Expected 1 metric, got %d", len(hpa.Spec.Metrics))
			Expect(hpa.Spec.Metrics[0].Type).To(Equal(autoscalingv2.ResourceMetricSourceType), "Expected metric type Resource, got %s", hpa.Spec.Metrics[0].Type)
			Expect(hpa.Spec.Metrics[0].Resource.Name).To(Equal(corev1.ResourceCPU), "Expected resource name CPU, got %s", hpa.Spec.Metrics[0].Resource.Name)
			Expect(hpa.Spec.Metrics[0].Resource.Target.Type).To(Equal(autoscalingv2.UtilizationMetricType), "Expected metric target type Utilization, got %s", hpa.Spec.Metrics[0].Resource.Target.Type)
			Expect(*hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).To(BeEquivalentTo(50), "Expected target utilization 50, got %d", *hpa.Spec.Metrics[0].Resource.Target.AverageUtilization)
		})
	})

	Context("Reconcile HPA", func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)
		_ = autoscalingv2.AddToScheme(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		logger := logr.Discard()

		// Prepare reconciler
		recon := &Wso2IsReconciler{
			Client: client,
			Log:    logger,
			Scheme: scheme,
		}

		instance := wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "",
			},
		}

		ctx := context.TODO()

		It("Should reconcile an HPA", func() {
			// Create and reconcile the HPA
			_, err := reconcileHpa(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during HPA reconciliation")

			hpa := &autoscalingv2.HorizontalPodAutoscaler{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-hpa", Namespace: instance.Namespace}, hpa)
			Expect(err).NotTo(HaveOccurred(), "HPA not found")
		})
	})
})
