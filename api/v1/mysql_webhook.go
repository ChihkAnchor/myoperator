/*
Copyright 2024 zhangzhikai.

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

package v1

import (
	"context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var mysqllog = logf.Log.WithName("mysql-webhook-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *MySQL) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-database-zhangzhikai-com-cn-v1-mysql,mutating=false,failurePolicy=fail,sideEffects=None,groups=database.zhangzhikai.com.cn,resources=mysqls,verbs=create;update,versions=v1,name=vmysql.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &MySQL{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MySQL) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	mysqllog.Info("validate create", "name", r.Name)
	return nil, r.validateMySQL()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MySQL) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	newMySQL := newObj.(*MySQL)
	oldMySQL := oldObj.(*MySQL)
	if !reflect.DeepEqual(oldMySQL.Spec, newMySQL.Spec) {
		return nil, apierrors.NewBadRequest("update is not allowed")
	}

	return nil, newMySQL.validateMySQL()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MySQL) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

func (r *MySQL) validateMySQL() error {
	var allErrs field.ErrorList

	// Size must be greater than 0
	if r.Spec.Size < 1 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("size"), r.Spec.Size, "size must be greater than 0"))
	}

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "database.zhangzhikai.com.cn", Kind: "MySQL"}, r.Name, allErrs)
}
