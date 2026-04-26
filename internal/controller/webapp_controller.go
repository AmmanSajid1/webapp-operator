/*
Copyright 2026.

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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/AmmanSajid1/webapp-operator/api/v1"
	appsv1apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebAppReconciler reconciles a WebApp object
type WebAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.amman.dev,resources=webapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.amman.dev,resources=webapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.amman.dev,resources=webapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/reconcile
func (r *WebAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// TODO(user): your logic here
	// 1. Fetch the WebApp custom resource
	var webapp appsv1.WebApp
	if err := r.Get(ctx, req.NamespacedName, &webapp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Define the Deployment we want Kubernetes to have
	deployment := &appsv1apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Name + "-deployment",
			Namespace: webapp.Namespace,
		},
		Spec: appsv1apps.DeploymentSpec{
			Replicas: &webapp.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": webapp.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": webapp.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webapp",
							Image: webapp.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// 3. Set WebApp as the owner of the Deployment
	if err := ctrl.SetControllerReference(&webapp, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// 4. Check if Deployment already exists
	var existing appsv1apps.Deployment
	err := r.Get(ctx, types.NamespacedName{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}, &existing)

	if err != nil && apierrors.IsNotFound(err) {
		log.Info("Creating Deployment", "name", deployment.Name)
		return ctrl.Result{}, r.Create(ctx, deployment)
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	needsUpdated := false

	// Check if replicas changed
	if existing.Spec.Replicas == nil || *existing.Spec.Replicas != webapp.Spec.Replicas {
		existing.Spec.Replicas = &webapp.Spec.Replicas
		needsUpdated = true
	}

	// Check if image changed
	currentImage := existing.Spec.Template.Spec.Containers[0].Image
	if currentImage != webapp.Spec.Image {
		existing.Spec.Template.Spec.Containers[0].Image = webapp.Spec.Image
		needsUpdated = true
	}

	// Update deployment if needed
	if needsUpdated {
		log.Info("Updating Deployment", "name", existing.Name)
		if err := r.Update(ctx, &existing); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("Deployment is up to date", "name", existing.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.WebApp{}).
		Named("webapp").
		Complete(r)
}
