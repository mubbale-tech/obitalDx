package domain

type HazardID string

const (
	HazardRadiation   HazardID = "radiation"
	HazardIsolation   HazardID = "isolation"
	HazardDistance    HazardID = "distance"
	HazardGravity     HazardID = "gravity"
	HazardEnvironment HazardID = "environment"
)

type NASAHazard struct {
	ID          HazardID `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	ShortName   string   `json:"short_name"`
	Icon        string   `json:"icon"`
	ColorClass  string   `json:"color_class"`
	Description string   `json:"description"`
	HRPRisks    []string `json:"hrp_risks"`
}

type CompoundSynergy struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Severity        TriageTier `json:"severity"`
	InvolvedHazards []HazardID `json:"involved_hazards"`
	Mechanism       string     `json:"mechanism"`
	ClinicalImpact  string     `json:"clinical_impact"`
}

func AllNASAHazards() []NASAHazard {
	return []NASAHazard{
		{
			ID:          HazardRadiation,
			Code:        "HZ-01",
			Name:        "Space Radiation (GCR, SPE & CSA CANS Secondary Neutrons)",
			ShortName:   "Space Radiation",
			Icon:        "☢️",
			ColorClass:  "purple",
			Description: "Galactic Cosmic Rays, Solar Particle Events, and secondary fast/thermal neutrons generated inside spacecraft shielding (monitored via CSA's Canadian Active Neutron Spectrometer).",
			HRPRisks:    []string{"Acute Radiation Sickness (ARS)", "Latent Viral Reactivation", "Radiation-Induced Insulin Resistance", "CNS Cognitive Decrements", "Cataracts & Vascular Endothelial Damage"},
		},
		{
			ID:          HazardIsolation,
			Code:        "HZ-02",
			Name:        "Isolation & Confinement",
			ShortName:   "Isolation & Confinement",
			Icon:        "🧠",
			ColorClass:  "indigo",
			Description: "Circadian desynchronization, confined habitat volume, sensory monotony, and high-workload operational stress.",
			HRPRisks:    []string{"Sleep Loss & Circadian Desynchrony", "Behavioral & Mood Drift", "Cortisol-Driven Immune Suppression", "Team Cohesion & Motor Slowing"},
		},
		{
			ID:          HazardDistance,
			Code:        "HZ-03",
			Name:        "Distance from Earth",
			ShortName:   "Distance from Earth",
			Icon:        "📡",
			ColorClass:  "cyan",
			Description: "Up to 20+ minutes one-way communication latency and no emergency medical evacuation capability beyond Low Earth Orbit.",
			HRPRisks:    []string{"Loss of Real-Time Flight Surgeon Triage", "Restricted Onboard Pharmacy & Diagnostics", "Autonomous Crew Medical Officer (CMO) Burden"},
		},
		{
			ID:          HazardGravity,
			Code:        "HZ-04",
			Name:        "Gravity Fields (0g, 1/6g Lunar & 3/8g Mars Transitions)",
			ShortName:   "Gravity Fields",
			Icon:        "🪐",
			ColorClass:  "blue",
			Description: "Cephalad fluid shift toward the head, mechanical unloading of weight-bearing bones/muscles, and cardiovascular baroreceptor deconditioning.",
			HRPRisks:    []string{"SANS (Spaceflight Associated Neuro-ocular Syndrome)", "1.0–1.5%/mo Bone Mineral Density Loss", "Orthostatic Intolerance & Plasma Volume Loss", "Renal Stone Formation (Hypercalciuria)"},
		},
		{
			ID:          HazardEnvironment,
			Code:        "HZ-05",
			Name:        "Hostile / Closed Environments",
			ShortName:   "Hostile / Closed Env",
			Icon:        "🛡️",
			ColorClass:  "amber",
			Description: "Closed-loop life support (ECLSS) elevated CO₂ pockets, cabin microbiota shifts, and toxic planetary dust (Lunar silica / Martian perchlorates).",
			HRPRisks:    []string{"Hypercapnic Cerebral Vasodilation (Worsens SANS)", "Lunar Dust Pneumonitis & Ocular Irritation", "Altered Microbial Virulence in Closed Loop"},
		},
	}
}

func GetHazardMap() map[HazardID]NASAHazard {
	m := make(map[HazardID]NASAHazard)
	for _, h := range AllNASAHazards() {
		m[h.ID] = h
	}
	return m
}
