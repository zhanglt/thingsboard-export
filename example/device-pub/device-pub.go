package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron"
)

type Message struct {
	DeviceName  string  `json:"name"`
	CmdName     string  `json:"cmd"`
	Randnum     float64 `json:"randnum"`
	Temperature float64 `json:"temperatue"`
	Humidity    float64 `json:"humidity"`
}
type MessageSlice struct {
	Messages []Message `json:"Messages"`
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
func RandFloat64(min, max int) float64 {
	if min >= max || min == 0 || max == 0 {
		return float64(max)
	}
	return Decimal(float64(rand.Intn(max-min)+min) + rand.Float64())
}
func main() {

	var client mqtt.Client
	var err error
	client, err = createClient("testID", "10.0.0.70:1883")
	if err != nil {
		fmt.Println("Client create:", err)
	}
	/*
		var msgRcvd := func(client *mqtt.Client, message mqtt.Message) {
				fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
			}

			if token := client.Subscribe("example/topic", 0, msgRcvd); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
			}

	*/

	crontab := cron.New()
	task := func() {
		fmt.Println("hello world")
		msg := Message{DeviceName: "MQ_DEVICE", CmdName: "randnum", Randnum: RandFloat64(10, 100), Temperature: RandFloat64(5, 40), Humidity: RandFloat64(10, 60)}
		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("JSON ERR:", err)
		}
		fmt.Println(string(b))
		//{"name":"MQ_DEVICE","cmd":"randnum","randnum":"163.7","temperature":"286.3","humidity":"323.8"}
		if token := client.Publish("DataTopic", 0, false, string(b)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}

	}
	// 添加定时任务, * * * * * 是 crontab,表示每分钟执行一次
	crontab.AddFunc("0/10 * * * * ?", task) //每10秒执行一次
	// 启动定时器
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {}

}

func createClient(clientID string, broker string) (mqtt.Client, error) {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername("admin")
	opts.SetPassword("public")

	client := mqtt.NewClient(opts)
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		fmt.Println("create client error", token.Error)
		return client, token.Error()
	}

	return client, nil
}
