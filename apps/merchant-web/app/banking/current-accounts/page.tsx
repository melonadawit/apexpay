"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type CurrentAccount } from "@/lib/api/client"

const columns: Column<CurrentAccount>[] = [
  { key: "account_number", label: "Account #" },
  { key: "account_name", label: "Name" },
  { key: "account_type", label: "Type" },
  { key: "partner_bank_name", label: "Bank" },
  {
    key: "status",
    label: "Status",
    render: (a) => (
      <span className={`px-2 py-0.5 rounded-full text-[11px] ${a.status === "active" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>
        {a.status}
      </span>
    ),
  },
  {
    key: "balance",
    label: "Balance",
    render: (a) => <span className="font-semibold">ETB {a.balance}</span>,
  },
  {
    key: "is_primary",
    label: "Primary",
    render: (a) => (a.is_primary ? "✓" : "—"),
  },
]

export default function CurrentAccountsPage() {
  return (
    <BankingPage
      titleEn="Current Accounts" titleAm="የአሁኑ ሂሳቦች"
      subtitle="Real partner-bank current accounts (CBE/Awash/Dashen) with balances and card/cheque status."
      columns={columns}
      loader={() => api.banking.currentAccounts()}
    />
  )
}
