package beherit

type Entity uint

type EntityManager struct {
	nextEntity Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{}
}

func (em *EntityManager) Create() Entity {
	entity := em.nextEntity
	em.nextEntity++
	return entity
}
