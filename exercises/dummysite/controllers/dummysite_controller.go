/*
Copyright 2021.

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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dummysitev1 "github.com/strax/dwk-2021/exercises/dummysite/api/v1"
)

// DummySiteReconciler reconciles a DummySite object
type DummySiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dummysite.strax.xyz,resources=dummysites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dummysite.strax.xyz,resources=dummysites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dummysite.strax.xyz,resources=dummysites/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;create;update;patch;watch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;create;update;patch;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DummySite object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *DummySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("dummysite", req.NamespacedName)
	var resource dummysitev1.DummySite
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		log.Error(err, "unable to fetch DummySite")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.V(1).Info("fetched DummySite", "WebsiteURL", resource.Spec.WebsiteURL)

	deployment, err := r.newDeploymentObject(&resource)
	if err != nil {
		log.Error(err, "could not construct deployment from template")
		return ctrl.Result{}, nil
	}
	if err := r.Create(ctx, deployment); err != nil {
		log.Error(err, "could not reconcile Deployment", "deployment", deployment)
		return ctrl.Result{}, nil
	}
	log.V(1).Info("created Deployment for DummySite", "name", req.NamespacedName)

	service, err := r.newServiceObject(&resource)
	if err != nil {
		log.Error(err, "could not construct service from template")
		return ctrl.Result{}, nil
	}
	if err := r.Create(ctx, service); err != nil {
		log.Error(err, "could not reconcile Service", "service", service)
		return ctrl.Result{}, nil
	}
	log.V(1).Info("created Service for DummySite", "name", req.NamespacedName)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dummysitev1.DummySite{}).
		Complete(r)
}

func (r *DummySiteReconciler) newDeploymentObject(dummySite *dummysitev1.DummySite) (*appsv1.Deployment, error) {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummySite.Name,
			Namespace: dummySite.Namespace,
			Labels: map[string]string{
				"app": dummySite.Name,
			},
			Annotations: make(map[string]string),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": dummySite.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": dummySite.Name,
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "mirror",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "mirror",
									ReadOnly:  true,
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:  "wget",
							Image: "jgoclawski/wget",
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "mirror",
									MountPath: "/mirror",
								},
							},
							Command: []string{"wget"},
							Args:    []string{"--mirror", "--directory-prefix=/mirror", "--no-host-directories", "--convert-links", "--adjust-extension", "--page-requisites", "--no-parent", dummySite.Spec.WebsiteURL},
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(dummySite, deployment, r.Scheme); err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *DummySiteReconciler) newServiceObject(dummySite *dummysitev1.DummySite) (*v1.Service, error) {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummySite.Name,
			Namespace: dummySite.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": dummySite.Name,
			},
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(dummySite, service, r.Scheme); err != nil {
		return nil, err
	}
	return service, nil
}
