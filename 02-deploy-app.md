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
- [Go](https://golang.org/doc/install)

## Instructions

You may have noticed, for this chapter I am going to use a different Pulumi supported language. My choice is Go, but
feel free to use the language you are most comfortable with.

### Step 1 - Change into the `02-deploy-app` directory

Use the `cd` command to change into the `02-deploy-app` directory.

```bash
cd 02-deploy-app
# Install go dependencies
go mod download
```

### Step 2 - Get the Kubernetes cluster outputs and container image

To retrieve the outputs of the different stacks, we use `StackReference`s. Please change the actual stack names to the
ones you used in the previous chapters.

```bash
pulumi config set infraStackRef
pulumi config set appImageRef
```

Pulumi will ask you now to create a new stack. You can name the stack whatever you want. If you run Pulumi with the
local login, please make sure to use for every stack a different name.

```bash
Please choose a stack, or create a new one:  [Use arrows to move, type to filter]
> <create a new stack>
Please choose a stack, or create a new one: <create a new stack>
Please enter your desired stack name: deploy   
```


### Step 3 - Deploy the application

> **Note:** If you run Pulumi for the first time, you will be asked to log in. Follow the instructions on the screen to
> login. You may need to create an account first, don't worry it is free.
> Alternatively you can use also the `pulumi login --local` command to login locally.

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
