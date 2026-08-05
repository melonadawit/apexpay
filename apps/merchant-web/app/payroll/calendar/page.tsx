"use client"
import * as React from "react"
import { motion } from "framer-motion"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100 text-neutral-700", success: "bg-green-500/15 text-green-700 border border-green-500/20", warning: "bg-amber-500/15 text-amber-700 border border-amber-500/20", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2.5 py-0.5 rounded-full text-[11px] font-medium border ${map[variant]}`}>{children}</span>
}

const mockCalendars = [
  { id: "cal_2026_07", name: "Monthly Payroll Calendar 2026-07", year: 2026, month: 7, pay_frequency: "monthly", cutoff_day: 25, disbursal_day: 30, pay_day: 31, cutoff_date: "2026-07-25", disbursal_date: "2026-07-30", pay_date: "2026-07-31", is_locked: true, locked_at: "2026-07-30T22:00:00Z", locked_by: "Finance Manager", total_gross: "200000", total_net: "150000", variance: "+5.2% vs Jun", status: "locked", runs: 1 },
  { id: "cal_2026_08", name: "Monthly Payroll Calendar 2026-08", year: 2026, month: 8, pay_frequency: "monthly", cutoff_day: 25, disbursal_day: 30, pay_day: 31, cutoff_date: "2026-08-25", disbursal_date: "2026-08-30", pay_date: "2026-08-31", is_locked: false, total_gross: "210000", total_net: "157500", variance: "+5% vs Jul", status: "draft", runs: 0 },
  { id: "cal_2026_09", name: "Monthly Payroll Calendar 2026-09", year: 2026, month: 9, pay_frequency: "monthly", cutoff_day: 25, disbursal_day: 30, pay_day: 30, cutoff_date: "2026-09-25", disbursal_date: "2026-09-30", pay_date: "2026-09-30", is_locked: false, status: "upcoming" },
]

export default function PayrollCalendarPage() {
  const [selected, setSelected] = React.useState(mockCalendars[0])
  const [showCreate, setShowCreate] = React.useState(false)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-3">Payroll Calendar • የደሞዝ ቀን መቁጠሪያ • Cutoff 25th Disbursal 30th Pay Last Day Lock After Disbursal • Ethiopia Business Practice</h1>
            <p className="text-sm text-muted-foreground mt-2">Per Ethiopia business practice monthly payroll cutoff 25th disbursal 30th pay date last day of month lock after disbursal per law, pay_frequency monthly/semimonthly/weekly/biweekly, year month cutoff_day disbursal_day pay_day cutoff_date disbursal_date pay_date is_locked locked_at locked_by • Outstanding modern UI glassmorphic Recharts</p>
          </div>
          <div className="flex gap-2">
            <button onClick={()=>setShowCreate(true)} className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Calendar • Monthly Weekly Semimonthly</button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Pay Calendars • 2026 • Ethiopia Business Practice • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockCalendars.map(cal => (
                <button key={cal.id} onClick={()=>setSelected(cal)} className={`w-full text-left rounded-xl border p-4 hover:bg-muted transition-colors ${selected.id===cal.id ? "bg-primary/10 border-primary/30 shadow-soft" : ""}`}>
                  <div className="flex justify-between items-start"><p className="font-medium text-sm">{cal.name}</p><Badge variant={cal.is_locked ? "success" : cal.status==="draft" ? "warning" : "default"}>{cal.is_locked ? "Locked • completed" : cal.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Year {cal.year} Month {cal.month} • Frequency {cal.pay_frequency} • Cutoff {cal.cutoff_day}th → {cal.cutoff_date} • Disbursal {cal.disbursal_day}th → {cal.disbursal_date} • Pay {cal.pay_day}th → {cal.pay_date} • Locked at {cal.locked_at || "Not locked"} by {cal.locked_by || "—"}</p>
                  <p className="text-[11px] mt-1">Total Gross {cal.total_gross || "—"} Net {cal.total_net || "—"} Variance {cal.variance || "—"} Runs {cal.runs ?? 0} • Cost center allocation CC-100 Engineering CC-200 Sales • Variance +5.2% vs Jun Recharts</p>
                </button>
              ))}
              <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Add Calendar • Monthly Weekly Semimonthly • Cutoff Disbursal Pay Date • Lock payroll after disbursal • Re-run amendment • Hold salary • Skip employee • Arrears • Bonus run taxable ex_gratia</button>
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Payroll Calendar Logic Outstanding per Ethiopia Law:</p>
              <ul className="list-disc list-inside mt-2 space-y-1 text-muted-foreground">
                <li>Monthly: cutoff 25th disbursal 30th pay date last day last day of month lock after disbursal per law</li>
                <li>Semimonthly: cutoff 15th & last day, disbursal 16th & 1st next month</li>
                <li>Weekly: cutoff Friday disbursal Monday pay Monday</li>
                <li>Biweekly: cutoff every 2 weeks Friday</li>
                <li>Lock after disbursal: is_locked true locked_at now locked_by finance manager • Prevents re-run amendment unless unlocked by admin with audit log payroll_audit_logs actor admin action unlock_calendar details locked_by IP inet request_id immutable</li>
                <li>Variance report vs last month +5.2% vs Jun OT increase + bonus Sales Q2 + new hires 2 • total_gross total_net total_tax • Recharts AreaChart trend Feb 160k Mar 170k Apr 180k May 185k Jun 190k Jul 200k +5.2% • Cost center breakdown Engineering 100k Sales 100k • Paid 280/300 LOP 20 • Proration avg 0.93</li>
              </ul>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <div className="flex justify-between items-center">
              <h3 className="font-semibold">Calendar Detail • {selected.name} • {selected.year}-{String(selected.month).padStart(2,"0")} • {selected.pay_frequency} • Cutoff 25th Disbursal 30th Pay Last Day Lock After Disbursal • Outstanding</h3>
              <div className="flex gap-2">
                <button className="rounded-xl border h-8 px-3 text-[11px]">Edit • Cutoff Disbursal Pay Date</button>
                <button className="rounded-xl border h-8 px-3 text-[11px]">Unlock • Audit log admin</button>
                <button className={`rounded-xl h-8 px-3 text-[11px] ${selected.is_locked ? "bg-neutral-200 text-neutral-500" : "bg-primary text-white"}`} disabled={selected.is_locked}>{selected.is_locked ? "Locked • completed at " + selected.locked_at : "Lock After Disbursal • is_locked true locked_at now locked_by"}</button>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-3 gap-4">
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Cutoff • መቁረጥ • 25th</p><p className="font-bold text-lg">{selected.cutoff_date}</p><p className="text-[11px]">Cutoff day {selected.cutoff_day}th • Attendance final • LOP proration • OT hours final • Variable pay final • No more attendance after cutoff • Payroll run can be created after cutoff • Before cutoff payroll run draft cannot calculate</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Disbursal • ክፍያ • 30th</p><p className="font-bold text-lg">{selected.disbursal_date}</p><p className="text-[11px]">Disbursal day {selected.disbursal_day}th • Payout batch pain.001 XML ISO20022 Document CstmrCdtTrfInitn + Pension CSV + ERCA CSV + Cost center report + Bank file • Ledger second journal Dr payroll_payable Cr clearing:bank • Status processing → completed after payout success worker</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Pay Day • የክፍያ ቀን • Last day</p><p className="font-bold text-lg">{selected.pay_date}</p><p className="text-[11px]">Pay day {selected.pay_day}th • Last day of month • Salary credited to employee bank account CBE/Awash/Dashen • Payslip PDF outstanding modern template logo QR pie chart YTD bilingual EN/AM • Email SMTP distribution • WhatsApp share • Lottie confetti 3s</p></div>
            </div>

            <div className="mt-6">
              <h4 className="font-semibold text-sm">Visual Stepper Pipeline • Outstanding Modern • Framer Motion pathLength</h4>
              <div className="mt-4 flex items-center gap-2 text-[11px]">
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">1</span><span>Cutoff 25th • Attendance final • LOP proration • OT final</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">2</span><span>Calculate Run • V2 formula engine O(n log n) + proration + OT + loans YTD</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-amber-500 text-white flex items-center justify-center">3</span><span>Pending Approval • Maker-checker dual &gt;100k net • Approver != submitter</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center">4</span><span>Disbursal 30th • Payout batch pain.001 XML • Pension CSV • ERCA CSV • Ledger M4</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className={`flex items-center gap-2 ${selected.is_locked ? "" : "opacity-50"}`}><span className={`h-6 w-6 rounded-full ${selected.is_locked ? "bg-green-600" : "bg-neutral-200"} text-white flex items-center justify-center`}>5</span><span>Locked • Pay Date Last Day • Salary credited • Payslip PDF QR • Email SMTP • Lock after disbursal is_locked true locked_at now locked_by • Prevents re-run amendment unless unlocked by admin with audit log</span></div>
              </div>
            </div>

            {showCreate && (
              <motion.div initial={{ opacity:0, y:10 }} animate={{ opacity:1, y:0 }} className="mt-6 rounded-xl border-2 border-primary/30 bg-primary/5 p-4">
                <h4 className="font-semibold text-sm">Create New Calendar • Monthly Weekly Semimonthly Biweekly • Outstanding Form</h4>
                <div className="mt-3 grid grid-cols-4 gap-3 text-xs">
                  <div><label className="text-muted-foreground">Name</label><input placeholder="Monthly Payroll Calendar 2026-10" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div><label className="text-muted-foreground">Pay Frequency • monthly/semimonthly/weekly/biweekly</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>monthly</option><option>semimonthly</option><option>weekly</option><option>biweekly</option></select></div>
                  <div><label className="text-muted-foreground">Year • 2026</label><input type="number" defaultValue={2026} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div><label className="text-muted-foreground">Month • 1-12 • null for weekly</label><input type="number" defaultValue={10} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div><label className="text-muted-foreground">Cutoff Day • 25th Ethiopia business practice</label><input type="number" defaultValue={25} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div><label className="text-muted-foreground">Disbursal Day • 30th</label><input type="number" defaultValue={30} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div><label className="text-muted-foreground">Pay Day • Last day • 31th</label><input type="number" defaultValue={31} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                  <div className="flex items-end gap-2"><button className="rounded-xl bg-primary text-white h-9 px-6">Create Calendar • cutoff_date disbursal_date pay_date auto-calc from cutoff_day disbursal_day pay_day year month • is_locked false • O(1) advisory lock</button><button onClick={()=>setShowCreate(false)} className="rounded-xl border h-9 px-4">Cancel</button></div>
                </div>
                <p className="mt-3 text-[11px] text-muted-foreground">Logic: cutoff_date = year-month-cutoff_day • disbursal_date = year-month-disbursal_day • pay_date = year-month-pay_day (if pay_day 31 and month Feb 28, then last day of month) • is_locked false initially • locked_at null • locked_by null • created_by current user • When payroll run disburse, set is_locked true locked_at now locked_by finance manager • Prevent re-run amendment unless unlocked by admin with audit log payroll_audit_logs actor admin action unlock_calendar details locked_by IP inet request_id immutable • Outstanding modern UI glassmorphic Recharts calendar view • Cost center allocation • Variance report</p>
              </motion.div>
            )}

            <div className="mt-6 rounded-xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Payroll Runs in This Calendar • 1 run • Ledger M4 per run book • 500 emps &lt;2s p99</h4>
              <div className="mt-3 rounded-xl border overflow-hidden">
                <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Run Ref</span><span>Period</span><span>Type</span><span>Status</span><span>Total Net</span><span>Action</span></div>
                <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>July2026_Regular</span><span>07/2026</span><span>regular</span><span><Badge variant="success">completed • Locked</Badge></span><span>ETB 150,000</span><span>View • Ledger M4 Dr salary 200k + Dr pension emplr 22k Cr payable 150k Cr tax 20k Cr pension 36k balanced</span></div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
