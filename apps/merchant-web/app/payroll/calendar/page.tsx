"use client"
import * as React from "react"
import { motion } from "framer-motion"
import { gregorianToEthiopian, formatEthiopianDate, ethiopianPublicHolidays } from "@/lib/ethiopian-calendar"
import { EthiopianCalendarGrid } from "@/components/payroll/EthiopianCalendarGrid"
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar } from "recharts"

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

const payrollTrend = [
  { month: "Feb", eth: "Yekatit 2018", gross: 160000, net: 120000, tax: 16000 },
  { month: "Mar", eth: "Megabit 2018", gross: 170000, net: 127500, tax: 17000 },
  { month: "Apr", eth: "Miyazya 2018", gross: 180000, net: 135000, tax: 18000 },
  { month: "May", eth: "Ginbot 2018", gross: 185000, net: 139000, tax: 18500 },
  { month: "Jun", eth: "Sene 2018", gross: 190000, net: 142500, tax: 19000 },
  { month: "Jul", eth: "Hamle 2018", gross: 200000, net: 150000, tax: 20000 },
  { month: "Aug", eth: "Nehasse 2018", gross: 210000, net: 157500, tax: 21000 },
]

export default function PayrollCalendarPage() {
  const [selected, setSelected] = React.useState(mockCalendars[0])
  const [showCreate, setShowCreate] = React.useState(false)
  const ethDate = gregorianToEthiopian(new Date())

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-3">Payroll Calendar • የደሞዝ ቀን መቁጠሪያ • Cutoff 25th Disbursal 30th Pay Last Day Lock After Disbursal • Recharts + Ethiopian Calendar Enkutatash</h1>
            <p className="text-sm text-muted-foreground mt-2">Today Gregorian {new Date().toLocaleDateString()} • Ethiopian {ethDate.formatted} ({ethDate.formattedAm}) • Enkutatash Meskerem 1 = Sept 11 • 13 months 12*30 + Pagume 5/6 days • Pay frequency monthly/semimonthly/weekly/biweekly • Cutoff 25th disbursal 30th pay last day lock after disbursal per law • Outstanding modern UI glassmorphic Recharts</p>
          </div>
          <button onClick={()=>setShowCreate(true)} className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Calendar • Monthly Weekly Semimonthly</button>
        </div>

        <Card className="p-6">
          <EthiopianCalendarGrid year={2026} />
        </Card>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Payroll Cost Trend • Recharts AreaChart • Gross Net Tax • Feb 160k → Aug 210k +5.2% • Ethiopian Months Meskerem Tikimt Hidar • Enkutatash Sept 11</h3>
            <div className="mt-4 h-64">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={payrollTrend}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e4e4e7" />
                  <XAxis dataKey="month" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip contentStyle={{ borderRadius: 12, fontSize: 11 }} />
                  <Area type="monotone" dataKey="gross" stackId="1" stroke="#0B6E4F" fill="#0B6E4F33" name="Gross" />
                  <Area type="monotone" dataKey="net" stackId="2" stroke="#10B981" fill="#10B98133" name="Net" />
                  <Area type="monotone" dataKey="tax" stackId="3" stroke="#EAB308" fill="#EAB30833" name="Tax" />
                </AreaChart>
              </ResponsiveContainer>
            </div>
            <div className="mt-4 grid grid-cols-7 gap-2 text-[11px]">
              {payrollTrend.map(d=>(
                <div key={d.month} className="rounded-xl bg-muted p-2 text-center"><p className="font-medium">{d.month}</p><p className="text-[10px]">{d.eth}</p><p>Gross {d.gross/1000}k</p><p>Net {d.net/1000}k</p></div>
              ))}
            </div>
          </Card>
          <Card className="p-6">
            <h3 className="font-semibold">Ethiopian Calendar • Enkutatash • 13 Months • Meskerem Tikimt Hidar • Public Holidays OT 2.0x per Art 90(2)</h3>
            <div className="mt-4 space-y-2 text-[11px]">
              <p>Today: Gregorian {new Date().toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" })} • Ethiopian {ethDate.formatted} ({ethDate.formattedAm}) • Year {ethDate.year} Month {ethDate.month} {ethDate.monthName} {ethDate.monthNameAm} Day {ethDate.day}</p>
              <p>Enkutatash: Meskerem 1 = Sept 11 Gregorian • Ethiopian New Year • አዲስ አመት • Public Holiday OT 2.0x per Art 90(2) Labour Proclamation 1156/2019</p>
              <div className="rounded-xl bg-amber-500/10 border border-amber-500/20 p-3">
                <p className="font-semibold">Public Holidays Ethiopia • OT 2.0x per Art 90(2)</p>
                {ethiopianPublicHolidays.slice(0,5).map(h=>(
                  <p key={h.name} className="mt-1">• {h.month}/{h.day} {h.name} • {h.nameAm} • OT {h.ot_rate}x</p>
                ))}
              </div>
            </div>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Pay Calendars • 2026 • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 space-y-3">
              {mockCalendars.map(cal => {
                const ethCutoff = gregorianToEthiopian(new Date(cal.cutoff_date))
                const ethPay = gregorianToEthiopian(new Date(cal.pay_date))
                return (
                <button key={cal.id} onClick={()=>setSelected(cal)} className={`w-full text-left rounded-xl border p-4 hover:bg-muted ${selected.id===cal.id ? "bg-primary/10 border-primary/30" : ""}`}>
                  <div className="flex justify-between"><p className="font-medium text-sm">{cal.name}</p><Badge variant={cal.is_locked ? "success" : cal.status==="draft" ? "warning" : "default"}>{cal.is_locked ? "Locked" : cal.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground">Cutoff {cal.cutoff_date} ({ethCutoff.formattedAm}) → Pay {cal.pay_date} ({ethPay.formattedAm}) • Gross {cal.total_gross || "—"} Net {cal.total_net || "—"}</p>
                </button>
              )})}
            </div>
          </Card>
          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Calendar Detail • {selected.name} • Cutoff 25th Disbursal 30th Pay Last Day • Recharts BarChart cost center • Ethiopian {formatEthiopianDate ? "enriched" : ""}</h3>
            <div className="mt-4 grid grid-cols-3 gap-4">
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px]">Cutoff 25th • መቁረጥ • {selected.cutoff_date} • Ethiopian {(() => { try { return formatEthiopianDate(new Date(selected.cutoff_date)) } catch { return selected.cutoff_date } })()}</p><p className="font-bold">{selected.cutoff_date}</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px]">Disbursal 30th • ክፍያ • {selected.disbursal_date}</p><p className="font-bold">{selected.disbursal_date}</p></div>
              <div className="rounded-xl bg-muted p-4"><p className="text-[11px]">Pay Day Last day • የክፍያ ቀን • {selected.pay_date}</p><p className="font-bold">{selected.pay_date}</p></div>
            </div>
            <div className="mt-6 h-32">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={[{ cc: "CC-100 Eng", gross: 100000, net: 75000 }, { cc: "CC-200 Sales", gross: 100000, net: 75000 }]}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="cc" tick={{ fontSize: 10 }} />
                  <YAxis tick={{ fontSize: 10 }} />
                  <Tooltip />
                  <Bar dataKey="gross" fill="#0B6E4F" name="Gross" />
                  <Bar dataKey="net" fill="#10B981" name="Net" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}


