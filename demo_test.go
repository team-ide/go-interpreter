package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNet(t *testing.T) {
	var err error
	fmt.Println("连接不存在的IP")
	start := time.Now()
	_, err = net.Dial("tcp", "192.168.81.248:10001")
	end := time.Now()
	fmt.Println("耗时：", end.UnixMilli()-start.UnixMilli())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("连接存在的IP")
	start = time.Now()
	_, err = net.Dial("tcp", "192.168.81.48:10001")
	end = time.Now()
	fmt.Println("耗时：", end.UnixMilli()-start.UnixMilli())
	if err != nil {
		fmt.Println(err)
	}
}
