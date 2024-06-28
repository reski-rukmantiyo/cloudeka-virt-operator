/*
Copyright 2024.

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
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/reski-rukmantiyo/cloudeka-virt-operator/api/v1alpha1"
	virtv1alpha1 "github.com/reski-rukmantiyo/cloudeka-virt-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdicontroller "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

// CloudekaDataVolumeReconciler reconciles a CloudekaDataVolume object
type CloudekaDataVolumeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=virt.cloudeka.ai,resources=cloudekadatavolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virt.cloudeka.ai,resources=cloudekadatavolumes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virt.cloudeka.ai,resources=cloudekadatavolumes/finalizers,verbs=update

func (r *CloudekaDataVolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling CloudekaDataVolume")

	// Fetch the CloudekaDataVolume cloudekaDataVolume
	cloudekaDataVolume := &virtv1alpha1.CloudekaDataVolume{}
	err := r.Get(ctx, req.NamespacedName, cloudekaDataVolume)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// url := "docker://docker.io/rrukmantiyo/kubevirt-images:ubuntu-22.04"

	osType := strings.ToLower(cloudekaDataVolume.Spec.Type)
	osVersion := strings.ToLower(cloudekaDataVolume.Spec.Version)
	url := fmt.Sprintf("%s:%s-%s", v1alpha1.Repository, osType, osVersion)

	// Define the desired DataVolume object
	datavolume := &cdicontroller.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cloudekaDataVolume.Name,
			Namespace: cloudekaDataVolume.Namespace,
			Annotations: map[string]string{
				"cdi.kubevirt.io/storage.bind.immediate.requested": "True",
			},
		},
		Spec: cdicontroller.DataVolumeSpec{
			Source: &cdicontroller.DataVolumeSource{
				Registry: &cdicontroller.DataVolumeSourceRegistry{
					URL: &url,
				},
			},
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(cloudekaDataVolume.Spec.Size),
					},
				},
			},
		},
	}

	// Set CloudekaDataVolume instance as the owner and controller
	if err := ctrl.SetControllerReference(cloudekaDataVolume, datavolume, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the DataVolume already exists
	found := &cdicontroller.DataVolume{}
	err = r.Get(ctx, client.ObjectKey{Name: datavolume.Name, Namespace: datavolume.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info("Creating a new DataVolume", "DataVolume.Namespace", datavolume.Namespace, "DataVolume.Name", datavolume.Name)
		err = r.Create(ctx, datavolume)
		if err != nil {
			return ctrl.Result{}, err
		}
		// DataVolume created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// DataVolume already exists - update status
	cloudekaDataVolume.Status.Phase = string(found.Status.Phase)
	cloudekaDataVolume.Status.Progress = string(found.Status.Progress)
	err = r.Status().Update(ctx, cloudekaDataVolume)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CloudekaDataVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&virtv1alpha1.CloudekaDataVolume{}).
		Owns(&cdicontroller.DataVolume{}).
		Complete(r)
}
