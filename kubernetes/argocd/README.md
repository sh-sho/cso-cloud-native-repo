# Argo CD

## Install

```sh
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Change the argocd-server service type to `LoadBalancer` and use Flexible LoadBalancer.

```sh
kubectl patch svc argocd-server \
    -n argocd \
    -p '{"metadata: {"annotations": {"oci.oraclecloud.com/load-balancer-type: "lb", "service.beta.kubernetes.io/oci-load-balancer-shape": "flexible", "service.beta.kubernetes.io/oci-load-balancer-shape-flex-min": "10", "service.beta.kubernetes.io/oci-load-balancer-shape-flex-max": "30"}}, "spec": {"type": "LoadBalancer"}}'
```

## References

- [https://argo-cd.readthedocs.io/en/stable/getting_started/](https://argo-cd.readthedocs.io/en/stable/getting_started/)
- [https://docs.oracle.com/ja-jp/iaas/Content/ContEng/Tasks/contengcreatingloadbalancers-subtopic.htm#flexible](https://docs.oracle.com/ja-jp/iaas/Content/ContEng/Tasks/contengcreatingloadbalancers-subtopic.htm#flexible)
