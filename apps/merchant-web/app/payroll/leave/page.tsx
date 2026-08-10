"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockBalances = [
  { employee: "Abebe Kebede • EMP001", leave_type: "annual", year: 2026, entitled: 14, used: 2, remaining: 12, carry_forward: 0, status: "active" },
  { employee: "Abebe Kebede • EMP001", leave_type: "sick", year: 2026, entitled: 180, used: 5, remaining: 175, carry_forward: 0, status: "active" },
  { employee: "Almaz Tadesse • EMP002", leave_type: "annual", year: 2026, entitled: 15, used: 5, remaining: 10, carry_forward: 2, status: "active" },
  { employee: "Sara Getachew • EMP006", leave_type: "maternity", year: 2026, entitled: 120, used: 0, remaining: 120, carry_forward: 0, status: "active" },
]
const mockRequests = [
  { id: "lreq_001", employee: "Abebe Kebede • EMP001", leave_type: "annual", start_date: "2026-07-10", end_date: "2026-07-12", days_requested: 2, reason: "Family event", status: "approved", approved_by: "HR Manager", approved_at: "2026-07-09" },
  { id: "lreq_002", employee: "Abebe Kebede • EMP001", leave_type: "sick", start_date: "2026-07-20", end_date: "2026-07-21", days_requested: 2, reason: "Fever - medical certificate attached", status: "pending", medical_file: "sick_cert_EMP001.pdf" },
  { id: "lreq_003", employee: "Sara Getachew • EMP006", leave_type: "maternity", start_date: "2026-08-01", end_date: "2026-11-28", days_requested: 120, reason: "Maternity leave 120 days (30 pre +90 post) per Art 86", status: "approved", approved_by: "HR Manager" },
]

export default function LeaveManagementPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Leave Management • የእረፍት አስተዳደር • Annual 14+1 up to 35 Sick 6 Months Maternity 120 Days (Art 77/82/86) • Ethiopian Calendar Enkutatash Meskerem 1 = Sept 11 • 13 Months</h1>
          <p className="text-sm text-muted-foreground mt-2">Manage Ethiopian labour-law leave (Annual, Sick, Maternity, Paternity) — request, approve, and track your balance.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Leave Balances • Entitled / Used / Remaining • Per Employee Per Year Per Type • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockBalances.map((b, i) => (
                <div key={i} className="rounded-xl border p-3 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-xs">{b.employee} • {b.leave_type} • {b.year}</p><Badge variant={b.leave_type==="annual" ? "success" : b.leave_type==="sick" ? "warning" : "default"}>{b.leave_type}</Badge></div>
                  <div className="mt-2 grid grid-cols-4 gap-2 text-[11px]"><span>Entitled {b.entitled}</span><span>Used {b.used}</span><span>Remaining {b.remaining}</span><span>Carry {b.carry_forward}</span></div>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(b.used/b.entitled)*100}%` }} /></div>
                  <p className="text-[10px] text-muted-foreground mt-1">Art 77 Annual 14+1 up to 35 • Art 82 Sick 180 days 30 days 100% 60 days 50% 90 days unpaid • Art 86 Maternity 120 days 30 pre +90 post</p>
                </div>
              ))}
            </div>
            <button className="mt-4 w-full rounded-xl border border-dashed h-10 text-xs">+ Add Balance • Entitled Used Remaining Carry Forward Year • Art 77/82/86</button>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <div className="flex justify-between items-center"><h3 className="font-semibold">Leave Requests • Start Date End Date Days Requested 0.5 Half Day • Status Pending/Approved/Rejected • Outstanding Pipeline Visual Stepper</h3><button className="rounded-xl bg-primary text-white h-9 px-4 text-xs">+ Request Leave • Annual Sick Maternity Paternity Marriage Mourning Unpaid Comp Off Study</button></div>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Employee</span><span>Type</span><span>Start → End</span><span>Days</span><span>Reason</span><span>Status</span><span>Action</span></div>
              {mockRequests.map(r => (
                <div key={r.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>{r.employee}</span>
                  <span><Badge variant={r.leave_type==="annual" ? "success" : r.leave_type==="maternity" ? "default" : "warning"}>{r.leave_type}</Badge></span>
                  <span>{r.start_date} → {r.end_date}</span>
                  <span>{r.days_requested} days • {r.days_requested===0.5 ? "Half day" : ""}</span>
                  <span className="text-[11px]">{r.reason} {r.medical_file ? `• Medical cert ${r.medical_file} MinIO for sick >3 days per Art 82` : ""}</span>
                  <span><Badge variant={r.status==="approved" ? "success" : r.status==="pending" ? "warning" : "danger"}>{r.status} • Approved by {r.approved_by || "—"} {r.approved_at || ""}</Badge></span>
                  <span className="flex gap-2"><button className="text-primary">Approve • Deduct from balance Used+=Requested Remaining=Entitled-Used floor zero</button><button className="text-red-500">Reject</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Leave Validation per Ethiopia Law • Outstanding Algorithm • O(n) where n=days between start and end</p>
              <ul className="list-disc list-inside mt-2 space-y-1 text-muted-foreground">
                <li>Annual insufficient balance check: requested {">"} remaining ? Error insufficient annual leave balance requested remaining entitled per Art 77 14 days first year +1 per year up to 35</li>
                <li>Sick max 6 months (180 days) per 12 months per Art 82 already exhausted 180 days? Error sick leave max 6 months per 12 months per Art 82 already exhausted • First 30 days 100% pay next 60 days 50% pay remaining 90 days unpaid job protected need medical certificate file_key MinIO presigned 15m TTL &lt;5MB pdf/jpg/png</li>
                <li>Maternity max 120 days consecutive cannot split per Art 86 30 prenatal +90 postnatal • Paternity max 3 days company policy beyond law unpaid no balance check but approval LOP</li>
                <li>Unpaid no balance check but approval will be LOP • Annual, maternity, paternity, marriage, mourning paid no LOP • Comp Off etc</li>
                <li>CalculateLeaveDays inclusive days between start and end excluding weekends optional O(n) where n=days between small ranges • For ET weekends Saturday Sunday? Per Art 75 rest day is Sunday? But many businesses Saturday/Sunday off For simplicity inclusive days count excluding public holidays would need holiday calendar Ethiopian calendar Enkutatash etc</li>
                <li>Payroll integration LOP calculation from leave: CalculateLOPFromLeave leaveRequests attendanceMonth O(n) for each request status approved leaveType unpaid LOP+=daysRequested sick if days{">"}30 and {"<="}90 then LOP 50% of excess {">"}30 if {">"}90 LOP 100% unpaid 90 days annual maternity paternity marriage mourning paid no LOP • Outstanding</li>
              </ul>
            </div>

            <div className="mt-6 rounded-xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Create Leave Request • Outstanding Form • Medical Certificate MinIO • 0.5 Half Day</h4>
              <div className="mt-3 grid grid-cols-4 gap-3 text-xs">
                <div><label>Employee • EMP001</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>EMP001 Abebe Kebede • Annual 14/12 remaining • Sick 180/175 remaining</option></select></div>
                <div><label>Leave Type • annual/sick/maternity/paternity/marriage/mourning/unpaid/comp_off/study</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>annual • Art 77 14+1 up to 35</option><option>sick • Art 82 6 months 30 days 100% 60 days 50% 90 days unpaid</option><option>maternity • Art 86 120 days 30 pre +90 post</option><option>paternity • Company policy 3 days</option></select></div>
                <div><label>Start Date • 2026-07-10</label><input type="date" defaultValue="2026-07-10" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label>End Date • 2026-07-12</label><input type="date" defaultValue="2026-07-12" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label>Days Requested • 2 • 0.5 half day</label><input type="number" defaultValue={2} step={0.5} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label>Reason • Family event</label><input placeholder="Family event • Medical • Maternity • etc" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label>Medical Certificate • MinIO for sick &gt;3 days per Art 82 • File &lt;5MB pdf/jpg/png presigned 15m</label><input type="file" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div className="flex items-end"><button className="rounded-xl bg-primary text-white h-9 px-6">Request Leave • Validate per Ethiopia law Art 77/82/86 • O(n) days between • Balance entitled used remaining</button></div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
