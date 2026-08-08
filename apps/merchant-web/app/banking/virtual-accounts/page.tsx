"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type VirtualAccount } from "@/lib/api/client"

const columns: Column<VirtualAccount>[] = [
  { key: "virtual_account_number", label: "Virtual Account #" },
  { key: "customer_id", label: "Customer" },
  { key: "purpose", label: "Purpose" },
  {
    key: "status",
    label: "Status",
    render: (v) => (
      <span className={`px-2 py-0.5 rounded-full text-[11px] ${v.status === "active" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>
        {v.status}
      </span>
    ),
  },
  { key: "bank_code", label: "Bank" },
  { key: "created_at", label: "Created" },
]

export default function VirtualAccountsPage() {
  return (
    <BankingPage
      title="Virtual Accounts • Smart Collect"
      subtitle="Automatically reconcile incoming payments using virtual accounts & UPI-IDs."
      columns={columns}
      loader={() => api.banking.virtualAccounts()}
    />
  )
}
