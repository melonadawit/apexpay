#!/usr/bin/env python3
"""Regenerate the ApexPay investor deck PDF.

Run from repo root:  python3 docs/investor/build_deck.py
Requires: reportlab, matplotlib (for chart PNGs). Reuses the chart PNGs under
docs/investor/charts/.
"""
import os
from reportlab.lib.pagesizes import A4, landscape
from reportlab.lib.units import mm
from reportlab.lib.colors import HexColor
from reportlab.pdfgen import canvas as pdfcanvas

GREEN = HexColor("#0B6E4F")
GREEN2 = HexColor("#10B981")
DARK = HexColor("#1a202c")
GREY = HexColor("#4a5568")
LIGHT = HexColor("#f7fafc")

OUT = os.path.join(os.path.dirname(__file__), "ApexPay_Investor_Deck.pdf")
CHARTS = os.path.join(os.path.dirname(__file__), "charts")

PAGE_W, PAGE_H = landscape(A4)


def header(c, title, sub):
    c.setFillColor(DARK)
    c.setFont("Helvetica-Bold", 20)
    c.drawString(18 * mm, PAGE_H - 18 * mm, title)
    c.setFont("Helvetica", 9.5)
    c.setFillColor(GREY)
    c.drawString(18 * mm, PAGE_H - 25 * mm, sub)
    c.setStrokeColor(GREEN)
    c.setLineWidth(1.2)
    c.line(18 * mm, PAGE_H - 29 * mm, PAGE_W - 18 * mm, PAGE_H - 29 * mm)


def bullets(c, x, y, items, size=10.5, gap=6.2, bold_prefix=False):
    yy = y
    c.setFillColor(DARK)
    for it in items:
        if bold_prefix:
            pfx, rest = it.split("|", 1)
            c.setFont("Helvetica-Bold", size)
            c.drawString(x, yy, pfx)
            c.setFont("Helvetica", size)
            x2 = x + c.stringWidth(pfx, "Helvetica-Bold", size) + 2
            c.drawString(x2, yy, rest)
        else:
            c.setFont("Helvetica", size)
            c.drawString(x, yy, "• " + it)
        yy -= size * 0.42 + gap
    return yy


def footer(c, page):
    c.setFillColor(GREY)
    c.setFont("Helvetica", 8)
    c.drawString(18 * mm, 8 * mm, "ApexPay — Investor Overview 2026")
    c.drawRightString(PAGE_W - 18 * mm, 8 * mm, f"{page}")

# Keep chart references here so we can place them if present.
CH = {k: os.path.join(CHARTS, k) for k in
      ["architecture.png", "modules.png", "revenue.png", "roadmap.png", "security.png"]}


def build():
    c = pdfcanvas.Canvas(OUT, pagesize=(PAGE_W, PAGE_H))

    # ---------- Page 1: Cover ----------
    c.setFillColor(DARK)
    c.rect(0, 0, PAGE_W, PAGE_H, stroke=0, fill=1)
    c.setFillColor(GREEN)
    c.rect(0, PAGE_H - 14 * mm, PAGE_W, 4 * mm, stroke=0, fill=1)
    c.setFillColor(HexColor("#ffffff"))
    c.setFont("Helvetica-Bold", 40)
    c.drawString(18 * mm, PAGE_H - 70 * mm, "ApexPay")
    c.setFont("Helvetica", 15)
    c.drawString(18 * mm, PAGE_H - 82 * mm,
                 "The all-in-one Financial Operating System for Ethiopia")
    c.setFont("Helvetica", 11)
    c.drawString(18 * mm, PAGE_H - 92 * mm,
                 "AI-native payments  ·  Real-time General Ledger  ·  Workforce & Payroll  ·  Procurement, Tax, FX & Budgeting")
    c.setFillColor(HexColor("#a0aec0"))
    c.setFont("Helvetica", 10)
    c.drawString(18 * mm, PAGE_H - 102 * mm, "Investor Overview — 2026")
    c.showPage()

    # ---------- Page 2: Executive Summary ----------
    header(c, "Executive Summary",
           "A full-stack, AI-native financial OS built for Ethiopia — one ledger, many modules.")
    y = PAGE_H - 40 * mm
    c.setFillColor(DARK)
    c.setFont("Helvetica", 11)
    p = ("ApexPay is a full-stack, AI-native financial operating system purpose-built for Ethiopia. "
         "It is not a single payment gateway but a composable platform that unifies payments, an "
         "immutable double-entry General Ledger, workforce and payroll, procurement, inventory, tax, "
         "FX, budgeting, and a conversational AI assistant — all under one API-first core that is "
         "legally compliant with National Bank of Ethiopia (NBE) directives and Ethiopian labour law.")
    c.drawString(18 * mm, y, p)
    y -= 30 * mm
    c.setFillColor(GREEN)
    c.setFont("Helvetica-Bold", 14)
    c.drawString(18 * mm, y, "Why ApexPay wins")
    y -= 8 * mm
    bullets(c, 18 * mm, y, [
        "Complete platform, not a feature.|  43 business modules work off a single ledger, so a sale, a payroll run, a vendor invoice, and a tax payment all reconcile automatically.",
        "AI assistant built in.|  A role-scoped Apex Assistant (RAG + tool-calling agents) answers finance, inventory and payroll questions for merchants, employees, vendors and employers — in English and Amharic.",
        "Lighter than legacy accounting.|  Real-time, API-first, Postgres-native GL — not a heavyweight desktop QuickBooks-style product.",
        "Ethiopia-native compliance.|  NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%, VAT/TOT, Fayda ID, Ethiopian calendar, labour proclamation 1156/2019.",
        "Live, verified data.|  The merchant dashboard is wired to the real API with seeded data; a 29-section end-to-end smoke runs green on a real Postgres in CI.",
    ], size=11, gap=6.5, bold_prefix=True)
    c.showPage()

    # ---------- Page 3: Market & Problem ----------
    header(c, "Market & Problem",
           "Ethiopia's economy is digitizing rapidly, yet business software remains fragmented.")
    y = PAGE_H - 40 * mm
    c.setFillColor(DARK)
    c.setFont("Helvetica", 11)
    c.drawString(18 * mm, y,
                 "Businesses still juggle a payment gateway, a payroll tool, spreadsheets, and desktop "
                 "accounting software that never talk to each other. A single transaction is re-entered "
                 "across systems, reconciliation is manual, payroll tax is error-prone, and there is no "
                 "unified view of cash, cost, or compliance.")
    y -= 24 * mm
    c.setFillColor(GREEN)
    c.setFont("Helvetica-Bold", 13)
    c.drawString(18 * mm, y, "Pain points we solve")
    y -= 7 * mm
    bullets(c, 18 * mm, y, [
        "Fragmentation — payments, payroll, accounting and procurement live in disconnected silos.",
        "Manual reconciliation — no shared ledger means cash never ties out automatically.",
        "Compliance burden — NBE, tax, pension and labour rules are complex and change often.",
        "Inflexible pricing — global gateways charge FX + cross-border fees on local ETB rails.",
        "No data gravity — no single source of truth for finance teams or AI to reason over.",
    ], size=10.5)
    c.showPage()

    # ---------- Page 4: Product & Architecture ----------
    header(c, "Product & Architecture",
           "One Go API core, one PostgreSQL ledger, four client surfaces, plus a self-service employee portal.")
    y = PAGE_H - 40 * mm
    c.setFillColor(GREEN)
    c.setFont("Helvetica-Bold", 12)
    c.drawString(18 * mm, y, "Core platform facts")
    y -= 6 * mm
    bullets(c, 18 * mm, y, [
        "API-first — REST under /v1, typed Go services, pgx + decimal (no float money).",
        "Real-time double-entry GL — journal entries, fiscal period close, depreciation, inventory COGS, tax schedules, multi-currency FX revaluation.",
        "Event-driven — outbox pattern + workers for reliable async jobs (webhooks, dunning, recon, notifications).",
        "Four client surfaces — Merchant Web, Admin Web, Checkout Web, Mobile (Flutter, Android APK built in CI).",
        "Merchant dashboard wired to live API — payments, subscriptions, refunds, payouts, payroll (runs, calendars, final settlements, reports, tax brackets) render real seeded data.",
        "i18n — per-user English/Amharic across API and UI.",
    ], size=10.5, gap=6.2)
    y -= 6 * mm
    # architecture diagram
    if os.path.exists(CH["architecture.png"]):
        c.drawImage(CH["architecture.png"], 20 * mm, y - 62 * mm, width=150 * mm, height=70 * mm,
                    preserveAspectRatio=True, mask="auto")
    c.showPage()

    # ---------- Page 5: Stakeholders & AI ----------
    header(c, "How Each Stakeholder Uses ApexPay",
           "A role-scoped conversational agent answers within each actor's scope, grounded in live ledger/data.")
    y = PAGE_H - 40 * mm
    c.setFillColor(DARK)
    c.setFont("Helvetica", 11)
    c.drawString(18 * mm, y,
                 "The Apex Assistant authenticates the caller, resolves whether they are a merchant, "
                 "employee, vendor or admin, and answers only within that scope — grounded in live "
                 "ledger/data via tools, never hallucinated.")
    y -= 18 * mm
    c.setFillColor(GREEN)
    c.setFont("Helvetica-Bold", 13)
    c.drawString(18 * mm, y, "Who uses it")
    y -= 7 * mm
    bullets(c, 18 * mm, y, [
        "Merchant / Owner|  payments, cash position, P&L, balance sheet, inventory, invoices, budget variance, tax due. Approves payroll, closes fiscal periods, posts manual journals.",
        "Employee|  own YTD pay, leave balance, expense claims, payslips — via a self-service employee portal.",
        "Vendor / Supplier|  self-service portal: their AP invoices and payment status only.",
        "Customer|  hosted checkout: Telebirr, CBE Birr, bank, card, EthSwitch QR; 2FA above the NBE threshold.",
        "Admin / Compliance|  onboarding queue, KYC/Fayda exam, compliance checks, approve/reject.",
    ], size=10.5, bold_prefix=True, gap=6.5)
    c.showPage()

    # ---------- Page 6: Quality, Security & Compliance ----------
    header(c, "Quality, Security & Compliance",
           "All 10 CI quality gates pass, plus a mobile APK build job.")
    y = PAGE_H - 40 * mm
    bullets(c, 18 * mm, y, [
        "Security|  gosec, trivy (SCA), gitleaks (no secrets), fin-privacy (no plain FIN in logs), no-float-money.",
        "Auditability|  append-only audit_logs with a DB trigger; append-only assistant threads & messages.",
        "Correctness|  SQL-param-cast lint prevents a whole class of pgx errors; double-entry balance validated on every journal.",
        "E2E|  docker-smoke runs 29 sections against real Postgres (auth, payments, payroll, GL, tax, FX, procurement, portals, AI, i18n).",
        "Mobile|  android.yml builds the merchant APK on a large CI runner and uploads it as an artifact.",
        "Compliance|  NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%, VAT/TOT, Fayda ID, labour proclamation 1156/2019.",
    ], size=10.5, bold_prefix=True, gap=6.5)
    y -= 16 * mm
    if os.path.exists(CH["security.png"]):
        c.drawImage(CH["security.png"], 20 * mm, y - 48 * mm, width=150 * mm, height=55 * mm,
                    preserveAspectRatio=True, mask="auto")
    c.showPage()

    # ---------- Page 7: Roadmap & Opportunity ----------
    header(c, "Delivery Roadmap & The Opportunity",
           "Eight milestones delivered; merchant data path now verified end-to-end.")
    y = PAGE_H - 40 * mm
    c.setFillColor(DARK)
    c.setFont("Helvetica", 11)
    c.drawString(18 * mm, y,
                 "From core payments to a hardened, AI-native operating system, ApexPay was built in "
                 "eight increments that each shipped working software: Core Payments, Real-time GL, "
                 "Workforce OS, AI Assistant, Procurement/AP, i18n & UI/UX, Budget & Portals, and "
                 "Security Hardening. Every increment kept build, vet, test, and the end-to-end smoke "
                 "suite green. The latest increment wired the remaining merchant dashboard pages to "
                 "live API data (payments detail, payroll calendar/final-settlements/reports/settings, "
                 "payouts, subscriptions detail, refunds detail).")
    y -= 24 * mm
    if os.path.exists(CH["roadmap.png"]):
        c.drawImage(CH["roadmap.png"], 20 * mm, y - 34 * mm, width=160 * mm, height=38 * mm,
                    preserveAspectRatio=True, mask="auto")
    y -= 44 * mm
    c.setFillColor(GREEN)
    c.setFont("Helvetica-Bold", 13)
    c.drawString(18 * mm, y, "The Opportunity for Investors")
    y -= 7 * mm
    bullets(c, 18 * mm, y, [
        "Defensible moat — an integrated ledger + AI assistant is hard to copy feature-by-feature.",
        "Land & expand — start with payments, expand to payroll, accounting, procurement, and finance ops.",
        "Live proof — the full data path is running against real seeded data, not mockups.",
        "Ethiopia-first — a locally compliant, locally-routed platform no global vendor matches.",
    ], size=10.5)

    c.save()
    print("Wrote", OUT)


if __name__ == "__main__":
    build()
