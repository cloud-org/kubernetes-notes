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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "crd-sample/api/v1"
)

// WebsiteReconciler reconciles a Website object
type WebsiteReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.ronething.cn,resources=websites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.ronething.cn,resources=websites/status,verbs=get;update;patch

func (r *WebsiteReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	//_ = context.Background()
	//_ = r.Log.WithValues("website", req.NamespacedName)

	// your logic here
	ctx := context.Background()
	l := r.Log.WithValues("meta", req.NamespacedName)

	//log.Println("namespaced name is ", req.NamespacedName) // meta.namespace/meta.name

	obj := &webappv1.Website{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		l.Info("Unable to fetch object " + err.Error())
	} else {
		l.Info(obj.Spec.FirstName + " " + obj.Spec.LastName)
	}

	obj.Status.Status = "Running"
	if err := r.Status().Update(ctx, obj); err != nil {
		l.Info("unable to update status err: " + err.Error())
	}

	return ctrl.Result{}, nil
}

func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Website{}).
		Complete(r)
}
