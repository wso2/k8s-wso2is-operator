package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineRoleBinding(m wso2v1beta1.Wso2Is) *rbacv1.RoleBinding {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-rolebinding",
			Namespace: m.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      m.Name + "-serviceaccount",
				Namespace: m.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     m.Name + "-role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	ctrl.SetControllerReference(&m, roleBinding, r.Scheme)
	return roleBinding
}

func reconcileRoleBinding(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	roleBindingDefinition := r.defineRoleBinding(instance)
	roleBinding := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-rolebinding", Namespace: instance.Namespace}, roleBinding)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("RoleBinding resource " + instance.Name + "-rolebinding" + " not found. Creating or re-creating role binding")
			err = r.Create(ctx, roleBindingDefinition)
			if err != nil {
				log.Error(err, "Failed to create new RoleBinding", "RoleBinding.Namespace", roleBindingDefinition.Namespace, "RoleBinding.Name", roleBindingDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get roleBinding resource " + instance.Name + "-rolebinding" + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes RoleBindings are not updated - see deployment section
		log.Info("Found RoleBinding")
	}
	return ctrl.Result{}, nil
}
