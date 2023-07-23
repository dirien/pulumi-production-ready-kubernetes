# Chapter 3 - Setup FluxCD

## Overview

In this chapter, we're going to setup FluxCD to deploy a simple application to the Kubernetes cluster. The application
has a Helm chart. We also will use `StackReference` to reference again, the get the outputs of our infrastructure stack
and also create programmatically a Kubernetes provider again.

New concepts in this chapter:

- The usage of the `Release` resource, to deploy a Helm chart
- The definition of a `CustomResource` to define a FluxCD `GitRepository` and `HelmRelease`

As programming language, I am going to use Go again, but feel free to use the language you are most comfortable with.

## Prerequisites

- The Kubernetes cluster from the [previous chapter](/00-cluster-setup.md)
- Pulumi CLI installed
- [Go](https://golang.org/doc/install)
- Fork the https://github.com/dirien/helloserver repository, if you want to change the application
- Optional: Install the [FluxCD CLI](https://fluxcd.io/flux/installation/#install-the-flux-cli)

## Instructions

### Step 1 - Change into the `03-fluxcd-setup` directory

```bash
cd 03-fluxcd-setup
# Install go dependencies
go mod download
```

## Step 2 - Get the Kubernetes cluster outputs

To retrieve the output from our cluster deployment, we use `StackReference`. Please change the actual stack names to the
ones you used in the previous chapters.

```bash
pulumi config set infraStackRef
```

Pulumi will ask you now to create a new stack. You can name the stack whatever you want. If you run Pulumi with the
local login, please make sure to use for every stack a different name.

```bash
Please choose a stack, or create a new one:  [Use arrows to move, type to filter]
> <create a new stack>
Please choose a stack, or create a new one: <create a new stack>
Please enter your desired stack name: deploy   
```

### Step 3 - Deploy the stack

> **Note:** If you run Pulumi for the first time, you will be asked to log in. Follow the instructions on the screen to
> login. You may need to create an account first, don't worry it is free.
> Alternatively you can use also the `pulumi login --local` command to login locally.

Change in the `Pulumi.yaml` the `gitrepo` section to match your forked repository.

Run `pulumi up` to deploy the stack.

```bash
pulumi up
```

### Step 4 - Change the default value of the Helm chart

If you forked the `https://github.com/dirien/helloserver`, you can change the default `tag` value in `
/helloserver/delivery/charts/hello-server/values.yaml` to `v0.1.1` and commit the change.

Now the GitOps pipeline should deploy the new version of the application, you can either wait for the next sync or use
the FluxCD CLI to trigger a sync.

```bash
fluc reconcile gitrepository hello[main.go](04-idp%2Fmain.go)-server-git-repo --kubeconfig kubeconfig -n default
flux reconcile helmrelease hello-server-helm-release --kubeconfig kubeconfig -n default
```

## Stretch Goals

- Can you deploy a second helm chart to the cluster? But this time create a `HelmRepository` resource to deploy the
  chart from a remote repository.
- Can you define a dependency between the two `HelmRelease` resources using the `dependsOn` property from `HelmRelease`?

## Learn More

- [Pulumi](https://www.pulumi.com/)
- [Kubernetes Pulumi Provider](https://www.pulumi.com/registry/packages/kubernetes/)
- [FluxCD](https://fluxcd.io/flux/)
