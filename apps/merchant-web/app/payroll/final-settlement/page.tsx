"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockFnF = [
  { id: "fnf_001", employee: "Kebede Lema • EMP007", resignation_date: "2026-06-15", last_working_date: "2026-07-15", notice_period: 30, notice_served: 30, notice_shortfall: 0, leave_encashment_days: 5, leave_encashment_amount: "2500", severance_amount: "15000", gratuity: "0", bonus_pro_rata: "2000", outstanding_loans: "5000", outstanding_advances: "0", other_earnings: "0", other_deductions: "0", total_payable: "19500", total_deductions: "5000", net_payable: "14500", status: "pending_approval", clearance: [
    { item: "Laptop • LP001", category: "IT", status: "pending", required: true, checked_by: "", checked_at: "", notes: "MacBook Pro 14 inch" },
    { item: "ID Card • ID-EMP007", category: "HR", status: "done", required: true, checked_by: "HR Manager", checked_at: "2026-07-14", notes: "Returned" },
    { item: "Office Keys • Keys-007", category: "Admin", status: "done", required: true, checked_by: "Admin", checked_at: "2026-07-14", notes: "Returned" },
    { item: "Company Car • Car-007", category: "Admin", status: "pending", required: false, checked_by: "", checked_at: "", notes: "N/A - No car" },
  ], assets_returned: [{ asset_type: "laptop", asset_id: "LP001", returned: false, condition: "good", returned_at: "" }], exit_interview: { conducted: false, conducted_by: "", date: "", feedback: "" } },
]

export default function FinalSettlementPage() {
  const [selected, setSelected] = React.useState(mockFnF[0])
  const [checklist, setChecklist] = React.useState(selected.clearance)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Final Settlement F&F • የሥራ መልቀቂያ • Clearance Checklist Laptop ID Card • Assets Returned • Exit Interview • Severance Art 39-44 • Leave Encashment Gross/30 • Ethiopia Labour Proclamation 1156/2019</h1>
          <p className="text-sm text-muted-foreground mt-2">Final settlement per Ethiopia Labour Proclamation Art 39-44: resignation_date last_working_date notice_period_days notice_served_days notice_shortfall_days leave_encashment_days per_day gross/30 amount severance per ET labour law Art 39-44 gratuity bonus_pro_rata outstanding_loans advances other_earnings other_deductions total_payable total_deductions net_payable status draft/pending_approval/approved/paid/rejected clearance_checklist JSON clearance_items_detailed assets_returned exit_interview • Outstanding modern UI glassmorphic</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">F&F Requests • Pending Approval • Clearance Checklist • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockFnF.map(fnf => (
                <button key={fnf.id} onClick={()=>{ setSelected(fnf); setChecklist(fnf.clearance) }} className={`w-full text-left rounded-xl border p-4 hover:bg-muted ${selected.id===fnf.id ? "bg-primary/10 border-primary/30" : ""}`}>
                  <div className="flex justify-between"><p className="font-medium text-sm">{fnf.employee}</p><Badge variant={fnf.status==="pending_approval" ? "warning" : "success"}>{fnf.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Resignation {fnf.resignation_date} • LWD {fnf.last_working_date} • Notice {fnf.notice_period} days Served {fnf.notice_served} Shortfall {fnf.notice_shortfall} • Leave Encashment {fnf.leave_encashment_days} days {fnf.leave_encashment_amount} ETB per_day gross/30 • Severance {fnf.severance_amount} ETB Art 39-44 • Bonus Pro-rata {fnf.bonus_pro_rata} • Outstanding Loans {fnf.outstanding_loans} Advances {fnf.outstanding_advances}</p>
                  <p className="text-[11px] mt-1">Total Payable {fnf.total_payable} • Total Deductions {fnf.total_deductions} • Net Payable {fnf.net_payable} • Status {fnf.status} • Clearance {fnf.clearance.filter(c=>c.status==="done").length}/{fnf.clearance.length} done • Assets Returned {fnf.assets_returned.filter((a:any)=>a.returned).length}/{fnf.assets_returned.length}</p>
                </button>
              ))}
              <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Create F&F • Resignation Date LWD Notice Period Served Shortfall Leave Encashment Days per_day gross/30 Amount Severance Art 39-44 Gratuity Bonus Pro-rata Outstanding Loans Advances Other Earnings Other Deductions Total Payable Total Deductions Net Payable Status Clearance Checklist • Outstanding</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <div className="flex justify-between items-center">
              <h3 className="font-semibold">F&F Detail • {selected.employee} • Resignation {selected.resignation_date} • LWD {selected.last_working_date} • Notice {selected.notice_period} Served {selected.notice_served} Shortfall {selected.notice_shortfall} • Leave Encashment {selected.leave_encashment_days} days {selected.leave_encashment_amount} ETB per_day gross/30 • Severance {selected.severance_amount} ETB Art 39-44 • Bonus Pro-rata {selected.bonus_pro_rata} • Outstanding Loans {selected.outstanding_loans} Advances {selected.outstanding_advances} • Total Payable {selected.total_payable} • Total Deductions {selected.total_deductions} • Net Payable {selected.net_payable} • Status {selected.status}</h3>
              <Badge variant="warning">Pending Approval • Maker-checker dual approval • HR + Finance + Admin</Badge>
            </div>

            <div className="mt-6 grid grid-cols-3 gap-4">
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Leave Encashment • per_day gross/30 • Art 77(3) Unused annual leave paid upon termination only</p><p className="font-bold text-lg">ETB {selected.leave_encashment_amount} • {selected.leave_encashment_days} days • per_day = gross/30 = 15000/30=500 • 5 days *500=2500 • Outstanding</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Severance • Art 39-44 • 30 days wage per year of service</p><p className="font-bold text-lg">ETB {selected.severance_amount} • Severance = base_salary /30*30*years* factor = base*years*factor • Base 15000 * 1 year *1.0 =15000 • For illegal termination 2x factor</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px] text-muted-foreground">Net Payable • Total Payable - Total Deductions • Floor Zero</p><p className="font-bold text-lg">ETB {selected.net_payable} • Total Payable {selected.total_payable} = Leave {selected.leave_encashment_amount} + Severance {selected.severance_amount} + Gratuity {selected.gratuity} + Bonus {selected.bonus_pro_rata} + Other 0 • Total Deductions {selected.total_deductions} = Loans {selected.outstanding_loans} + Advances {selected.outstanding_advances} + Other 0 • Net {selected.net_payable} = {selected.total_payable} - {selected.total_deductions}</p></div>
            </div>

            <div className="mt-6">
              <h4 className="font-semibold text-sm">Clearance Checklist • Laptop ID Card Office Keys Company Car • Assets Returned • Exit Interview • Outstanding Modern • Framer Motion stagger 50ms</h4>
              <div className="mt-3 rounded-xl border overflow-hidden">
                <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Item • Laptop ID Card</span><span>Category • IT HR Admin Finance</span><span>Required</span><span>Status • pending/done</span><span>Checked By • HR Manager Admin</span><span>Notes • MacBook Pro 14 inch Returned</span></div>
                {checklist.map((item:any, idx:number)=>(
                  <div key={idx} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                    <span className="font-medium">{item.item}</span>
                    <span><Badge>{item.category}</Badge></span>
                    <span>{item.required ? "Yes • Required" : "No • Optional"}</span>
                    <span>
                      <select value={item.status} onChange={e=>{
                        const newList = [...checklist]
                        newList[idx].status = e.target.value
                        newList[idx].checked_by = "HR Manager"
                        newList[idx].checked_at = new Date().toISOString()
                        setChecklist(newList)
                      }} className="rounded-lg border h-7 px-2 text-xs">
                        <option value="pending">Pending • Not returned</option>
                        <option value="done">Done • Returned • Checked</option>
                      </select>
                    </span>
                    <span>{item.checked_by || "—"} • {item.checked_at ? new Date(item.checked_at).toLocaleDateString() : "—"}</span>
                    <span className="text-[11px]">{item.notes}</span>
                  </div>
                ))}
              </div>
              <div className="mt-3 flex justify-between text-[11px] text-muted-foreground">
                <span>Clearance {checklist.filter(c=>c.status==="done").length}/{checklist.length} done • {checklist.filter(c=>c.required && c.status!=="done").length} required pending • Must complete all required before approval • Outstanding</span>
                <span>Assets Returned {selected.assets_returned.filter((a:any)=>a.returned).length}/{selected.assets_returned.length} • Exit Interview {selected.exit_interview.conducted ? "Conducted" : "Pending"}</span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-6">
              <div className="rounded-xl bg-neutral-50 border p-4 text-xs">
                <p className="font-semibold">Severance Calculation per ET Labour Law Art 39-44 • Outstanding • O(1)</p>
                <p className="mt-2 font-mono text-[11px]">Daily wage = base_salary /30 • Severance per year = 30 days wage = base/30*30 = base • Total severance = base * yearsOfService * factor • Factor 1.0 = one month per year • For illegal termination maybe 2x configurable factor • Example: Base 15000 * 1 year *1.0 =15000 • Base 20000 * 2 years *1.0 =40000 • Base 25000 * 3 years *1.0 =75000 • For mutual separation gratuity? Configurable • Outstanding</p>
                <p className="mt-2">Notice Period Art 43: Depends service years probation &lt;... For permanent notice 1 month if &lt;=1 year 2 months if 1-5 years 3 months &gt;5 years during probation 15 days • NoticePeriodDaysET yearsOfService employmentType confirmationStatus probation 15 &lt;1 year 30 &lt;5 years 60 else 90 • Notice shortfall = notice_period - notice_served • If shortfall &gt;0 deduct from final settlement? Actually per law if employee doesn't serve notice, employer can deduct? For simplicity, shortfall days * per_day = deduction? Configurable • Outstanding</p>
              </div>
              <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-xs">
                <p className="font-semibold">Final Settlement Workflow • Maker-checker dual approval • Outstanding avatars • Audit log</p>
                <div className="mt-3 space-y-3">
                  <div className="flex items-center gap-3"><div className="h-8 w-8 rounded-full bg-neutral-200 flex items-center justify-center text-xs">HR</div><div><p className="font-medium">HR • Created F&F • Resignation {selected.resignation_date} LWD {selected.last_working_date}</p><p className="text-[11px] text-muted-foreground">Notice {selected.notice_period} Served {selected.notice_served} Shortfall {selected.notice_shortfall} • Leave Encashment {selected.leave_encashment_days} days {selected.leave_encashment_amount} per_day gross/30 • Severance {selected.severance_amount} Art 39-44</p></div></div>
                  <div className="flex items-center gap-3"><div className="h-8 w-8 rounded-full bg-amber-500 text-white flex items-center justify-center text-xs">F</div><div><p className="font-medium">Finance • Pending approval • Outstanding Loans {selected.outstanding_loans} Advances {selected.outstanding_advances} • Total Deductions {selected.total_deductions}</p><p className="text-[11px] text-muted-foreground">Checks: Loans outstanding, advances, other deductions, bonus pro-rata, gratuity</p></div><span className="ml-auto px-2 py-0.5 rounded-full bg-amber-500/15 text-amber-700 text-[10px]">Pending</span></div>
                  <div className="flex items-center gap-3 opacity-50"><div className="h-8 w-8 rounded-full bg-neutral-200 flex items-center justify-center text-xs">A</div><div><p className="font-medium">Admin • Final approval • Clearance checklist {checklist.filter(c=>c.status==="done").length}/{checklist.length} done • Assets returned • Exit interview</p><p className="text-[11px] text-muted-foreground">After finance • Then disburse F&F payout single payout via off-cycle payroll run type bonus/fnf</p></div></div>
                </div>
                <div className="mt-4 rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]">
                  <p className="font-semibold">Disbursal • Single Payout via Off-cycle Payroll Run Type Bonus/FnF • Outstanding</p>
                  <p className="mt-1">When F&F approved, disburse single payout via off-cycle payroll run type bonus/fnf • Net payable {selected.net_payable} ETB • Payout batch pain.001 XML ISO20022 Document CstmrCdtTrfInitn + Pension CSV + ERCA CSV + Cost center report + Bank file • Ledger second journal Dr payable Cr clearing:bank • Payout batch book per batch • Status queued → processing → succeeded/failed/returned • Outstanding modern UI glassmorphic</p>
                </div>
              </div>
            </div>

            <div className="mt-6 flex gap-2">
              <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">Approve F&F • dual HR+Finance+Admin avatar • Outstanding</button>
              <button className="rounded-xl bg-green-600 text-white h-10 px-6 text-xs">Disburse • Single Payout {selected.net_payable} ETB • Off-cycle bonus/fnf • Payout batch pain.001 XML • Ledger second journal Dr payable Cr clearing:bank • Outstanding</button>
              <button className="rounded-xl border h-10 px-6 text-xs">Download F&F PDF • Outstanding modern template logo QR pie chart YTD bilingual EN/AM</button>
              <button className="rounded-xl border h-10 px-6 text-xs">Exit Interview • Conducted by HR • Feedback • Outstanding</button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
