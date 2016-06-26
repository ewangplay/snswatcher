package main

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"strings"
	"time"
)

type NSQProducer struct {
	producer *nsq.Producer
	isReady  bool
}

func NewNSQProducer() (*NSQProducer, error) {
	nsqProducer := &NSQProducer{}
	nsqProducer.producer = nil
	nsqProducer.isReady = false

	return nsqProducer, nil
}

func (this *NSQProducer) Init() error {
	var outputStr string
	var nsqNetAddr string

	serviceIp, serviceIPIsSet := g_config.Get("nsq_service.nsqd.addr")
	servicePort, servicePortIsSet := g_config.Get("nsq_service.nsqd.port")

	if !serviceIPIsSet || !servicePortIsSet {
		outputStr = "没有配置NSQ服务的网络地址"
		LOG_ERROR(outputStr)
		return fmt.Errorf(outputStr)
	}

	nsqNetAddr = fmt.Sprintf("%s:%s", serviceIp, servicePort)

	cfg := nsq.NewConfig()
	p, err := nsq.NewProducer(nsqNetAddr, cfg)
	if err != nil {
		LOG_ERROR("创建NSQ Producer失败，失败原因：%v", err)
		return err
	}

	this.producer = p
	this.isReady = true

	return nil
}

func (this *NSQProducer) Release() {
	if this.isReady {
		if this.producer != nil {
			this.producer.Stop()
		}
		this.isReady = false
	}
}

func (this *NSQProducer) IsReady() bool {
	return this.isReady
}

func (this *NSQProducer) GetProducer() *nsq.Producer {
	if !this.IsReady() {
		if err := this.Init(); err != nil {
			return nil
		}
	}

	return this.producer
}

func (this *NSQProducer) Publish(topic string, msg string) (err error) {
	retry_count := RETRY_MAX_COUNT
	for {
		err = this.GetProducer().Publish(topic, []byte(msg))
		if err != nil {
			LOG_ERROR("发布消息[%v]到NSQ系统失败，失败原因：%v", msg, err)

			if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "connection refused") {
				LOG_ERROR("到NSQ消息系统的连接可能异常断开，进行重新连接: 第%v次", RETRY_MAX_COUNT-retry_count+1)

				this.Release()

				if retry_count > 0 {
					time.Sleep(3 * time.Second)
					retry_count--
					continue
				}
			}
		}

		//succ or fail
		break
	}

	if err != nil {
		LOG_ERROR("发布消息[%v]到NSQ队列Topic[%v]失败，失败原因：%v", msg, topic, err)
	} else {
		LOG_INFO("发布消息[%v]到NSQ队列Topic[%v]成功", msg, topic)
	}

	return
}
