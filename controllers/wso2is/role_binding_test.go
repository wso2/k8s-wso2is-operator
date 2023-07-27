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

var _ = Describe("RoleBinding", func() {
	It("should define the correct RoleBinding", func() {
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
		roleBinding := r.defineRoleBinding(m)

		Expect(roleBinding).NotTo(BeNil(), "RoleBinding is nil")
		Expect(roleBinding.Name).To(Equal(m.Name+"-rolebinding"), "Expected RoleBinding Name %s, got %s", m.Name+"-rolebinding", roleBinding.Name)
		Expect(roleBinding.Namespace).To(Equal(m.Namespace), "Expected RoleBinding Namespace %s, got %s", m.Namespace, roleBinding.Namespace)
		Expect(roleBinding.RoleRef.Kind).To(Equal("Role"), "Expected RoleBinding RoleRef Kind %s, got %s", "Role", roleBinding.RoleRef.Kind)
		Expect(roleBinding.RoleRef.Name).To(Equal(m.Name+"-role"), "Expected RoleBinding RoleRef Name %s, got %s", m.Name+"-role", roleBinding.RoleRef.Name)
		Expect(roleBinding.RoleRef.APIGroup).To(Equal("rbac.authorization.k8s.io"), "Expected RoleBinding RoleRef APIGroup %s, got %s", "rbac.authorization.k8s.io", roleBinding.RoleRef.APIGroup)
		Expect(roleBinding.Subjects).To(HaveLen(1), "Expected RoleBinding Subjects length %d, got %d", 1, len(roleBinding.Subjects))
		Expect(roleBinding.Subjects[0].Name).To(Equal(m.Name+"-serviceaccount"), "Expected RoleBinding Subject Name %s, got %s", m.Name+"-serviceaccount", roleBinding.Subjects[0].Name)
		Expect(roleBinding.Subjects[0].Kind).To(Equal("ServiceAccount"), "Expected RoleBinding Subject Kind %s, got %s", "ServiceAccount", roleBinding.Subjects[0].Kind)
	})

	Context("RoleBinding reconciliation", func() {
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

		It("Should reconcile a RoleBinding", func() {
			// Create and reconcile the RoleBinding
			_, err := reconcileRoleBinding(recon, instance, logger, nil, ctx)
			Expect(err).NotTo(HaveOccurred(), "Error during RoleBinding reconciliation")

			roleBinding := &rbacv1.RoleBinding{}
			err = recon.Client.Get(ctx, types.NamespacedName{Name: instance.Name + "-rolebinding", Namespace: instance.Namespace}, roleBinding)
			Expect(err).NotTo(HaveOccurred(), "RoleBinding not got")
		})
	})

})
