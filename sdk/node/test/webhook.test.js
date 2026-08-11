const { test } = require("node:test");
const assert = require("node:assert");
const { createHmac } = require("node:crypto");
const { verifyWebhookSignature } = require("../dist/webhook");

test("accepts a valid HMAC signature", () => {
  const secret = "whsec_super_secret_123";
  const body = '{"event_type":"payment.succeeded","status":"succeeded"}';
  const sig = createHmac("sha256", secret).update(body).digest("hex");
  assert.equal(verifyWebhookSignature(secret, body, sig), true);
});

test("rejects a signature with the wrong secret", () => {
  const body = '{"event_type":"payment.succeeded"}';
  const sig = createHmac("sha256", "wrong-secret").update(body).digest("hex");
  assert.equal(verifyWebhookSignature("right-secret", body, sig), false);
});

test("rejects a tampered body", () => {
  const secret = "whsec_secret";
  const body = '{"event_type":"payment.succeeded"}';
  const sig = createHmac("sha256", secret).update(body).digest("hex");
  const tampered = '{"event_type":"payment.failed"}';
  assert.equal(verifyWebhookSignature(secret, tampered, sig), false);
});
