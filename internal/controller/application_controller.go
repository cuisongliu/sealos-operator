/*
Copyright 2025 cuisongliu.

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

package controller

import (
	"context"
	"errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1beta1 "github.com/cuisongliu/sealos-operator/api/v1beta1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Finalizer string
}

// +kubebuilder:rbac:groups=apps.github.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.github.com,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.github.com,resources=applications/finalizers,verbs=update

func (r *ApplicationReconciler) doFinalizerOperationsForSetting(ctx context.Context, app *appsv1beta1.Application) error {
	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	appFinalizer := r.Finalizer
	// Fetch the Setting instance
	// The purpose is check if the Custom Resource for the Kind Setting
	// is applied on the cluster if not we return nil to stop the reconciliation
	application := &appsv1beta1.Application{}

	err := r.Get(ctx, req.NamespacedName, application)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if application.GetDeletionTimestamp() != nil && !application.GetDeletionTimestamp().IsZero() {
		if err = r.doFinalizerOperationsForSetting(ctx, application); err != nil {
			return ctrl.Result{}, err
		}
		if controllerutil.ContainsFinalizer(application, appFinalizer) {
			controllerutil.RemoveFinalizer(application, appFinalizer)
		}
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Update(ctx, application)
		}); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err = r.statusReconcile(ctx, application); err != nil {
					lg := log.FromContext(ctx)
					lg.Error(err, "Failed to update application status")
				}
			}
		}
	}()

	if application.GetDeletionTimestamp().IsZero() || application.GetDeletionTimestamp() == nil {
		controllerutil.AddFinalizer(application, appFinalizer)
		if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Update(ctx, application)
		}); err != nil {
			return ctrl.Result{}, err
		}
		return r.reconcile(ctx, application)
	}

	return ctrl.Result{}, errors.New("reconcile error from Finalizer")
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	r.Scheme = mgr.GetScheme()
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("app-controller")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1beta1.Application{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Named("application").
		Complete(r)
}

func (r *ApplicationReconciler) reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	return ctrl.Result{}, nil
}
