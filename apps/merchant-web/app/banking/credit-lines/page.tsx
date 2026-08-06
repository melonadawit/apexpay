"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockCreditLines = [
  { id: "cl_001", credit_limit: "2000000", available_credit: "1500000", utilized_credit: "500000", interest_rate: "18.00", status: "active", credit_score: 750, approved_at: "2026-01-15", created_at: "2026-01-10", disbursements: 3, total_disbursed: "500000", total_repaid: "100000", outstanding: "400000" },
  { id: "cl_002", credit_limit: "5000000", available_credit: "5000000", utilized_credit: "0", interest_rate: "16.00", status: "approved", credit_score: 800, approved_at: "2026-02-01", created_at: "2026-01-20", disbursements: 0, total_disbursed: "0", total_repaid: "0", outstanding: "0" },
]

const mockDisbursements = [
  { id: "ld_001", credit_line_id: "cl_001", amount: "200000", currency: "ETB", purpose: "working_capital", status: "disbursed", disbursed_at: "2026-02-01", due_date: "2026-05-01", repaid_amount: "100000", outstanding_amount: "100000", ledger_book_id: "lbk_loan_001", created_at: "2026-02-01" },
  { id: "ld_002", credit_line_id: "cl_001", amount: "300000", currency: "ETB", purpose: "inventory", status: "disbursed", disbursed_at: "2026-03-15", due_date: "2026-06-15", repaid_amount: "0", outstanding_amount: "300000", ledger_book_id: "lbk_loan_002", created_at: "2026-03-15" },
]

export default function CreditLinesPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Credit Lines • የብድር መስመሮች • Instant Loans Digital Lending Collateral-free Credit Lines — RazorpayX Capital Line of Credit — RazorpayX Parity • P0 • Ethiopia Business Banking Core</h1>
            <p className="text-sm text-muted-foreground mt-2">Instant Loans Digital Lending Collateral-free Credit Lines — RazorpayX Capital Line of Credit — RazorpayX Parity • P0 • Ethiopia Business Banking Core • Credit Lines Merchant ID Credit Limit up to 2Cr ETB Equivalent Available Credit Utilized Credit Interest Rate 18% per Annum Status Draft/Pending Approval/Approved/Active/Suspended/Closed/Rejected Credit Score 300-900 Based on TPV Payroll Data etc Approved By Approved At Created At Updated At Index Merchant Status • Loan Disbursements Credit Line ID Merchant ID Amount Currency Purpose Working Capital/Inventory/Payroll etc Status Pending/Approved/Disbursed/Repaid/Defaulted/Cancelled Disbursed At Due Date Repaid Amount Outstanding Amount Ledger Book ID Created By Created At Index Credit Line ID Status Merchant Status • Outstanding modern UI glassmorphic Recharts • Credit Scoring Based on TPV Payroll Data etc • Outstanding • RazorpayX Capital Line of Credit Company Registration Current Accounts Vendor Payments Payroll Partners for Startups & SMEs • Get In-depth Reporting into Cash Flow Trends • Get Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History with Razorpay • Instant Loans Without Collaterals • Collateral-free Credit Lines • Embedded Lending Products • Digital Lending 2.0 • Razorpay Capital Line of Credit</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Request Credit Line • Credit Limit up to 2Cr ETB Equivalent • Available Credit Utilized Credit Interest Rate 18% per Annum • Credit Scoring Based on TPV Payroll Data • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Credit Lines • Credit Limit up to 2Cr ETB Equivalent • Available Credit Utilized Credit Interest Rate 18% per Annum • Status Draft/Pending Approval/Approved/Active/Suspended/Closed/Rejected • Credit Score 300-900 Based on TPV Payroll Data etc • Outstanding • RazorpayX Capital Line of Credit</h3>
            <div className="mt-4 space-y-3">
              {mockCreditLines.map(cl => (
                <div key={cl.id} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">Credit Line {cl.id} • Limit {cl.credit_limit} ETB • Available {cl.available_credit} • Utilized {cl.utilized_credit} • Interest {cl.interest_rate}% • Status {cl.status} • Credit Score {cl.credit_score} • Approved {cl.approved_at} • Created {cl.created_at} • Disbursements {cl.disbursements} • Total Disbursed {cl.total_disbursed} • Total Repaid {cl.total_repaid} • Outstanding {cl.outstanding} • Outstanding per RazorpayX Capital Line of Credit</p><Badge variant={cl.status==="active" ? "success" : cl.status==="approved" ? "success" : "warning"}>{cl.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Credit Limit {cl.credit_limit} ETB up to 2Cr ETB Equivalent Available Credit {cl.available_credit} Utilized Credit {cl.utilized_credit} Interest Rate {cl.interest_rate}% per Annum Status {cl.status} Credit Score {cl.credit_score} 300-900 Based on TPV Payroll Data etc Approved At {cl.approved_at} Created At {cl.created_at} Disbursements {cl.disbursements} Total Disbursed {cl.total_disbursed} Total Repaid {cl.total_repaid} Outstanding {cl.outstanding} • Outstanding per RazorpayX Capital Line of Credit • Credit Scoring Based on TPV Payroll Data etc • Outstanding • RazorpayX Capital Line of Credit Company Registration Current Accounts Vendor Payments Payroll Partners for Startups & SMEs • Get In-depth Reporting into Cash Flow Trends • Get Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History with Razorpay • Instant Loans Without Collaterals • Collateral-free Credit Lines • Embedded Lending Products • Digital Lending 2.0 • Razorpay Capital Line of Credit</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(parseInt(cl.utilized_credit)/parseInt(cl.credit_limit))*100}%` }} /></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Request Credit Line • Credit Limit up to 2Cr ETB Equivalent • Available Credit Utilized Credit Interest Rate 18% per Annum • Credit Scoring Based on TPV Payroll Data • Outstanding • RazorpayX Capital Line of Credit</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Loan Disbursements • Credit Line ID • Amount • Currency • Purpose Working Capital/Inventory/Payroll etc • Status Pending/Approved/Disbursed/Repaid/Defaulted/Cancelled • Disbursed At • Due Date • Repaid Amount • Outstanding Amount • Ledger Book ID • Created By • Created At • Outstanding • RazorpayX Capital Line of Credit • Instant Loans Without Collaterals</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Credit Line ID • Amount • Currency • Purpose Working Capital/Inventory/Payroll etc • Status Pending/Approved/Disbursed/Repaid/Defaulted/Cancelled</span><span>Disbursed At • Due Date • Repaid Amount • Outstanding Amount • Ledger Book ID • Created By • Created At</span><span>Action • Repay • Default • Cancel • Outstanding • Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History</span></div>
              {mockDisbursements.map(d => (
                <div key={d.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>Credit Line {d.credit_line_id} • Amount {d.amount} {d.currency} • Purpose {d.purpose} • Status {d.status} • Disbursed At {d.disbursed_at} • Due Date {d.due_date} • Repaid {d.repaid_amount} • Outstanding {d.outstanding_amount} • Ledger Book {d.ledger_book_id} • Created {d.created_at}</span>
                  <span>Amount {d.amount} {d.currency} • Purpose {d.purpose} • Status {d.status} • Disbursed At {d.disbursed_at} • Due Date {d.due_date} • Repaid {d.repaid_amount} • Outstanding {d.outstanding_amount} • Ledger Book {d.ledger_book_id} • Created {d.created_at} • Outstanding per RazorpayX Capital Line of Credit • Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History</span>
                  <span className="flex flex-col gap-1"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Repay • Repaid Amount {d.repaid_amount} • Outstanding {d.outstanding_amount} • Ledger Book {d.ledger_book_id} • Created {d.created_at} • Outstanding • Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History</button><button className="rounded-xl border h-7 px-3 text-[10px]">View • Ledger Book {d.ledger_book_id} • Created {d.created_at} • Outstanding • Instant Loans Without Collaterals</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Credit Scoring Based on TPV Payroll Data etc • Outstanding • RazorpayX Capital Line of Credit • Credit Scoring • TCO Win 30-60% Over Stitched Stack • Engineering Cost Saved on Integrations</p>
              <p className="mt-2">Credit Score 300-900 Based on TPV Payroll Data etc • TPV Today 125430 • Success Rate 96.2% • Active Links 12 • Recent Payments • Ledger M4 per run book • Payroll Runs • Cost Center Allocation • Variance Report • Payroll Register • Cost Center Report • Compliance Pension CSV + ERCA CSV + Bank File • Outstanding modern UI glassmorphic • Recharts • Credit Scoring Based on TPV Payroll Data etc • Outstanding • RazorpayX Capital Line of Credit Company Registration Current Accounts Vendor Payments Payroll Partners for Startups & SMEs • Get In-depth Reporting into Cash Flow Trends • Get Instant Loans Without Collaterals • Collateral-free Instant Loan Only After 3 Months Transactions History with Razorpay • Instant Loans Without Collaterals • Collateral-free Credit Lines • Embedded Lending Products • Digital Lending 2.0 • Razorpay Capital Line of Credit • Outstanding modern UI glassmorphic • Recharts • Credit Scoring Based on TPV Payroll Data etc</p>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
