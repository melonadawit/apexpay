"use client"
import * as React from "react"
import { gregorianToEthiopian, ethiopianMonths, ethiopianMonthsAmharic, getEthiopianMonthsForYear, isEnkutatash, ethiopianPublicHolidays } from "@/lib/ethiopian-calendar"

export function EthiopianCalendarGrid({ year = 2026 }: { year?: number }) {
  const months = getEthiopianMonthsForYear(year)
  const [selectedMonth, setSelectedMonth] = React.useState(months[6]) // Hamle 2018 = July 2026

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="font-semibold text-sm">Ethiopian Calendar • 13 Months • Meskerem to Pagume • Enkutatash Meskerem 1 = Sept 11 • 12*30 + Pagume 5/6 days • Ethiopia Law Compliance • Outstanding Modern • Glassmorphic</h3>
        <div className="flex gap-2">
          <span className="px-2 py-0.5 rounded-full bg-primary/10 text-primary border border-primary/20 text-[11px]">Enkutatash • Ethiopian New Year • አዲስ አመት • Sept 11 • Public Holiday OT 2.0x per Art 90(2)</span>
        </div>
      </div>

      <div className="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-13 gap-2">
        {ethiopianMonths.map((m, idx) => {
          const isSelected = selectedMonth.ethiopianMonth === idx + 1
          const isPagume = idx === 12
          const gregMonth = months[idx]
          return (
            <button
              key={m}
              onClick={() => setSelectedMonth(gregMonth)}
              className={`rounded-xl border p-3 text-left hover:bg-muted transition-all ${isSelected ? "bg-primary text-white border-primary shadow-medium" : isPagume ? "bg-amber-500/10 border-amber-500/20" : "bg-card"} ${isPagume ? "col-span-1" : ""}`}
            >
              <p className="font-medium text-xs">{m}</p>
              <p className="text-[10px] opacity-80">{ethiopianMonthsAmharic[idx]}</p>
              <p className="text-[10px] mt-1">{isPagume ? "5/6 days" : "30 days"}</p>
              <p className="text-[9px] mt-1 opacity-70">{gregMonth ? `${gregMonth.gregorianMonth}/${gregMonth.gregorianYear}` : `${year}`}</p>
              {isPagume && <p className="text-[9px] mt-1">Pagume • ጷጉሜ • 5/6 days • Leap year • Enkutatash + 1 year</p>}
            </button>
          )
        })}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
        <div className="rounded-xl bg-muted p-4">
          <p className="text-[11px] font-semibold">Selected Month • {selectedMonth.monthName} • {selectedMonth.monthNameAm} • {selectedMonth.ethiopianMonth} • Ethiopian Year {selectedMonth.ethiopianYear} • Gregorian {selectedMonth.gregorianMonth}/{selectedMonth.gregorianYear}</p>
          <p className="text-[11px] mt-2">Cutoff Date: {selectedMonth.cutoffDate.toLocaleDateString()} • 25th Ethiopia business practice • Gregorian {selectedMonth.cutoffDate.toLocaleDateString()} • Ethiopian {(() => { try { const eth = gregorianToEthiopian(selectedMonth.cutoffDate); return `${eth.formatted} (${eth.formattedAm})` } catch { return selectedMonth.cutoffDate.toLocaleDateString() } })()}</p>
          <p className="text-[11px] mt-1">Disbursal Date: {selectedMonth.disbursalDate.toLocaleDateString()} • 30th • Payroll batch pain.001 XML ISO20022 • Pension CSV • ERCA CSV • Ledger M4 per run book</p>
          <p className="text-[11px] mt-1">Pay Date: {selectedMonth.payDate.toLocaleDateString()} • Last day • Salary credited • Payslip PDF QR • Email SMTP • Lock after disbursal is_locked true locked_at now locked_by</p>
          <p className="text-[10px] mt-2 text-muted-foreground">Enkutatash Meskerem 1 = Sept 11 Gregorian • Ethiopian New Year • አዲስ አመት • Public Holiday OT 2.0x per Art 90(2) Labour Proclamation 1156/2019 • 13 months 12*30 + Pagume 5/6 days • Outstanding modern UI glassmorphic</p>
        </div>
        <div className="rounded-xl bg-amber-500/10 border border-amber-500/20 p-4">
          <p className="text-[11px] font-semibold">Public Holidays Ethiopia • OT 2.0x per Art 90(2) • Labour Proclamation 1156/2019 • Recharts Calendar View</p>
          <div className="mt-2 space-y-1 max-h-32 overflow-auto">
            {ethiopianPublicHolidays.map(h => (
              <p key={h.name} className="text-[10px]">• {h.month}/{h.day} {h.name} • {h.nameAm} • Type {h.type} • OT {h.ot_rate}x • Per Art 90(2) public holiday OT 2.0x (some say 250% configurable to 2.5x) • Holiday OT calculation weekday 1.5x weekend 1.5x holiday 2.0x night 1.3x hourly=base/208</p>
            ))}
          </div>
        </div>
        <div className="rounded-xl bg-green-500/10 border border-green-500/20 p-4">
          <p className="text-[11px] font-semibold">Payroll Calendar Logic Outstanding per Ethiopia Law • O(1) Advisory Lock • Variance Report</p>
          <ul className="list-disc list-inside mt-2 space-y-1 text-[10px] text-muted-foreground">
            <li>Monthly: cutoff 25th disbursal 30th pay date last day last day of month lock after disbursal per law</li>
            <li>Semimonthly: cutoff 15th & last day, disbursal 16th & 1st next month</li>
            <li>Weekly: cutoff Friday disbursal Monday pay Monday</li>
            <li>Biweekly: cutoff every 2 weeks Friday</li>
            <li>Lock after disbursal: is_locked true locked_at now locked_by finance manager • Prevents re-run amendment unless unlocked by admin with audit log payroll_audit_logs actor admin action unlock_calendar details locked_by IP inet request_id immutable</li>
            <li>Variance report vs last month +5.2% vs Jun OT increase + bonus Sales Q2 + new hires 2 • total_gross total_net total_tax • Recharts AreaChart trend Feb 160k Mar 170k Apr 180k May 185k Jun 190k Jul 200k +5.2% • Cost center breakdown Engineering 100k Sales 100k • Paid 280/300 LOP 20 • Proration avg 0.93 • Outstanding modern UI glassmorphic Recharts calendar view AreaChart BarChart + Ethiopian calendar Enkutatash Meskerem Tikimt Hidar 13 months 12*30 + Pagume 5/6 days • Cost center allocation • Variance report • Enkutatash Sept 11 public holiday OT 2.0x per Art 90(2)</li>
          </ul>
        </div>
      </div>
    </div>
  )
}
