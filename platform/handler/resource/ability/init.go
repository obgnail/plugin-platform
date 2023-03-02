package ability

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"io/ioutil"
	"net/http"
)

type Ability struct {
	source   *protocol.PlatformMessage
	distinct *protocol.PlatformMessage
	ability  string
	mapper   RouteMapper
}

func NewAbility(sourceMessage, distinctMessage *protocol.PlatformMessage) *Ability {
	ability := &Ability{
		source:   sourceMessage,
		distinct: distinctMessage,
		ability:  sourceMessage.GetResource().GetAbility().GetAbility(),
		mapper:   &DefaultRouMapper{},
	}
	return ability
}

func (a *Ability) Execute() {
	var content []byte
	var err error

	defer a.buildMsg(content, err)

	req := a.source.GetResource().GetAbility().GetContent()
	reqObj, err := a.mapper.Map(a.ability, req)
	if err != nil {
		return
	}

	respObj, err := new(http.Client).Do(reqObj)
	if err != nil {
		return
	}
	defer respObj.Body.Close()
	content, err = ioutil.ReadAll(respObj.Body)
	if err != nil {
		return
	}
}

func (a *Ability) buildMsg(content []byte, err error) {
	msg := &protocol.AbilityMessage{
		Ability: a.ability,
		Content: content,
	}
	if err != nil {
		log.ErrorDetails(err)
		e := common_type.NewPluginError(common_type.CallAbilityFailure, err.Error(),
			common_type.CallAbilityError.Error())
		msg.Error = message.BuildErrorMessage(e)
	}
	a.distinct.Resource.Ability = msg
}
