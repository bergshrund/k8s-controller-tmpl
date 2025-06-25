package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/rs/zerolog/log"
)

func SetupEnvTest(t *testing.T) (*envtest.Environment, *kubernetes.Clientset, func()) {

	ctx := context.Background()
	env := &envtest.Environment{}

	cfg, err := env.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	clientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("test-deployment-%d", i),
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "test"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nginx",
								Image: "nginx:latest",
							},
						},
					},
				},
			},
		}

		_, err := clientset.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	cleanup := func() {
		err := env.Stop()
		if err != nil {
			log.Error().Err(err).Msg("Failed to stop envtest environment. Please check if resources were cleaned up properly.")
		}
	}

	return env, clientset, cleanup
}

func int32Ptr(i int32) *int32 { return &i }
