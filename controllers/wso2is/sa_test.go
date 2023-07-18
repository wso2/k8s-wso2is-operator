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

var _ = Describe("ServiceAccount", func() {
	It("should define the correct ServiceAccount", func() {
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
		serviceAccount := r.defineServiceAccount(m)

		Expect(serviceAccount).NotTo(BeNil(), "ServiceAccount is nil")
		Expect(serviceAccount.Name).To(Equal(m.Name+"-serviceaccount"), "Expected ServiceAccount Name %s, got %s", m.Name+"-serviceaccount", serviceAccount.Name)
		Expect(serviceAccount.Namespace).To(Equal(m.Namespace), "Expected ServiceAccount Namespace %s, got %s", m.Namespace, serviceAccount.Namespace)
	})

	Context("ServiceAccount reconciliation", func() {
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
				Name:      "wso2is",
				Namespace: "wso2-iam-system",
			},
			Spec: wso2v1beta1.Wso2IsSpec{
				TomlConfigFile: "",
			},
		}

		ctx := context.TODO()

		It("Should reconcile a ServiceAccount", func() {
			// Create and reconcile the ServiceAccount
			_, err := reconcileSva(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during ServiceAccount reconciliation")

			serviceAccount := &corev1.ServiceAccount{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-serviceaccount", Namespace: instance.Namespace}, serviceAccount)
			Expect(err).NotTo(HaveOccurred(), "ServiceAccount not got")
		})
	})

})
