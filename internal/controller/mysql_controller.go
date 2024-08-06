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

package controller

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasev1alpha1 "github.com/ChihkAnchor/myoperator/api/v1"
)

// MySQLReconciler reconciles a MySQL object
type MySQLReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=database.example.com,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.example.com,resources=mysqls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=database.example.com,resources=mysqls/finalizers,verbs=update

func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the MySQL instance
	mysql := &databasev1alpha1.MySQL{}
	if err := r.Get(ctx, req.NamespacedName, mysql); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define a new Deployment
	dep := r.deploymentForMySQL(mysql)
	// Set MySQL instance as the owner and controller
	if err := controllerutil.SetControllerReference(mysql, dep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: dep.Name, Namespace: dep.Namespace}, found)
	if err != nil && client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	} else if err != nil && client.IgnoreNotFound(err) == nil {
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err := r.Create(ctx, dep); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// Update the Deployment if necessary
		if !reflect.DeepEqual(dep.Spec, found.Spec) {
			log.Info("Updating Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			found.Spec = dep.Spec
			if err := r.Update(ctx, found); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForMySQL returns a MySQL Deployment object
func (r *MySQLReconciler) deploymentForMySQL(m *databasev1alpha1.MySQL) *appsv1.Deployment {
	labels := map[string]string{"app": "mysql", "mysql_cr": m.Name}
	replicas := m.Spec.Size
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "mysql",
						Image: "mysql:5.7",
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "password", // 注意：这是一个简单的示例，实际应用中应使用更安全的方式管理密码
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
					}},
				},
			},
		},
	}
	return dep
}

// SetupWithManager sets up the controller with the Manager.
func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.MySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
