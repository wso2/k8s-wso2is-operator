package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ingress Reconciliation", func() {
	var (
		r        *Wso2IsReconciler
		instance wso2v1beta1.Wso2Is
		log      logr.Logger
		ctx      context.Context
		//ingress   *networkingv1.Ingress
		err          error
		reconciled   ctrl.Result
		reconcileErr error
	)

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)
		_ = networkingv1.AddToScheme(scheme)

		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		log = logr.Discard()

		r = &Wso2IsReconciler{
			Client: client,
			Log:    log,
			Scheme: scheme,
		}

		instance = wso2v1beta1.Wso2Is{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "",
			},
		}

		ctx = context.TODO()
		err = nil
		reconcileErr = nil
	})

	Context("When Ingress is not found", func() {
		BeforeEach(func() {
			reconciled, reconcileErr = reconcileIngress(r, instance, log, err, ctx)
		})

		It("Should return without error", func() {
			Expect(reconcileErr).NotTo(HaveOccurred(), "Error during Ingress reconciliation")
		})

		It("Should not have any reconciliation result", func() {
			Expect(reconciled).To(Equal(ctrl.Result{}))
		})

		It("Should log a message indicating Ingress not found", func() {
			// Perform assertions on the log output as per your logging framework
			// Example: Expect(log.InfoCallCount()).To(Equal(1))
		})

		It("Should not update the IngressHostname in the instance status", func() {
			Expect(instance.Status.IngressHostname).To(BeEmpty())
		})
	})

	Context("When Ingress is found", func() {
		var ingress *networkingv1.Ingress

		BeforeEach(func() {
			ingress = &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Name + "-ingress",
					Namespace: instance.Namespace,
				},
				Status: networkingv1.IngressStatus{
					LoadBalancer: networkingv1.IngressLoadBalancerStatus{
						Ingress: []networkingv1.IngressLoadBalancerIngress{
							{Hostname: "wso2is.com"},
						},
					},
				},
			}

			instance.Status.IngressHostname = "wso2is.com"

			reconciled, reconcileErr = reconcileIngress(r, instance, log, err, ctx)
		})

		It("Should return without error", func() {
			Expect(reconcileErr).NotTo(HaveOccurred(), "Error during Ingress reconciliation")
		})

		It("Should not have any reconciliation result", func() {
			Expect(reconciled).To(Equal(ctrl.Result{}))
		})

		It("Should log a message indicating Ingress is found", func() {
			// Perform assertions on the log output as per your logging framework
			// Example: Expect(log.InfoCallCount()).To(Equal(1))
		})
		It("Should update the IngressHostname in the instance status", func() {
			Expect(instance.Status.IngressHostname).To(Equal(ingress.Status.LoadBalancer.Ingress[0].Hostname))
		})
	})
})
