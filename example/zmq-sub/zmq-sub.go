package main

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/pebbe/zmq4/examples/kvsimple"

	"fmt"
)

func main() {
	socket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		fmt.Println("newsocker err:", err)
	}
	socket.SetSubscribe("events")
	var er = socket.Connect("tcp://10.0.0.70:5563")
	fmt.Println("er", er)
	if er != nil {
		fmt.Println(er)

	}
	//fmt.Println(updates)
	kvmap := make(map[string]*kvsimple.Kvmsg)

	sequence := int64(0)
	for ; true; sequence++ {

		kvmsg, err := kvsimple.RecvKvmsg(socket)
		if err != nil {
			fmt.Println("err:", err)
			break //  Interrupted
		}

		kvmsg.Store(kvmap)
		fmt.Println(kvmsg)

	}
	fmt.Printf("Interrupted\n%d messages in\n", sequence)
}
