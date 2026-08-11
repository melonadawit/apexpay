#!/usr/bin/env python3
"""Build a one-page ApexPay investor one-pager PDF."""
import os
from reportlab.lib.pagesizes import A4
from reportlab.lib.units import mm
from reportlab.lib.colors import HexColor, white
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.platypus import (SimpleDocTemplate, Paragraph, Spacer, HRFlowable)

GREEN = HexColor("#0B6E4F")
DARK = HexColor("#1a202c")
GREY = HexColor("#4a5568")

OUT = os.path.join(os.path.dirname(__file__), "ApexPay_One_Pager.pdf")


def build():
    styles = getSampleStyleSheet()
    title = ParagraphStyle("title", parent=styles["Title"], fontName="Helvetica-Bold",
                           textColor=DARK, fontSize=26, spaceAfter=2, leading=30)
    tag = ParagraphStyle("tag", parent=styles["Normal"], fontName="Helvetica",
                         textColor=GREEN, fontSize=13, spaceAfter=10)
    h = ParagraphStyle("h", parent=styles["Heading2"], fontName="Helvetica-Bold",
                       textColor=GREEN, fontSize=14, spaceBefore=12, spaceAfter=5)
    body = ParagraphStyle("body", parent=styles["BodyText"], fontName="Helvetica",
                          textColor=DARK, fontSize=10.5, leading=15, spaceAfter=4)
    small = ParagraphStyle("small", parent=body, fontSize=9, textColor=GREY)

    doc = SimpleDocTemplate(OUT, pagesize=A4,
                            leftMargin=20 * mm, rightMargin=20 * mm,
                            topMargin=18 * mm, bottomMargin=18 * mm,
                            title="ApexPay — Investor One-Pager")
    story = []
    story.append(Paragraph("ApexPay", title))
    story.append(Paragraph("The all-in-one Financial Operating System for Ethiopia.", tag))
    story.append(HRFlowable(width="100%", thickness=1.2, color=GREEN))
    story.append(Paragraph(
        "AI-native, API-first platform unifying payments, a real-time double-entry General Ledger, "
        "workforce &amp; payroll, procurement/AP, inventory, tax, FX, budgeting, and a role-scoped "
        "AI assistant — one compliant, Ethiopia-native core.", body))
    story.append(Spacer(1, 4))

    story.append(Paragraph("What was verified this cycle (all green on master)", h))
    story.append(Paragraph(
        "• <b>Merchant dashboard on live data.</b> Payments (list + transaction detail with ledger "
        "journals), subscriptions (+ detail with plan/customer/invoices), refunds, payouts, and "
        "payroll (runs, calendars, final settlements, reports, tax brackets) all render from a real "
        "seeded Postgres, not mockups.", body))
    story.append(Paragraph(
        "• <b>Developer portal.</b> Real, DB-backed API-key management (list/create/revoke) — created "
        "keys are immediately usable as Bearer tokens — plus webhooks endpoints &amp; deliveries.", body))
    story.append(Paragraph(
        "• <b>Embedded finance.</b> Lending, escrow, corporate cards, credit lines, and virtual "
        "accounts all return live seeded data.", body))
    story.append(Paragraph(
        "• <b>Mobile (Flutter).</b> Tests pass 3/3; analyzer clean; Android APK built as a CI artifact "
        "on a large runner.", body))
    story.append(Paragraph(
        "• <b>Security audit.</b> Money-safety, FIN-privacy, secret-handling, and webhook SSRF controls "
        "verified. A real gap was found and closed: audit_logs is now truly append-only via a DB "
        "trigger, and CI now fails if the trigger is ever removed.", body))

    story.append(Paragraph("Why ApexPay wins", h))
    story.append(Paragraph(
        "• <b>Complete platform, not a feature.</b> 43 modules share one ledger, so a sale, a payroll "
        "run, a vendor invoice, and a tax payment all reconcile automatically.", body))
    story.append(Paragraph(
        "• <b>AI assistant built in.</b> Role-scoped, grounded in live ledger/data, English + Amharic.", body))
    story.append(Paragraph(
        "• <b>Ethiopia-native compliance.</b> NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%, "
        "VAT/TOT, Fayda ID, labour proclamation 1156/2019.", body))
    story.append(Paragraph(
        "• <b>Production-grade.</b> 10/10 CI gates green, including a 29-section end-to-end smoke on "
        "real Postgres; money is always decimal.Decimal.", body))

    story.append(Paragraph("The opportunity", h))
    story.append(Paragraph(
        "Defensible moat (integrated ledger + AI assistant), land-and-expand from payments to the full "
        "finance stack, and <b>live proof</b> — the entire data path is running against real data, not "
        "mockups.", body))
    story.append(Spacer(1, 6))
    story.append(HRFlowable(width="100%", thickness=0.8, color=GREEN))
    story.append(Paragraph(
        "References: docs/SAD.md · docs/SDD.md · docs/SECURITY_AUDIT.md · full investor deck.", small))

    doc.build(story)
    print("Wrote", OUT)


if __name__ == "__main__":
    build()
