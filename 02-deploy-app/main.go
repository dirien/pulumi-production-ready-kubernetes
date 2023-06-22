package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		infraStackRef, err := pulumi.NewStackReference(ctx, config.Get(ctx, "infraStackRef"), nil)
		if err != nil {
			return err
		}
		appImageRef, err := pulumi.NewStackReference(ctx, config.Get(ctx, "appImageRef"), nil)
		if err != nil {
			return err
		}
		k8sProvider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
			Kubeconfig:            infraStackRef.GetStringOutput(pulumi.String("kubeconfig")),
			EnableServerSideApply: pulumi.Bool(true),
		})

		appLabels := pulumi.StringMap{
			"app": pulumi.String("devopsdaysams"),
		}
		deployment, err := appsv1.NewDeployment(ctx, "app-dep", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Annotations: pulumi.StringMap{
					"pulumi.com/skipAwait": pulumi.String("true"),
				},
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("devopsdaysams"),
								Image: appImageRef.GetStringOutput(pulumi.String("imageName")),
								Ports: corev1.ContainerPortArray{
									corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(3000),
										Name:          pulumi.String("http"),
									},
								},
							}},
					},
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		service, err := corev1.NewService(ctx, "app-svc", &corev1.ServiceArgs{
			Spec: &corev1.ServiceSpecArgs{
				Selector: appLabels,
				Type:     corev1.ServiceSpecTypeLoadBalancer,
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(3000),
					},
				},
			},
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		ctx.Export("name", deployment.Metadata.Elem().Name())
		ctx.Export("url", service.Status.LoadBalancer().ToLoadBalancerStatusPtrOutput())

		return nil
	})
}
