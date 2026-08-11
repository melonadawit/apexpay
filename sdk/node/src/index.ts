// ApexPay Node.js SDK — payments-only.
//
// A minimal client for merchants who want just the payment gateway for their
// e-commerce / online store (like Chapa, Arif Pay, or Telebirr), without payroll,
// HR, ledger, or any other module. See docs/PAYMENTS_ONLY_GUIDE.md.

import { ApexPayError } from "./errors";
import { verifyWebhookSignature } from "./webhook";

export type Rail =
  | "telebirr"
  | "cbe_birr"
  | "bank"
  | "card_acquirer"
  | "ethswitch";

export type PaymentStatus =
  | "created"
  | "pending"
  | "processing"
  | "succeeded"
  | "failed"
  | "canceled";

export interface ApexPayConfig {
  /** Test key `sk_test_...` or live key `sk_live_...`. */
  apiKey: string;
  /** Base URL of the API. Defaults to the local dev API. */
  baseUrl?: string;
}

export interface InitializeRequest {
  /** Your unique reference — use your order id. */
  tx_ref: string;
  amount: string;
  currency?: string;
  method?: Rail;
  description?: string;
  customer_email?: string;
  return_url?: string;
  callback_url?: string;
  /** Set for safe retries. */
  idempotencyKey?: string;
}

export interface Payment {
  id: string;
  tx_ref: string;
  amount: string;
  currency: string;
  status: PaymentStatus;
  checkout_url: string;
  connector_id: string;
  requires_2fa: boolean;
  fee_amount: string;
  net_amount: string;
  [key: string]: unknown;
}

export interface PaymentLink {
  id: string;
  amount: string;
  currency: string;
  status: string;
  public_token: string;
  checkout_url: string;
  expires_at?: string | null;
  share?: { whatsapp?: string; telegram?: string };
  [key: string]: unknown;
}

export interface WebhookEvent {
  event_type: string;
  payment_id?: string;
  tx_ref?: string;
  status?: string;
  [key: string]: unknown;
}

export class ApexPay {
  private readonly apiKey: string;
  private readonly baseUrl: string;

  constructor(config: ApexPayConfig) {
    if (!config.apiKey) {
      throw new Error("ApexPay: apiKey is required (sk_test_... or sk_live_...)");
    }
    this.apiKey = config.apiKey;
    this.baseUrl = (config.baseUrl || "http://localhost:8080").replace(/\/$/, "");
  }

  private async request<T>(method: string, path: string, body?: unknown, idempotencyKey?: string): Promise<T> {
    const headers: Record<string, string> = {
      Authorization: `Bearer ${this.apiKey}`,
      "Content-Type": "application/json",
    };
    if (idempotencyKey) headers["Idempotency-Key"] = idempotencyKey;

    const res = await fetch(`${this.baseUrl}/v1/${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
    const json = await res.json().catch(() => ({}));
    if (!res.ok) {
      const err = json as { error?: { message?: string; code?: string } };
      throw new ApexPayError(
        err.error?.message || (json as { message?: string }).message || `ApexPay request failed (${res.status})`,
        res.status,
        err.error?.code,
      );
    }
    // Unwrap the { success, data } envelope.
    return (json as { data: T }).data;
  }

  /** Initialize a payment and get a checkout_url to redirect the customer to. */
  async initialize(req: InitializeRequest): Promise<Payment> {
    return this.request<Payment>(
      "POST",
      "transactions/initialize",
      {
        tx_ref: req.tx_ref,
        amount: req.amount,
        currency: req.currency || "ETB",
        method: req.method,
        description: req.description,
        customer_email: req.customer_email,
        return_url: req.return_url,
        callback_url: req.callback_url,
      },
      req.idempotencyKey,
    );
  }

  /** Server-side verification of a payment by your tx_ref. */
  async verify(txRef: string): Promise<Payment> {
    return this.request<Payment>("GET", `transactions/verify/${encodeURIComponent(txRef)}`);
  }

  /** Get a single payment by id. */
  async getPayment(id: string): Promise<Payment> {
    return this.request<Payment>("GET", `transactions/${encodeURIComponent(id)}`);
  }

  /** Create a shareable payment link (hosted checkout). */
  async createPaymentLink(input: { amount: string; currency?: string; description?: string }): Promise<PaymentLink> {
    return this.request<PaymentLink>("POST", "payment_links", {
      amount: input.amount,
      currency: input.currency || "ETB",
      description: input.description,
    });
  }

  /**
   * Verify an inbound webhook HMAC signature.
   * ApexPay signs each delivery with X-ApexPay-Signature = HMAC-SHA256(signingSecret, rawBody).
   */
  verifyWebhookSignature(signingSecret: string, rawBody: string, signature: string): boolean {
    return verifyWebhookSignature(signingSecret, rawBody, signature);
  }
}
