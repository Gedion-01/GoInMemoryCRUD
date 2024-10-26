package server

import (
	"sync"

	"github.com/Gedion-01/Go-Crud-Challenge/types"
)

type memoryDB struct {
	persons []types.Person
	mu      sync.RWMutex
}

func newDB() *memoryDB {
	return &memoryDB{
		persons: []types.Person{},
	}
}

func (m *memoryDB) set(person *types.Person) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newPerson := *person
	m.persons = append(m.persons, newPerson)
}

func (m *memoryDB) get(key string) (*types.Person, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, person := range m.persons {
		if person.ID == key {
			return &person, true
		}
	}
	return nil, false
}

func (m *memoryDB) delete(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, person := range m.persons {
		if person.ID == key {
			m.persons = append(m.persons[:i], m.persons[i+1:]...)
			return true
		}
	}
	return false
}

func (m *memoryDB) update(key string, person *types.UpdatePersonParams) (*types.Person, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, p := range m.persons {
		if p.ID == key {
			if person.Name != "" {
				m.persons[i].Name = person.Name
			}
			if person.Age != "" {
				m.persons[i].Age = person.Age
			}
			if person.Hobbies != nil {
				m.persons[i].Hobbies = person.Hobbies
			}
			return &m.persons[i], true
		}
	}
	return nil, false
}

func (m *memoryDB) all() []types.Person {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.persons
}
