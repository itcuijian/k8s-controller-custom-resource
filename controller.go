package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	clientset "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned"
	networkscheme "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned/scheme"
	informer "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/informers/externalversions/samplecrd/v1"
	listers "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/listers/samplecrd/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "sample-controller"

const (
	SuccessSynced         = "Synced"
	ErrorResourceExists   = "ErrorResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	MessageResourceSynced = "Foo synced successful"
)

type Controller struct {
	kubeclientset    kubernetes.Interface
	networkclientset clientset.Interface

	networksLister listers.NetworkLister
	networksSynced cache.InformerSynced

	workquue workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	networkclientset clientset.Interface,
	networkInformer informer.NetworkInformer,
) *Controller {
	utilruntime.Must(networkscheme.AddToScheme(scheme.Scheme))

	glog.V(4).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		networkclientset: networkclientset,
		networksLister:   networkInformer.Lister(),
		networksSynced:   networkInformer.Informer().HasSynced,
		workquue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")

	networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workquue.ShutDown()

	glog.Info("Starting Network control loop")
	glog.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(stopCh, c.networksSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
}

func (c *Controller) processNextWorkItem() bool {
	return true
}
