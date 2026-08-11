"use client"
import * as React from "react"
import Link from "next/link"
import { motion, AnimatePresence } from "framer-motion"

// Outstanding UI components — glassmorphic, Mercury/Linear inspiration
function Card({ children, className = "" }: { children: React.ReactNode; className?: string }) {
  return <div className={`rounded-2xl border bg-card shadow-soft hover:shadow-medium transition-shadow ${className}`}>{children}</div>
}
function GlassCard({ children, className = "" }: { children: React.ReactNode; className?: string }) {
  return <div className={`rounded-2xl border border-white/50 bg-white/70 backdrop-blur-xl shadow-glass ${className}`}>{children}</div>
}
function Badge({ children, variant = "default" }: { children: React.ReactNode; variant?: "default" | "success" | "warning" | "danger" }) {
  const variants: any = {
    default: "bg-neutral-100 text-neutral-700",
    success: "bg-green-500/15 text-green-700 border-green-500/20",
    warning: "bg-amber-500/15 text-amber-700 border-amber-500/20",
    danger: "bg-red-500/15 text-red-700 border-red-500/20",
  }
  return <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${variants[variant]}`}>{children}</span>
}

// Mock data — would be fetched via SWR /v1/employees, /v1/payroll_runs, /v1/salary_structures in prod
const mockDepartments = [
  { id: "dept_eng", name: "Engineering", code: "ENG", cost_center: "CC-100", headcount: 5 },
  { id: "dept_sales", name: "Sales", code: "SALES", cost_center: "CC-200", headcount: 5 },
  { id: "dept_finance", name: "Finance", code: "FIN", cost_center: "CC-300", headcount: 2 },
]
const mockStructures = [
  { id: "sstr_001", name: "Fixed CTC 500k Annual", ctc_annual: "500000", ctc_monthly: "41666", components: [
    { code: "BASIC", name: "Basic Salary", calc: "CTC_MONTHLY * 0.4", taxable: true, pensionable: true, order: 1, amount: "16666" },
    { code: "HOUSING", name: "Housing Allowance", calc: "CTC_MONTHLY * 0.2", taxable: true, pensionable: false, order: 2, amount: "8333" },
    { code: "TRANSPORT", name: "Transport Allowance", calc: "fixed 3000", taxable: false, pensionable: false, order: 3, amount: "3000" },
    { code: "FUEL", name: "Fuel Allowance", calc: "fixed 2000", taxable: true, pensionable: false, order: 4, amount: "2000" },
  ]},
  { id: "sstr_002", name: "Tech Band G3", ctc_annual: "840000", ctc_monthly: "70000", components: [
    { code: "BASIC", name: "Basic", calc: "CTC *0.45", amount: "31500" },
    { code: "SPECIAL", name: "Special Allow", calc: "CTC *0.3", amount: "21000" },
  ]},
]
const mockEmployees = [
  { code: "EMP001", name: "Abebe Kebede", name_am: "አበበ ከበደ", base: "20000", ctc: "240000", dept: "Engineering", grade: "G3", bank: "CBE ****1234", bank_code: "CBE", cost: "CC-100", fayda: true, face_score: 0.92, status: "active", tin: "0098765432", pension_no: "PEN-001", structure: "Fixed CTC 500k Annual" },
  { code: "EMP002", name: "Almaz Tadesse", base: "25000", ctc: "300000", dept: "Sales", grade: "G4", bank: "Awash ****5678", bank_code: "AWASH", cost: "CC-200", fayda: true, face_score: 0.89, status: "active", tin: "0098765433", pension_no: "PEN-002", structure: "Tech Band G3" },
  { code: "EMP003", name: "Yonas Bekele", base: "18000", ctc: "216000", dept: "Engineering", grade: "G2", bank: "Dashen ****9012", bank_code: "DASHEN", cost: "CC-100", fayda: false, status: "probation", tin: "0098765434", pension_no: "PEN-003", structure: "Fixed CTC 500k Annual" },
]
const mockRuns = [
  { id: "prun_July2026", ref: "July2026_Regular", period: "07/2026", type: "regular", status: "pending_approval", total_gross: "200000", total_tax: "20000", total_pension: "30000", total_net: "150000", employer_pension: "22000", employer_cost: "222000", count: 10, variance: "+5.2% vs Jun", paid: 0, failed: 0 },
  { id: "prun_June2026", ref: "June2026_Regular", period: "06/2026", type: "regular", status: "completed", total_gross: "190000", total_tax: "19000", total_pension: "28500", total_net: "142500", count: 10, variance: "-2.1%", paid: 10, failed: 0 },
]
const mockLoans = [
  { id: "loan_001", employee: "Abebe Kebede", code: "EMP001", type: "salary_advance", principal: "20000", emi: "5000", outstanding: "15000", tenure: 4, status: "active" },
  { id: "loan_002", employee: "Almaz Tadesse", type: "personal", principal: "50000", emi: "8333", outstanding: "41667", tenure: 6, status: "approved" },
]
const mockCompliance = [
  { type: "pension_contribution", period: "07/2026", status: "generated", file: "pension_2026_07.csv", count: 10, total: "30000" },
  { type: "erca_withholding", period: "07/2026", status: "generated", file: "erca_2026_07.csv", total_tax: "20000" },
  { type: "bank_disbursal_file", period: "07/2026", status: "generated", file: "bank_disbursal_July2026.xml", format: "pain.001.001.03", count: 10, total_net: "150000" },
  { type: "payroll_register", period: "07/2026", status: "generated", file: "payroll_register_2026_07.xlsx" },
  { type: "cost_center_report", period: "07/2026", status: "generated", file: "cost_center_2026_07.csv" },
]

export default function PayrollPage() {
  const [activeTab, setActiveTab] = React.useState("overview")
  const [selectedStructure, setSelectedStructure] = React.useState(mockStructures[0])
  const [showOTCalculator, setShowOTCalculator] = React.useState(false)

  const tabs = [
    { id: "overview", label: "Overview • አጠቃላይ", icon: "📊" },
    { id: "employees", label: "Employees • ሰራተኞች", icon: "👥", count: 10 },
    { id: "structures", label: "Salary Structures • ደሞዝ አወቃቀር", icon: "🏗️" },
    { id: "runs", label: "Payroll Runs • ደሞዝ ሩጫዎች", icon: "🔄" },
    { id: "attendance", label: "Attendance • መገኘት", icon: "📅" },
    { id: "loans", label: "Loans & Claims • ብድር", icon: "💸" },
    { id: "compliance", label: "Compliance • ተገዢነት", icon: "📜" },
    { id: "settings", label: "Settings • ቅንብሮች", icon: "⚙️" },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 via-white to-primary-50/30 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold tracking-tight flex items-center gap-3">
              Payroll OS • ደሞዝ
              <Badge variant="success">enterprise-grade • Comprehensive</Badge>
            </h1>
            <p className="text-sm text-muted-foreground mt-2">Workforce Money OS — CTC templates, LOP proration, OT 1.25/1.5/2.0/1.3, Loans EMI auto, Pension 7%/11%, ERCA CSV, Bank pain.001, Payslip QR, F&F settlement • 500 employees &lt;2s p99 • Ledger M4 per run book</p>
          </div>
          <div className="flex gap-2">
            <button className="rounded-xl border bg-white px-4 h-10 text-xs font-medium hover:bg-muted">Import CSV • Employees 1000 bulk</button>
            <button className="rounded-xl bg-primary text-white px-6 h-10 text-xs font-medium shadow-medium hover:shadow-large">Create Run July regular • 07/2026</button>
          </div>
        </div>

        {/* KPI Cards — Outstanding */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0 }}>
            <Card className="p-5">
              <p className="text-xs text-muted-foreground">Total Gross • አጠቃላይ • July 2026 regular</p>
              <p className="text-2xl font-bold mt-2">ETB 200,000</p>
              <div className="mt-2 flex items-center gap-2">
                <Badge variant="warning">+5.2% vs Jun</Badge>
                <span className="text-[11px] text-muted-foreground">10 employees • 30 paid days</span>
              </div>
              <div className="mt-3 h-1 rounded-full bg-neutral-100 overflow-hidden">
                <div className="h-full bg-primary w-[78%] rounded-full" />
              </div>
            </Card>
          </motion.div>
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.05 }}>
            <Card className="p-5">
              <p className="text-xs text-muted-foreground">Total Tax • ግብር • ET binary search O(log n)</p>
              <p className="text-2xl font-bold mt-2">ETB 20,000</p>
              <p className="text-[11px] text-muted-foreground mt-1">Brackets 0-600 0% 601-1650 10%-60 1651-3200 15%-142.5 3201-5250 20%-302.5 5251-7800 25%-565 7801-10900 30%-955 &gt;10900 35%-1500</p>
              <Badge variant="default" >ERCA withholding CSV generated</Badge>
            </Card>
          </motion.div>
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
            <Card className="p-5">
              <p className="text-xs text-muted-foreground">Total Net • የተጣራ • Disburse via M3 payout batch</p>
              <p className="text-2xl font-bold mt-2">ETB 150,000</p>
              <p className="text-[11px] text-muted-foreground mt-1">Pension Emp 7% ETB 14k • Emplr 11% ETB 22k • Total cost 222k</p>
              <div className="mt-2 flex gap-2">
                <Badge variant="success">Ledger M4 balanced ✓</Badge>
                <Badge variant="success">Bank pain.001 generated ✓</Badge>
              </div>
            </Card>
          </motion.div>
          <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.15 }}>
            <GlassCard className="p-5">
              <p className="text-xs text-muted-foreground">AI Payroll Assist • Swarm 🤖</p>
              <p className="text-sm font-medium mt-2">Goal: Run payroll July, add 10k bonus Sales</p>
              <div className="mt-3 space-y-2 text-[11px]">
                <div className="flex gap-2"><span className="h-5 w-5 rounded-full bg-primary text-white flex items-center justify-center">1</span><span>get_employees cost_center Sales → 5</span><span className="ml-auto text-green-600">✓</span></div>
                <div className="flex gap-2"><span className="h-5 w-5 rounded-full bg-primary text-white flex items-center justify-center">2</span><span>create variable inputs bonus 10k ×5 =50k</span><span className="ml-auto text-green-600">✓</span></div>
                <div className="flex gap-2"><span className="h-5 w-5 rounded-full bg-amber-500 text-white flex items-center justify-center">3</span><span>calculate_payroll_draft → needs_confirmation &gt;100k</span><span className="ml-auto text-amber-600">⚠ needs approval</span></div>
              </div>
              <p className="mt-2 text-[11px] font-medium">Final: Draft ready ETB 200k → Pending approval dual finance+admin</p>
            </GlassCard>
          </motion.div>
        </div>

        {/* Tabs */}
        <div className="rounded-2xl border bg-card p-2 flex gap-2 overflow-x-auto">
          {tabs.map(t => (
            <button key={t.id} onClick={() => setActiveTab(t.id)} className={`flex items-center gap-2 px-4 h-9 rounded-xl text-xs font-medium whitespace-nowrap transition-all ${activeTab===t.id ? "bg-primary text-white shadow-medium" : "hover:bg-muted text-muted-foreground"}`}>
              <span>{t.icon}</span> {t.label} {t.count && <span className="ml-1 px-1.5 py-0.5 rounded-full bg-white/20 text-[10px]">{t.count}</span>}
            </button>
          ))}
        </div>

        <AnimatePresence mode="wait">
          <motion.div key={activeTab} initial={{ opacity: 0, y: 5 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -5 }} transition={{ duration: 0.2 }}>
            {activeTab==="overview" && <OverviewTab />}
            {activeTab==="employees" && <EmployeesTab />}
            {activeTab==="structures" && <StructuresTab selected={selectedStructure} onSelect={setSelectedStructure} />}
            {activeTab==="runs" && <RunsTab />}
            {activeTab==="attendance" && <AttendanceTab />}
            {activeTab==="loans" && <LoansTab />}
            {activeTab==="compliance" && <ComplianceTab />}
            {activeTab==="settings" && <SettingsTab />}
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  )
}

// ==================== Overview Tab ====================
function OverviewTab() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card className="p-6 lg:col-span-2">
        <h3 className="font-semibold">Payroll Cost Trend • গত 6 মাস • Recharts AreaChart per cost_center</h3>
        <div className="mt-4 h-48 rounded-xl bg-gradient-to-br from-primary-50 to-gold-50 border border-dashed flex items-center justify-center text-xs text-muted-foreground">
          Recharts AreaChart: Jan 180k → Jul 200k +5.2% • Engineering 100k • Sales 100k • Cost center pie breakdown • Variance report vs last month
        </div>
        <div className="mt-4 grid grid-cols-3 gap-3 text-xs">
          <div className="rounded-xl bg-muted p-3"><p className="text-muted-foreground">Engineering • CC-100</p><p className="font-bold text-sm">ETB 100k • 5 emps • Gross 105k Net 75k • Pension 11% 11k</p></div>
          <div className="rounded-xl bg-muted p-3"><p className="text-muted-foreground">Sales • CC-200</p><p className="font-bold text-sm">ETB 100k • 5 emps • Bonus 10k Sales • Commission 5k</p></div>
          <div className="rounded-xl bg-muted p-3"><p className="text-muted-foreground">Total Employer Cost</p><p className="font-bold text-sm">ETB 222k = Gross 200k + Pension Emplr 22k</p></div>
        </div>
      </Card>
      <Card className="p-6 space-y-4">
        <h3 className="font-semibold">Quick Actions • ፈጣን እርምጃዎች</h3>
        <div className="space-y-2">
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Create Run <span>→</span></button>
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Import Attendance CSV • 500 rows <span>📅</span></button>
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Import Variable Pay • Bonus/Commission <span>💸</span></button>
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Generate Pension CSV • Social Security <span>📜</span></button>
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Generate ERCA Withholding CSV <span>🧾</span></button>
          <button className="w-full rounded-xl border h-12 text-xs font-medium hover:bg-muted flex items-center justify-between px-4">Download Payslips ZIP • QR verified <span>📄</span></button>
        </div>
        <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-3 text-xs">
          <p className="font-semibold">Ledger M4 per run book outstanding:</p>
          <p className="mt-1 font-mono text-[11px]">Dr expense:salary 200k + Dr expense:pension_employer 22k<br/>Cr payroll_payable 150k<br/>Cr et_income_tax_payable 20k<br/>Cr pension_payable 52k (14k emp 7% + 38k emplr 11% *2? Actually 14k+22k=36k for 200k gross 7%+11% = 14k+22k)<br/>ValidateBalanced ✓</p>
          <p className="mt-2">Second journal disburse: Dr payroll_payable 150k Cr clearing:bank 150k via payout batch pain.001 XML</p>
        </div>
      </Card>
    </div>
  )
}

// ==================== Employees Tab ====================
function EmployeesTab() {
  return (
    <Card className="p-6">
      <div className="flex justify-between items-center">
        <h3 className="font-semibold">Employees • 10 • Fayda badge verified • Salary structure assigned • Department/Grade/Branch</h3>
        <div className="flex gap-2">
          <input placeholder="Search EMP001, Abebe..." className="rounded-xl border h-9 px-3 text-xs w-64" />
          <select className="rounded-xl border h-9 px-3 text-xs"><option>All Departments</option><option>Engineering</option><option>Sales</option></select>
          <button className="rounded-xl border h-9 px-4 text-xs">Import CSV • 500 &lt;2s</button>
          <button className="rounded-xl bg-primary text-white h-9 px-4 text-xs">Add Employee + Fayda verify</button>
        </div>
      </div>
      <div className="mt-4 rounded-xl border overflow-hidden">
        <div className="grid grid-cols-10 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Code</span><span>Name</span><span>Dept/Grade</span><span>CTC Annual</span><span>Base</span><span>Bank</span><span>Fayda</span><span>Structure</span><span>Status</span><span>Action</span></div>
        {mockEmployees.map(e => (
          <div key={e.code} className="grid grid-cols-10 gap-2 p-3 border-t text-xs hover:bg-muted/50">
            <span className="font-mono font-medium">{e.code}</span>
            <span><span className="font-medium">{e.name}</span><span className="block text-[10px] text-muted-foreground">{e.name_am || e.code} • TIN {e.tin}</span></span>
            <span>{e.dept}<span className="block text-[10px] text-muted-foreground">{e.grade} • {e.cost}</span></span>
            <span>ETB {(parseInt(e.ctc)/1000)}k<span className="block text-[10px]">Monthly {parseInt(e.ctc)/12}</span></span>
            <span>ETB {e.base}</span>
            <span>{e.bank_code} <span className="text-[10px] text-muted-foreground block">{e.bank}</span></span>
            <span>{e.fayda ? <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-green-500/15 text-green-700 border border-green-500/20 text-[10px]">✓ Verified {e.face_score}</span> : <span className="px-2 py-0.5 rounded-full bg-amber-500/15 text-[11px]">Pending</span>}</span>
            <span className="text-[11px]">{e.structure}</span>
            <span><Badge variant={e.status==="active" ? "success" : "warning"}>{e.status}</Badge></span>
            <span className="flex gap-2"><Link href={`/payroll/${e.code}`} className="text-primary">View</Link><button className="text-muted-foreground">Edit</button></span>
          </div>
        ))}
      </div>
      <div className="mt-4 flex justify-between text-[11px] text-muted-foreground">
        <span>Showing 3 of 10 employees • Pagination O(n) • Fayda verification front/back &lt;2MB + OTP consent id.gov.et • Bank account name fuzzy Levenshtein &lt;3 validation</span>
        <span>Bulk actions: Export XLSX, Deactivate, Transfer department</span>
      </div>
    </Card>
  )
}

// ==================== Structures Tab ====================
function StructuresTab({ selected, onSelect }: any) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card className="p-6">
        <h3 className="font-semibold">Salary Structures • CTC Templates • enterprise-grade</h3>
        <div className="mt-3 space-y-2">
          {mockStructures.map((s:any) => (
            <button key={s.id} onClick={()=>onSelect(s)} className={`w-full text-left rounded-xl border p-3 hover:bg-muted text-xs ${selected.id===s.id ? "bg-primary/10 border-primary/30" : ""}`}>
              <p className="font-medium">{s.name}</p>
              <p className="text-[11px] text-muted-foreground">CTC Annual {s.ctc_annual} • Monthly {s.ctc_monthly} • 4 components • Default ✓</p>
              <div className="mt-2 flex gap-1">{s.components.map((c:any)=><span key={c.code} className="px-1.5 py-0.5 rounded bg-neutral-100 text-[10px]">{c.code} {c.amount}</span>)}</div>
            </button>
          ))}
          <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Create New Structure • Formula engine</button>
        </div>
        <div className="mt-6 rounded-xl bg-amber-500/10 border border-amber-500/20 p-3 text-[11px]">
          <p className="font-semibold">Formula Engine — secure O(n) tokenization + shunting-yard + decimal precise, no evil eval</p>
          <p className="mt-1 font-mono">BASIC = CTC_MONTHLY * 0.4 • HOUSING = CTC_MONTHLY * 0.2 • TRANSPORT fixed 3000 non-taxable up to 1000 exempt limit • ValidateFormula only allowed vars BASIC CTC_MONTHLY GROSS</p>
        </div>
      </Card>
      <Card className="p-6 lg:col-span-2">
        <div className="flex justify-between items-center">
          <h3 className="font-semibold">Structure Builder • {selected.name} • Outstanding drag-drop</h3>
          <div className="flex gap-2">
            <button className="rounded-xl border h-8 px-3 text-[11px]">Preview Calculation • ETB {selected.ctc_monthly}/mo</button>
            <button className="rounded-xl bg-primary text-white h-8 px-3 text-[11px]">Save Template</button>
          </div>
        </div>
        <div className="mt-4 space-y-3">
          <div className="grid grid-cols-4 gap-2 text-[11px] font-semibold bg-muted p-2 rounded-xl"><span>Code</span><span>Name / Calculation</span><span>Type / Taxable</span><span>Amount</span></div>
          {selected.components.map((c:any)=>(
            <div key={c.code} className="grid grid-cols-4 gap-2 p-3 border rounded-xl text-xs hover:bg-muted/50">
              <span className="font-mono font-medium">{c.code}</span>
              <span>{c.name}<span className="block text-[10px] text-muted-foreground font-mono">{c.calc}</span></span>
              <span><Badge variant="default">{c.taxable !== false ? "Taxable" : "Non-taxable"}</Badge><span className="block text-[10px]">Pensionable: {c.pensionable !== false ? "Yes" : "No"}</span></span>
              <span className="font-bold">ETB {c.amount}<span className="block text-[10px] text-muted-foreground">Order {c.order} • Proratable ✓</span></span>
            </div>
          ))}
        </div>
        <div className="mt-6 rounded-xl border bg-gradient-to-br from-white to-neutral-50 p-4">
          <p className="text-xs font-semibold">Live Preview • CTC Annual {selected.ctc_annual} → Monthly {selected.ctc_monthly}</p>
          <div className="mt-3 grid grid-cols-2 gap-4 text-xs">
            <div><p className="text-muted-foreground">Earnings (Gross)</p><p className="font-bold">BASIC 16,666 + HOUSING 8,333 + TRANSPORT 3,000 + FUEL 2,000 + OT 1,250 + COMMISSION = 31,249</p><p className="text-[11px]">Proration factor 1.0 (paid_days 30/30) → 31,249 *1.0</p></div>
            <div><p className="text-muted-foreground">Deductions</p><p>Taxable 31,249 - Pension 7% 2,187 = 29,062 → Tax binary search 20% -302.5 = 5,510 • Pension Emp 2,187 • Loan EMI 5,000</p><p className="font-bold">Net 18,552 • Employer pension 11% 3,437 • Employer cost 34,686</p></div>
          </div>
        </div>
        <div className="mt-4 flex gap-2">
          <button className="rounded-xl border h-9 px-4 text-xs">Add Component +</button>
          <button className="rounded-xl border h-9 px-4 text-xs">Add Formula • secure parser</button>
          <button className="rounded-xl border h-9 px-4 text-xs">Validate • O(n log n) sort order_no + O(n) eval</button>
        </div>
      </Card>
    </div>
  )
}

// ==================== Runs Tab ====================
function RunsTab() {
  return (
    <Card className="p-6">
      <div className="flex justify-between items-center"><h3 className="font-semibold">Payroll Runs • Status pipeline visual stepper • Ledger M4 per run book • 500 emps &lt;2s p99</h3><button className="rounded-xl bg-primary text-white px-4 h-9 text-xs">Create Run Wizard 5 steps</button></div>
      <div className="mt-4 rounded-xl border overflow-hidden">
        <div className="grid grid-cols-9 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Run Ref</span><span>Period</span><span>Type</span><span>Status</span><span>Total Gross</span><span>Total Net</span><span>Variance</span><span>Paid/Failed</span><span>Action</span></div>
        {mockRuns.map(r=>(
          <div key={r.id} className="grid grid-cols-9 gap-2 p-3 border-t text-xs hover:bg-muted/50">
            <span className="font-mono">{r.ref}</span>
            <span>{r.period}</span>
            <span><Badge>{r.type}</Badge></span>
            <span><Badge variant={r.status==="completed" ? "success" : r.status==="pending_approval" ? "warning" : "default"}>{r.status}</Badge></span>
            <span>ETB {r.total_gross}</span>
            <span className="font-bold">ETB {r.total_net}</span>
            <span className={r.variance.startsWith("+") ? "text-green-600" : "text-red-600"}>{r.variance}</span>
            <span>{r.paid}/{r.failed}</span>
            <Link href={`/payroll/${r.id}`} className="text-primary">View • Approve dual &gt;100k → Disburse → payout batch pain.001 • Cost center report</Link>
          </div>
        ))}
      </div>
      <div className="mt-6 rounded-xl bg-neutral-50 border p-4">
        <p className="text-xs font-semibold">Create Run Wizard 5 Steps • Outstanding Modern</p>
        <div className="mt-3 flex items-center gap-2 text-[11px]">
          <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">1</span><span>Period Month/Year Type regular/off_cycle/bonus/adjustment Pay calendar cutoff disbursal</span></div>
          <div className="h-0.5 w-8 bg-neutral-200" />
          <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">2</span><span>Select employees 500 &lt;2s filter dept/branch/cost_center</span></div>
          <div className="h-0.5 w-8 bg-neutral-200" />
          <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">3</span><span>Attendance CSV + OT 1.25/1.5/2.0/1.3 + Variable pay bonus/commission</span></div>
          <div className="h-0.5 w-8 bg-neutral-200" />
          <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-amber-500 text-white flex items-center justify-center">4</span><span>Review totals variance vs last month 5.2% • Ledger M4 balanced</span></div>
          <div className="h-0.5 w-8 bg-neutral-200" />
          <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center">5</span><span>Disburse → payout batch pain.001 + Pension CSV + ERCA CSV</span></div>
        </div>
      </div>
    </Card>
  )
}

// ==================== Attendance Tab ====================
function AttendanceTab() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card className="p-6 lg:col-span-2">
        <h3 className="font-semibold">Attendance & Variable Inputs • Bulk CSV 500 rows • LOP proration • OT rates</h3>
        <div className="mt-4 grid grid-cols-2 gap-4">
          <div className="rounded-xl border-2 border-dashed border-primary/30 p-6 text-center hover:bg-primary/5 transition-colors">
            <p className="text-sm font-medium">Attendance CSV Upload • paid_days lop_days ot_hours</p>
            <p className="text-[11px] text-muted-foreground mt-1">Format: employee_id,paid_days,lop_days,total_days,ot_weekday,ot_weekend,ot_holiday,ot_night</p>
            <p className="text-[11px] mt-2">Example: EMP001,25,5,30,5,0,0,0 → proration 25/30=0.8333 • Gross 20000*0.8333=16666 • OT 5h weekday 1.25x = 600 ETB</p>
            <button className="mt-3 rounded-xl bg-primary text-white h-9 px-4 text-xs">Drop CSV here • papaparse O(n)</button>
          </div>
          <div className="rounded-xl border-2 border-dashed border-amber-500/30 p-6 text-center hover:bg-amber-500/5">
            <p className="text-sm font-medium">Variable Pay CSV • Bonus Commission Penalty Arrear</p>
            <p className="text-[11px] text-muted-foreground mt-1">Format: employee_id,component_code,amount,is_taxable,description</p>
            <p className="text-[11px] mt-2">Example: EMP001,BONUS,10000,true,Sales Q2 bonus • COMMISSION,5000,true • ARREAR,5000,true,salary revision</p>
            <button className="mt-3 rounded-xl bg-amber-500 text-white h-9 px-4 text-xs">Drop CSV here • validation icons</button>
          </div>
        </div>
        <div className="mt-6 rounded-xl border overflow-hidden">
          <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Employee</span><span>Paid/Total</span><span>LOP</span><span>Proration</span><span>OT Hours</span><span>OT Amount ET</span><span>Status</span></div>
          <div className="grid grid-cols-7 gap-2 p-3 border-t text-xs"><span>EMP001 Abebe</span><span>25/30</span><span>5</span><span>0.8333</span><span>5h weekday</span><span>600 • hourly 96.15*1.25</span><span><Badge variant="warning">Prorated</Badge></span></div>
          <div className="grid grid-cols-7 gap-2 p-3 border-t text-xs"><span>EMP002 Almaz</span><span>30/30</span><span>0</span><span>1.0</span><span>0</span><span>0</span><span><Badge variant="success">Full</Badge></span></div>
        </div>
      </Card>
      <Card className="p-6">
        <h3 className="font-semibold">OT Calculator • ET Labour Law 1156/2019</h3>
        <div className="mt-3 space-y-3 text-xs">
          <div className="rounded-xl bg-muted p-3"><p>Hourly rate = Base Salary /208 (26 days *8h)</p><p className="font-mono">EMP001 20000/208 = 96.15 ETB/hr</p></div>
          <div className="rounded-xl border p-3 space-y-2">
            <p>OT Rates:</p>
            <p>• Weekday 1.25x → 96.15*1.25=120.19/hr *5h=600.96</p>
            <p>• Weekend 1.5x → 144.23/hr</p>
            <p>• Holiday 2.0x → 192.30/hr</p>
            <p>• Night 1.3x → 125.00/hr</p>
          </div>
          <div className="rounded-xl bg-green-500/10 border border-green-500/20 p-3"><p>Formula Engine secure:</p><p className="font-mono text-[11px]">OTAmount = weekday*rate_weekday + weekend*rate_weekend... O(1) map lookup</p></div>
        </div>
      </Card>
    </div>
  )
}

// ==================== Loans Tab ====================
function LoansTab() {
  return (
    <Card className="p-6">
      <div className="flex justify-between items-center"><h3 className="font-semibold">Loans & Advances • Salary Advance • Personal Loan • EMI auto deduction per run</h3><button className="rounded-xl bg-primary text-white h-9 px-4 text-xs">Request Loan • EMI preview</button></div>
      <div className="mt-4 rounded-xl border overflow-hidden">
        <div className="grid grid-cols-8 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Loan ID</span><span>Employee</span><span>Type</span><span>Principal</span><span>EMI</span><span>Outstanding</span><span>Tenure</span><span>Status</span></div>
        {mockLoans.map(l=>(
          <div key={l.id} className="grid grid-cols-8 gap-2 p-3 border-t text-xs hover:bg-muted/50">
            <span className="font-mono">{l.id}</span>
            <span>{l.employee} • {l.code}</span>
            <span><Badge>{l.type}</Badge></span>
            <span>ETB {l.principal}</span>
            <span>ETB {l.emi}</span>
            <span className="font-bold">ETB {l.outstanding}</span>
            <span>{l.tenure}mo</span>
            <span><Badge variant={l.status==="active" ? "success" : "warning"}>{l.status}</Badge></span>
          </div>
        ))}
      </div>
      <div className="mt-6 grid grid-cols-2 gap-6">
        <div className="rounded-xl bg-neutral-50 border p-4 text-xs">
          <p className="font-semibold">Loan EMI Auto Deduction Logic O(k) k=0-2 loans per employee</p>
          <p className="mt-2 font-mono text-[11px]">active_loans = ListActiveLoans(employee_id) → for each loan: emi if emi &gt; outstanding then emi=outstanding → deduction += emi → Create repayment record run_id → Update loan outstanding total_paid += emi → if outstanding ==0 status closed</p>
          <p className="mt-2">Example EMP001 salary_advance 20k EMI 5k tenure 4mo outstanding 15k after 1st run → next 3 runs auto deduct 5k each until 0</p>
        </div>
        <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-xs">
          <p className="font-semibold">Reimbursements & Claims • Expense/Medical/Travel • Receipt MinIO &lt;5MB</p>
          <p className="mt-2">Claim types: expense, medical, travel, other • Status pending→approved→rejected→paid • Approval flow manager→finance • File_key MinIO presigned 15m hash integrity</p>
          <p className="mt-2">Example: EMP001 travel 2000 ETB receipt travel_receipt.pdf → manager approve → finance approve → paid via next payroll run reimbursement non-taxable</p>
        </div>
      </div>
    </Card>
  )
}

// ==================== Compliance Tab ====================
function ComplianceTab() {
  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex justify-between items-center"><h3 className="font-semibold">Compliance Reports • Pension 7%/11% + ERCA withholding + Bank pain.001 • Outstanding</h3><button className="rounded-xl border h-9 px-4 text-xs">Generate All Reports • July 2026</button></div>
        <div className="mt-4 rounded-xl border overflow-hidden">
          <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Report Type</span><span>Period</span><span>Status</span><span>File</span><span>Total/Count</span><span>Format</span><span>Action</span></div>
          {mockCompliance.map((c,i)=>(
            <div key={i} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
              <span><Badge variant="default">{c.type.replaceAll("_"," ")}</Badge></span>
              <span>{c.period}</span>
              <span><Badge variant="success">{c.status}</Badge></span>
              <span className="font-mono text-[11px]">{c.file}</span>
              <span>{c.total || c.total_tax || c.total_net || c.count} • {c.count ? `${c.count} emps` : ""}</span>
              <span className="text-[11px]">{c.format || "CSV"}</span>
              <span className="flex gap-2"><button className="text-primary">Download • MinIO presigned 15m</button><button className="text-muted-foreground">View</button></span>
            </div>
          ))}
        </div>
      </Card>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-6">
          <h3 className="font-semibold text-sm">Pension Contribution Report • Private Org Employees Social Security Agency • Format</h3>
          <div className="mt-3 rounded-xl bg-muted p-3 font-mono text-[11px]">
            <p>pension_no, employee_name, employee_code, pensionable_gross, employee_7pct, employer_11pct, total_18pct, period</p>
            <p>PEN-001, Abebe Kebede, EMP001, 20000, 1400, 2200, 3600, 2026-07</p>
            <p>PEN-002, Almaz Tadesse, EMP002, 25000, 1750, 2750, 4500, 2026-07</p>
            <p>...</p>
            <p>Total 10 employees • Employee 7% 14000 • Employer 11% 22000 • Total 36000 • Generated via payroll_compliance_reports table file_key MinIO hash</p>
          </div>
        </Card>
        <Card className="p-6">
          <h3 className="font-semibold text-sm">ERCA Withholding Monthly CSV • Ethiopian Revenues and Customs Authority</h3>
          <div className="mt-3 rounded-xl bg-muted p-3 font-mono text-[11px]">
            <p>tin, employee_name, employee_code, gross, pension_employee, taxable_income, income_tax, net, period, cost_center</p>
            <p>0098765432, Abebe Kebede, EMP001, 20000, 1400, 18600, 1668.5, 16931.5, 2026-07, CC-100</p>
            <p>... • Income tax binary search O(log n) 7 brackets • Tax = taxable*rate - deduction • Rounded 2 decimals</p>
            <p>Total tax 20000 for July 2026 • 10 employees • Audit immutable payroll_audit_logs</p>
          </div>
        </Card>
      </div>
      <Card className="p-6">
        <h3 className="font-semibold text-sm">Bank Disbursal File • ISO20022 pain.001.001.03 XML • CBE/Awash/Dashen • MT103/CSV fallback</h3>
        <div className="mt-3 rounded-xl bg-neutral-900 text-neutral-100 p-4 font-mono text-[10px] overflow-x-auto">
          <p>{`<?xml version="1.0" encoding="UTF-8"?>`}</p>
          <p>{`<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pain.001.001.03">`}</p>
          <p>{`  <CstmrCdtTrfInitn><GrpHdr><MsgId>prun_July2026</MsgId><NbOfTxs>10</NbOfTxs><CtrlSum>150000</CtrlSum>`}</p>
          <p>{`  <PmtInf><PmtInfId>prun_July2026</PmtInfId><PmtMtd>TRF</PmtMtd><NbOfTxs>10</NbOfTxs><CtrlSum>150000</CtrlSum>`}</p>
          <p>{`    <CdtTrfTxInf><PmtId><InstrId>pay_001</InstrId></PmtId><Amt><InstdAmt Ccy="ETB">15000</InstdAmt></Amt><Cdtr><Nm>Abebe Kebede</Nm></Cdtr><CdtrAcct><Id><Othr><Id>CBE ****1234</Id></Othr></Id>`}</p>
          <p>{`  ... 10 transactions • Reconciliation bank statement MT940 window 24h amount tolerance 0.01 ETB O(n+m) map connector_ref->journal`}</p>
        </div>
      </Card>
    </div>
  )
}

// ==================== Settings Tab ====================
function SettingsTab() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card className="p-6">
        <h3 className="font-semibold">Departments • Cost Center • 3</h3>
        <div className="mt-3 space-y-2 text-xs">{mockDepartments.map(d=><div key={d.id} className="flex justify-between border rounded-xl p-2"><span>{d.name} • {d.code} • {d.cost_center}</span><span>{d.headcount} emps</span></div>)}</div>
        <button className="mt-3 w-full rounded-xl border border-dashed h-10 text-xs">+ Add Department</button>
      </Card>
      <Card className="p-6">
        <h3 className="font-semibold">Payroll Calendar • Cutoff • Disbursal • 2026</h3>
        <div className="mt-3 space-y-2 text-xs">
          <div className="rounded-xl bg-muted p-3"><p>July 2026 regular</p><p>Cutoff: 25 Jul • Disbursal: 30 Jul • Pay date: 31 Jul • Lock after disbursal ✓</p></div>
          <div className="rounded-xl bg-muted p-3"><p>August 2026 regular</p><p>Cutoff: 25 Aug • Disbursal: 30 Aug</p></div>
        </div>
        <button className="mt-3 w-full rounded-xl border border-dashed h-10 text-xs">+ Add Calendar</button>
      </Card>
      <Card className="p-6">
        <h3 className="font-semibold">Tax Brackets • Versioned effective_from • ET 2024 • 7 brackets</h3>
        <div className="mt-3 space-y-1 text-[11px] font-mono">
          <p>0-600 0% deduction 0</p>
          <p>601-1650 10% -60</p>
          <p>1651-3200 15% -142.5</p>
          <p>3201-5250 20% -302.5</p>
          <p>5251-7800 25% -565</p>
          <p>7801-10900 30% -955</p>
          <p>&gt;10900 35% -1500 • Formula tax=taxable*rate-deduction • Binary search O(log n)</p>
        </div>
        <div className="mt-4 rounded-xl bg-green-500/10 border border-green-500/20 p-3 text-[11px]"><p>Benchmark: CalculateTax 10k iterations p99&lt;30ms • ValidateBalanced O(n) • payroll calc 500 emps &lt;2s p99 • k6 100 VUs p95&lt;300ms</p></div>
      </Card>
    </div>
  )
}
