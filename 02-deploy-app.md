# Chapter 2 - Deploy the Application to Kubernetes

## Overview

In this chapter, you will learn some intermediate Pulumi concepts:

- [Stack References](https://www.pulumi.com/docs/intro/concepts/stack/#stackreferences) to share outputs
  between the different stacks.
- Create programmatically a Kubernetes provider using the `Provider` resource. This is useful if you want to use a
  different Kubernetes provider than the one you used to create the cluster.
- Deploy a different Kubernetes resources to the cluster, without the use of YAML files.
- The usage of the `dependsOn` property, to make sure the Kubernetes resources are created in the correct order.

## Prerequisites

- The Kubernetes cluster from the [previous chapter](/00-cluster-setup.md)
- The nodejs application from the [previous chapter](/01-app-setup.md)
- Pulumi CLI installed

## Instructions

You may have noticed, for this chapter I am going to use a different Pulumi supported language. My choice is Go, but
feel free to use the language you are most comfortable with.

### Step 1 - Change into the `02-deploy-app` directory

Use the `cd` command to change into the `02-deploy-app` directory.

```bash
cd 02-deploy-app
```

### Step 2 - Get the Kubernetes cluster outputs and container image

To retrieve the outputs of the different stacks, we use `StackReference`s. Please change the actual stack names to the
ones you used in the previous chapters.

```go
package main
..
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		infraStackRef, err := pulumi.NewStackReference(ctx, "dirien/devopsdaysams/dev", nil)
		if err != nil {
			return err
		}
		appImageRef, err := pulumi.NewStackReference(ctx, "dirien/devopsdaysams-app/dev", nil)
		if err != nil {
			return err
		}
		...
	}
}

```

### Step 3 - Deploy the application

Now we can deploy the application to the cluster. Run `pulumi up` to deploy the application.

```bash
pulumi up
```

## Stretch Goals

- Can you deploy a `namespace` via Pulumi to the cluster and add the `deployment` and `service` to this namespace?
- What Pulumi Resource you would use to deploy a Helm Chart? Can you deploy the `ingress-nginx` Helm Chart to the
  cluster?

## Learn More

- [Pulumi](https://www.pulumi.com/)
- [Kubernetes Pulumi Provider](https://www.pulumi.com/registry/packages/kubernetes/)
