// Typed client for the ApexPay merchant API. All calls go through the same-origin
// server-side proxy (/api/proxy/...), which injects the merchant API key and forwards
// to the Go API. The browser never holds the key.

export type PaymentStatus =
  | "created"
  | "pending"
  | "processing"
  | "succeeded"
  | "failed"
  | "canceled"
  | "refunded"
  | "partially_refunded"

export type Payment = {
  id: string
  merchant_id: string
  tx_ref: string
  amount: string
  currency: string
  status: PaymentStatus
  method?: string
  description?: string
  customer_email?: string
  connector_id: string
  connector_ref?: string
  fee_amount?: string
  net_amount?: string
  requires_2fa?: boolean
  two_fa_verified?: boolean
  succeeded_at?: string | null
  created_at?: string
}

export type DashboardSummary = {
  tpv_today: string
  tpv_7_days: string
  success_count_7_days: number
  total_count_7_days: number
  success_rate_7_days: number
  active_links: number
  refunded_count_7_days: number
}

export type PaymentLink = {
  id: string
  amount: string
  currency: string
  description?: string
  status: string
  public_token: string
  created_at?: string
}

// Internal request helper — unwraps the proxy and normalizes errors.
async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`/api/proxy/${path}`, init)
  if (!res.ok) {
    let detail = ""
    try {
      const body = await res.json()
      detail = body.message || body.error || ""
    } catch {
      /* ignore */
    }
    throw new Error(detail || `Request failed (${res.status})`)
  }
  return res.json()
}

function get<T>(path: string): Promise<T> {
  return request<T>(path)
}

function post<T>(path: string, body?: unknown): Promise<T> {
  return request<T>(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  })
}

export const api = {
  // Dashboard
  summary: () => get<DashboardSummary>("dashboard"),
  payments: (limit = 25) => get<Payment[]>(`transactions?limit=${limit}`),
  links: () => get<PaymentLink[]>("payment_links"),
  banks: () => get<Array<{ code: string; name: string; name_am?: string }>>("banks"),

  // Payments
  initialize: (payload: unknown, idempotencyKey?: string) =>
    post<unknown>("transactions/initialize", payload).then((r) => r),

  // Hosted checkout links
  createLink: (payload: unknown) => post<PaymentLink>("payment_links", payload),
}
