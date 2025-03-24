package component

import (
	"fmt"

	"github.com/dwethmar/beherit/entity"
)

type byEntity struct {
	components map[entity.Entity][]*Component
}

func newByEntity() *byEntity {
	return &byEntity{
		components: make(map[entity.Entity][]*Component),
	}
}

func (b byEntity) List(e entity.Entity) []*Component {
	return b.components[e]
}

func (b byEntity) Add(c *Component) {
	b.components[c.Entity] = append(b.components[c.Entity], c)
}

func (b byEntity) Update(n *Component) error {
	for i, c := range b.components[n.Entity] {
		if n.ID == c.ID {
			b.components[c.Entity][i] = n
			return nil
		}
	}
	return fmt.Errorf("component not found: %d", n.ID)
}

func (b byEntity) Remove(e entity.Entity, id uint) error {
	for i, c := range b.components[e] {
		if c.ID == id {
			b.components[e] = append(b.components[e][:i], b.components[e][i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("component not found: %d", id)
}
