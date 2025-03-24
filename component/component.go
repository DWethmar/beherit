package component

import (
	"fmt"

	"github.com/dwethmar/beherit/entity"
)

//go:generate go tool stringer -type=Type
type Type string

type Component struct {
	ID     uint
	Entity entity.Entity
	Type   Type
	Data   any
}

// Manager is a component manager.
type Manager struct {
	byEntity  *byEntity
	byType    *byType
	factories map[Type]func() *Component
}

// NewManager creates a new component manager.
func NewManager() *Manager {
	return &Manager{
		byEntity:  newByEntity(),
		byType:    newByType(),
		factories: make(map[Type]func() *Component),
	}
}

// Register registers a component factory.
func (cm *Manager) Register(factory func() *Component) {
	cm.factories[factory().Type] = factory
}

// Create creates a new component.
func (cm *Manager) Create(entity entity.Entity, componentType Type) (Component, error) {
	factory, ok := cm.factories[componentType]
	if !ok {
		return Component{}, fmt.Errorf("factory not found for component type %s", componentType)
	}
	component := Component{
		ID:     uint(len(cm.byType.List(componentType)) + 1),
		Entity: entity,
		Type:   componentType,
		Data:   factory(),
	}
	return component, nil
}

// Update updates a component.
func (cm *Manager) Update(c Component) error {
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
