package handlers

import (
	"orbitodx/internal/domain"
	"orbitodx/internal/engine"
	"orbitodx/internal/store"
	"orbitodx/views"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

type ServerHandlers struct {
	State *store.AppState
}

func NewServerHandlers(state *store.AppState) *ServerHandlers {
	return &ServerHandlers{State: state}
}

func renderTempl(c fiber.Ctx, comp templ.Component) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return comp.Render(c.Context(), c.Response().BodyWriter())
}

// HandleDashboard serves the main OrbitoDX spacecraft cockpit UI
func (h *ServerHandlers) HandleDashboard(c fiber.Ctx) error {
	boot := domain.InitialBootPayload{
		Evaluation: h.State.GetEvaluation(),
		Missions:   domain.AllMissionProfiles(),
		Hazards:    domain.AllNASAHazards(),
		Scenarios:  engine.AllScenarioPresets(),
	}
	return renderTempl(c, views.Dashboard(boot))
}

// HandleGetState returns the current clinical triage state as JSON
func (h *ServerHandlers) HandleGetState(c fiber.Ctx) error {
	return c.JSON(h.State.GetEvaluation())
}

// HandleEvaluateVitals evaluates manually submitted astronaut vitals
func (h *ServerHandlers) HandleEvaluateVitals(c fiber.Ctx) error {
	var input domain.VitalsInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid vitals payload: " + err.Error(),
		})
	}
	res := h.State.UpdateVitals(input)
	return c.JSON(res)
}

// HandleScenario applies one of the 1-Click Demo Scenarios
func (h *ServerHandlers) HandleScenario(c fiber.Ctx) error {
	id := c.Params("id")
	res, ok := h.State.ApplyScenario(id)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "scenario not found: " + id,
		})
	}
	return c.JSON(res)
}

// HandleSwitchMission switches between ISS, Artemis Lunar Base, and Mars Transit Spacecraft
func (h *ServerHandlers) HandleSwitchMission(c fiber.Ctx) error {
	id := domain.MissionID(c.Params("id"))
	res := h.State.SwitchMission(id)
	return c.JSON(res)
}

// HandleToggleStep toggles completion of a countermeasure protocol checklist step
func (h *ServerHandlers) HandleToggleStep(c fiber.Ctx) error {
	stepID := c.Params("stepId")
	res := h.State.ToggleProtocolStep(stepID)
	return c.JSON(res)
}

// HandleHealth returns a lightweight health probe for VPS systemd / reverse proxies
func (h *ServerHandlers) HandleHealth(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "orbitodx",
		"engine":  "NASA-HRP-5-Hazard + CSA-CANS-v3.4",
	})
}
