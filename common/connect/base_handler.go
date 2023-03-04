package connect

import (
	"github.com/golang/protobuf/proto"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"sync"
	"time"
)

var _ MessageHandler = (*Connection)(nil)

// Connection 改造了 MessageHandler 中的核心函数 OnMessage, 并且为 send 函数添加了同异步发送, 超时功能。
// 其余函数委托给 ConnectionHandler 实现。
type Connection struct {
	Zmq *ZmqEndpoint
	ConnectionHandler

	spins sync.Map
}

func NewConnection(zmq *ZmqEndpoint, handler ConnectionHandler) *Connection {
	s := &Connection{Zmq: zmq, ConnectionHandler: handler}
	s.Zmq.SetHandler(s)
	return s
}

func (c *Connection) GetZmq() *ZmqEndpoint {
	return c.Zmq
}

func (c *Connection) Connect() common_type.PluginError {
	return c.Zmq.Connect()
}

// Send 同步发送
func (c *Connection) Send(msg *protocol.PlatformMessage, timeout time.Duration) (
	result *protocol.PlatformMessage, err common_type.PluginError) {
	spin := newSpin(msg.Header.SeqNo, msg, timeout, nil)
	c.spins.Store(spin.id, spin)

	defer c.spins.Delete(spin.id)

	if err = c.SendOnly(msg); err != nil {
		return nil, err
	}

	spin.wait()

	return spin.result, spin.err
}

// SendAsync 异步发送
func (c *Connection) SendAsync(msg *protocol.PlatformMessage, timeout time.Duration, callback CallBack) {
	spin := newSpin(msg.Header.SeqNo, msg, timeout, callback)
	c.spins.Store(spin.id, spin)

	go func() {
		spin.wait()
		c.spins.Delete(spin.id)
	}()

	if err := c.SendOnly(msg); err != nil {
		spin.onError(err)
	}
}

func (c *Connection) SendOnly(msg *protocol.PlatformMessage) (err common_type.PluginError) {
	msgBytes, e := proto.Marshal(msg)
	if e != nil {
		log.ErrorDetails(errors.Trace(e))
		return common_type.NewPluginError(common_type.ProtoMarshalFailure, common_type.ProtoMarshalFailureError.Error(), e.Error())
	}
	if e = c.Zmq.Send(msgBytes); e != nil {
		log.ErrorDetails(errors.Trace(e))
		return common_type.NewPluginError(common_type.EndpointSendErr, common_type.EndpointSendError.Error(), e.Error())
	}
	return nil
}

func (c *Connection) OnMessage(endpoint *EndpointInfo, content []byte) {
	msg := &protocol.PlatformMessage{}
	var e common_type.PluginError
	if err := proto.Unmarshal(content, msg); err != nil {
		log.ErrorDetails(errors.Trace(e))
		e = common_type.NewPluginError(common_type.ProtoUnmarshalFailure, common_type.ProtoUnmarshalFailureError.Error(), err.Error())
	}
	if spin, ok := c.spins.Load(msg.Header.SeqNo); ok {
		spin.(*syncSpin).onResult(msg)
	}
	c.ConnectionHandler.OnMsg(endpoint, msg, e)
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
