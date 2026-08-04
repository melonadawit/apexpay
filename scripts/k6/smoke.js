import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend, Counter } from 'k6/metrics';

// Custom metrics per MVP §11
const ledgerPostTrend = new Trend('ledger_post_seconds');
const faydaVerifyTrend = new Trend('fayda_verify_duration');
const routingFallbackCounter = new Counter('routing_fallback_used_total');
const payrollCalcTrend = new Trend('payroll_calc_duration');

// K6 smoke + soak per NFR: init p95 <300ms local ex-rail, p99 <150ms staging ex-rail, payroll 500 emp <2s, RAG <1.5s

export const options = {
  stages: [
    { duration: '30s', target: 20 }, // ramp to 20 VUs
    { duration: '1m', target: 100 }, // soak 100 VUs 1m
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    'http_req_duration{scenario:default}': ['p(95)<300', 'p(99)<500'],
    'http_req_failed': ['rate<0.01'],
    'ledger_post_seconds': ['p(99)<30'],
    'payroll_calc_duration': ['p(99)<2000'],
  },
};

const BASE = __ENV.API_URL || 'http://localhost:8080';
const API_KEY = __ENV.API_KEY || 'sk_test_abc123';

export default function () {
  const headers = { 'Authorization': `Bearer ${API_KEY}`, 'Content-Type': 'application/json', 'X-Request-Id': `k6-${__VU}-${__ITER}` };

  group('Onboarding + Fayda', () => {
    const kycRes = http.post(`${BASE}/v1/onboarding/kyc`, JSON.stringify({
      merchant_id: 'mer_test', legal_name: 'Test PLC', business_type: 'plc', registration_number: 'MT/AA/123', tin_number: '0023456789',
      industry_category: 'e-commerce', business_description: 'test', region: 'Addis Ababa', city: 'Addis Ababa', office_address_full: 'Bole',
      contact_person_name: 'Abebe', contact_person_role: 'owner', contact_email: 'test@test.et', contact_phone: '0911111111', expected_monthly_tpv: '500000'
    }), { headers });
    check(kycRes, { 'kyc created 201': (r) => r.status === 201 || r.status === 200 });

    // Fayda mock OTP 123456
    const faydaInit = http.post(`${BASE}/v1/onboarding/fayda/verify/init`, JSON.stringify({
      merchant_id: 'mer_test', owner_id: 'own_1', kyc_profile_id: 'kyc_1', fin: '123456789012', method: 'otp', front_file_key: 'merchants/mer_test/kyc/fayda_front.jpg', back_file_key: 'merchants/mer_test/kyc/fayda_back.jpg'
    }), { headers });
    check(faydaInit, { 'fayda otp sent': (r) => r.status === 201 });

    const faydaConfirm = http.post(`${BASE}/v1/onboarding/fayda/verify/confirm`, JSON.stringify({ request_id: 'fayda_test_123', otp: '123456' }), { headers });
    faydaVerifyTrend.add(faydaConfirm.timings.duration);
    check(faydaConfirm, { 'fayda verified': (r) => r.status === 200 || r.status === 201 });
  });

  group('Payments + Routing + 2FA', () => {
    const txRef = `txr_k6_${__VU}_${__ITER}`;
    const initRes = http.post(`${BASE}/v1/transactions/initialize`, JSON.stringify({
      tx_ref: txRef, amount: '250.00', currency: 'ETB', method: 'telebirr', description: 'k6 smoke', customer_email: 'cust@test.et'
    }), { headers });

    check(initRes, {
      'initialize 201': (r) => r.status === 201,
      'has checkout_url': (r) => r.json().checkout_url !== undefined,
      'requires_2fa boolean present': (r) => r.json().requires_2fa !== undefined,
    });

    // Verify
    const verifyRes = http.get(`${BASE}/v1/transactions/verify/${txRef}`, { headers });
    check(verifyRes, { 'verify 200': (r) => r.status === 200 });

    // Methods ranked
    const methodsRes = http.get(`${BASE}/v1/methods?amount=1000&currency=ETB`, { headers });
    check(methodsRes, {
      'methods ranked 200': (r) => r.status === 200,
      'methods has chosen': (r) => JSON.stringify(r.body).includes('chosen'),
    });
    if (methodsRes.body && methodsRes.body.includes('fallback')) routingFallbackCounter.add(1);
  });

  group('Refunds + Payouts + Payroll + RAG + Swarm', () => {
    // Refund
    const refundRes = http.post(`${BASE}/v1/refunds`, JSON.stringify({ payment_id: 'pay_test', refund_ref: `ref_k6_${__VU}_${__ITER}`, amount: '100.00', currency: 'ETB', fee_policy: 'pro_rata' }), { headers });
    check(refundRes, { 'refund 201': (r) => r.status === 201 || r.status === 200 });

    // Payroll calc - mock but measure duration
    const payrollStart = Date.now();
    const payrollRes = http.post(`${BASE}/v1/payroll_runs`, JSON.stringify({ merchant_id: 'mer_test', run_ref: `prun_k6_${__VU}_${__ITER}`, period_month: 7, period_year: 2026, type: 'regular' }), { headers });
    const payrollCalc = http.post(`${BASE}/v1/payroll_runs/${payrollRes.json().id || 'prun_test'}/calculate`, null, { headers });
    payrollCalcTrend.add(Date.now() - payrollStart);
    check(payrollCalc, { 'payroll calc 200': (r) => r.status === 200 || r.status === 201 });

    // RAG ask
    const ragRes = http.post(`${BASE}/v1/compliance/ask`, JSON.stringify({ query: 'When is 2FA required?', lang: 'en', top_k: 5 }), { headers });
    check(ragRes, {
      'rag 200': (r) => r.status === 200,
      'rag has citation': (r) => r.body.includes('citation') || r.body.includes('5000'),
      'rag p95 <1.5s': (r) => r.timings.duration < 1500,
    });

    // Swarm
    const swarmRes = http.post(`${BASE}/v1/swarm/run`, JSON.stringify({ goal: 'Create link 100 ETB for coffee' }), { headers });
    check(swarmRes, { 'swarm 201': (r) => r.status === 201 });
  });

  sleep(1);
}
