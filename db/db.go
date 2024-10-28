package db

import (
	"sync"

	"github.com/Gedion-01/Go-Crud-Challenge/types"
)

type PersonStore interface {
	Set(person *types.Person)
	Get(key string) (*types.Person, bool)
	Delete(key string) bool
	Update(key string, person *types.CreatePersonParams) (*types.Person, bool)
	All() *[]types.Person
}

type MemoryDB struct {
	persons []types.Person
	mu      sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		persons: []types.Person{},
	}
}

func (m *MemoryDB) Set(person *types.Person) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newPerson := *person
	m.persons = append(m.persons, newPerson)
}

func (m *MemoryDB) Get(key string) (*types.Person, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, person := range m.persons {
		if person.ID == key {
			return &person, true
		}
	}
	return nil, false
}

func (m *MemoryDB) Delete(key string) bool {
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

func (m *MemoryDB) Update(key string, person *types.CreatePersonParams) (*types.Person, bool) {
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

func (m *MemoryDB) All() *[]types.Person {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return &m.persons
}
