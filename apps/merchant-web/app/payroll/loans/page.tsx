"use client"
import * as React from "react"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from "recharts"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockLoans = [
  { id: "loan_001", employee: "Abebe Kebede • EMP001", type: "salary_advance", principal: "20000", interest_rate: "0", tenure: 4, emi: "5000", total_paid: "5000", outstanding: "15000", status: "active", disbursed_at: "2026-06-01", next_due: "2026-08-01", reason: "Family emergency advance", schedule: [
    { installment_no: 1, due_date: "2026-07-01", emi: "5000", principal: "5000", interest: "0", outstanding_after: "15000", status: "paid", paid_at: "2026-07-01", run_id: "prun_July2026" },
    { installment_no: 2, due_date: "2026-08-01", emi: "5000", principal: "5000", interest: "0", outstanding_after: "10000", status: "pending" },
    { installment_no: 3, due_date: "2026-09-01", emi: "5000", principal: "5000", interest: "0", outstanding_after: "5000", status: "pending" },
    { installment_no: 4, due_date: "2026-10-01", emi: "5000", principal: "5000", interest: "0", outstanding_after: "0", status: "pending" },
  ]},
  { id: "loan_002", employee: "Almaz Tadesse • EMP002", type: "personal", principal: "50000", interest_rate: "5", tenure: 6, emi: "8583", total_paid: "8583", outstanding: "41417", status: "active", disbursed_at: "2026-06-15", next_due: "2026-08-15", reason: "Personal loan for housing", schedule: [
    { installment_no: 1, due_date: "2026-07-15", emi: "8583", principal: "8333", interest: "250", outstanding_after: "41667", status: "paid" },
    { installment_no: 2, due_date: "2026-08-15", emi: "8583", principal: "8333", interest: "250", outstanding_after: "33334", status: "pending" },
  ]},
]

export default function LoansPage() {
  const [selected, setSelected] = React.useState(mockLoans[0])
  const chartData = selected.schedule.map((s:any)=>({ name: `#${s.installment_no} ${s.due_date}`, principal: parseInt(s.principal), interest: parseInt(s.interest), outstanding: parseInt(s.outstanding_after) }))
  const pieData = [
    { name: "Principal", value: parseInt(selected.principal), color: "#0B6E4F" },
    { name: "Interest", value: selected.interest_rate!=="0" ? Math.round(parseInt(selected.principal)*parseInt(selected.interest_rate)/100*selected.tenure/12) : 0, color: "#EAB308" },
    { name: "Paid", value: parseInt(selected.total_paid), color: "#10B981" },
    { name: "Outstanding", value: parseInt(selected.outstanding), color: "#E4E4E7" },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Loans & Advances • ብድር • EMI Schedule Repayment Tracking UI • Payroll Deduction Auto O(k) • Recharts Bar Principal vs Interest • Ethiopia Business Practice</h1>
            <p className="text-sm text-muted-foreground mt-2">Loan types personal/salary_advance/housing/education/medical/other principal interest_rate tenure_months emi_amount total_paid outstanding status draft/pending_approval/approved/active/closed/rejected/written_off disbursed_at next_due_date approved_by reason meta JSON + loan_repayments run_id amount principal_component interest_component outstanding_after status pending/paid/failed + loan_emi_schedule installment_no due_date emi_amount principal_component interest_component outstanding_after status pending/paid/overdue/skipped paid_at run_id per payroll run auto deduction O(k) per employee k=0-2 active loans • Outstanding modern UI glassmorphic Recharts bar principal vs interest pie deductions loan 40% tax 30% pension 20%</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Request Loan • EMI Preview Formula Simple Interest • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Loans • Active • EMI Auto Deduction Per Run • O(k) • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockLoans.map(loan => (
                <button key={loan.id} onClick={()=>setSelected(loan)} className={`w-full text-left rounded-xl border p-4 hover:bg-muted ${selected.id===loan.id ? "bg-primary/10 border-primary/30" : ""}`}>
                  <div className="flex justify-between"><p className="font-medium text-sm">{loan.employee} • {loan.type}</p><Badge variant={loan.status==="active" ? "success" : "warning"}>{loan.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Principal {loan.principal} ETB • Interest {loan.interest_rate}% • Tenure {loan.tenure}mo • EMI {loan.emi} • Paid {loan.total_paid} • Outstanding {loan.outstanding} • Next due {loan.next_due}</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(parseInt(loan.total_paid)/parseInt(loan.principal))*100}%` }} /></div>
                </button>
              ))}
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">EMI Schedule • Repayment Tracking UI • Recharts Bar Principal vs Interest • {selected.employee} • Principal {selected.principal} • EMI {selected.emi} • Outstanding {selected.outstanding} • Pie Deductions Loan 40% Tax 30% Pension 20%</h3>
            <div className="mt-4 grid grid-cols-2 gap-4">
              <div className="h-64">
                <p className="text-xs font-semibold mb-2">Bar Chart Principal vs Interest per Installment • Outstanding Recharts • O(n) per loan n=tenure months</p>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" tick={{ fontSize: 9 }} />
                    <YAxis tick={{ fontSize: 10 }} />
                    <Tooltip />
                    <Bar dataKey="principal" fill="#0B6E4F" name="Principal" />
                    <Bar dataKey="interest" fill="#EAB308" name="Interest" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
              <div className="h-64">
                <p className="text-xs font-semibold mb-2">Pie Chart • Principal vs Interest vs Paid vs Outstanding • Outstanding modern template QR verification</p>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie data={pieData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label>
                      {pieData.map((entry, index) => <Cell key={`cell-${index}`} fill={entry.color} />)}
                    </Pie>
                    <Tooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="mt-6 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-8 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>No</span><span>Due Date</span><span>EMI</span><span>Principal</span><span>Interest</span><span>Outstanding After</span><span>Status</span><span>Run ID</span></div>
              {selected.schedule.map((s:any)=>(
                <div key={s.installment_no} className="grid grid-cols-8 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>#{s.installment_no}</span><span>{s.due_date}</span><span className="font-bold">ETB {s.emi}</span><span>ETB {s.principal}</span><span>ETB {s.interest} • {selected.interest_rate}%</span><span>ETB {s.outstanding_after}</span><span><Badge variant={s.status==="paid" ? "success" : "warning"}>{s.status}</Badge></span><span>{s.run_id || "—"}</span>
                </div>
              ))}
            </div>

            <div className="mt-6 grid grid-cols-2 gap-4 text-xs">
              <div className="rounded-xl bg-muted p-4"><p className="font-semibold">EMI Calculation Logic Outstanding • Simple Interest for Demo Real Amortized • O(1)</p><p className="mt-2 font-mono text-[11px]">EMI = principal / tenure_months Round2 • If interest_rate !=0: interestTotal = principal * rate/100 * tenure/12 • EMI = (principal + interestTotal) / tenure Round2 • Example: principal 20000 interest 0% tenure 4 EMI 5000 • Principal 50000 interest 5% tenure 6 interestTotal 1250 EMI 8541? Mock 8583 • Outstanding after each installment = outstanding_before - principal_component • Total paid += EMI • Outstanding -= EMI • Closed if outstanding ≤0</p></div>
              <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-4"><p className="font-semibold">Repayment Tracking UI • Outstanding • Audit log • Ledger M4 • Beyond RazorpayX</p><p className="mt-2 text-[11px]">Loan repayments table loan_repayments id loan_id run_id employee_id amount principal_component interest_component outstanding_after status created_at + loan_emi_schedule installment_no due_date emi_amount principal_component interest_component outstanding_after status paid_at run_id per payroll run • O(n) per loan n=tenure • Repayment history per loan per employee • Chart Recharts bar principal vs interest • Pie deductions loan 40% tax 30% pension 20% • Outstanding modern template QR verification • Audit payroll_audit_logs actor_type hr/finance/admin/employee action loan_requested loan_approved loan_repayment details JSON IP inet request_id immutable</p></div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
