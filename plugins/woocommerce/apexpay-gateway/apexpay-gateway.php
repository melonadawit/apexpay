<?php
/**
 * Plugin Name:       ApexPay for WooCommerce
 * Plugin URI:        https://apexpay.et
 * Description:       Accept Telebirr, CBE Birr, bank, card, and EthSwitch QR payments via ApexPay. Payments-only — no payroll/HR required.
 * Version:           0.1.0
 * Author:            ApexPay
 * Text Domain:       apexpay-gateway
 * Requires at least: 6.0
 * Requires PHP:      8.0
 * License:           Proprietary
 *
 * @package ApexPay\WooCommerce
 */

defined( 'ABSPATH' ) || exit;

/**
 * Register the gateway with WooCommerce.
 */
function apexpay_init_gateway_class(): void {
    if ( ! class_exists( 'WC_Payment_Gateway' ) ) {
        return; // WooCommerce not active
    }

    class WC_Gateway_ApexPay extends \WC_Payment_Gateway {

        /** @var string ApexPay secret key (sk_test_... / sk_live_...). */
        public $secret_key = '';

        /** @var string Webhook signing secret. */
        public $webhook_secret = '';

        /** @var string Default rail: telebirr | cbe_birr | bank | card_acquirer | ethswitch. */
        public $default_method = 'telebirr';

        /** @var string API base URL. */
        public $base_url = 'https://api.apexpay.et';

        public function __construct() {
            $this->id                 = 'apexpay';
            $this->has_fields         = false;
            $this->method_title       = __( 'ApexPay', 'apexpay-gateway' );
            $this->method_description = __( 'Accept Telebirr, CBE Birr, bank, card and EthSwitch QR via ApexPay.', 'apexpay-gateway' );

            $this->supports = array( 'products', 'refunds' );

            $this->init_form_fields();
            $this->init_settings();

            $this->title          = $this->get_option( 'title' );
            $this->description    = $this->get_option( 'description' );
            $this->secret_key     = $this->get_option( 'secret_key' );
            $this->webhook_secret = $this->get_option( 'webhook_secret' );
            $this->default_method = $this->get_option( 'default_method' );
            $this->base_url       = $this->get_option( 'base_url' );

            add_action( 'woocommerce_update_options_payment_gateways_' . $this->id, array( $this, 'process_admin_options' ) );
            add_action( 'woocommerce_api_apexpay_webhook', array( $this, 'handle_webhook' ) );
            add_action( 'woocommerce_thankyou_apexpay', array( $this, 'thankyou' ) );
        }

        public function init_form_fields(): void {
            $this->form_fields = array(
                'enabled' => array(
                    'title'   => __( 'Enable/Disable', 'apexpay-gateway' ),
                    'type'    => 'checkbox',
                    'label'   => __( 'Enable ApexPay payments', 'apexpay-gateway' ),
                    'default' => 'no',
                ),
                'title' => array(
                    'title'       => __( 'Title', 'apexpay-gateway' ),
                    'type'        => 'text',
                    'default'     => __( 'Pay via Telebirr / CBE Birr', 'apexpay-gateway' ),
                ),
                'description' => array(
                    'title' => __( 'Description', 'apexpay-gateway' ),
                    'type'  => 'textarea',
                ),
                'secret_key' => array(
                    'title'       => __( 'ApexPay Secret Key', 'apexpay-gateway' ),
                    'type'        => 'password',
                    'description' => __( 'sk_test_... or sk_live_... from the Developer portal.', 'apexpay-gateway' ),
                ),
                'webhook_secret' => array(
                    'title'       => __( 'Webhook Signing Secret', 'apexpay-gateway' ),
                    'type'        => 'password',
                    'description' => __( 'Used to verify ApexPay webhook signatures.', 'apexpay-gateway' ),
                ),
                'default_method' => array(
                    'title'   => __( 'Default payment method', 'apexpay-gateway' ),
                    'type'    => 'select',
                    'options' => array(
                        'telebirr'       => 'Telebirr',
                        'cbe_birr'       => 'CBE Birr',
                        'bank'           => 'Bank / IPS',
                        'card_acquirer'  => 'Card',
                        'ethswitch'      => 'EthSwitch QR',
                    ),
                    'default' => 'telebirr',
                ),
                'base_url' => array(
                    'title'   => __( 'API Base URL', 'apexpay-gateway' ),
                    'type'    => 'text',
                    'default' => 'https://api.apexpay.et',
                ),
            );
        }

        /**
         * Called when the customer clicks "Place order".
         *
         * @param int $order_id
         * @return array
         */
        public function process_payment( $order_id ): array {
            $order = wc_get_order( $order_id );
            require_once __DIR__ . '/includes/ApexPay.php';

            $apexpay = new \ApexPay\ApexPay( $this->secret_key, $this->base_url );

            try {
                $payment = $apexpay->initialize(
                    array(
                        'tx_ref'         => 'woo-' . $order_id,
                        'amount'         => number_format( (float) $order->get_total(), 2, '.', '' ),
                        'currency'       => $order->get_currency() ?: 'ETB',
                        'method'         => $this->default_method,
                        'customer_email' => $order->get_billing_email(),
                        'callback_url'   => WC()->api_request_url( 'apexpay_webhook' ),
                        'return_url'     => $this->get_return_url( $order ),
                    ),
                    'woo-' . $order_id,
                );
            } catch ( \Exception $e ) {
                wc_add_notice( 'ApexPay: ' . $e->getMessage(), 'error' );
                return array( 'result' => 'failure' );
            }

            // Mark as on-hold; we confirm via webhook before completing.
            $order->update_status( 'on-hold', __( 'Awaiting ApexPay payment.', 'apexpay-gateway' ) );
            $order->update_meta_data( '_apexpay_payment_id', $payment['id'] );
            $order->save();

            return array(
                'result'   => 'success',
                'redirect' => $payment['checkout_url'],
            );
        }

        /**
         * Handle the ApexPay webhook. Authoritative order fulfilment.
         */
        public function handle_webhook(): void {
            $raw_body   = file_get_contents( 'php://input' );
            $signature  = $_SERVER['HTTP_X_APEXPAY_SIGNATURE'] ?? '';

            require_once __DIR__ . '/includes/ApexPay.php';
            if ( ! \ApexPay\ApexPay::verifyWebhookSignature( $this->webhook_secret, $raw_body, $signature ) ) {
                status_header( 401 );
                exit( 'invalid signature' );
            }

            $event = json_decode( $raw_body, true );
            if ( ( $event['event_type'] ?? '' ) === 'payment.succeeded' ) {
                $order = wc_get_order( (int) str_replace( 'woo-', '', $event['tx_ref'] ?? '' ) );
                if ( $order && $order->needs_payment() ) {
                    $order->payment_complete( $event['payment_id'] ?? '' );
                }
            }
            status_header( 200 );
            exit( 'ok' );
        }

        /** Optional refund support. */
        public function can_refund_pending(): bool {
            return false;
        }
    }
}

add_action( 'plugins_loaded', function () {
    apexpay_init_gateway_class();
    add_filter( 'woocommerce_payment_gateways', function ( $methods ) {
        if ( class_exists( 'WC_Gateway_ApexPay' ) ) {
            $methods[] = 'WC_Gateway_ApexPay';
        }
        return $methods;
    } );
} );
