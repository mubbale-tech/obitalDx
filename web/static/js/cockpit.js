// OrbitoDX Spacecraft Clinical Cockpit Controller (Alpine.js + Chart.js v4 + Offline Canvas Fallback)
let telemetryChartInstance = null;

window.orbitodxApp = function () {
  const bootEl = document.getElementById("orbitodx-boot-data");
  const boot = bootEl ? JSON.parse(bootEl.textContent) : {};

  return {
    state: boot.evaluation || {},
    missions: boot.missions || [],
    hazards: boot.hazards || [],
    scenarios: boot.scenarios || [],
    form: Object.assign({}, (boot.evaluation && boot.evaluation.vitals) || {}),
    activeSystem: "all",
    showAlertsOnly: false,
    selectedHazard: null,
    chartMode: "cardio",
    earthModalOpen: false,
    pitchBarOpen: true,
    activePitchBeat: 0,
    commSecondsRemaining: 1724, // 28m 44s round-trip in Mars mode
    commCountdownDisplay: "28m 44s (No Live Ground Link)",

    pitchBeats: [
      {
        label: "1. Hook (20m Delay)",
        talkingPoint:
          "Beat 1 — On a Mars mission, a 14–20m one-way delay (~29m round-trip) means astronauts cannot wait for Houston Flight Surgeons when vitals shift.",
      },
      {
        label: "2. Solution (5 Hazards + CANS)",
        talkingPoint:
          "Beat 2 — OrbitoDX maps vitals & CSA CANS secondary neutron dosimetry (~30% ISS / 42% Mars dose) directly to NASA HRP's 5 Hazards.",
      },
      {
        label: "3. Live Mars Anomaly Demo",
        talkingPoint:
          "Beat 3 — Triggering Mars Solar Neutron Spike flips indicators to RED, explains why in plain English, and prescribes Storm Shelter + protocols.",
      },
      {
        label: "4. NASA & CSA Validity",
        talkingPoint:
          "Beat 4 — Grounded in NASA HRP (SANS, Bone Loss, Baroreceptor Drift, Viral Reactivation) & CSA CANS fast-neutron spallation physics.",
      },
      {
        label: "5. Earth Rural Spin-Off",
        talkingPoint:
          "Beat 5 — The same 15 MB offline-first binary empowers remote Northern clinics and expedition medics to triage without specialist access.",
      },
    ],

    init() {
      this.syncFormFromState();
      this.$nextTick(() => {
        this.renderTelemetryChart();
      });

      // Live Houston Round-Trip Comm Delay Countdown ticker
      setInterval(() => {
        if (this.state.mission && this.state.mission.comm_delay_one_way_sec > 60) {
          this.commSecondsRemaining =
            this.commSecondsRemaining > 1 ? this.commSecondsRemaining - 1 : 1724;
          const mins = Math.floor(this.commSecondsRemaining / 60);
          const secs = this.commSecondsRemaining % 60;
          this.commCountdownDisplay = `${mins}m ${secs < 10 ? "0" : ""}${secs}s until earliest Earth reply`;
        } else {
          this.commCountdownDisplay = "Real-Time Ground Link (< 3s)";
        }
      }, 1000);
    },

    syncFormFromState() {
      if (this.state && this.state.vitals) {
        this.form = Object.assign({}, this.state.vitals);
      }
    },

    async switchMission(missionId) {
      try {
        const res = await fetch(`/api/mission/${missionId}`, { method: "POST" });
        if (res.ok) {
          this.state = await res.json();
          this.syncFormFromState();
          this.commSecondsRemaining = Math.round(
            (this.state.mission.comm_delay_one_way_sec || 0) * 2
          );
          this.renderTelemetryChart();
        }
      } catch (e) {
        console.error("Failed to switch mission:", e);
      }
    },

    async triggerScenario(scenarioId) {
      try {
        const res = await fetch(`/api/scenario/${scenarioId}`, { method: "POST" });
        if (res.ok) {
          this.state = await res.json();
          this.syncFormFromState();
          this.commSecondsRemaining = Math.round(
            (this.state.mission.comm_delay_one_way_sec || 0) * 2
          );
          if (scenarioId === "mars_neutron_spike") {
            this.chartMode = "radiation";
          } else if (scenarioId === "bone_cardio_drift") {
            this.chartMode = "bone";
          } else if (scenarioId === "circadian_sleep_deficit") {
            this.chartMode = "neuro";
          } else {
            this.chartMode = "cardio";
          }
          this.renderTelemetryChart();
        }
      } catch (e) {
        console.error("Failed to trigger scenario:", e);
      }
    },

    async submitVitals() {
      try {
        const payload = Object.assign({}, this.form, {
          mission_id: this.state.mission.id,
          scenario_id: "custom",
        });
        const res = await fetch("/api/evaluate", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          this.state = await res.json();
          this.syncFormFromState();
          this.renderTelemetryChart();
        }
      } catch (e) {
        console.error("Failed to evaluate vitals:", e);
      }
    },

    async toggleStep(stepId) {
      try {
        const res = await fetch(`/api/protocol/${stepId}/toggle`, {
          method: "POST",
        });
        if (res.ok) {
          this.state = await res.json();
        }
      } catch (e) {
        console.error("Failed to toggle protocol step:", e);
      }
    },

    toggleSystem(sysId) {
      this.activeSystem = this.activeSystem === sysId ? "all" : sysId;
      if (this.activeSystem === "head") {
        this.chartMode = "neuro";
      } else if (this.activeSystem === "heart") {
        this.chartMode = "cardio";
      } else if (this.activeSystem === "skeleton") {
        this.chartMode = "bone";
      } else if (this.activeSystem === "blood") {
        this.chartMode = "radiation";
      }
      this.renderTelemetryChart();
    },

    activeSystemLabel() {
      const map = {
        all: "ALL 4 ANATOMICAL SYSTEMS",
        head: "1. HEAD & EYES (NEURO-OCULAR / SANS)",
        heart: "2. HEART & CHEST (CARDIO / RESPIRATORY)",
        skeleton: "3. SKELETON & MUSCLES (BONE / ARED)",
        blood: "4. BLOOD, IMMUNE & CSA CANS RADIATION",
      };
      return map[this.activeSystem] || "ALL SYSTEMS";
    },

    systemTier(sysId) {
      const list = this.state.body_systems || [];
      const found = list.find((s) => s.id === sysId);
      return found ? found.tier : "GREEN";
    },

    systemStrokeColor(sysId) {
      const tier = this.systemTier(sysId);
      if (tier === "RED") return "#f43f5e";
      if (tier === "AMBER") return "#fbbf24";
      return "#10b981";
    },

    systemFillColor(sysId) {
      const tier = this.systemTier(sysId);
      if (tier === "RED") return "rgba(244, 63, 94, 0.32)";
      if (tier === "AMBER") return "rgba(251, 191, 36, 0.28)";
      return "rgba(16, 185, 129, 0.22)";
    },

    filteredAlertCards() {
      const all = this.state.alert_cards || [];
      return all.filter((c) => {
        if (this.activeSystem !== "all" && c.body_system !== this.activeSystem) {
          return false;
        }
        if (this.showAlertsOnly && c.tier === "GREEN") {
          return false;
        }
        return true;
      });
    },

    setChartMode(mode) {
      this.chartMode = mode;
      this.renderTelemetryChart();
    },

    chartTitle() {
      const titles = {
        cardio: "30-Day Cardiovascular & Glycemic Telemetry (Systolic BP, Heart Rate & Blood Sugar)",
        radiation: "30-Day CSA CANS Neutron Dosimetry vs. Total Daily Exposure (mSv/day)",
        neuro: "30-Day Neuro-Ocular SANS Risk Index, Sleep Duration (hrs) & Mood Score",
        bone: "30-Day ARED Resistive Exercise Compliance (%) vs. Est. Bone Loss Rate (%/mo)",
      };
      return titles[this.chartMode] || titles.cardio;
    },

    runPitchBeat(idx) {
      this.activePitchBeat = idx;
      if (idx === 0) {
        this.earthModalOpen = false;
        this.switchMission("mars");
        document.getElementById("module-a")?.scrollIntoView({ behavior: "smooth" });
      } else if (idx === 1) {
        this.earthModalOpen = false;
        this.activeSystem = "all";
        document.getElementById("module-b")?.scrollIntoView({ behavior: "smooth" });
      } else if (idx === 2) {
        this.earthModalOpen = false;
        this.triggerScenario("mars_neutron_spike");
        document.getElementById("module-d")?.scrollIntoView({ behavior: "smooth" });
      } else if (idx === 3) {
        this.earthModalOpen = false;
        this.activeSystem = "blood";
        this.setChartMode("radiation");
        document.getElementById("module-e")?.scrollIntoView({ behavior: "smooth" });
      } else if (idx === 4) {
        this.earthModalOpen = true;
      }
    },

    renderTelemetryChart() {
      const canvas = document.getElementById("telemetry30DayChart");
      if (!canvas) return;
      const series = this.state.telemetry_30_days || [];
      if (!series.length) return;

      const labels = series.map((p) => p.label);
      let datasets = [];

      if (this.chartMode === "cardio") {
        datasets = [
          {
            label: "Systolic BP (mmHg)",
            data: series.map((p) => p.systolic_bp),
            borderColor: "#22d3ee",
            backgroundColor: "rgba(34, 211, 238, 0.12)",
            fill: true,
            tension: 0.35,
            borderWidth: 2.2,
          },
          {
            label: "Blood Sugar (mg/dL)",
            data: series.map((p) => p.blood_sugar),
            borderColor: "#f43f5e",
            backgroundColor: "rgba(244, 63, 94, 0.08)",
            fill: false,
            tension: 0.35,
            borderWidth: 2.2,
          },
          {
            label: "Heart Rate (bpm)",
            data: series.map((p) => p.heart_rate),
            borderColor: "#10b981",
            backgroundColor: "transparent",
            borderDash: [4, 4],
            tension: 0.35,
            borderWidth: 2,
          },
        ];
      } else if (this.chartMode === "radiation") {
        datasets = [
          {
            label: "Total Daily Radiation (mSv/d)",
            data: series.map((p) => p.total_rad_msv),
            borderColor: "#f43f5e",
            backgroundColor: "rgba(244, 63, 94, 0.15)",
            fill: true,
            tension: 0.3,
            borderWidth: 2.4,
          },
          {
            label: "CSA CANS Secondary Neutron Dose (mSv/d)",
            data: series.map((p) => p.neutron_rad_msv),
            borderColor: "#c084fc",
            backgroundColor: "rgba(192, 132, 252, 0.22)",
            fill: true,
            tension: 0.3,
            borderWidth: 2.4,
          },
          {
            label: "Cumulative 30-Day Dose (÷10 mSv scale)",
            data: series.map((p) => +(p.cumulative_rad_msv / 10).toFixed(2)),
            borderColor: "#38bdf8",
            borderDash: [4, 3],
            fill: false,
            tension: 0.3,
            borderWidth: 1.8,
          },
        ];
      } else if (this.chartMode === "neuro") {
        datasets = [
          {
            label: "SANS Risk Index (0-100)",
            data: series.map((p) => p.sans_risk_index),
            borderColor: "#fbbf24",
            backgroundColor: "rgba(251, 191, 36, 0.14)",
            fill: true,
            tension: 0.35,
            borderWidth: 2.2,
          },
          {
            label: "Sleep Hours (×10 scale)",
            data: series.map((p) => +(p.sleep_hours * 10).toFixed(1)),
            borderColor: "#22d3ee",
            fill: false,
            tension: 0.35,
            borderWidth: 2.2,
          },
          {
            label: "Mood Readiness (×10 scale)",
            data: series.map((p) => p.mood_score * 10),
            borderColor: "#10b981",
            borderDash: [4, 4],
            fill: false,
            tension: 0.35,
            borderWidth: 2,
          },
        ];
      } else {
        datasets = [
          {
            label: "ARED Exercise Compliance (%)",
            data: series.map((p) => p.ared_compliance),
            borderColor: "#10b981",
            backgroundColor: "rgba(16, 185, 129, 0.12)",
            fill: true,
            tension: 0.35,
            borderWidth: 2.2,
          },
          {
            label: "Est. Bone Loss Rate (%/mo × 50 scale)",
            data: series.map((p) => +(p.bone_loss_rate_pct * 50).toFixed(1)),
            borderColor: "#f43f5e",
            backgroundColor: "rgba(244, 63, 94, 0.12)",
            fill: true,
            tension: 0.35,
            borderWidth: 2.2,
          },
        ];
      }

      if (window.Chart) {
        if (telemetryChartInstance) {
          telemetryChartInstance.destroy();
        }
        telemetryChartInstance = new window.Chart(canvas, {
          type: "line",
          data: { labels, datasets },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: { mode: "index", intersect: false },
            plugins: {
              legend: {
                labels: {
                  color: "#cbd5e1",
                  font: { family: "JetBrains Mono, monospace", size: 11 },
                  boxWidth: 14,
                },
              },
              tooltip: {
                backgroundColor: "rgba(2, 6, 23, 0.95)",
                titleColor: "#38bdf8",
                bodyColor: "#f8fafc",
                borderColor: "#334155",
                borderWidth: 1,
              },
            },
            scales: {
              x: {
                ticks: {
                  color: "#64748b",
                  font: { family: "JetBrains Mono, monospace", size: 10 },
                  maxTicksLimit: 10,
                },
                grid: { color: "rgba(30, 41, 59, 0.5)" },
              },
              y: {
                ticks: {
                  color: "#94a3b8",
                  font: { family: "JetBrains Mono, monospace", size: 10 },
                },
                grid: { color: "rgba(30, 41, 59, 0.65)" },
              },
            },
          },
        });
      } else {
        // Offline Native 2D Canvas Fallback Renderer
        this.drawOfflineFallbackChart(canvas, labels, datasets);
      }
    },

    drawOfflineFallbackChart(canvas, labels, datasets) {
      const ctx = canvas.getContext("2d");
      const rect = canvas.parentElement.getBoundingClientRect();
      canvas.width = rect.width || 640;
      canvas.height = rect.height || 260;
      const w = canvas.width;
      const h = canvas.height;
      ctx.clearRect(0, 0, w, h);

      const padL = 42,
        padR = 16,
        padT = 28,
        padB = 28;
      const plotW = w - padL - padR;
      const plotH = h - padT - padB;

      let allVals = [];
      datasets.forEach((d) => (allVals = allVals.concat(d.data)));
      const minV = Math.min(0, ...allVals);
      const maxV = Math.max(10, ...allVals) * 1.1;

      ctx.strokeStyle = "#1e293b";
      ctx.lineWidth = 1;
      for (let i = 0; i <= 4; i++) {
        const y = padT + (plotH * i) / 4;
        ctx.beginPath();
        ctx.moveTo(padL, y);
        ctx.lineTo(w - padR, y);
        ctx.stroke();
      }

      datasets.forEach((ds) => {
        ctx.strokeStyle = ds.borderColor;
        ctx.lineWidth = 2;
        ctx.beginPath();
        ds.data.forEach((val, idx) => {
          const x = padL + (plotW * idx) / (ds.data.length - 1);
          const y = padT + plotH - ((val - minV) / (maxV - minV)) * plotH;
          if (idx === 0) ctx.moveTo(x, y);
          else ctx.lineTo(x, y);
        });
        ctx.stroke();
      });
    },
  };
};
