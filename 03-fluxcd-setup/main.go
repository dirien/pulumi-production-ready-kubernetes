package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apiextensions"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	config "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		infraStackRef, err := pulumi.NewStackReference(ctx, "dirien/devopsdaysams/dev", nil)
		if err != nil {
			return err
		}

		k8sProvider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
			Kubeconfig:            infraStackRef.GetStringOutput(pulumi.String("kubeconfig")),
			EnableServerSideApply: pulumi.Bool(true),
		})

		release, err := helm.NewRelease(ctx, "flux2", &helm.ReleaseArgs{
			Chart:           pulumi.String("flux2"),
			Namespace:       pulumi.String("flux-system"),
			CreateNamespace: pulumi.Bool(true),
			RepositoryOpts: &helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://fluxcd-community.github.io/helm-charts"),
			},
			Version: pulumi.String("2.8.0"),
		}, pulumi.Provider(k8sProvider))

		if err != nil {
			return err
		}

		gitRepo, err := apiextensions.NewCustomResource(ctx, "hello-server-repo", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("source.toolkit.fluxcd.io/v1"),
			Kind:       pulumi.String("GitRepository"),
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("hello-server-git-repo"),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": pulumi.Map{
					"url": pulumi.String(config.Get(ctx, "gitrepo")),
					"ref": pulumi.Map{
						"branch": pulumi.String("main"),
					},
					"interval": pulumi.String("1m"),
				},
			},
		}, pulumi.Provider(k8sProvider), pulumi.DependsOn([]pulumi.Resource{release}))
		if err != nil {
			return err
		}

		_, err = apiextensions.NewCustomResource(ctx, "hello-server-helm-release", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("helm.toolkit.fluxcd.io/v2beta1"),
			Kind:       pulumi.String("HelmRelease"),
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("hello-server-helm-release"),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": pulumi.Map{
					"interval": pulumi.String("1m"),
					"chart": pulumi.Map{
						"spec": pulumi.Map{
							"chart":    pulumi.String("./delivery/charts/hello-server"),
							"interval": pulumi.String("1m"),
							"sourceRef": pulumi.Map{
								"kind": gitRepo.Kind,
								"name": gitRepo.Metadata.Name(),
							},
						},
					},
					"values": pulumi.Map{
						"service": pulumi.Map{
							"type": pulumi.String("LoadBalancer"),
						},
					},
				},
			},
		}, pulumi.Provider(k8sProvider), pulumi.DependsOn([]pulumi.Resource{release, gitRepo}))
		if err != nil {
			return err
		}
		return nil
	})
}
