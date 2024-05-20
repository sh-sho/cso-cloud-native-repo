# vpa walk through

## prerequire

install metrics server: [cso-cloud-native-repo/kubernetes/argocd-apps/metrics-server.yaml](../argocd-apps/metrics-server.yaml)

## walk through

Autoscaler のソースコードをクローンします。（これ以外インストール方法ないの？Helm とか...）

```sh
cd kubernetes/vpa
git clone https://github.com/kubernetes/autoscaler.git
```

VPA をインストールします。（[参考: hack/vpa-up.sh の中身](#参考-hackvpa-upsh-の中身)）

```sh
cd autoscaler
./hack/vpa-up.sh
```

実行結果

```sh
customresourcedefinition.apiextensions.k8s.io/verticalpodautoscalercheckpoints.autoscaling.k8s.io created
customresourcedefinition.apiextensions.k8s.io/verticalpodautoscalers.autoscaling.k8s.io created
clusterrole.rbac.authorization.k8s.io/system:metrics-reader created
clusterrole.rbac.authorization.k8s.io/system:vpa-actor created
clusterrole.rbac.authorization.k8s.io/system:vpa-status-actor created
clusterrole.rbac.authorization.k8s.io/system:vpa-checkpoint-actor created
clusterrole.rbac.authorization.k8s.io/system:evictioner created
clusterrolebinding.rbac.authorization.k8s.io/system:metrics-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-actor created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-status-actor created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-checkpoint-actor created
clusterrole.rbac.authorization.k8s.io/system:vpa-target-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-target-reader-binding created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-evictioner-binding created
serviceaccount/vpa-admission-controller created
serviceaccount/vpa-recommender created
serviceaccount/vpa-updater created
clusterrole.rbac.authorization.k8s.io/system:vpa-admission-controller created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-admission-controller created
clusterrole.rbac.authorization.k8s.io/system:vpa-status-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-status-reader-binding created
deployment.apps/vpa-updater created
deployment.apps/vpa-recommender created
Generating certs for the VPA Admission Controller in /tmp/vpa-certs.
Certificate request self-signature ok
subject=CN = vpa-webhook.kube-system.svc
Uploading certs to the cluster.
secret/vpa-tls-certs created
Deleting /tmp/vpa-certs.
deployment.apps/vpa-admission-controller created
service/vpa-webhook created
```

サンプルアプリケーションをデプロイします。

```sh
kubectl -n example -f examples/hamster.yaml
```

hamster のデプロイメントは、以下の通りリソースリクエスト（CPU: 100m, Memory: 50Mi）が設定してあることを確認。

```yaml
resources:
  requests:
    cpu: 100m
    memory: 50Mi
```

VPA リソースが作成されていることを確認します。

```sh
kubectl get vpa -n example
```

実行結果

```sh
NAME          MODE   CPU    MEM       PROVIDED   AGE
hamster-vpa   Auto   587m   262144k   True       83s
```

リソースリクエストの内容を確認してみると、Pod の再起動が発生し、スケールアップされていることが確認できます。

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    vpaObservedContainers: hamster
    vpaUpdates:
      "Pod resources updated by hamster-vpa: container 0: memory request,
      cpu request"
  creationTimestamp: "2024-05-20T06:08:30Z"
  generateName: hamster-c6967774f-
  labels:
    app: hamster
    pod-template-hash: c6967774f
  name: hamster-c6967774f-t4q4d
  namespace: example
  ownerReferences:
    - apiVersion: apps/v1
      blockOwnerDeletion: true
      controller: true
      kind: ReplicaSet
      name: hamster-c6967774f
      uid: b3c960dc-bffb-4752-a636-2b3ee6836be2
  resourceVersion: "33972071"
  uid: 7bc4eaa5-0494-4b11-b14d-34f211f66e22
spec:
  containers:
    - args:
        - -c
        - while true; do timeout 0.5s yes >/dev/null; sleep 0.5s; done
      command:
        - /bin/sh
      image: registry.k8s.io/ubuntu-slim:0.1
      imagePullPolicy: IfNotPresent
      name: hamster
      resources:
        requests:
          cpu: 587m # こちら
          memory: 262144k # こちら
      terminationMessagePath: /dev/termination-log
      terminationMessagePolicy: File
      volumeMounts:
        - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
          name: kube-api-access-tlz9s
          readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  nodeName: 10.0.10.230
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  tolerations:
    - effect: NoExecute
      key: node.kubernetes.io/not-ready
      operator: Exists
      tolerationSeconds: 300
    - effect: NoExecute
      key: node.kubernetes.io/unreachable
      operator: Exists
      tolerationSeconds: 300
  volumes:
    - name: kube-api-access-tlz9s
      projected:
        defaultMode: 420
        sources:
          - serviceAccountToken:
              expirationSeconds: 3607
              path: token
          - configMap:
              items:
                - key: ca.crt
                  path: ca.crt
              name: kube-root-ca.crt
          - downwardAPI:
              items:
                - fieldRef:
                    apiVersion: v1
                    fieldPath: metadata.namespace
                  path: namespace
status:
  conditions:
    - lastProbeTime: null
      lastTransitionTime: "2024-05-20T06:09:03Z"
      status: "True"
      type: PodReadyToStartContainers
    - lastProbeTime: null
      lastTransitionTime: "2024-05-20T06:09:00Z"
      status: "True"
      type: Initialized
    - lastProbeTime: null
      lastTransitionTime: "2024-05-20T06:09:03Z"
      status: "True"
      type: Ready
    - lastProbeTime: null
      lastTransitionTime: "2024-05-20T06:09:03Z"
      status: "True"
      type: ContainersReady
    - lastProbeTime: null
      lastTransitionTime: "2024-05-20T06:09:00Z"
      status: "True"
      type: PodScheduled
  containerStatuses:
    - containerID: cri-o://86a78dbee076dc5880d521fc64b5bf86b76aa321b89d1165d5a9cdf1dc3fcf2d
      image: registry.k8s.io/ubuntu-slim:0.1
      imageID: 42caf9b4247ea7e4527144f345cd790f98a27b144b3afb91fa15fc6dcd687783
      lastState: {}
      name: hamster
      ready: true
      restartCount: 0
      started: true
      state:
        running:
          startedAt: "2024-05-20T06:09:02Z"
  hostIP: 10.0.10.230
  hostIPs:
    - ip: 10.0.10.230
  phase: Running
  podIP: 10.0.10.90
  podIPs:
    - ip: 10.0.10.90
  qosClass: Burstable
  startTime: "2024-05-20T06:09:00Z"
```

## 参考: hack/vpa-up.sh の中身

`hack/vpa-up.sh`

```sh
#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

$SCRIPT_ROOT/hack/vpa-process-yamls.sh create $*
```

`hack/vpa-process-yamls.sh`

```sh
# ... 省略 ...
COMPONENTS="vpa-v1-crd-gen vpa-rbac updater-deployment recommender-deployment admission-controller-deployment"

function script_path {
  # Regular components have deployment yaml files under /deploy/.  But some components only have
  # test deployment yaml files that are under hack/e2e. Check the main deploy directory before
  # using the e2e subdirectory.
  if test -f "${SCRIPT_ROOT}/deploy/${1}.yaml"; then
    echo "${SCRIPT_ROOT}/deploy/${1}.yaml"
  else
    echo "${SCRIPT_ROOT}/hack/e2e/${1}.yaml"
  fi
}
# ... 省略 ...
```

Argo CD とかで CD 基盤を作る時は、vpa-v1-crd-gen.yaml, vpa-rbac.yaml, updater-deployment.yaml, recommender-deployment.yaml, admission-controller-deployment.yaml をデプロイすれば十分に見える。
