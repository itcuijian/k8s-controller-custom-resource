package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned"
	"github.com/itcuijim/k8s-controller-custom-resource/pkg/client/informers/externalversions"
	"github.com/itcuijim/k8s-controller-custom-resource/signals"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// 设置系统信号处理者，实现优雅退出
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	networkClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clietnset: %s", err.Error())
	}

	networkInformerFactory := externalversions.NewSharedInformerFactory(networkClient, time.Second*30)

	controller := NewController(kubeClient, networkClient, networkInformerFactory.Samplecrd().V1().Networks())

	go networkInformerFactory.Start(stopCh)

	if err := controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster")
}
