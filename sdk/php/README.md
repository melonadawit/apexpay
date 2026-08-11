# ApexPay PHP SDK

A payments-only PHP client for merchants who want just the payment gateway for their
e-commerce / online store (like Chapa, Arif Pay, or Telebirr) — no payroll, HR, or
ledger required.

Requires **PHP 8.0+** with the cURL extension.

## Quickstart

```php
require 'src/ApexPay.php';

use ApexPay\ApexPay;

$apexpay = new ApexPay('sk_test_...', 'https://api.apexpay.et');

// 1. Initialize a payment
$payment = $apexpay->initialize([
    'tx_ref'          => 'order-1001',
    'amount'          => '2500.00',
    'currency'        => 'ETB',
    'method'          => 'telebirr',   // telebirr | cbe_birr | bank | card_acquirer | ethswitch
    'customer_email'  => 'buyer@example.com',
    'return_url'      => 'https://store.example/checkout/return',
    'callback_url'    => 'https://store.example/api/apexpay-webhook',
]);
// redirect the customer to $payment['checkout_url']

// 2. (On return) verify server-side
$verified = $apexpay->verify('order-1001');
if ($verified['status'] === 'succeeded') {
    // mark order paid
}

// 3. Or create a shareable payment link
$link = $apexpay->createPaymentLink(['amount' => '1500.00', 'description' => 'Order #1001']);
```

## Webhooks

```php
use ApexPay\ApexPay;

$rawBody   = file_get_contents('php://input');
$signature = $_SERVER['HTTP_X_APEXPAY_SIGNATURE'] ?? '';

if (!ApexPay::verifyWebhookSignature('your-signing-secret', $rawBody, $signature)) {
    http_response_code(401);
    exit('invalid signature');
}
// trusted event — mark the order paid (idempotently)
```

## API

- `initialize(array $req): array` → payment with `checkout_url`
- `verify(string $txRef): array` → payment with `status`
- `getPayment(string $id): array`
- `createPaymentLink(array $input): array`
- `ApexPay::verifyWebhookSignature($secret, $rawBody, $signature): bool`

## WooCommerce

See `examples/woocommerce_checkout.php` for a reference checkout + webhook flow you
can adapt into a WooCommerce/custom-cart plugin.

See [`docs/PAYMENTS_ONLY_GUIDE.md`](../../docs/PAYMENTS_ONLY_GUIDE.md) for the full
integration modes and setup.
