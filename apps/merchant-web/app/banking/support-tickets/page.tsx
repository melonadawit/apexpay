"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type SupportTicket } from "@/lib/api/client"

const columns: Column<SupportTicket>[] = [
  { key: "id", label: "ID" },
  { key: "subject", label: "Subject" },
  { key: "priority", label: "Priority" },
  {
    key: "status",
    label: "Status",
    render: (t) => (
      <span className={`px-2 py-0.5 rounded-full text-[11px] ${
        t.status === "resolved" ? "bg-green-500/15 text-green-700"
        : t.status === "open" ? "bg-amber-500/15 text-amber-700"
        : "bg-blue-500/15 text-blue-700"
      }`}>
        {t.status}
      </span>
    ),
  },
  { key: "assigned_to", label: "Assigned" },
  { key: "created_at", label: "Created" },
]

export default function SupportTicketsPage() {
  return (
    <BankingPage
      title="Support Tickets • የድጋፍ ጥያቄዎች"
      subtitle="Priority support with SLA and dedicated relationship manager."
      columns={columns}
      loader={() => api.banking.supportTickets()}
    />
  )
}
