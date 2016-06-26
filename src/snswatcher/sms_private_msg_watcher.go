package main

import (
	"encoding/json"
	"jzlservice/smssender"
	"strings"
	"time"
	"unicode"
)

func (this *WatcherMgr) pullSMSPrivateMsgWorker() {
	var result *smssender.SMSMOMessage
	var err error

	smsProviders, ok := g_config.Get("service.smssender.watcher")
	if !ok || smsProviders == "" {
		LOG_ERROR("没有设置要监控的第三方短信服务，负责拉取上行短信的工作线程将退出")
		return
	}

	smsProvider_list := strings.FieldsFunc(smsProviders, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, smsProvider := range smsProvider_list {
				smsSender := g_snsClient.GetSMSSender(smsProvider)
				if smsSender != nil {
					category_list, _ := smsSender.GetCategory()
					for _, category := range category_list {
						result, err = smsSender.GetMOMessage(category)
						if err == nil {
							err = this.processSMSPrivateMsg(result)
							if err != nil {
								LOG_ERROR("处理上行短信失败。失败原因：%v", err)
							}
						}
					}
				}
			}
		}
	}
}

func (this *WatcherMgr) processSMSPrivateMsg(result *smssender.SMSMOMessage) error {
	if result.Status == "0" {
		if result.Message != "0" {
			data, err := json.Marshal(result.Data)
			if err != nil {
				LOG_ERROR("序列化上行短信数据失败，失败原因: %v", err)
				return err
			}

			topic := "SMSPrivateMsg"
			return g_nsqProducer.Publish(topic, string(data))
		}
	} else {
		LOG_INFO("上行短信数据有误，错误信息: %v", result.Message)
	}
	return nil
}
