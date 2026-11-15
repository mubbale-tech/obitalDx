package engine

import (
	"fmt"
	"math"
	"orbitodx/internal/domain"
)

func AllScenarioPresets() []domain.ScenarioPreset {
	return []domain.ScenarioPreset{
		{
			ID:          "nominal",
			Name:        "Nominal Baseline",
			Badge:       "ALL GREEN • LEO / ISS",
			Tier:        domain.TierGreen,
			Description: "Optimal spaceflight homeostasis: balanced blood pressure, normoglycemia, 7.8h sleep, 95% ARED exercise compliance, and nominal CSA CANS neutron background.",
			Vitals: domain.VitalsInput{
				MissionID:      domain.MissionISS,
				ScenarioID:     "nominal",
				SystolicBP:     118,
				DiastolicBP:    76,
				HeartRate:      64,
				BloodSugar:     92,
				BodyTemp:       36.7,
				SleepHours:     7.8,
				MoodScore:      9,
				RadiationMSv:   0.42,
				AREDCompliance: 95,
				CabinCO2:       2.7,
			},
		},
		{
			ID:          "mars_neutron_spike",
			Name:        "Mars Solar Neutron Spike",
			Badge:       "RED ESCALATION • 14m DELAY",
			Tier:        domain.TierRed,
			Description: "Deep-space Solar Particle Event triggers a secondary fast-neutron surge (CSA CANS 8.4 mSv/d), stress hyperglycemia (152 mg/dL), sleep disruption (4.5h), and early viral reactivation fever (37.7°C) during a 28m 44s round-trip comm blackout.",
			Vitals: domain.VitalsInput{
				MissionID:      domain.MissionMars,
				ScenarioID:     "mars_neutron_spike",
				SystolicBP:     138,
				DiastolicBP:    88,
				HeartRate:      94,
				BloodSugar:     152,
				BodyTemp:       37.7,
				SleepHours:     4.5,
				MoodScore:      4,
				RadiationMSv:   8.40,
				AREDCompliance: 72,
				CabinCO2:       3.8,
			},
		},
		{
			ID:          "bone_cardio_drift",
			Name:        "Microgravity Bone & Cardio Drift",
			Badge:       "AMBER / RED • 0g UNLOADING",
			Tier:        domain.TierRed,
			Description: "30-day progressive cardiovascular deconditioning (BP 92/58 mmHg orthostatic drop, HR 90 bpm) compounded by skipped ARED resistive exercise (32% compliance), accelerating femoral/lumbar bone loss to 1.68%/month.",
			Vitals: domain.VitalsInput{
				MissionID:      domain.MissionMars,
				ScenarioID:     "bone_cardio_drift",
				SystolicBP:     92,
				DiastolicBP:    58,
				HeartRate:      90,
				BloodSugar:     104,
				BodyTemp:       36.8,
				SleepHours:     6.4,
				MoodScore:      6,
				RadiationMSv:   1.84,
				AREDCompliance: 32,
				CabinCO2:       3.6,
			},
		},
		{
			ID:          "circadian_sleep_deficit",
			Name:        "Circadian Sleep Deficit",
			Badge:       "AMBER WATCH • SANS + CO₂",
			Tier:        domain.TierAmber,
			Description: "Chronic circadian misalignment (4.1h sleep, Mood 3/10) combined with elevated closed-habitat CO₂ (4.7 mmHg), amplifying intracranial venous congestion and Spaceflight Associated Neuro-ocular Syndrome (SANS).",
			Vitals: domain.VitalsInput{
				MissionID:      domain.MissionArtemis,
				ScenarioID:     "circadian_sleep_deficit",
				SystolicBP:     128,
				DiastolicBP:    84,
				HeartRate:      82,
				BloodSugar:     124,
				BodyTemp:       37.1,
				SleepHours:     4.1,
				MoodScore:      3,
				RadiationMSv:   1.35,
				AREDCompliance: 78,
				CabinCO2:       4.7,
			},
		},
	}
}

func GetScenarioByID(id string) (domain.ScenarioPreset, bool) {
	for _, s := range AllScenarioPresets() {
		if s.ID == id {
			return s, true
		}
	}
	return AllScenarioPresets()[0], false
}

// Generate30DayTelemetry builds a deterministic, physiologically realistic 30-day mission
// time series (Sol 1 through Sol 30) that smoothly converges on the live Sol 30 VitalsInput.
func Generate30DayTelemetry(mission domain.MissionProfile, current domain.VitalsInput) []domain.TelemetryDayPoint {
	points := make([]domain.TelemetryDayPoint, 30)

	// Anchor Sol 1 at healthy launch/early-orbit baseline for the given mission
	baseSys := 118.0
	baseDia := 76.0
	baseHR := 64.0
	baseGlu := 90.0
	baseTemp := 36.7
	baseSleep := 7.8
	baseMood := 9.0
	baseRad := mission.BaselineRadMSvDay
	baseARED := 94.0
	baseCO2 := mission.DefaultCabinCO2

	cumRad := 0.0

	for i := 0; i < 30; i++ {
		sol := i + 1
		// Progress factor from 0.0 (Sol 1) to 1.0 (Sol 30), accelerating in the final 10 Sols
		t := float64(i) / 29.0
		curve := math.Pow(t, 1.45)

		// Deterministic subtle daily physiological oscillation
		wave := math.Sin(float64(sol)*0.85) * (1.0 - t*0.5)

		sys := int(math.Round(baseSys + (float64(current.SystolicBP)-baseSys)*curve + wave*2.0))
		dia := int(math.Round(baseDia + (float64(current.DiastolicBP)-baseDia)*curve + wave*1.5))
		hr := int(math.Round(baseHR + (float64(current.HeartRate)-baseHR)*curve - wave*2.0))
		glu := int(math.Round(baseGlu + (float64(current.BloodSugar)-baseGlu)*curve + wave*2.5))
		temp := math.Round((baseTemp+(current.BodyTemp-baseTemp)*curve+wave*0.06)*10) / 10
		sleep := math.Round((baseSleep+(current.SleepHours-baseSleep)*curve+wave*0.25)*10) / 10
		mood := int(math.Round(baseMood + (float64(current.MoodScore)-baseMood)*curve))
		if mood < 1 {
			mood = 1
		} else if mood > 10 {
			mood = 10
		}

		ared := int(math.Round(baseARED + (float64(current.AREDCompliance)-baseARED)*curve))
		co2 := math.Round((baseCO2+(current.CabinCO2-baseCO2)*curve+wave*0.1)*10) / 10

		// Radiation: if current radiation is a sudden Solar Particle Event spike (>3.5 mSv/d),
		// model background GCR for Sols 1-26 and a steep solar storm ramp on Sols 27-30.
		var dailyRad float64
		if current.RadiationMSv >= 3.5 {
			if sol < 26 {
				dailyRad = baseRad + math.Abs(wave)*0.15
			} else {
				stormProgress := float64(sol-25) / 5.0
				dailyRad = baseRad + (current.RadiationMSv-baseRad)*math.Pow(stormProgress, 1.6)
			}
		} else {
			dailyRad = baseRad + (current.RadiationMSv-baseRad)*curve + math.Abs(wave)*0.08
		}
		if sol == 30 {
			sys = current.SystolicBP
			dia = current.DiastolicBP
			hr = current.HeartRate
			glu = current.BloodSugar
			temp = current.BodyTemp
			sleep = current.SleepHours
			mood = current.MoodScore
			dailyRad = current.RadiationMSv
			ared = current.AREDCompliance
			co2 = current.CabinCO2
		}

		dailyRad = math.Round(dailyRad*100) / 100
		neutronShare := mission.CANSNeutronFraction
		if dailyRad >= 4.0 {
			neutronShare = math.Min(0.54, neutronShare+0.12)
		}
		neutronRad := math.Round(dailyRad*neutronShare*100) / 100
		cumRad = math.Round((cumRad+dailyRad)*100) / 100

		// Compute SANS Risk Index (0-100) driven by 0g cephalad shift, cabin CO2 vasodilation, sleep debt, and systolic BP
		sansBase := 24.0
		if mission.GravityG == 0.0 {
			sansBase = 38.0
		} else if mission.GravityG < 0.25 {
			sansBase = 28.0
		}
		if co2 > 3.2 {
			sansBase += (co2 - 3.2) * 18.0
		}
		if sleep < 6.5 {
			sansBase += (6.5 - sleep) * 6.5
		}
		if sys > 125 {
			sansBase += float64(sys-125) * 0.7
		}
		sansIdx := int(math.Min(98, math.Max(12, math.Round(sansBase))))

		// Compute Monthly Bone Loss Rate (% per month) driven by gravity unloading and ARED compliance
		unloadingFactor := 1.65 // 0g unmitigated loss ~1.5-1.7%/mo
		if mission.GravityG > 0.1 {
			unloadingFactor = 1.05 // 1/6g partial unloading
		}
		mitigation := (float64(ared) / 100.0) * 1.25
		boneLoss := math.Round(math.Max(0.22, unloadingFactor-mitigation)*100) / 100

		// Compute Daily Readiness Score (0-100)
		readiness := 96
		if sleep < 6.0 {
			readiness -= int((6.5 - sleep) * 7)
		}
		if dailyRad > 2.0 {
			readiness -= int( math.Min(28, (dailyRad-1.5)*4.2) )
		}
		if ared < 75 {
			readiness -= (75 - ared) / 3
		}
		if sys < 98 || sys > 135 {
			readiness -= 12
		}
		if glu > 125 {
			readiness -= 10
		}
		if readiness < 34 {
			readiness = 34
		}

		points[i] = domain.TelemetryDayPoint{
			Sol:              sol,
			Label:            fmt.Sprintf("Sol %d", sol),
			SystolicBP:       sys,
			DiastolicBP:      dia,
			HeartRate:        hr,
			BloodSugar:       glu,
			BodyTemp:         temp,
			SleepHours:       sleep,
			MoodScore:        mood,
			TotalRadMSv:      dailyRad,
			NeutronRadMSv:    neutronRad,
			CumulativeRadMSv: cumRad,
			SANSRiskIndex:    sansIdx,
			BoneLossRatePct:  boneLoss,
			AREDCompliance:   ared,
			CabinCO2:         co2,
			ReadinessScore:   readiness,
		}
	}

	return points
}
