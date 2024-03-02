package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	//"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// parseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func parseFlags() string {
	var namespace string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&namespace, "namespace", "events", "namespace to be used for this utility")

	// Actually parse the flags
	flag.Parse()

	// Return the configuration path
	return namespace
}

func main() {
	ns := parseFlags()
	logf.SetLogger(zap.New())
	log.Println("Starting event generator")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error creating client config: %v", err)
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Fatalf("Error creating manager: %v", err)
	}

	recorder := mgr.GetEventRecorderFor("event-generator")

	for {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx-pod",
				Namespace: ns,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "nginx-container",
						Image: "nginx:latest",
					},
				},
			},
		}

		// Create the pod using client
		err = mgr.GetClient().Create(context.Background(), pod)
		if err != nil {
			log.Fatalf("Error creating pod: %v", err)
		}

		log.Printf("Pod created: %s/%s", pod.Namespace, pod.Name)

		for count := 1; count < 26; count++ {
			log.Println(fmt.Sprintf("Event count: %d", count))
			recorder.Event(pod, v1.EventTypeNormal, fmt.Sprintf("Event count: %d", count), "The pod is now running")
		}

		err = mgr.GetClient().Delete(context.Background(), pod)
		if err != nil {
			log.Fatalf("Error deleting pod: %v", err)
		}
	}
}
