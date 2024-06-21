package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)

	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop) // 关闭通道，通知接收者
		<-c
		os.Exit(1) // 第二次直接异常退出
	}()

	return stop
}
