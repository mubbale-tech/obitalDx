package engine

import (
	"math"
	"orbitodx/internal/domain"
)

// CalculateCANSNeutron computes CSA Canadian Active Neutron Spectrometer (CANS) telemetry
// based on the active mission environment and daily total radiation exposure.
func CalculateCANSNeutron(mission domain.MissionProfile, vitals domain.VitalsInput, series []domain.TelemetryDayPoint) domain.CANSNeutronTelemetry {
	fraction := mission.CANSNeutronFraction
	// During a Solar Particle Event / high radiation spike (> 3.5 mSv/day), secondary fast-neutron
	// spallation inside aluminum hull shielding increases the neutron fraction further.
	if vitals.RadiationMSv >= 4.0 {
		fraction = math.Min(0.54, fraction+0.12)
	} else if vitals.RadiationMSv >= 2.2 {
		fraction = math.Min(0.46, fraction+0.05)
	}

	dailyNeutron := math.Round(vitals.RadiationMSv*fraction*100) / 100
	neutronPct := int(math.Round(fraction * 100))

	// Breakdown across neutron energy regimes tracked by CSA CANS (Bubble / Silicon / Bonner sphere equivalents):
	// Fast neutrons (100 keV - 15 MeV) carry highest biological RBE (w_R = 10 to 20).
	fastShare := 0.58
	epithermalShare := 0.27
	thermalShare := 0.15
	if mission.ID == domain.MissionArtemis {
		// Lunar regolith moderates some fast neutrons into epithermal albedo neutrons
		fastShare = 0.52
		epithermalShare = 0.31
		thermalShare = 0.17
	}

	cumTotal := 0.0
	cumNeutron := 0.0
	for _, pt := range series {
		cumTotal += pt.TotalRadMSv
		cumNeutron += pt.NeutronRadMSv
	}
	cumTotal = math.Round(cumTotal*10) / 10
	cumNeutron = math.Round(cumNeutron*10) / 10

	shieldingStatus := "NOMINAL — Polyethylene & Water-Wall Attenuation Optimal"
	spectrometerStatus := "CSA CANS-II ACTIVE —Nominal Bubble/Si Telemetry"
	if vitals.RadiationMSv >= 4.0 {
		shieldingStatus = "CRITICAL — Secondary Fast-Neutron Hull Spallation Detected; Storm Shelter Required"
		spectrometerStatus = "CSA CANS-II HIGH-FLUX ALERT — SPE Fast Neutron Surge (>1 MeV)"
	} else if vitals.RadiationMSv >= 2.0 {
		shieldingStatus = "ELEVATED — Increased GCR Secondary Neutron Scatter; Limit Crew Quarters Time in Forward Node"
		spectrometerStatus = "CSA CANS-II WATCH — Elevated Epithermal/Fast Albedo Flux"
	}

	return domain.CANSNeutronTelemetry{
		DailyTotalMSv:        math.Round(vitals.RadiationMSv*100) / 100,
		DailyNeutronMSv:      dailyNeutron,
		NeutronPercentage:    neutronPct,
		ThermalNeutronFlux:   math.Round(dailyNeutron*thermalShare*100) / 100,
		EpithermalFlux:       math.Round(dailyNeutron*epithermalShare*100) / 100,
		FastNeutronFlux:      math.Round(dailyNeutron*fastShare*100) / 100,
		Cumulative30DayMSv:   cumTotal,
		CumulativeNeutronMSv: cumNeutron,
		NASALimit30DayMSv:    250.0,
		ShieldingStatus:      shieldingStatus,
		SpectrometerStatus:   spectrometerStatus,
	}
}
