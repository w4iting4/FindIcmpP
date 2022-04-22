package main

import (
	"FindIcmpP/flag"
	"fmt"
	"time"
)

func main() {
	start := time.Now().Unix()
	flag.ParseFlag()
	fmt.Println("程序运行时间：" + string(time.Now().Unix()-start) + "s")
	//utils.Chcp65001()
}
