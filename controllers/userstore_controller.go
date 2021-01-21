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

package controllers

import (
	"context"
	"crypto/tls"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"

	wso2v1beta1 "github.com/wso2/k8s-wso2is-operator/api/v1beta1"

	b64 "encoding/base64"
	"encoding/json"
)

// UserstoreReconciler reconciles a Userstore object
type UserstoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wso2.wso2.com,resources=userstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wso2.wso2.com,resources=userstores/status,verbs=get;update;patch

func (r *UserstoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("userstore", req.NamespacedName)

	//isInstance := wso2v1beta1.Wso2Is{}
	usInstance := wso2v1beta1.Userstore{}

	// Check if WSO2 custom resource is present
	err := r.Get(ctx, req.NamespacedName, &usInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("UserStore resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get UserStore Instance")
		return ctrl.Result{}, err
	}

	GenerateUserstore(usInstance, log)

	return ctrl.Result{}, nil
}

func SpecToJson(spec wso2v1beta1.UserstoreSpec, log logr.Logger) string {
	a, err := json.Marshal(spec)
	if err != nil {
		log.Error(err, "Failed to get parse json")
	}
	return string(a)
}

func GenerateUserstore(instance wso2v1beta1.Userstore, log logr.Logger) {
	url := "https://" + instance.Auth.Host + "/api/server/v1/userstores"
	method := "POST"

	log.Info("Reading from: " + url)

	payload := strings.NewReader(SpecToJson(instance.Spec, log))

	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: instance.Spec.InsecureSkipVerify}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Error(err, "An error has occured")
		return
	}
	encodedToken := b64.StdEncoding.EncodeToString([]byte(instance.Auth.Username + ":" + instance.Auth.Password))
	req.Header.Add("Authorization", "Basic "+encodedToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Error(err, "Unable to send request")
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 201 {
		log.Info("UserStore has been successfully created")
	} else {
		log.Error(err, "Error "+strconv.Itoa(res.StatusCode)+" has occurred during UserStore creation")
	}
}

func (r *UserstoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wso2v1beta1.Userstore{}).
		Complete(r)
}
