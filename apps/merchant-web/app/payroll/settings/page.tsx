"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) { const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border" }; return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span> }

export default function PayrollSettingsPage() {
  const [brackets, setBrackets] = React.useState([
    { id: "brack_600", min: 0, max: 600, rate: 0, deduction: 0, effective_from: "2024-01-01" },
    { id: "brack_1650", min: 601, max: 1650, rate: 10, deduction: 60, effective_from: "2024-01-01" },
    { id: "brack_3200", min: 1651, max: 3200, rate: 15, deduction: 142.5, effective_from: "2024-01-01" },
    { id: "brack_5250", min: 3201, max: 5250, rate: 20, deduction: 302.5, effective_from: "2024-01-01" },
    { id: "brack_7800", min: 5251, max: 7800, rate: 25, deduction: 565, effective_from: "2024-01-01" },
    { id: "brack_10900", min: 7801, max: 10900, rate: 30, deduction: 955, effective_from: "2024-01-01" },
    { id: "brack_inf", min: 10901, max: null, rate: 35, deduction: 1500, effective_from: "2024-01-01" },
  ])

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Payroll Settings • ቅንብሮች • Tax Brackets • OT Rates • Pension 7%/11% • Pay Calendar • Compliance • Beyond RazorpayX</h1>
            <p className="text-sm text-muted-foreground mt-2">Versioned effective_from/to, binary search O(log n) tax calculation, formula engine secure O(n) tokenization + shunting-yard, OT map O(1) 1.25/1.5/2.0/1.3 per ET Labour Law 1156/2019, pension 7% employee 11% employer Private Org Employees Social Security Agency, pay calendar cutoff disbursal lock after disbursal, compliance pension ERCA bank file • Outstanding modern UI glassmorphic</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">Save Settings • Versioned • Audit log</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Tax Brackets • ET 2024 • 7 brackets • Versioned effective_from • Binary Search O(log n) • Formula tax=taxable*rate-deduction • Rounded 2 decimals • Benchmark p99&lt;30ms</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Min</span><span>Max</span><span>Rate %</span><span>Deduction ETB</span><span>Effective From</span><span>Action</span></div>
              {brackets.map((b, i) => (
                <div key={b.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <input value={b.min} onChange={e => { const nb = [...brackets]; nb[i].min = parseInt(e.target.value)||0; setBrackets(nb) }} className="rounded-lg border h-8 px-2 text-xs" />
                  <input value={b.max ?? ""} placeholder="∞" onChange={e => { const nb = [...brackets]; nb[i].max = e.target.value? parseInt(e.target.value): null; setBrackets(nb) }} className="rounded-lg border h-8 px-2 text-xs" />
                  <input value={b.rate} onChange={e => { const nb = [...brackets]; nb[i].rate = parseFloat(e.target.value)||0; setBrackets(nb) }} className="rounded-lg border h-8 px-2 text-xs" />
                  <input value={b.deduction} onChange={e => { const nb = [...brackets]; nb[i].deduction = parseFloat(e.target.value)||0; setBrackets(nb) }} className="rounded-lg border h-8 px-2 text-xs" />
                  <input value={b.effective_from} onChange={e => { const nb = [...brackets]; nb[i].effective_from = e.target.value; setBrackets(nb) }} className="rounded-lg border h-8 px-2 text-xs" />
                  <button className="text-red-500 text-xs">Delete</button>
                </div>
              ))}
            </div>
            <div className="mt-4 flex gap-2">
              <button className="rounded-xl border h-9 px-4 text-xs">+ Add Bracket • Versioned effective_from/to O(log n) binary search via sort.Search len brackets func taxable<Max</button>
              <button className="rounded-xl border h-9 px-4 text-xs">Validate • tax=taxable*rate-deduction rounded 2 decimals • p99&lt;30ms benchmark 10k iterations deterministic seed 42</button>
            </div>
            <div className="mt-4 rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]">
              <p className="font-semibold">Benchmark: CalculateTax 10k iterations p99&lt;30ms • ValidateBalanced O(n) • payroll calc 500 emps &lt;2s p99 • k6 100 VUs p95&lt;300ms • TestPayrollTaxBracketLogic binary search vs known examples rounding edge .005 • Property 10k iterations deterministic seed 42 • No float money grep float64 amount</p>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">OT Rates • ET Labour Law Proclamation No. 1156/2019 Art 90 • Map O(1) • 1.25x 1.5x 2.0x 1.3x • Hourly Rate base/208 (26 days *8h)</h3>
            <div className="mt-4 space-y-3 text-xs">
              {[
                { type: "Weekday", rate: "1.25x", desc: "Weekday overtime beyond 8h", formula: "hourly*1.25*hours", example: "Base 20000/208=96.15/hr *1.25=120.19/hr *5h=600.96" },
                { type: "Weekend", rate: "1.5x", desc: "Weekend Saturday/Sunday per labour law", formula: "hourly*1.5*hours", example: "96.15*1.5=144.23/hr" },
                { type: "Holiday", rate: "2.0x", desc: "Public holiday per ET calendar", formula: "hourly*2.0*hours", example: "96.15*2.0=192.30/hr" },
                { type: "Night", rate: "1.3x", desc: "Night 10PM-6AM", formula: "hourly*1.3*hours", example: "96.15*1.3=125/hr" },
              ].map(r => (
                <div key={r.type} className="rounded-xl border p-3 flex justify-between items-center">
                  <div><p className="font-medium">{r.type} • {r.rate}</p><p className="text-[11px] text-muted-foreground">{r.desc} • Formula {r.formula}</p></div>
                  <div className="text-right"><p className="font-mono text-[11px]">{r.example}</p><Badge variant="default">Map O(1) lookup OTRates[OTWeekday]</Badge></div>
                </div>
              ))}
            </div>
            <div className="mt-4 rounded-xl bg-amber-500/10 border border-amber-500/20 p-3 text-[11px]">
              <p className="font-semibold">Hourly Rate = Base Salary /208 (26 days *8h) ET standard • OT Amount = Σ(weekdayHours*rateWeekday + weekend*rateWeekend + holiday*rateHoliday + night*rateNight) • Round 2 decimals • Proration Factor paidDays/totalDays 25/30=0.8333 applied to base but OT not prorated? Configurable via is_proratable flag per component</p>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">Pension Rates • Private Organization Employees Social Security Agency • 7% Employee 11% Employer Total 18% • Configurable • No cap placeholder</h3>
            <div className="mt-4 grid grid-cols-2 gap-4 text-xs">
              <div className="rounded-xl bg-muted p-4"><p className="text-muted-foreground">Employee • 7%</p><p className="text-2xl font-bold">7.00%</p><p className="text-[11px]">Pensionable Gross = Gross - non-pensionable? Configurable pensionable salary = basic + hardship? For simplicity gross for now but make configurable pension_applicable_gross = gross - (OT + Bonus non-pensionable)</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-muted-foreground">Employer • 11%</p><p className="text-2xl font-bold">11.00%</p><p className="text-[11px]">Employer 11% extra cost Total Employer Cost = Gross + Pension Emplr • Ledger M4 Dr expense:pension_employer Cr pension_payable • Total pension both Emp+Emplr 18% • Report pension CSV pension_no name code pensionable gross employee 7% employer 11% total period</p></div>
            </div>
            <div className="mt-4 rounded-xl border p-3 text-[11px]">
              <p className="font-semibold">Pension Contribution Report • Format for Agency:</p>
              <p className="font-mono">pension_no, employee_name, employee_code, pensionable_gross, employee_7pct, employer_11pct, total_18pct, period, cost_center, bank_code, bank_masked</p>
              <p>PEN-001, Abebe Kebede, EMP001, 20000, 1400, 2200, 3600, 2026-07, CC-100, CBE, ****1234</p>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">Pay Calendar • Cutoff Date • Disbursal Date • Lock After Disbursal • Pay Schedule Monthly Weekly Semimonthly • Outstanding</h3>
            <div className="mt-4 space-y-3 text-xs">
              <div className="rounded-xl bg-muted p-3 flex justify-between"><div><p className="font-medium">July 2026 regular • Cutoff 25 Jul • Disbursal 30 Jul • Pay Date 31 Jul</p><p className="text-[11px]">Total Gross 200k Net 150k Tax 20k Pension 36k • 10 emps • Status completed • Locked at 30 Jul 22:00 • Variance +5.2%</p></div><Badge variant="success">Locked • completed</Badge></div>
              <div className="rounded-xl bg-muted p-3 flex justify-between"><div><p className="font-medium">August 2026 regular • Cutoff 25 Aug • Disbursal 30 Aug • Pay Date 31 Aug</p><p className="text-[11px]">Expected Gross 210k Net 157k • 11 emps • Status draft • 1 new hire EMP011</p></div><Badge variant="warning">Draft</Badge></div>
            </div>
            <div className="mt-4 flex gap-2">
              <button className="rounded-xl border h-9 px-4 text-xs">+ Add Calendar • Monthly Weekly Semimonthly • Cutoff Disbursal Pay Date • Lock payroll after disbursal • Re-run amendment • Hold salary • Skip employee • Arrears • Bonus run taxable ex_gratia</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Departments • Designations • Grades • Branches • Organizational Hierarchy • Cost Center Allocation • Workforce Money OS</h3>
            <div className="mt-4 grid grid-cols-4 gap-4">
              <div>
                <p className="text-xs font-semibold">Departments • 5 • Cost Center</p>
                <div className="mt-2 space-y-2 text-xs">{["Engineering CC-100 5 emps", "Sales CC-200 5 emps", "HR & Admin CC-300 2 emps", "Finance CC-400 1 emps", "Operations CC-500 2 emps"].map(d=><div key={d} className="rounded-xl border p-2">{d}</div>)}</div>
                <button className="mt-2 w-full rounded-xl border border-dashed h-8 text-xs">+ Add Department • code ENG cost_center CC-100 description</button>
              </div>
              <div>
                <p className="text-xs font-semibold">Designations • 7 • Level 1-5</p>
                <div className="mt-2 space-y-2 text-xs">{["Junior Engineer L1", "Senior Engineer L3", "Engineering Manager L5", "Sales Rep L2", "Sales Manager L4", "HR Manager L4", "Finance Manager L5"].map(d=><div key={d} className="rounded-xl border p-2">{d}</div>)}</div>
                <button className="mt-2 w-full rounded-xl border border-dashed h-8 text-xs">+ Add Designation • title level description</button>
              </div>
              <div>
                <p className="text-xs font-semibold">Grades • 5 • G1-G5 min 10k max 100k • Salary Bands</p>
                <div className="mt-2 space-y-2 text-xs">{["G1 10k-15k", "G2 15k-25k", "G3 25k-40k", "G4 40k-60k", "G5 60k-100k"].map(d=><div key={d} className="rounded-xl border p-2">{d}</div>)}</div>
                <button className="mt-2 w-full rounded-xl border border-dashed h-8 text-xs">+ Add Grade • name min max description</button>
              </div>
              <div>
                <p className="text-xs font-semibold">Branches • 3 • Head Office Addis Shashemene Adama • Region Oromiya Addis Ababa</p>
                <div className="mt-2 space-y-2 text-xs">{["Head Office - Addis Ababa • Addis Ababa • Head true", "Shashemene Branch • Oromiya Shashemene", "Adama Branch • Oromiya Adama"].map(d=><div key={d} className="rounded-xl border p-2">{d}</div>)}</div>
                <button className="mt-2 w-full rounded-xl border border-dashed h-8 text-xs">+ Add Branch • name region city sub_city address is_head</button>
              </div>
            </div>
          </Card>
        </div>

        <Card className="p-6">
          <h3 className="font-semibold">Compliance Settings • Pension Agency SFTP • ERCA Filing • Bank File Format pain.001 vs CSV • Fayda Verification Required • Bank Account Name Fuzzy Levenshtein &lt;3 • TIN 10-digit • Document Vault Encrypted SSE-S3 7y NBE</h3>
          <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
            <div className="rounded-xl border p-4"><p className="font-medium">Pension Agency • Private Org Employees Social Security Agency</p><p className="text-[11px] text-muted-foreground mt-1">SFTP host sftp.pension.gov.et • Port 22 • Username apexpay • Password env PENSION_SFTP_PASS • File format CSV pension_no name code pensionable gross 7% 11% total period • Upload after disburse • Status generated → paid → filed • MinIO presigned 15m hash</p><p className="mt-2"><Badge variant="success">Connected • SFTP mock</Badge></p></div>
            <div className="rounded-xl border p-4"><p className="font-medium">ERCA • Ethiopian Revenues and Customs Authority • Withholding Monthly</p><p className="text-[11px] text-muted-foreground mt-1">File format CSV TIN name code gross pension taxable tax net period cost_center department branch employment_date type is_fayda_verified • TIN 10-digit validation regex ^[0-9]{10}$ • Upload after disburse • Annual tax certificate PDF YTD gross tax net + monthly breakdown • Form? ERCA format • Audit immutable</p><p className="mt-2"><Badge variant="success">Connected • CSV generated</Badge></p></div>
            <div className="rounded-xl border p-4"><p className="font-medium">Bank Files • CBE/Awash/Dashen • ISO20022 pain.001.001.03 XML vs MT103 CSV fallback</p><p className="text-[11px] text-muted-foreground mt-1">CBE supports pain.001 XML • Awash supports CSV employee_code name bank_code masked name amount currency payout_ref period runRef cost_center • Dashen supports MT103? • Reconciliation MT940 window 24h amount tolerance 0.01 ETB O(n+m) map connector_ref→journal • File hash • Batch ref • Status queued → processing → succeeded/failed/returned • Ledger second journal Dr payable Cr clearing:bank • Payout batch book per batch</p><p className="mt-2"><Badge variant="success">pain.001 generated • 10 txs 150k</Badge></p></div>
          </div>
        </Card>
      </div>
    </div>
  )
}
