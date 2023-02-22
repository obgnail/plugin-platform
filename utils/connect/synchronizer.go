package connect

import "sync"

type SyncSpin struct {
	MsgSeqNo uint64
	Input    []byte
	Result   []byte
	Error    error

	Async       bool
	AsyncObject interface{}

	Handler        func(input, result []byte, asyncObject interface{}, err error)
	TimeoutHandler func(input []byte, asyncObject interface{}, err error)
	WaitPeriod     int64

	resultChan  chan []byte
	timeoutChan chan bool
}

type Synchronized struct {
	timeout int64
	spins   sync.Map
}

func NewSynchronized(timeout int64) *Synchronized {
	return &Synchronized{timeout: timeout}
}

func (s *Synchronized) OnSend(spin *SyncSpin) {

}

func (s *Synchronized) OnMessage(spin *SyncSpin) {

}
