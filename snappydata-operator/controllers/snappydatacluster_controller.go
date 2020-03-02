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
	"github.com/prometheus/common/log"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	snappydatav1beta1 "snappydata-operator/api/v1beta1"
)

// SnappyDataClusterReconciler reconciles a SnappyDataCluster object
type SnappyDataClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// name of our custom finalizer
var snappyFinalizerString = "finalizers.snappydata.cloud.tibco.com"

// +kubebuilder:rbac:groups=snappydata.cloud.tibco.com,resources=snappydataclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=snappydata.cloud.tibco.com,resources=snappydataclusters/status,verbs=get;update;patch
func (r *SnappyDataClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log.WithValues("snappydatacluster", req.NamespacedName)

	// your logic here
	l.Info("Inside Reconcile...")
	snappy := snappydatav1beta1.SnappyDataCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, &snappy); err != nil {
		log.Error(err, "failed to get SnappyDataCluster resource")
		// Ignore NotFound errors as they will be retried automatically if the
		// resource is created in future.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if snappy.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(snappy.ObjectMeta.Finalizers, snappyFinalizerString) {
			snappy.ObjectMeta.Finalizers = append(snappy.ObjectMeta.Finalizers, snappyFinalizerString)
			if err := r.Update(context.Background(), &snappy); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(snappy.ObjectMeta.Finalizers, snappyFinalizerString) {
			// TODO: Add the deletion logic here
			return r.handleDelete(&snappy)
		}
		return ctrl.Result{}, nil
	}

	// // start the required headless services
	// l.Info("Starting headless services")
	// if result, err := r.startHeadLessServices(&snappy); err != nil {
	// 	return result, err
	// }

	// // start the required LB services
	// l.Info("Starting load balancer services")
	// if result, err := r.startLBServices(&snappy); err != nil {
	// 	return result, err
	// }

	// l.Info("Starting locator statefulset")
	// if result, err := r.startLocatorStatefulset(&snappy); err != nil {
	// 	return result, err
	// }

	// l.Info("Starting server statefulset")
	// if result, err := r.startServerStatefulset(&snappy); err != nil {
	// 	return result, err
	// }

	// l.Info("Starting lead statefulset")
	// if result, err := r.startLeaderStatefulset(&snappy); err != nil {
	// 	return result, err
	// }

	if result, err := r.reconcileSnappyObjects(&snappy); err != nil {
		return result, err
	}

	// update the cluster status
	snappy.Status.State = snappydatav1beta1.RunningCluster
	if err := r.Status().Update(ctx, &snappy); err != nil {
		return ctrl.Result{}, err
	}

	// l.Info("snappy cluster ", "details", snappy)

	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) reconcileSnappyObjects(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	l := r.Log.WithValues("snappydatacluster", snappy.Namespace)
	// TODO: Add the deletion logic here
	l.Info("Inside reconcileSnappyObjects...")

	switch snappy.Spec.Size {
	case snappydatav1beta1.SmallCluster:
		snappy.Spec.Configuration = SmallClusterConfiguration
	case snappydatav1beta1.MediumCluster:
		snappy.Spec.Configuration = MediumClusterConfiguration
	case snappydatav1beta1.LargeCluster:
		snappy.Spec.Configuration = LargeClusterConfiguration
	}

	// start the required headless services
	l.Info("Starting headless services")
	if result, err := r.startHeadLessServices(snappy); err != nil {
		return result, err
	}

	// start the required LB services
	l.Info("Starting load balancer services")
	if result, err := r.startLBServices(snappy); err != nil {
		return result, err
	}

	l.Info("Starting locator statefulset")
	if result, err := r.startLocatorStatefulset(snappy); err != nil {
		return result, err
	}

	l.Info("Starting server statefulset")
	if result, err := r.startServerStatefulset(snappy); err != nil {
		return result, err
	}

	l.Info("Starting lead statefulset")
	if result, err := r.startLeaderStatefulset(snappy); err != nil {
		return result, err
	}
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) handleDelete(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log.WithValues("snappydatacluster", snappy.Namespace)
	// TODO: Add the deletion logic here
	l.Info("Inside handleDelete...cluster is being deleted")
	leader := &appsv1.StatefulSet{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: snappy.ObjectMeta.Name + "-lead", Namespace: snappy.Namespace}, leader)
	if err != nil && errors.IsNotFound(err) {
	} else if err == nil {
		l.Info("Deleting leader")
		deletepolicy := metav1.DeletePropagationForeground
		err = r.Delete(ctx, leader, &client.DeleteOptions{PropagationPolicy: &deletepolicy})
		l.Info("Deleted leader")
	}

	// remove our finalizer from the list and update it.
	snappy.ObjectMeta.Finalizers = removeString(snappy.ObjectMeta.Finalizers, snappyFinalizerString)
	if err := r.Update(context.Background(), snappy); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// SetupWithManager sets up SnappyData controller and also informs the manager that
// controller owns some objects such as StatefulSets and Services so that it will automatically
// call Reconcile on the underlying SnappyDataCluster when those objects change
func (r *SnappyDataClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&snappydatav1beta1.SnappyDataCluster{}).
		// Owns(&appsv1.StatefulSet{}).
		// Owns(&corev1.Service{}).
		Complete(r)
}
