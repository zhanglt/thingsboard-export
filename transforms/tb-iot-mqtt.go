package transforms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	tbIOTMQTTProtocol       = "tbIOTMQTTProtocol"
	tbIOTMQTTHost           = "tbIOTMQTTHost"
	tbIOTMQTTPort           = "tbIOTMQTTPort"
	tbUserName              = "TBIoTMQTTUserName"
	topic                   = "topic"
	tbIOTThingName          = "tbIOTThingName"
	tbIOTRootCAFilename     = "CaCertPath"
	tbIOTCertFilename       = "MQTTCert"
	tbIOTPrivateKeyFilename = "MQTTKey"
	tbSkipCertVerify        = "SkipCertVerify"
	tbPersistOnError        = "PersistOnError"
	tbDeviceNames           = "TbDeviceNames"
)

var log logger.LoggingClient

// TBMQTTConfig holds TB IoT specific information
type TBMQTTConfig struct {
	MQTTConfig     sdkTransforms.MqttConfig
	IoTProtocol    string
	IoTHost        string
	IoTPort        string
	IoTUserName    string
	IoTDevice      string
	IoTTopic       string
	DeviceNames    string
	PersistOnError bool
	KeyCertPair    *sdkTransforms.KeyCertPair
}

func getNewClient(skipVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &http.Client{Timeout: 10 * time.Second, Transport: tr}
}

func getAppSetting(settings map[string]string, name string) string {
	value, ok := settings[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""

}

// LoadTBMQTTConfig Loads the mqtt configuration necessary to connect to Thingsboard
func LoadTBMQTTConfig(sdk *appsdk.AppFunctionsSDK) (*TBMQTTConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	var IoTProtocol, ioTHost, iotPort, IoTUserName, iotDevice, mqttCert, mqttKey, ioTTopic, deviceNames string
	var skipCertVerify, persistOnError bool
	var errSkip, errPersist error

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		IoTProtocol = getAppSetting(appSettings, tbIOTMQTTProtocol)
		ioTHost = getAppSetting(appSettings, tbIOTMQTTHost)
		iotPort = getAppSetting(appSettings, tbIOTMQTTPort)
		IoTUserName = getAppSetting(appSettings, tbUserName)
		iotDevice = getAppSetting(appSettings, tbIOTThingName)
		mqttCert = getAppSetting(appSettings, tbIOTCertFilename)
		mqttKey = getAppSetting(appSettings, tbIOTPrivateKeyFilename)
		ioTTopic = getAppSetting(appSettings, topic)
		deviceNames = getAppSetting(appSettings, tbDeviceNames)
		skipCertVerify, errSkip = strconv.ParseBool(getAppSetting(appSettings, tbSkipCertVerify))
		persistOnError, errPersist = strconv.ParseBool(getAppSetting(appSettings, tbPersistOnError))

		if errSkip != nil {
			log.Error("Unable to parse " + tbSkipCertVerify + " value")
		}
		if errPersist != nil {
			log.Error("Unable to parse " + tbPersistOnError + " value")
		}
	} else {
		return nil, errors.New("No application-specific settings found")
	}

	config := TBMQTTConfig{}
	config.IoTProtocol = IoTProtocol
	config.IoTHost = ioTHost
	config.IoTPort = iotPort
	config.IoTUserName = IoTUserName
	config.IoTDevice = iotDevice
	config.IoTTopic = ioTTopic
	config.DeviceNames = deviceNames
	config.PersistOnError = persistOnError

	pair := &sdkTransforms.KeyCertPair{
		KeyFile:  mqttKey,
		CertFile: mqttCert,
	}

	mqttConfig := sdkTransforms.MqttConfig{
		SkipCertVerify: skipCertVerify,
	}

	log.Debug(fmt.Sprintf("Read SkipCertVerify from configuration: %t", config.MQTTConfig.SkipCertVerify))
	log.Debug(fmt.Sprintf("Read PersistOnError from configuration: %t", config.PersistOnError))

	config.KeyCertPair = pair
	config.MQTTConfig = mqttConfig

	return &config, nil
}

// NewTBMQTTSender return a mqtt sender capable of sending the event's value to the given MQTT broker
func NewTBMQTTSender(logging logger.LoggingClient, config *TBMQTTConfig) *sdkTransforms.MQTTSender {

	logging.Debug(config.IoTTopic)

	port, err := strconv.Atoi(config.IoTPort)
	if err != nil {
		// falling back to default TB IoT port
		port = 1883
	}

	addressable := models.Addressable{
		Protocol:  config.IoTProtocol,
		Address:   config.IoTHost,
		Port:      port,
		Publisher: config.IoTDevice,
		User:      config.IoTUserName,
		Password:  "",
		Topic:     config.IoTTopic,
	}

	mqttSender := sdkTransforms.NewMQTTSender(logging, addressable, config.KeyCertPair, config.MQTTConfig, config.PersistOnError)

	return mqttSender
}
