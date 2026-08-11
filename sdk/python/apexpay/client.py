"""ApexPay payments-only Python client."""
from __future__ import annotations

import hashlib
import hmac
from typing import Any, Dict, Optional

try:
    import requests
except ImportError:  # pragma: no cover
    requests = None  # type: ignore[assignment]


class ApexPayError(Exception):
    """Raised when the ApexPay API returns an error or a request fails."""

    def __init__(self, message: str, status_code: Optional[int] = None, code: Optional[str] = None):
        super().__init__(message)
        self.status_code = status_code
        self.code = code


def verify_webhook_signature(signing_secret: str, raw_body: str, signature: str) -> bool:
    """Verify an ApexPay webhook HMAC signature.

    ApexPay signs each delivery with ``X-ApexPay-Signature`` =
    ``HMAC-SHA256(signing_secret, raw_body)``.
    """
    expected = hmac.new(
        signing_secret.encode("utf-8"),
        raw_body.encode("utf-8"),
        hashlib.sha256,
    ).hexdigest()
    return hmac.compare_digest(expected, signature)


class ApexPay:
    """Minimal ApexPay client for payments-only integrations."""

    def __init__(self, api_key: str, base_url: str = "http://localhost:8080"):
        if not api_key:
            raise ApexPayError("ApexPay: api_key is required (sk_test_... or sk_live_...)")
        if requests is None:
            raise ApexPayError("ApexPay: the 'requests' library is required (pip install requests)")
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")

    def _request(
        self,
        method: str,
        path: str,
        body: Optional[Dict[str, Any]] = None,
        idempotency_key: Optional[str] = None,
    ) -> Dict[str, Any]:
        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
        }
        if idempotency_key:
            headers["Idempotency-Key"] = idempotency_key

        resp = requests.request(method, f"{self.base_url}/v1/{path}", headers=headers, json=body, timeout=30)
        try:
            payload = resp.json()
        except ValueError:
            raise ApexPayError("ApexPay returned an invalid response", resp.status_code)

        if resp.status_code >= 400:
            err = payload.get("error", {}) if isinstance(payload, dict) else {}
            msg = err.get("message") or payload.get("message") or f"ApexPay request failed ({resp.status_code})"
            raise ApexPayError(msg, resp.status_code, err.get("code"))

        # Unwrap the { success, data } envelope.
        return payload.get("data", {}) if isinstance(payload, dict) else {}

    def initialize(self, tx_ref: str, amount: str, currency: str = "ETB",
                   method: Optional[str] = None, description: Optional[str] = None,
                   customer_email: Optional[str] = None, return_url: Optional[str] = None,
                   callback_url: Optional[str] = None, idempotency_key: Optional[str] = None) -> Dict[str, Any]:
        """Initialize a payment and return a payment with a checkout_url."""
        return self._request(
            "POST",
            "transactions/initialize",
            {
                "tx_ref": tx_ref,
                "amount": amount,
                "currency": currency,
                "method": method,
                "description": description,
                "customer_email": customer_email,
                "return_url": return_url,
                "callback_url": callback_url,
            },
            idempotency_key,
        )

    def verify(self, tx_ref: str) -> Dict[str, Any]:
        """Server-side verification of a payment by your tx_ref."""
        from urllib.parse import quote
        return self._request("GET", f"transactions/verify/{quote(tx_ref)}")

    def get_payment(self, payment_id: str) -> Dict[str, Any]:
        """Get a single payment by id."""
        from urllib.parse import quote
        return self._request("GET", f"transactions/{quote(payment_id)}")

    def create_payment_link(self, amount: str, currency: str = "ETB", description: Optional[str] = None) -> Dict[str, Any]:
        """Create a shareable payment link (hosted checkout)."""
        return self._request(
            "POST",
            "payment_links",
            {"amount": amount, "currency": currency, "description": description},
        )

    verify_webhook_signature = staticmethod(verify_webhook_signature)
