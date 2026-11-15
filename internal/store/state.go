package store

import (
	"orbitodx/internal/domain"
	"orbitodx/internal/engine"
	"sync"
)

type AppState struct {
	mu             sync.RWMutex
	currentVitals  domain.VitalsInput
	completedSteps map[string]bool
	streakDays     int
}

func NewAppState() *AppState {
	// Start on the Nominal Baseline or Mars Transit preset
	presets := engine.AllScenarioPresets()
	return &AppState{
		currentVitals:  presets[0].Vitals,
		completedSteps: make(map[string]bool),
		streakDays:     29,
	}
}

func (s *AppState) GetEvaluation() domain.EvaluationResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stepsCopy := make(map[string]bool, len(s.completedSteps))
	for k, v := range s.completedSteps {
		stepsCopy[k] = v
	}
	return engine.EvaluateCrewHealth(s.currentVitals, stepsCopy, s.streakDays)
}

func (s *AppState) UpdateVitals(v domain.VitalsInput) domain.EvaluationResult {
	s.mu.Lock()
	if v.MissionID == "" {
		v.MissionID = s.currentVitals.MissionID
	}
	if v.ScenarioID == "" {
		v.ScenarioID = "custom"
	}
	s.currentVitals = v
	s.mu.Unlock()
	return s.GetEvaluation()
}

func (s *AppState) ApplyScenario(scenarioID string) (domain.EvaluationResult, bool) {
	preset, ok := engine.GetScenarioByID(scenarioID)
	if !ok {
		return s.GetEvaluation(), false
	}
	s.mu.Lock()
	s.currentVitals = preset.Vitals
	s.completedSteps = make(map[string]bool)
	s.mu.Unlock()
	return s.GetEvaluation(), true
}

func (s *AppState) SwitchMission(missionID domain.MissionID) domain.EvaluationResult {
	profile := domain.GetMissionProfile(missionID)
	s.mu.Lock()
	s.currentVitals.MissionID = profile.ID
	// If switching mission manually, align background radiation and cabin CO2 if not in a solar storm spike
	if s.currentVitals.RadiationMSv < 3.5 {
		s.currentVitals.RadiationMSv = profile.BaselineRadMSvDay
	}
	if s.currentVitals.CabinCO2 < 4.2 {
		s.currentVitals.CabinCO2 = profile.DefaultCabinCO2
	}
	s.mu.Unlock()
	return s.GetEvaluation()
}

func (s *AppState) ToggleProtocolStep(stepID string) domain.EvaluationResult {
	s.mu.Lock()
	s.completedSteps[stepID] = !s.completedSteps[stepID]
	// Reward streak increment if at least 2 steps completed
	count := 0
	for _, done := range s.completedSteps {
		if done {
			count++
		}
	}
	if count >= 2 {
		s.streakDays = 30
	} else {
		s.streakDays = 29
	}
	s.mu.Unlock()
	return s.GetEvaluation()
}
