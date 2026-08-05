"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockTaxPayments = [
  { id: "tax_001", tax_type: "vat", amount: "7500", currency: "ETB", period_month: 7, period_year: 2026, due_date: "2026-08-15", status: "pending", challan_file_key: null, payment_reference: null, created_at: "2026-07-25" },
  { id: "tax_002", tax_type: "withholding", amount: "1000", currency: "ETB", period_month: 7, period_year: 2026, due_date: "2026-08-15", status: "paid", challan_file_key: "tax/challan_withholding_2026_07.pdf", challan_file_hash: "hash_challan_withholding_2026_07", payment_reference: "UTR-CBE-123456", paid_at: "2026-08-10", created_at: "2026-07-25" },
  { id: "tax_003", tax_type: "paye", amount: "20000", currency: "ETB", period_month: 7, period_year: 2026, due_date: "2026-08-15", status: "pending", challan_file_key: null, created_at: "2026-07-25" },
  { id: "tax_004", tax_type: "pension", amount: "36000", currency: "ETB", period_month: 7, period_year: 2026, due_date: "2026-08-15", status: "paid", challan_file_key: "tax/challan_pension_2026_07.pdf", payment_reference: "UTR-CBE-789012", paid_at: "2026-08-10" },
]

export default function TaxPaymentsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Tax Payments • የግብር ክፍያዎች • GST TDS Advance Tax in Under 30 Seconds Pre-filled Forms Automated TDS Payments Paid On Time Challans Shared in Inbox for Filing Collaborate with Accountant from Dashboard • RazorpayX Parity • P0 • Ethiopia VAT 15% TOT Withholding 2% PAYE Pension 7%/11%</h1>
            <p className="text-sm text-muted-foreground mt-2">Automate tax payments with pre-filled tax payment forms GST TDS & Advance Tax in under 30 seconds from one dashboard pay GST TDS & Advance Tax automated TDS payments paid on time and challans shared in inbox for filing collaborate with accountant from dashboard for easy and timely tax disbursal never miss payment again tax payments are automated paid on time and challans shared in inbox for tax filing automated tax payments paid on time and challans shared TDS challans available on dashboard for auditing and tax-filing TDS payments are automated paid on time and challans shared payslips auto-generated salaries are credited directly employees access via self-service portal TDS challans shared in inbox for tax filing purposes calculate payroll and disburse salaries in few clicks RazorpayX Payroll automates payments and filings of compliances like TDS PF ESI PT and much more tax payments are automated paid on time and challans shared • Ethiopia: VAT 15% on goods/services TOT 2% or 10% Withholding Tax 2% for services per Income Tax Proclamation 286/2002 PAYE Employment Income Tax per ET brackets 0-600 0% etc 601-1650 10%-60 1651-3200 15%-142.5 3201-5250 20%-302.5 5251-7800 25%-565 7801-10900 30%-955 >10900 35%-1500 Pension 7% Employee 11% Employer Total 18% per Private Organization Employees Pension Proclamation 1268/2022 • Outstanding modern UI glassmorphic Recharts</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Tax Payment • VAT TOT Withholding PAYE Pension • Pre-filled Forms • Challans Inbox • Accountant Collaboration • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Tax Payments • VAT 15% TOT 2%/10% Withholding 2% PAYE Pension 7%/11% • Pre-filled Forms • Challans Shared in Inbox for Filing • Collaborate with Accountant • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Tax Type • VAT TOT Withholding PAYE Pension</span><span>Amount</span><span>Period</span><span>Due Date</span><span>Status • Draft/Pending/Paid/Failed</span><span>Challan • File Key Hash • Payment Reference • Paid At</span><span>Action • Pay • Download Challan • Accountant Collaboration</span></div>
              {mockTaxPayments.map(tax => (
                <div key={tax.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span><Badge variant={tax.tax_type==="vat" ? "default" : tax.tax_type==="withholding" ? "warning" : tax.tax_type==="paye" ? "danger" : "success"}>{tax.tax_type}</Badge><span className="block text-[10px]">{tax.tax_type==="vat" ? "VAT 15% on goods/services per Ethiopia VAT Proclamation" : tax.tax_type==="withholding" ? "Withholding Tax 2% for services per Income Tax Proclamation 286/2002" : tax.tax_type==="paye" ? "PAYE Employment Income Tax per ET brackets 0-600 0% etc binary search O(log n)" : tax.tax_type==="pension" ? "Pension 7% Employee 11% Employer Total 18% per Private Organization Employees Pension Proclamation 1268/2022" : tax.tax_type}</span></span>
                  <span className="font-bold">ETB {tax.amount} {tax.currency}</span>
                  <span>{tax.period_month}/{tax.period_year}</span>
                  <span>{tax.due_date} • Due Date • Before due date pay to avoid late fees • Schedule payouts and avoid losses due to late fees • Outstanding</span>
                  <span><Badge variant={tax.status==="paid" ? "success" : tax.status==="pending" ? "warning" : "default"}>{tax.status}</Badge><span className="block text-[10px]">Challan {tax.challan_file_key || "—"} • Payment Ref {tax.payment_reference || "—"} • Paid At {tax.paid_at || "—"} • Created {tax.created_at}</span></span>
                  <span className="text-[11px]"><span className="font-mono text-[10px] block">Challan File Key {tax.challan_file_key || "—"} • Hash {tax.challan_file_hash || "—"} • MinIO presigned 15m • Encrypted SSE-S3 • 7y retention NBE • File key tax/challan_{tax.tax_type}_{tax.period_year}_{tax.period_month}.pdf • No plain FIN logs • Hash integrity • Outstanding modern UI glassmorphic</span><span className="block">Payment Reference {tax.payment_reference || "—"} • UTR-CBE-123456 • Bank Payment Reference • UTR • Outstanding • For Reconciliation • Audit Trail • Tracking Who Created Approved and Processed Every Single Payout for Easy Reconciliation and Auditing</span></span>
                  <span className="flex flex-col gap-1"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Pay • Automated TDS Payments Paid On Time Challans Shared in Inbox for Filing • Pre-filled Tax Payment Forms • Pay GST TDS & Advance Tax in Under 30 Seconds • Outstanding</button><button className="rounded-xl border h-7 px-3 text-[10px]">Download Challan • Challan Shared in Inbox for Filing • Collaborate with Accountant from Dashboard • Easy Timely Tax Disbursal Never Miss Payment</button><button className="rounded-xl border h-7 px-3 text-[10px]">View • Accountant Collaboration • CA Portal • Dedicated CA Portal • Collaborate with Accountant</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Tax Payments Automated Pre-filled Forms Challans Inbox Accountant Collaboration • Outstanding Pipeline Visual Stepper • Maker-checker • VAT 15% TOT Withholding 2% • Ethiopia Law Compliance Gold</p>
              <div className="mt-3 flex items-center gap-2">
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">1</span><span>Create Tax Payment Draft • Tax Type VAT/TOT/Withholding/PAYE/Pension/Corporate Tax/Excise/Other Amount Currency Period Month Year Due Date Status draft/pending_approval/pending/paid/failed/cancelled Challan File Key File Hash Payment Reference Paid At Created By Approved By • Pre-filled Forms Based on Payroll Data TDS Monthly VAT 15% TOT Withholding 2% • Outstanding modern UI glassmorphic</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-amber-500 text-white flex items-center justify-center">2</span><span>Pending Approval • Maker-checker dual approval risk>=70 or TPV>1M or payroll net >100k requires 2 approvers approver != submitter per NBE controls • Approval Flow • Outstanding avatars HR Finance Admin • 2FA mandatory >5000 ETB per ONPS/10/2025 • Maker-checker >50k payout >100k payroll</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center">3</span><span>Paid • Automated TDS Payments Paid On Time Challans Shared in Inbox for Filing • Generate Challan via Partner Bank API Mock • Schedule Tax Payments • Notification Before Due Date • Inbox Challans Shared for Filing • Audit Trail • Tracking Who Created Approved and Processed Every Single Payout for Easy Reconciliation and Auditing • Outstanding • Tax Payments Automated Pre-filled Forms Challans Inbox Accountant Collaboration VAT 15% TOT Withholding 2% • Ethiopia Law Compliance Gold</span></div>
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">VAT 15% • TOT 2%/10% • Withholding Tax 2% for Services per Income Tax Proclamation 286/2002 • PAYE Employment Income Tax per ET Brackets 0-600 0% etc • Pension 7% Employee 11% Employer Total 18% • Outstanding • Ethiopia Law Compliance Gold</h3>
            <div className="mt-4 space-y-3 text-xs">
              <div className="rounded-xl bg-muted p-3"><p className="font-medium">VAT 15% on Goods/Services per Ethiopia VAT Proclamation • TOT 2% or 10%? Withholding Tax 2% for Services per Income Tax Proclamation 286/2002 • PAYE Employment Income Tax per ET Brackets 0-600 0% 601-1650 10%-60 1651-3200 15%-142.5 3201-5250 20%-302.5 5251-7800 25%-565 7801-10900 30%-955 >10900 35%-1500 • Pension 7% Employee 11% Employer Total 18% per Private Organization Employees Pension Proclamation 1268/2022 • Outstanding • Ethiopia Law Compliance Gold</p><p className="text-[11px] text-muted-foreground mt-1">VAT 15% on goods/services • TOT 2% for goods 10% for services? Actually Ethiopia VAT 15% for goods and services above threshold, TOT 2% for goods 10% for services if not VAT registered? Withholding Tax 2% for services per Income Tax Proclamation 286/2002 Art 87? 2% for services, 30% for imports? • Outstanding modern UI glassmorphic • Recharts • Tax Brackets Binary Search O(log n) 7 brackets • Formula tax=taxable*rate-deduction rounded 2 decimals • ValidateTIN 10-digit numeric regex ^[0-9]{{10}}$ • Taxable Income ET Gross - Pension Employee 7% - tax_exempt_allowances • YTD cumulative annual certificate Form equivalent India Form16</p></div>
              <div className="rounded-xl border p-3"><p className="font-medium">Pre-filled Tax Forms Based on Payroll Data • TDS Monthly • VAT 15% • TOT • Withholding Tax 2% • Outstanding</p><p className="text-[11px]">Payroll Data: total_gross 200k total_tax 20k total_pension 14k/22k total_net 150k • VAT 15% on goods/services above threshold • TOT 2% for goods 10% for services if not VAT registered • Withholding Tax 2% for services per Income Tax Proclamation 286/2002 Art 87? 2% for services • Pre-filled Forms: Tax Type VAT/TOT/Withholding/PAYE/Pension Amount Currency Period Month Year Due Date Status draft/pending_approval/pending/paid/failed/cancelled Challan File Key File Hash Payment Reference Paid At Created By Approved By • Outstanding modern UI glassmorphic • Recharts • Tax Brackets Binary Search O(log n)</p></div>
              <div className="rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]">
                <p className="font-semibold">Accountant Collaboration • CA Portal • Dedicated CA Portal • Collaborate with Accountant • Outstanding • RazorpayX Parity • P0</p>
                <p className="mt-1">Tax Accountants table id merchant_id user_id role compliance/accountant/auditor/viewer permissions jsonb created_at unique merchant_id user_id + API POST /v1/tax_accountants + GET /v1/tax_accountants + CA Portal: Dedicated CA Portal • Collaborate with Accountant from Dashboard for Easy and Timely Tax Disbursal Never Miss Payment Again • Tax Payments Are Automated Paid On Time and Challans Shared • Real-time Insights • Outstanding modern UI glassmorphic • Recharts • Tax Brackets Binary Search O(log n)</p>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
