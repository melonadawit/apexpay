"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type CorporateCard } from "@/lib/api/client"

const columns: Column<CorporateCard>[] = [
  { key: "card_number_masked", label: "Card" },
  { key: "card_type", label: "Type" },
  { key: "card_network", label: "Network" },
  { key: "cardholder_name", label: "Cardholder" },
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
    key: "credit_limit",
    label: "Limit",
    render: (c) => <span>ETB {c.credit_limit}</span>,
  },
  {
    key: "forex_markup_percent",
    label: "Forex",
    render: (c) => <span>{c.forex_markup_percent}%</span>,
  },
]

export default function CorporateCardsPage() {
  return (
    <BankingPage
      title="Corporate Cards • የኩባንያ ካርዶች"
      subtitle="Virtual + physical corporate cards with spending controls and forex markup."
      columns={columns}
      loader={() => api.banking.corporateCards()}
    />
  )
}
