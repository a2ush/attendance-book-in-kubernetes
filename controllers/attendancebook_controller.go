/*
Copyright 2022 a2ush.

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
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	officev1alpha1 "github.com/a2ush/attendance-book-in-kubernetes/api/v1alpha1"
)

// AttendanceBookReconciler reconciles a AttendanceBook object
type AttendanceBookReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=office.a2ush.dev,resources=attendancebooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=office.a2ush.dev,resources=attendancebooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=office.a2ush.dev,resources=attendancebooks/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AttendanceBook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *AttendanceBookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	instance := &officev1alpha1.AttendanceBook{}

	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating new resource %s/%s\n", instance.Namespace, instance.Name)
		// err = r.Create(ctx, instance)
		// if err != nil {
		// 	return reconcile.Result{}, err
		// }
		// r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created resource %s/%s", instance.Namespace, instance.Name))
	}
	// if err != nil {
	// 	log.Println(err)
	// 	return ctrl.Result{}, err
	// }

	// if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	log.Printf("Deleting resource %s/%s\n", instance.Namespace, instance.Name)
	// 	err = r.Delete(ctx, instance)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// 	r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s/%s", instance.Namespace, instance.Name))
	// 	return ctrl.Result{}, nil
	// }

	desire := instance.DeepCopy()

	if desire.Status.Attendance != instance.Spec.Attendance || desire.Status.Reason != instance.Spec.Reason {
		desire.Status.Attendance = instance.Spec.Attendance
		desire.Status.Reason = instance.Spec.Reason
		err = r.Status().Update(ctx, desire)
		if err != nil {
			return reconcile.Result{}, err
		}
		log.Printf("Updated resource %s/%s\n", instance.Namespace, instance.Name)
		r.Recorder.Event(instance, corev1.EventTypeNormal, "Updated", fmt.Sprintf("Attendance: %s, Reason: %s", desire.Status.Attendance, desire.Status.Reason))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AttendanceBookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&officev1alpha1.AttendanceBook{}).
		Complete(r)
}
