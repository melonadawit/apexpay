"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type BankVerification } from "@/lib/api/client"

const columns: Column<BankVerification>[] = [
  { key: "id", label: "ID" },
  { key: "bank_code", label: "Bank Code" },
  { key: "account_number_masked", label: "Account" },
  { key: "account_name", label: "Account Name" },
  { key: "verification_method", label: "Method" },
  {
    key: "status",
    label: "Status",
    render: (v) => (
      <span
        className={`px-2 py-0.5 rounded-full text-[11px] border ${
          v.status === "verified"
            ? "bg-green-500/15 text-green-700"
            : v.status === "pending"
            ? "bg-amber-500/15 text-amber-700"
            : "bg-red-500/15 text-red-700"
        }`}
      >
        {v.status}
      </span>
    ),
  },
  { key: "created_at", label: "Created" },
]

export default function BankVerificationPage() {
  return (
    <BankingPage
      titleEn="Bank Account Verification" titleAm="ባንክ ሂሳብ ማረጋገጫ"
      subtitle="Penny testing / fund account validation — 1 ETB deposit returns validated bank details."
      columns={columns}
      loader={() => api.banking.bankVerifications()}
    />
  )
}
