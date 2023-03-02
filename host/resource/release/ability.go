package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.Ability = (*Ability)(nil)

type Ability struct {
	plugin common_type.IPlugin
	sender common.Sender
}

func NewAbility(plugin common_type.IPlugin, sender common.Sender) common_type.Ability {
	return &Ability{plugin: plugin, sender: sender}
}

func (a *Ability) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return a.sender.Send(a.plugin, platformMessage)
}

func (a *Ability) buildMessage(abilityMessageMessage *protocol.AbilityMessage) *protocol.PlatformMessage {
	msg := message.GetInitMessage(nil, nil)
	msg.Resource = &protocol.ResourceMessage{Ability: abilityMessageMessage}
	return msg
}

func (a *Ability) Call(kind string, req []byte) (result []byte, err common_type.PluginError) {
	msg := &protocol.AbilityMessage{
		Ability: kind,
		Content: req,
		Error:   nil,
	}

	resp, err := a.sendMsgToHost(a.buildMessage(msg))
	if err != nil {
		return nil, err
	}
	result = resp.GetResource().GetAbility().GetContent()
	if e := resp.GetResource().GetAbility().GetError(); e != nil {
		return result, common_type.NewPluginError(int(e.Code), e.GetError(), e.GetMsg())
	}
	return result, nil
}
