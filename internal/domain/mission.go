package domain

type MissionID string

const (
	MissionISS     MissionID = "iss"
	MissionArtemis MissionID = "artemis"
	MissionMars    MissionID = "mars"
)

type MissionProfile struct {
	ID                      MissionID `json:"id"`
	Name                    string    `json:"name"`
	ShortCode               string    `json:"short_code"`
	Subtitle                string    `json:"subtitle"`
	Habitat                 string    `json:"habitat"`
	DistanceKm              string    `json:"distance_km"`
	GravityG                float64   `json:"gravity_g"`
	GravityLabel            string    `json:"gravity_label"`
	CommDelayOneWaySec      float64   `json:"comm_delay_one_way_sec"`
	CommDelayOneWayLabel    string    `json:"comm_delay_one_way_label"`
	CommDelayRoundTripLabel string    `json:"comm_delay_round_trip_label"`
	BaselineRadMSvDay       float64   `json:"baseline_rad_msv_day"`
	CANSNeutronFraction     float64   `json:"cans_neutron_fraction"`
	CANSNeutronNote         string    `json:"cans_neutron_note"`
	DefaultCabinCO2         float64   `json:"default_cabin_co2"`
	AutonomousMode          string    `json:"autonomous_mode"`
	PrimaryEnvRisk          string    `json:"primary_env_risk"`
	ActiveHazards           []string  `json:"active_hazards"`
}

func AllMissionProfiles() []MissionProfile {
	return []MissionProfile{
		{
			ID:                      MissionISS,
			Name:                    "ISS — Low Earth Orbit",
			ShortCode:               "LEO / ISS-EXP",
			Subtitle:                "Magnetosphere-Partially Shielded Microgravity Laboratory",
			Habitat:                 "U.S. Destiny Lab & Columbus Orbital Segment (408 km)",
			DistanceKm:              "408 km",
			GravityG:                0.00,
			GravityLabel:            "0.00g Microgravity",
			CommDelayOneWaySec:      0.8,
			CommDelayOneWayLabel:    "0.8s One-Way",
			CommDelayRoundTripLabel: "1.6s Real-Time Ground Link",
			BaselineRadMSvDay:       0.42,
			CANSNeutronFraction:     0.30,
			CANSNeutronNote:         "CSA CANS telemetry confirms secondary neutrons contribute ~30% of total ISS crew dose Equivalent via hull spallation.",
			DefaultCabinCO2:         2.8,
			AutonomousMode:          "LEVEL 1 — Live Flight Surgeon Teleguidance",
			PrimaryEnvRisk:          "Cephalad fluid shift, bone/muscle unloading, South Atlantic Anomaly proton passes, closed-loop CO₂ pockets.",
			ActiveHazards:           []string{"gravity", "radiation", "environment", "isolation"},
		},
		{
			ID:                      MissionArtemis,
			Name:                    "Artemis Lunar Base",
			ShortCode:               "LUNAR / SHACKLETON",
			Subtitle:                "1/6g Surface Habitat + EVA Regolith Dust Exposure",
			Habitat:                 "Artemis Surface Habitat & Pressurized Rover (South Pole)",
			DistanceKm:              "384,400 km",
			GravityG:                0.16,
			GravityLabel:            "0.16g (1/6g Lunar)",
			CommDelayOneWaySec:      1.3,
			CommDelayOneWayLabel:    "1.3s One-Way",
			CommDelayRoundTripLabel: "2.6s Near-Real-Time Relay",
			BaselineRadMSvDay:       1.15,
			CANSNeutronFraction:     0.36,
			CANSNeutronNote:         "Deep-space GCRs interacting with lunar regolith produce upward-scattered albedo neutrons (~36% of dose equivalent).",
			DefaultCabinCO2:         3.1,
			AutonomousMode:          "LEVEL 2 — Crew Medical Officer First-Response",
			PrimaryEnvRisk:          "Ultrafine electrostatic lunar silica dust (<2 µm) after EVAs, partial-gravity vestibular remodeling, unshielded solar particle events.",
			ActiveHazards:           []string{"radiation", "environment", "gravity", "isolation", "distance"},
		},
		{
			ID:                      MissionMars,
			Name:                    "Mars Transit Spacecraft",
			ShortCode:               "DEEP SPACE / MTV-01",
			Subtitle:                "0g Interplanetary Cruise + High GCR/Neutron Flux + 14-Min One-Way Delay",
			Habitat:                 "Deep-Space Transit Habitat (185M km from Earth)",
			DistanceKm:              "185,000,000 km",
			GravityG:                0.00,
			GravityLabel:            "0.00g Deep-Space Cruise",
			CommDelayOneWaySec:      862.0,
			CommDelayOneWayLabel:    "14m 22s One-Way",
			CommDelayRoundTripLabel: "28m 44s Round-Trip Delay",
			BaselineRadMSvDay:       1.84,
			CANSNeutronFraction:     0.42,
			CANSNeutronNote:         "High-energy HZE ions & Solar Particle Events generate intense secondary fast-neutron cascades inside aluminum hull (~42% effective dose).",
			DefaultCabinCO2:         3.5,
			AutonomousMode:          "LEVEL 3 — Full Autonomous Onboard Triage (No Live Houston Link)",
			PrimaryEnvRisk:          "28m 44s round-trip communication blackout for acute emergencies; compounding GCR/Neutron radiation, circadian drift, SANS, and bone loss.",
			ActiveHazards:           []string{"distance", "radiation", "gravity", "isolation", "environment"},
		},
	}
}

func GetMissionProfile(id MissionID) MissionProfile {
	for _, m := range AllMissionProfiles() {
		if m.ID == id {
			return m
		}
	}
	return AllMissionProfiles()[2] // Default to Mars Transit for high-impact demo
}
