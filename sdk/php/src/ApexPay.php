<?php
/**
 * ApexPay PHP SDK — payments-only client.
 *
 * A minimal client for merchants who want just the payment gateway for their
 * e-commerce / online store (like Chapa, Arif Pay, or Telebirr), without payroll,
 * HR, ledger, or any other module. See docs/PAYMENTS_ONLY_GUIDE.md.
 *
 * Requires PHP 8.0+ with the cURL extension.
 */

namespace ApexPay;

class ApexPayError extends \Exception
{
    public ?int $statusCode;
    public ?string $errorCode;

    public function __construct(string $message, ?int $statusCode = null, ?string $errorCode = null)
    {
        parent::__construct($message);
        $this->statusCode = $statusCode;
        $this->errorCode = $errorCode;
    }
}

class ApexPay
{
    private string $apiKey;
    private string $baseUrl;

    public function __construct(string $apiKey, string $baseUrl = 'http://localhost:8080')
    {
        if ($apiKey === '') {
            throw new ApexPayError('ApexPay: apiKey is required (sk_test_... or sk_live_...)');
        }
        $this->apiKey = $apiKey;
        $this->baseUrl = rtrim($baseUrl, '/');
    }

    /**
     * @param string $method  GET | POST
     * @param string $path    e.g. "transactions/initialize"
     * @param array|null $body
     * @param string|null $idempotencyKey
     * @return array  the unwrapped `data` payload
     */
    private function request(string $method, string $path, ?array $body = null, ?string $idempotencyKey = null): array
    {
        $url = $this->baseUrl . '/v1/' . $path;
        $headers = [
            'Authorization: Bearer ' . $this->apiKey,
            'Content-Type: application/json',
        ];
        if ($idempotencyKey !== null) {
            $headers[] = 'Idempotency-Key: ' . $idempotencyKey;
        }

        $ch = curl_init($url);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_HTTPHEADER => $headers,
            CURLOPT_TIMEOUT => 30,
        ]);
        if ($body !== null) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($body));
        }

        $raw = curl_exec($ch);
        if ($raw === false) {
            $err = curl_error($ch);
            curl_close($ch);
            throw new ApexPayError('ApexPay request failed: ' . $err);
        }
        $status = (int) curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        $json = json_decode($raw, true);
        if (!is_array($json)) {
            throw new ApexPayError('ApexPay returned an invalid response', $status);
        }
        if ($status >= 400) {
            $msg = $json['error']['message'] ?? $json['message'] ?? ('ApexPay request failed (' . $status . ')');
            throw new ApexPayError((string) $msg, $status, $json['error']['code'] ?? null);
        }
        // Unwrap the { success, data } envelope.
        return $json['data'] ?? [];
    }

    /**
     * Initialize a payment and get a checkout_url to redirect the customer to.
     *
     * @param array{
     *   tx_ref: string,
     *   amount: string,
     *   currency?: string,
     *   method?: string,
     *   description?: string,
     *   customer_email?: string,
     *   return_url?: string,
     *   callback_url?: string,
     *   idempotencyKey?: string
     * } $req
     * @return array
     */
    public function initialize(array $req): array
    {
        return $this->request('POST', 'transactions/initialize', [
            'tx_ref' => $req['tx_ref'],
            'amount' => $req['amount'],
            'currency' => $req['currency'] ?? 'ETB',
            'method' => $req['method'] ?? null,
            'description' => $req['description'] ?? null,
            'customer_email' => $req['customer_email'] ?? null,
            'return_url' => $req['return_url'] ?? null,
            'callback_url' => $req['callback_url'] ?? null,
        ], $req['idempotencyKey'] ?? null);
    }

    /** Server-side verification of a payment by your tx_ref. */
    public function verify(string $txRef): array
    {
        return $this->request('GET', 'transactions/verify/' . rawurlencode($txRef));
    }

    /** Get a single payment by id. */
    public function getPayment(string $id): array
    {
        return $this->request('GET', 'transactions/' . rawurlencode($id));
    }

    /** Create a shareable payment link (hosted checkout). */
    public function createPaymentLink(array $input): array
    {
        return $this->request('POST', 'payment_links', [
            'amount' => $input['amount'],
            'currency' => $input['currency'] ?? 'ETB',
            'description' => $input['description'] ?? null,
        ]);
    }

    /**
     * Verify an inbound webhook HMAC signature.
     * ApexPay signs each delivery with X-ApexPay-Signature = HMAC-SHA256(secret, rawBody).
     */
    public static function verifyWebhookSignature(string $signingSecret, string $rawBody, string $signature): bool
    {
        $expected = hash_hmac('sha256', $rawBody, $signingSecret);
        return hash_equals($expected, $signature);
    }
}
