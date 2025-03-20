package beherit

import "fmt"

type ComponentType string

func (c ComponentType) String() string { return string(c) }

type Component struct {
	ID     uint
	Entity Entity
	Type   ComponentType
	Data   any
}

type ComponentManager struct {
	components map[ComponentType][]Component
	factories  map[ComponentType]func() any
}

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make(map[ComponentType][]Component),
		factories:  make(map[ComponentType]func() any),
	}
}

func (cm *ComponentManager) Register(componentType ComponentType, factory func() any) {
	cm.factories[componentType] = factory
}

func (cm *ComponentManager) Create(entity Entity, componentType ComponentType) (Component, error) {
	factory, ok := cm.factories[componentType]
	if !ok {
		return Component{}, fmt.Errorf("factory not found for component type %s", componentType)
	}
	component := Component{
		ID:     uint(len(cm.components[componentType])),
		Entity: entity,
		Type:   componentType,
		Data:   factory(),
	}
	cm.components[componentType] = append(cm.components[componentType], component)
	return component, nil
}

func (cm *ComponentManager) Update(c Component) error {
	components := cm.components[c.Type]
	for i, component := range components {
		if component.ID == c.ID {
			cm.components[c.Type][i] = c
			return nil
		}
	}
	return fmt.Errorf("component %d not found", c.ID)
}

func (cm *ComponentManager) Get(componentType ComponentType, id uint) (Component, error) {
	components := cm.components[componentType]
	for _, component := range components {
		if component.ID == id {
			return component, nil
		}
	}
	return Component{}, fmt.Errorf("component %d not found", id)
}

func (cm *ComponentManager) Components(componentType ComponentType) []Component {
	return cm.components[componentType]
}

func (cm *ComponentManager) Remove(componentType ComponentType, id uint) {
	components := cm.components[componentType]
	for i, component := range components {
		if component.ID == id {
			cm.components[componentType] = append(components[:i], components[i+1:]...)
			break
		}
	}
}
