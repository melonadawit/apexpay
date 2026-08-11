#!/usr/bin/env python3
"""Regenerate the ApexPay investor roadmap chart."""
import os
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.patches import FancyBboxPatch

OUT = os.path.join(os.path.dirname(__file__), "charts", "roadmap.png")
GREEN = "#0B6E4F"
GREEN2 = "#10B981"
GREY = "#4a5568"

milestones = [
    ("Core Payments", "initialise/2FA/verify + ledger", True),
    ("Real-time GL", "double-entry, fiscal close", True),
    ("Workforce OS", "payroll runs, leave, claims", True),
    ("AI Assistant", "RAG + tool-calling, i18n", True),
    ("Procurement/AP", "vendors, POs, AP aging", True),
    ("i18n & UI/UX", "EN + Amharic, 4 frontends", True),
    ("Budget & Portals", "budget, vendor/customer portal", True),
    ("Security Hardening", "10/10 CI gates green", True),
    ("Live Data Verified", "dashboard on real seeded API", True),
]

fig, ax = plt.subplots(figsize=(14.7, 5.8), dpi=100)
ax.set_xlim(0, 1)
ax.set_ylim(0, 1)
ax.axis("off")

n = len(milestones)
box_w = 0.94 / n
gap = 0.015
start = 0.03
y = 0.55

for i, (title, desc, done) in enumerate(milestones):
    x = start + i * (box_w + gap)
    color = GREEN if done else GREY
    box = FancyBboxPatch((x, y), box_w, 0.30, boxstyle="round,pad=0.004",
                         linewidth=0, facecolor=color, alpha=0.92)
    ax.add_patch(box)
    ax.text(x + box_w / 2, y + 0.20, title, ha="center", va="center",
            fontsize=8.5, color="white", fontweight="bold")
    ax.text(x + box_w / 2, y + 0.09, desc, ha="center", va="center",
            fontsize=6, color="white", alpha=0.92, wrap=True)

ax.text(0.5, 0.90, "ApexPay Delivery Roadmap", ha="center", fontsize=16,
        fontweight="bold", color=GREY)
ax.text(0.5, 0.83, "Every increment shipped working software with build / vet / test / smoke green — the final box is live-data verification.",
        ha="center", fontsize=8, color=GREY)

plt.tight_layout()
plt.savefig(OUT, bbox_inches="tight", facecolor="white")
print("Wrote", OUT)
