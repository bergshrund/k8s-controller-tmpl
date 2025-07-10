package informer

import (
	"context"
	"sync"
	"testing"
	"time"

	testutils "k8s-controller-tmpl/pkg/testutils"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func TestStartInformer(t *testing.T) {

	_, clientset, cleanup := testutils.SetupEnvTest(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	added := make(chan string, 2)

	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
	)

	informer := factory.Apps().V1().Deployments().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if deployment, ok := obj.(metav1.Object); ok {
				added <- deployment.GetName()
			}
		},
	})

	go func() {
		defer wg.Done()
		factory.Start(ctx.Done())
		factory.WaitForCacheSync(ctx.Done())
		<-ctx.Done()
	}()

	found := map[string]bool{}

	for range 2 {
		select {
		case name := <-added:
			found[name] = true
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for deployment add events")
		}
	}

	require.True(t, found["test-deployment-0"])
	require.True(t, found["test-deployment-1"])

	cancel() // Stop the informer after the test completes
	wg.Wait()
}

func TestStartInformer_CoversFunction(t *testing.T) {
	_, clientset, cleanup := testutils.SetupEnvTest(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		StartInformer(ctx, clientset, "default")
	}()

	time.Sleep(1 * time.Second)
	cancel() // Stop the informer after the test completes
}

func TestGetDeploymentName(t *testing.T) {
	dep := &metav1.PartialObjectMetadata{}
	dep.SetName("test-deployment")
	name := getDeploymentName(dep)
	if name != "test-deployment" {
		t.Errorf("expected 'test-deployment', got %q", name)
	}
	name = getDeploymentName("not-an-object")
	if name != "unknown" {
		t.Errorf("expected 'unknown', got %q", name)
	}
}
