"use client"
import * as React from "react"
import { motion } from "framer-motion"

function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}
function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }

const mockItems = [
  { name: "Abebe Kebede", code: "EMP001", dept: "Engineering", structure: "Fixed CTC 500k", gross: "21250", ctc: "20000", ot: "1250", ot_hours: "5h weekday", taxable: "19850", tax: "1800", pension_emp: "1400", pension_emplr: "2200", net: "16800", paid_days: "25/30", proration: "0.8333", ytd_gross: "140000", ytd_tax: "12000", ytd_net: "98000", status: "calculated", earnings: [{ code: "BASIC", name: "Basic", amount: "16666" }, { code: "HOUSING", name: "Housing", amount: "8333" }, { code: "OT", name: "OT", amount: "1250" }], deductions: [{ code: "TAX", amount: "1800" }, { code: "PENSION", amount: "1400" }, { code: "LOAN", amount: "5000" }], hold: false },
  { name: "Almaz Tadesse", code: "EMP002", dept: "Sales", gross: "35000", taxable: "33250", tax: "3500", pension_emp: "1750", pension_emplr: "2750", net: "24750", paid_days: "30/30", proration: "1.0", ytd_gross: "210000", ytd_tax: "21000", status: "calculated", bonus: "10000", commission: "5000" },
]

export default function PayrollRunDetail({ params }: { params: { id: string } }) {
  const [activeEmployee, setActiveEmployee] = React.useState(mockItems[0])
  const [showBreakdown, setShowBreakdown] = React.useState<string | null>(null)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-2xl font-bold">Payroll Run Detail • {params.id} • July 2026 Regular • 07/2026</h1>
            <div className="mt-2 flex items-center gap-2 text-xs">
              <span className="px-3 py-1 rounded-full bg-neutral-200">draft → calculating → pending_approval • current</span>
              <Badge variant="warning">pending_approval • Needs dual if &gt;100k net • Maker-checker</Badge>
              <Badge variant="success">Ledger M4 balanced ✓</Badge>
              <Badge variant="default">Variance +5.2% vs Jun</Badge>
            </div>
          </div>
          <div className="flex gap-2">
            <button className="rounded-xl border bg-white h-10 px-4 text-xs">Hold Salary • LOP</button>
            <button className="rounded-xl border bg-white h-10 px-4 text-xs">Add Variable Pay</button>
            <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs shadow-medium">Approve • dual finance+admin avatar</button>
            <button className="rounded-xl bg-green-600 text-white h-10 px-6 text-xs shadow-medium">Disburse → payout batch pain.001 + Pension CSV + ERCA CSV</button>
          </div>
        </div>

        {/* KPI sticky footer top */}
        <div className="grid grid-cols-5 gap-4">
          <Card className="p-4"><p className="text-[11px] text-muted-foreground">Total Gross • አጠቃላይ</p><p className="text-xl font-bold">ETB 200,000 • 10 emps • Paid days 280/300 • LOP 20 days • OT 25h</p></Card>
          <Card className="p-4"><p className="text-[11px] text-muted-foreground">Total Deductions</p><p className="text-xl font-bold">ETB 50,000 • Tax 20k + Pension Emp 14k + Loans 10k + Other 6k</p></Card>
          <Card className="p-4"><p className="text-[11px] text-muted-foreground">Total Net • የተጣራ</p><p className="text-xl font-bold">ETB 150,000 • Disburse via Bank IPS CBE/Awash</p></Card>
          <Card className="p-4"><p className="text-[11px] text-muted-foreground">Employer Cost</p><p className="text-xl font-bold">ETB 222,000 = Gross 200k + Pension Emplr 22k • Cost center allocation</p></Card>
          <Card className="p-4"><p className="text-[11px] text-muted-foreground">Compliance</p><p className="text-xs"><Badge variant="success">Pension CSV generated ✓</Badge><Badge variant="success">ERCA CSV ✓</Badge><Badge variant="success">Bank pain.001 ✓</Badge></p></Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="lg:col-span-2 overflow-hidden">
            <div className="p-4 flex justify-between items-center border-b"><h3 className="font-semibold text-sm">Payroll Items • 10 • Earnings breakdown + Deductions breakdown + Employer 11% + YTD + Proration factor</h3><div className="flex gap-2"><input placeholder="Search employee..." className="rounded-xl border h-8 px-3 text-xs w-48" /><button className="rounded-xl border h-8 px-3 text-xs">Export XLSX</button></div></div>
            <div className="overflow-auto max-h-[600px]">
              <div className="grid grid-cols-10 gap-2 bg-muted p-3 text-[11px] font-semibold sticky top-0"><span>Employee</span><span>Gross</span><span>OT/Bonus</span><span>Taxable</span><span>Tax ET binary O(log n)</span><span>Pension 7/11</span><span>Net</span><span>Paid/LOP Factor</span><span>YTD</span><span>Status</span></div>
              {mockItems.map((r,i)=>(
                <div key={i} className={`grid grid-cols-10 gap-2 p-3 border-t text-xs hover:bg-muted/50 cursor-pointer ${activeEmployee.code===r.code ? "bg-primary/5 border-l-4 border-l-primary" : ""}`} onClick={()=>setActiveEmployee(r)}>
                  <span className="font-medium">{r.name}<span className="block text-[10px] text-muted-foreground">{r.code} • {r.dept} • {r.structure || "Basic"}</span></span>
                  <span>{r.gross}<span className="block text-[10px]">CTC {r.ctc || r.gross}</span></span>
                  <span>{r.ot || r.bonus}+{r.ot_hours || ""}<span className="block text-[10px]">{r.commission ? `Comm ${r.commission}` : ""}</span></span>
                  <span>{r.taxable}</span>
                  <span>{r.tax}<span className="block text-[10px]">{parseInt(r.tax)>2000 ? "bracket 5251-7800 25%-565" : "1651-3200 15%-142.5"}</span></span>
                  <span>{r.pension_emp}/{r.pension_emplr}<span className="block text-[10px]">Emp 7% Emplr 11% total {parseInt(r.pension_emp)+parseInt(r.pension_emplr)}</span></span>
                  <span className="font-bold">{r.net}</span>
                  <span>{r.paid_days}<span className="block text-[10px]">Factor {r.proration}</span></span>
                  <span className="text-[10px]">Gross {r.ytd_gross} Tax {r.ytd_tax} Net {r.ytd_net}</span>
                  <span><Badge variant="success">{r.status}</Badge></span>
                </div>
              ))}
              <div className="grid grid-cols-10 gap-2 p-3 bg-muted font-bold text-xs sticky bottom-0"><span>Total 10 emps</span><span>200,000</span><span>OT 5k Bonus 15k Comm 5k</span><span>185,000</span><span>20,000</span><span>14k/22k pension total 36k</span><span>150,000 net</span><span>Paid 280/300 LOP 20 Factor avg 0.93</span><span>YTD Gross 1.2M YTD Net 980k</span><span>calculated</span></div>
            </div>
            <div className="p-3 flex gap-2 border-t">
              <button className="rounded-xl border h-9 px-4 text-xs">Download Payslips PDF ZIP • QR verified outstanding modern</button>
              <button className="rounded-xl border h-9 px-4 text-xs">ET Report CSV ERCA • JSON • Pension CSV • Cost center CSV</button>
              <button className="rounded-xl border h-9 px-4 text-xs">Bank File pain.001 XML • MT940 reconciliation window 24h tolerance 0.01</button>
            </div>
          </Card>

          <div className="space-y-4">
            <Card className="p-4">
              <h3 className="font-semibold text-sm">Payslip PDF Preview • {activeEmployee.name} • {activeEmployee.code} • Outstanding Modern</h3>
              <motion.div initial={{ opacity:0, y:5 }} animate={{ opacity:1, y:0 }} key={activeEmployee.code} className="mt-3 rounded-xl border p-4 bg-gradient-to-br from-white to-neutral-50 shadow-soft max-w-[400px]">
                <div className="flex justify-between"><span className="font-bold">Apex Trading PLC • አፔክስ</span><span className="text-xs">July 2026 • 07/2026 • Regular</span></div>
                <div className="mt-2 text-xs space-y-1">
                  <p>Employee: {activeEmployee.name} • {activeEmployee.code} • {activeEmployee.dept} • Fayda ****1234 ✓ face 0.92 • Bank CBE ****1234 • TIN 0098765432 • Pension PEN-001</p>
                  <p>CTC Monthly {activeEmployee.ctc || activeEmployee.gross} • Paid Days {activeEmployee.paid_days} • LOP 5 • Factor {activeEmployee.proration}</p>
                  <div className="mt-2 grid grid-cols-2 gap-2 text-[11px] border rounded-xl p-2 bg-muted/50">
                    <div><p className="font-semibold">Earnings • Gross {activeEmployee.gross}</p>{activeEmployee.earnings?.map((e:any)=><p key={e.code}>{e.code} {e.name}: ETB {e.amount}</p>) || <p>Basic 20000 + OT 1250 (5h weekday 1.25x hourly 96.15) = Gross 21250</p>}</div>
                    <div><p className="font-semibold">Deductions</p><p>Taxable 19850 • Tax 1800 (bracket 1651-3200 15%-142.5) • Pension Emp 7% 1400 • Loan EMI 5000 • Other 0</p><p>Employer Pension 11% 2200 • Total Employer Cost 22200</p></div>
                  </div>
                  <p className="font-bold text-base">Net Pay ETB {activeEmployee.net} • አጠቃላይ የተጣራ</p>
                  <p className="text-[10px]">YTD Gross {activeEmployee.ytd_gross} • YTD Tax {activeEmployee.ytd_tax} • YTD Net {activeEmployee.ytd_net}</p>
                  <div className="mt-2 h-12 bg-muted/80 rounded flex items-center justify-center text-[10px] gap-2">
                    <span>Pie chart deductions • Tax 30% Pension 20% Loan 40%</span>
                    <span className="h-8 w-8 border border-dashed rounded flex items-center justify-center">QR</span>
                    <span>QR verification https://apexpay.et/verify/payslip/{params.id}/{activeEmployee.code} signed JWT</span>
                  </div>
                  <p className="text-[9px] text-muted-foreground mt-2">This is computer generated payslip no signature required • Password protected DOB DDMM + last4 • Bilingual EN/AM • Digitally verified via QR • Ledger M4 per run book Dr expense:salary Cr payroll_payable</p>
                </div>
              </motion.div>
              <div className="mt-3 flex gap-2">
                <button className="rounded-xl bg-primary text-white h-8 px-3 text-[11px]">Download PDF • gofpdf + barcode/qr</button>
                <button className="rounded-xl border h-8 px-3 text-[11px]">Email • SMTP</button>
                <button className="rounded-xl border h-8 px-3 text-[11px]">WhatsApp Share • share_plus</button>
              </div>
            </Card>

            <Card className="p-4">
              <h3 className="font-semibold text-sm">Approval Flow • Maker-checker dual &gt;100k net • Outstanding avatars</h3>
              <div className="mt-3 space-y-3">
                <div className="flex items-center gap-3"><div className="h-8 w-8 rounded-full bg-neutral-200 flex items-center justify-center text-xs">HR</div><div><p className="text-xs font-medium">Meron HR • Created run</p><p className="text-[11px] text-muted-foreground">2 min ago • draft → calculating → pending_approval</p></div><span className="ml-auto text-green-600 text-xs">✓</span></div>
                <div className="flex items-center gap-3"><div className="h-8 w-8 rounded-full bg-amber-500 text-white flex items-center justify-center text-xs">F</div><div><p className="text-xs font-medium">Finance • Pending approval</p><p className="text-[11px] text-muted-foreground">Needs 2nd approver &gt;100k • Total net 150k ETB</p></div><span className="ml-auto px-2 py-0.5 rounded-full bg-amber-500/15 text-amber-700 text-[10px]">Pending</span></div>
                <div className="flex items-center gap-3 opacity-50"><div className="h-8 w-8 rounded-full bg-neutral-200 flex items-center justify-center text-xs">A</div><div><p className="text-xs font-medium">Admin • Final approval</p><p className="text-[11px] text-muted-foreground">After finance • Then disburse payout batch</p></div></div>
              </div>
              <div className="mt-4 rounded-xl bg-blue-500/10 border border-blue-500/20 p-3 text-[11px]">
                <p>Audit log O(1) advisory lock pg_advisory_xact_lock(hashtext(book_id)) • payroll_audit_logs actor finance action approve_run details approved_by • IP inet request_id • Immutable</p>
              </div>
            </Card>

            <Card className="p-4">
              <h3 className="font-semibold text-sm">Compliance & Bank File • Generated</h3>
              <div className="mt-2 space-y-2 text-[11px]">
                <div className="flex justify-between"><span>Pension CSV • Social Security Agency</span><span className="text-green-600">✓ generated 10 emps 36k</span></div>
                <div className="flex justify-between"><span>ERCA Withholding CSV</span><span className="text-green-600">✓ 20k tax</span></div>
                <div className="flex justify-between"><span>Bank Disbursal pain.001 XML</span><span className="text-green-600">✓ 150k CBE 10 txs</span></div>
                <div className="flex justify-between"><span>Cost Center Report</span><span className="text-green-600">✓ Engineering 100k Sales 100k</span></div>
              </div>
              <div className="mt-3 flex gap-2">
                <button className="rounded-xl border h-8 px-3 text-[11px]">Download All ZIP • MinIO presigned 15m</button>
                <button className="rounded-xl border h-8 px-3 text-[11px]">View Ledger M4 • Dr salary 200k + Dr pension emplr 22k Cr payable 150k Cr tax 20k Cr pension 36k balanced</button>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
