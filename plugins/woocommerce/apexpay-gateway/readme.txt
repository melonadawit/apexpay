=== ApexPay for WooCommerce ===
Contributors: apexpay
Tags: payments, ethiopia, telebirr, cbe, etb, gateway
Requires at least: 6.0
Tested up to: 6.7
Requires PHP: 8.0
Stable tag: 0.1.0
License: Proprietary

Accept Telebirr, CBE Birr, bank, card, and EthSwitch QR payments via ApexPay.

== Description ==

ApexPay for WooCommerce adds a payment gateway that lets your store accept
Telebirr, CBE Birr, bank/IPS, card, and EthSwitch QR payments. It is
payments-only — no payroll or HR required.

Features:
* Hosted checkout — customers pay on ApexPay's secure page.
* Webhook-based order confirmation (HMAC-signed).
* Optional refund support scaffold.

== Installation ==

1. Upload the `apexpay-gateway` folder to `/wp-content/plugins/`.
2. Activate the plugin.
3. Go to WooCommerce → Settings → Payments → ApexPay.
4. Enter your secret key (`sk_test_...` / `sk_live_...`) and webhook secret.
5. Save. Set the webhook URL in the ApexPay developer portal to:
   `https://yourstore.com/?wc-api=apexpay_webhook`

== Changelog ==

= 0.1.0 =
* Initial release (skeleton).
