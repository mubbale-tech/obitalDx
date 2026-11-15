package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"orbitodx/internal/domain"
	"orbitodx/internal/store"
	"orbitodx/web"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func setupTestApp() *fiber.App {
	appState := store.NewAppState()
	h := NewServerHandlers(appState)
	app := fiber.New()
	app.Get("/static*", static.New("", static.Config{
		FS: web.StaticFS(),
	}))
	app.Get("/", h.HandleDashboard)
	app.Get("/api/health", h.HandleHealth)
	app.Get("/api/state", h.HandleGetState)
	app.Post("/api/evaluate", h.HandleEvaluateVitals)
	app.Post("/api/scenario/:id", h.HandleScenario)
	app.Post("/api/mission/:id", h.HandleSwitchMission)
	app.Post("/api/protocol/:stepId/toggle", h.HandleToggleStep)
	return app
}

func TestDashboardAndEndpoints(t *testing.T) {
	app := setupTestApp()

	// 1. Test GET / renders HTML with Modules A-E
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on GET /, got %d", resp.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	html := string(bodyBytes)
	for _, marker := range []string{"module-a", "module-b", "module-c", "module-d", "module-e", "CSA CANS", "mars_neutron_spike", "Mars Solar Neutron Spike", "Microgravity Bone"} {
		if !strings.Contains(html, marker) {
			t.Errorf("expected dashboard HTML to contain %q", marker)
		}
	}

	// 1b. Test static assets return 200 OK
	for _, asset := range []string{"/static/css/app.css", "/static/js/cockpit.js"} {
		sReq := httptest.NewRequest("GET", asset, nil)
		sResp, sErr := app.Test(sReq)
		if sErr != nil || sResp.StatusCode != 200 {
			t.Fatalf("expected 200 on %s, got status=%d err=%v", asset, sResp.StatusCode, sErr)
		}
	}

	// 2. Test POST /api/scenario/mars_neutron_spike
	req2 := httptest.NewRequest("POST", "/api/scenario/mars_neutron_spike", nil)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("POST /api/scenario/mars_neutron_spike failed: %v", err)
	}
	var eval domain.EvaluationResult
	if err := json.NewDecoder(resp2.Body).Decode(&eval); err != nil {
		t.Fatalf("failed to decode scenario response: %v", err)
	}
	if eval.OverallTier != domain.TierRed {
		t.Errorf("expected mars_neutron_spike to return RED, got %s", eval.OverallTier)
	}
	if eval.Mission.ID != domain.MissionMars {
		t.Errorf("expected mission to be mars, got %s", eval.Mission.ID)
	}

	// 3. Test POST /api/protocol/rad-step-1/toggle
	initialScore := eval.ReadinessScore
	req3 := httptest.NewRequest("POST", "/api/protocol/rad-step-1/toggle", nil)
	resp3, err := app.Test(req3)
	if err != nil {
		t.Fatalf("POST /api/protocol/rad-step-1/toggle failed: %v", err)
	}
	var evalAfterStep domain.EvaluationResult
	_ = json.NewDecoder(resp3.Body).Decode(&evalAfterStep)
	if evalAfterStep.ReadinessScore <= initialScore {
		t.Errorf("expected readiness score to increase after completing protocol step: before=%d, after=%d", initialScore, evalAfterStep.ReadinessScore)
	}

	// 4. Test POST /api/evaluate with custom vitals
	custom := domain.VitalsInput{
		MissionID:      domain.MissionArtemis,
		SystolicBP:     116,
		DiastolicBP:    74,
		HeartRate:      66,
		BloodSugar:     90,
		BodyTemp:       36.7,
		SleepHours:     8.0,
		MoodScore:      9,
		RadiationMSv:   1.15,
		AREDCompliance: 92,
		CabinCO2:       2.9,
	}
	payload, _ := json.Marshal(custom)
	req4 := httptest.NewRequest("POST", "/api/evaluate", bytes.NewReader(payload))
	req4.Header.Set("Content-Type", "application/json")
	resp4, err := app.Test(req4)
	if err != nil || resp4.StatusCode != 200 {
		t.Fatalf("POST /api/evaluate failed: %v (status %d)", err, resp4.StatusCode)
	}
}
