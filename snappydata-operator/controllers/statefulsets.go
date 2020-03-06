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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	snappydatav1beta1 "snappydata-operator/api/v1beta1"
)

func (r *SnappyDataClusterReconciler) createStatefulsetIfNotExists(snappy *snappydatav1beta1.SnappyDataCluster,
	statefulset *appsv1.StatefulSet) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log

	// set the owner reference for the statefulset
	// TODO: handle the case when SetControllerReference throws error

	// l.Info("set owner ref for statefulset", "statefulset", *statefulset)
	if err := ctrl.SetControllerReference(snappy, statefulset, r.Scheme); err != nil {
		l.Error(err, "unable to set owner reference for statefulset", "statefulset", *statefulset)
		return ctrl.Result{}, err
	}

	found := &appsv1.StatefulSet{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: statefulset.Name, Namespace: statefulset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		l.Info("Creating statefulset", "statefulset", statefulset.Name+"."+statefulset.Namespace)
		r.Recorder.Eventf(snappy, corev1.EventTypeNormal, "Info", "Creating statefulset %s in namespace %s", statefulset.Name, statefulset.Namespace)
		if err := r.Create(ctx, statefulset); err != nil {
			l.Error(err, "unable to create statefulset", "statefulset", *statefulset)
			r.Recorder.Eventf(snappy, corev1.EventTypeWarning, "Warning",
				"Unable to create statefulset %s in namespace %s due to error: %s",
				statefulset.Name, statefulset.Namespace, err.Error())
			return ctrl.Result{}, err
		}
		return reconcile.Result{}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(statefulset.Spec, found.Spec) {
		found.Spec = statefulset.Spec
		l.Info("Updating statefulset", "statefulset", statefulset.Name+"."+statefulset.Namespace)
		r.Recorder.Eventf(snappy, corev1.EventTypeNormal, "Info", "Updating statefulset %s in namespace %s", statefulset.Name, statefulset.Namespace)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) createStatefulset(snappy *snappydatav1beta1.SnappyDataCluster,
	statefulset *appsv1.StatefulSet) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log

	// set the owner reference for the statefulset
	// TODO: handle the case when SetControllerReference throws error

	// l.Info("set owner ref for statefulset", "statefulset", *statefulset)
	if err := ctrl.SetControllerReference(snappy, statefulset, r.Scheme); err != nil {
		l.Error(err, "unable to set owner reference for statefulset", "statefulset", *statefulset)
		return ctrl.Result{}, err
	}
	// l.Info("creating statefulset", "statefulset", *statefulset)
	if err := r.Create(ctx, statefulset); err != nil {
		l.Error(err, "unable to create statefulset", "statefulset", *statefulset)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func getVolumeMounts(volumeclaimtemplates []v1.PersistentVolumeClaim) []v1.VolumeMount {
	length := len(volumeclaimtemplates)
	vms := make([]v1.VolumeMount, length)
	// update the names for PVCs
	for index := range volumeclaimtemplates {
		claim := &volumeclaimtemplates[index]
		claim.Name = "disk" + strconv.Itoa(index)
		// first PVC assumed to be datadirectory for Snappy
		if index == 0 {
			vms[index] = v1.VolumeMount{
				Name:      "disk0",
				MountPath: "/opt/snappydata/work",
			}
		} else {
			vms[index] = v1.VolumeMount{
				Name:      claim.Name,
				MountPath: "/" + claim.Name,
			}
		}
	}
	return vms
}

func (r *SnappyDataClusterReconciler) startLocatorStatefulset(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	// expectedreplicas := int32(1)
	expectedreplicas := snappy.Spec.Configuration.Locators.Count
	terminationgraceperiodseconds := int64(10)
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"app": snappy.ObjectMeta.Name + "-locator"}}
	lifecycle := corev1.Lifecycle{
		PreStop: &corev1.Handler{Exec: &corev1.ExecAction{Command: []string{"/opt/snappydata/sbin/snappy-locators.sh", "stop"}}},
	}
	livenessprobe := corev1.Probe{
		Handler: corev1.Handler{Exec: &corev1.ExecAction{
			Command: []string{
				"/bin/sh", "-c", "/opt/snappydata/sbin/snappy-locators.sh status | grep -e running -e waiting"}}},
		InitialDelaySeconds: int32(360),
	}

	vm := v1.VolumeMount{
		Name:      "disk0",
		MountPath: "/opt/snappydata/work",
	}

	volumemounts := []v1.VolumeMount{vm}
	// volumemounts := getVolumeMounts(snappy.Spec.Configuration.Locators.VolumeClaimTemplates)

	userconf := snappy.Spec.Configuration.Locators.Conf
	r.Log.Info("startLocatorStatefulset userconf is " + userconf)
	podtemplate := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"app": snappy.ObjectMeta.Name + "-locator"},
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationgraceperiodseconds,
			ImagePullSecrets:              []corev1.LocalObjectReference{{Name: "regcred"}},
			Containers: []corev1.Container{
				{
					Name:            snappy.ObjectMeta.Name + "-locator",
					Image:           snappy.Spec.Image,
					ImagePullPolicy: snappy.Spec.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "locator",
							ContainerPort: int32(10334),
						},
						{
							Name:          "jdbc",
							ContainerPort: int32(1527),
						},
					},
					LivenessProbe: &livenessprobe,
					Command:       []string{"/bin/bash", "-c", "start locator " + userconf},
					Lifecycle:     &lifecycle,
					Resources:     snappy.Spec.Configuration.Locators.Resources,
					VolumeMounts:  volumemounts,
				},
			},
		},
	}

	for index := range snappy.Spec.Configuration.Locators.VolumeClaimTemplates {
		claim := &snappy.Spec.Configuration.Locators.VolumeClaimTemplates[index]
		claim.Name = "disk" + strconv.Itoa(index)
	}

	locator := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-locator",
			Namespace: snappy.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &expectedreplicas,
			ServiceName:          snappy.ObjectMeta.Name + "-locator-headless",
			Selector:             &selector,
			Template:             podtemplate,
			VolumeClaimTemplates: snappy.Spec.Configuration.Locators.VolumeClaimTemplates,
		},
	}

	// r.createStatefulset(snappy, &locator)
	r.createStatefulsetIfNotExists(snappy, &locator)
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) startServerStatefulset(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	// expectedreplicas := int32(1)
	expectedreplicas := snappy.Spec.Configuration.Servers.Count
	terminationgraceperiodseconds := int64(10)
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"app": snappy.ObjectMeta.Name + "-server"}}
	lifecycle := corev1.Lifecycle{
		PreStop: &corev1.Handler{Exec: &corev1.ExecAction{Command: []string{"/opt/snappydata/sbin/snappy-servers.sh", "stop"}}},
	}
	livenessprobe := corev1.Probe{
		Handler: corev1.Handler{Exec: &corev1.ExecAction{
			Command: []string{
				"/bin/sh", "-c", "/opt/snappydata/sbin/snappy-servers.sh status | grep -e running -e waiting"}}},
		InitialDelaySeconds: int32(360),
	}

	// TODO: we can just check whether service is ready by listing the services here instead of code like this
	waitforservicearg := "--get-ip " + snappy.ObjectMeta.Name + "-server-public --wait-for " + snappy.ObjectMeta.Name + "-locator-headless 10334"
	r.Log.Info("waitforservicearg is ", "waitforservicearg", waitforservicearg)

	vm := v1.VolumeMount{
		Name:      "disk0",
		MountPath: "/opt/snappydata/work",
	}

	volumemounts := []v1.VolumeMount{vm}
	userconf := snappy.Spec.Configuration.Servers.Conf
	podtemplate := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"app": snappy.ObjectMeta.Name + "-server"},
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationgraceperiodseconds,
			ImagePullSecrets:              []corev1.LocalObjectReference{{Name: "regcred"}},
			Containers: []corev1.Container{
				{
					Name:            snappy.ObjectMeta.Name + "-server",
					Image:           snappy.Spec.Image,
					ImagePullPolicy: snappy.Spec.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "jdbc",
							ContainerPort: int32(1527),
						},
					},
					LivenessProbe: &livenessprobe,
					Command: []string{"/bin/bash", "-c", "start server " + waitforservicearg +
						" -locators=" + snappy.ObjectMeta.Name + "-locator-headless:10334 " + userconf},
					Lifecycle:    &lifecycle,
					Resources:    snappy.Spec.Configuration.Servers.Resources,
					VolumeMounts: volumemounts,
				},
			},
		},
	}

	for index := range snappy.Spec.Configuration.Servers.VolumeClaimTemplates {
		claim := &snappy.Spec.Configuration.Servers.VolumeClaimTemplates[index]
		claim.Name = "disk" + strconv.Itoa(index)
	}

	server := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-server",
			Namespace: snappy.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &expectedreplicas,
			ServiceName:          snappy.ObjectMeta.Name + "-server-headless",
			Selector:             &selector,
			Template:             podtemplate,
			VolumeClaimTemplates: snappy.Spec.Configuration.Servers.VolumeClaimTemplates,
		},
	}

	// r.createStatefulset(snappy, &server)
	r.createStatefulsetIfNotExists(snappy, &server)
	return ctrl.Result{}, nil
}

func (r *SnappyDataClusterReconciler) startLeaderStatefulset(snappy *snappydatav1beta1.SnappyDataCluster) (ctrl.Result, error) {
	// expectedreplicas := int32(1)
	expectedreplicas := snappy.Spec.Configuration.Leads.Count
	terminationgraceperiodseconds := int64(10)
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"app": snappy.ObjectMeta.Name + "-lead"}}
	lifecycle := corev1.Lifecycle{
		PreStop: &corev1.Handler{Exec: &corev1.ExecAction{Command: []string{"/opt/snappydata/sbin/snappy-leads.sh", "stop"}}},
	}
	livenessprobe := corev1.Probe{
		Handler: corev1.Handler{Exec: &corev1.ExecAction{
			Command: []string{
				"/bin/sh", "-c", "/opt/snappydata/sbin/snappy-leads.sh status | grep -e running -e waiting"}}},
		InitialDelaySeconds: int32(360),
	}
	// TODO: we can just check whether service is ready by listing the services here instead of code like this
	waitforservicearg := "--get-ip " + snappy.ObjectMeta.Name + "-lead-public --wait-for " + snappy.ObjectMeta.Name + "-server-headless 1527"
	r.Log.Info("waitforservicearg is ", "waitforservicearg", waitforservicearg)

	vm := v1.VolumeMount{
		Name:      "disk0",
		MountPath: "/opt/snappydata/work",
	}

	volumemounts := []v1.VolumeMount{vm}

	userconf := snappy.Spec.Configuration.Leads.Conf
	podtemplate := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"app": snappy.ObjectMeta.Name + "-lead"},
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationgraceperiodseconds,
			ImagePullSecrets:              []corev1.LocalObjectReference{{Name: "regcred"}},
			Containers: []corev1.Container{
				{
					Name:            snappy.ObjectMeta.Name + "-lead",
					Image:           snappy.Spec.Image,
					ImagePullPolicy: snappy.Spec.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "sparkui",
							ContainerPort: int32(5050),
						},
					},
					LivenessProbe: &livenessprobe,
					Command: []string{"/bin/bash", "-c", "start lead " + waitforservicearg +
						" -locators=" + snappy.ObjectMeta.Name + "-locator-headless:10334 " + userconf},
					Lifecycle:    &lifecycle,
					Resources:    snappy.Spec.Configuration.Leads.Resources,
					VolumeMounts: volumemounts,
				},
			},
		},
	}

	for index := range snappy.Spec.Configuration.Leads.VolumeClaimTemplates {
		claim := &snappy.Spec.Configuration.Leads.VolumeClaimTemplates[index]
		claim.Name = "disk" + strconv.Itoa(index)
	}

	leader := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snappy.ObjectMeta.Name + "-lead",
			Namespace: snappy.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &expectedreplicas,
			ServiceName:          snappy.ObjectMeta.Name + "-lead-headless",
			Selector:             &selector,
			Template:             podtemplate,
			VolumeClaimTemplates: snappy.Spec.Configuration.Servers.VolumeClaimTemplates,
		},
	}

	// r.createStatefulset(snappy, &leader)
	r.createStatefulsetIfNotExists(snappy, &leader)
	return ctrl.Result{}, nil
}
