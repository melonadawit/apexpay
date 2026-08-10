"use client"
import * as React from "react"
import { Loader2, Inbox, AlertTriangle, Search, RefreshCw } from "lucide-react"
import { getQueue, getExam, reviewMerchant, type QueueItem, type MerchantExam } from "@/lib/admin"

export default function AdminPage() {
  const [queue, setQueue] = React.useState<QueueItem[]>([])
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState("")
  const [selectedId, setSelectedId] = React.useState<string | null>(null)
  const [exam, setExam] = React.useState<MerchantExam | null>(null)
  const [examLoading, setExamLoading] = React.useState(false)
  const [notice, setNotice] = React.useState("")

  const loadQueue = async () => {
    setLoading(true)
    setError("")
    try {
      const q = await getQueue()
      setQueue(q)
      if (q.length && !selectedId) setSelectedId(q[0].merchant_id)
    } catch (e) {
      setError((e as Error).message || "Could not load the onboarding queue.")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    void loadQueue()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  React.useEffect(() => {
    if (!selectedId) return
    setExamLoading(true)
    getExam(selectedId)
      .then(setExam)
      .catch(() => setExam(null))
      .finally(() => setExamLoading(false))
  }, [selectedId])

  const review = async (action: string) => {
    if (!selectedId) return
    setNotice("")
    try {
      await reviewMerchant(selectedId, action, action === "request_info" ? "More information requested." : `${action} by admin`)
      setNotice(`Merchant ${action.replace("_", " ")}`)
      void loadQueue()
    } catch (e) {
      setNotice((e as Error).message || "Action failed.")
    }
  }

  const selected = queue.find((q) => q.merchant_id === selectedId)

  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <header className="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">Admin — Onboarding Queue</h1>
          <p className="text-sm text-muted-foreground">Review, verify, and approve merchant onboarding.</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">{queue.length} merchants</span>
          <button onClick={loadQueue} className="inline-flex items-center gap-1.5 rounded-xl border bg-white px-3 h-9 text-sm hover:bg-muted">
            <RefreshCw className="h-3.5 w-3.5" /> Refresh
          </button>
        </div>
      </header>

      <div className="grid grid-cols-12 gap-6">
        {/* Queue list */}
        <div className="col-span-12 lg:col-span-4 space-y-2">
          <h3 className="font-semibold text-sm">Queue</h3>
          {loading && (
            <div className="flex items-center gap-2 p-4 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" /> Loading…
            </div>
          )}
          {error && (
            <div className="flex flex-col items-start gap-2 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
              <span className="inline-flex items-center gap-1.5 font-medium"><AlertTriangle className="h-4 w-4" /> Unable to load</span>
              <span className="text-xs">{error}</span>
              <button onClick={loadQueue} className="mt-1 inline-flex items-center gap-1 rounded-lg border px-3 py-1.5 text-xs"><RefreshCw className="h-3 w-3" /> Retry</button>
            </div>
          )}
          {!loading && !error && queue.length === 0 && (
            <div className="flex flex-col items-center gap-2 rounded-xl border bg-white p-8 text-center text-muted-foreground">
              <Inbox className="h-8 w-8 text-muted-foreground/50" />
              <p className="text-sm font-medium">No merchants in the queue</p>
            </div>
          )}
          {!loading && !error && queue.map((m) => (
            <button
              key={m.merchant_id}
              onClick={() => setSelectedId(m.merchant_id)}
              className={`w-full text-left rounded-xl border p-3 ${selectedId === m.merchant_id ? "border-primary bg-primary/5" : "border-black/10 bg-white hover:bg-muted"}`}
            >
              <p className="text-sm font-semibold">{m.legal_name}</p>
              <p className="text-xs text-muted-foreground">{m.onboarding_status.replace("_", " ")} • Risk {m.risk_score}</p>
              <div className="mt-1 flex gap-1">
                <span className={`text-[10px] px-2 py-0.5 rounded-full ${m.fayda_verified ? "bg-green-100 text-green-700" : "bg-amber-100 text-amber-700"}`}>
                  Fayda {m.fayda_verified ? "✓" : "pending"}
                </span>
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-neutral-100">{m.risk_tier}</span>
              </div>
            </button>
          ))}
        </div>

        {/* Exam view */}
        <div className="col-span-12 lg:col-span-8">
          {examLoading && (
            <div className="flex items-center gap-2 rounded-xl border bg-white p-8 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" /> Loading merchant…
            </div>
          )}
          {!examLoading && !exam && (
            <div className="flex flex-col items-center gap-2 rounded-xl border bg-white p-8 text-center text-muted-foreground">
              <Search className="h-8 w-8 text-muted-foreground/50" />
              <p className="text-sm font-medium">Select a merchant to view their onboarding review.</p>
            </div>
          )}
          {!examLoading && exam && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-xl border bg-white p-4 space-y-4">
                <h4 className="font-semibold">{exam.legal_name || "Merchant"}</h4>
                <div className="space-y-2 text-sm">
                  <p><span className="text-muted-foreground">Status:</span> {exam.onboarding_status?.replace("_", " ")}</p>
                  <p>
                    <span className="text-muted-foreground">Risk:</span>{" "}
                    <span className={`px-2 py-0.5 rounded-full text-xs ${(exam.risk_score || 0) < 50 ? "bg-green-100" : "bg-amber-100"}`}>
                      {exam.risk_score}/100 {exam.risk_tier}
                    </span>
                  </p>
                  <p><span className="text-muted-foreground">Merchant ID:</span> {exam.merchant_id}</p>
                </div>
                <div className="rounded-xl border p-3 text-xs">
                  <p className="font-semibold">Bank Settlement</p>
                  {(exam.banks ?? []).length === 0 ? (
                    <p className="text-muted-foreground">No bank accounts on record.</p>
                  ) : (
                    exam.banks!.map((b, i) => (
                      <p key={i} className="text-muted-foreground">{b.bank_code || "Bank"} • {b.account_number_masked || "—"}</p>
                    ))
                  )}
                </div>
              </div>

              <div className="rounded-xl border bg-white p-4 space-y-3">
                <h4 className="font-semibold">Compliance Checks</h4>
                {(exam.compliance_checks ?? []).length === 0 ? (
                  <p className="text-sm text-muted-foreground">No checks recorded yet.</p>
                ) : (
                  exam.compliance_checks!.map((c, i) => (
                    <div key={i} className="flex items-center justify-between text-sm border-b last:border-0 py-2">
                      <span>{c.check_type || c.id || "check"}</span>
                      <span className={`text-xs px-2 py-0.5 rounded-full ${c.status === "passed" ? "bg-green-100 text-green-700" : "bg-amber-100 text-amber-700"}`}>{c.status}</span>
                      {typeof c.score === "number" && <span className="text-xs text-muted-foreground">{c.score}</span>}
                    </div>
                  ))
                )}
                <div className="flex gap-2 pt-2">
                  <button onClick={() => review("request_info")} className="flex-1 rounded-xl border h-10 text-sm">Request Info</button>
                  <button onClick={() => review("approve")} className="flex-1 rounded-xl bg-green-600 text-white h-10 text-sm">Approve</button>
                  <button onClick={() => review("reject")} className="flex-1 rounded-xl bg-red-600 text-white h-10 text-sm">Reject</button>
                </div>
                {notice && <p className="text-xs text-muted-foreground">{notice}</p>}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
