package ability

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"sync"
)

type instanceAbilities struct {
	instanceID string
	Abilities  map[string]*types.Ability // map[abilityID]*common.Ability
}

type Ability struct {
	m sync.Map // map[instanceID]*instanceAbilities
}

func NewAbility() *Ability {
	return &Ability{}
}

// TODO 检测是否完整实现了Ability的所有函数
func (a *Ability) Check(abilities []*types.Ability) bool {
	return true
}

func (a *Ability) Register(instanceID string, abilities []*types.Ability) {
	abilitiesMap := make(map[string]*types.Ability)
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

func (a *Ability) Search(instanceID, abilityID string) (*types.Ability, error) {
	ins, ok := a.m.Load(instanceID)
	if !ok {
		return nil, fmt.Errorf("instance %s has no ability", instanceID)
	}
	abs := ins.(*instanceAbilities).Abilities
	ab, ok := abs[abilityID]
	if !ok {
		return nil, fmt.Errorf("instance %s has no such ability: %s", instanceID, abilityID)
	}
	return ab, nil
}

func (a *Ability) GetConfig(instanceID, abilityID string) (map[string]string, error) {
	ab, err := a.Search(instanceID, abilityID)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return ab.Config, nil
}

func (a *Ability) Execute(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) (chan *common_type.AbilityResponse, error) {
	ab, err := a.Search(instanceID, abilityID)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if ab.Type != abilityType {
		return nil, fmt.Errorf("ability id or type err: %s-%s", abilityID, abilityType)
	}
	function, ok := ab.Function[abilityFuncKey]
	if !ok {
		return nil, fmt.Errorf("ability has no such function: %s", abilityFuncKey)
	}

	c := handler.CallPluginFunction(instanceID, abilityID, abilityType, function, arg)
	return c, nil
}

func (a *Ability) SyncExecute(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) ([]byte, error) {
	respChan, err := a.Execute(instanceID, abilityID, abilityType, abilityFuncKey, arg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	resp := <-respChan
	if resp.Err != nil {
		return resp.Data, fmt.Errorf(resp.Err.Error() + resp.Err.Msg())
	}
	return resp.Data, nil
}
