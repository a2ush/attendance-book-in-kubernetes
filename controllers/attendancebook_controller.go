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
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	officev1alpha1 "github.com/a2ush/attendance-book-in-kubernetes/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AttendanceBookReconciler reconciles a AttendanceBook object
type AttendanceBookReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var specified_namespace string
var employeeList []string

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
		r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s", req.NamespacedName.String()))
		return reconcile.Result{}, nil
	}

	// Delete object due to not-registerd employee.
	fmt.Println(employeeList)

	// Delete object if deployed namespace is not correct.
	if req.NamespacedName.Namespace != specified_namespace {
		log.Println("Delete due to other namespace.")

		uid := instance.GetUID()
		resourceVersion := instance.GetResourceVersion()
		cond := metav1.Preconditions{
			UID:             &uid,
			ResourceVersion: &resourceVersion,
		}
		err = r.Delete(ctx, instance, &client.DeleteOptions{
			Preconditions: &cond,
		})
		if err != nil {
			log.Println(err)
		}
		r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s due to the namespace that is not allowed to deploy.", req.NamespacedName.String()))

		return reconcile.Result{}, nil
	}

	desire := instance.DeepCopy()

	if desire.Status.Attendance != instance.Spec.Attendance || desire.Status.Reason != instance.Spec.Reason {

		// Status.Attendance is "" when first created.
		if desire.Status.Attendance == "" {
			if err = r.changeStatus(desire, instance); err != nil {
				return reconcile.Result{}, err
			}
			r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created resource. Attendance/Reason: %s/%s", desire.Status.Attendance, desire.Status.Reason))

		} else {
			if err = r.changeStatus(desire, instance); err != nil {
				return reconcile.Result{}, err
			}
			r.Recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated resource. Attendance/Reason: %s/%s", desire.Status.Attendance, desire.Status.Reason))
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AttendanceBookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	specified_namespace = GetNamespace()
	employeeList = ReadEmployeeList("mnt/employee-list/Employeelist")

	return ctrl.NewControllerManagedBy(mgr).
		For(&officev1alpha1.AttendanceBook{}).
		Complete(r)
}

func (r *AttendanceBookReconciler) changeStatus(desire, instance *officev1alpha1.AttendanceBook) error {
	desire.Status.Attendance = instance.Spec.Attendance
	desire.Status.Reason = instance.Spec.Reason
	err := r.Status().Update(context.TODO(), desire)

	return err
}

func GetNamespace() string {
	specified_namespace, found := os.LookupEnv("SPECIFIED_NAMESPACE")
	if !found {
		specified_namespace = "default"
	}
	return specified_namespace
}

func ReadEmployeeList(filename string) []string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	employeeList := make([]string, 0, 10)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmp := scanner.Text()
		log.Println(tmp)
		employeeList = append(employeeList, tmp)
	}

	return employeeList
}
