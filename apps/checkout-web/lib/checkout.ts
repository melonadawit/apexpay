// Typed client for the public hosted-checkout API. All calls go through the same-origin
// Next.js proxy (`/api/proxy/...`), so no CORS and no backend host leaks to the browser.

export type PaymentLink = {
  token: string
  amount: string
  currency: string
  description: string
  status: string
}

export type InitResult = {
  id: string
  tx_ref: string
  amount: string
  currency: string
  status: string
  connector_id: string
  requires_2fa: boolean
  fee_amount: string
  net_amount: string
}

export type StatusResult = {
  id: string
  tx_ref: string
  status: string
  connector_id: string
  requires_2fa: boolean
  two_fa_verified: boolean
  succeeded_at: string | null
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, init)
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

export async function fetchLink(token: string): Promise<PaymentLink> {
  return request(`/api/proxy/checkout/${encodeURIComponent(token)}`)
}

export async function initializePayment(
  token: string,
  method: string,
  customerEmail?: string
): Promise<InitResult> {
  return request(`/api/proxy/checkout/${encodeURIComponent(token)}/initialize`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ method, customer_email: customerEmail || "" }),
  })
}

export async function fetchStatus(token: string, txRef: string): Promise<StatusResult> {
  return request(
    `/api/proxy/checkout/${encodeURIComponent(token)}/status/${encodeURIComponent(txRef)}`
  )
}

export async function verify2FA(token: string, paymentID: string, otp: string): Promise<void> {
  return request(`/api/proxy/checkout/${encodeURIComponent(token)}/2fa/${encodeURIComponent(paymentID)}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ otp }),
  })
}
