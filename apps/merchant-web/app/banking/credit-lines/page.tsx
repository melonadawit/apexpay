"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type CreditLine } from "@/lib/api/client"

const columns: Column<CreditLine>[] = [
  { key: "id", label: "ID" },
  {
    key: "credit_limit",
    label: "Limit",
    render: (c) => <span className="font-semibold">ETB {c.credit_limit}</span>,
  },
  {
    key: "available_credit",
    label: "Available",
    render: (c) => <span>ETB {c.available_credit}</span>,
  },
  {
    key: "utilized_credit",
    label: "Utilized",
    render: (c) => <span>ETB {c.utilized_credit}</span>,
  },
  { key: "interest_rate", label: "Interest %" },
  {
    key: "status",
    label: "Status",
    render: (c) => (
      <span className={`px-2 py-0.5 rounded-full text-[11px] ${c.status === "active" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>
        {c.status}
      </span>
    ),
  },
  {
    key: "credit_score",
    label: "Score",
    render: (c) => (c.credit_score ? String(c.credit_score) : "—"),
  },
]

export default function CreditLinesPage() {
  return (
    <BankingPage
      titleEn="Credit Lines" titleAm="የብድር መስመሮች"
      subtitle="Collateral-free credit lines based on TPV / payroll scoring, 18% p.a."
      columns={columns}
      loader={() => api.banking.creditLines()}
    />
  )
}
