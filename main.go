package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"

	tbTransforms "github.com/zhanglt/thingsboard-export/transforms"
)

const (
	serviceKey = "TBExport"
)

func main() {
	// 关闭安全模式
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) 创建一个EdgeX SDK实例并初始化
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	// 2) 从App SDK加载thingsboard特定的MQTT配置

	config, err := tbTransforms.LoadTBMQTTConfig(edgexSdk)
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to load AWS MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	// 3) 从配置中获取设备名称过滤器
	deviceNamesCleaned := util.DeleteEmptyAndTrim(strings.FieldsFunc(config.DeviceNames, util.SplitComma))
	edgexSdk.LoggingClient.Debug(fmt.Sprintf("Device names read %s\n", deviceNamesCleaned))
	fmt.Println("deviceNamesCleaned:", deviceNamesCleaned)
	// 4) 管道配置，是每次触发事件时要执行的函数集合。
	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(deviceNamesCleaned).FilterByDeviceName,
		tbTransforms.NewConversion().TransformToTB,
		tbTransforms.NewTBMQTTSender(edgexSdk.LoggingClient, config).MQTTSend,
		//printTBDataToConsole,
		//text

	)
	// 5) 启动SDK并开始监听触发管道的事件。
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func text(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	fmt.Println("test....pip")
	return false, nil
}
func printTBDataToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}
	for i := 0; i < len(params); i++ {
		fmt.Println("=======:", params[i].(string))
	}

	// Leverage the built in logging service in EdgeX
	edgexcontext.LoggingClient.Debug("Printed to console")

	edgexcontext.Complete([]byte(params[0].(string)))
	return false, nil

}
