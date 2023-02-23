package connect

import (
	"github.com/golang/protobuf/proto"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"sync"
	"time"
)

var _ MessageHandler = (*BaseHandler)(nil)

// BaseHandler 改造了 MessageHandler 中的核心函数 OnMessage, 并且为 send 函数添加了同异步发送, 超时功能。
// 其余函数委托给 FurtherHandler 实现。
type BaseHandler struct {
	Zmq *ZmqEndpoint
	FurtherHandler

	spins sync.Map
}

func NewBaseHandler(zmq *ZmqEndpoint, handler FurtherHandler) *BaseHandler {
	s := &BaseHandler{Zmq: zmq, FurtherHandler: handler}
	s.Zmq.SetHandler(s)
	return s
}

func (s *BaseHandler) GetZmq() *ZmqEndpoint {
	return s.Zmq
}

// Send 同步发送
func (s *BaseHandler) Send(msg *protocol.PlatformMessage, timeout time.Duration) (
	result *protocol.PlatformMessage, err common_type.PluginError) {

	spin := newSpin(msg.Header.SeqNo, msg, timeout, nil)
	s.spins.Store(spin.id, spin)
	defer s.spins.Delete(spin.id)

	msgBytes, e := proto.Marshal(msg)
	if e != nil {
		err = common_type.NewPluginError(common_type.ProtoMarshalFailure, common_type.ProtoMarshalFailureError.Error(), e.Error())
		return nil, err
	}
	if e = s.Zmq.Send(msgBytes); e != nil {
		err = common_type.NewPluginError(common_type.EndpointSendErr, common_type.EndpointSendError.Error(), e.Error())
		return nil, err
	}

	spin.wait()

	return spin.result, spin.err
}

// SendAsync 异步发送
func (s *BaseHandler) SendAsync(msg *protocol.PlatformMessage, timeout time.Duration, callback CallBack) {

	spin := newSpin(msg.Header.SeqNo, msg, timeout, callback)
	s.spins.Store(spin.id, spin)

	go func() {
		spin.wait()
		s.spins.Delete(spin.id)
	}()

	msgBytes, e := proto.Marshal(msg)
	if e != nil {
		err := common_type.NewPluginError(common_type.ProtoMarshalFailure, common_type.ProtoMarshalFailureError.Error(), e.Error())
		spin.onError(err)
		return
	}
	if e = s.Zmq.Send(msgBytes); e != nil {
		err := common_type.NewPluginError(common_type.EndpointSendErr, common_type.EndpointSendError.Error(), e.Error())
		spin.onError(err)
		return
	}
}

func (s *BaseHandler) OnMessage(endpoint *EndpointInfo, content []byte) {
	msg := &protocol.PlatformMessage{}
	var e common_type.PluginError
	if err := proto.Unmarshal(content, msg); err != nil {
		e = common_type.NewPluginError(common_type.ProtoUnmarshalFailure, common_type.ProtoUnmarshalFailureError.Error(), err.Error())
	}
	if spin, ok := s.spins.Load(msg.Header.SeqNo); ok {
		spin.(*syncSpin).onResult(msg)
	}
	s.FurtherHandler.OnMsg(endpoint, msg, e)
}

type CallBack func(input, result *protocol.PlatformMessage, err common_type.PluginError)

// syncSpin: onError、onResult、onTimeout时执行回调
type syncSpin struct {
	id uint64

	input  *protocol.PlatformMessage
	result *protocol.PlatformMessage
	err    common_type.PluginError

	timeout  time.Duration
	callback CallBack

	exit chan struct{}
}

func newSpin(id uint64, input *protocol.PlatformMessage, timeout time.Duration, callback CallBack) *syncSpin {
	return &syncSpin{
		id:       id,
		input:    input,
		result:   nil,
		err:      nil,
		timeout:  timeout,
		callback: callback,
		exit:     make(chan struct{}, 1),
	}
}

func (s *syncSpin) onError(err common_type.PluginError) {
	s.err = err
	s.exit <- struct{}{}
}

func (s *syncSpin) onResult(msg *protocol.PlatformMessage) {
	s.result = msg
	s.exit <- struct{}{}
}

func (s *syncSpin) wait() {
	select {
	case <-time.After(s.timeout):
		s.err = common_type.NewPluginError(common_type.MsgTimeOut, common_type.MsgTimeOutError.Error(), "timeout")
	case <-s.exit:
	}

	if s.callback != nil {
		go s.callback(s.input, s.result, s.err)
	}
}
