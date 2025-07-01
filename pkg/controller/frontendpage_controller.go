package controller

import (
	"context"
	"reflect"

	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	frontendv1alpha1 "k8s-controller-tmpl/pkg/apis/frontend/v1alpha1"
)

type FrontendpageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func buildConfigMap(page *frontendv1alpha1.FrontendPage) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Data: map[string]string{
			"content": page.Spec.Contents,
		},
	}
}

func buildDeployment(page *frontendv1alpha1.FrontendPage) *appsv1.Deployment {
	replicas := int32(page.Spec.Replicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": page.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": page.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "frontend",
							Image: page.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "contents",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "contents",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: page.Name},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *FrontendpageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var page frontendv1alpha1.FrontendPage

	err := r.Get(ctx, req.NamespacedName, &page)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get FrontendPage: %s/%s", req.Namespace, req.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Processing configmap
	cm := buildConfigMap(&page)
	if err := ctrl.SetControllerReference(&page, cm, r.Scheme); err != nil {
		log.Error().Err(err).Msgf("Failed to set controller reference for ConfigMap: %s/%s", req.Namespace, req.Name)
		return ctrl.Result{}, err
	}

	log.Info().Msgf("Reconciling ConfigMap for FrontendPage: %s/%s", req.Namespace, req.Name)
	var existingCM corev1.ConfigMap

	if err := r.Get(ctx, req.NamespacedName, &existingCM); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
	} else if !reflect.DeepEqual(existingCM.Data, cm.Data) {
		log.Info().Msgf("Updating ConfigMap for FrontendPage: %s/%s", req.Namespace, req.Name)

		existingCM.Data = cm.Data
		if err := r.Update(ctx, &existingCM); err != nil {
			log.Error().Err(err).Msgf("Failed to update ConfigMap: %s/%s", req.Namespace, req.Name)
			return ctrl.Result{}, err
		}
	} else {
		log.Info().Msgf("ConfigMap for FrontendPage is up-to-date: %s/%s", req.Namespace, req.Name)
	}

	// Processing deployment
	deployment := buildDeployment(&page)
	if err := ctrl.SetControllerReference(&page, deployment, r.Scheme); err != nil {
		log.Error().Err(err).Msgf("Failed to set controller reference for Deployment: %s/%s", req.Namespace, req.Name)
		return ctrl.Result{}, err
	}
	var existingDeployment appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &existingDeployment); err != nil {
		if !errors.IsNotFound(err) {
			log.Error().Err(err).Msgf("Failed to get Deployment: %s/%s", req.Namespace, req.Name)
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, deployment); err != nil {
			log.Error().Err(err).Msgf("Failed to create Deployment: %s/%s", req.Namespace, req.Name)
			return ctrl.Result{}, err
		}
	} else if !reflect.DeepEqual(existingDeployment.Spec, deployment.Spec) {
		updated := false

		log.Info().Msgf("Updating Deployment for FrontendPage: %s/%s", req.Namespace, req.Name)

		if *existingDeployment.Spec.Replicas != *deployment.Spec.Replicas {
			existingDeployment.Spec.Replicas = deployment.Spec.Replicas
			updated = true
		}

		if existingDeployment.Spec.Template.Spec.Containers[0].Image != deployment.Spec.Template.Spec.Containers[0].Image {
			existingDeployment.Spec.Template.Spec.Containers[0].Image = deployment.Spec.Template.Spec.Containers[0].Image
			updated = true
		}

		if updated {
			if err := r.Update(ctx, &existingDeployment); err != nil {
				if errors.IsConflict(err) {
					// Requeue to try again with the latest version
					return ctrl.Result{Requeue: true}, nil
				}
				log.Error().Err(err).Msgf("Failed to update Deployment: %s/%s", req.Namespace, req.Name)
				return ctrl.Result{}, err
			}
		}

	} else {
		log.Info().Msgf("Deployment for FrontendPage is up-to-date: %s/%s", req.Namespace, req.Name)
	}

	return ctrl.Result{}, nil
}

func AddFrontendController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1alpha1.FrontendPage{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(&FrontendpageReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
}
