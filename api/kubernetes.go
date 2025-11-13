package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
	namespace string
}

type Function struct {
	Name        string            `json:"name"`
	Runtime     string            `json:"runtime"`
	Handler     string            `json:"handler"`
	Code        string            `json:"code"`
	Environment map[string]string `json:"environment,omitempty"`
	MinReplicas int32             `json:"minReplicas,omitempty"`
	MaxReplicas int32             `json:"maxReplicas,omitempty"`
	Triggers    []Trigger         `json:"triggers,omitempty"`
	Status      FunctionStatus    `json:"status,omitempty"`
}

type Trigger struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type FunctionStatus struct {
	State          string    `json:"state"`
	Endpoint       string    `json:"endpoint,omitempty"`
	Replicas       int32     `json:"replicas"`
	LastDeployment time.Time `json:"lastDeployment,omitempty"`
}

type FunctionMetrics struct {
	Invocations   int64   `json:"invocations"`
	ColdStarts    int64   `json:"coldStarts"`
	AvgDuration   float64 `json:"avgDuration"`
	ErrorRate     float64 `json:"errorRate"`
	CostEstimate  float64 `json:"costEstimate"`
}

func NewKubernetesClient() (*KubernetesClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	namespace := "kube-serverless"
	// Try to get namespace from environment
	if ns := corev1.Namespace(); ns != "" {
		namespace = ns
	}

	return &KubernetesClient{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func (k *KubernetesClient) Ping() error {
	_, err := k.clientset.CoreV1().Namespaces().Get(context.Background(), k.namespace, metav1.GetOptions{})
	return err
}

func (k *KubernetesClient) ListFunctions(ctx context.Context) ([]Function, error) {
	deployments, err := k.clientset.AppsV1().Deployments(k.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=kube-serverless",
	})
	if err != nil {
		return nil, err
	}

	functions := make([]Function, 0, len(deployments.Items))
	for _, dep := range deployments.Items {
		function := k.deploymentToFunction(&dep)
		functions = append(functions, function)
	}

	return functions, nil
}

func (k *KubernetesClient) CreateFunction(ctx context.Context, fn *Function) error {
	// Set defaults
	if fn.MinReplicas == 0 {
		fn.MinReplicas = 0
	}
	if fn.MaxReplicas == 0 {
		fn.MaxReplicas = 10
	}

	// Create ConfigMap for function code
	if err := k.createFunctionConfigMap(ctx, fn); err != nil {
		return err
	}

	// Create Deployment
	if err := k.createFunctionDeployment(ctx, fn); err != nil {
		return err
	}

	// Create Service
	if err := k.createFunctionService(ctx, fn); err != nil {
		return err
	}

	// Create HPA for auto-scaling
	if err := k.createFunctionHPA(ctx, fn); err != nil {
		return err
	}

	return nil
}

func (k *KubernetesClient) GetFunction(ctx context.Context, name string) (*Function, error) {
	deployment, err := k.clientset.AppsV1().Deployments(k.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	function := k.deploymentToFunction(deployment)
	return &function, nil
}

func (k *KubernetesClient) UpdateFunction(ctx context.Context, fn *Function) error {
	// Update ConfigMap
	if err := k.updateFunctionConfigMap(ctx, fn); err != nil {
		return err
	}

	// Update Deployment
	deployment, err := k.clientset.AppsV1().Deployments(k.namespace).Get(ctx, fn.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deployment.Spec.Template.Spec.Containers[0].Image = k.getRuntimeImage(fn.Runtime)
	deployment.Spec.Template.Spec.Containers[0].Env = k.buildEnvVars(fn)

	_, err = k.clientset.AppsV1().Deployments(k.namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

func (k *KubernetesClient) DeleteFunction(ctx context.Context, name string) error {
	// Delete Deployment
	if err := k.clientset.AppsV1().Deployments(k.namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	// Delete Service
	if err := k.clientset.CoreV1().Services(k.namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	// Delete ConfigMap
	if err := k.clientset.CoreV1().ConfigMaps(k.namespace).Delete(ctx, name+"-code", metav1.DeleteOptions{}); err != nil {
		return err
	}

	// Delete HPA
	if err := k.clientset.AutoscalingV2().HorizontalPodAutoscalers(k.namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func (k *KubernetesClient) InvokeFunction(ctx context.Context, name string, body io.Reader) ([]byte, error) {
	// Get service endpoint
	service, err := k.clientset.CoreV1().Services(k.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// In a real implementation, we would make an HTTP request to the function
	// For now, return a placeholder response
	return []byte(fmt.Sprintf(`{"function": "%s", "status": "invoked"}`, name)), nil
}

func (k *KubernetesClient) GetFunctionMetrics(ctx context.Context, name string) (*FunctionMetrics, error) {
	// In a real implementation, this would query Prometheus
	// For now, return placeholder metrics
	return &FunctionMetrics{
		Invocations:  100,
		ColdStarts:   5,
		AvgDuration:  0.250,
		ErrorRate:    0.01,
		CostEstimate: 0.05,
	}, nil
}

func (k *KubernetesClient) createFunctionConfigMap(ctx context.Context, fn *Function) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fn.Name + "-code",
			Namespace: k.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       fn.Name,
				"app.kubernetes.io/managed-by": "kube-serverless",
			},
		},
		Data: map[string]string{
			"handler": fn.Handler,
			"code":    fn.Code,
		},
	}

	_, err := k.clientset.CoreV1().ConfigMaps(k.namespace).Create(ctx, cm, metav1.CreateOptions{})
	return err
}

func (k *KubernetesClient) updateFunctionConfigMap(ctx context.Context, fn *Function) error {
	cm, err := k.clientset.CoreV1().ConfigMaps(k.namespace).Get(ctx, fn.Name+"-code", metav1.GetOptions{})
	if err != nil {
		return err
	}

	cm.Data["handler"] = fn.Handler
	cm.Data["code"] = fn.Code

	_, err = k.clientset.CoreV1().ConfigMaps(k.namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}

func (k *KubernetesClient) createFunctionDeployment(ctx context.Context, fn *Function) error {
	replicas := fn.MinReplicas

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fn.Name,
			Namespace: k.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       fn.Name,
				"app.kubernetes.io/managed-by": "kube-serverless",
				"function":                     fn.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"function": fn.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"function":                     fn.Name,
						"app.kubernetes.io/name":       fn.Name,
						"app.kubernetes.io/managed-by": "kube-serverless",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "function",
							Image: k.getRuntimeImage(fn.Runtime),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Name:          "http",
								},
							},
							Env: k.buildEnvVars(fn),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "function-code",
									MountPath: "/function",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "function-code",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fn.Name + "-code",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := k.clientset.AppsV1().Deployments(k.namespace).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

func (k *KubernetesClient) createFunctionService(ctx context.Context, fn *Function) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fn.Name,
			Namespace: k.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       fn.Name,
				"app.kubernetes.io/managed-by": "kube-serverless",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
					Name:       "http",
				},
			},
			Selector: map[string]string{
				"function": fn.Name,
			},
		},
	}

	_, err := k.clientset.CoreV1().Services(k.namespace).Create(ctx, service, metav1.CreateOptions{})
	return err
}

func (k *KubernetesClient) createFunctionHPA(ctx context.Context, fn *Function) error {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fn.Name,
			Namespace: k.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       fn.Name,
				"app.kubernetes.io/managed-by": "kube-serverless",
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       fn.Name,
			},
			MinReplicas: &fn.MinReplicas,
			MaxReplicas: fn.MaxReplicas,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: int32Ptr(80),
						},
					},
				},
			},
		},
	}

	_, err := k.clientset.AutoscalingV2().HorizontalPodAutoscalers(k.namespace).Create(ctx, hpa, metav1.CreateOptions{})
	return err
}

func (k *KubernetesClient) getRuntimeImage(runtime string) string {
	runtimeImages := map[string]string{
		"nodejs18": "node:18-alpine",
		"python39": "python:3.9-alpine",
		"go119":    "golang:1.19-alpine",
	}

	if image, ok := runtimeImages[runtime]; ok {
		return image
	}

	return "node:18-alpine" // default
}

func (k *KubernetesClient) buildEnvVars(fn *Function) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name:  "FUNCTION_NAME",
			Value: fn.Name,
		},
		{
			Name:  "FUNCTION_HANDLER",
			Value: fn.Handler,
		},
		{
			Name:  "RUNTIME",
			Value: fn.Runtime,
		},
	}

	for key, value := range fn.Environment {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars
}

func (k *KubernetesClient) deploymentToFunction(dep *appsv1.Deployment) Function {
	function := Function{
		Name:    dep.Name,
		Runtime: k.getEnvVar(dep.Spec.Template.Spec.Containers[0].Env, "RUNTIME"),
		Handler: k.getEnvVar(dep.Spec.Template.Spec.Containers[0].Env, "FUNCTION_HANDLER"),
		Status: FunctionStatus{
			State:    "running",
			Replicas: dep.Status.ReadyReplicas,
		},
	}

	if dep.Spec.Replicas != nil {
		function.MinReplicas = *dep.Spec.Replicas
	}

	return function
}

func (k *KubernetesClient) getEnvVar(envVars []corev1.EnvVar, name string) string {
	for _, env := range envVars {
		if env.Name == name {
			return env.Value
		}
	}
	return ""
}

func int32Ptr(i int32) *int32 {
	return &i
}
