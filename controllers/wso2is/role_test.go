package wso2is

import (
	"context"
	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Role", func() {
	It("should define the correct Role", func() {
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
		role := r.defineRole(m)

		Expect(role).NotTo(BeNil(), "Role is nil")
		Expect(role.Name).To(Equal(m.Name+"-role"), "Expected Role Name %s, got %s", m.Name+"-role", role.Name)
		Expect(role.Namespace).To(Equal(m.Namespace), "Expected Role Namespace %s, got %s", m.Namespace, role.Namespace)
	})

	Context("Role reconciliation", func() {
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)
		_ = rbacv1.AddToScheme(scheme)

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

		It("Should reconcile a Role", func() {
			// Create and reconcile the Role
			_, err := reconcileRole(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during Role reconciliation")

			role := &rbacv1.Role{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-role", Namespace: instance.Namespace}, role)
			Expect(err).NotTo(HaveOccurred(), "Role not got")
		})
	})

})
