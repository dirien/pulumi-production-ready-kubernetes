package main

import (
	"fmt"
	"github.com/port-labs/pulumi-port/sdk/go/port"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"os"
	"strings"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		infraStackRef, err := pulumi.NewStackReference(ctx, config.Get(ctx, "infraStackRef"), nil)
		if err != nil {
			return err
		}

		kapsuleID := infraStackRef.GetStringOutput(pulumi.String("kapsuleID"))

		cloudProviderBlueprint, err := port.NewBlueprint(ctx, "cloud-provider-blueprint", &port.BlueprintArgs{
			Title:       pulumi.String("Cloud Provider"),
			Icon:        pulumi.String("Environment"),
			Identifier:  pulumi.String("cloud-provider-blueprint"),
			Description: pulumi.String("A blueprint to describe a cloud provider"),
			Properties: port.BlueprintPropertiesArgs{
				StringProps: port.BlueprintPropertiesStringPropsMap{
					"provider": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Provider"),
						Enums: pulumi.StringArray{
							pulumi.String("Scaleway"),
						},
						EnumColors: pulumi.StringMap{
							"Scaleway": pulumi.String("purple"),
						},
					},
					"type": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Type"),
						Enums: pulumi.StringArray{
							pulumi.String("Dev"),
							pulumi.String("Staging"),
							pulumi.String("Prod"),
						},
					},
					"region": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Region"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Some hardcoded cloud provider configurations
		for _, name := range []string{"Dev", "Staging", "Prod"} {
			_, err = port.NewEntity(ctx, fmt.Sprintf("cloud-provider-%s-entity", name), &port.EntityArgs{
				Blueprint:  cloudProviderBlueprint.ID(),
				Identifier: pulumi.String(fmt.Sprintf("cloud-provider-%s", strings.ToLower(name))),
				Title:      pulumi.String(name),
				Properties: port.EntityPropertiesArgs{
					StringProps: pulumi.StringMap{
						"provider": pulumi.String("Scaleway"),
						"type":     pulumi.String(name),
						"region":   pulumi.String("nl-ams"),
					},
				},
				Icon: pulumi.String("Environment"),
			})
			if err != nil {
				return err
			}
		}

		gitOpsBackendBlueprint, err := port.NewBlueprint(ctx, "gitops-backend-blueprint", &port.BlueprintArgs{
			Title:      pulumi.String("GitOps Backend"),
			Identifier: pulumi.String("gitops-backend-blueprint"),
			Icon:       pulumi.String("CICD"),
			Properties: port.BlueprintPropertiesArgs{
				StringProps: port.BlueprintPropertiesStringPropsMap{
					"backend": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("GitOps Backend"),
						Enums: pulumi.StringArray{
							pulumi.String("ArgoCD"),
							pulumi.String("FluxCD"),
						},
						EnumColors: pulumi.StringMap{
							"ArgoCD": pulumi.String("orange"),
							"FluxCD": pulumi.String("blue"),
						},
					},
					"version": port.BlueprintPropertiesStringPropsArgs{
						Title: pulumi.String("Version"),
					},
				},
			},
		})

		clusterBlueprint, err := port.NewBlueprint(ctx, "k8s-blueprint", &port.BlueprintArgs{
			Title:       pulumi.String("Kubernetes"),
			Icon:        pulumi.String("Cluster"),
			Identifier:  pulumi.String("k8s-blueprint"),
			Description: pulumi.String("A blueprint to describe a Kubernetes cluster"),
			Properties: port.BlueprintPropertiesArgs{
				StringProps: port.BlueprintPropertiesStringPropsMap{
					"name": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Name"),
					},
					"version": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Version"),
						Enums: pulumi.StringArray{
							pulumi.String("1.27"),
							pulumi.String("1.26"),
							pulumi.String("1.25"),
						},
					},
				},
				BooleanProps: port.BlueprintPropertiesBooleanPropsMap{
					"auto-upgrade": port.BlueprintPropertiesBooleanPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Auto upgrade"),
						Default:  pulumi.Bool(true),
					},
				},
				NumberProps: port.BlueprintPropertiesNumberPropsMap{
					"nodes": port.BlueprintPropertiesNumberPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Nodes"),
						Minimum:  pulumi.Float64(1),
						Maximum:  pulumi.Float64(10),
						Default:  pulumi.Float64(3),
					},
				},
			},
			MirrorProperties: port.BlueprintMirrorPropertiesMap{},
			CalculationProperties: port.BlueprintCalculationPropertiesMap{
				"cloud-provider-portal": port.BlueprintCalculationPropertiesArgs{
					Title:       pulumi.String("Cloud portal"),
					Calculation: pulumi.String("'https://console.scaleway.com/kubernetes/clusters/' + .identifier + '/overview'"),
					Type:        pulumi.String("string"),
					Icon:        pulumi.String("Link"),
					Format:      pulumi.String("url"),
				},
			},
			Relations: port.BlueprintRelationsMap{
				"runs-on": port.BlueprintRelationsArgs{
					Title:    pulumi.String("Runs on"),
					Target:   cloudProviderBlueprint.ID(),
					Required: pulumi.Bool(true),
					Many:     pulumi.Bool(false),
				},
				"gitops-backend": port.BlueprintRelationsArgs{
					Title:  pulumi.String("GitOps Backend"),
					Target: gitOpsBackendBlueprint.ID(),
					Many:   pulumi.Bool(false),
				},
			},
		})
		if err != nil {
			return err
		}

		clusterEntity, err := port.NewEntity(ctx, "k8s-entity", &port.EntityArgs{
			Blueprint:  clusterBlueprint.ID(),
			Identifier: kapsuleID,
			Icon:       pulumi.String("Cluster"),
			Title:      infraStackRef.GetStringOutput(pulumi.String("kapsuleName")),
			Properties: port.EntityPropertiesArgs{
				StringProps: pulumi.StringMap{
					"name":    infraStackRef.GetStringOutput(pulumi.String("kapsuleName")),
					"version": infraStackRef.GetStringOutput(pulumi.String("kapsuleVersion")),
				},
				BooleanProps: pulumi.BoolMap{
					"auto-upgrade": infraStackRef.GetOutput(pulumi.String("kapsuleAutoUpgrade")).AsBoolOutput(),
				},
				NumberProps: pulumi.Float64Map{
					"nodes": infraStackRef.GetOutput(pulumi.String("kapusuleNodeCount")).AsFloat64Output(),
				},
			},
			Relations: port.EntityRelationsArgs{
				SingleRelations: pulumi.StringMap{
					"runs-on": pulumi.String("cloud-provider-prod"),
				},
			},
		})

		gitopsPlatform, err := port.NewBlueprint(ctx, "gitops-platform-blueprint", &port.BlueprintArgs{
			Title: pulumi.String("GitOps Platform"),
			Icon:  pulumi.String("Home"),
			Properties: port.BlueprintPropertiesArgs{
				StringProps: port.BlueprintPropertiesStringPropsMap{
					"name": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Name"),
					},
				},
			},
			Identifier: pulumi.String("gitops-platform-blueprint"),
			MirrorProperties: port.BlueprintMirrorPropertiesMap{
				"platform-region": port.BlueprintMirrorPropertiesArgs{
					Title: pulumi.String("Region"),
					Path:  pulumi.String("consists-of.runs-on.region"),
				},
				"platform-type": port.BlueprintMirrorPropertiesArgs{
					Title: pulumi.String("Stage"),
					Path:  pulumi.String("consists-of.runs-on.type"),
				},
				"gitops-backend": port.BlueprintMirrorPropertiesArgs{
					Title: pulumi.String("GitOps Backend"),
					Path:  pulumi.String("consists-of.gitops-backend.backend"),
				},
			},
			Relations: port.BlueprintRelationsMap{
				"consists-of": port.BlueprintRelationsArgs{
					Title:    pulumi.String("Consists of"),
					Target:   clusterBlueprint.ID(),
					Required: pulumi.Bool(true),
					Many:     pulumi.Bool(false),
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = port.NewEntity(ctx, "gitops-platform-entity", &port.EntityArgs{
			Blueprint:  gitopsPlatform.ID(),
			Identifier: pulumi.String("gitops-platform-entity"),
			Title:      pulumi.String("GitOps Platform Scaleway"),
			Icon:       pulumi.String("Home"),
			Properties: port.EntityPropertiesArgs{
				StringProps: pulumi.StringMap{
					"name": pulumi.String("Scaleway"),
				},
			},
			Relations: port.EntityRelationsArgs{
				SingleRelations: pulumi.StringMap{
					"consists-of": clusterEntity.Identifier,
				},
			},
		})
		if err != nil {
			return err
		}

		helmReleaseBlueprint, err := port.NewBlueprint(ctx, "helm-release-blueprint", &port.BlueprintArgs{
			Title:      pulumi.String("Helm Release"),
			Icon:       pulumi.String("Package"),
			Identifier: pulumi.String("helm-release-blueprint"),
			Relations: port.BlueprintRelationsMap{
				"deployed-on": port.BlueprintRelationsArgs{
					Title:  pulumi.String("Deployed on"),
					Target: clusterBlueprint.ID(),
					Many:   pulumi.Bool(false),
				},
				"gitops-platform": port.BlueprintRelationsArgs{
					Title:  pulumi.String("GitOps Platform"),
					Target: gitopsPlatform.ID(),
					Many:   pulumi.Bool(false),
				},
			},
			MirrorProperties: port.BlueprintMirrorPropertiesMap{
				"clusterName": port.BlueprintMirrorPropertiesArgs{
					Title: pulumi.String("Cluster Name"),
					Path:  pulumi.String("deployed-on.name"),
				},
			},
			Properties: port.BlueprintPropertiesArgs{
				NumberProps: port.BlueprintPropertiesNumberPropsMap{
					"releaseRevision": port.BlueprintPropertiesNumberPropsArgs{
						Title:   pulumi.String("Release Revision"),
						Default: pulumi.Float64(0),
					},
				},
				StringProps: port.BlueprintPropertiesStringPropsMap{
					"name": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Name"),
					},
					"version": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Version"),
					},
					"chart": port.BlueprintPropertiesStringPropsArgs{
						Required: pulumi.Bool(true),
						Title:    pulumi.String("Chart"),
					},
					"namespace": port.BlueprintPropertiesStringPropsArgs{
						Title:   pulumi.String("Namespace"),
						Default: pulumi.String("default"),
					},
					"status": port.BlueprintPropertiesStringPropsArgs{
						Title:   pulumi.String("Status"),
						Default: pulumi.String("False"),
					},
					"message": port.BlueprintPropertiesStringPropsArgs{
						Title:   pulumi.String("Message"),
						Default: pulumi.String(""),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Deploy the HelmRelease and HelmRepository CRs to the cluster and let Flux reconcile them.
		k8sProvider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
			Kubeconfig:            infraStackRef.GetStringOutput(pulumi.String("kubeconfig")),
			EnableServerSideApply: pulumi.Bool(true),
		})

		portRepo, err := apiextensions.NewCustomResource(ctx, "port-repo", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("source.toolkit.fluxcd.io/v1beta2"),
			Kind:       pulumi.String("OCIRepository"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("port-labs"),
				Namespace: pulumi.String("flux-system"),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"interval": pulumi.String("1m"),
					"url":      pulumi.String("oci://ghcr.io/port-labs/charts/port-k8s-exporter"),
					"ref": kubernetes.UntypedArgs{
						"tag": pulumi.String("0.1.19"),
					},
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		k8sExporterNS, err := corev1.NewNamespace(ctx, "port-k8s-exporter", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("port-k8s-exporter"),
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		_, err = apiextensions.NewCustomResource(ctx, "port-repo-release", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("helm.toolkit.fluxcd.io/v2beta1"),
			Kind:       pulumi.String("HelmRelease"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("port-k8s-exporter"),
				Namespace: k8sExporterNS.Metadata.Name(),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": pulumi.Map{
					"install": pulumi.Map{
						"createNamespace": pulumi.Bool(true),
					},
					"interval": pulumi.String("1m"),
					"chart": pulumi.Map{
						"spec": pulumi.Map{
							"chart":    pulumi.String("port-k8s-exporter"),
							"interval": pulumi.String("1m"),
							"version":  pulumi.String("0.1.18"),
							"sourceRef": pulumi.Map{
								"kind":      portRepo.Kind,
								"name":      portRepo.Metadata.Name(),
								"namespace": portRepo.Metadata.Namespace(),
							},
						},
					},
					"values": pulumi.Map{
						"deleteDependents": pulumi.Bool(true),
						"secret": pulumi.Map{
							"secrets": pulumi.Map{
								"portClientId":     pulumi.String(os.Getenv("PORT_CLIENT_ID")),
								"portClientSecret": pulumi.ToSecret(pulumi.String(os.Getenv("PORT_CLIENT_SECRET"))),
							},
						},
						"configMap": pulumi.Map{
							"config": pulumi.Sprintf(`
resources:
- kind: apps/v1/Deployment
  selector:
    query: if (.metadata.labels."app.kubernetes.io/part-of") then (.metadata.labels."app.kubernetes.io/part-of" | startswith("flux")) else false end
  port:
   entity:
    mappings:
    - identifier: '"%s"'
      blueprint: '"k8s-blueprint"'
      relations:
       gitops-backend: .metadata.annotations."meta.helm.sh/release-name"
    - identifier: .metadata.annotations."meta.helm.sh/release-name"
      title: '"FluxCD"'
      blueprint: '"gitops-backend-blueprint"'
      properties:
       backend: '"FluxCD"'
       version: .metadata.labels."app.kubernetes.io/version"
- kind: helm.toolkit.fluxcd.io/v2beta1/HelmRelease
  port:
   entity:
    mappings:
    - identifier: .metadata.name
      blueprint: '"helm-release-blueprint"'
      properties:
        name: .metadata.name
        version: .status.lastAppliedRevision
        chart: (.status.helmChart | split("/") | last)
        releaseRevision: .status.lastReleaseRevision
        status: (.status.conditions | last | .status)
        message: (.status.conditions | last | .message)
      relations:
        deployed-on: '"%s"'
        gitops-platform: '"gitops-platform-entity"'`, kapsuleID, kapsuleID),
						},
					},
				},
			},
		}, pulumi.Provider(k8sProvider), pulumi.Provider(k8sProvider), pulumi.DependsOn([]pulumi.Resource{helmReleaseBlueprint, gitopsPlatform, clusterBlueprint, gitOpsBackendBlueprint, portRepo, k8sExporterNS}))
		if err != nil {
			return err
		}
		return nil
	})
}
