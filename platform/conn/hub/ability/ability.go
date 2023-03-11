package ability

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"sync"
)

type instanceAbilities struct {
	instanceID string
	Abilities  map[string]*common.Ability // map[abilityID]*common.Ability
}

type Ability struct {
	m sync.Map // map[instanceID]*instanceAbilities
}

func NewAbility() *Ability {
	return &Ability{}
}

// TODO 检测是否完整实现了Ability的所有函数
func (a *Ability) Check(abilities []*common.Ability) bool {
	return true
}

func (a *Ability) Register(instanceID string, abilities []*common.Ability) {
	abilitiesMap := make(map[string]*common.Ability)
	for _, _ability := range abilities {
		abilitiesMap[_ability.Id] = _ability
	}

	instance, ok := a.m.Load(instanceID)
	if !ok {
		i := &instanceAbilities{
			instanceID: instanceID,
			Abilities:  abilitiesMap,
		}
		a.m.Store(instanceID, i)
	} else {
		ab := instance.(*instanceAbilities).Abilities
		for _, _ability := range abilities {
			ab[_ability.Id] = _ability
		}
		i := &instanceAbilities{
			instanceID: instanceID,
			Abilities:  ab,
		}
		a.m.Store(instanceID, i)
	}
}

func (a *Ability) Cancel(instanceID string) {
	a.m.Delete(instanceID)
}

func (a *Ability) Execute(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) (chan *common_type.AbilityResponse, error) {

	ins, ok := a.m.Load(instanceID)
	if !ok {
		return nil, fmt.Errorf("instance %s has no ability", instanceID)
	}
	abs := ins.(*instanceAbilities).Abilities
	ab, ok := abs[abilityID]
	if !ok {
		return nil, fmt.Errorf("instance %s has no such ability: %s", instanceID, abilityID)
	}

	if ab.AbilityType != abilityType {
		return nil, fmt.Errorf("ability id or type err: %s-%s", abilityID, abilityType)
	}
	function, ok := ab.Function[abilityFuncKey]
	if !ok {
		return nil, fmt.Errorf("ability has no such function: %s", abilityFuncKey)
	}

	c := handler.CallPluginFunction(instanceID, abilityID, abilityType, function, arg)
	return c, nil
}
