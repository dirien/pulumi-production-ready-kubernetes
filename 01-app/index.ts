import * as pulumi from "@pulumi/pulumi";
import * as scaleway from "@ediri/scaleway";
import * as docker from "@pulumi/docker";

const containerRegistry = new scaleway.RegistryNamespace("devopsdaysams-registry", {
    isPublic: true,
})

const config = new pulumi.Config();

// Build and publish the image.
const image = new docker.Image("devopsdaysams-image", {
    build: {
        context: "app",
        platform: "linux/amd64",
        builderVersion
    },
    imageName: containerRegistry.endpoint.apply(s => `${s}/myapp`),
    registry: {
        server: containerRegistry.endpoint,
        username: "nologin",
        password: config.requireSecret("secret_key"),
    }
}, {dependsOn: containerRegistry});

export const imageRegistry = containerRegistry.endpoint;
export const imageName = image.imageName;
