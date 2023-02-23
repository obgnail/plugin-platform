package local

import (
	common "github.com/obgnail/plugin-platform/common/common_type"
)

var _ common.Ability = (*Ability)(nil)

type Ability struct {
	LayoutCard common.LayoutCard
	Notify     common.Notify
	Field      common.Field
}

func NewAbility(plugin common.IPlugin) common.Ability {
	//return &Ability{
	//	LayoutCard: NewLayoutCard(plugin),
	//	Notify:     NewNotify(plugin),
	//	Field:      NewField(plugin),
	//}
	return nil
}

func (a *Ability) GetNotify() common.Notify {
	return a.Notify
}

func (a *Ability) GetLayoutCard() common.LayoutCard {
	return a.LayoutCard
}

func (a *Ability) GetField() common.Field {
	return a.Field
}
