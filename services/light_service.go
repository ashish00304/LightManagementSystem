package services

import (
	"fmt"
	"light-management/models"
	"light-management/repository"
	"sync"
	"time"
)

type LightService interface {
	GetAll() []models.Light
	Create(name string) *models.Light
	TurnOn(id string, duration int) (*models.Light, error)
	TurnOff(id string) (*models.Light, error)
}

type lightService struct {
	repo   repository.LightRepository
	timers map[string]*time.Timer
	mu     sync.Mutex
}

func NewLightService(repo repository.LightRepository) LightService {
	return &lightService{
		repo:   repo,
		timers: make(map[string]*time.Timer),
	}
}

func (s *lightService) GetAll() []models.Light {
	return s.repo.GetAll()
}

func (s *lightService) Create(name string) *models.Light {
	return s.repo.Create(name)
}

func (s *lightService) TurnOn(id string, duration int) (*models.Light, error) {
	_, exists := s.repo.GetByID(id)
	if !exists {
		return nil, fmt.Errorf("light not found")
	}

	s.repo.UpdateStatus(id, "ON")

	s.mu.Lock()
	defer s.mu.Unlock()

	if timer, ok := s.timers[id]; ok {
		timer.Stop()
		delete(s.timers, id)
	}

	if duration > 0 {
		timer := time.AfterFunc(time.Duration(duration)*time.Second, func() {
			s.repo.UpdateStatus(id, "OFF")
			s.mu.Lock()
			delete(s.timers, id)
			s.mu.Unlock()
			fmt.Printf("Light %s automatically turned OFF.\n", id)
		})
		s.timers[id] = timer
	}

	light, _ := s.repo.GetByID(id)
	return light, nil
}

func (s *lightService) TurnOff(id string) (*models.Light, error) {
	_, exists := s.repo.GetByID(id)
	if !exists {
		return nil, fmt.Errorf("light not found")
	}

	s.repo.UpdateStatus(id, "OFF")

	s.mu.Lock()
	defer s.mu.Unlock()

	if timer, ok := s.timers[id]; ok {
		timer.Stop()
		delete(s.timers, id)
	}

	light, _ := s.repo.GetByID(id)
	return light, nil
}
