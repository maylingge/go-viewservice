package main


import (
//	"fmt"
	"viewservice"
//	"net/rpc"
//	"time"
	"server"
)

type PingInfo struct {
	Num int
	Id	string
}

func main() {
	viewaddr := ":1234"
	viewservice.StartViewService(viewaddr) 
//	<-time.After(10*time.Second)	
//	fmt.Println(viewservice.GetView()) 
//	
//	client, err := rpc.DialHTTP("tcp", ":1234")
//	if err != nil {
//		fmt.Println("failed to dial rpc server")
//		fmt.Println(err)
//		return
//	}
//	
//	args := PingInfo{Num: 0, Id: "First"}
//	var reply viewservice.View
//	err = client.Call("ViewService.Ping", &args, &reply)
//	if err != nil {
//		fmt.Println("failed to ping rpc server")
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(reply)
//	fmt.Println(viewservice.GetView())
//	args.Num = reply.Num
//	err = client.Call("ViewService.Ping", &args, &reply)
//	if err != nil {
//		fmt.Println("failed to ping rpc server")
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(reply)
//	fmt.Println(viewservice.GetView())
//
//	args.Num = 0
//	args.Id = "Second"
//	err = client.Call("ViewService.Ping", &args, &reply)
//	if err != nil {
//		fmt.Println("failed to ping rpc server")
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(reply)
//	fmt.Println(viewservice.GetView())
//	


	go func() {
		server.Start(":2346", viewaddr)
	}()
	server.Start(":2348", viewaddr)
}
