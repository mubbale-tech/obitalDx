package domain

type TriageTier string

const (
	TierGreen TriageTier = "GREEN"
	TierAmber TriageTier = "AMBER"
	TierRed   TriageTier = "RED"
)

type ProtocolStep struct {
	ID          string `json:"id"`
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Instruction string `json:"instruction"`
	Equipment   string `json:"equipment"`
	Timeframe   string `json:"timeframe"`
	Completed   bool   `json:"completed"`
}

type ClinicalAlertCard struct {
	ID                  string         `json:"id"`
	Title               string         `json:"title"`
	BiomarkerName       string         `json:"biomarker_name"`
	BodySystem          BodySystemID   `json:"body_system"`
	BodySystemLabel     string         `json:"body_system_label"`
	CurrentValueDisplay string         `json:"current_value_display"`
	NominalRangeDisplay string         `json:"nominal_range_display"`
	Tier                TriageTier     `json:"tier"`
	TierBadgeText       string         `json:"tier_badge_text"`
	HazardTags          []NASAHazard   `json:"hazard_tags"`
	CSACANSLinked       bool           `json:"csa_cans_linked"`
	WhyExamNeeded       string         `json:"why_exam_needed"`
	HowHazardsCause     string         `json:"how_hazards_cause"`
	EarthClinicAnalogy  string         `json:"earth_clinic_analogy"`
	ProtocolTitle       string         `json:"protocol_title"`
	ProtocolSteps       []ProtocolStep `json:"protocol_steps"`
}

type BodySystemSummary struct {
	ID             BodySystemID `json:"id"`
	Name           string       `json:"name"`
	Subtitle       string       `json:"subtitle"`
	Tier           TriageTier   `json:"tier"`
	ActiveAlerts   int          `json:"active_alerts"`
	HeadlineMetric string       `json:"headline_metric"`
	SecondaryNote  string       `json:"secondary_note"`
	PrimaryHazards []string     `json:"primary_hazards"`
}

type ScenarioPreset struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Badge       string      `json:"badge"`
	Tier        TriageTier  `json:"tier"`
	Description string      `json:"description"`
	Vitals      VitalsInput `json:"vitals"`
}

type EvaluationResult struct {
	Timestamp               string               `json:"timestamp"`
	Mission                 MissionProfile       `json:"mission"`
	Vitals                  VitalsInput          `json:"vitals"`
	OverallTier             TriageTier           `json:"overall_tier"`
	OverallHeadline         string               `json:"overall_headline"`
	OverallSubtext          string               `json:"overall_subtext"`
	GreenCount              int                  `json:"green_count"`
	AmberCount              int                  `json:"amber_count"`
	RedCount                int                  `json:"red_count"`
	ReadinessScore          int                  `json:"readiness_score"`
	StreakDays              int                  `json:"streak_days"`
	StreakTargets           []StreakTarget       `json:"streak_targets"`
	BodySystems             []BodySystemSummary  `json:"body_systems"`
	AlertCards              []ClinicalAlertCard  `json:"alert_cards"`
	CompoundSynergies       []CompoundSynergy    `json:"compound_synergies"`
	CANSNeutron             CANSNeutronTelemetry `json:"cans_neutron"`
	Telemetry30Days         []TelemetryDayPoint  `json:"telemetry_30_days"`
	StoreAndForwardPacketID string               `json:"store_and_forward_packet_id"`
}
