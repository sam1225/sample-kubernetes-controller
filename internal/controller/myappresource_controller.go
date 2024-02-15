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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sv1alpha1 "github.com/sam1225/sample-kubernetes-controller/api/v1alpha1"
)

// MyAppResourceReconciler reconciles a MyAppResource object
type MyAppResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8s.myappresources.io,resources=myappresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.myappresources.io,resources=myappresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.myappresources.io,resources=myappresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyAppResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *MyAppResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)
	log := log.FromContext(ctx)
	log.Info("Enter Reconcile", "req", req)

	// TODO(user): your logic here

	var myAppResource k8sv1alpha1.MyAppResource
	log.Info("Fetching MyAppResource")
	if err := r.Get(ctx, req.NamespacedName, &myAppResource); err != nil {
		log.Info("Unable to fetch MyAppResource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.reconcileRedis(ctx, myAppResource)

	r.reconcileDeployment(ctx, myAppResource)
	r.reconcileService(ctx, myAppResource)

	return ctrl.Result{}, nil
}

func (r *MyAppResourceReconciler) reconcileDeployment(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	// log := log.FromContext(ctx)

	deploy := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: myAppResource.Name, Namespace: myAppResource.Namespace}, deploy)
	if err != nil {
		if errors.IsNotFound(err) {
			return r.createDeployment(ctx, myAppResource)
		}

		return err
	}

	// Deployment found - next, check for updates
	r.updateDeployment(ctx, myAppResource, *deploy)

	return nil
}

func (r *MyAppResourceReconciler) createDeployment(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	log := log.FromContext(ctx)

	replicas := myAppResource.Spec.Replicas

	image := fmt.Sprintf("%v:%v", myAppResource.Spec.Image.Repository, myAppResource.Spec.Image.Tag)

	memoryQuantity, _ := resource.ParseQuantity(myAppResource.Spec.Resources.MemoryLimit)
	cpuQuantity, _ := resource.ParseQuantity(myAppResource.Spec.Resources.CPURequest)

	cacheServer := fmt.Sprintf("tcp://%s-redis:6379", myAppResource.Name)

	podEnvVar := []corev1.EnvVar{
		{Name: "PODINFO_UI_COLOR", Value: myAppResource.Spec.UI.Color},
		{Name: "PODINFO_UI_MESSAGE", Value: myAppResource.Spec.UI.Message},
	}

	if myAppResource.Spec.Redis.Enabled {
		podEnvVar = []corev1.EnvVar{
			{Name: "PODINFO_UI_COLOR", Value: myAppResource.Spec.UI.Color},
			{Name: "PODINFO_UI_MESSAGE", Value: myAppResource.Spec.UI.Message},
			{Name: "PODINFO_CACHE_SERVER", Value: cacheServer},
		}
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myAppResource.Name,
			Namespace: myAppResource.Namespace,
			Labels:    map[string]string{"app.kubernetes.io/name": myAppResource.Name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": myAppResource.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app.kubernetes.io/name": myAppResource.Name, "app": myAppResource.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           image,
						Name:            myAppResource.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command:         []string{"./podinfo", "--port=9898", "--level=info", "--random-delay=false", "--random-error=false"},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 9898,
								Protocol:      corev1.ProtocolTCP,
								Name:          "http",
							},
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: memoryQuantity,
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: cpuQuantity,
							},
						},
						Env: podEnvVar,
					}},
				},
			},
		},
	}

	// Making MyAppResource owner of the Deployment so that it can manage its lifecycle
	if err := controllerutil.SetControllerReference(&myAppResource, deploy, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference for Deployment")
		return err
	}

	log.Info("Creating Deployment")
	if err := r.Create(ctx, deploy); err != nil {
		log.Error(err, "Failed to create Deployment")
		return err
	}
	return nil
}

func (r *MyAppResourceReconciler) updateDeployment(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource, deploy appsv1.Deployment) error {
	log := log.FromContext(ctx)

	updateFlag := false

	cacheServer := fmt.Sprintf("tcp://%s-redis:6379", myAppResource.Name)

	cacheServerFlag := false
	for _, envVar := range deploy.Spec.Template.Spec.Containers[0].Env {
		if envVar.Name == "PODINFO_CACHE_SERVER" {
			cacheServerFlag = true
		}
	}

	if myAppResource.Spec.Redis.Enabled != cacheServerFlag {
		podEnvVar := []corev1.EnvVar{
			{Name: "PODINFO_UI_COLOR", Value: myAppResource.Spec.UI.Color},
			{Name: "PODINFO_UI_MESSAGE", Value: myAppResource.Spec.UI.Message},
		}

		if myAppResource.Spec.Redis.Enabled {
			podEnvVar = []corev1.EnvVar{
				{Name: "PODINFO_UI_COLOR", Value: myAppResource.Spec.UI.Color},
				{Name: "PODINFO_UI_MESSAGE", Value: myAppResource.Spec.UI.Message},
				{Name: "PODINFO_CACHE_SERVER", Value: cacheServer},
			}
		}

		deploy.Spec.Template.Spec.Containers[0].Env = podEnvVar

		updateFlag = true
	}

	desiredReplicas := myAppResource.Spec.Replicas
	if *deploy.Spec.Replicas != desiredReplicas {
		log.Info("Updating Deployment replicas", "current", *deploy.Spec.Replicas, "desired", desiredReplicas)
		deploy.Spec.Replicas = &desiredReplicas

		updateFlag = true
	}

	desiredMemoryLimit, _ := resource.ParseQuantity(myAppResource.Spec.Resources.MemoryLimit)
	if deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] != desiredMemoryLimit {
		log.Info("Updating MemoryLimit")
		deploy.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] = desiredMemoryLimit

		updateFlag = true
	}

	desiredCpuRequest, _ := resource.ParseQuantity(myAppResource.Spec.Resources.CPURequest)
	if deploy.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU] != desiredCpuRequest {
		log.Info("Updating CpuRequest")
		deploy.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU] = desiredCpuRequest

		updateFlag = true
	}

	desiredImage := fmt.Sprintf("%v:%v", myAppResource.Spec.Image.Repository, myAppResource.Spec.Image.Tag)
	if deploy.Spec.Template.Spec.Containers[0].Image != desiredImage {
		log.Info("Updating Image")
		deploy.Spec.Template.Spec.Containers[0].Image = desiredImage

		updateFlag = true
	}

	desiredPodUIColor := myAppResource.Spec.UI.Color
	if deploy.Spec.Template.Spec.Containers[0].Env[0].Value != desiredPodUIColor {
		log.Info("Updating POD UI Color")
		deploy.Spec.Template.Spec.Containers[0].Env[0].Value = desiredPodUIColor

		updateFlag = true
	}

	desiredPodUIMessage := myAppResource.Spec.UI.Message
	if deploy.Spec.Template.Spec.Containers[0].Env[1].Value != desiredPodUIMessage {
		log.Info("Updating POD UI Message")
		deploy.Spec.Template.Spec.Containers[0].Env[1].Value = desiredPodUIMessage

		updateFlag = true
	}

	if updateFlag {
		if err := r.Update(ctx, &deploy); err != nil {
			log.Error(err, "Failed to update Deployment")
			return err
		}
	}

	return nil
}

func (r *MyAppResourceReconciler) reconcileService(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	// log := log.FromContext(ctx)

	mySvc := &corev1.Service{}

	err := r.Get(ctx, types.NamespacedName{Name: myAppResource.Name, Namespace: myAppResource.Namespace}, mySvc)
	if err != nil {
		if errors.IsNotFound(err) {
			return r.createService(ctx, myAppResource)
		}

		return err
	}

	return nil
}

func (r *MyAppResourceReconciler) createService(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	log := log.FromContext(ctx)

	mySvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myAppResource.Name,
			Namespace: myAppResource.Namespace,
			Labels:    map[string]string{"app.kubernetes.io/name": myAppResource.Name},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       9898,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{"app": myAppResource.Name},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	// Making MyAppResource owner of the Service so that it can manage its lifecycle
	if err := controllerutil.SetControllerReference(&myAppResource, mySvc, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference for Service")
		return err
	}

	log.Info("Creating Service")
	if err := r.Create(ctx, mySvc); err != nil {
		log.Error(err, "Failed to create Service")
		return err
	}
	return nil
}

func (r *MyAppResourceReconciler) reconcileRedis(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	log := log.FromContext(ctx)

	redisName := fmt.Sprintf("%s-redis", myAppResource.Name)

	redisDeploy := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: redisName, Namespace: myAppResource.Namespace}, redisDeploy)
	if err != nil {
		if errors.IsNotFound(err) {
			if myAppResource.Spec.Redis.Enabled {
				r.createRedisDeployment(ctx, myAppResource)
			}
		}

		return err
	}

	redisService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: redisName, Namespace: myAppResource.Namespace}, redisService)
	if err != nil {
		if errors.IsNotFound(err) {
			if myAppResource.Spec.Redis.Enabled {
				r.createRedisService(ctx, myAppResource)
			}
		}

		return err
	}

	if !myAppResource.Spec.Redis.Enabled {
		if err := r.Delete(ctx, redisDeploy); err != nil {
			log.Error(err, "Failed to delete Redis Deployment")
			return err
		}

		if err := r.Delete(ctx, redisService); err != nil {
			log.Error(err, "Failed to delete Redis Service")
			return err
		}
	}

	return nil
}

func (r *MyAppResourceReconciler) createRedisDeployment(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	log := log.FromContext(ctx)

	replicas := int32(1)
	redisName := fmt.Sprintf("%s-redis", myAppResource.Name)

	redisDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisName,
			Namespace: myAppResource.Namespace,
			Labels:    map[string]string{"app.kubernetes.io/name": myAppResource.Name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": redisName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app.kubernetes.io/name": myAppResource.Name, "app": redisName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "redis:7.0.7",
						Name:            redisName,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 6379,
								Protocol:      corev1.ProtocolTCP,
								Name:          "redis",
							},
						},
					}},
				},
			},
		},
	}

	// Making MyAppResource owner of the Redis Deployment so that it can manage its lifecycle
	if err := controllerutil.SetControllerReference(&myAppResource, redisDeploy, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference for Redis Deployment")
		return err
	}

	log.Info("Creating Redis Deployment")
	if err := r.Create(ctx, redisDeploy); err != nil {
		log.Error(err, "Failed to create Redis Deployment")
		return err
	}

	return nil
}

func (r *MyAppResourceReconciler) createRedisService(ctx context.Context, myAppResource k8sv1alpha1.MyAppResource) error {
	log := log.FromContext(ctx)

	redisName := fmt.Sprintf("%s-redis", myAppResource.Name)

	redisSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisName,
			Namespace: myAppResource.Namespace,
			Labels:    map[string]string{"app.kubernetes.io/name": myAppResource.Name},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "redis",
					Port:       6379,
					TargetPort: intstr.FromString("redis"),
				},
			},
			Selector: map[string]string{"app": redisName},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	// Making MyAppResource owner of the Redis Service so that it can manage its lifecycle
	if err := controllerutil.SetControllerReference(&myAppResource, redisSvc, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference for Redis Service")
		return err
	}

	log.Info("Creating Redis Service")
	if err := r.Create(ctx, redisSvc); err != nil {
		log.Error(err, "Failed to create Redis Service")
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.MyAppResource{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
