package engine

import (
	"fmt"
	"math"
	"orbitodx/internal/domain"
	"time"
)

func EvaluateCrewHealth(vitals domain.VitalsInput, completedSteps map[string]bool, streakDays int) domain.EvaluationResult {
	mission := domain.GetMissionProfile(vitals.MissionID)
	hazards := domain.GetHazardMap()
	series := Generate30DayTelemetry(mission, vitals)
	cans := CalculateCANSNeutron(mission, vitals, series)
	latestSol := series[len(series)-1]

	var cards []domain.ClinicalAlertCard
	var synergies []domain.CompoundSynergy

	// =========================================================================
	// 1. HEAD & EYES: SANS (Spaceflight Associated Neuro-ocular Syndrome)
	// =========================================================================
	sansTier := domain.TierGreen
	sansBadge := "NOMINAL • OPTIC NERVE HOMEOSTASIS"
	if latestSol.SANSRiskIndex >= 72 || vitals.CabinCO2 >= 4.5 {
		sansTier = domain.TierRed
		sansBadge = "RED • ELEVATED SANS / OPTIC DISC EDEMA RISK"
	} else if latestSol.SANSRiskIndex >= 48 || vitals.CabinCO2 >= 3.5 || mission.GravityG == 0.0 && vitals.SystolicBP >= 128 {
		sansTier = domain.TierAmber
		sansBadge = "AMBER • CEPHALAD FLUID SHIFT WATCH"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-sans",
		Title:               "SANS (Spaceflight Associated Neuro-ocular Syndrome) & Intracranial Pressure",
		BiomarkerName:       "SANS Risk Index & Cephalad Venous Congestion",
		BodySystem:          domain.SystemHead,
		BodySystemLabel:     "Head & Eyes (Neuro-Ocular)",
		CurrentValueDisplay: fmt.Sprintf("%d / 100 Index (ppCO₂ %.1f mmHg)", latestSol.SANSRiskIndex, vitals.CabinCO2),
		NominalRangeDisplay: "< 45 Index • Cabin ppCO₂ < 3.2 mmHg",
		Tier:                sansTier,
		TierBadgeText:       sansBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardGravity],
			hazards[domain.HazardEnvironment],
			hazards[domain.HazardDistance],
		},
		CSACANSLinked: false,
		WhyExamNeeded: "In microgravity, astronauts cannot feel intracranial pressure rising until visual acuity blurs or hyperopic (farsighted) shifts occur. Early detection prevents permanent optic nerve sheath distension, globe flattening, and choroidal folds.",
		HowHazardsCause: "Without Earth's gravity pulling blood downward (Gravity Fields hazard), ~2 liters of fluid shift toward the head, preventing normal venous drainage from the brain and eyes. Elevated cabin CO₂ (Hostile/Closed Environment hazard) dilates cerebral blood vessels, compounding optic nerve pressure.",
		EarthClinicAnalogy: "Mirrors point-of-care ocular ultrasound screening for idiopathic intracranial hypertension (IIH) or traumatic brain edema in remote Arctic/rural clinics lacking an MRI.",
		ProtocolTitle:      "NASA HRP SANS & Cephalad Venous Decongestion Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "sans-step-1",
				Order:       1,
				Title:       "Perform Autonomous Ocular Ultrasound & Fundoscopy",
				Instruction: "Measure Optic Nerve Sheath Diameter (ONSD > 5.8 mm threshold) and capture retinal OCT scan for store-and-forward Flight Surgeon review.",
				Equipment:   "Onboard Ultrasound2 + Handheld OCT",
				Timeframe:   "Within 45 mins",
				Completed:   completedSteps["sans-step-1"],
			},
			{
				ID:          "sans-step-2",
				Order:       2,
				Title:       "Initiate Lower Body Negative Pressure (LBNP) Countermeasure",
				Instruction: "Apply -25 to -40 mmHg lower-body vacuum suction to draw venous blood volume back toward the pelvis and legs, relieving jugular/ocular pressure.",
				Equipment:   "Chibis / LBNP Countermeasure Garment",
				Timeframe:   "30 mins session",
				Completed:   completedSteps["sans-step-2"],
			},
			{
				ID:          "sans-step-3",
				Order:       3,
				Title:       "Scrub Cabin CO₂ & Adjust Visual Correction Lenses",
				Instruction: "Verify ECLSS CDRA/LiOH scrubber airflow in crew sleep quarters to bring ppCO₂ below 3.0 mmHg; switch to +0.75D Space Anticipation Glasses if near vision blurs.",
				Equipment:   "ECLSS Air Handler + Adjustable Focus Eyewear",
				Timeframe:   "Immediate",
				Completed:   completedSteps["sans-step-3"],
			},
		},
	})

	// =========================================================================
	// 2. HEAD & EYES: Sleep Duration, Circadian Rhythm & Behavioral Mood
	// =========================================================================
	sleepTier := domain.TierGreen
	sleepBadge := "NOMINAL • CIRCADIAN & COGNITIVE SYNC"
	if vitals.SleepHours < 5.0 || vitals.MoodScore <= 3 {
		sleepTier = domain.TierRed
		sleepBadge = "RED • ACUTE SLEEP DEBT & COGNITIVE DRIFT"
	} else if vitals.SleepHours < 6.5 || vitals.MoodScore <= 6 {
		sleepTier = domain.TierAmber
		sleepBadge = "AMBER • CIRCADIAN PHASE DESYNCHRONY"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-sleep-mood",
		Title:               "Circadian Sleep Architecture & Behavioral Readiness",
		BiomarkerName:       "Daily Sleep Duration & Mood/Stress Score",
		BodySystem:          domain.SystemHead,
		BodySystemLabel:     "Head & Eyes (Behavioral & CNS)",
		CurrentValueDisplay: fmt.Sprintf("%.1f hrs/24h • Mood %d/10", vitals.SleepHours, vitals.MoodScore),
		NominalRangeDisplay: "7.0–8.5 hrs Sleep • Mood 7–10/10",
		Tier:                sleepTier,
		TierBadgeText:       sleepBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardIsolation],
			hazards[domain.HazardDistance],
			hazards[domain.HazardEnvironment],
		},
		CSACANSLinked: false,
		WhyExamNeeded: "Astronauts frequently underestimate their own fatigue during high-stakes operations. Sleep below 6 hours degrades Psychomotor Vigilance Task (PVT) reaction speed by >35%, increasing robotic arm, docking, or EVA error risk while elevating stress cortisol.",
		HowHazardsCause: "Absence of a natural 24-hour solar light/dark cycle, 90-minute orbital sunrises or deep-space monotony (Isolation & Confinement), acoustic fan noise, and 20-minute comm delays isolating crew from family (Distance from Earth) disrupt suprachiasmatic nucleus melatonin signaling.",
		EarthClinicAnalogy: "Directly applicable to shift-work fatigue monitoring for Arctic expedition crews, offshore medical staff, and emergency responders in isolated regions.",
		ProtocolTitle:      "NASA Behavioral Health & SLA Phototherapy Reset Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "sleep-step-1",
				Order:       1,
				Title:       "Execute Solid-State Lighting Assembly (SLA) Blue-Light Phase Reset",
				Instruction: "Set habitat overhead LEDs to 480nm blue-enriched high-melanopic spectrum (6500K, 500 lux) for 45 minutes upon wake; switch to blue-depleted 2700K (<50 lux) 90 minutes pre-sleep.",
				Equipment:   "Habitat Dynamic Solid-State Lighting",
				Timeframe:   "Next Wake/Sleep Cycle",
				Completed:   completedSteps["sleep-step-1"],
			},
			{
				ID:          "sleep-step-2",
				Order:       2,
				Title:       "Run 3-Minute NASA Psychomotor Vigilance Self-Test (Reaction Check)",
				Instruction: "Verify neurocognitive motor readiness prior to operating life-support, EVA suits, or propulsion interfaces; defer high-risk manual tasks if lapse count > 3.",
				Equipment:   "OrbitoDX Crew Tablet PVT Module",
				Timeframe:   "Within 30 mins",
				Completed:   completedSteps["sleep-step-2"],
			},
			{
				ID:          "sleep-step-3",
				Order:       3,
				Title:       "Administer Non-Sedating Melatonin & Acoustic Sleep Pod Check",
				Instruction: "Take 0.5–2.0 mg timed oral melatonin 60 minutes before scheduled sleep block and verify sleep station acoustic baffling (<45 dBA) and thermal cooling (19.5°C).",
				Equipment:   "Crew Medical Kit (Oral Melatonin)",
				Timeframe:   "Pre-Sleep Block",
				Completed:   completedSteps["sleep-step-3"],
			},
		},
	})

	// =========================================================================
	// 3. HEART & CHEST: Cardiovascular Blood Pressure, Heart Rate & Orthostasis
	// =========================================================================
	cardioTier := domain.TierGreen
	cardioBadge := "NOMINAL • HEMODYNAMIC STABILITY"
	if vitals.SystolicBP < 95 || vitals.SystolicBP >= 140 || vitals.DiastolicBP < 60 || vitals.DiastolicBP >= 90 || vitals.HeartRate >= 92 {
		cardioTier = domain.TierRed
		if vitals.SystolicBP < 95 {
			cardioBadge = "RED • HYPOTENSIVE ORTHOSTATIC DECONDITIONING"
		} else {
			cardioBadge = "RED • SYMPATHETIC HEMODYNAMIC STRAIN"
		}
	} else if vitals.SystolicBP < 104 || vitals.SystolicBP >= 128 || vitals.HeartRate >= 80 {
		cardioTier = domain.TierAmber
		cardioBadge = "AMBER • BARORECEPTOR & PLASMA VOLUME DRIFT"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-cardio",
		Title:               "Cardiovascular Adaptation, Blood Pressure & Orthostatic Tolerance",
		BiomarkerName:       "Arterial Blood Pressure & Resting Heart Rate",
		BodySystem:          domain.SystemHeart,
		BodySystemLabel:     "Heart & Chest (Cardiovascular)",
		CurrentValueDisplay: fmt.Sprintf("%d/%d mmHg • HR %d bpm", vitals.SystolicBP, vitals.DiastolicBP, vitals.HeartRate),
		NominalRangeDisplay: "105–125 / 65–82 mmHg • HR 55–78 bpm",
		Tier:                cardioTier,
		TierBadgeText:       cardioBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardGravity],
			hazards[domain.HazardRadiation],
			hazards[domain.HazardDistance],
		},
		CSACANSLinked: false,
		WhyExamNeeded: "A drop in blood pressure with rising heart rate signals shrinking blood plasma volume and weakened carotid baroreceptors—putting an astronaut at severe risk of fainting (presyncope) during Lunar/Mars gravity landing or emergency egress. Conversely, elevated BP + tachycardia indicates acute sympathetic stress or vascular stiffness.",
		HowHazardsCause: "In 0g or 1/6g (Gravity Fields hazard), fluid pooling in the chest tricks the heart's stretch receptors into thinking the body has excess fluid, triggering diuresis that drops blood plasma volume by up to 15–20% and causes left ventricular muscle atrophy. Space radiation also induces oxidative endothelial stiffening.",
		EarthClinicAnalogy: "Guides non-physician rural paramedics in differentiating hypovolemic shock, dehydration, and autonomic dysreflexia without invasive arterial lines.",
		ProtocolTitle:      "NASA Cardiovascular Volume Restoration & Baroreceptor Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "cardio-step-1",
				Order:       1,
				Title:       "Execute Isotonic Fluid & Sodium Loading Bolus",
				Instruction: "Consume 500 mL water with 4x 1g NaCl electrolyte tablets (or NASA Astro-Ade isotonic packet) to rapidly expand intravascular plasma volume.",
				Equipment:   "Potable Water Dispenser + NaCl Med-Kit",
				Timeframe:   "Within 30 mins",
				Completed:   completedSteps["cardio-step-1"],
			},
			{
				ID:          "cardio-step-2",
				Order:       2,
				Title:       "Perform 25-Min High-Intensity Interval Cycle Ergometer (CEVIS) Block",
				Instruction: "Execute 4×4 min aerobic intervals at 80% VO₂max peak heart rate to stimulate baroreceptor sensitivity and maintain left ventricular stroke volume.",
				Equipment:   "CEVIS / T2 Treadmill Harness",
				Timeframe:   "Today's Exercise Window",
				Completed:   completedSteps["cardio-step-2"],
			},
			{
				ID:          "cardio-step-3",
				Order:       3,
				Title:       "Record 3-Lead Resting ECG & Compression Garment Readiness",
				Instruction: "Capture 60-second rhythm strip to rule out spaceflight atrial/ventricular ectopy; stage Kentavr/Thigh-Cuff compression garments if systolic BP remains < 98 mmHg.",
				Equipment:   "Bio-Monitor Wearable ECG Lead",
				Timeframe:   "Within 1 hour",
				Completed:   completedSteps["cardio-step-3"],
			},
		},
	})

	// =========================================================================
	// 4. HEART & CHEST: Cabin CO2 Hypercapnia & Lunar Dust Respiratory Triage
	// =========================================================================
	respTier := domain.TierGreen
	respBadge := "NOMINAL • ECLSS ATMOSPHERE & AIRWAY CLEAR"
	if vitals.CabinCO2 >= 4.5 {
		respTier = domain.TierRed
		respBadge = "RED • HYPERCAPNIC CABIN CO₂POCKET (>4.5 mmHg)"
	} else if vitals.CabinCO2 >= 3.4 || mission.ID == domain.MissionArtemis {
		respTier = domain.TierAmber
		if mission.ID == domain.MissionArtemis && vitals.CabinCO2 < 3.4 {
			respBadge = "AMBER • LUNAR SILICA REGOLITH AIRWAY WATCH"
		} else {
			respBadge = "AMBER • ELEVATED CABIN ppCO₂ (3.4–4.4 mmHg)"
		}
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-respiratory",
		Title:               "Closed-Loop Cabin CO₂ & Planetary Dust Respiratory Triage",
		BiomarkerName:       "Cabin ppCO₂ & Alveolar Ventilation Risk",
		BodySystem:          domain.SystemHeart,
		BodySystemLabel:     "Heart & Chest (Respiratory & ECLSS)",
		CurrentValueDisplay: fmt.Sprintf("ppCO₂ %.1f mmHg (%s)", vitals.CabinCO2, mission.ShortCode),
		NominalRangeDisplay: "ppCO₂ < 3.0 mmHg (Earth Sea Level ~0.3 mmHg)",
		Tier:                respTier,
		TierBadgeText:       respBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardEnvironment],
			hazards[domain.HazardGravity],
		},
		CSACANSLinked: false,
		WhyExamNeeded: "Without buoyancy-driven air convection in microgravity, exhaled CO₂ forms invisible pockets around an astronaut's face—causing frontal headaches, slowed decision-making, and cerebral vasodilation. On Artemis Lunar missions, electrostatically charged <2 µm sharp silica dust clings to EVA suits and can inflame deep lung alveoli.",
		HowHazardsCause: "Closed-loop life support (Hostile/Closed Environments) recycles air at 10x Earth's normal CO₂ concentration, while 0g eliminates warm-air rising currents (Gravity Fields). On the Moon, unweathered micrometeorite-shattered regolith contains reactive iron nanophase particles.",
		EarthClinicAnalogy: "Parallels indoor air quality, carbon monoxide/CO₂ monitoring, and silicosis/wildfire smoke triage in mining towns and remote occupational health posts.",
		ProtocolTitle:      "ECLSS Atmospheric Scrubbing & Airway Protection Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "resp-step-1",
				Order:       1,
				Title:       "Activate Personal Ventilation Fan & Inspect CDRA Zeolite Bed",
				Instruction: "Position personal airflow nozzle across crew workstation/bunk to disperse stagnant exhaled CO₂ boundary layers; switch secondary LiOH/amine scrubber bed online.",
				Equipment:   "ECLSS CDRA / Personal Air Duct",
				Timeframe:   "Immediate",
				Completed:   completedSteps["resp-step-1"],
			},
			{
				ID:          "resp-step-2",
				Order:       2,
				Title:       "Verify Pulse Oximetry (SpO₂) & Peak Expiratory Flow",
				Instruction: "Confirm SpO₂ ≥ 97% and perform handheld spirometry check after any Lunar EVA airlock ingress or filter maintenance.",
				Equipment:   "Finger Pulse Oximeter + Spirometer",
				Timeframe:   "Within 20 mins",
				Completed:   completedSteps["resp-step-2"],
			},
			{
				ID:          "resp-step-3",
				Order:       3,
				Title:       "Run HEPA Electrostatic Dust Cycle & Saline Nasal Lavage (Lunar Mode)",
				Instruction: "If operating in Artemis Lunar habitat, run suitlock HEPA purge for 15 mins and perform isotonic saline nasal rinse to clear ultrafine regolith particulates.",
				Equipment:   "Airlock HEPA Scrubber + Saline Kit",
				Timeframe:   "Post-EVA / Daily",
				Completed:   completedSteps["resp-step-3"],
			},
		},
	})

	// =========================================================================
	// 5. SKELETON & MUSCLES: Bone Density Loss & ARED Resistive Compliance
	// =========================================================================
	boneTier := domain.TierGreen
	boneBadge := "NOMINAL • OSTEOBLAST MECHANO-STIMULATION"
	if vitals.AREDCompliance < 50 || latestSol.BoneLossRatePct >= 1.25 {
		boneTier = domain.TierRed
		boneBadge = "RED • RAPID OSTEOCLAST BONE LOSS & HYPERCALCIURIA RISK"
	} else if vitals.AREDCompliance < 80 || latestSol.BoneLossRatePct >= 0.70 {
		boneTier = domain.TierAmber
		boneBadge = "AMBER • SUB-THRESHOLD ARED RESISTIVE LOADING"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-bone-muscle",
		Title:               "Spaceflight Bone Mineral Density (BMD) & ARED Exercise Triage",
		BiomarkerName:       "Est. Monthly Bone Loss & ARED Load Compliance",
		BodySystem:          domain.SystemSkeleton,
		BodySystemLabel:     "Skeleton & Muscles (Musculoskeletal)",
		CurrentValueDisplay: fmt.Sprintf("-%.2f%% BMD/mo • ARED %d%%", latestSol.BoneLossRatePct, vitals.AREDCompliance),
		NominalRangeDisplay: "< 0.45% BMD/mo • ARED Compliance ≥ 85%",
		Tier:                boneTier,
		TierBadgeText:       boneBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardGravity],
			hazards[domain.HazardRadiation],
			hazards[domain.HazardDistance],
		},
		CSACANSLinked: false,
		WhyExamNeeded: "In microgravity, astronauts lose 1.0% to 1.5% of weight-bearing bone mass every single month—10x faster than terrestrial osteoporosis! Worse, dissolved bone calcium floods the bloodstream and kidneys, creating a high risk of mission-disabling kidney stones (nephrolithiasis) millions of kilometers from Earth.",
		HowHazardsCause: "Mechanical unloading in 0g and 1/6g (Gravity Fields) shuts down bone-building osteoblasts while upregulating bone-resorbing osteoclasts in the proximal femur, pelvis, and lumbar spine. Deep-space radiation further damages bone marrow osteoprogenitor stem cells.",
		EarthClinicAnalogy: "Directly empowers physical therapy and fracture/renal-stone prevention for bedridden, immobilized, or frail elderly patients in rural home-care settings.",
		ProtocolTitle:      "NASA ARED High-Load Resistive & Renal Stone Prevention Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "bone-step-1",
				Order:       1,
				Title:       "Complete ARED Axial Loading Block (Squat, Deadlift & Heel Raise)",
				Instruction: "Perform 4 sets × 8 reps Back Squat, Romanian Deadlift, and Smith-bar Heel Raises at 1.5x–1.7x body-weight vacuum cylinder resistance to trigger femoral/spinal piezoelectrical bone remodeling.",
				Equipment:   "ARED (Advanced Resistive Exercise Device)",
				Timeframe:   "60 mins Daily Block",
				Completed:   completedSteps["bone-step-1"],
			},
			{
				ID:          "bone-step-2",
				Order:       2,
				Title:       "Administer Vitamin D3 (2000 IU) + Potassium Citrate Renal Prophylaxis",
				Instruction: "Take daily Cholecalciferol (2000 IU) with Vitamin K2; if ARED compliance < 65%, initiate oral Potassium Citrate (20 mEq) to alkalinize urine and prevent calcium oxalate stone crystallization.",
				Equipment:   "Nutritional Countermeasure Pack",
				Timeframe:   "With Mid-Sol Meal",
				Completed:   completedSteps["bone-step-2"],
			},
			{
				ID:          "bone-step-3",
				Order:       3,
				Title:       "Maintain > 2.5 L/day Hydration Target & Urine Specific Gravity Check",
				Instruction: "Verify urine test strip specific gravity < 1.015 to dilute urinary calcium excretion during skeletal unloading phases.",
				Equipment:   "Colorimetric Urine Dipstick + Water Dispenser",
				Timeframe:   "Evening Check-In",
				Completed:   completedSteps["bone-step-3"],
			},
		},
	})

	// =========================================================================
	// 6. BLOOD, IMMUNE & RADIATION: CSA CANS Neutron Radiation Dosimetry
	// =========================================================================
	radTier := domain.TierGreen
	radBadge := "NOMINAL • SHIELDED GCR/NEUTRON BACKGROUND"
	if vitals.RadiationMSv >= 3.5 || cans.Cumulative30DayMSv >= 90.0 {
		radTier = domain.TierRed
		radBadge = "RED • CSA CANS SOLAR NEUTRON STORM ALERT"
	} else if vitals.RadiationMSv >= 1.5 || cans.Cumulative30DayMSv >= 45.0 {
		radTier = domain.TierAmber
		radBadge = "AMBER • ELEVATED DEEP-SPACE NEUTRON DOSE"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-radiation-cans",
		Title:               "CSA CANS Neutron Radiation Dosimetry & Shielding Triage",
		BiomarkerName:       "Total Daily Dose & CSA CANS Secondary Neutron Fraction",
		BodySystem:          domain.SystemBlood,
		BodySystemLabel:     "Blood, Immune & Radiation (CSA CANS)",
		CurrentValueDisplay: fmt.Sprintf("%.2f mSv/d (CANS Neutrons: %.2f mSv • %d%%)", cans.DailyTotalMSv, cans.DailyNeutronMSv, cans.NeutronPercentage),
		NominalRangeDisplay: fmt.Sprintf("< 1.20 mSv/d • 30-Day Cum: %.1f / %.0f mSv Limit", cans.Cumulative30DayMSv, cans.NASALimit30DayMSv),
		Tier:                radTier,
		TierBadgeText:       radBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardRadiation],
			hazards[domain.HazardDistance],
			hazards[domain.HazardEnvironment],
		},
		CSACANSLinked: true,
		WhyExamNeeded: "Neutrons have zero electrical charge, meaning they pass undetected through standard charged-particle dosimeters yet carry up to 20x the biological DNA-damaging power (RBE) of gamma rays! As established by the Canadian Space Agency's CANS (Canadian Active Neutron Spectrometer), secondary neutrons account for ~30% of the ISS radiation dose and over 40% in deep-space habitats.",
		HowHazardsCause: "When high-energy Galactic Cosmic Rays (GCR) or Solar Particle Event (SPE) protons slam into a spacecraft's aluminum hull or lunar soil (Space Radiation hazard), atomic nuclei shatter ('spallation') and spray secondary fast neutrons inward into the crew cabin.",
		EarthClinicAnalogy: "Directly supports radiation safety triage for nuclear medicine clinics, high-altitude aviation crews, and emergency radiological response teams.",
		ProtocolTitle:      "CSA CANS Neutron Attenuation & Storm Shelter Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "rad-step-1",
				Order:       1,
				Title:       "Relocate Crew to Hydrogen-Rich Polyethylene / Water-Wall Storm Shelter",
				Instruction: "Move immediately from outer aluminum hull modules to the central crew shelter surrounded by potable water bags, food stowage, and HDPE blocks (hydrogen nuclei effectively slow and capture fast neutrons).",
				Equipment:   "Aft Water-Wall Storm Shelter Bay",
				Timeframe:   "Immediate (< 10 mins)",
				Completed:   completedSteps["rad-step-1"],
			},
			{
				ID:          "rad-step-2",
				Order:       2,
				Title:       "Cross-Check CSA CANS Fast-Neutron Spectrum (100 keV – 15 MeV)",
				Instruction: "Log CANS silicon/moderator channel ratio (`Fast Flux: " + fmt.Sprintf("%.2f", cans.FastNeutronFlux) + " mSv/d`) and wear personal AstroRad vest shielding blood-forming bone marrow organs (pelvis/sternum).",
				Equipment:   "CSA CANS-II Console + AstroRad Vest",
				Timeframe:   "Within 15 mins",
				Completed:   completedSteps["rad-step-2"],
			},
			{
				ID:          "rad-step-3",
				Order:       3,
				Title:       "Administer Radioprotective Antioxidant Countermeasure Pack",
				Instruction: "Initiate oral N-Acetylcysteine (NAC 600 mg), Vitamin C/E, and CoQ10 protocol to scavenge neutron-induced reactive oxygen species (ROS) in hematopoietic stem cells.",
				Equipment:   "Radiation Medical Countermeasure Kit",
				Timeframe:   "Within 30 mins",
				Completed:   completedSteps["rad-step-3"],
			},
		},
	})

	// =========================================================================
	// 7. BLOOD, IMMUNE & RADIATION: Blood Sugar & Immune / Viral Reactivation
	// =========================================================================
	immTier := domain.TierGreen
	immBadge := "NOMINAL • METABOLIC & IMMUNE HOMEOSTASIS"
	if vitals.BloodSugar >= 140 || vitals.BloodSugar < 68 || vitals.BodyTemp >= 37.6 {
		immTier = domain.TierRed
		if vitals.BodyTemp >= 37.6 {
			immBadge = "RED • IMMUNE DYSREGULATION & VIRAL REACTIVATION FEVER"
		} else {
			immBadge = "RED • SPACEFLIGHT HYPERGLYCEMIC INSULIN RESISTANCE"
		}
	} else if vitals.BloodSugar >= 110 || vitals.BodyTemp >= 37.3 {
		immTier = domain.TierAmber
		immBadge = "AMBER • SUBCLINICAL INSULIN RESISTANCE / IMMUNE STRESS"
	}

	cards = append(cards, domain.ClinicalAlertCard{
		ID:                  "card-metabolic-immune",
		Title:               "Spaceflight Glucose Metabolism & Latent Viral Immune Triage",
		BiomarkerName:       "Blood Glucose, Core Temp & Viral Reactivation Risk",
		BodySystem:          domain.SystemBlood,
		BodySystemLabel:     "Blood, Immune & Radiation (Metabolic/Immune)",
		CurrentValueDisplay: fmt.Sprintf("%d mg/dL Glucose • %.1f°C Core Temp", vitals.BloodSugar, vitals.BodyTemp),
		NominalRangeDisplay: "75–105 mg/dL Glucose • 36.3–37.2°C Temp",
		Tier:                immTier,
		TierBadgeText:       immBadge,
		HazardTags: []domain.NASAHazard{
			hazards[domain.HazardRadiation],
			hazards[domain.HazardIsolation],
			hazards[domain.HazardGravity],
		},
		CSACANSLinked: true,
		WhyExamNeeded: "Even healthy astronauts develop pre-diabetic 'spaceflight insulin resistance' (elevated blood sugar) due to muscle unloading and radiation oxidative stress. Simultaneously, T-cell exhaustion allows dormant viruses (Epstein-Barr, Varicella-Zoster/Shingles, CMV) to reactivate, signaled by low-grade fever (≥37.5°C) and fatigue.",
		HowHazardsCause: "Muscle unloading in 0g reduces GLUT4 glucose uptake (Gravity Fields), while neutron/GCR radiation (Space Radiation) and confinement cortisol spikes (Isolation & Confinement) impair pancreatic beta-cell signaling and suppress cytotoxic CD8+ T-lymphocytes.",
		EarthClinicAnalogy: "Provides point-of-care sepsis, diabetic ketoacidosis, and immunocompromised fever triage for community health workers in resource-limited rural clinics.",
		ProtocolTitle:      "NASA Metabolic Glycemic Stabilization & Immune Viral Protocol",
		ProtocolSteps: []domain.ProtocolStep{
			{
				ID:          "imm-step-1",
				Order:       1,
				Title:       "Perform Point-of-Care Capillary Metabolic & Salivary Cortisol/Viral Assay",
				Instruction: "Run fingerstick ketone/HbA1c check and wet-lab salivary PCR swab for latent EBV/VZV shedding; restrict high-glycemic simple carbohydrate rations.",
				Equipment:   "i-STAT Blood Analyzer + Saliva Swab Kit",
				Timeframe:   "Within 30 mins",
				Completed:   completedSteps["imm-step-1"],
			},
			{
				ID:          "imm-step-2",
				Order:       2,
				Title:       "Execute 20-Min Postprandial GLUT4 Muscle Glucose Uptake Walk/Cycle",
				Instruction: "Perform light aerobic cycling within 30 minutes after meals to stimulate non-insulin-dependent skeletal muscle GLUT4 glucose clearance.",
				Equipment:   "CEVIS Ergometer",
				Timeframe:   "Post-Meal Window",
				Completed:   completedSteps["imm-step-2"],
			},
			{
				ID:          "imm-step-3",
				Order:       3,
				Title:       "Initiate Antiviral Prophylaxis & Antipyretic Protocol (If Temp ≥ 37.6°C)",
				Instruction: "If low-grade fever persists with vesicular rash or malaise, Crew Medical Officer initiates Valacyclovir 1000 mg PO q8h and queues encrypted biomedical telemetry packet for Flight Surgeon.",
				Equipment:   "CMO Antimicrobial Locker",
				Timeframe:   "If Temp ≥ 37.6°C",
				Completed:   completedSteps["imm-step-3"],
			},
		},
	})

	// =========================================================================
	// DETECT COMPOUND HAZARD SYNERGIES (Multi-Hazard Compounding Loops)
	// =========================================================================
	if radTier != domain.TierGreen && (sleepTier != domain.TierGreen || immTier != domain.TierGreen) {
		synergies = append(synergies, domain.CompoundSynergy{
			ID:              "syn-rad-immune",
			Title:           "COMPOUND LOOP: Neutron Radiation × Sleep Debt → Accelerated Viral & Glycemic Drift",
			Severity:        domain.TierRed,
			InvolvedHazards: []domain.HazardID{domain.HazardRadiation, domain.HazardIsolation, domain.HazardDistance},
			Mechanism:       "CSA CANS secondary neutron oxidative stress combined with circadian sleep loss (<6h) synergistically elevates serum cortisol, suppressing NK/T-cell immunity and spiking blood glucose to " + fmt.Sprintf("%d mg/dL.", vitals.BloodSugar),
			ClinicalImpact:  "Doubles latent herpesvirus (EBV/VZV) reactivation probability while delaying DNA double-strand break repair during Mars comm delay.",
		})
	}
	if sansTier != domain.TierGreen && (respTier != domain.TierGreen || cardioTier != domain.TierGreen) {
		synergies = append(synergies, domain.CompoundSynergy{
			ID:              "syn-sans-co2",
			Title:           "COMPOUND LOOP: 0g Cephalad Fluid Shift × Cabin CO₂ → Amplified SANS Optic Edema",
			Severity:        sansTier,
			InvolvedHazards: []domain.HazardID{domain.HazardGravity, domain.HazardEnvironment},
			Mechanism:       fmt.Sprintf("Microgravity venous pooling in the head compounded by %.1f mmHg cabin CO₂ cerebral vasodilation increases intracranial pressure.", vitals.CabinCO2),
			ClinicalImpact:  "Accelerates optic nerve sheath distension and globe flattening; requires simultaneous LBNP suction and ECLSS CO₂ scrubbing.",
		})
	}
	if boneTier != domain.TierGreen && cardioTier != domain.TierGreen {
		synergies = append(synergies, domain.CompoundSynergy{
			ID:              "syn-bone-cardio",
			Title:           "COMPOUND LOOP: Musculoskeletal Unloading × Hypovolemia → Gravity-Transition Syncope & Renal Stone Risk",
			Severity:        domain.TierRed,
			InvolvedHazards: []domain.HazardID{domain.HazardGravity, domain.HazardDistance},
			Mechanism:       fmt.Sprintf("Low ARED resistive loading (%d%%) plus reduced plasma volume (BP %d/%d mmHg) concentrates urinary calcium while weakening lower-limb vascular tone.", vitals.AREDCompliance, vitals.SystolicBP, vitals.DiastolicBP),
			ClinicalImpact:  "High risk of both acute renal colic in deep space and orthostatic fainting upon planetary landing.",
		})
	}
	if len(synergies) == 0 {
		synergies = append(synergies, domain.CompoundSynergy{
			ID:              "syn-nominal",
			Title:           "ACTIVE MULTI-HAZARD EQUILIBRIUM: All 5 NASA Hazards Controlled Within Countermeasure Envelope",
			Severity:        domain.TierGreen,
			InvolvedHazards: []domain.HazardID{domain.HazardGravity, domain.HazardRadiation, domain.HazardIsolation, domain.HazardEnvironment, domain.HazardDistance},
			Mechanism:       "ARED resistive loading, ECLSS CO₂ scrubbing, circadian lighting, and polyethylene shielding are successfully counteracting physiological drift.",
			ClinicalImpact:  "Crew maintains 90%+ mission readiness with low autonomous intervention burden.",
		})
	}

	// =========================================================================
	// AGGREGATE BODY SYSTEM SUMMARIES & OVERALL TRIAGE TIER
	// =========================================================================
	greenCount, amberCount, redCount := 0, 0, 0
	for _, c := range cards {
		switch c.Tier {
		case domain.TierGreen:
			greenCount++
		case domain.TierAmber:
			amberCount++
		case domain.TierRed:
			redCount++
		}
	}

	overallTier := domain.TierGreen
	overallHeadline := "ALL CREW SYSTEMS NOMINAL — HOMEOSTASIS MAINTAINED"
	overallSubtext := "Vital signs and CSA CANS neutron dosimetry are within NASA HRP spaceflight norms. Continue scheduled daily countermeasures."
	if redCount > 0 {
		overallTier = domain.TierRed
		overallHeadline = fmt.Sprintf("RED ESCALATION — %d CRITICAL ANOMALIES REQUIRE AUTONOMOUS ACTION", redCount)
		if mission.CommDelayOneWaySec > 60 {
			overallSubtext = fmt.Sprintf("MARS COMM DELAY ACTIVE (%s Round-Trip): Do NOT wait for Houston Flight Surgeon. Execute onboard CMO countermeasure protocols immediately; store-and-forward packet queued.", mission.CommDelayRoundTripLabel)
		} else {
			overallSubtext = "Immediate Crew Medical Officer intervention required. Telemetry packet flagged for priority Flight Surgeon console review."
		}
	} else if amberCount > 0 {
		overallTier = domain.TierAmber
		overallHeadline = fmt.Sprintf("AMBER WATCH — %d EARLY DRIFT INDICATORS DETECTED", amberCount)
		overallSubtext = "Early physiological or behavioral drift spotted. Execute preventative self-intervention countermeasures within 2–6 hours to prevent Red escalation."
	}

	headTier := maxTier(cards[0].Tier, cards[1].Tier)
	heartTier := maxTier(cards[2].Tier, cards[3].Tier)
	skeletonTier := cards[4].Tier
	bloodTier := maxTier(cards[5].Tier, cards[6].Tier)

	bodySystems := []domain.BodySystemSummary{
		{
			ID:             domain.SystemHead,
			Name:           "1. Head & Eyes",
			Subtitle:       "Neuro-Ocular (SANS) & Behavioral",
			Tier:           headTier,
			ActiveAlerts:   countNonGreen(cards[0].Tier, cards[1].Tier),
			HeadlineMetric: fmt.Sprintf("SANS %d/100 • Sleep %.1fh • Mood %d/10", latestSol.SANSRiskIndex, vitals.SleepHours, vitals.MoodScore),
			SecondaryNote:  "Cephalad fluid shift, optic nerve sheath & circadian phase",
			PrimaryHazards: []string{"Gravity Fields", "Isolation & Confinement", "Hostile Env"},
		},
		{
			ID:             domain.SystemHeart,
			Name:           "2. Heart & Chest",
			Subtitle:       "Cardiovascular & Respiratory",
			Tier:           heartTier,
			ActiveAlerts:   countNonGreen(cards[2].Tier, cards[3].Tier),
			HeadlineMetric: fmt.Sprintf("BP %d/%d mmHg • HR %d bpm • CO₂ %.1f", vitals.SystolicBP, vitals.DiastolicBP, vitals.HeartRate, vitals.CabinCO2),
			SecondaryNote:  "Baroreceptor reflex, plasma volume & cabin ppCO₂/dust",
			PrimaryHazards: []string{"Gravity Fields", "Hostile/Closed Env", "Distance"},
		},
		{
			ID:             domain.SystemSkeleton,
			Name:           "3. Skeleton & Muscles",
			Subtitle:       "Musculoskeletal & ARED Loading",
			Tier:           skeletonTier,
			ActiveAlerts:   countNonGreen(cards[4].Tier),
			HeadlineMetric: fmt.Sprintf("BMD -%.2f%%/mo • ARED %d%%", latestSol.BoneLossRatePct, vitals.AREDCompliance),
			SecondaryNote:  "Femoral/lumbar osteoclast activity & renal stone risk",
			PrimaryHazards: []string{"Gravity Fields", "Space Radiation"},
		},
		{
			ID:             domain.SystemBlood,
			Name:           "4. Blood, Immune & Radiation",
			Subtitle:       "CSA CANS Neutrons, Glucose & Immunity",
			Tier:           bloodTier,
			ActiveAlerts:   countNonGreen(cards[5].Tier, cards[6].Tier),
			HeadlineMetric: fmt.Sprintf("CANS %.2f mSv/d • Glu %d mg/dL • %.1f°C", cans.DailyNeutronMSv, vitals.BloodSugar, vitals.BodyTemp),
			SecondaryNote:  fmt.Sprintf("CSA CANS Neutron Share: %d%% of dose • Viral/Metabolic index", cans.NeutronPercentage),
			PrimaryHazards: []string{"Space Radiation (CANS)", "Isolation", "Distance"},
		},
	}

	// Calculate Readiness Score with bonus from completed protocol steps
	completedCount := 0
	for _, v := range completedSteps {
		if v {
			completedCount++
		}
	}
	readiness := latestSol.ReadinessScore + (completedCount * 3)
	if readiness > 100 {
		readiness = 100
	}

	streakTargets := []domain.StreakTarget{
		{
			ID:        "target-vitals",
			Label:     "Autonomous Morning Vitals & Bio-Monitor Sync",
			Detail:    "BP, HR, Glucose, Core Temp & Sleep logged on Sol 30",
			Completed: true,
			Points:    25,
		},
		{
			ID:        "target-cans",
			Label:     "CSA CANS Neutron Spectrometer Dosimetry Check",
			Detail:    fmt.Sprintf("Neutron fraction verified (%d%% • %.2f mSv/d)", cans.NeutronPercentage, cans.DailyNeutronMSv),
			Completed: true,
			Points:    25,
		},
		{
			ID:        "target-ared",
			Label:     "ARED Resistive Exercise & Hydration Target (≥80%)",
			Detail:    fmt.Sprintf("Current compliance: %d%% axial load", vitals.AREDCompliance),
			Completed: vitals.AREDCompliance >= 80 || completedSteps["bone-step-1"],
			Points:    25,
		},
		{
			ID:        "target-countermeasure",
			Label:     "Active Triage Countermeasure Protocol Execution",
			Detail:    fmt.Sprintf("%d onboard protocol steps verified complete", completedCount),
			Completed: overallTier == domain.TierGreen || completedCount >= 2,
			Points:    25,
		},
	}

	packetID := fmt.Sprintf("MEDB-%s-SOL30-%d", string(mission.ID), int(math.Abs(float64(vitals.SystolicBP*17+vitals.BloodSugar*13)))%9000+1000)

	return domain.EvaluationResult{
		Timestamp:               time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		Mission:                 mission,
		Vitals:                  vitals,
		OverallTier:             overallTier,
		OverallHeadline:         overallHeadline,
		OverallSubtext:          overallSubtext,
		GreenCount:              greenCount,
		AmberCount:              amberCount,
		RedCount:                redCount,
		ReadinessScore:          readiness,
		StreakDays:              streakDays,
		StreakTargets:           streakTargets,
		BodySystems:             bodySystems,
		AlertCards:              cards,
		CompoundSynergies:       synergies,
		CANSNeutron:             cans,
		Telemetry30Days:         series,
		StoreAndForwardPacketID: packetID,
	}
}

func maxTier(tiers ...domain.TriageTier) domain.TriageTier {
	best := domain.TierGreen
	for _, t := range tiers {
		if t == domain.TierRed {
			return domain.TierRed
		}
		if t == domain.TierAmber {
			best = domain.TierAmber
		}
	}
	return best
}

func countNonGreen(tiers ...domain.TriageTier) int {
	c := 0
	for _, t := range tiers {
		if t != domain.TierGreen {
			c++
		}
	}
	return c
}
