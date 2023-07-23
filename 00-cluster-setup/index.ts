import * as pulumi from "@pulumi/pulumi";
import * as scaleway from "@ediri/scaleway";

const clusterConfig = new pulumi.Config("cluster")

const kapsule = new scaleway.K8sCluster("scalewayworkshop-cluster", {
    version: clusterConfig.require("version"),
    cni: "cilium",
    deleteAdditionalResources: true,
    tags: [
        "pulumi",
        "workshop",
    ],
    autoUpgrade: {
        enable: clusterConfig.requireBoolean("auto_upgrade"),
        maintenanceWindowStartHour: 3,
        maintenanceWindowDay: "monday"
    },
});

const nodeConfig = new pulumi.Config("node")

const pool = new scaleway.K8sPool("scalewayworkshop-node-pool", {
    nodeType: nodeConfig.require("node_type"),
    size: nodeConfig.requireNumber("node_count"),
    autoscaling: nodeConfig.requireBoolean("auto_scale"),
    autohealing: nodeConfig.requireBoolean("auto_heal"),
    clusterId: kapsule.id,
});

export const kapsuleName = kapsule.name;
export const kapsuleVersion = kapsule.version;
export const kapsuleAutoUpgrade = kapsule.autoUpgrade.enable;
export const kapusuleNodeCount = pool.size;
export const kapsuleNodeType = pool.nodeType;
export const kapsuleID = kapsule.id;
export const region = kapsule.region;
export const kubeconfig = pulumi.secret(kapsule.kubeconfigs[0].configFile);
