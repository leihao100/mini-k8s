package main

import (
	"MiniK8S/pkg/api/config"
	apiserver "MiniK8S/pkg/apiServer"
	"MiniK8S/pkg/kubelet"
	"context"
	"time"
)

func main() {
	//go func() {
	api := apiserver.NewApiServer()
	ctx, cancel := context.WithCancel(context.Background())
	api.Run(cancel)
	//<-context.Done()
	//}()
	time.Sleep(3 * time.Second)
	k := kubelet.NewKubelet(config.Node{})
	k.Run(ctx, cancel)
	k.SendMessage()
}

//func main() {
//	// 创建一个带有3秒超时的Context
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel() // 确保在main函数返回前取消上下文
//
//	// 启动一个goroutine，模拟一个工作
//	go doWork(ctx)
//
//	// 等待5秒以观察doWork是否会因超时而结束
//	time.Sleep(5 * time.Second)
//	fmt.Println("Main function finished")
//}
//
//func doWork(ctx context.Context) {
//	fmt.Println("Work started")
//	select {
//	case <-time.After(1 * time.Second):
//		// 模拟一个需要5秒才能完成的工作
//		fmt.Println("Work completed")
//	case <-ctx.Done():
//		// 当Context被取消或超时时
//		fmt.Println("Work cancelled:", ctx.Err())
//	}
//}
