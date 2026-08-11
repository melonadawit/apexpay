# ApexPay SDKs & Integrations

Official client libraries and plugins for merchants who want **just the payment
gateway** for their e-commerce / online store — like Chapa, Arif Pay, or Telebirr.
No payroll, HR, or ledger required (those modules are optional on your account).

## Available SDKs

| Language | Package | Path | Status |
|---|---|---|---|
| Node.js / TypeScript | `@apexpay/node` | [`sdk/node`](./node) | ✅ build + tests + live API e2e |
| PHP | `apexpay/php` | [`sdk/php`](./php) | ✅ php -l + live API e2e |
| Python | `apexpay` | [`sdk/python`](./python) | ✅ pytest + live API e2e |
| Go | (example) | [`examples/go`](../examples/go) | ✅ go vet/build |

## Plugins

| Platform | Path | Status |
|---|---|---|
| WooCommerce | [`plugins/woocommerce/apexpay-gateway`](../plugins/woocommerce/apexpay-gateway) | ✅ php -l clean (requires a WP install to run) |

## Integration modes

1. **Hosted checkout / payment link** — no code, share on WhatsApp/Telegram.
2. **Embedded JS SDK (`checkout.js`)** — public-key snippet for SPAs/mobile-web.
3. **Direct REST API** — `initialize` → redirect → `verify` for full cart/backend control.
4. **Webhook** — HMAC-signed `payment.succeeded` events for authoritative fulfilment.

See [`docs/PAYMENTS_ONLY_GUIDE.md`](../docs/PAYMENTS_ONLY_GUIDE.md) for the full guide
with exact request/response shapes.

> **Note:** The SDKs are working scaffolds verified against the live API. Publishing
> to npm / Packagist / PyPI / the WordPress plugin repo requires real registry
> credentials and is the remaining step to make them install-and-go drop-ins.
