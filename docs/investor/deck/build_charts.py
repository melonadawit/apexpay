#!/usr/bin/env python3
"""Generate all charts and graphics for the detailed ApexPay investor deck."""
import os
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch, Wedge
import numpy as np

OUT = os.path.join(os.path.dirname(__file__), "charts")
os.makedirs(OUT, exist_ok=True)

GREEN = "#0B6E4F"
GREEN2 = "#10A37A"
GOLD = "#EAB308"
AMBER = "#F59E0B"
RED = "#EF4444"
BLUE = "#3B82F6"
NAVY = "#1a202c"
GREY = "#4a5568"
LIGHT = "#f7fafc"
GREY2 = "#e2e8f0"

plt.rcParams.update({
    "font.family": "DejaVu Sans",
    "axes.edgecolor": "#cbd5e1",
    "axes.labelcolor": GREY,
    "xtick.color": GREY,
    "ytick.color": GREY,
    "text.color": NAVY,
    "axes.titlecolor": NAVY,
    "figure.facecolor": "white",
    "axes.facecolor": "white",
})


def save(fig, name):
    fig.savefig(os.path.join(OUT, name), bbox_inches="tight", dpi=150, facecolor="white")
    plt.close(fig)
    print("wrote", name)


# ---------------------------------------------------------------- market size
def market_size():
    years = ["2023", "2024", "2025", "2026E", "2027E", "2028E"]
    value = [6.5, 9.0, 12.5, 17.0, 23.0, 30.0]  # USD bn digital payments Ethiopia/CET
    fig, ax = plt.subplots(figsize=(7, 4.2))
    bars = ax.bar(years, value, color=[GREEN, GREEN, GREEN, GREEN2, GREEN2, GOLD], width=0.6)
    ax.set_title("Digital Payments Value, Ethiopia & East Africa (USD bn)")
    ax.set_ylabel("USD billions")
    for b, v in zip(bars, value):
        ax.text(b.get_x() + b.get_width()/2, v + 0.6, f"${v:g}", ha="center", fontsize=9, fontweight="bold", color=NAVY)
    ax.grid(axis="y", color=GREY2, linewidth=0.6)
    ax.spines[["top", "right"]].set_visible(False)
    ax.set_ylim(0, 33)
    save(fig, "market_size.png")


# ---------------------------------------------------------------- TAM SAM SOM
def tam_sam_som():
    fig, ax = plt.subplots(figsize=(7, 4.2))
    labels = ["TAM\nEthiopia digital finance\n$17B by 2026", "SAM\nFormal MSME + payroll\n$3.4B", "SOM\nApexPay 5-yr serviceable\n$180M"]
    vals = [17.0, 3.4, 0.18]
    colors = [GREY2, GREEN, GOLD]
    bars = ax.barh([2, 1, 0], vals, color=colors, height=0.6)
    ax.set_yticks([2, 1, 0]); ax.set_yticklabels(labels, fontsize=8.5)
    for b, v in zip(bars, vals):
        ax.text(v + 0.2, b.get_y()+b.get_height()/2, f"${v:,.2f}B" if v>=1 else f"${v*1000:,.0f}M",
                va="center", fontsize=9, fontweight="bold", color=NAVY)
    ax.set_title("TAM / SAM / SOM")
    ax.set_xscale("log"); ax.set_xlim(0.05, 60)
    ax.spines[["top", "right"]].set_visible(False)
    ax.grid(axis="x", color=GREY2, linewidth=0.5)
    save(fig, "tam_sam_som.png")


# ---------------------------------------------------------------- TPV growth
def tpv_growth():
    years = ["2026", "2027", "2028", "2029", "2030"]
    tpv = [30, 95, 210, 390, 650]  # ETB bn
    fig, ax = plt.subplots(figsize=(7, 4.2))
    x = np.arange(len(years))
    ax.plot(x, tpv, marker="o", color=GREEN, linewidth=2.5, markersize=7)
    ax.fill_between(x, tpv, color=GREEN, alpha=0.12)
    ax.set_xticks(x); ax.set_xticklabels(years)
    ax.set_ylabel("TPV (ETB billions)")
    ax.set_title("Projected Total Payment Volume (TPV)")
    for i, v in enumerate(tpv):
        ax.annotate(f"ETB {v}B", (i, v), textcoords="offset points", xytext=(0, 10), ha="center", fontsize=9, fontweight="bold")
    ax.grid(axis="y", color=GREY2, linewidth=0.6)
    ax.spines[["top", "right"]].set_visible(False)
    ax.set_ylim(0, 720)
    save(fig, "tpv_growth.png")


# ---------------------------------------------------------------- revenue projection
def revenue_projection():
    years = ["2026", "2027", "2028", "2029", "2030"]
    rev = [1.8, 6.2, 14.5, 27.0, 44.0]  # USD M
    ebitda = [0.1, 0.8, 3.0, 7.5, 14.0]  # USD M
    fig, ax = plt.subplots(figsize=(7, 4.2))
    x = np.arange(len(years)); w=0.38
    ax.bar(x-w/2, rev, w, color=GREEN, label="Revenue")
    ax.bar(x+w/2, ebitda, w, color=GOLD, label="EBITDA")
    for i,(a,b) in enumerate(zip(rev,ebitda)):
        ax.text(i-w/2, a+1, f"{a:.1f}", ha="center", fontsize=8, fontweight="bold")
        ax.text(i+w/2, b+1, f"{b:.1f}", ha="center", fontsize=8, fontweight="bold")
    ax.set_xticks(x); ax.set_xticklabels(years)
    ax.set_ylabel("USD millions")
    ax.set_title("Revenue & EBITDA Projection (USD M)")
    ax.legend(frameon=False, fontsize=9)
    ax.grid(axis="y", color=GREY2, linewidth=0.6)
    ax.spines[["top","right"]].set_visible(False)
    ax.set_ylim(0, 50)
    save(fig, "revenue_projection.png")


# ---------------------------------------------------------------- revenue mix
def revenue_mix():
    labels = ["Payment fees", "Payroll SaaS", "Embedded finance", "FX & banking", "Data / AI"]
    sizes = [34, 26, 18, 14, 8]
    colors = [GREEN, GREEN2, GOLD, BLUE, GREY]
    fig, ax = plt.subplots(figsize=(6.4, 4.2))
    wedges, texts, autotxts = ax.pie(sizes, labels=labels, colors=colors, autopct="%d%%", startangle=140,
                                     pctdistance=0.78, wedgeprops=dict(width=0.42, edgecolor="white"))
    for t in autotxts: t.set_fontsize(9); t.set_color("white"); t.set_fontweight("bold")
    ax.set_title("Revenue Mix (mature state)")
    ax.legend(wedges, labels, loc="lower center", bbox_to_anchor=(0.5,-0.25), ncol=2, fontsize=8, frameon=False)
    save(fig, "revenue_mix.png")


# ---------------------------------------------------------------- unit economics
def unit_economics():
    metrics = ["Take rate\n(avg)", "Gross margin\n(payments)", "Gross margin\n(platform)", "Contribution\nmargin", "Net revenue\nretention"]
    vals = [2.9, 55, 72, 38, 112]
    colors = [GOLD, GREEN, GREEN2, BLUE, GREEN]
    fig, ax = plt.subplots(figsize=(7, 4.2))
    bars = ax.bar(metrics, vals, color=colors, width=0.55)
    for b, v in zip(bars, vals):
        ax.text(b.get_x()+b.get_width()/2, v+2, (f"{v:.1f}%" if v<100 else f"{v:.0f}%"), ha="center", fontsize=9, fontweight="bold")
    ax.set_ylim(0, 130)
    ax.set_ylabel("Percent")
    ax.set_title("Unit Economics & Retention")
    ax.grid(axis="y", color=GREY2, linewidth=0.6)
    ax.spines[["top","right"]].set_visible(False)
    save(fig, "unit_economics.png")


# ---------------------------------------------------------------- cost structure
def cost_structure():
    cats = ["Engineering", "Sales & Marketing", "Ops & Compliance", "G&A", "Infrastructure", "R&D / AI"]
    vals = [38, 22, 14, 10, 9, 7]
    colors = [GREEN, GOLD, BLUE, GREEN2, AMBER, GREY]
    fig, ax = plt.subplots(figsize=(7, 4.2))
    wedges, texts, autotxts = ax.pie(vals, labels=cats, colors=colors, autopct="%d%%", startangle=120,
                                     wedgeprops=dict(edgecolor="white"))
    for t in autotxts: t.set_fontsize(8); t.set_color("white"); t.set_fontweight("bold")
    ax.set_title("Cost Structure (% of spend)")
    save(fig, "cost_structure.png")


# ---------------------------------------------------------------- KPI ratios dashboard
def kpi_ratios():
    fig, axes = plt.subplots(2, 3, figsize=(9, 4.8))
    data = [
        ("Gross Margin", 72, GREEN),
        ("Net Margin", 34, GREEN2),
        ("CAC Payback (months)", 9, GOLD),
        ("LTV:CAC", 5.2, BLUE),
        ("Monthly Burn (USD k)", 240, RED),
        ("Runway (months)", 22, AMBER),
    ]
    for ax, (label, val, color) in zip(axes.flat, data):
        # gauge style
        ax.add_patch(Wedge((0,0), 0.9, 0, 180, facecolor=GREY2, width=0.18, edgecolor="white"))
        ax.add_patch(Wedge((0,0), 0.9, 0, min(180, val*180/100) if label not in("Monthly Burn","Runway","CAC Payback") else 180, facecolor=color, width=0.18, edgecolor="white"))
        ax.text(0, -0.35, label, ha="center", fontsize=7, color=GREY)
        ax.text(0, 0.05, f"{val:,}%" if label not in("LTV:CAC","Monthly Burn","Runway","CAC Payback") else (f"{val:.1f}x" if label=="LTV:CAC" else f"${val:,}k" if label=="Monthly Burn" else f"{val} mo" if label=="Runway" else f"{val} mo"),
                ha="center", fontsize=13, fontweight="bold", color=color)
        ax.set_xlim(-1.1,1.1); ax.set_ylim(-0.7,1.1); ax.axis("off")
    fig.suptitle("Key Ratios & KPIs", fontsize=12, fontweight="bold")
    save(fig, "kpi_ratios.png")


# ---------------------------------------------------------------- payment mix
def payment_mix():
    labels = ["Telebirr", "CBE Birr", "Bank / IPS", "Card", "EthSwitch QR"]
    vals = [38, 27, 18, 10, 7]
    colors = [GREEN, GREEN2, GOLD, BLUE, AMBER]
    fig, ax = plt.subplots(figsize=(6.4, 4.0))
    wedges, texts, autotxts = ax.pie(vals, labels=labels, colors=colors, autopct="%d%%", startangle=90,
                                     wedgeprops=dict(edgecolor="white"))
    for t in autotxts: t.set_fontsize(9); t.set_color("white"); t.set_fontweight("bold")
    ax.set_title("Payment Method Mix (Ethiopia)")
    save(fig, "payment_mix.png")


# ---------------------------------------------------------------- module coverage
def module_coverage():
    cats = ["Payments", "General Ledger", "Payroll", "Procurement/AP", "Inventory", "Tax", "FX", "Budget/FP&A", "Embedded Finance", "AI Assistant"]
    done = [100, 100, 100, 100, 100, 100, 100, 100, 100, 100]
    fig, ax = plt.subplots(figsize=(7, 4.2))
    y = np.arange(len(cats))
    ax.barh(y, done, color=GREEN, height=0.6)
    for i, v in enumerate(done):
        ax.text(v+1, i, f"{v}%", va="center", fontsize=8, fontweight="bold", color=NAVY)
    ax.set_yticks(y); ax.set_yticklabels(cats, fontsize=8)
    ax.set_xlim(0, 110)
    ax.set_title("Module Coverage — 43 business modules on one ledger")
    ax.spines[["top","right"]].set_visible(False)
    ax.grid(axis="x", color=GREY2, linewidth=0.5)
    save(fig, "module_coverage.png")


# ---------------------------------------------------------------- competitive matrix
def competitive_matrix():
    # x = integration depth, y = Ethiopia-native
    points = {
        "ApexPay": (0.92, 0.94, GREEN, 300),
        "Global gateways": (0.45, 0.25, BLUE, 120),
        "Legacy accounting": (0.30, 0.60, GREY, 120),
        "Point payroll tools": (0.28, 0.75, AMBER, 110),
        "Fintech challengers": (0.60, 0.55, RED, 110),
    }
    fig, ax = plt.subplots(figsize=(7, 4.4))
    for name, (x, y, c, s) in points.items():
        ax.scatter(x, y, s=s, color=c, edgecolor="white", linewidth=1.5, zorder=3)
        ax.annotate(name, (x, y), textcoords="offset points", xytext=(0, 8), ha="center", fontsize=9, fontweight="bold", color=NAVY)
    ax.set_xlim(0,1); ax.set_ylim(0,1)
    ax.set_xlabel("Integration depth / full-stack")
    ax.set_ylabel("Ethiopia-native compliance")
    ax.set_title("Competitive Positioning")
    ax.axvline(0.5, color=GREY2, linewidth=0.8)
    ax.axhline(0.5, color=GREY2, linewidth=0.8)
    ax.spines[["top","right"]].set_visible(False)
    save(fig, "competitive_matrix.png")


# ---------------------------------------------------------------- user flow helpers
def flow_box(ax, x, y, w, h, text, color=GREEN, fc=None, fs=7.5):
    fc = fc or "white"
    box = FancyBboxPatch((x, y), w, h, boxstyle="round,pad=0.004", linewidth=1.2,
                         edgecolor=color, facecolor=fc)
    ax.add_patch(box)
    ax.text(x+w/2, y+h/2, text, ha="center", va="center", fontsize=fs, color=NAVY, wrap=True)


def flow_arrow(ax, x1, y1, x2, y2, color=GREEN):
    ax.add_patch(FancyArrowPatch((x1, y1), (x2, y2), arrowstyle="-|>", mutation_scale=12,
                                 color=color, linewidth=1.6))


def user_flow_onboarding():
    fig, ax = plt.subplots(figsize=(8, 4.2)); ax.set_xlim(0,10); ax.set_ylim(0,6); ax.axis("off")
    steps = ["Merchant signs up", "KYC / Fayda ID\n+ business docs", "Risk scoring &\ncompliance review", "Onboarding\napproved", "Activation\n(test + live keys)", "First payment\non rails"]
    xs = [0.2, 2.0, 3.8, 5.6, 7.4, 9.0]
    for i,(x, s) in enumerate(zip(xs, steps)):
        flow_box(ax, x, 2.2, 1.1, 1.6, s, color=GREEN, fc=("white" if i%2==0 else LIGHT))
        if i < len(xs)-1:
            flow_arrow(ax, x+1.1, 3.0, xs[i+1], 3.0)
    ax.set_title("Merchant Onboarding Flow", fontsize=11, fontweight="bold")
    save(fig, "flow_onboarding.png")


def user_flow_payment():
    fig, ax = plt.subplots(figsize=(8, 4.4)); ax.set_xlim(0,10); ax.set_ylim(0,6); ax.axis("off")
    # top row: initiate
    flow_box(ax, 0.2, 4.4, 1.4, 1.0, "Initialize\n(tx_ref)", GREEN)
    flow_box(ax, 2.0, 4.4, 1.4, 1.0, "Routing\n(best rail)", GREEN2)
    flow_box(ax, 3.8, 4.4, 1.4, 1.0, "2FA above\nthreshold", GOLD)
    flow_arrow(ax,1.6,4.9,2.0,4.9); flow_arrow(ax,3.4,4.9,3.8,4.9)
    # connector
    flow_box(ax, 5.8, 4.4, 1.4, 1.0, "Connector\n(Telebirr…)", BLUE)
    flow_arrow(ax,5.2,4.9,5.8,4.9)
    # verify
    flow_box(ax, 7.8, 4.4, 1.4, 1.0, "Verify", GREEN)
    flow_arrow(ax,7.2,4.9,7.8,4.9)
    # ledger row
    flow_box(ax, 3.0, 2.0, 2.2, 1.0, "Ledger journal\n(debit/credit)", GREEN, LIGHT)
    flow_box(ax, 5.8, 2.0, 1.8, 1.0, "Outbox event", AMBER, LIGHT)
    flow_box(ax, 8.0, 2.0, 1.6, 1.0, "Webhook\n+ notify", BLUE, LIGHT)
    flow_arrow(ax,4.6,4.4,4.6,3.0); flow_arrow(ax,4.6,3.0,5.8,2.5); flow_arrow(ax,7.6,2.5,8.0,2.5)
    ax.set_title("Payment Lifecycle (one transaction, one ledger)", fontsize=11, fontweight="bold")
    save(fig, "flow_payment.png")


def user_flow_payroll():
    fig, ax = plt.subplots(figsize=(8, 4.4)); ax.set_xlim(0,10); ax.set_ylim(0,6); ax.axis("off")
    steps = ["Create run", "Attendance /\nOT inputs", "Calculate\n(ET tax, proration)", "Approve\n(maker-checker)", "Disburse\n+ bank file", "Ledger + reports\n(pension, ERCA)"]
    xs = [0.2, 2.0, 3.8, 5.6, 7.4, 9.0]
    for i,(x, s) in enumerate(zip(xs, steps)):
        flow_box(ax, x, 2.2, 1.1, 1.7, s, color=GREEN, fc=("white" if i%2==0 else LIGHT))
        if i < len(xs)-1:
            flow_arrow(ax, x+1.1, 3.05, xs[i+1], 3.05)
    ax.set_title("Payroll Run Flow (compliance: ET tax + pension 7%/11%)", fontsize=11, fontweight="bold")
    save(fig, "flow_payroll.png")


# ---------------------------------------------------------------- architecture
def architecture():
    fig, ax = plt.subplots(figsize=(8, 4.6)); ax.set_xlim(0,10); ax.set_ylim(0,7); ax.axis("off")
    flow_box(ax, 0.4, 6.2, 2.4, 0.6, "Merchant Web", BLUE)
    flow_box(ax, 3.0, 6.2, 2.4, 0.6, "Admin / Checkout", BLUE)
    flow_box(ax, 5.6, 6.2, 2.4, 0.6, "Mobile (Flutter)", BLUE)
    flow_box(ax, 8.2, 6.2, 1.6, 0.6, "Portals", BLUE)
    # API core
    flow_box(ax, 1.2, 4.2, 7.6, 1.2, "Go API core (chi)  ·  payments / GL / payroll / procurement / tax / FX / AI",
             GREEN, fc=LIGHT, fs=9)
    for x in [1.6, 4.2, 6.8, 9.0]:
        pass
    flow_arrow(ax,1.6,6.2,2.2,5.4); flow_arrow(ax,4.2,6.2,4.2,5.4); flow_arrow(ax,6.8,6.2,6.2,5.4); flow_arrow(ax,8.9,6.2,7.4,5.4)
    # stores
    flow_box(ax, 0.4, 2.0, 2.2, 0.9, "PostgreSQL 17\n(source of truth)", NAVY, LIGHT)
    flow_box(ax, 3.0, 2.0, 2.0, 0.9, "Redis\n(cache)", NAVY, LIGHT)
    flow_box(ax, 5.4, 2.0, 1.8, 0.9, "MinIO\n(docs)", NAVY, LIGHT)
    flow_box(ax, 7.6, 2.0, 2.0, 0.9, "Workers\n(outbox, jobs)", NAVY, LIGHT)
    flow_arrow(ax,2.6,4.2,1.8,2.9); flow_arrow(ax,4.4,4.2,4.0,2.9); flow_arrow(ax,5.6,4.2,6.2,2.9); flow_arrow(ax,6.4,4.2,8.5,2.9)
    ax.set_title("System Architecture — one money core, thin clients", fontsize=11, fontweight="bold")
    save(fig, "architecture.png")


# ---------------------------------------------------------------- roadmap
def roadmap_timeline():
    fig, ax = plt.subplots(figsize=(8, 4.2)); ax.set_xlim(0,10); ax.set_ylim(0,5); ax.axis("off")
    phases = [
        ("Core Payments", "2025", GREEN),
        ("GL + Payroll", "2025", GREEN2),
        ("AI Assistant", "2026", GOLD),
        ("Procurement/Tax/FX", "2026", BLUE),
        ("Embedded Finance", "2026", AMBER),
        ("Scale + Banking", "2027", RED),
    ]
    for i,(label, yr, color) in enumerate(phases):
        x = 0.4 + i*1.6
        flow_box(ax, x, 2.2, 1.3, 1.3, f"{label}\n{yr}", color, fc="white", fs=7)
        if i < len(phases)-1:
            flow_arrow(ax, x+1.3, 2.85, x+1.6, 2.85)
    ax.text(5, 4.4, "Delivery Roadmap", ha="center", fontsize=12, fontweight="bold", color=NAVY)
    save(fig, "roadmap_timeline.png")


# ---------------------------------------------------------------- financial ratios table (payroll economics)
def payroll_economics():
    cats = ["Median salary (ETB/mo)", "Employer pension 11%", "Employee pension 7%", "Payroll tax (avg)", "Formal employees (M)"]
    # two-panel: cost of compliance + formalization opportunity
    fig, axes = plt.subplots(1, 2, figsize=(8, 4.2))
    ax = axes[0]
    v = [18, 11, 7, 15]
    lbls = ["Pension\n(employer)", "Pension\n(employee)", "Income tax\n(avg)", "Total burden %"]
    ax.bar(lbls, v, color=[GREEN, GREEN2, GOLD, RED], width=0.55)
    for b,val in zip(ax.patches, v): ax.text(b.get_x()+b.get_width()/2, val+0.4, f"{val}%", ha="center", fontsize=8, fontweight="bold")
    ax.set_title("Ethiopia payroll burden", fontsize=10); ax.set_ylim(0,20)
    ax.spines[["top","right"]].set_visible(False); ax.grid(axis="y", color=GREY2, linewidth=0.5)
    ax = axes[1]
    years = ["2024","2025","2026E","2027E"]
    formal = [2.1, 2.4, 2.8, 3.3]
    ax.plot(years, formal, marker="o", color=GREEN, linewidth=2)
    ax.fill_between(years, formal, color=GREEN, alpha=0.12)
    ax.set_title("Formal employment (millions)", fontsize=10)
    ax.grid(axis="y", color=GREY2, linewidth=0.5); ax.spines[["top","right"]].set_visible(False)
    save(fig, "payroll_economics.png")


# ---------------------------------------------------------------- financial summary bar
def financial_summary():
    years = ["2026","2027","2028","2029","2030"]
    rev = [1.8,6.2,14.5,27.0,44.0]
    cash = [-1.2, 0.5, 4.5, 10.0, 18.5]
    fig, ax = plt.subplots(figsize=(7,4.0))
    x = np.arange(len(years)); w=0.38
    ax.bar(x-w/2, rev, w, color=GREEN, label="Revenue")
    ax.bar(x+w/2, cash, w, color=GOLD, label="Operating cash flow")
    ax.axhline(0, color=GREY, linewidth=0.8)
    ax.set_xticks(x); ax.set_xticklabels(years)
    ax.set_ylabel("USD millions")
    ax.set_title("Financial Summary (USD M)")
    ax.legend(frameon=False, fontsize=8)
    ax.grid(axis="y", color=GREY2, linewidth=0.5)
    ax.spines[["top","right"]].set_visible(False)
    ax.set_ylim(-4, 48)
    save(fig, "financial_summary.png")


def main():
    market_size()
    tam_sam_som()
    tpv_growth()
    revenue_projection()
    revenue_mix()
    unit_economics()
    cost_structure()
    kpi_ratios()
    payment_mix()
    module_coverage()
    competitive_matrix()
    user_flow_onboarding()
    user_flow_payment()
    user_flow_payroll()
    architecture()
    roadmap_timeline()
    payroll_economics()
    financial_summary()
    print("ALL CHARTS DONE")


if __name__ == "__main__":
    main()
