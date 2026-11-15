package domain

type BodySystemID string

const (
	SystemAll      BodySystemID = "all"
	SystemHead     BodySystemID = "head"
	SystemHeart    BodySystemID = "heart"
	SystemSkeleton BodySystemID = "skeleton"
	SystemBlood    BodySystemID = "blood"
)

type VitalsInput struct {
	MissionID      MissionID `json:"mission_id"`
	ScenarioID     string    `json:"scenario_id"`
	SystolicBP     int       `json:"systolic_bp"`     // mmHg (Normal: 105-125)
	DiastolicBP    int       `json:"diastolic_bp"`    // mmHg (Normal: 65-82)
	HeartRate      int       `json:"heart_rate"`      // bpm  (Normal: 55-78)
	BloodSugar     int       `json:"blood_sugar"`     // mg/dL (Normal: 75-105 fasting)
	BodyTemp       float64   `json:"body_temp"`       // °C (Normal: 36.3-37.2)
	SleepHours     float64   `json:"sleep_hours"`     // hrs/24h (Normal: 7.0-8.5)
	MoodScore      int       `json:"mood_score"`      // 1-10 (Normal: 7-10)
	RadiationMSv   float64   `json:"radiation_msv"`   // mSv/day (ISS ~0.42, Lunar ~1.15, Mars ~1.84, SPE Spike >4.0)
	AREDCompliance int       `json:"ared_compliance"` // % of daily resistive load target (Normal: 85-100%)
	CabinCO2       float64   `json:"cabin_co2"`       // mmHg partial pressure CO2 (Nominal: <3.0, Amber: 3.5-4.5, Red: >4.5)
}

type CANSNeutronTelemetry struct {
	DailyTotalMSv        float64 `json:"daily_total_msv"`
	DailyNeutronMSv      float64 `json:"daily_neutron_msv"`
	NeutronPercentage    int     `json:"neutron_percentage"`
	ThermalNeutronFlux   float64 `json:"thermal_neutron_flux"`   // < 0.5 eV
	EpithermalFlux       float64 `json:"epithermal_flux"`        // 0.5 eV - 100 keV
	FastNeutronFlux      float64 `json:"fast_neutron_flux"`      // 100 keV - 15 MeV (highest biological RBE weight w_R up to 20)
	Cumulative30DayMSv   float64 `json:"cumulative_30d_msv"`
	CumulativeNeutronMSv float64 `json:"cumulative_neutron_msv"`
	NASALimit30DayMSv    float64 `json:"nasa_limit_30d_msv"` // 250 mSv 30-day BFO limit
	ShieldingStatus      string  `json:"shielding_status"`
	SpectrometerStatus   string  `json:"spectrometer_status"`
}

type TelemetryDayPoint struct {
	Sol              int     `json:"sol"`
	Label            string  `json:"label"`
	SystolicBP       int     `json:"systolic_bp"`
	DiastolicBP      int     `json:"diastolic_bp"`
	HeartRate        int     `json:"heart_rate"`
	BloodSugar       int     `json:"blood_sugar"`
	BodyTemp         float64 `json:"body_temp"`
	SleepHours       float64 `json:"sleep_hours"`
	MoodScore        int     `json:"mood_score"`
	TotalRadMSv      float64 `json:"total_rad_msv"`
	NeutronRadMSv    float64 `json:"neutron_rad_msv"`
	CumulativeRadMSv float64 `json:"cumulative_rad_msv"`
	SANSRiskIndex    int     `json:"sans_risk_index"`    // 0-100 composite neuro-ocular pressure proxy
	BoneLossRatePct  float64 `json:"bone_loss_rate_pct"` // % per month loss rate (e.g., 0.35% with ARED vs 1.65% without)
	AREDCompliance   int     `json:"ared_compliance"`
	CabinCO2         float64 `json:"cabin_co2"`
	ReadinessScore   int     `json:"readiness_score"`
}

type StreakTarget struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Detail    string `json:"detail"`
	Completed bool   `json:"completed"`
	Points    int    `json:"points"`
}
