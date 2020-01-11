package cmd

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetDocker() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	gatewayContainerConfig := container.Config{
		Image: "Gateway",
		Env: []string{
			"Two env are yet to be added",
		},
	}
	gatewayContainerHostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: "/build",
			},
		},
	}
	if err := runContainer(gatewayContainerConfig, gatewayContainerHostConfig); err != nil {
		return err
	}

	configStoreContainerConfig := container.Config{
		Image: "Config-Store",
		Env: []string{
			"Two env are yet to be added",
		},
	}
	configStoreContainerHostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: "/build",
			},
		},
	}
	if err := runContainer(configStoreContainerConfig, configStoreContainerHostConfig); err != nil {
		return err
	}

	runnerContainerConfig := container.Config{
		Image: "Gateway",
		Env: []string{
			"Two env are yet to be added",
		},
	}
	runnerContainerHostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: "/build",
			},
		},
	}
	if err := runContainer(runnerContainerConfig, runnerContainerHostConfig); err != nil {
		return err
	}

	return nil
}

func runContainer(containerConfig container.Config, containerHostConfig container.HostConfig) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	resp, err := cli.ContainerCreate(ctx,
		&containerConfig,
		&containerHostConfig, nil, "")
	if err != nil {
		return err
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func SetKubernetes() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
