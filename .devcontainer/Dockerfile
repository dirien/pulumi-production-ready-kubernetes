FROM --platform=linux/amd64 pulumi/pulumi:3.122.0

RUN <<EOF
## Install pulumi and set to PATH
curl -fsSL https://get.pulumi.com | sh
PATH="/root/.pulumi/bin:${PATH}"

## Install pulumictl and set to PATH
mkdir -p /root/pulumictl && cd /root/pulumictl/
wget https://github.com/pulumi/pulumictl/releases/download/v0.0.42/pulumictl-v0.0.42-linux-amd64.tar.gz -O /root/pulumictl/pulumictl-v0.0.42-linux-amd64.tar.gz
tar -xvf /root/pulumictl/pulumictl-v0.0.42-linux-amd64.tar.gz -C /root/pulumictl/
PATH="//root/pulumictl/:${PATH}"

## Install fluxctl
curl -s https://fluxcd.io/install.sh | bash

## Install Scaleway CLI
curl -o /usr/local/bin/scw -L "https://github.com/scaleway/scaleway-cli/releases/download/v2.18.0/scaleway-cli_2.18.0_linux_amd64"
chmod +x /usr/local/bin/scw
EOF
