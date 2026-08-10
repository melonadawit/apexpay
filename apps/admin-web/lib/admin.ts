// Typed client for the admin API. All calls go through the same-origin proxy
// (/api/proxy/...), so no CORS and no backend host leak to the browser.

export type QueueItem = {
  merchant_id: string
  legal_name: string
  email: string
  onboarding_status: string
  risk_score: number
  risk_tier: string
  fayda_verified: boolean
  submitted_at?: string
  created_at: string
}

export type ComplianceCheck = {
  id?: string
  check_type?: string
  status?: string
  score?: number
}

export type BankAccount = {
  id?: string
  account_number_masked?: string
  bank_code?: string
  account_name?: string
  verification_status?: string
}

export type MerchantExam = {
  merchant_id: string
  legal_name: string
  onboarding_status: string
  risk_score: number
  risk_tier: string
  compliance_checks?: ComplianceCheck[]
  banks?: BankAccount[]
  owners?: any[]
  documents?: any[]
  onboarding_reviews?: any[]
}

export type ConnectorHealth = {
  connector_id: string
  name?: string
  status?: string
  latency_ms?: number
}

export type ReconBreak = {
  tx_ref?: string
  merchant_id?: string
  amount?: string
  status?: string
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
  const json = await res.json()
  return json.data as T
}

export async function getQueue(): Promise<QueueItem[]> {
  return request<QueueItem[]>("/api/proxy/admin/onboarding/queue")
}

export async function getExam(merchantId: string): Promise<MerchantExam> {
  return request<MerchantExam>(`/api/proxy/admin/merchants/${encodeURIComponent(merchantId)}/exam`)
}

export async function reviewMerchant(merchantId: string, action: string, comment: string) {
  return request(`/api/proxy/admin/onboarding/${encodeURIComponent(merchantId)}/review`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ action, comment }),
  })
}

export async function getConnectorHealth(): Promise<ConnectorHealth[]> {
  return request<ConnectorHealth[]>("/api/proxy/admin/connectors/health")
}

export async function getReconBreaks(): Promise<ReconBreak[]> {
  return request<ReconBreak[]>("/api/proxy/admin/recon/breaks")
}
