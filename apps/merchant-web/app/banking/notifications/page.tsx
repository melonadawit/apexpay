"use client"
import { BankingPage, type Column } from "@/components/banking/page"
import { api, type Notification } from "@/lib/api/client"

const columns: Column<Notification>[] = [
  { key: "type", label: "Type" },
  { key: "title", label: "Title" },
  { key: "message", label: "Message" },
  {
    key: "is_read",
    label: "Read",
    render: (n) => (n.is_read ? "✓" : <span className="px-2 py-0.5 rounded-full bg-blue-500/15 text-blue-700 text-[11px]">New</span>),
  },
  { key: "created_at", label: "Created" },
]

export default function NotificationsPage() {
  return (
    <BankingPage
      title="Notifications • ማሳወቂያዎች"
      subtitle="Bulk payouts, payroll approvals, escrow, forex alerts and more."
      columns={columns}
      loader={() => api.banking.notifications()}
    />
  )
}
