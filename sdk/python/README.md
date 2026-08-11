# apexpay — Python SDK

A payments-only Python client for merchants who want just the payment gateway for
their e-commerce / online store (like Chapa, Arif Pay, or Telebirr) — no payroll,
HR, or ledger required.

Requires **Python 3.8+** and `requests`.

## Install

```bash
pip install requests
# then add sdk/python to your path, or pip install ./sdk/python
```

## Quickstart

```python
from apexpay import ApexPay

apexpay = ApexPay("sk_test_...", "https://api.apexpay.et")

# 1. Initialize a payment
payment = apexpay.initialize(
    tx_ref="order-1001",
    amount="2500.00",
    currency="ETB",
    method="telebirr",          # telebirr | cbe_birr | bank | card_acquirer | ethswitch
    customer_email="buyer@example.com",
    callback_url="https://store.example/api/apexpay-webhook",
)
# redirect the customer to payment["checkout_url"]

# 2. (On return) verify server-side
verified = apexpay.verify("order-1001")
if verified["status"] == "succeeded":
    # mark order paid
    pass

# 3. Or create a shareable payment link
link = apexpay.create_payment_link("1500.00", description="Order #1001")
```

## Webhooks

```python
from apexpay import verify_webhook_signature

raw_body = request.get_data(as_text=True)          # exact bytes ApexPay sent
signature = request.headers.get("X-ApexPay-Signature", "")
if not verify_webhook_signature("your-signing-secret", raw_body, signature):
    # reject — signature invalid
    pass
# trusted event — mark the order paid (idempotently)
```

## API

- `initialize(tx_ref, amount, ...)` → payment with `checkout_url`
- `verify(tx_ref)` → payment with `status`
- `get_payment(id)` → payment
- `create_payment_link(amount, ...)` → payment link
- `verify_webhook_signature(secret, raw_body, signature)` → bool

See [`docs/PAYMENTS_ONLY_GUIDE.md`](../../docs/PAYMENTS_ONLY_GUIDE.md) for the full
integration modes and setup.
