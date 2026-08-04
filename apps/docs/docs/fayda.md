# Fayda ID Verification — Front/Back <2MB + Selfie + OTP Consent per id.gov.et

Per National ID Program Ethiopia — Fayda is 12-digit unique identification number issued by NIDP to residents who fulfill required procedures [OpenG2P](https://docs.openg2p.org/social-registry/features/id-integration/fayda-id-integration), [Trinsic](https://www.trinsic.id/solutions/ethiopia)

## Fayda FIN/FAN Concepts

- **FIN** — Fayda Identification Number 12-digit unique identification number sent directly via SMS [Aratek](https://www.aratek.co/news/national-id-ethiopia-a-gateway-to-digital-id-empowerment)
- **FAN** — Fayda Alias Number alias 16 chars for privacy per id.gov.et/authentication
- **Fayda Card** — physical ID card contains number + 2D barcode for offline use per [Aratek](https://www.aratek.co/news/national-id-ethiopia-a-gateway-to-digital-id-empowerment), QR code scanned via Fayda App Scanner for authenticity per [National ID Ethiopia](https://nationalidethiopia.com/fayda-national-id/)

## Verification Methods per id.gov.et/api + id.gov.et/authentication

### 1. OTP Based - TOTP

Per spec: "Using the FAYDA's authentication service a registered authentication partner can request for OTP authentication. Before performing OTP based authentication the Partner needs to request for an OTP using the individual's ID and use it for OTP based authentication." [Confluence](https://nidp.atlassian.net/wiki/spaces/FAPIQ/pages/633733136/Fayda+Platform+API+Specification)

Flow:
1. Partner API Key request with PartnerCode, UseCaseDescription, SupportingInfo, Status etc. Once approved by Partner Manager, Partner provided PartnerAPIKey containing PartnerCode, policy group and policy, issuedOn, validTill, isActive etc per [Confluence](https://nidp.atlassian.net/wiki/spaces/FAPIQ/pages/633733136/Fayda+Platform+API+Specification)
2. Request OTP using FIN/FAN — customer receives verification code OTP for added security, giving consent [id.gov.et/authentication](https://id.gov.et/authentication)
3. Verify OTP with demographics: name language eng/amh value Milkon Bulcha / ሚልኮን ቡልቻ, gender male/ወንድ, age 25, dob 25/11/1990, fullAddress Woreda01 Yeka Addis Ababa per [Confluence](https://nidp.atlassian.net/wiki/spaces/FAPIQ/pages/633733136/Fayda+Platform+API+Specification)

ApexPay implements: `POST /v1/onboarding/fayda/verify/init` FIN 12-digit + front/back <2MB + selfie + OTP consent + consent IP logged + request_id ULID + fayda_transaction_id mock_tx, `POST /v1/onboarding/fayda/verify/confirm` OTP 6-digit mock 123456 always success face 0.92 per MockVerifier

### 2. Offline QR — FaydaEncode

"You can get your Fayda verified without the need for network connectivity, simply by either using the QR code (we also call the QR FaydaEncode) on your Fayda Credential (id.et/credentials)." [id.gov.et/authentication](https://id.gov.et/authentication)

QR contains masked data + signature NIDP public key verify offline — ApexPay `POST /v1/onboarding/fayda/verify/qr` QRData FIN_LAST4|NAME|DOB|SIG signature valid

### 3. OIDC with Fayda eSignet

"OIDC with Fayda eSignet: OpenID Connect is an identity layer built on top of OAuth 2.0 protocol, providing a way to verify user identity in addition to handling access authorization. When integrated with Fayda eSignet, OIDC enables secure, streamlined user authentication and access management tailored specifically for Ethiopian users. Fayda eSignet serves as the OpenID provider, managing user identity verification and authorization" [id.gov.et/api](https://id.gov.et/api)

ApexPay: OIDC eSignet exchange code for id_token JWT verify FIN

### 4. VeriFayda 1.0 and 2.0 eKYC

"Online use cases leverage advanced digital solutions. Veri Fayda 1.0 and 2.0 enable seamless eKYC with features like one-time passwords (OTP), facial recognition via selfies, and fingerprint scans, ensuring robust and secure digital identity verification." [id.gov.et/authentication](https://id.gov.et/authentication)

ApexPay face_match threshold 0.85, demographics_match boolean, face_match_score 0.92

## Privacy Best Practice per DATABASE v1.1.0

- Never store plain FIN — store `fin_hash = sha256(salt+fin)` + `fin_last4` per crypto.go HashFIN + Last4
- FIN_last4 only returned to UI masked ****-****-1234 per MaskFINLast4
- Front/back images <2MB per NIDP spec — each file must be under 2MB per [National ID Ethiopia](https://nationalidethiopia.com/fayda-national-id/) "Remember, each file must be under 2MB"
- Encrypted at rest AES-256 SSE-S3 MinIO presigned 15m TTL, hash integrity sha256 file_hash unique index per merchant
- Consent timestamp + IP logged per id.gov.et "Service providers can perform e-KYC by requesting the customer's Fayda Alias Number (FAN) or Fayda Identification Number (FIN) on our secure portal. Customers will receive a verification code (OTP) for added security, giving consent" [id.gov.et/authentication]
- No plain FIN in logs — grep logs test in CI ensures no 12-digit pattern

## Outstanding UI — FaydaCapture.tsx

- Video 1280x720 facingMode environment front/back user selfie
- Corner guides L shape 8x8 border-l-4 border-t-4 white rounded-tl-xl animate-pulse
- Glare detection canvas 100x100 sampled 400ms brightness >200 ratio>15% warning "Move to shade • ጥላ ውስጥ ይሂዱ" with brightness value — optimal sampling every 4th pixel O(n/4)
- Capture toBlob JPEG 0.85 <2MB per NIDP, hash integrity sha256, encrypted MinIO SSE-S3 presigned 15m

Next: [API Reference — 21 Paths + 2FA >5000 ETB per ONPS/10/2025](/docs/api-reference)
