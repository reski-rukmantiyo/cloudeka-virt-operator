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

package v1alpha1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var cloudekadatavolumelog = logf.Log.WithName("cloudekadatavolume-resource")

var (
	allowedTypes    []string
	allowedVersions map[string][]string
	Repository      string
)

func loadConfig() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("unable to load configuration: %v", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	configMap, err := clientset.CoreV1().ConfigMaps("cloudeka-system").Get(context.TODO(), "cloudekadatavolume-config", metav1.GetOptions{})
	if err != nil {
		return err
	}

	Repository = configMap.Data["repository"]
	if Repository == "" {
		return errors.New("repository is required")
	}

	err = json.Unmarshal([]byte(configMap.Data["types"]), &allowedTypes)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(configMap.Data["versions"]), &allowedVersions)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	err := loadConfig()
	if err != nil {
		cloudekadatavolumelog.Error(err, "failed to load configuration")
	}
}

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *CloudekaDataVolume) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-virt-cloudeka-ai-v1alpha1-cloudekadatavolume,mutating=true,failurePolicy=fail,sideEffects=None,groups=virt.cloudeka.ai,resources=cloudekadatavolumes,verbs=create;update,versions=v1alpha1,name=mcloudekadatavolume.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CloudekaDataVolume{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CloudekaDataVolume) Default() {
	cloudekadatavolumelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

//+kubebuilder:webhook:path=/validate-virt-cloudeka-ai-v1alpha1-cloudekadatavolume,mutating=false,failurePolicy=fail,sideEffects=None,groups=virt.cloudeka.ai,resources=cloudekadatavolumes,verbs=create;update;delete,versions=v1alpha1,name=vcloudekadatavolume.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CloudekaDataVolume{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CloudekaDataVolume) ValidateCreate() (admission.Warnings, error) {
	cloudekadatavolumelog.Info("validate create", "name", r.Name)

	// Validate Type field
	osType := strings.ToLower(r.Spec.Type)
	if !contains(allowedTypes, osType) {
		return nil, fmt.Errorf("invalid type: %s, must be one of %v", r.Spec.Type, allowedTypes)
	}

	// Validate Version field
	versions, exists := allowedVersions[osType]
	if !exists || !contains(versions, r.Spec.Version) {
		if r.Spec.Version == "" {
			return nil, fmt.Errorf("version is required for type: %s, must be one of %v", r.Spec.Type, versions)
		} else {
			return nil, fmt.Errorf("invalid version: %s for type: %s, must be one of %v", r.Spec.Version, r.Spec.Type, versions)
		}
	}

	return nil, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CloudekaDataVolume) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	cloudekadatavolumelog.Info("cannot update for ", "name", r.Name, " values")

	warnings := admission.Warnings{
		"cannot update immutability object - CloudekaDataVolume",
	}

	return warnings, fmt.Errorf("Cannot update")
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CloudekaDataVolume) ValidateDelete() (admission.Warnings, error) {
	cloudekadatavolumelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
