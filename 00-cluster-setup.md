# Chapter 0 - Create a Kubernetes Cluster

## Cloud Providers used in this workshop

- [Scaleway](https://www.scaleway.com/en/)

## Overview

In order to set up a GitOps workflow, we are going to need a Kubernetes Cluster. The goal of this chapter is to firstly
create a Kubernetes Cluster using Pulumi.

We are not only going to create a Scaleway Kubernetes Cluster, we will create a Scaleway Container Registry too.

### Modern Infrastructure As Code with Pulumi

Pulumi is an open-source infrastructure-as-code tool for creating, deploying and managing cloud
infrastructure. Pulumi works with traditional infrastructures like VMs, networks, and databases and modern
architectures, including containers, Kubernetes clusters, and serverless functions. Pulumi supports dozens of public,
private, and hybrid cloud service providers.

Pulumi is a multi-language infrastructure as Code tool using imperative languages to create a declarative
infrastructure description.

You have a wide range of programming languages available, and you can use the one you and your team are the most
comfortable with. Currently, (6/2023) Pulumi supports the following languages:

* Node.js (JavaScript / TypeScript)

* Python

* Go

* Java

* .NET (C#, VB, F#)

* YAML

The workshop examples are written in `typescript` and `Go`, but feel free to use the language you are most comfortable
with.

## Prerequisites

To successful complete this chapter, you must meet all of these requirements:

- [Scaleway account](https://console.scaleway.com/register)
- [Scaleway CLI](https://www.scaleway.com/en/cli/)
- [The Pulumi CLI](https://www.pulumi.com/docs/get-started/install/) should be present on your machine
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

## Instructions

### Step 1 - Clone the repo

Go to GitHub and fork/clone the [Production Ready Kubernetes Workshop](pulumi-production-ready-kubernetes)
repo and then change into the directory.

If you use SSH to clone:

```bash
git clone git@github.com:dirien/pulumi-production-ready-kubernetes.git
cd pulumi-production-ready-kubernetes
```

To clone with HTTP:

```bash
git clone https://github.com/dirien/pulumi-production-ready-kubernetes.git
cd pulumi-production-ready-kubernetes
```

### Step 2 - Configure Scaleway CLI

1. Run `scw init` and follow the instructions to configure the Scaleway CLI, you will get all the informations you need
   from the workshop host.
1. Test your configuration by running following command:

```bash
scw k8s cluster list
ID   NAME  STATUS  VERSION  REGION  PROJECT ID  TAGS  CNI  DESCRIPTION  CLUSTER URL  CREATED AT  UPDATED AT  TYPE
```

### Step 3 - Set up Pulumi

Change into the `00-cluster-setup` directory.

```bash
cd 00-cluster-setup
```

Most important part of a Pulumi program is the `Pulumi.yaml`. Here you can define and modify various settings. From
the runtime of the programming language you are using to changing the default config values.

- Change the region in the `Pulumi.yaml` file to your preferred region
- Change the zone in the `Pulumi.yaml` file to your preferred zone
- Change the node type in the `Pulumi.yaml` file to your preferred node type
- And much more options to configure

```yaml
...
config:
  scaleway:region: "nl-ams"
  scaleway:zone: "nl-ams-1"
  cluster:version: "1.27"
  cluster:auto_upgrade: true
  node:node_type: "PLAY2-NANO"
  node:auto_scale: false
  node:node_count: 3
  node:auto_heal: true
```

### Step 4 - Run Pulumi Up

```bash
pulumi up
```

If the preview looks good, select `yes` to deploy the cluster

```bash
Previewing update (dev)

View in Browser (Ctrl+O): https://app.pulumi.com/dirien/devopsdaysams/dev/previews/aaf2d43a-2ee9-4d8c-9b8c-c66d211dbd21

     Type                                  Name                     Plan       
 +   pulumi:pulumi:Stack                   devopsdaysams-dev        create     
 +   ├─ scaleway:index:KubernetesCluster   devopsdaysams-cluster    create     
 +   └─ scaleway:index:KubernetesNodePool  devopsdaysams-node-pool  create     


Outputs:
    kapsuleName: "devopsdaysams-cluster-70aa1ed"
    kubeconfig : output<string>
    region     : output<string>

Resources:
    + 3 to create

Do you want to perform this update?  [Use arrows to move, type to filter]
  yes
> no
  details
  [experimental] yes, using Update Plans (https://pulumi.com/updateplans)
```

If the deployment is successful, you should see the following output. The duration of the deployment can take a few
minutes.

```bash
...
Resources:
    + 3 created

Duration: 5m6s
```

### Step 5 - Configure Kubectl

With the `pulumi stack output` command, you can retrieve any output value from the stack. In this case, we are going to
retrieve the kubeconfig to use with `kubectl`.

```bash
pulumi stack output kubeconfig --show-secrets > kubeconfig 
```

### Step 6 - Verify the cluster

Now that we have the kubeconfig, we can verify the cluster is up and running. Not that we need this, but it is always
good to verify.

```bash
kubectl --kubeconfig kubeconfig get nodes
```

You should see a similar output:

```bash
NAME                                             STATUS   ROLES    AGE    VERSION
scw-devopsdaysams-cl-devopsdaysams-node-0dad28   Ready    <none>   111s   v1.27.2
scw-devopsdaysams-cl-devopsdaysams-node-11277b   Ready    <none>   114s   v1.27.2
scw-devopsdaysams-cl-devopsdaysams-node-516e39   Ready    <none>   111s   v1.27.2
```

Congratulations! You have successfully deployed a Kubernetes cluster on Scaleway using Pulumi. Please leave the cluster
up and running for [Chapter 1 - Containerize an Application](./01-app-setup.md)

## Stretch Goals

- Can you create a second node pool with a different node type? Add this node pool to the existing cluster.
- Can you add an admission plugin `AlwaysPullImages` to the cluster? You can check all the available plugins with
  following cli command `scw k8s version get <version`.

## Learn More

- [Pulumi](https://www.pulumi.com/)
- [Scaleway Pulumi Provider](https://www.pulumi.com/registry/packages/scaleway/)
