package main

import (
	"fmt"
	"go_dev/day7/example4/balance"
	"math/rand"
	"time"
)

func main(){
	// 先注册实例
	var insts []*balance.Instance
	for i:=0; i<10; i++ {
		host := fmt.Sprintf("127.0.%d.%d", rand.Intn(255), rand.Intn(255))
		one := balance.NewInstance(host, 8080)
		insts = append(insts, one)
	}

	balanceName := "random"

	for {
		if len(insts) > 0 {
			inst, err := balance.DoBalance(balanceName, insts)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(inst)
			time.Sleep(time.Second)
		}
	}

}
