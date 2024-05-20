# hpa walk through

## prerequire

install metrics server: [cso-cloud-native-repo/kubernetes/argocd-apps/metrics-server.yaml](../argocd-apps/metrics-server.yaml)

## walk through

```sh
kubectl apply -f https://k8s.io/examples/application/php-apache.yaml -n example
```

Horizontal Pod Autoscaler を作成します。

```sh
kubectl -n example autoscale deployment php-apache --cpu-percent=20 --min=1 --max=10
```

HPA リソースを watch しておきます。

```sh
$ kubectl get hpa -n example -w
NAME         REFERENCE               TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   0%/20%    1         10        1          53s
```

当該コンテナに対して、負荷をかけます

```sh
$ kubectl -n example run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh
/ # while true; do wget -q -O- http://php-apache; done
OK!OK!OK!OK!OK!OK!OK!OK!OK!OK!, ...
```

HPA リソースが以下のように変化することを確認します。

```sh
k get hpa -w
NAME         REFERENCE               TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   0%/20%    1         10        1          20m
php-apache   Deployment/php-apache   217%/20%   1         10        1          20m
php-apache   Deployment/php-apache   217%/20%   1         10        4          21m
php-apache   Deployment/php-apache   244%/20%   1         10        8          21m
php-apache   Deployment/php-apache   90%/20%    1         10        10         21m
php-apache   Deployment/php-apache   70%/20%    1         10        10         21m
php-apache   Deployment/php-apache   55%/20%    1         10        10         22m
php-apache   Deployment/php-apache   23%/20%    1         10        10         22m
php-apache   Deployment/php-apache   0%/20%     1         10        10         22m
```

ターゲットの使用率が 20%なのに対して、現在の CPU 使用率が 244%に増加していることがわかる。このとき、ターゲットの CPU 使用率となるように Pod の数を HPA がスケールしてくれる。
