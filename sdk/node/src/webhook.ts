// Standalone webhook signature verification helper.
// Can be used without instantiating the full client (e.g. in a Cloud Function).
import { createHmac, timingSafeEqual } from "node:crypto";

/**
 * Verify an ApexPay webhook signature.
 *
 * @param signingSecret your webhook endpoint's signing secret
 * @param rawBody the raw request body string (exact bytes ApexPay sent)
 * @param signature the value of the `X-ApexPay-Signature` header
 */
export function verifyWebhookSignature(
  signingSecret: string,
  rawBody: string,
  signature: string,
): boolean {
  const expected = createHmac("sha256", signingSecret).update(rawBody).digest("hex");
  const a = Buffer.from(signature, "hex");
  const b = Buffer.from(expected, "hex");
  if (a.length !== b.length) return false;
  return timingSafeEqual(a, b);
}
