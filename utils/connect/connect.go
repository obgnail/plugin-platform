package connect

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/utils/errors"
	"sync"
)

type EndpointInfo struct {
	ID   string
	Role Role
	Name string
}

type zmqMessage struct {
	endpointID string
	data       []byte
}

type ZmqEndpoint struct {
	address    string
	info       *EndpointInfo
	socketType SocketType

	packer  MessagePacker
	handler MessageHandler

	socket zmq4.Socket

	ctx    context.Context
	cancel context.CancelFunc

	sendChan chan *zmqMessage

	endpoints sync.Map // map[ID]*EndpointInfo
}

func New(id, name, addr string, socketType SocketType, role Role) *ZmqEndpoint {
	ctx, cancel := context.WithCancel(context.Background())
	p := &ZmqEndpoint{
		info:       &EndpointInfo{ID: id, Role: role, Name: name},
		address:    addr,
		socketType: socketType,
		ctx:        ctx,
		cancel:     cancel,
		sendChan:   make(chan *zmqMessage, 1024),
	}
	return p
}

func (p *ZmqEndpoint) SetPacker(packer MessagePacker) *ZmqEndpoint {
	p.packer = packer
	return p
}

func (p *ZmqEndpoint) SetHandler(h MessageHandler) *ZmqEndpoint {
	p.handler = h
	return p
}

func (p *ZmqEndpoint) Close() error {
	p.cancel()
	if p.socket != nil {
		if err := p.socket.Close(); err != nil {
			return errors.Trace(err)
		}
	}
	p.socket = nil
	p.endpoints.Range(func(key, value any) bool {
		p.endpoints.Delete(key)
		return true
	})
	p.handler.OnDisconnect()
	return nil
}

func (p *ZmqEndpoint) GetEndpoint(id string) (*EndpointInfo, bool) {
	end, ok := p.endpoints.Load(id)
	if !ok {
		return nil, false
	}
	return end.(*EndpointInfo), true
}

func (p *ZmqEndpoint) AddEndpoint(ep *EndpointInfo) {
	p.endpoints.Store(ep.ID, ep)
}

func (p *ZmqEndpoint) DeleteEndpoint(id string) {
	p.endpoints.Delete(id)
}

func (p *ZmqEndpoint) ExistEndpoint(id string) bool {
	_, ok := p.endpoints.Load(id)
	return ok
}

func (p *ZmqEndpoint) ListEndpoints() []*EndpointInfo {
	var result []*EndpointInfo
	p.endpoints.Range(func(key, value any) bool {
		result = append(result, value.(*EndpointInfo))
		return true
	})
	return result
}

func (p *ZmqEndpoint) Send(rawData []byte) common.PluginError {
	_, target, processData, err := p.packer.Unpack(rawData)
	if err != nil {
		return err
	}
	p.sendChan <- &zmqMessage{endpointID: target.ID, data: processData}
	return nil
}

func (p *ZmqEndpoint) SendTo(id string, content []byte) common.PluginError {
	target, ok := p.GetEndpoint(id)
	if !ok {
		errStr := fmt.Sprintf("%s SendMsgError(getTarget by endpointID[%s]): TargetEndpointNotFound", p.info, id)
		return common.NewPluginError(common.TargetEndpointNotFound, common.TargetEndpointNotFoundError.Error(), errStr)
	}
	p.sendChan <- &zmqMessage{endpointID: target.ID, data: content}
	return nil
}

func (p *ZmqEndpoint) Publish(rawData []byte) common.PluginError {
	for _, target := range p.ListEndpoints() {
		processData, err := p.packer.Pack(p.info, target, rawData)
		if err != nil {
			return err
		}
		if err = p.SendTo(target.ID, processData); err != nil {
			return err
		}
	}
	return nil
}

func (p *ZmqEndpoint) Connect() common.PluginError {
	if p.handler == nil {
		return common.NewPluginError(common.SocketListenOrDialFailure,
			common.SocketListenOrDialFailureError.Error(), "hasNoHandler")
	}
	if p.packer == nil {
		return common.NewPluginError(common.SocketListenOrDialFailure,
			common.SocketListenOrDialFailureError.Error(), "hasNoIdentifier")
	}

	var err error
	switch p.socketType {
	case SocketTypeRouter:
		p.socket = zmq4.NewRouter(p.ctx, zmq4.WithID(zmq4.SocketIdentity(p.info.ID)))
		err = p.socket.Listen(p.address)
	case SocketTypeDealer:
		p.socket = zmq4.NewDealer(p.ctx, zmq4.WithID(zmq4.SocketIdentity(p.info.ID)))
		err = p.socket.Dial(p.address)
	}
	if err != nil {
		p.Close()
		return common.NewPluginError(common.SocketListenOrDialFailure,
			common.SocketListenOrDialFailureError.Error(), "socketListenOrDialFailureError")
	}

	go p.startReceiver()
	go p.startSender()
	go p.handler.OnConnect()
	return nil
}

func (p *ZmqEndpoint) startReceiver() {
	for {
		select {
		case <-p.ctx.Done():
			break
		default:
			raw, err := p.socket.Recv()
			if err != nil {
				p.handler.OnError(common.NewPluginError(common.EndpointReceiveErr,
					common.EndpointReceiveError.Error(), err.Error()))
				continue
			}
			rawData := raw.Frames[1]
			source, _, processData, err := p.packer.Unpack(rawData)
			if err != nil {
				p.handler.OnError(common.NewPluginError(common.EndpointIdentifyErr,
					common.EndpointIdentifyError.Error(), err.Error()))
				continue
			}
			p.AddEndpoint(source)
			msg := &zmqMessage{endpointID: string(raw.Frames[0]), data: processData}
			go p.handler.OnMessage(source, msg.data)
		}
	}
}

func (p *ZmqEndpoint) startSender() {
	for {
		select {
		case <-p.ctx.Done():
			break
		case msg := <-p.sendChan:
			content := zmq4.NewMsgFrom([]byte(msg.endpointID), msg.data)
			err := p.socket.Send(content)
			if err != nil {
				p.handler.OnError(common.NewPluginError(common.EndpointSendErr,
					common.EndpointSendError.Error(), err.Error()))
				continue
			}
		}
	}
}
