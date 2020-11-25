module thingsboard-export

go 1.14

require (
	github.com/alecthomas/gozmq v0.0.0-20140622232202-d1b01a2df6b2
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/edgexfoundry/app-functions-sdk-go v1.2.0
	github.com/edgexfoundry/go-mod-core-contracts v0.1.58
	github.com/pebbe/zmq4 v1.0.0
	github.com/robfig/cron v1.2.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/zhanglt/thingsboard-export/transforms v0.0.0
	gopkg.in/zeromq/goczmq.v4 v4.1.0 // indirect

)

replace github.com/zhanglt/thingsboard-export/transforms => ./transforms
