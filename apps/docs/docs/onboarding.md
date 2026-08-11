# Merchant Onboarding — NBE-Grade 6-Step Wizard Outstanding (Like Mercury + Linear + Stripe Atlas)

Per NBE ONPS/02/2020 Payment System Operator Directive + ONPS/09/2023 + ONPS/10/2025 + PayAtlas ET PSP KYC + ApexPay/ApexPay onboarding reference.

## Required Docs Checklist Per PayAtlas + ApexPay/ApexPay + NBE

| Document | Required For | Notes | NBE Ref |
|---|---|---|---|
| Company Registration | KYC | Must be notarized; official documents Amharic or English | ONPS/02/2020 § business registration license |
| Passport or National ID of UBO | KYC | Clear copies + notarization often requested — **Fayda front/back <2MB + selfie + OTP consent id.gov.et for Ethiopian nationals** | KYC io Tiered |
| TIN (Tax Identification Number) | KYC | 10-digit ET TIN mandatory proof tax registration | PayAtlas |
| Bank Account Statement | Risk Review | Recent statements verify financial stability + settlement account must match legal name per validation fuzzy Levenshtein <3 | PayAtlas |
| Processing History | Risk Review | Optional but beneficial high-volume merchants | PayAtlas |
| Website URL & Policy Docs | Website Review | Must include refund, privacy, terms & conditions pages per PayAtlas ET PSP | PayAtlas + NBE ONPS/10/2025 refund policy requirement |
| Local Business License | KYC | Required certain regulated industries education/health | PayAtlas |

ApexPay model: Businesses must first pass KYC checks, provide business info like TIN, business licences, address, contact info, receive approval before accepting payments. Simple integration 6 lines copy-paste. [GxCamp](https://gxcamp.com/blog-details/the-first-ethiopias-online-payment-gateway-service--ApexPay-payment--)

## 6-Step Wizard Outstanding UI/UX (Like Stripe Atlas + Mercury)

### Step 1: Business Info • የንግድ መረጃ

Form fields with smart validation:
- Legal Name • ህጋዊ ስም * (must == account name fuzzy <3)
- Trade Name • የንግድ ስም
- Business Type • የንግድ አይነት *: sole_proprietorship, plc, **share_company min 5 shareholders** per Commercial Code, 10 if multi-system per ONPS/02/2020 §, max 40% single person ownership
- TIN Number • ቲን ቁጥር 10-digit * pattern `^[0-9]{10}$` ET TIN validation
- Registration No • ምዝገባ ቁጥር *
- Industry • ዘርፍ *: e-commerce, education, tech, food, logistics — **Restricted gambling/crypto/adult blocked per NBE** 🚫 — PayAtlas Restricted Industries gambling, cryptocurrencies, adult content strict restrictions or outright bans
- Business Description • የንግድ መግለጫ * per PayAtlas
- Website URL • ድህረ ገጽ + Expected Monthly TPV • ወርሃዊ ግምት ETB * — if TPV>1M ETB → dual approval required maker-checker
- Region • ክልል *: Addis Ababa, Oromia... City • ከተማ *: Addis Ababa, Full Office Address • ሙሉ አድራሻ * + GPS pin optional map picker outstanding
- Contact Person Name/Role/Email/Phone

Outstanding UI: glassmorphic nav `backdrop-blur-xl bg-white/70` border white/50 shadow-glass, card border 1px rgba(0,0,0,0.06) elevated hover medium, stepper animated line pathLength motion.div, DonutProgress 64px progress 0-100% animated

### Step 2: Owners & Fayda • ባለቤቶች እና ፋይዳ

UBO >10% per NBE + authorized signatory mandatory at least one:
- Full Name EN + AM • ሙሉ ስም + Full Name AM • አማርኛ
- Role: owner, shareholder, director, authorized_rep, ubo, contact_person
- Ownership % 0-100 sum 100% for share company
- Nationality default ET, ID Type fayda/passport, Phone Email, Is PEP checkbox, Is Authorized Signatory checkbox
- Fayda Verification: FIN 12-digit or FAN 16 alias + front image camera overlay corner guides animated pulse + back image + selfie liveness + OTP 6-digit mock 123456 + consent timestamp IP logged per id.gov.et

**FaydaCapture Outstanding Component:**
- Video 1280x720 facingMode environment front/back user selfie, corner guides L shape 8x8 border-l-4 border-t-4 white rounded-tl-xl animate-pulse
- Glare detection canvas 100x100 sampled 400ms brightness >200 ratio>15% warning "Move to shade • ጥላ ውስጥ ይሂዱ" with brightness value
- Capture toBlob JPEG 0.85 <2MB per NIDP spec, hash integrity sha256, encrypted MinIO SSE-S3 presigned 15m, FIN never logged only last4 + hash
- Offline QR alternative FaydaEncode scan without network + OIDC eSignet OIDC flow id.gov.et/api PartnerCode PartnerAPIKey UseCaseDescription SupportingInfo Status etc per spec

### Step 3: Bank Account • የባንክ ሂሳብ Settlement

- Bank select combobox with logos from GET /v1/banks CBE/Awash/Dashen/Abyssinia/Berhan/Wegagen/NIB/United/Coop/Oromia/Bunna/Lion/Zemen/CBO — 14 ET banks seeded
- Account Name must == Legal Name fuzzy Levenshtein <3 check O(n*m) DP optimal per banking verification name_match or require override note
- Account Number masked ****1234 display per privacy, hash sha256 for lookup per DATABASE privacy rule, account_number_hash index
- Is Settlement Default boolean, verification_status pending/verified/failed, verification_method bank_letter/micro_deposit/manual

Outstanding UI: bank combobox with logos, account name auto-check must match legal

### Step 4: Documents Vault • ሰነዶች ማከማቻ

Dropzone outstanding:
- Dashed border 2px rounded-2xl pulse on drag scale 0.98 pulseGlow boxShadow 0 0 20px rgba(11,110,79,0.2) → 0 0 30px rgba(11,110,79,0.4) infinite
- File thumbs PDF icon + JPG thumb 64px, upload progress slim top bar, status chip check animate, required checklist progress donut 0-100%
- Required docs computed per business_type + KYC level via RequiredDocs() unique O(n) per domain.go: base tin_certificate business_license proof_of_address bank_letter + if PLC/share_company company_registration memorandum_articles board_resolution shareholder_list + if Level2/3 fayda_card_front fayda_card_back + website_screenshot refund_policy_doc
- File upload via presigned POST TTL 15m directly to MinIO `merchants/{merchant_id}/kyc/{doc_type}_{id}.pdf` no server buffering O(n) streaming sha256, hash integrity file_hash unique index per merchant doc, mime whitelist pdf/jpg/png size <5MB Fayda <2MB per NIDP, status pending → ocr_done → verified|rejected by compliance ops with rejection_reason, OCR extracts fields registration_no TIN expiry side-by-side in admin viewer, MinIO versioning + encryption SSE-S3 retention 7y per NBE, expires_at business_license not expired validation

### Step 5: Compliance Preview • ተገዢነት ቅድመ-እይታ

Automated checks per NBE ONPS + PayAtlas + ApexPay model, risk_score weighted sum optimal:
- tin_validation format 10 digits valid
- business_license_validation not expired
- bank_account_validation name match Levenshtein 1
- fayda_verification OTP verified face 0.92
- restricted_industry E-commerce allowed not gambling/crypto/adult per PayAtlas Restricted Industries
- website_policy_check refund/privacy/terms found via crawl
- aml_screening No sanctions hit, pep_check, sanctions
- risk_scoring 42/100 Medium TPV 500k ETB, if risk high >=70 or TPV>1M ETB → dual approval required maker-checker per NBE capacity assessment doc 30-60 days pilot analogy

Outstanding UI: risk gauge chart 42/100, checks green/red cards progress bars, Kanban board outstanding

### Step 6: Review & Submit • ግምገማ እና አስገባ

Summary cards glass:
- Business Legal TIN Industry Region
- Owners & Fayda FIN ****-****-1234 Verified 0.92
- Bank CBE ****1234 Default settlement
- Documents 4/6 uploaded hash_company_reg
- Compliance Preview risk 42 medium, dual approval required if high risk or TPV>1M

Consent checkboxes:
- I confirm business info true per NBE ONPS/02/2020 directive and capacity to manage payment gateway per capital ETB 3M requirement ONPS/02/2020 § minimum capital payment gateway operator ETB 3M
- I consent Fayda verification via id.gov.et VeriFayda 2.0 / OIDC eSignet with OTP consent, understanding FIN stored as sha256 hash + last4 only, front/back images encrypted at rest AES-256, presigned URLs 15m TTL, no plain PII in logs
- I agree website has refund, privacy, terms pages per PayAtlas ET PSP requirement, business not in restricted industries gambling/crypto/adult per PayAtlas Restricted Industries
- I understand risk scoring medium, dual approval if high risk or TPV>1M ETB, test keys immediately, live keys after pilot 30-60 days analogy NBE, 2FA mandatory >5000 ETB per ONPS/10/2025

After submit: compliance team reviews in Kanban board outstanding Submitted → In Review → Fayda Pending → Compliance Check → Approved, Timeline vertical like Linear, email with confetti animation, Ledger merchant operating book created + 6 accounts seeded, outbox merchant.activated triggers email

### Tracking Timeline • የጊዜ መስመር

Outstanding vertical like Linear:
- Business info saved draft • 1 min ago
- Fayda verification OTP sent • 123456 mock • face score 0.92
- Bank added • CBE ****1234 verified
- Docs 4/6 uploaded 66% • hash integrity sha256
- Compliance checks 8 green/red • risk 42 medium • dual approval needed if high
- Submitted • pending_kyc → kyc_in_review → fayda_pending → compliance_check → pending_approval → active • Progress 100% • Test keys immediately, live after pilot 30-60 days

Next: [Fayda ID Front/Back OTP](/docs/fayda) with camera overlay guides corners glare detection
