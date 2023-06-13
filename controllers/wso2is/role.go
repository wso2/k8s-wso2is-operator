package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineRole(m wso2v1beta1.Wso2Is) *rbacv1.Role {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      variables.RoleName,
			Namespace: m.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"get", "list"},
				Resources: []string{"endpoints"},
			},
		},
	}
	ctrl.SetControllerReference(&m, role, r.Scheme)
	return role
}

func reconcileRole(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	roleDefinition := r.defineRole(instance)
	role := &rbacv1.Role{}
	err = r.Get(ctx, types.NamespacedName{Name: variables.RoleName, Namespace: instance.Namespace}, role)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Role resource " + variables.RoleName + " not found. Creating or re-creating role")
			err = r.Create(ctx, roleDefinition)
			if err != nil {
				log.Error(err, "Failed to create new Role", "Role.Namespace", roleDefinition.Namespace, "Role.Name", roleDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get role resource " + variables.RoleName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes Roles are not updated - see deployment section
		log.Info("Found Role")
	}
	return ctrl.Result{}, nil
}
