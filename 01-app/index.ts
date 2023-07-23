import * as pulumi from "@pulumi/pulumi";
import * as scaleway from "@ediri/scaleway";
import * as docker from "@pulumi/docker";

const containerRegistry = new scaleway.RegistryNamespace("scalewayworkshop-registry", {
    isPublic: true,
})

const config = new pulumi.Config();

// Build and publish the image.
const image = new docker.Image("scalewayworkshop-image", {
    build: {
        context: "app",
        platform: "linux/amd64",
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
