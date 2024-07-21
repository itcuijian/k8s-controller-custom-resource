package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	samplecrdv1 "github.com/itcuijim/k8s-controller-custom-resource/pkg/apis/samplecrd/v1"
	clientset "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned"
	networkscheme "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/clientset/versioned/scheme"
	informer "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/informers/externalversions/samplecrd/v1"
	listers "github.com/itcuijim/k8s-controller-custom-resource/pkg/client/listers/samplecrd/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	MessageResourceExists = "Resource %q already exists and is not managed by Network"
	MessageResourceSynced = "Network synced successful"
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

	networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueNetwork,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNetwork := oldObj.(*samplecrdv1.Network)
			newNetwork := newObj.(*samplecrdv1.Network)
			// 每经过 resyncPeriod 指定时间，为维护本地缓存都会使用最近一次 LIST 返回结果强制更新一次，
			// 从而保证本地缓存的有效性
			// 该更新事件对应的 Network 对象实际上并没有发生变化，即新旧两个 Network 对象对应的 ResourceVersion 是一样的
			// 此时也就是不需要做进一步的处理，直接 return 即可
			if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
				return
			}
			controller.enqueueNetwork(newObj)
		},
		DeleteFunc: controller.enqueueNetworkForDelete,
	})

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

	// 启动 threadiness 个 goroutine 来处理 Network 资源
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker 是一个死循环，调用 processNextWorkItem 来读取和处理工作队列里的信息
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workquue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workquue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workquue.Forget(obj)
			runtime.HandleError(fmt.Errorf("excepted string in workquue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workquue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)

		return nil
	}(obj)
	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	network, err := c.networksLister.Networks(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("Network does not exists in local cache: %s/%s, will delete it from Neutron ...", namespace, name)
			glog.Infof("Deleting network: %s/%s ...", namespace, name)
			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))

		// TODO: 添加删除 Network 资源对象的逻辑

		return err
	}

	glog.Infof("Try to process network: %#v ...", network)

	// TODO: 调谐过程-获取 network 预期状态与实际状态，比较两者的不同，并根据比较结果进行相关操作

	c.recorder.Event(network, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)

	return nil
}

// enqueueNetwork 获取一个 Network 资源对象并将它转换成“命名空间/名称”的字符串，然后将该字符串放进工作队列里
func (c *Controller) enqueueNetwork(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	c.workquue.AddRateLimited(key)
}

// enqueueNetworkForDelete 获取一个被删除的 Network 资源对象并将它转换成“命名空间/名称”的字符串，然后将该字符串放进工作队列里
func (c *Controller) enqueueNetworkForDelete(obj interface{}) {
	var key string
	var err error

	if key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	c.workquue.AddRateLimited(key)
}
