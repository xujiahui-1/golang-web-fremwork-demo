package main

import (
	"context"
	"fmt"
	"golang-web-fremwork-demo/framework"
	"log"
	"time"
)

// FooControllerHandler 创建一个设置超时的处理器
func FooControllerHandler(c *framework.Context) error {

	// 首先 携程的处理需要chan
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	// 生成一个超时的 Context
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()
	go func() {
		defer func() {
			//使用recover()函数判断有没有panic
			if p := recover(); p != nil {
				//获取到panic写入channel
				panicChan <- p
			}
		}()
		// 写个超时
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")
		finish <- struct{}{}
	}()
	//监听chan
	select {
	case p := <-panicChan:
		//panic 情况
		//写入的时候防止并发问题，上锁
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		// 正常情况下，finish chan中存入了空结构体
		fmt.Println("finish")
	case <-durationCtx.Done():
		//超时情况
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}

	return nil
}
