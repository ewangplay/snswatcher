package main

import ()

const (
	RETRY_MAX_COUNT = 3
)

type SNSWatcherImpl struct {
}

func (this *SNSWatcherImpl) Ping() (r string, err error) {
	LOG_INFO("请求ping方法")
	return "pong", nil
}
