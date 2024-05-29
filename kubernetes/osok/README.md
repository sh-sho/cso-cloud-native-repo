# OCI Service Operator for Kubernetes

## toc

<!-- @import "[TOC]" {cmd="toc" depthFrom=3 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Install Operator SDK](#install-operator-sdk)

<!-- /code_chunk_output -->

## setup

[https://github.com/oracle/oci-service-operator/blob/main/docs/instal＃lation.md#install-operator-sdk](https://github.com/oracle/oci-service-operator/blob/main/docs/installation.md#install-operator-sdk) に従って実施していく

### Install Operator SDK

[https://sdk.operatorframework.io/docs/installation/](https://sdk.operatorframework.io/docs/installation/)

CPU アーキテクチャと OS を環境変数に置いておく

```sh
export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')
```

バイナリをダウンロードする

```sh
export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.34.2
curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}
```

GPG キー（公開鍵）をインストールしておく

```sh
gpg --keyserver keyserver.ubuntu.com --recv-keys 052996E2A20B5C7E
```

ダウンロードしたバイナリを検証しておく

```sh
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt.asc
gpg -u "Operator SDK (release) <cncf-operator-sdk@cncf.io>" --verify checksums.txt.asc
```

```sh
grep operator-sdk_${OS}_${ARCH} checksums.txt | sha256sum -c -
```

```sh
chmod +x operator-sdk_${OS}_${ARCH} && sudo mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk
```

インストールが完了したか確認

```sh
operator-sdk --help
```

OK

```sh
CLI tool for building Kubernetes extensions and tools.

Usage:
  operator-sdk [flags]
  operator-sdk [command]

Examples:
The first step is to initialize your project:
    operator-sdk init [--plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]]

<PLUGIN KEYS> is a comma-separated list of plugin keys from the following table
and <PROJECT VERSION> a supported project version for these plugins.

                                   Plugin keys | Supported project versions
-----------------------------------------------+----------------------------
           ansible.sdk.operatorframework.io/v1 |                          3
              declarative.go.kubebuilder.io/v1 |                       2, 3
       deploy-image.go.kubebuilder.io/v1-alpha |                          3
                          go.kubebuilder.io/v2 |                       2, 3
                          go.kubebuilder.io/v3 |                          3
                          go.kubebuilder.io/v4 |                          3
               grafana.kubebuilder.io/v1-alpha |                          3
              helm.sdk.operatorframework.io/v1 |                          3
 hybrid.helm.sdk.operatorframework.io/v1-alpha |                          3
            quarkus.javaoperatorsdk.io/v1-beta |                          3

For more specific help for the init command of a certain plugins and project version
configuration please run:
    operator-sdk init --help --plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]

Default plugin keys: "go.kubebuilder.io/v4"
Default project version: "3"


Available Commands:
  alpha            Alpha-stage subcommands
  bundle           Manage operator bundle metadata
  cleanup          Clean up an Operator deployed with the 'run' subcommand
  completion       Load completions for the specified shell
  create           Scaffold a Kubernetes API or webhook
  edit             Update the project configuration
  generate         Invokes a specific generator
  help             Help about any command
  init             Initialize a new project
  olm              Manage the Operator Lifecycle Manager installation in your cluster
  pkgman-to-bundle Migrates packagemanifests to bundles
  run              Run an Operator in a variety of environments
  scorecard        Runs scorecard
  version          Print the operator-sdk version

Flags:
  -h, --help                     help for operator-sdk
      --plugins strings          plugin keys to be used for this subcommand execution
      --project-version string   project version (default "3")
      --verbose                  Enable verbose logging

Use "operator-sdk [command] --help" for more information about a command.
```

### Install OLM(Operator Lifecycle Manager)

[https://github.com/oracle/oci-service-operator/blob/main/docs/installation.md#install-operator-lifecycle-manager-olm](https://github.com/oracle/oci-service-operator/blob/main/docs/installation.md#install-operator-lifecycle-manager-olm)

OLM をインストールする。
最新版は、[https://github.com/operator-framework/operator-lifecycle-manager/releases](https://github.com/operator-framework/operator-lifecycle-manager/releases) から確認

```sh
operator-sdk olm install --version 0.28.0
```

### Deploy OSOK

Instance Principal 使うので、こっち

[https://github.com/oracle/oci-service-operator/blob/main/docs/installation.md#enable-instance-principal](https://github.com/oracle/oci-service-operator/blob/main/docs/installation.md#enable-instance-principal)

Kubernetes クラスタにデプロイする
最新版は、[https://github.com/oracle/oci-service-operator/releases](https://github.com/oracle/oci-service-operator/releases) から確認

```sh
operator-sdk run bundle iad.ocir.io/oracle/oci-service-operator-bundle:1.1.9
```

（おまけ）アンデプロイする

```sh
operator-sdk cleanup oci-service-operator
```

### Deploy Streaming

```sh
kubectl apply -f streaming.yaml
```
