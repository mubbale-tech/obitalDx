package engine

import (
	"orbitodx/internal/domain"
	"testing"
)

func TestScenarioEvaluations(t *testing.T) {
	presets := AllScenarioPresets()
	if len(presets) != 4 {
		t.Fatalf("expected 4 scenario presets, got %d", len(presets))
	}

	// 1. Nominal Baseline -> Overall GREEN
	nom := EvaluateCrewHealth(presets[0].Vitals, map[string]bool{}, 29)
	if nom.OverallTier != domain.TierGreen {
		t.Errorf("expected Nominal Baseline to be GREEN, got %s", nom.OverallTier)
	}
	if nom.CANSNeutron.NeutronPercentage != 30 {
		t.Errorf("expected ISS CANS neutron share to be 30%%, got %d%%", nom.CANSNeutron.NeutronPercentage)
	}

	// 2. Mars Solar Neutron Spike -> Overall RED + CSA CANS high flux
	mars := EvaluateCrewHealth(presets[1].Vitals, map[string]bool{}, 29)
	if mars.OverallTier != domain.TierRed {
		t.Errorf("expected Mars Solar Neutron Spike to be RED, got %s", mars.OverallTier)
	}
	if mars.RedCount < 3 {
		t.Errorf("expected at least 3 RED cards in Mars Solar Neutron Spike, got %d", mars.RedCount)
	}
	if mars.CANSNeutron.DailyNeutronMSv <= 3.0 {
		t.Errorf("expected elevated CANS daily neutron dose > 3.0 mSv, got %.2f", mars.CANSNeutron.DailyNeutronMSv)
	}

	// 3. Microgravity Bone & Cardio Drift -> RED
	bone := EvaluateCrewHealth(presets[2].Vitals, map[string]bool{}, 29)
	if bone.OverallTier != domain.TierRed {
		t.Errorf("expected Microgravity Bone & Cardio Drift to be RED, got %s", bone.OverallTier)
	}

	// 4. Circadian Sleep Deficit -> AMBER or RED
	sleep := EvaluateCrewHealth(presets[3].Vitals, map[string]bool{}, 29)
	if sleep.OverallTier == domain.TierGreen {
		t.Errorf("expected Circadian Sleep Deficit to trigger non-GREEN alert")
	}
}
