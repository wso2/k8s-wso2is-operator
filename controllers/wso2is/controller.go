/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wso2is

import (
	"context"
	"github.com/wso2/k8s-wso2is-operator/variables"
	"k8s.io/apimachinery/pkg/types"

	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.wso2.com,resources=wso2is/status,verbs=get;update;patch

func (r *Wso2IsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log
	logger.Info("\n-----------------------\nTriggered Reconcile Method\n-----------------------\n")
	logger.Info("Triggered resource : ", "NamespacedName", req.NamespacedName)

	labelSelector := client.MatchingLabels{
		"app": "wso2is",
	}
	instanceList := &wso2v1beta1.Wso2IsList{}
	err := r.List(ctx, instanceList, labelSelector)
	if err != nil {
		logger.Error(err, "Failed to list YourResource objects")
		return ctrl.Result{}, err
	}
	instance := wso2v1beta1.Wso2Is{}
	if len(instanceList.Items) != 0 {
		instance = instanceList.Items[0]
	} else {
		logger.Info("WSO2IS resource not found. Ignoring since object must be deleted")
		return ctrl.Result{}, err
	}

	_, err = reconcileSva(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileSvc(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileVolume(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileConfigMap(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileRole(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileRoleBinding(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileSecret(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: not sure of following.

	_, err = reconcileIngress(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = reconcileStatefulSet(r, instance, logger, err, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	//updateStatus(r, instance, logger, err, ctx)

	return ctrl.Result{}, nil
}

func (r *Wso2IsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wso2v1beta1.Wso2Is{}).
		Watches(
			&source.Kind{Type: &corev1.ConfigMap{}},
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
				if a.GetNamespace() == variables.Wso2IsNamespace {
					// Check if the ConfigMap has the required label
					// This prevents unnecessary reconcile triggers.
					labels := a.GetLabels()
					if labels["app"] == "wso2is" {
						return []reconcile.Request{
							{
								NamespacedName: types.NamespacedName{
									Name:      a.GetName(),
									Namespace: a.GetNamespace(),
								},
							},
						}
					}
				}
				return []reconcile.Request{}
			}),
		).
		Complete(r)
}
