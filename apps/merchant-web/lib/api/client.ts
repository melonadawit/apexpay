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
  checkout_url?: string
  fee_amount?: string
  net_amount?: string
  requires_2fa?: boolean
  two_fa_verified?: boolean
  succeeded_at?: string | null
  created_at?: string
}

export type PaymentEntry = {
  direction: string
  amount: string
  currency: string
  account_code: string
  account_name: string
}

export type PaymentJournal = {
  id: string
  book_id: string
  posting_key: string
  memo: string
  created_at: string
  entries: PaymentEntry[]
}

export type PaymentDetail = {
  payment: Payment
  journals: PaymentJournal[]
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

function put<T>(path: string, body?: unknown): Promise<T> {
  return request<T>(path, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: body ? JSON.stringify(body) : undefined,
  })
}

function del<T>(path: string): Promise<T> {
  return request<T>(path, { method: "DELETE" })
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

// ---- HRIS ----
export type Team = { id: string; name: string; department_id?: string; manager_id?: string; description?: string; created_at: string }
export type Contract = { id: string; employee_id: string; contract_type: string; start_date: string; end_date?: string; salary_amount: string; salary_currency: string; probation_months: number; notice_days: number; status: string; doc_key?: string; signed_at?: string }
export type Shift = { id: string; name: string; start_time: string; end_time: string; break_min: number }
export type AttendanceClock = { id: string; employee_id: string; shift_id?: string; clock_date: string; punch_in?: string; punch_out?: string; hours: string; status: string; source: string; note?: string }
export type OnboardingTask = { id: string; employee_id: string; task: string; category: string; due_in_days: number; status: string; assigned_to?: string }
export type PerformanceReview = { id: string; employee_id: string; reviewer_id?: string; period: string; rating?: number; goals?: string; comments?: string; status: string }

// ---- Risk ----
export type RiskRule = { id: string; name: string; rule_type: string; parameters?: Record<string, unknown>; action: string; severity: string; enabled: boolean }
export type RiskFlag = { id: string; entity_type: string; entity_id: string; rule_name: string; severity: string; action: string; reason: string; status: string; details?: Record<string, unknown> }
export type RiskEval = { findings: Array<{ rule_name: string; severity: string; action: string; reason: string }>; block: boolean; approved: boolean }

// ---- Treasury ----
export type TreasuryPosition = { accounts: Array<{ account_id: string; account_number: string; account_name: string; account_type: string; bank_code: string; balance: string; available_balance: string; currency: string }>; total_balance: string; total_available: string; currency: string; generated_at: string }
export type Transfer = { id: string; from_account_id: string; to_account_id: string; amount: string; currency: string; purpose?: string; status: string; created_at: string }
export type Forecast = { id: string; forecast_date: string; horizon_days: number; inflow_today: string; inflow_30d: string; inflow_60d: string; inflow_90d: string; outflow_today: string; outflow_30d: string; outflow_60d: string; outflow_90d: string; net_90d: string; generated_at: string }

// ---- Invoicing ----
export type Invoice = { id: string; invoice_number: string; customer_name: string; customer_email?: string; issue_date: string; due_date: string; currency: string; subtotal: string; tax_amount: string; withholding_amount: string; total_amount: string; amount_paid: string; status: string; hosted_token?: string; dunning_stage: number; line_items: Array<{ description: string; quantity: string; unit_price: string; line_total: string }>; created_at: string }
export type AgingBucket = { bucket: string; count: number; amount: string }

// ---- Team & approvals ----
export type Member = { user_id: string; email: string; name: string; role: string; permissions: string[]; created_at: string }
export type Approval = { id: string; resource_type: string; resource_id: string; action: string; summary: string; amount: string; currency: string; status: string; required_approvals: number; approvals: Array<{ user_id: string; name: string; role: string; decision: string; decided_at: string }>; created_at: string }

// ---- Compliance console ----
export type ComplianceStatus = { merchant_id: string; onboarding_status: string; kyc_expiry_date?: string; license_expiry?: string; fayda_verified: boolean; risk_tier: string; next_erca_due?: string; next_pension_due?: string; annual_tax_filing_due?: string; aml_due?: string; overall_status: string; notes?: string }
export type ComplianceCheck = { id: string; check_type: string; status: string; detail: string; checked_at: string }

// ---- Notification preferences ----
export type NotifyPref = { event_type: string; email: boolean; sms: boolean; push: boolean; in_app: boolean }

// ---- Fixed assets ----
export type FixedAsset = { id: string; asset_name: string; category: string; acquisition_date: string; cost: string; salvage_value: string; useful_life_years: number; depreciation_method: string; depreciation_rate?: string; accumulated_depreciation: string; net_book_value: string; status: string }
export type DepreciationEntry = { id: string; asset_id: string; period: string; amount: string; book_value_after: string }

// ---- Analytics ----
export type RevenueDaily = { stat_date: string; revenue: string; tpv: string; payment_count: number; success_count: number; failed_count: number; refund_amount: string }
export type Cohort = { cohort_month: string; customers: number; month1_retention: string; month2_retention: string; month3_retention: string; mrr: string }

// ---- 2FA ----
export type TwoFAEnroll = { secret: string; otpauth_url: string; issuer: string }

// ---- Inventory & Sales ----
export type Product = { id: string; name: string; description?: string; sku?: string; price: string; cost_price: string; currency: string; vat_category: string; stock_qty: string; low_stock_threshold: string; status: string }
export type Order = { id: string; order_number: string; customer_name?: string; customer_email?: string; status: string; subtotal: string; tax_amount: string; total_amount: string; items: Array<{ product_id?: string; description: string; quantity: string; unit_price: string; line_total: string }>; created_at: string }
export type StockMovement = { id: string; product_id: string; qty: string; direction: string; reference?: string; note?: string; created_at: string }

// ---- Disputes ----
export type Dispute = { id: string; payment_id?: string; amount: string; currency: string; reason_code: string; status: string; evidence: Array<{ file_key: string; description?: string }>; resolution?: string; fee: string; created_at: string }

// ---- Loyalty ----
export type LoyaltyTier = { id: string; name: string; min_spend: string; cashback_percent: string }
export type LoyaltyAccount = { id: string; customer_email?: string; customer_phone?: string; points: string; tier_id?: string; tier_name?: string; total_spend: string }
export type CashbackTx = { id: string; payment_id?: string; amount: string; type: string; created_at: string }

// ---- Lending ----
export type Loan = { id: string; amount: string; currency: string; purpose: string; status: string; interest_rate: string; due_date?: string; repaid_amount: string; outstanding_amount: string; created_at: string }

// ---- Accounting & Bookkeeping ----
export type Account = { code: string; name: string; category: string; normal_side: string; is_active: boolean }
export type TrialBalanceRow = { code: string; name: string; debit: string; credit: string }
export type StatementLine = { label: string; amount: string; kind: string }
export type FinancialStatement = { title: string; period: string; lines: StatementLine[] }
export type CashFlowLine = { label: string; amount: string; kind: string }

export const api = {
  // Dashboard
  summary: () => get<DashboardSummary>("dashboard"),
  payments: (limit = 25) => get<Payment[]>(`transactions?limit=${limit}`),
  payment: (id: string) => get<PaymentDetail>(`transactions/${id}`),
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

  // Workforce OS (HRIS)
  hris: {
    teams: () => get<Team[]>("hris/teams"),
    createTeam: (p: unknown) => post<Team>("hris/teams", p),
    contracts: () => get<Contract[]>("hris/contracts"),
    createContract: (p: unknown) => post<Contract>("hris/contracts", p),
    shifts: () => get<Shift[]>("hris/shifts"),
    attendance: (from = "", to = "") => get<AttendanceClock[]>(`hris/attendance${from ? `?from=${from}&to=${to}` : ""}`),
    clockIn: (p: unknown) => post<AttendanceClock>("hris/attendance/clock-in", p),
    clockOut: (p: unknown) => post<AttendanceClock>("hris/attendance/clock-out", p),
    onboardingTasks: () => get<OnboardingTask[]>("hris/onboarding-tasks"),
    reviews: () => get<PerformanceReview[]>("hris/performance-reviews"),
  },

  // Risk & Fraud
  risk: {
    rules: () => get<RiskRule[]>("risk/rules"),
    createRule: (p: unknown) => post<RiskRule>("risk/rules", p),
    flags: (status = "") => get<RiskFlag[]>(`risk/flags${status ? `?status=${status}` : ""}`),
    evaluate: (p: unknown) => post<RiskEval>("risk/evaluate", p),
  },

  // Treasury
  treasury: {
    position: () => get<TreasuryPosition>("treasury/position"),
    transfers: () => get<Transfer[]>("treasury/transfers"),
    createTransfer: (p: unknown) => post<Transfer>("treasury/transfers", p),
    forecast: () => get<Forecast>("treasury/forecast"),
    generateForecast: () => post<Forecast>("treasury/forecast"),
  },

  // Invoicing
  invoices: {
    list: (status = "") => get<Invoice[]>(`invoices${status ? `?status=${status}` : ""}`),
    create: (p: unknown) => post<Invoice>("invoices", p),
    get: (id: string) => get<Invoice>(`invoices/${id}`),
    aging: () => get<AgingBucket[]>("invoices/aging"),
    issue: (id: string) => post<unknown>(`invoices/${id}/issue`),
  },

  // Team & approvals
  team: {
    members: () => get<Member[]>("team/members"),
    invite: (p: unknown) => post<Member>("team/members", p),
    updateRole: (userID: string, role: string) => put(`team/members/${userID}/role`, { role }),
    remove: (userID: string) => del(`team/members/${userID}`),
    approvals: (status = "") => get<Approval[]>(`team/approvals${status ? `?status=${status}` : ""}`),
    decide: (id: string, decision: string) => post<unknown>(`team/approvals/${id}/decide`, { decision }),
  },

  // Compliance console
  compliance: {
    status: () => get<ComplianceStatus>("compliance-console/status"),
    checks: () => get<ComplianceCheck[]>("compliance-console/checks"),
    addCheck: (p: unknown) => post<unknown>("compliance-console/checks", p),
    ask: (query: string) => post<{ answer?: string; citations?: unknown[] }>("compliance/ask", { query }),
  },

  // Notification preferences
  notificationsPrefs: {
    list: () => get<NotifyPref[]>("notifications/preferences"),
    update: (p: NotifyPref) => put<NotifyPref>("notifications/preferences", p),
  },

  // Fixed assets
  fixedAssets: {
    list: () => get<FixedAsset[]>("fixed-assets"),
    create: (p: unknown) => post<FixedAsset>("fixed-assets", p),
    depreciate: (id: string) => post<DepreciationEntry>(`fixed-assets/${id}/depreciate`),
    entries: () => get<DepreciationEntry[]>("fixed-assets/depreciation"),
  },

  // Analytics
  analytics: {
    revenue: () => get<RevenueDaily[]>("analytics/revenue"),
    methods: () => get<Array<{ method: string; count: number; success: number; revenue: string }>>("analytics/methods"),
    cohorts: () => get<Cohort[]>("analytics/cohorts"),
  },

  // Real 2FA (TOTP)
  twofa: {
    enroll: (account: string) => post<TwoFAEnroll>("2fa/enroll", { account }),
    verify: (secret: string, code: string) => post<{ verified: boolean }>("2fa/verify", { secret, code }),
  },

  // Accounting & Bookkeeping
  accounting: {
    accounts: () => get<Account[]>("accounting/accounts"),
    trialBalance: () => get<TrialBalanceRow[]>("accounting/trial-balance"),
    profitLoss: (from = "", to = "") => get<FinancialStatement>(`accounting/profit-loss${from ? `?from=${from}&to=${to}` : ""}`),
    balanceSheet: () => get<FinancialStatement>("accounting/balance-sheet"),
    cashFlow: (from = "", to = "") => get<CashFlowLine[]>(`accounting/cash-flow${from ? `?from=${from}&to=${to}` : ""}`),
  },

  // Inventory & Sales
  inventory: {
    products: () => get<Product[]>("inventory/products"),
    createProduct: (p: unknown) => post<Product>("inventory/products", p),
    addStock: (id: string, qty: string, note = "") => post<StockMovement>(`inventory/products/${id}/stock`, { qty, note }),
    orders: (status = "") => get<Order[]>(`inventory/orders${status ? `?status=${status}` : ""}`),
    createOrder: (p: unknown) => post<Order>("inventory/orders", p),
    markPaid: (id: string, paymentId: string) => post<unknown>(`inventory/orders/${id}/mark-paid`, { payment_id: paymentId }),
    stockMovements: () => get<StockMovement[]>("inventory/stock-movements"),
  },

  // Disputes
  disputes: {
    list: (status = "") => get<Dispute[]>(`disputes${status ? `?status=${status}` : ""}`),
    create: (p: unknown) => post<Dispute>("disputes", p),
    evidence: (id: string, evidence: unknown[]) => post<unknown>(`disputes/${id}/evidence`, { evidence }),
    decide: (id: string, decision: string, resolution = "", fee = "0") => post<unknown>(`disputes/${id}/decide`, { decision, resolution, fee }),
  },

  // Loyalty
  loyalty: {
    tiers: () => get<LoyaltyTier[]>("loyalty/tiers"),
    createTier: (p: unknown) => post<LoyaltyTier>("loyalty/tiers", p),
    accounts: () => get<LoyaltyAccount[]>("loyalty/accounts"),
    earn: (id: string, amount: string, paymentId = "") => post<CashbackTx>(`loyalty/accounts/${id}/earn`, { amount, payment_id: paymentId }),
    transactions: () => get<CashbackTx[]>("loyalty/transactions"),
  },

  // Lending
  lending: {
    loans: (status = "") => get<Loan[]>(`lending${status ? `?status=${status}` : ""}`),
    create: (p: unknown) => post<Loan>("lending", p),
    repay: (id: string, amount: string) => post<unknown>(`lending/${id}/repay`, { amount }),
  },

  // Payroll & Workforce
  payroll: {
    employees: () => get<unknown[]>("payroll/employees"),
    createEmployee: (p: unknown) => post<unknown>("payroll/employees", p),
    employeeRevisions: (employeeId: string) => get<unknown[]>(`payroll/employees/${employeeId}/revisions`),
    runs: () => get<unknown[]>("payroll/payroll_runs"),
    calendars: (year = 2026) => get<unknown[]>(`payroll/calendars?year=${year}`),
    lockCalendar: (id: string) => post<unknown>(`payroll/calendars/${id}/lock`),
    unlockCalendar: (id: string) => post<unknown>(`payroll/calendars/${id}/unlock`),
    finalSettlements: () => get<unknown[]>("payroll/final_settlements"),
    auditLogs: () => get<unknown[]>("payroll/payroll_audit_logs"),
    costCenterReport: (year = 2026, month = 0) => get<unknown>(`payroll/payroll_reports/cost_center?year=${year}${month ? `&month=${month}` : ""}`),
    varianceReport: (year = 2026, month = 0) => get<unknown>(`payroll/payroll_reports/variance?year=${year}${month ? `&month=${month}` : ""}`),
    payrollRegister: (runId = "") => get<unknown>(`payroll/payroll_reports/payroll_register${runId ? `?run_id=${runId}` : ""}`),
    createRun: (p: unknown) => post<unknown>("payroll/payroll_runs", p),
    calculate: (id: string) => post<unknown>(`payroll/payroll_runs/${id}/calculate`),
    approve: (id: string) => post<unknown>(`payroll/payroll_runs/${id}/approve`),
    disburse: (id: string) => post<unknown>(`payroll/payroll_runs/${id}/disburse`),
    runItems: (id: string) => get<unknown[]>(`payroll/payroll_runs/${id}/items`),
    leaveRequests: () => get<unknown[]>("payroll/leave_requests"),
    leaveBalances: () => get<unknown[]>("payroll/leave_balances"),
    claims: () => get<unknown[]>("payroll/claims"),
    loans: () => get<unknown[]>("payroll/loans"),
    departments: () => get<unknown[]>("payroll/departments"),
    salaryStructures: () => get<unknown[]>("payroll/salary_structures"),
  },

  // Refunds
  refunds: {
    get: (id: string) => get<unknown>(`refunds/${id}`),
    byPayment: (paymentId: string) => get<unknown[]>(`refunds/payment/${paymentId}`),
    create: (p: unknown) => post<unknown>("refunds", p),
  },

  // Payouts
  payouts: {
    batches: () => get<unknown[]>("payout_batches"),
    getBatch: (id: string) => get<unknown>(`payout_batches/${id}`),
    approveBatch: (id: string) => post<unknown>(`payout_batches/${id}/approve`),
    beneficiaries: () => get<unknown[]>("beneficiaries"),
    createBeneficiary: (p: unknown) => post<unknown>("beneficiaries", p),
  },

  // Subscriptions
  subscriptions: {
    plans: () => get<unknown[]>("subscription_plans"),
    subscriptions: () => get<unknown[]>("subscriptions"),
    customers: () => get<unknown[]>("customers"),
  },

  // Procurement & AP
  procurement: {
    vendors: () => get<unknown[]>("procurement/vendors"),
    createVendor: (p: unknown) => post<unknown>("procurement/vendors", p),
    purchaseOrders: () => get<unknown[]>("procurement/purchase-orders"),
    createPurchaseOrder: (p: unknown) => post<unknown>("procurement/purchase-orders", p),
    receive: (id: string) => post<unknown>(`procurement/purchase-orders/${id}/receive`),
    invoices: () => get<unknown[]>("procurement/invoices"),
    createInvoice: (p: unknown) => post<unknown>("procurement/invoices", p),
    aging: () => get<unknown[]>("procurement/aging"),
  },

  // Budgeting / FP&A
  budget: {
    budgets: (period = "") => get<unknown[]>(`budget/budgets${period ? `?period=${period}` : ""}`),
    setBudget: (p: unknown) => post<unknown>("budget/budgets", p),
    variance: (period = "") => get<unknown>(`budget/variance${period ? `?period=${period}` : ""}`),
  },

  // Tax schedules
  tax: {
    schedule: () => get<unknown>("tax/schedule"),
    postToLedger: (period = "") => post<unknown>("tax/schedule/post", { period }),
  },

  // FX revaluation
  fx: {
    revalue: (period = "") => post<unknown>("fx/revalue", { period }),
  },

  // Accounting GL (journals + periods)
  gl: {
    journalEntries: () => get<unknown[]>("accounting/journal-entries"),
    postJournalEntry: (p: unknown) => post<unknown>("accounting/journal-entries", p),
    periods: () => get<unknown[]>("accounting/periods"),
    closePeriod: (period: string) => post<unknown>("accounting/periods/close", { period }),
    reopenPeriod: (period: string) => post<unknown>("accounting/periods/reopen", { period }),
  },

  // Assistant
  assistant: {
    chat: (message: string) => post<unknown>("assistant/chat", { message }),
    thread: (id: string) => get<unknown>(`assistant/threads/${id}`),
  },

  // Portals
  portal: {
    createToken: (p: unknown) => post<unknown>("portal/token", p),
    me: () => get<unknown>("portal/me"),
  },
}
