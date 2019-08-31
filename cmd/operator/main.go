// Copyright 2019 The Liftbridge Operator Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teivah/onecontext"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"

	"github.com/liftbridge-io/liftbridge-operator/pkg/version"
)

const (
	appName      = "liftbridge-operator"
	appNamespace = appName
)

func main() {
	// Parse command-line flags.
	pathToKubeconfig := flag.String("path-to-kubeconfig", "", "the path to the kubeconfig file to use")
	logLevel := flag.String("log-level", log.InfoLevel.String(), "the log level to use")
	flag.Parse()

	// Create a root context.
	rootContext := rootContext()

	// Configure logging.
	if v, err := log.ParseLevel(*logLevel); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	} else {
		log.SetLevel(v)
	}

	// Create a Kubernetes client.
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", *pathToKubeconfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Print useful information about Liftbridge Operator and the Kubernetes cluster.
	if v, err := kubeClient.Discovery().ServerVersion(); err != nil {
		log.Fatalf("Failed to check Kubernetes version: %v", err)
	} else {
		log.Infof("Liftbridge Operator %s is starting (Kubernetes version: %s)", version.Version, v.String())
	}

	// Attempt to retrieve our identity.
	id, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	// Create an event recorder.
	eb := record.NewBroadcaster()
	eb.StartLogging(log.Tracef)
	eb.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: kubeClient.CoreV1().Events(""),
	})
	er := eb.NewRecorder(scheme.Scheme, corev1.EventSource{
		Component: appName,
	})

	// Setup a lease to be used for leader election.
	rl, _ := resourcelock.New(
		resourcelock.LeasesResourceLock,
		appNamespace,
		appName,
		kubeClient.CoreV1(),
		kubeClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: er,
		},
	)

	// Perform leader election so that at most a single instance of Liftbridge Operator is active at any given moment.
	leaderelection.RunOrDie(context.Background(), leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				log.Debugf("Started leading")
				ctx, fn := onecontext.Merge(rootContext, ctx)
				defer fn()
				run(ctx)
			},
			OnStoppedLeading: func() {
				log.Fatalf("Leader election lost")
			},
			OnNewLeader: func(identity string) {
				log.Debugf("Current leader: %s", identity)
			},
		},
	})
}

// rootContext returns a context meant to be used as the root context.
// It is canceled when a SIGINT or SIGTERM signal is received.
// If multiple signals are received, Liftbridge Operator exits with an error.
func rootContext() context.Context {
	ctx, fn := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		log.Info("Liftbridge Operator is stopping")
		fn()
		<-ch
		log.Fatal("Liftbridge Operator failed to shutdown gracefully")
	}()
	return ctx
}

func run(ctx context.Context) {
	// Wait for the context to be canceled.
	<-ctx.Done()
	// The shutdown process is now complete.
	log.Info("Liftbridge Operator has stopped")
	// The leader election process started a background goroutine that is trying to renew the leader election lock.
	// Hence, we must manually exit now that we know the shutdown process is complete.
	os.Exit(0)
}
