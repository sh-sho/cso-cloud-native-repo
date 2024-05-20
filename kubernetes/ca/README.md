# ca walk through

## prerequire

特になし。Virtual Nodes だとダメみたい。

## walk through

Cluster Autoscaler 用のリソースをデプロイする

`cluster-autoscaler.yaml`

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-addon: cluster-autoscaler.addons.k8s.io
    k8s-app: cluster-autoscaler
  name: cluster-autoscaler
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-autoscaler
  labels:
    k8s-addon: cluster-autoscaler.addons.k8s.io
    k8s-app: cluster-autoscaler
rules:
  - apiGroups: [""]
    resources: ["events", "endpoints"]
    verbs: ["create", "patch"]
  - apiGroups: [""]
    resources: ["pods/eviction"]
    verbs: ["create"]
  - apiGroups: [""]
    resources: ["pods/status"]
    verbs: ["update"]
  - apiGroups: [""]
    resources: ["endpoints"]
    resourceNames: ["cluster-autoscaler"]
    verbs: ["get", "update"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["watch", "list", "get", "patch", "update"]
  - apiGroups: [""]
    resources:
      - "pods"
      - "services"
      - "replicationcontrollers"
      - "persistentvolumeclaims"
      - "persistentvolumes"
    verbs: ["watch", "list", "get"]
  - apiGroups: ["extensions"]
    resources: ["replicasets", "daemonsets"]
    verbs: ["watch", "list", "get"]
  - apiGroups: ["policy"]
    resources: ["poddisruptionbudgets"]
    verbs: ["watch", "list"]
  - apiGroups: ["apps"]
    resources: ["statefulsets", "replicasets", "daemonsets"]
    verbs: ["watch", "list", "get"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses", "csinodes"]
    verbs: ["watch", "list", "get"]
  - apiGroups: ["batch", "extensions"]
    resources: ["jobs"]
    verbs: ["get", "list", "watch", "patch"]
  - apiGroups: ["coordination.k8s.io"]
    resources: ["leases"]
    verbs: ["create"]
  - apiGroups: ["coordination.k8s.io"]
    resourceNames: ["cluster-autoscaler"]
    resources: ["leases"]
    verbs: ["get", "update"]
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["watch", "list"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["csidrivers", "csistoragecapacities"]
    verbs: ["watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cluster-autoscaler
  namespace: kube-system
  labels:
    k8s-addon: cluster-autoscaler.addons.k8s.io
    k8s-app: cluster-autoscaler
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create", "list", "watch"]
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames:
      ["cluster-autoscaler-status", "cluster-autoscaler-priority-expander"]
    verbs: ["delete", "get", "update", "watch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cluster-autoscaler
  labels:
    k8s-addon: cluster-autoscaler.addons.k8s.io
    k8s-app: cluster-autoscaler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-autoscaler
subjects:
  - kind: ServiceAccount
    name: cluster-autoscaler
    namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cluster-autoscaler
  namespace: kube-system
  labels:
    k8s-addon: cluster-autoscaler.addons.k8s.io
    k8s-app: cluster-autoscaler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cluster-autoscaler
subjects:
  - kind: ServiceAccount
    name: cluster-autoscaler
    namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cluster-autoscaler
  namespace: kube-system
  labels:
    app: cluster-autoscaler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cluster-autoscaler
  template:
    metadata:
      labels:
        app: cluster-autoscaler
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8085"
    spec:
      serviceAccountName: cluster-autoscaler
      containers:
        - image: lhr.ocir.io/oracle/oci-cluster-autoscaler:1.29.0-10
          name: cluster-autoscaler
          resources:
            limits:
              cpu: 100m
              memory: 300Mi
            requests:
              cpu: 100m
              memory: 300Mi
          command:
            - ./cluster-autoscaler
            - --v=4
            - --stderrthreshold=info
            - --cloud-provider=oci
            - --max-node-provision-time=25m
            - --nodes=1:2:ocid1.nodepool.oc1.uk-london-1.aaaaaaaaaqqfcvolpdrndsvgp4fnyly52sntduhbt7gkt5guynhezlqk2kla
            - --scale-down-delay-after-add=10m
            - --scale-down-unneeded-time=10m
            - --unremovable-node-recheck-timeout=5m
            - --balance-similar-node-groups
            - --balancing-ignore-label=displayName
            - --balancing-ignore-label=hostname
            - --balancing-ignore-label=internal_addr
            - --balancing-ignore-label=oci.oraclecloud.com/fault-domain
          imagePullPolicy: "Always"
          env:
            - name: OKE_USE_INSTANCE_PRINCIPAL
              value: "true"
            - name: OCI_SDK_APPEND_USER_AGENT
              value: "oci-oke-cluster-autoscaler"
```

SA, CR, CRB, Role, RB は割とどうでも良くて、大事なのは Deployment。

まずは、CA 用の Deployment のイメージを選択する。今回は、ロンドンかつ v1.29.1 なので `lhr.ocir.io/oracle/oci-cluster-autoscaler:1.29.0-10` を選択。

```yaml
image: lhr.ocir.io/oracle/oci-cluster-autoscaler:1.29.0-10
```

オートスケールの対象となるノードプールを作成し、それを `--nodes` に追加

```yaml
--nodes=1:2:ocid1.nodepool.oc1.uk-london-1.aaaaaaaaaqqfcvolpdrndsvgp4fnyly52sntduhbt7gkt5guynhezlqk2kla
```

あとは、デフォルトでとりあえず OK。

Cluster Autoscaler をデプロイする

```sh
kubectl apply -f ca/cluster-autoscaler.yaml
```

Cluster Autoscaler が正常にデプロイされたかどうかログをみて確認します。

```sh
kubectl -n kube-system logs -f deployment.apps/cluster-autoscaler
```

以下のようなログが継続的に出力されれば OK

```sh
^CI0520 12:49:07.633076       1 static_autoscaler.go:290] Starting main loop
I0520 12:49:07.633905       1 oci_manager.go:451] did not find node pool for reference: {AvailabilityDomain:UK-LONDON-1-AD-1 Name:10.0.10.230 CompartmentID:ocid1.compartment.oc1..aaaaaaaac5t2rwhyzq6fm6pepwdn6gc434ymir7sgvh4ac4sd5sre4altoka InstanceID:ocid1.instance.oc1.uk-london-1.anwgiljrssl65iqc2kmgxagkfepoxgts4cfq53o7n5wq7mg3fbn737h5iqcq NodePoolID:ocid1.nodepool.oc1.uk-london-1.aaaaaaaarbie6tyvecy7jgyhybjv5tbjjja2fomiexkfgobn3nasgahvhftq InstancePoolID: PrivateIPAddress:10.0.10.230 PublicIPAddress: Shape:VM.Standard.E4.Flex}
I0520 12:49:07.633951       1 oci_manager.go:451] did not find node pool for reference: {AvailabilityDomain:UK-LONDON-1-AD-1 Name:10.0.10.230 CompartmentID:ocid1.compartment.oc1..aaaaaaaac5t2rwhyzq6fm6pepwdn6gc434ymir7sgvh4ac4sd5sre4altoka InstanceID:ocid1.instance.oc1.uk-london-1.anwgiljrssl65iqc2kmgxagkfepoxgts4cfq53o7n5wq7mg3fbn737h5iqcq NodePoolID:ocid1.nodepool.oc1.uk-london-1.aaaaaaaarbie6tyvecy7jgyhybjv5tbjjja2fomiexkfgobn3nasgahvhftq InstancePoolID: PrivateIPAddress:10.0.10.230 PublicIPAddress: Shape:VM.Standard.E4.Flex}
I0520 12:49:07.634004       1 filter_out_schedulable.go:63] Filtering out schedulables
I0520 12:49:07.634015       1 filter_out_schedulable.go:120] 0 pods marked as unschedulable can be scheduled.
I0520 12:49:07.634035       1 filter_out_schedulable.go:83] No schedulable pods
I0520 12:49:07.634042       1 filter_out_daemon_sets.go:40] Filtering out daemon set pods
I0520 12:49:07.634048       1 filter_out_daemon_sets.go:49] Filtered out 0 daemon set pods, 0 unschedulable pods left
I0520 12:49:07.634070       1 static_autoscaler.go:547] No unschedulable pods
I0520 12:49:07.634089       1 static_autoscaler.go:570] Calculating unneeded nodes
I0520 12:49:07.634100       1 oci_manager.go:451] did not find node pool for reference: {AvailabilityDomain:UK-LONDON-1-AD-1 Name:10.0.10.230 CompartmentID:ocid1.compartment.oc1..aaaaaaaac5t2rwhyzq6fm6pepwdn6gc434ymir7sgvh4ac4sd5sre4altoka InstanceID:ocid1.instance.oc1.uk-london-1.anwgiljrssl65iqc2kmgxagkfepoxgts4cfq53o7n5wq7mg3fbn737h5iqcq NodePoolID:ocid1.nodepool.oc1.uk-london-1.aaaaaaaarbie6tyvecy7jgyhybjv5tbjjja2fomiexkfgobn3nasgahvhftq InstancePoolID: PrivateIPAddress:10.0.10.230 PublicIPAddress: Shape:VM.Standard.E4.Flex}
I0520 12:49:07.634115       1 pre_filtering_processor.go:57] Node 10.0.10.230 should not be processed by cluster autoscaler (no node group config)
I0520 12:49:07.634126       1 pre_filtering_processor.go:67] Skipping 10.0.10.245 - node group min size reached (current: 1, min: 1)
I0520 12:49:07.634171       1 static_autoscaler.go:617] Scale down status: lastScaleUpTime=2024-05-20 11:48:17.576354703 +0000 UTC m=-3579.919975677 lastScaleDownDeleteTime=2024-05-20 11:48:17.576354703 +0000 UTC m=-3579.919975677 lastScaleDownFailTime=2024-05-20 11:48:17.576354703 +0000 UTC m=-3579.919975677 scaleDownForbidden=false scaleDownInCooldown=false
I0520 12:49:07.634203       1 static_autoscaler.go:642] Starting scale down
```

`I0520 12:49:07.634100       1 oci_manager.go:451] did not find node pool for reference` は、Autoscale 対象外のノードプールなので問題なし

ノードの数が増えることを確認するために、watch しておく

```sh
kubectl get nodes -w
```

nginx のデプロイメントをデプロイします。

```sh
kubectl -n example apply -f ca/nginx.yaml
```

Pod の数を 2 -> 100 へ増やします。

```sh
kubectl -n example scale deployment nginx-deploymeny --replicas=100
```

watch しているノード数が増えていることが確認できます。（Compute Instance のプロビジョニングのためそれ相応に時間がかかります）

```sh
NAME          STATUS   ROLES   AGE     VERSION
10.0.10.230   Ready    node    14d     v1.29.1
10.0.10.245   Ready    node    3h35m   v1.29.1
10.0.10.245   Ready    node    3h36m   v1.29.1
10.0.10.230   Ready    node    14d     v1.29.1
10.0.10.230   Ready    node    14d     v1.29.1
10.0.10.246   NotReady   <none>   0s      v1.29.1
10.0.10.246   NotReady   <none>   0s      v1.29.1
10.0.10.246   NotReady   <none>   0s      v1.29.1
10.0.10.246   NotReady   <none>   0s      v1.29.1
10.0.10.246   NotReady   <none>   2s      v1.29.1
10.0.10.246   NotReady   <none>   10s     v1.29.1
10.0.10.246   NotReady   node     25s     v1.29.1
10.0.10.246   NotReady   node     32s     v1.29.1
10.0.10.245   Ready      node     3h39m   v1.29.1
10.0.10.246   NotReady   node     2m4s    v1.29.1
10.0.10.246   NotReady   node     2m34s   v1.29.1
10.0.10.246   Ready      node     2m45s   v1.29.1
10.0.10.246   Ready      node     2m45s   v1.29.1
10.0.10.246   Ready      node     2m46s   v1.29.1
10.0.10.246   Ready      node     2m47s   v1.29.1
```
