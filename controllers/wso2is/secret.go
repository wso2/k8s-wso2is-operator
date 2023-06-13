package wso2is

import (
	"context"

	"github.com/go-logr/logr"
	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	"github.com/wso2/k8s-wso2is-operator/variables"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Wso2IsReconciler) defineSecret(m wso2v1beta1.Wso2Is) *corev1.Secret {
	//append all secrets to here
	secretsMap := map[string][]byte{}
	//mount keystore secrets to Kubernetes secrets
	for _, element := range m.Spec.KeystoreMounts {
		secretsMap[element.Name] = []byte(element.Data)
	}
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      variables.SecretName,
			Namespace: m.Namespace,
		},
		Data: secretsMap,
	}
	ctrl.SetControllerReference(&m, secret, r.Scheme)
	return secret
}
func reconcileSecret(r *Wso2IsReconciler, instance wso2v1beta1.Wso2Is, log logr.Logger, err error, ctx context.Context) (ctrl.Result, error) {
	secretDefinition := r.defineSecret(instance)
	secret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: variables.SecretName, Namespace: instance.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Secret resource " + variables.SecretName + " not found. Creating or re-creating secret")
			err = r.Create(ctx, secretDefinition)
			if err != nil {
				log.Error(err, "Failed to create new Secret", "Secret.Namespace", secretDefinition.Namespace, "Secret.Name", secretDefinition.Name)
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get secret resource " + variables.SecretName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes Secrets are not updated - see deployment section

		log.Info("Found Secret")
	}
	return ctrl.Result{}, nil
}
