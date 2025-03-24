package entity

type Entity uint

type Manager struct {
	nextEntity Entity
}

func NewManager(nextEntity Entity) *Manager {
	return &Manager{
		nextEntity: nextEntity,
	}
}

func (em *Manager) Create() Entity {
	entity := em.nextEntity
	em.nextEntity++
	return entity
}
