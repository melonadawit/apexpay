// Typed errors for the ApexPay SDK.

export class ApexPayError extends Error {
  constructor(message: string, readonly statusCode?: number, readonly code?: string) {
    super(message);
    this.name = "ApexPayError";
  }
}
