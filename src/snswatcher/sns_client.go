package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"jzlservice/smssender"
	"strings"
	"unicode"
)

type SMSSenderClient struct {
	smsProvider string
}

func (this *SMSSenderClient) GetReport(category int16) (r *smssender.SMSReport, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_DEBUG("获取发送短信的状态报告开始...")

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		LOG_ERROR("没有提供短信服务提供商")
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败")
		return nil, err
	}

	r, err = client.GetReport(category)
	if err != nil {
		LOG_ERROR("获取发送短信的状态报告失败。 失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取发送短信的状态报告成功。状态报告为：%v", r)

	return r, nil
}

func (this *SMSSenderClient) GetMOMessage(category int16) (r *smssender.SMSMOMessage, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取上行的短信开始...")

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = fmt.Sprintf("没有提供短信服务提供商")
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败")
		return nil, err
	}

	r, err = client.GetMOMessage(category)
	if err != nil {
		LOG_ERROR("获取上行的短信失败。 失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取上行的短信成功。上行短信：%v", r)

	return r, nil
}

func (this *SMSSenderClient) GetCategory() (r []int16, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取短信通道类别开始...")

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = fmt.Sprintf("没有提供短信服务提供商")
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败")
		return nil, err
	}

	r, err = client.GetCategory()
	if err != nil {
		LOG_ERROR("获取短信通道类别失败。 失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取短信类别列表成功：%v", r)

	return r, nil
}

type SNSClient struct {
	smsSenderClient map[string]*SMSSenderClient
}

func NewSNSClient() (*SNSClient, error) {
	snsClient := &SNSClient{}

	snsClient.smsSenderClient = make(map[string]*SMSSenderClient)
	channels, err := GetSMSChannels()
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		snsClient.smsSenderClient[channel] = &SMSSenderClient{smsProvider: channel}
	}

	return snsClient, nil
}

func (this *SNSClient) GetSMSSender(provider string) *SMSSenderClient {
	r, ok := this.smsSenderClient[provider]
	if ok && r != nil {
		return r
	} else {
		return nil
	}
}

func GetSMSChannels() (r []string, err error) {
	var outputStr string

	smsProviders, ok := g_config.Get("service.smssender.watcher")
	if !ok || smsProviders == "" {
		outputStr = fmt.Sprintf("没有设置要监控的第三方短信服务")
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	r = strings.FieldsFunc(smsProviders, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	return
}
