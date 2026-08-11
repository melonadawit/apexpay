"""ApexPay Python SDK — payments-only.

A minimal client for merchants who want just the payment gateway for their
e-commerce / online store (like Chapa, Arif Pay, or Telebirr), without payroll,
HR, ledger, or any other module. See docs/PAYMENTS_ONLY_GUIDE.md.

Requires Python 3.8+ and the `requests` library (pip install requests).
"""

from .client import ApexPay, ApexPayError, verify_webhook_signature

__all__ = ["ApexPay", "ApexPayError", "verify_webhook_signature"]
__version__ = "0.1.0"
