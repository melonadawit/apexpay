<?php
/**
 * Example: ApexPay payment gateway for a WooCommerce-style checkout.
 *
 * This is a reference flow you can adapt into a WooCommerce / custom cart plugin:
 *   1. Initialize a payment for the order
 *   2. Redirect the customer to checkout_url
 *   3. On return, verify server-side
 *   4. Receive a webhook and mark the order paid (authoritative)
 */

require __DIR__ . '/../src/ApexPay.php';

use ApexPay\ApexPay;
use ApexPay\ApexPayError;

$apiKey   = 'sk_test_...';          // or sk_live_...
$webhookSecret = 'whsec_...';       // from the developer portal

$apexpay = new ApexPay($apiKey, 'https://api.apexpay.et');

try {
    // --- 1. Initialize the payment for this order ---
    $payment = $apexpay->initialize([
        'tx_ref'          => 'order-' . $orderId,   // your order id
        'amount'          => '2500.00',
        'currency'        => 'ETB',
        'method'          => 'telebirr',            // telebirr | cbe_birr | bank | card_acquirer | ethswitch
        'customer_email'  => 'buyer@example.com',
        'return_url'      => 'https://store.example/checkout/return',
        'callback_url'    => 'https://store.example/api/apexpay-webhook',
        'idempotencyKey'  => 'order-' . $orderId,   // safe retries
    ]);

    // --- 2. Redirect the customer ---
    header('Location: ' . $payment['checkout_url']);
    exit;

} catch (ApexPayError $e) {
    error_log("ApexPay initialize failed: {$e->getMessage()} ({$e->getCode()})");
    // render an error page
}

// --- 3 & 4. Return + webhook handler ---
// On return: $payment = $apexpay->verify($orderId);
// In your webhook route, verify the signature then mark the order paid idempotently:
/*
$rawBody = file_get_contents('php://input');
$signature = $_SERVER['HTTP_X_APEXPAY_SIGNATURE'] ?? '';
if (!ApexPay::verifyWebhookSignature($webhookSecret, $rawBody, $signature)) {
    http_response_code(401);
    exit('invalid signature');
}
$event = json_decode($rawBody, true);
if (($event['event_type'] ?? '') === 'payment.succeeded') {
    // mark order paid (idempotently)
}
http_response_code(200);
*/
