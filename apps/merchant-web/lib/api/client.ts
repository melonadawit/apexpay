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

export type CurrentAccount = {
  id: string
  account_number: string
  account_name: string
  account_type: string
  currency: string
  bank_code: string
  partner_bank_name: string
  status: string
  balance: string
  available_balance: string
  overdraft_limit: string
  is_primary: boolean
  is_lite: boolean
  is_virtual: boolean
  cheque_book_issued: boolean
  debit_card_issued: boolean
}

export type CreditLine = {
  id: string
  credit_limit: string
  available_credit: string
  utilized_credit: string
  interest_rate: string
  status: string
  credit_score?: number
}

export type ForexRequest = {
  id: string
  from_currency: string
  to_currency: string
  from_amount: string
  to_amount: string
  rate_used: string
  forex_fee_percent: string
  forex_fee_amount: string
  purpose: string
  status: string
  nbe_approval_status?: string
  created_at: string
}

export type ForexRate = {
  from_currency: string
  to_currency: string
  rate: string
  buy_rate: string
  sell_rate: string
  source: string
  last_updated_at: string
}

export type VirtualAccount = {
  id: string
  virtual_account_number: string
  customer_id: string
  purpose: string
  status: string
  bank_code: string
  created_at: string
}

export type Notification = {
  id: string
  type: string
  title: string
  message: string
  is_read: boolean
  action_url?: string
  created_at: string
}

export type CorporateCard = {
  id: string
  card_number_masked: string
  card_type: string
  card_network: string
  cardholder_name: string
  cardholder_email: string
  status: string
  credit_limit: string
  available_credit: string
  forex_markup_percent: string
  cashback_percent: string
}

export type EscrowAccount = {
  id: string
  agreement_id?: string
  account_number?: string
  account_name?: string
  amount: string
  status: string
  order_id?: string
  platform_fee: string
  seller_amount: string
}

export type SupportTicket = {
  id: string
  subject: string
  priority: string
  status: string
  assigned_to?: string
  created_at: string
}

export type BankVerification = {
  id: string
  bank_code: string
  account_number_masked: string
  account_name: string
  verification_method: string
  status: string
  created_at: string
}

export type VendorInvoice = {
  id: string
  vendor_id?: string
  invoice_number: string
  invoice_date: string
  due_date?: string
  amount: string
  currency: string
  tax_amount: string
  withholding_tax_amount: string
  total_amount: string
  status: string
  ocr_confidence: number
  vendor_name?: string
  file_key?: string
  created_at: string
}

export type PettyCashBudget = {
  id: string
  budget_name: string
  amount: string
  assigned_to?: string
  status: string
  spent_amount: string
  remaining_amount: string
  created_at: string
}

export type PettyCashExpense = {
  id: string
  budget_id: string
  amount: string
  description: string
  receipt_file_key?: string
  status: string
  created_at: string
}

export type TaxPayment = {
  id: string
  tax_type: string
  amount: string
  currency: string
  period_month?: number
  period_year?: number
  due_date?: string
  status: string
  payment_reference?: string
  paid_at?: string
  created_at: string
}

export type PayoutLink = {
  id: string
  amount: string
  currency: string
  public_token: string
  recipient_name?: string
  recipient_phone?: string
  recipient_email?: string
  purpose?: string
  status: string
  expires_at: string
  created_at: string
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

  // Banking modules (session-authenticated via proxy)
  banking: {
    currentAccounts: () => get<CurrentAccount[]>("banking/current_accounts"),
    creditLines: () => get<CreditLine[]>("banking/credit_lines"),
    forexRates: () => get<ForexRate[]>("banking/forex/rates"),
    forexRequests: () => get<ForexRequest[]>("banking/forex/requests"),
    virtualAccounts: () => get<VirtualAccount[]>("banking/virtual_accounts"),
    notifications: () => get<Notification[]>("banking/notifications"),
    corporateCards: () => get<CorporateCard[]>("banking/corporate_cards"),
    escrow: () => get<EscrowAccount[]>("banking/escrow"),
    supportTickets: () => get<SupportTicket[]>("banking/support_tickets"),
    bankVerifications: () => get<BankVerification[]>("banking/bank_verifications"),
    vendorInvoices: () => get<VendorInvoice[]>("banking/vendor_invoices"),
    createVendorInvoice: (payload: unknown) => post<VendorInvoice>("banking/vendor_invoices", payload),
    pettyCashBudgets: () => get<PettyCashBudget[]>("banking/petty_cash_budgets"),
    createPettyCashBudget: (payload: unknown) => post<PettyCashBudget>("banking/petty_cash_budgets", payload),
    pettyCashExpenses: () => get<PettyCashExpense[]>("banking/petty_cash_expenses"),
    createPettyCashExpense: (payload: unknown) => post<PettyCashExpense>("banking/petty_cash_expenses", payload),
    taxPayments: () => get<TaxPayment[]>("banking/tax_payments"),
    createTaxPayment: (payload: unknown) => post<TaxPayment>("banking/tax_payments", payload),
    payoutLinks: () => get<PayoutLink[]>("banking/payout_links"),
    createPayoutLink: (payload: unknown) => post<PayoutLink>("banking/payout_links", payload),
    relationshipManagers: () => get<Array<{ id: string; rm_user_id: string; status: string; assigned_at: string }>>("banking/relationship_managers"),
    accountingIntegrations: () => get<Array<{ id: string; provider: string; status: string; last_sync_status: string; last_sync_error: string; created_at: string }>>("banking/accounting_integrations"),
  },
}
