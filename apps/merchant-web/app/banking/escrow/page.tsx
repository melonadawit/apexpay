"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type EscrowAccount } from "@/lib/api/client"

const columns: Column<EscrowAccount>[] = [
  { key: "account_number", label: "Escrow Account" },
  { key: "account_name", label: "Name" },
  {
    key: "amount",
    label: "Amount",
    render: (e) => <span className="font-semibold">ETB {e.amount}</span>,
  },
  {
    key: "status",
    label: "Status",
    render: (e) => (
      <span className={`px-2 py-0.5 rounded-full text-[11px] ${e.status === "released" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>
        {e.status}
      </span>
    ),
  },
  {
    key: "platform_fee",
    label: "Platform Fee",
    render: (e) => <span>ETB {e.platform_fee}</span>,
  },
  {
    key: "seller_amount",
    label: "Seller",
    render: (e) => <span>ETB {e.seller_amount}</span>,
  },
]

export default function EscrowPage() {
  return (
    <BankingPage
      titleEn="Escrow" titleAm="ኢስክሮ"
      subtitle="Automated marketplace escrow that holds and releases funds under defined conditions."
      columns={columns}
      loader={() => api.banking.escrow()}
    />
  )
}
