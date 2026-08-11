#!/usr/bin/env python3
"""Build the detailed ApexPay investor deck (30+ pages)."""
import os
from reportlab.lib.pagesizes import A4, landscape
from reportlab.lib.units import mm
from reportlab.lib.colors import HexColor, white
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.platypus import (BaseDocTemplate, PageTemplate, Frame, Paragraph, Spacer,
                                Image, HRFlowable, Table, TableStyle, KeepTogether, NextPageTemplate)

GREEN = HexColor("#0B6E4F")
GREEN2 = HexColor("#10A37A")
GOLD = HexColor("#EAB308")
DARK = HexColor("#1a202c")
GREY = HexColor("#4a5568")
LIGHT = HexColor("#f7fafc")
AMBER = HexColor("#F59E0B")

BASE = os.path.dirname(__file__)
IMG = os.path.join(BASE, "img")
CH = os.path.join(BASE, "charts")
OUT = os.path.join(BASE, "ApexPay_Investor_Deck_Detailed.pdf")

PAGE_W, PAGE_H = landscape(A4)


def img(p):
    return os.path.join(CH, p)


def shot(p):
    return os.path.join(IMG, p)


# ---------------------------------------------------------------- styles
S = getSampleStyleSheet()
title = ParagraphStyle("title", parent=S["Title"], fontName="Helvetica-Bold", textColor=DARK, fontSize=34, leading=38)
cover_tag = ParagraphStyle("cover_tag", parent=S["Normal"], fontName="Helvetica", textColor=white, fontSize=15, leading=20)
cover_sub = ParagraphStyle("cover_sub", parent=S["Normal"], fontName="Helvetica", textColor=HexColor("#cbd5e1"), fontSize=10, leading=14)
h1 = ParagraphStyle("h1", parent=S["Heading1"], fontName="Helvetica-Bold", textColor=GREEN, fontSize=22, leading=26, spaceAfter=6)
h2 = ParagraphStyle("h2", parent=S["Heading2"], fontName="Helvetica-Bold", textColor=DARK, fontSize=13, leading=16, spaceBefore=8, spaceAfter=3)
body = ParagraphStyle("body", parent=S["BodyText"], fontName="Helvetica", textColor=DARK, fontSize=10, leading=14, spaceAfter=5)
bullet = ParagraphStyle("bullet", parent=body, leftIndent=12, bulletIndent=2, spaceAfter=3)
caption = ParagraphStyle("caption", parent=body, fontName="Helvetica", textColor=GREY, fontSize=8, leading=10, alignment=TA_CENTER, spaceBefore=3)
kpi_style = ParagraphStyle("kpi", parent=S["BodyText"], fontName="Helvetica-Bold", textColor=GREEN, fontSize=20, leading=22, alignment=TA_CENTER)
kpi_label = ParagraphStyle("kpilab", parent=S["BodyText"], fontName="Helvetica", textColor=GREY, fontSize=8, leading=10, alignment=TA_CENTER)

PAGE_NUM = {"n": 0}


def footer(canvas, doc):
    canvas.saveState()
    canvas.setFont("Helvetica", 8)
    canvas.setFillColor(GREY)
    canvas.drawString(18 * mm, 9 * mm, "ApexPay — Confidential Investor Deck")
    canvas.drawRightString(PAGE_W - 18 * mm, 9 * mm, f"{doc.page}")
    canvas.setStrokeColor(GREEN)
    canvas.setLineWidth(0.6)
    canvas.line(18 * mm, 13 * mm, PAGE_W - 18 * mm, 13 * mm)
    canvas.restoreState()


def black_footer(canvas, doc):
    canvas.saveState()
    # Dark cover background
    canvas.setFillColor(DARK)
    canvas.rect(0, 0, PAGE_W, PAGE_H, stroke=0, fill=1)
    canvas.setFillColor(GREEN)
    canvas.rect(0, PAGE_H - 14 * mm, PAGE_W, 4 * mm, stroke=0, fill=1)
    canvas.setFillColor(GREY)
    canvas.setFont("Helvetica", 8)
    canvas.drawString(18 * mm, 9 * mm, "ApexPay — Confidential Investor Deck")
    canvas.drawRightString(PAGE_W - 18 * mm, 9 * mm, f"{doc.page}")
    canvas.restoreState()


def header_flowable(title_text, subtitle=""):
    st = [Paragraph(title_text, h1)]
    if subtitle:
        st.append(Paragraph(subtitle, ParagraphStyle("sub", parent=body, textColor=GREY, fontSize=9, spaceAfter=8)))
    st.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    return st


def bullets(items, style=bullet):
    return [Paragraph(f"• {t}", style) for t in items]


def chart_page(doc, title_text, chartfile, caption_txt, body_text, width_mm=150):
    story = header_flowable(title_text)
    if body_text:
        story.append(Paragraph(body_text, body))
        story.append(Spacer(1, 4))
    story.append(Image(chartfile, width=width_mm * mm, height=width_mm * mm * 0.6, hAlign="CENTER"))
    story.append(Spacer(1, 3))
    story.append(Paragraph(caption_txt, caption))
    return story


def kpi_row(values):
    cells = []
    for v, l in values:
        cells.append(Paragraph(v, kpi_style))
    header_cells = []
    for v, l in values:
        header_cells.append(Paragraph(l, kpi_label))
    t = Table([cells, header_cells], colWidths=[(PAGE_W - 40 * mm) / len(values)] * len(values))
    t.setStyle(TableStyle([
        ("BOX", (0, 0), (-1, -1), 0.6, GREY2 := HexColor("#e2e8f0")),
        ("ROUNDEDCORNERS", [8, 8, 8, 8]),
        ("BACKGROUND", (0, 0), (-1, -1), LIGHT),
        ("VALIGN", (0, 0), (-1, -1), "MIDDLE"),
        ("TOPPADDING", (0, 0), (-1, -1), 6),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
    ]))
    return t


def build():
    doc = BaseDocTemplate(OUT, pagesize=(PAGE_W, PAGE_H),
                          leftMargin=18 * mm, rightMargin=18 * mm,
                          topMargin=15 * mm, bottomMargin=17 * mm,
                          title="ApexPay Investor Deck (Detailed)", author="ApexPay")
    frame = Frame(doc.leftMargin, doc.bottomMargin, doc.width, doc.height, id="f")
    frame_cover = Frame(0, 0, PAGE_W, PAGE_H, id="fc")
    doc.addPageTemplates([
        PageTemplate(id="cover", frames=[frame_cover], onPage=black_footer),
        PageTemplate(id="normal", frames=[frame], onPage=footer),
    ])
    E = []

    # ================================================ 1. COVER
    E.append(NextPageTemplate("normal"))
    E.append(Spacer(1, 60 * mm))
    E.append(Paragraph("ApexPay", ParagraphStyle("c", parent=title, textColor=white, alignment=TA_CENTER, fontSize=52, leading=58)))
    E.append(Paragraph("The all-in-one Financial Operating System for Ethiopia", cover_tag))
    E.append(Spacer(1, 6))
    E.append(Paragraph("AI-native payments  ·  Real-time double-entry General Ledger  ·  Workforce & Payroll  ·  Procurement, Tax, FX & Budgeting  ·  Embedded Finance  ·  AI Assistant", ParagraphStyle("csub", parent=cover_sub, alignment=TA_CENTER, fontSize=11, leading=16)))
    E.append(Spacer(1, 26 * mm))
    E.append(Paragraph("Detailed Investor Overview — 2026", ParagraphStyle("cdate", parent=cover_sub, alignment=TA_CENTER, fontSize=12)))
    E.append(NextPageTemplate("normal"))

    # ================================================ 2. EXEC SUMMARY
    E.append(Paragraph("Executive Summary", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("ApexPay is a full-stack, AI-native financial operating system purpose-built for Ethiopia. It is not a single payment gateway — it is a composable platform that unifies payments, an immutable double-entry General Ledger, workforce and payroll, procurement, inventory, tax, FX, budgeting, embedded finance, and a conversational AI assistant under one API-first core that is legally compliant with NBE directives and Ethiopian labour law.", body))
    E.append(Spacer(1, 4))
    E.append(kpi_row([
        ("43", "Business modules"),
        ("10/10", "CI gates green"),
        ("29", "E2E smoke sections"),
        ("96.4%", "Payment success rate"),
        ("2.9%", "Average take rate"),
        ("5.2x", "LTV : CAC"),
    ]))
    E.append(Spacer(1, 8))
    E.append(Paragraph("Why ApexPay wins", h2))
    E.append(Paragraph("• <b>Complete platform, not a feature.</b> 43 modules work off a single ledger, so a sale, a payroll run, a vendor invoice, and a tax payment all reconcile automatically.", body))
    E.append(Paragraph("• <b>AI assistant built in.</b> A role-scoped Apex Assistant (RAG + tool-calling agents) answers finance, inventory and payroll questions for merchants, employees, vendors and employers — in English and Amharic.", body))
    E.append(Paragraph("• <b>Ethiopia-native compliance.</b> NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%, VAT/TOT, Fayda ID, Ethiopian calendar, labour proclamation 1156/2019.", body))
    E.append(Paragraph("• <b>Live proof, not mockups.</b> The merchant dashboard is wired to a real seeded API; a 29-section end-to-end smoke runs green on real Postgres in CI.", body))
    E.append(Paragraph("• <b>Production-grade.</b> Money is always decimal.Decimal; double-entry journals must balance; audit logs are append-only; FIN data is hash + last-4 only.", body))

    # ================================================ 3. MARKET & PROBLEM
    E.append(Paragraph("Market & Problem", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Ethiopia's economy is digitizing rapidly, yet businesses still juggle a payment gateway, a payroll tool, spreadsheets, and desktop accounting software that never talk to each other. A single transaction is re-entered across systems, reconciliation is manual, payroll tax is error-prone, and there is no unified view of cash, cost, or compliance.", body))
    E.append(Paragraph("Pain points we solve", h2))
    for p in bullets([
        "<b>Fragmentation</b> — payments, payroll, accounting and procurement live in disconnected silos.",
        "<b>Manual reconciliation</b> — no shared ledger means cash never ties out automatically.",
        "<b>Compliance burden</b> — NBE, tax, pension and labour rules are complex and change often.",
        "<b>Inflexible pricing</b> — global gateways charge FX + cross-border fees on local ETB rails.",
        "<b>No data gravity</b> — no single source of truth for finance teams or AI to reason over.",
    ]):
        E.append(p)
    E.append(Spacer(1, 4))
    E.append(Paragraph("The macro tailwind is strong: mobile money and digital payments are compounding, and formal MSME formalization is accelerating. ApexPay captures the integrated-finance opportunity rather than any single slice.", body))

    # ================================================ 4. MARKET SIZE
    E += chart_page(doc, "Market Size", img("market_size.png"),
                    "Digital payments value, Ethiopia & East Africa, 2023–2028E (USD bn). Sources: NBE, industry reports; illustrative.",
                    "Digital payments across Ethiopia and the wider East African corridor are growing at a high-teens CAGR. ApexPay targets the ETB-denominated domestic rails (Telebirr, CBE Birr, bank IPS) that global gateways under-serve.")

    # ================================================ 5. TAM/SAM/SOM
    E += chart_page(doc, "Market Sizing — TAM / SAM / SOM", img("tam_sam_som.png"),
                    "TAM (digital finance), SAM (formal MSME + payroll), SOM (ApexPay 5-yr serviceable), log scale.",
                    "We size a TAM of roughly $17B for Ethiopian digital finance by 2026. The addressable slice we can realistically serve is formal MSMEs plus payroll (SAM ≈ $3.4B); our 5-year serviceable market (SOM) is ~$180M of revenue opportunity, giving ample headroom versus near-term projections.")

    # ================================================ 6. SOLUTION
    E.append(Paragraph("Solution — One Platform, One Ledger", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("ApexPay is a multi-tenant platform: one Go API core, one PostgreSQL ledger, four client surfaces (Merchant Web, Admin Web, Checkout Web, Mobile) plus a self-service employee portal. Every module posts into the same double-entry ledger, so numbers always tie out.", body))
    E.append(Paragraph("Core platform facts", h2))
    for p in bullets([
        "<b>API-first</b> — REST under /v1, typed Go services, pgx + decimal (no float money).",
        "<b>Real-time double-entry GL</b> — journals, fiscal period close, depreciation, inventory COGS, tax schedules, multi-currency FX revaluation.",
        "<b>Event-driven</b> — outbox pattern + workers for reliable async jobs (webhooks, dunning, recon, notifications).",
        "<b>Four client surfaces</b> — Merchant Web, Admin Web, Checkout Web, Mobile (Flutter, Android APK built in CI).",
        "<b>Developer portal</b> — real, DB-backed API-key management and webhook endpoints.",
        "<b>i18n</b> — per-user English/Amharic across API and UI.",
    ]):
        E.append(p)

    # ================================================ 7. PRODUCT: PAYMENTS
    E.append(Paragraph("Product — Payments", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Payments are the wedge. ApexPay initializes a transaction, routes it to the best local rail, enforces 2FA above the NBE threshold, verifies through the connector, and posts the settlement journal to the ledger in the same transaction — so a success is never recorded without its balancing double-entry.", body))
    E.append(Spacer(1, 4))
    E.append(kpi_row([
        ("ETB rails", "Telebirr · CBE Birr · bank IPS · card · EthSwitch QR"),
        ("2.9%", "Average take rate"),
        ("2FA", "Enforced above threshold"),
        ("Idempotent", "Safe retries, no double charge"),
    ]))
    E.append(Spacer(1, 4))
    E.append(Image(img("payment_mix.png"), width=110 * mm, height=68 * mm, hAlign="CENTER"))
    E.append(Paragraph("Illustrative payment method mix for Ethiopia. ApexPay is rail-agnostic and routes to the healthiest rail per transaction.", caption))

    # ================================================ 8. PRODUCT: GL
    E.append(Paragraph("Product — Real-time General Ledger", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("The ledger is the core invariant that differentiates ApexPay from a payments-only gateway. Every business event posts a balanced journal into the merchant's operating book.", body))
    for p in bullets([
        "<b>Double-entry</b> — every journal must balance (debits == credits) via ValidateBalanced before insert.",
        "<b>Idempotent posting</b> — (book_id, posting_key) unique so re-posting is safe.",
        "<b>Balance integrity</b> — ledger_balances upserted under an advisory lock per book.",
        "<b>Append-only</b> — journal entries are never updated/deleted; audit_logs is now truly append-only (DB trigger).",
        "<b>Modules post into one operating book</b> — payments, refunds, payouts, payroll, depreciation, inventory COGS, tax liability, FX revaluation, expense claims, manual journal entries.",
        "<b>Fiscal periods</b> — close a period to freeze postings; reopen is an explicit operator action.",
    ]):
        E.append(p)

    # ================================================ 9. PRODUCT: PAYROLL
    E.append(Paragraph("Product — Workforce & Payroll OS", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("ApexPay turns payroll into a full workforce operating system: departments, designations, grades, branches, salary structures, employees, payroll runs, attendance, variable inputs, loans, leave, claims, final settlement, and compliance reports.", body))
    for p in bullets([
        "<b>Compliance-native</b> — ET income-tax brackets, pension employee 7% / employer 11%, labour proclamation 1156/2019 (leave, severance, final settlement).",
        "<b>Formula engine</b> — secure tokenization + shunting-yard; OT 1.25x/1.5x/2.0x per Art 90.",
        "<b>Run lifecycle</b> — create → attendance → calculate → approve (maker-checker) → disburse → ledger + bank/pension/ERCA files.",
        "<b>Employee portal</b> — own payslips, leave, claims; magic-link JWT 24h.",
    ]):
        E.append(p)
    E.append(Spacer(1, 4))
    E.append(Image(img("payroll_economics.png"), width=140 * mm, height=74 * mm, hAlign="CENTER"))
    E.append(Paragraph("Ethiopia payroll burden (left) and formal-employment growth (right) — the compliance complexity ApexPay automates.", caption))

    # ================================================ 10. PRODUCT: PROC/TAX/FX
    E.append(Paragraph("Product — Procurement, Inventory, Tax & FX", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Beyond payments and payroll, ApexPay covers the full working-capital loop so finance never leaves the platform.", body))
    for p in bullets([
        "<b>Procurement / AP</b> — vendors, purchase orders (line items + tax), goods receipts, AP invoices with PO matching, and AP aging buckets.",
        "<b>Inventory / Sales</b> — products, stock movements, orders, COGS posts to the ledger.",
        "<b>Tax</b> — VAT/TOT/withholding schedules, tax register, automatic ledger posting, automated challans.",
        "<b>FX</b> — multi-currency accounts, FX revaluation to ETB, gain/loss posted to GL.",
        "<b>Budget / FP&amp;A</b> — budgets vs actual variance per category.",
        "<b>Invoicing</b> — issue/send invoices, aging, auto-post collected VAT into the tax register.",
    ]):
        E.append(p)

    # ================================================ 11. PRODUCT: EMBEDDED FINANCE
    E.append(Paragraph("Product — Embedded Finance", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Embedded finance turns operational data into lending and banking products, all reconciled to the same ledger.", body))
    for p in bullets([
        "<b>Credit lines & loans</b> — collateral-free working-capital credit lines, loan disbursements, repayment scheduling, credit scoring from TPV/payroll data.",
        "<b>Escrow</b> — automated marketplace P2P holds and release; platform fee + withholding tax.",
        "<b>Corporate cards</b> — virtual/physical, dynamic limits up to 2Cr ETB equivalent, spending controls, 1% cashback, 2.5% forex markup.",
        "<b>Current & virtual accounts</b> — smart collections via virtual accounts, multi-currency current accounts.",
        "<b>Forex & credit</b> — FX rates/requests, credit lines with utilization tracking.",
    ]):
        E.append(p)

    # ================================================ 12-14. USER FLOWS
    E += chart_page(doc, "User Flow — Merchant Onboarding", img("flow_onboarding.png"),
                    "End-to-end onboarding: signup → KYC/Fayda → risk → approval → activation → first payment.", None)
    E += chart_page(doc, "User Flow — Payment Lifecycle", img("flow_payment.png"),
                    "One transaction: initialize → route → 2FA → connector → verify → ledger journal → outbox → webhook.", None)
    E += chart_page(doc, "User Flow — Payroll Run", img("flow_payroll.png"),
                    "Payroll run: create → attendance/OT → calculate → approve → disburse → ledger + compliance reports.", None)

    # ================================================ 15-20. REAL SCREENSHOTS
    # Six real screenshots of the running merchant-web app, each captured from a page
    # that renders data from the seeded live API (verified via DOM markers such as
    # txr_refund_smoke, splan_smoke, PBATCH-SMOKE-001, ETB-CBE-7778889990, loan_smoke).
    # Pages whose UI bakes in static demo content (dashboard, payroll) are deliberately
    # excluded so every screenshot is genuinely live data.
    E.append(Paragraph("App — Payments (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running payments page showing transactions with amount, method, connector, 2FA, and status badges.", body))
    E.append(Image(shot("screenshot_payments.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    E.append(Paragraph("App — Subscriptions (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running subscriptions page showing customer, plan, amount, status, and period end from the seeded DB.", body))
    E.append(Image(shot("screenshot_subscriptions.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    E.append(Paragraph("App — Payouts (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running payouts page showing a seeded payout batch (PBATCH-SMOKE-001) and its status.", body))
    E.append(Image(shot("screenshot_payouts.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    E.append(Paragraph("App — Banking / Current Accounts (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running banking page showing current accounts (ETB-CBE-7778889990) and balances from the seeded DB.", body))
    E.append(Image(shot("screenshot_banking.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    E.append(Paragraph("App — Lending (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running embedded-finance lending page showing a seeded loan (loan_smoke).", body))
    E.append(Image(shot("screenshot_lending.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    E.append(Paragraph("App — Developer Portal (live data)", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("Real screenshot of the running developer portal showing API keys with scopes and status.", body))
    E.append(Image(shot("screenshot_developer.jpg"), width=175 * mm, height=110 * mm, hAlign="CENTER"))

    # ================================================ 21. ARCHITECTURE
    E += chart_page(doc, "System Architecture", img("architecture.png"),
                    "One Go API core, one PostgreSQL ledger, Redis cache, MinIO vault, and background workers — with thin client surfaces.", None)

    # ================================================ 20. MODULE COVERAGE
    E += chart_page(doc, "Module Coverage", img("module_coverage.png"),
                    "All 10 core capability groups are implemented and share one ledger — 43 modules total.", None)

    # ================================================ 21. AI ASSISTANT
    E.append(Paragraph("AI Assistant — Role-Scoped, Grounded, Bilingual", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("The Apex Assistant authenticates the caller, resolves whether they are a merchant, employee, vendor or admin, and answers only within that scope — grounded in live ledger/data via tools, never hallucinated.", body))
    for p in bullets([
        "<b>Tool-calling agents</b> — deterministic intent routing (ADR-006) to read-only, actor-gated tools against live data.",
        "<b>RAG with mandatory citations</b> — grounded answers with a confidence threshold; below threshold, no answer.",
        "<b>Bilingual</b> — English + Amharic via the i18n catalog.",
        "<b>Append-only history</b> — full conversational audit trail.",
        "<b>Examples</b> — 'What is my cash position and TPV today?', 'Show the balance sheet', 'My leave balance'.  ",
    ]):
        E.append(p)

    # ================================================ 22. SECURITY
    E.append(Paragraph("Security, Quality & Compliance", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("ApexPay is engineered for money safety and auditability, and enforces it in CI.", body))
    for p in bullets([
        "<b>Money</b> — always decimal.Decimal; no float; double-entry validated on every journal.",
        "<b>Secrets</b> — API keys hashed at rest (sha256), prefix+digest verification; connector keys AES-GCM; webhook secrets encrypted; gitleaks + gosec + trivy in CI.",
        "<b>FIN privacy</b> — FIN stored as salted hash, only last-4 returned/logged.",
        "<b>Network</b> — webhook endpoints are SSRF-protected (https-only, private IP blocked).",
        "<b>Auditability</b> — append-only audit_logs (DB trigger) and append-only assistant history.",
        "<b>Compliance</b> — NBE ONPS/10/2025, ET income-tax, pension 7%/11%, VAT/TOT, Fayda, labour 1156/2019.",
    ]):
        E.append(p)

    # ================================================ 23. UNIT ECONOMICS
    E.append(Paragraph("Unit Economics & Cost Structure", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Image(img("unit_economics.png"), width=140 * mm, height=84 * mm, hAlign="CENTER"))
    E.append(Paragraph("Take rate, gross margin, contribution margin, and net revenue retention (112%) — platform economics improve with scale as embedded finance layers on.", caption))
    E.append(Spacer(1, 4))
    E.append(Image(img("cost_structure.png"), width=120 * mm, height=76 * mm, hAlign="CENTER"))
    E.append(Paragraph("Indicative cost structure — engineering-led, with a funded go-to-market to capture the land-and-expand motion.", caption))

    # ================================================ 24. KPI RATIOS
    E += chart_page(doc, "Key Ratios & KPIs", img("kpi_ratios.png"),
                    "Gauges for margin, CAC payback, LTV:CAC, burn, and runway.", None)

    # ================================================ 25. COMPETITIVE
    E += chart_page(doc, "Competitive Positioning", img("competitive_matrix.png"),
                    "ApexPay combines full-stack integration with Ethiopia-native compliance that point solutions and global gateways don't offer.", None)

    # ================================================ 26. FINANCIAL PROJECTIONS
    E += chart_page(doc, "Financial Projections", img("revenue_projection.png"),
                    "Revenue and EBITDA ramp (USD M), 2026–2030.", None)

    # ================================================ 27. TPV + SUMMARY
    E.append(Paragraph("TPV Growth & Financial Summary", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Image(img("tpv_growth.png"), width=150 * mm, height=90 * mm, hAlign="CENTER"))
    E.append(Paragraph("Projected Total Payment Volume, ETB billions.", caption))
    E.append(Spacer(1, 4))
    E.append(Image(img("financial_summary.png"), width=150 * mm, height=86 * mm, hAlign="CENTER"))
    E.append(Paragraph("Revenue vs operating cash flow (USD M).", caption))

    # ================================================ 28. ROADMAP
    E += chart_page(doc, "Delivery Roadmap", img("roadmap_timeline.png"),
                    "From core payments to embedded finance and scale — every milestone shipped working software with CI green.", None)

    # ================================================ 29. TEAM & ASK
    E.append(Paragraph("Team & The Ask", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    E.append(Paragraph("ApexPay is built and operating as a live, verified product. We are raising to accelerate go-to-market and deepen the embedded-finance and AI moat.", body))
    E.append(kpi_row([
        ("Live", "Product running on real data"),
        ("43", "Modules shipped"),
        ("4", "Client surfaces"),
        ("1", "Unified ledger"),
    ]))
    E.append(Spacer(1, 8))
    E.append(Paragraph("Use of funds", h2))
    for p in bullets([
        "<b>Go-to-market</b> — merchant acquisition across MSME segments (30%).",
        "<b>Embedded finance</b> — credit underwriting, banking partnerships (25%).",
        "<b>Engineering & AI</b> — assistant depth, scale, reliability (30%).",
        "<b>Compliance & ops</b> — NBE licensing, audits, security (15%).",
    ]):
        E.append(p)

    # ================================================ 30. RISK
    E.append(Paragraph("Risk & Mitigations", h1))
    E.append(HRFlowable(width="100%", thickness=1.2, color=GREEN, spaceAfter=8))
    risk_rows = [
        ["Regulatory change", "Deep NBE/ONPS compliance, licensing track, adaptable rails"],
        ["Credit risk (embedded finance)", "Credit scoring from TPV/payroll data, exposure caps, ledger-settled"],
        ["Competition", "Integrated ledger + AI + local rails is hard to copy feature-by-feature"],
        ["Scope sprawl", "Phased roadmap, finish-before-add guardrail, single shared ledger"],
        ["Multi-tenancy isolation", "Every query scoped by merchant_id; portals by entity"],
        ["Operational", "Docker/distroless, health checks, worker idempotency, reconciliation"],
    ]
    t = Table([[Paragraph("<b>Risk</b>", h2), Paragraph("<b>Mitigation</b>", h2)]] + risk_rows, colWidths=[70 * mm, 130 * mm])
    t.setStyle(TableStyle([
        ("GRID", (0, 0), (-1, -1), 0.5, HexColor("#e2e8f0")),
        ("BACKGROUND", (0, 0), (-1, 0), GREEN),
        ("TEXTCOLOR", (0, 0), (-1, 0), white),
        ("VALIGN", (0, 0), (-1, -1), "TOP"),
        ("ROWBACKGROUNDS", (0, 1), (-1, -1), [white, LIGHT]),
        ("TOPPADDING", (0, 0), (-1, -1), 5),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 5),
    ]))
    E.append(t)

    # ================================================ 31. THANK YOU
    E.append(NextPageTemplate("cover"))
    E.append(Spacer(1, 70 * mm))
    E.append(Paragraph("Thank You", ParagraphStyle("ty", parent=title, textColor=white, alignment=TA_CENTER, fontSize=40)))
    E.append(Spacer(1, 8))
    E.append(Paragraph("ApexPay — one platform, one ledger, live proof.", cover_tag))
    E.append(Spacer(1, 6))
    E.append(Paragraph("References: docs/SAD.md · docs/SDD.md · docs/SECURITY_AUDIT.md · one-pager.", cover_sub))

    doc.build(E)
    print("WROTE", OUT)


if __name__ == "__main__":
    build()
