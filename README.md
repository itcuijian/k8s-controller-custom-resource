### 编译启动

编译命令：

```
$ go build -mod=vendor -o samplecrd-controller .
```

启动命令：

```
$ ./samplecrd-controller -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```
