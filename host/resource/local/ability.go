package local

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.Ability = (*Ability)(nil)

type Ability struct {
	LayoutCard common_type.LayoutCard
	Notify     common_type.Notify
}

func NewAbility(plugin common_type.IPlugin) common_type.Ability {
	//return &Ability{
	//	LayoutCard: NewLayoutCard(Plugin),
	//	Notify:     NewNotify(Plugin),
	//	Field:      NewField(Plugin),
	//}
	return nil
}

func (a *Ability) GetNotify() common_type.Notify {
	return a.Notify
}

func (a *Ability) GetLayoutCard() common_type.LayoutCard {
	return a.LayoutCard
}
