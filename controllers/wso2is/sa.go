package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineServiceAccount(m wso2v1beta1.Wso2Is) *corev1.ServiceAccount {
	svc := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-serviceaccount",
			Namespace: m.Namespace,
		},
	}
	ctrl.SetControllerReference(&m, svc, r.Scheme)
	return svc
}

func reconcileSva(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	svaDefinition := r.defineServiceAccount(instance)
	sva := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-serviceaccount", Namespace: instance.Namespace}, sva)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ServiceAccount resource " + instance.Name + "-serviceaccount" + " not found. Creating or re-creating service account")
			err = r.Create(ctx, svaDefinition)
			if err != nil {
				log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", svaDefinition.Namespace, "ServiceAccount.Name", svaDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get service account resource " + instance.Name + "-serviceaccount" + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes ServiceAccounts are not updated - see deployment section
		log.Info("Found ServiceAccount")
	}
	return ctrl.Result{}, nil
}
