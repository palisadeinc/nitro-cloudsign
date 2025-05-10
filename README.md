# CloudSign Nitro

Collection of tools and configuration required to run cloud sign inside AWS Nitro.

# Compatibility

| Servitor Version | Cloud Sign Version |
| ---------------- | ------------------ |
| 1.0.2            | 1.9.2              |
| 1.0.1            | 1.9.1              |

# Setup

```
                                             ┌────────────┐
                                             │            │
                                             │  servitor  │     ┌────────────┐
                                             │            │     │            │
┌───────────────────────────────────────┐    └──────▲─────┘     │  Palisade  │
│                     AWS Nitro Enclave │           │           │            │
│                                       │           │           └──────▲─────┘
│                                       │           │                  │
│    ┌─────────────┐      ┌─────────┐   │      ┌────┼────┐             │
│    │             │      │         │   │      │         │             │
│    │  CloudSign  ┼──────►  proxy  ┼───┼──────►  proxy  ┼─────────────┤
│    │             │      │         │   │      │         │             │
│    └─────────────┘      └─────────┘   │      └─────────┘             │
│                                       │                              │
│                                       │                   ┌──────────▼──────────┐
└───────────────────────────────────────┘                   │  Postgres Database  │
                                                            └─────────────────────┘
```

The proxy used is tools/service/tacos.service. The linux/amd64 binary is embedded as tools/tacos.
Servitor is a golang binary in servitor/ directory.

# EC2 Setup

Use m5.xlarge EC2 Instance with Amazon Linux 2023 AMI.

Configure Docker:

```shell
sudo yum install -y docker
sudo systemctl enable docker --now
sudo usermod -aG docker ssm-user
newgrp docker  # Activate without relogin
```

Setup AWS Nitro:

```shell
sudo dnf install -y aws-nitro-enclaves-cli aws-nitro-enclaves-cli-devel
sudo usermod -aG ne ssm-user
newgrp ne
```

Ensure Nitro allocator `/etc/nitro_enclaves/allocator.yaml` has allocated at least 1G memory and 2 CPU Cores. Make sure you restart the nitro-enclaves-allocator service if you made any modifications in this file.

Ensure Nitro version is 1.4+.

```shell
nitro-cli --version  # Should show v1.3+
```

Start nitro allocator services:

```shell
sudo systemctl enable nitro-enclaves-allocator docker
sudo systemctl restart nitro-enclaves-allocator
```
