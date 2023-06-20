# Chapter 1 - Containerize an Application

## Overview

Now that we created a Kubernetes cluster, we can deploy our application to it. In this chapter we will containerize a
simple node.js application, push it to a container registry.

The Container Registry, we will use is from Scaleway too! The Pulumi code will create a registry for us, build the image
and push it to the registry.

## Prerequisites

- [Node.js](https://nodejs.org/en/download/) installed
- [Docker](https://docs.docker.com/get-docker/) installed
- Pulumi CLI installed

## Instructions

### Step 1 - Change into the `01-app` directory

Change into the `01-app` directory using:

```bash
cd 01-app
```

### Step 2 - Set the Scaleway secret key

Before we can run the Pulumi code, we need to pass the Scaleway secret key to Pulumi. We can do this by using the Pulumi
CLI:

```bash
pulumi config set secret_key <value> --secret
```

### Step 3 - Inspect the Pulumi code

One important thing to note is that we need pass the `platform` explicitly to the `docker.Image` resource. This is
because we may use a different platform than the one we are running on. In this case, I am running on a Mac, but I want
to build a Linux image.

### Step 4 - Build and push the image

Now we can run `pulumi up` to build the image and push it to the registry. This will take a few minutes.

```bash
pulumi up
```

### Step 5 (Optional) - Test the application

If you want you can run the image locally to see if it works:

```bash
docker run -p 3000:3000 -d <imageName output from pulumi>/myapp:latest
WARNING: The requested image's platform (linux/amd64) does not match the detected host platform (linux/arm64/v8) and no specific platform was requested
Server started on port 3000
```

#### Step 5.1 - Curl the GET endpoint

```bash
curl localhost:3000
```

You should see the following output:

```bash
Hello DevOpsDays Amsterdam!
```

#### Step 5.2 - Curl the POST endpoint

```bash
curl -X POST -H "Content-Type: application/json" -d '{"message":"This is a test message"}' http://localhost:3000
```

We expect following output:

```bash
Received message: This is a test message
```

Congratulations! You have successfully containerized an application and pushed it to a container registry. Please leave
the cluster up and running for [Chapter 2 - Deploy an Application](./02-deploy-app.md)

## Stretch Goals

- Can you add `json` support to the application?
- Can you add a `health` endpoint to the application?
- Can you change the `FROM` image a `chainguard` image, for enhanced security and smaller image size?

## Learn More

- [Pulumi](https://www.pulumi.com/)
- [Scaleway Pulumi Provider](https://www.pulumi.com/registry/packages/scaleway/)
- [Docker Pulumi Provider](https://www.pulumi.com/registry/packages/docker/)
