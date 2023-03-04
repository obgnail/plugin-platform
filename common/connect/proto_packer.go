package connect

import (
	"github.com/golang/protobuf/proto"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var _ MessagePacker = (*ProtoPacker)(nil)

type ProtoPacker struct{}

func (i *ProtoPacker) Unpack(data []byte) (*EndpointInfo, *EndpointInfo, []byte, common_type.PluginError) {
	if len(data) == 0 {
		return nil, nil, nil, genError("empty msg")
	}

	message := &protocol.PlatformMessage{}
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, nil, nil, genError(err.Error())
	}

	source := getEndpointInfo(message.Header.Source)
	target := getEndpointInfo(message.Header.Distinct)
	return source, target, data, nil
}

func (i *ProtoPacker) Pack(source, target *EndpointInfo, data []byte) ([]byte, common_type.PluginError) {
	msg := &protocol.PlatformMessage{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, genError(err.Error())
	}

	msg.Header.Source.ID = source.ID
	msg.Header.Distinct.ID = target.ID

	newData, err := proto.Marshal(msg)
	if err != nil {
		return nil, genError(err.Error())
	}
	return newData, nil
}

func getEndpointInfo(node *protocol.RouterNode) *EndpointInfo {
	ret := &EndpointInfo{
		ID:   node.ID,
		Name: node.Tags["name"],
		Role: Role(node.Tags["role"]),
	}
	return ret
}

func genError(err string) common_type.PluginError {
	return common_type.NewPluginError(common_type.ProtoUnmarshalFailure, err)
}
