package component

import "fmt"

type byType struct {
	components map[Type][]*Component
}

func newByType() *byType {
	return &byType{
		components: make(map[Type][]*Component),
	}
}

func (b byType) List(t Type) []*Component {
	return b.components[t]
}

func (b byType) Get(t Type, id uint) (*Component, error) {
	for _, c := range b.components[t] {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("component not found: %d", id)
}

func (b byType) Add(c *Component) {
	b.components[c.Type] = append(b.components[c.Type], c)
}

func (b byType) Update(n *Component) error {
	for i, c := range b.components[n.Type] {
		if n.ID == c.ID {
			b.components[n.Type][i] = n
			return nil
		}
	}
	return fmt.Errorf("component not found: %d", n.ID)
}

func (b byType) Remove(t Type, id uint) error {
	for i, c := range b.components[t] {
		if c.ID == id {
			b.components[t] = append(b.components[t][:i], b.components[t][i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("component not found: %d", id)
}
