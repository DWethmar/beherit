package component

import (
	"fmt"

	"github.com/dwethmar/beherit/entity"
)

//go:generate go tool stringer -type=Type
type Type string

type Component struct {
	ID     uint          `json:"id"`
	Entity entity.Entity `json:"entity"`
	Type   Type          `json:"type"`
	Data   any           `json:"data"`
}

// Manager is a component manager.
type Manager struct {
	byEntity *byEntity
	byType   *byType
	nextID   uint
}

// NewManager creates a new component manager.
func NewManager(nextID uint) *Manager {
	return &Manager{
		byEntity: newByEntity(),
		byType:   newByType(),
		nextID:   nextID,
	}
}

// Create creates a new component.
func (cm *Manager) Add(c Component) error {
	if c.ID == 0 {
		c.ID = cm.nextID
		cm.nextID++
	}
	// check if component already exists
	o, err := cm.byType.Get(c.Type, c.ID)
	if err == nil {
		return fmt.Errorf("component already exists: %d", o.ID)
	}
	cm.byEntity.Add(&c)
	cm.byType.Add(&c)
	return nil
}

// Update updates a component data
func (cm *Manager) Update(c Component) error {
	// check if exists
	_, err := cm.byType.Get(c.Type, c.ID)
	if err != nil {
		return err
	}
	if err := cm.byEntity.Update(&c); err != nil {
		return err
	}
	if err := cm.byType.Update(&c); err != nil {
		return err
	}
	return nil
}

// Get gets a component.
func (cm *Manager) Get(componentType Type, id uint) (Component, error) {
	c, err := cm.byType.Get(componentType, id)
	if err != nil {
		return Component{}, err
	}
	return *c, nil
}

// ListByEntity gets all components of an entity.
func (cm *Manager) ListByEntity(e entity.Entity, componentType Type) []Component {
	l := cm.byEntity.List(e)
	if len(l) == 0 {
		return []Component{}
	}
	result := make([]Component, 0)
	for _, c := range l {
		if c.Type == componentType {
			result = append(result, *c)
		}
	}
	return result
}

// ListByEntity gets all components of an entity.
func (cm *Manager) FirstByEntity(e entity.Entity, componentType Type) (Component, bool) {
	l := cm.byEntity.List(e)
	if len(l) == 0 {
		return Component{}, false
	}
	for _, c := range l {
		if c.Type == componentType {
			return *c, true
		}
	}
	return Component{}, false
}

// List gets all components of a type.
func (cm *Manager) List(componentType Type) []Component {
	l := cm.byType.List(componentType)
	if len(l) == 0 {
		return []Component{}
	}
	result := make([]Component, len(l))
	for i, c := range l {
		result[i] = *c
	}
	return result
}

// Remove removes a component.
func (cm *Manager) Remove(componentType Type, id uint) error {
	c, err := cm.byType.Get(componentType, id)
	if err != nil {
		return err
	}
	if err := cm.byType.Remove(componentType, id); err != nil {
		return err
	}
	if err := cm.byEntity.Remove(c.Entity, id); err != nil {
		return err
	}
	return nil
}

type Factory struct {
	factories map[Type]func() *Component
}

func NewFactory() *Factory {
	return &Factory{
		factories: make(map[Type]func() *Component),
	}
}

func (f *Factory) Register(factory func() *Component) {
	f.factories[factory().Type] = factory
}

func (f *Factory) Create(componentType Type) (*Component, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component type is empty")
	}
	factory, ok := f.factories[componentType]
	if !ok {
		return nil, fmt.Errorf("factory not found for component type %q", componentType)
	}
	return factory(), nil
}
