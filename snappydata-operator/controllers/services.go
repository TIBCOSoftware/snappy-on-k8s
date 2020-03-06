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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	snappydatav1beta1 "snappydata-operator/api/v1beta1"
)

func (r *SnappyDataClusterReconciler) createServiceIfNotExists(snappy *snappydatav1beta1.SnappyDataCluster,
	service *corev1.Service) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log

	// set the owner reference for the service
	// TODO: handle the case when SetControllerReference throws error

	// l.Info("set owner ref for service", "service", *service)
	if err := ctrl.SetControllerReference(snappy, service, r.Scheme); err != nil {
		l.Error(err, "unable to set owner reference for service", "service", *service)
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	found := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		l.Info("Creating service", "service", service.Namespace+"."+service.Name)
		r.Recorder.Eventf(snappy, corev1.EventTypeNormal, "Info", "Creating service %s in namespace %s", service.Name, service.Namespace)
		if err := r.Create(ctx, service); err != nil {
			l.Error(err, "unable to create service", "service", *service)
			r.Recorder.Eventf(snappy, corev1.EventTypeWarning, "Warning",
				"Unable to create service %s in namespace %s due to error: %s",
				service.Name, service.Namespace, err.Error())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(service.Spec, found.Spec) {
		found.Spec = service.Spec
		l.Info("Updating service", "service", service.Namespace+"."+service.Name)
		r.Recorder.Eventf(snappy, corev1.EventTypeNormal, "Info", "Updating service %s in namespace %s", service.Name, service.Namespace)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil

}

func (r *SnappyDataClusterReconciler) createService(snappy *snappydatav1beta1.SnappyDataCluster,
	service *corev1.Service) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log

	// set the owner reference for the service
	// TODO: handle the case when SetControllerReference throws error

	// l.Info("set owner ref for service", "service", *service)
	if err := ctrl.SetControllerReference(snappy, service, r.Scheme); err != nil {
		l.Error(err, "unable to set owner reference for service", "service", *service)
		return ctrl.Result{}, err
	}

	// l.Info("creating service", "service", *service)
	if err := r.Create(ctx, service); err != nil {
		l.Error(err, "unable to create service", "service", *service)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) startHeadLessServices(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	// locator
	locatorservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-locator-headless",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "locator",
					Port:       10334,
					TargetPort: intstr.FromInt(10334),
				},
				{
					Name:       "jdbc",
					Port:       1527,
					TargetPort: intstr.FromInt(1527),
				},
			},
			ClusterIP: "None",
			Selector:  map[string]string{"app": snappy.ObjectMeta.Name + "-locator"},
		},
	}
	//TODO: handle errors
	// r.createService(snappy, &locatorservice)
	r.createServiceIfNotExists(snappy, &locatorservice)

	// server
	serverservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-server-headless",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jdbc",
					Port:       1527,
					TargetPort: intstr.FromInt(1527),
				},
			},
			ClusterIP: "None",
			Selector:  map[string]string{"app": snappy.ObjectMeta.Name + "-server"},
		},
	}
	// r.createService(snappy, &serverservice)
	r.createServiceIfNotExists(snappy, &serverservice)

	// lead
	leadservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-lead-headless",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "sparkui",
					Port:       5050,
					TargetPort: intstr.FromInt(5050),
				},
			},
			ClusterIP: "None",
			Selector:  map[string]string{"app": snappy.ObjectMeta.Name + "-lead"},
		},
	}
	// r.createService(snappy, &leadservice)
	r.createServiceIfNotExists(snappy, &leadservice)
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) startLBServices(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	// locator
	locatorlbservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-locator-public",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jdbc",
					Port:       1527,
					TargetPort: intstr.FromInt(1527),
				},
			},
			Type:     "LoadBalancer",
			Selector: map[string]string{"app": snappy.ObjectMeta.Name + "-locator"},
		},
	}
	//TODO: handle errors
	// r.createService(snappy, &locatorlbservice)
	r.createServiceIfNotExists(snappy, &locatorlbservice)

	// server
	serverlbservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-server-public",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jdbc",
					Port:       1527,
					TargetPort: intstr.FromInt(1527),
				},
			},
			Type:     "LoadBalancer",
			Selector: map[string]string{"app": snappy.ObjectMeta.Name + "-server"},
		},
	}
	//TODO: handle errors
	// r.createService(snappy, &serverlbservice)
	r.createServiceIfNotExists(snappy, &serverlbservice)

	// lead
	leadlbservice := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-lead-public",
			Namespace: snappy.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "sparkui",
					Port:       5050,
					TargetPort: intstr.FromInt(5050),
				},
				{
					Name:       "jobserver",
					Port:       8090,
					TargetPort: intstr.FromInt(8090),
				},
				{
					Name:       "zeppelin-interpreter",
					Port:       3768,
					TargetPort: intstr.FromInt(3768),
				},
			},
			Type:     "LoadBalancer",
			Selector: map[string]string{"app": snappy.ObjectMeta.Name + "-lead"},
		},
	}
	//TODO: handle errors
	// r.createService(snappy, &leadlbservice)
	r.createServiceIfNotExists(snappy, &leadlbservice)
	return ctrl.Result{}, nil
}
