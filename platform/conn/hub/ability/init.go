package ability

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/platform/service/common"
)

var ability *Ability

func InitAbility() {
	ability = NewAbility()
}

func CheckAbility(abilities []*common.Ability) bool {
	return ability.Check(abilities)
}

func RegisterAbility(instanceID string, abilities []*common.Ability) {
	ability.Register(instanceID, abilities)
}

func CancelAbility(instanceID string) {
	ability.Cancel(instanceID)
}

func Execute(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) (chan *common_type.AbilityResponse, error) {
	return ability.Execute(instanceID, abilityID, abilityType, abilityFuncKey, arg)
}
