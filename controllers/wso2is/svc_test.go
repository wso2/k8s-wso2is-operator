package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wso2IsReconciler", func() {
	It("should define a service", func() {
		// Prepare scheme for the fake client
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)

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
		service := r.defineService(m)

		Expect(service).ToNot(BeNil())
		Expect(service.Name).To(Equal(m.Name + "-service"))
		Expect(service.Namespace).To(Equal(m.Namespace))
	})

	Context("Service reconciliation", func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)

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
				Name:      "TestInstance",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "",
			},
		}

		ctx := context.TODO()

		It("Should reconcile a Service", func() {
			_, err := reconcileSvc(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during Service reconciliation")

			service := &corev1.Service{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-service", Namespace: instance.Namespace}, service)
			Expect(err).NotTo(HaveOccurred(), "Service not got")
		})
	})

})
