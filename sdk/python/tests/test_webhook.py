import hmac
import hashlib

from apexpay import verify_webhook_signature


def test_valid_signature():
    secret = "whsec_secret_123"
    body = '{"event_type":"payment.succeeded","status":"succeeded"}'
    sig = hmac.new(secret.encode(), body.encode(), hashlib.sha256).hexdigest()
    assert verify_webhook_signature(secret, body, sig) is True


def test_wrong_secret():
    body = '{"event_type":"payment.succeeded"}'
    sig = hmac.new(b"wrong-secret", body.encode(), hashlib.sha256).hexdigest()
    assert verify_webhook_signature("right-secret", body, sig) is False


def test_tampered_body():
    secret = "whsec_secret"
    body = '{"event_type":"payment.succeeded"}'
    sig = hmac.new(secret.encode(), body.encode(), hashlib.sha256).hexdigest()
    assert verify_webhook_signature(secret, '{"event_type":"payment.failed"}', sig) is False


if __name__ == "__main__":
    test_valid_signature()
    test_wrong_secret()
    test_tampered_body()
    print("all python webhook tests passed")
