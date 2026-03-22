package repository

import (
	"fmt"
	"light-management/models"
	"sync"
)

type LightRepository interface {
	GetAll() []models.Light
	Create(name string) *models.Light
	GetByID(id string) (*models.Light, bool)
	UpdateStatus(id, status string) bool
}

type memoryLightRepo struct {
	sync.RWMutex
	lights map[string]*models.Light
}

func NewMemoryLightRepository() LightRepository {
	return &memoryLightRepo{
		lights: make(map[string]*models.Light),
	}
}

func (r *memoryLightRepo) GetAll() []models.Light {
	r.RLock()
	defer r.RUnlock()
	var lightsList []models.Light
	for _, l := range r.lights {
		lightsList = append(lightsList, *l)
	}
	return lightsList
}

func (r *memoryLightRepo) Create(name string) *models.Light {
	r.Lock()
	defer r.Unlock()

	id := fmt.Sprintf("%d", len(r.lights)+1)
	light := &models.Light{
		ID:     id,
		Name:   name,
		Status: "OFF",
	}
	r.lights[id] = light
	return light
}

func (r *memoryLightRepo) GetByID(id string) (*models.Light, bool) {
	r.RLock()
	defer r.RUnlock()

	l, exists := r.lights[id]
	if !exists {
		return nil, false
	}
	
	lightCopy := *l
	return &lightCopy, true
}

func (r *memoryLightRepo) UpdateStatus(id, status string) bool {
	r.Lock()
	defer r.Unlock()

	if l, exists := r.lights[id]; exists {
		l.Status = status
		return true
	}
	return false
}
