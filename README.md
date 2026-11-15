# OrbitoDX

Autonomous crew health triage and countermeasure software for deep-space missions.

## Background

In Low Earth Orbit, astronauts have near-instant communication with Flight Surgeons on the ground and can return to Earth in hours if a medical issue escalates. On lunar and Mars missions, that safety net goes away. A Mars transit introduces a one-way communication delay of up to 20 minutes (40 minutes round-trip), while NASA's Five Hazards of Human Spaceflight interact and compound one another. Because not every crew member is a medical doctor, the crew needs onboard software that can flag early physiological or behavioral drift, explain what is happening in plain English, and walk them through the right countermeasure before a minor shift becomes an emergency.

OrbitoDX is a clinical decision-support dashboard built around data and risk models from the NASA Human Research Program (HRP) and the Canadian Space Agency (CSA). It takes crew vital signs and habitat telemetry, maps them to the five spaceflight hazards and neutron radiation exposure data (based on CSA's Canadian Active Neutron Spectrometer, CANS), and outputs color-coded triage protocols.

## What's Inside

- **Mission Context Switcher:** Switch between ISS (Low Earth Orbit), Artemis Lunar Base (1/6g + lunar regolith dust), and Mars Transit (0g + high GCR/neutron flux + 14m 22s one-way communication delay). Switching missions updates baseline radiation, gravity unloading, and communication latency assumptions.
- **Interactive Body Map:** An SVG anatomical view broken into four systems:
  1. *Head & Eyes* (SANS / neuro-ocular pressure, sleep hours, mood and stress)
  2. *Heart & Chest* (blood pressure, heart rate, orthostatic tolerance, cabin CO2 and lunar dust)
  3. *Skeleton & Muscles* (monthly bone mineral density loss estimate and ARED resistive exercise compliance)
  4. *Blood, Immune & Radiation* (blood glucose, core temperature, latent viral reactivation risk, and cumulative neutron dose via CSA CANS)
- **Vitals Check-In & Anomaly Simulator:** A manual form for logging crew vitals alongside four 1-click presets (*Nominal Baseline*, *Mars Solar Neutron Spike*, *Microgravity Bone & Cardio Drift*, and *Circadian Sleep Deficit*) for testing and live demos.
- **3-Tier Triage & Plain-Language Protocols:** Categorizes indicators into `GREEN` (nominal), `AMBER` (self-intervention / watch), and `RED` (escalate to Crew Medical Officer / store-and-forward packet to Flight Surgeon). Each card breaks down why the check matters, which spaceflight hazards are driving the change, and a step-by-step checklist of countermeasures.
- **30-Day Trend Charts & Readiness Score:** Time-series views across Sol 1 to Sol 30 and a daily readiness score that updates as countermeasures are completed.

## Tech Stack

We built this to compile down to a single static binary with no external database or runtime dependencies, keeping memory usage low enough to run on constrained hardware or offline edge devices (such as remote clinics on Earth).

- **Backend:** Go (`gofiber/fiber/v3`)
- **Templates:** `a-h/templ` (compiled to Go code)
- **Frontend:** Tailwind CSS, Alpine.js, Chart.js
- **Packaging:** Go `embed.FS` bundles the compiled CSS and JS directly into the binary

## Project Layout

```text
cmd/orbitodx/        Server entrypoint and CLI flags
internal/domain/     Data structures for missions, hazards, vitals, and triage cards
internal/engine/     Triage rules, CSA CANS neutron calculations, and 30-day trend generator
internal/handlers/   HTTP routes and JSON API handlers
internal/store/      In-memory session state
views/               templ templates (layout and Modules A-E)
web/                 Embedded static assets (compiled CSS and cockpit JS)
deploy/              systemd unit file and reverse proxy config for VPS hosting
```

## Running Locally

If you have Go 1.25+ installed, the generated `templ` files and compiled CSS are already committed in the repo, so you can build and run directly:

```bash
go build -o bin/orbitodx ./cmd/orbitodx
./bin/orbitodx --port 8080
```

Then open `http://localhost:8080` in your browser.

### Running Tests

```bash
go test -v ./internal/...
```

### Rebuilding Templates or CSS

If you edit `.templ` files or Tailwind styles:

```bash
make generate
make css
make build-vps
```

## VPS Deployment

Build the static Linux binary:

```bash
make build-vps
```

Copy `bin/orbitodx` and `deploy/orbitodx.service` to your server:

```bash
scp bin/orbitodx root@your-server-ip:/opt/orbitodx/orbitodx
scp deploy/orbitodx.service root@your-server-ip:/etc/systemd/system/orbitodx.service
ssh root@your-server-ip "chmod +x /opt/orbitodx/orbitodx && systemctl daemon-reload && systemctl enable --now orbitodx"
```

## References

- NASA Human Research Program (HRP) — *5 Hazards of Human Spaceflight* (Space Radiation, Isolation & Confinement, Distance from Earth, Gravity Fields, Hostile/Closed Environments)
- Canadian Space Agency (CSA) — *Canadian Active Neutron Spectrometer (CANS)* and secondary neutron dosimetry on the ISS (~30% of crew dose equivalent)
