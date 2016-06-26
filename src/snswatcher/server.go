package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/ewangplay/jzlconfig"
	"github.com/outmana/log4jzl"
	"jzlservice/snswatcher"
	"os"
)

//global object
var g_config jzlconfig.JZLConfig
var g_logger *log4jzl.Log4jzl
var g_snsClient *SNSClient
var g_nsqProducer *NSQProducer
var g_watcherMgr *WatcherMgr

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [--config path_to_config_file]")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	var err error

	//parse command line
	var configFile string
	flag.Usage = Usage
	flag.StringVar(&configFile, "config", "snswatcher.conf", "specified config filename")
	flag.Parse()

	fmt.Println("config file: ", configFile)

	//read config file
	if err = g_config.Read(configFile); err == nil {
		fmt.Println(g_config)
	} else {
		fmt.Println("Read config file fail.", err)
		os.Exit(1)
	}

	//init logger
	g_logger, err = log4jzl.New("snswatcher")
	if err != nil {
		fmt.Println("Open log file fail.", err)
		os.Exit(1)
	}

	//init log level object
	g_logLevel, err = NewLogLevel()
	if err != nil {
		LOG_ERROR("Craete Log level error: %v", err)
		os.Exit(1)
	}

	//init SNS client
	g_snsClient, err = NewSNSClient()
	if err != nil {
		fmt.Println("create SNSClient object fail.", err)
		os.Exit(1)
	}

	//init nsq client
	g_nsqProducer, err = NewNSQProducer()
	if err != nil {
		fmt.Println("create NSQProducer object fail.", err)
		os.Exit(1)
	}
	g_nsqProducer.Init()
	defer g_nsqProducer.Release()

	//init watcher  manager
	g_watcherMgr, err = NewWatcherMgr()
	if err != nil {
		fmt.Println("create WatcherMgr object fail.", err)
		os.Exit(1)
	}
	g_watcherMgr.Init()
	g_watcherMgr.Run()
	defer g_watcherMgr.Release()

	//format the server listening newwork address
	var networkAddr string
	serviceIp, serviceIPIsSet := g_config.Get("service.addr")
	servicePort, servicePortIsSet := g_config.Get("service.port")
	if serviceIPIsSet && servicePortIsSet {
		networkAddr = fmt.Sprintf("%s:%s", serviceIp, servicePort)
	} else {
		networkAddr = "127.0.0.1:19090"
	}

	//startup snswatcher service
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(networkAddr)
	if err != nil {
		fmt.Println("create socket listening fail.", err)
		os.Exit(1)
	}
	handler := &SNSWatcherImpl{}
	processor := snswatcher.NewSNSWatcherProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	fmt.Println("snswatcher server working on", networkAddr)
	LOG_INFO("snswatcher服务启动，监听地址：%v", networkAddr)

	server.Serve()
}
