package main

import ()

type WatcherMgr struct {
}

func NewWatcherMgr() (*WatcherMgr, error) {
	watcherMgr := &WatcherMgr{}

	return watcherMgr, nil
}

func (this *WatcherMgr) Init() error {
	return nil
}

func (this *WatcherMgr) Release() {
}

func (this *WatcherMgr) Run() {
	go this.pullSMSStatusWorker()
	go this.pullSMSPrivateMsgWorker()
}
