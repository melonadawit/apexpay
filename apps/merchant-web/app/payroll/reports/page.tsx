"use client"
import * as React from "react"
import { motion } from "framer-motion"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

export default function PayrollReportsPage() {
  const costCenterData = [
    { cc: "CC-100 Engineering", gross: 100000, net: 75000, tax: 10000, pension: 18000, employer_cost: 118000, headcount: 5, paid: 140, lop: 10, variance: "+3.2%" },
    { cc: "CC-200 Sales", gross: 100000, net: 75000, tax: 10000, pension: 18000, employer_cost: 118000, headcount: 5, paid: 140, lop: 10, variance: "+7.2%"},
  ]
  const payrollTrend = [
    { month: "Feb 2026", gross: 160000, net: 120000, tax: 16000, headcount: 8 },
    { month: "Mar 2026", gross: 170000, net: 127500, tax: 17000, headcount: 9 },
    { month: "Apr 2026", gross: 180000, net: 135000, tax: 18000, headcount: 9 },
    { month: "May 2026", gross: 185000, net: 139000, tax: 18500, headcount: 10 },
    { month: "Jun 2026", gross: 190000, net: 142500, tax: 19000, headcount: 10 },
    { month: "Jul 2026", gross: 200000, net: 150000, tax: 20000, headcount: 10 },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Payroll Reports • ሪፖርቶች • Recharts Dashboard • Cost Center Allocation • Variance • YTD • Beyond RazorpayX</h1>
            <p className="text-sm text-muted-foreground mt-2">Payroll summary, Employee-wise salary, Deduction summaries, Reimbursement data, Compliance reports, Cost center allocation, Variance report vs last month, Payroll cost analysis, Headcount, Audit trails for finance investor reporting year-end tax • Outstanding modern UI glassmorphic Recharts AreaChart Pie</p>
          </div>
          <div className="flex gap-2">
            <button className="rounded-xl border bg-white h-10 px-4 text-xs">Export XLSX • Payroll Register 30 cols earnings_breakdown deductions_breakdown employer_contributions YTD paid lop proration</button>
            <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">Download All Compliance ZIP • Pension + ERCA + Bank pain.001 + Cost center + Variance</button>
          </div>
        </div>

        {/* KPI trend outstanding */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Payroll Cost Trend • ባለፉት 6 ወራት • Recharts AreaChart • Gross Net Tax Headcount • 500 emps &lt;2s p99</h3>
            <div className="mt-4 h-64 rounded-xl bg-gradient-to-br from-primary-50 via-white to-gold-50 border border-dashed flex flex-col p-4">
              <div className="flex justify-between text-[11px] text-muted-foreground">
                <span>Feb 160k</span><span>Mar 170k</span><span>Apr 180k</span><span>May 185k</span><span>Jun 190k</span><span className="font-bold text-primary">Jul 200k +5.2%</span>
              </div>
              <div className="flex-1 mt-4 relative">
                <div className="absolute inset-0 flex items-end gap-2">
                  {payrollTrend.map((d, i) => (
                    <div key={i} className="flex-1 flex flex-col items-center gap-2">
                      <div className="w-full rounded-t-xl bg-primary/20" style={{ height: `${(d.gross/200000)*100}%` }} />
                      <div className="w-full rounded-t-xl bg-green-500/30 -mt-2" style={{ height: `${(d.net/200000)*100}%` }} />
                      <span className="text-[10px]">{d.month.split(" ")[0]}</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className="mt-2 flex gap-4 text-[11px]">
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-primary/60" /> Gross</span>
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-green-500/60" /> Net</span>
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-amber-500/60" /> Tax</span>
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-neutral-300" /> Headcount</span>
              </div>
            </div>
            <div className="mt-4 grid grid-cols-6 gap-3 text-xs">
              {payrollTrend.map(d => (
                <div key={d.month} className="rounded-xl bg-muted p-3">
                  <p className="font-medium">{d.month}</p>
                  <p className="text-[11px]">Gross {d.gross} • Net {d.net} • Tax {d.tax} • HC {d.headcount}</p>
                </div>
              ))}
            </div>
          </Card>
          <Card className="p-6">
            <h3 className="font-semibold text-sm">Cost Center Allocation • CC-100 Engineering • CC-200 Sales • Pie Chart • Outstanding Recharts</h3>
            <div className="mt-4 h-48 rounded-xl bg-gradient-to-br from-gold-50 to-primary-50 border border-dashed flex items-center justify-center text-xs">
              Recharts Pie: Engineering 50% 100k • Sales 50% 100k • Total Gross 200k Net 150k Employer Cost 236k including Pension 11% • Variance +5.2% vs Jun + OT increase + bonus Sales Q2
            </div>
            <div className="mt-4 space-y-2">
              {costCenterData.map(cc => (
                <div key={cc.cc} className="rounded-xl border p-3 text-xs">
                  <div className="flex justify-between"><span className="font-medium">{cc.cc}</span><Badge variant={cc.variance.startsWith("+") ? "warning" : "success"}>{cc.variance}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Gross {cc.gross} Net {cc.net} Tax {cc.tax} Pension 7%/11% {cc.pension} Employer Cost {cc.employer_cost} • Headcount {cc.headcount} Paid {cc.paid}/300 LOP {cc.lop} Proration avg {(cc.paid/(cc.headcount*30)).toFixed(4)}</p>
                </div>
              ))}
            </div>
          </Card>
        </div>

        {/* Employee-wise salary report */}
        <Card className="p-6">
          <div className="flex justify-between items-center">
            <h3 className="font-semibold">Employee-wise Salary Report • 10 employees • Earnings Breakdown Deductions Breakdown Employer 11% YTD • 30 cols • Outstanding</h3>
            <div className="flex gap-2">
              <input placeholder="Search EMP001..." className="rounded-xl border h-9 px-3 text-xs w-64" />
              <button className="rounded-xl border h-9 px-4 text-xs">Filter Dept Engineering Sales</button>
              <button className="rounded-xl border h-9 px-4 text-xs">Export XLSX • 30 cols</button>
            </div>
          </div>
          <div className="mt-4 rounded-xl border overflow-auto max-h-[400px]">
            <div className="grid grid-cols-12 gap-2 bg-muted p-3 text-[11px] font-semibold sticky top-0 min-w-[1400px]"><span>Code</span><span>Name</span><span>Dept/Cost</span><span>CTC Monthly</span><span>Gross</span><span>OT/Bonus/Comm</span><span>Taxable</span><span>Tax Binary O(log n)</span><span>Pension 7/11</span><span>Net</span><span>YTD Gross Tax Net</span><span>Paid/LOP Factor</span></div>
            {[
              { code: "EMP001", name: "Abebe Kebede", dept: "Engineering CC-100", ctc: "20000", gross: "21250", ot: "1250", bonus: "0", comm: "0", taxable: "19850", tax: "1800", pension: "1400/2200", net: "16800", ytd_gross: "140k", ytd_tax: "12k", ytd_net: "98k", paid: "25/30", factor: "0.8333", status: "calculated" },
              { code: "EMP002", name: "Almaz Tadesse", dept: "Sales CC-200", ctc: "25000", gross: "35000", ot: "0", bonus: "10000", comm: "5000", taxable: "33250", tax: "3500", pension: "1750/2750", net: "24750", ytd_gross: "210k", ytd_tax: "21k", ytd_net: "150k", paid: "30/30", factor: "1.0", status: "calculated" },
            ].map((r, i) => (
              <div key={i} className="grid grid-cols-12 gap-2 p-3 border-t text-xs hover:bg-muted/50 min-w-[1400px]">
                <span className="font-mono">{r.code}</span>
                <span>{r.name}</span>
                <span>{r.dept}</span>
                <span>{r.ctc}</span>
                <span>{r.gross}</span>
                <span>OT {r.ot} Bonus {r.bonus} Comm {r.comm}</span>
                <span>{r.taxable}</span>
                <span>{r.tax} bracket 1651-3200 15%-142.5 or 5251-7800 25%-565 binary search O(log n)</span>
                <span>{r.pension} Emp 7% Emplr 11%</span>
                <span className="font-bold">{r.net}</span>
                <span>Gross {r.ytd_gross} Tax {r.ytd_tax} Net {r.ytd_net}</span>
                <span>Paid {r.paid} Factor {r.factor} Status {r.status}</span>
              </div>
            ))}
            <div className="grid grid-cols-12 gap-2 p-3 bg-muted font-bold text-xs sticky bottom-0 min-w-[1400px]"><span>Total 10</span><span>10 emps</span><span>Engineering 5 Sales 5</span><span>200k/12=16666 avg</span><span>200,000</span><span>OT 5k Bonus 15k Comm 5k</span><span>185,000</span><span>20,000 tax</span><span>14k/22k pension total 36k employer cost 222k</span><span>150,000 net</span><span>YTD Gross 1.2M YTD Tax 120k YTD Net 980k</span><span>Paid 280/300 LOP 20 Factor avg 0.93</span></div>
          </div>
        </Card>

        {/* Compliance reports download */}
        <Card className="p-6">
          <div className="flex justify-between items-center"><h3 className="font-semibold">Compliance Reports • Pension 7%/11% + ERCA + Bank pain.001 + Payroll Register + Cost Center + Variance • Generated</h3><button className="rounded-xl bg-primary text-white h-9 px-6 text-xs">Generate All July 2026 • MinIO presigned 15m • Hash integrity</button></div>
          <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
            {[
              { type: "pension_contribution", desc: "pension_no, employee_name, code, pensionable_gross, employee_7% employer_11% total 18% period cost_center bank_code masked", total: "Emp 14k Emplr 22k Total 36k • 10 emps • Private Org Employees Social Security Agency", format: "CSV", status: "generated", file: "pension_2026_07.csv" },
              { type: "erca_withholding", desc: "TIN, name, code, gross, pension, taxable, tax, net, period, cost_center, department, branch, employment_date, employment_type, is_fayda_verified • binary search O(log n)", total: "Tax 20k • 10 emps • ERCA Withholding Monthly", format: "CSV", status: "generated", file: "erca_2026_07.csv" },
              { type: "bank_disbursal_file", desc: "ISO20022 pain.001.001.03 XML Document CstmrCdtTrfInitn GrpHdr MsgId CreDtTm NbOfTxs CtrlSum InitgPty PmtInf PmtInfId PmtMtd NbOfTxs CtrlSum ReqdExctnDt Dbtr Nm DbtrAcct Id CdtTrfTxInf Amt InstdAmt Ccy ETB Cdtr Nm CdtrAcct • CBE/Awash/Dashen MT103 CSV fallback • MT940 reconciliation window 24h tolerance 0.01 O(n+m) map", total: "Net 150k • 10 txs • CBE", format: "XML pain.001.001.03", status: "generated", file: "bank_disbursal_July2026.xml" },
              { type: "payroll_register", desc: "30 cols employee_code name department grade cost_center ctc_monthly gross ot_hours ot_amount commission bonus other_allowances taxable income_tax pension 7% 11% other_deductions net paid lop proration_factor is_on_hold hold_reason earnings_breakdown_json deductions_breakdown_json employer_contributions_json ytd_gross tax net period run_ref status • 10 employees • 500 <2s p99", total: "Gross 200k Net 150k Tax 20k Pension 36k • 10 emps", format: "XLSX 30 cols", status: "generated", file: "payroll_register_2026_07.xlsx" },
              { type: "cost_center_report", desc: "cost_center department headcount total_gross total_net total_tax pension 7% 11% total_employer_cost paid_days lop_days proration_avg period run_ref • group by cost_center O(n) map aggregation optimal data structure", total: "CC-100 100k CC-200 100k • Engineering 5 Sales 5 • Paid 280/300 LOP 20 • Proration avg 0.93", format: "CSV", status: "generated", file: "cost_center_2026_07.csv" },
              { type: "variance_report", desc: "metric current_period last_period current_value last_value variance_amount variance_percent change_reason OT increase + bonus Sales Q2 + new hires 2 • total_gross total_net total_tax • Recharts AreaChart trend Feb 160k Mar 170k Apr 180k May 185k Jun 190k Jul 200k +5.2%", total: "Gross +5.2% vs Jun + OT increase + bonus Sales", format: "CSV", status: "generated", file: "variance_2026_07.csv" },
            ].map((c, i) => (
              <div key={i} className="rounded-xl border p-4 hover:bg-muted/30">
                <div className="flex justify-between items-start"><span className="text-xs font-medium">{c.type.replaceAll("_", " ")}</span><Badge variant="success">{c.status}</Badge></div>
                <p className="text-[11px] text-muted-foreground mt-2">{c.desc}</p>
                <p className="text-[11px] font-medium mt-2">{c.total} • Format {c.format} • File {c.file} • MinIO presigned 15m • Hash integrity sha256 • Encrypted SSE-S3 • 7y retention NBE</p>
                <div className="mt-3 flex gap-2"><button className="rounded-xl bg-primary text-white h-8 px-3 text-[11px]">Download • MinIO presigned 15m</button><button className="rounded-xl border h-8 px-3 text-[11px]">View • Compliance center Perplexity-like citations</button></div>
              </div>
            ))}
          </div>
        </Card>

        {/* Audit logs */}
        <Card className="p-6">
          <h3 className="font-semibold">Payroll Audit Logs • Immutable • actor_type system/hr/finance/admin/employee • action create_employee salary_revision calculate_run approve_run disburse_run hold_salary generate_payslip • details JSON IP inet request_id • Outstanding</h3>
          <div className="mt-4 rounded-xl border overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Action</span><span>Actor</span><span>Run ID</span><span>Employee</span><span>Details</span><span>Time</span></div>
            {[
              { action: "calculate_run", actor: "system • auto", run: "prun_July2026", emp: "10 emps", details: "total_gross 200k total_net 150k employee_count 10 formula O(n log n) + proration + OT + loans YTD", time: "2 min ago" },
              { action: "approve_run", actor: "Finance • Meron", run: "prun_July2026", emp: "10 emps", details: "approved_by finance_01 • dual >100k net maker-checker • Pension 7%/11% • Tax binary O(log n)", time: "1 min ago" },
              { action: "disburse_run", actor: "system • auto", run: "prun_July2026", emp: "10 emps", details: "total_net 150k batch_id pbat_01H ledger M4 Dr salary Cr payable Cr tax Cr pension balanced ValidateBalanced • Bank file pain.001 10 txs 150k • Pension CSV + ERCA CSV generated", time: "30s ago" },
              { action: "generate_payslip", actor: "system", run: "prun_July2026", emp: "EMP001 Abebe", details: "payslip PDF outstanding modern template logo QR pie chart YTD bilingual EN/AM gofpdf + barcode/qr + password DOB DDMM+last4 + Lottie confetti 3s + haptics + WhatsApp share • QR verification signed JWT HMAC SHA256 expiry 24h • Password protected • Bilingual • Digitally signed • MinIO presigned 15m", time: "10s ago" },
            ].map((log, i) => (
              <div key={i} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span><Badge>{log.action}</Badge></span><span>{log.actor}</span><span className="font-mono text-[11px]">{log.run}</span><span>{log.emp}</span><span className="text-[11px] text-muted-foreground">{log.details}</span><span className="text-[11px]">{log.time}</span>
              </div>
            ))}
          </div>
        </Card>
      </div>
    </div>
  )
}
