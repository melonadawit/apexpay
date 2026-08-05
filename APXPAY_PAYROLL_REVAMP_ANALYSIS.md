# ApexPay Payroll — From Basic Run to RazorpayX-Grade Comprehensive Money OS
**Date:** 2026-08-05
**Module:** `services/api/internal/payroll` + `apps/merchant-web/app/payroll`
**Goal:** Make payroll as comprehensive as RazorpayX Payroll (India) but Ethiopia-native and more powerful with AI

---

## 1. Current State — What You Built Well

### ✅ Already Excellent
- **Migration 0008:** employees (employee_code unique per merchant, Fayda hash, bank masked/hash, cost_center), payroll_tax_brackets versioned effective_from/to, payroll_runs book_id per run (Ledger M4), payroll_items gross/taxable/income_tax/pension 7%/11%/net, payroll_claims expense/medical/travel.
- **Domain:** RunType regular/off_cycle/bonus/adjustment, Status FSM draft->calculating->pending_approval->approved->processing->completed->failed->voided, EmploymentType permanent/contract/part_time/intern.
- **Service:** `CalculateTax` binary search O(log n) over sorted brackets `tax=taxable*rate-deduction` rounded 2 decimals, `CalculateRun` O(n) active employees only, pension 7% employee 11% employer, hourly_rate = base/208 (26 days *8h) ET standard, OT rates map weekday 1.25 weekend 1.5 holiday 2.0 night 1.3. `ApproveRun` dual >100k net, `DisburseRun` ledger M4 Dr expense:salary Gross Cr payroll_payable Net Cr et_income_tax_payable Tax Cr pension_payable Pension + second journal Dr payroll_payable Cr clearing:bank Net via payout batch.
- **Repo:** `CreateRunBookTx` Tx ledger_books + ledger_journals + ledger_entries atomic, `BulkCreateItems` Tx loop, `ListEmployees` active filter, `GetTaxBrackets` ordered by Min ASC effective date versioned.
- **Handler:** POST /employees, GET /employees, POST /payroll_runs, POST /{id}/calculate, approve, disburse, GET items.
- **UI:** merchant-web payroll page 3 cards Total Gross/Total Tax/Total Net, employees 10 Fayda badge + bank masked + cost_center, runs table status pipeline stepper pending_approval needs dual, run detail 8 cols Employee Gross OT Taxable Income Tax Pension 7%/11% Net Status sticky footer totals, buttons Approve Disburse Download Payslips PDF ZIP ET Report CSV ERCA JSON, payslip preview modern template logo QR pie chart breakdown outstanding.

**Grade: B+ calc engine, C product.**

### ❌ What RazorpayX Has That You Don't (Yet)

From research [1](https://www.positioniseverything.net/razorpayx-payroll-pricing-reviews-2026/)[2](https://www.techbloat.com/razorpayx-payroll-pricing-reviews-2026.html)[4](https://www.techjockey.com/detail/razorpayx)[5](https://www.techimply.com/profile/razorpayx-payroll):

RazorpayX Payroll Toolkit:
- Salary structure setup: Basic, HRA, Special Allowance, Employer PF, Gratuity, PT, Income Tax deductions → configurable components Earnings/Deductions/Reimbursements
- Monthly payroll runs: attendance inputs, LOP, arrears/salary revisions, variable pay, bonuses, deductions, F&F settlement support, employee lifecycle onboarding/offboarding
- Automatic salary, tax, deduction calc, compliance with tax & labour laws India PF ESI PT LWF, TDS, Form16, ESIC challan, PT challan generation + payment filing automated (not just report)
- Payslip generation + distribution online, employee self-service portal + mobile apps + WhatsApp payslips + tax declarations + investment declarations regime old/new
- Benefits: healthcare, retirement, PTO, reimbursements, loans/advances, expense claims
- Leave & Attendance management, attendance tracking integration
- 100% compliance automation: calculates, pays, files TDS, PF, ESIC, PT multi-state, Form16 digitally signed
- Direct bank transfers multi-bank bulk file, instant reimbursements
- Employee master report, payroll cost analysis, deduction summaries, reimbursement data, compliance reports, audit trails
- Custom user roles workflow, scheduling payments in advance
- 360° payroll management from salary disbursement to compliance.

You have: base_salary only, no structure, no HRA/allowances breakdown, no LOP, no attendance, no revisions, no reimbursements/loans, no F&F, no self-service portal, no benefits, no compliance filing only ledger, no Form16 equivalent, no department/grade/designation, no probation, no onboarding checklist, no offboarding.

---

## 2. Ethiopian vs India Compliance — Must Localize, Not Copy-Paste

| India (RazorpayX) | Ethiopia (ApexPay) | Notes |
|---|---|---|
| PF 12% employee 12% employer | Pension 7% employee 11% employer Private Organization Employees Social Security Agency — different rates for government? Need config | Must generate pension contribution report monthly CSV for agency, employer 11% cap? No cap placeholder but need rule engine |
| Professional Tax state slabs | No PT in ET, but Employment Income Tax federal 7 brackets 0-600 0% 601-1650 10%-60 1651-3200 15%-142.5 3201-5250 20%-302.5 5251-7800 25%-565 7801-10900 30%-955 >10900 35%-1500 | Your brackets seeded but deduction values slightly off vs actual 2024: you use 60,142.5,302.5,565,955,1500 correct per law? Verify ERCA official. Need versioned effective_from |
| TDS salary old/new regime + Form16 | ERCA withholding income tax monthly, annual employee TIN, need ERCA monthly withholding CSV + annual certificate equivalent | No Form16 but need annual income tax certificate for employee |
| ESI 0.75%/3.25% | No ESI, but medical allowance, private insurance optional | Add benefit type |
| Gratuity, LWF | No gratuity law same but severance per Ethiopian labor law: if illegal termination, severance = 30 days wage per year? Need implement final settlement module | Important for F&F |
| Leave: CL/SL/PL | ET: Annual leave 14 days first year +1 per year up to 35, sick leave 6 months (3mo 100% then 50%?), maternity 120 days (30 pre +90 post), etc | Need leave management |
| HRA | No HRA, but Housing allowance, Transport allowance, Fuel allowance, Per diem, Hardship, Overtime 1.25/1.5/2.0/1.3 night | Your OT rates correct per Labour Proclamation No. 1156/2019 |
| Bonus: Diwali | Bonus: Ethiopian 13th month? Many give extra salary for Ethiopian New Year Enkutatash or upon profit. Plus Sales commission | Keep bonus type commission bonus ex_gratia 13th_month |

**Ethiopian Unique Needs RazorpayX Doesn't Have:**
- Fayda ID verification for employee KYC
- Cost center allocation per employee per project (workforce money OS vision per FULL_PLATFORM_SPEC)
- Pension number + TIN 10-digit validation
- Bank letter + bank account name match fuzzy Levenshtein <3 (already claimed)
- Salary in ETB only but need USD contractor payroll via payout rail Bank IPS
- Payroll approval maker-checker dual >100k net + audit log per NBE PO Operator financial controls software≠license
- Ledger per run book + payout batch integration for disburse via Telebirr/CBE/Bank — RazorpayX uses RazorpayX business banking, you use your own M3/M4 — better!
- Amharic payslip + QR verification

---

## 3. Gap Analysis — Current vs RazorpayX-Grade

| Category | Current ApexPay | RazorpayX Comprehensive | Gap Severity |
|---|---|---|---|
| **Employee Master** | employee_code, name, name_am, email, phone, TIN, Fayda hash, pension_no, bank masked/hash bank_code base_salary employment_date type cost_center status | + Department, Designation, Grade, Location/Branch, Date of Joining, Probation end, Employment status (confirmed, probation, notice), Reporting manager, Work location, Salary structure assigned, UAN/PF equivalent pension, ESIC? insurance No, employment history, documents vault (contract, TIN certificate, Fayda front/back, bank letter), KYC status, Fayda verified badge | **High** |
| **Salary Structure** | base_salary only | CTC Template: Earnings Basic 40-50% of CTC, HRA (ET: Housing), Special Allowance, Transport, Fuel, Medical, LTA, Bonus structure, Employer contributions Pension 11% (in template but not cost to CTC? In ET employee 7% deduction only, employer 11% extra cost), Deductions Income Tax, Pension 7%, Other Deductions, Reimbursements Conveyance etc, Overtime structure, Variable pay commission structure, Formula engine `Basic = CTC * 0.4` etc, revisions history with effective date | **Critical** |
| **Attendance & Leave** | None | Leave Types Annual, Sick, Maternity, Paternity, Unpaid, Comp Off, attendance input Paid Days, LOP Days, Proration logic `net = (gross / total_days)*paid_days`, Monthly variable inputs OT hours, LOP reversal, leave balance | **Critical** |
| **Payroll Runs** | regular/off_cycle/bonus/adjustment draft->calculating->pending_approval->approved->processing->completed + total_gross/net/tax/pension | + Pay schedule monthly/weekly/semimonthly, Pay calendar, Payroll cutoff date, Lock payroll after disbursal, Re-run/amendment, Hold salary for employee, Skip employee, Arrears calculation `new_basic - old_basic * pending_months`, Bonus run with taxable? ex_gratia handling, Adjustment run negative adjustments, Supplementary run | **High** |
| **Variable Pay** | commission, bonus, other_allowances in payroll_items but not input per employee per run UI | Needs variable pay import CSV employee_code, component, amount, notes, bulk upload 1000 rows O(n), validation | **High** |
| **Deductions** | income_tax pension_employee other_deductions | + Loan deduction, Advance deduction, Salary advance recovery, Penalty, Insurance premium, Savings, Union due, Custom deduction, Non-taxable vs taxable classification, pre-tax vs post-tax | **High** |
| **Compliance** | tax brackets seeded, pension 7/11 calc, ledger M4 journal | + Pension contribution report monthly `Format for Private Org Employees Social Security Agency: employee pension No, name, gross, employee 7% employer 11% total 18%` CSV + pension challan file generation + ERCA withholding monthly CSV `TIN, name, gross, taxable, tax, pension, net, month` + Annual tax certificate Form? + Audit trail immutable, compliance calendar reminders | **Critical** for NBE payroll operator |
| **Disbursal** | DisburseRun -> ledger second journal Dr payroll_payable Cr clearing:bank + payout batch creation via payout service but not atomic + payout batch disburses via mock connector | + Bank file generation pain.001 ISO20022 XML or MT103 CSV per bank Code CBE/Awash/Dashen, batch ref, file hash, status tracking file uploaded to bank SFTP/Manual upload, reconciliation bank statement MT940 window 24h amount tolerance 0.01, failure retry per employee status paid/failed/returned, salary on hold escrow book | **High** |
| **Payslip** | PDF modern template logo QR verification breakdown pie - lib/pdf/payslip.ts placeholder | + Outstanding modern template with company logo, employee code, DOJ, designation, department, cost center, bank masked, period, earnings table basic HRA allowance OT commission bonus gross, deductions table tax pension loan etc net, employer contribution table 11% pension, YTD summary Gross Tax Net, QR verification `https://apexpay.et/verify/payslip/{id}` signed JWT, digital signature, email distribution, WhatsApp share, password protected PDF (DOB+last4), Amharic + English bilingual | **Medium** |
| **Self-Service** | None | Employee portal login magic link JWT 24h, view payslips, download Form16 equivalent tax certificate, view YTD, update PAN/TIN? Investment declarations? Not needed ET but allow bank update request + docs upload, leave apply, reimbursement claim expense/travel/medical receipt MinIO file_key, approval workflow manager->finance->disburse, payslip WhatsApp, push notification payroll processed | **Critical** for RazorpayX parity |
| **Benefits / Claims** | payroll_claims table claim_type expense/medical/travel/other amount receipt_file_key status pending/approved/rejected/paid | + Loan: employee request loan amount reason interest? Repayment schedule EMIs deduction auto per payroll run, Advance salary advance up to 50% net recovery next months, Reimbursements expense with receipt attachment approval flow, Bonus ex_gratia 13th month | **High** |
| **Final Settlement (F&F)** | None | Full & Final Settlement: resignation date, last working day, notice period days, LOP for notice, leave encashment annual balance * per day gross/30, gratuity/severance per ET labour law illegal termination severance formula `salary * years * factor` + earning deduction outstanding + loan outstanding + payslip F&F + clearance checklist, workflow hr->finance->approved->disburse single payout, ledger M4 adjusted | **High** |
| **Reporting & Analytics** | Total Gross/Tax/Net cards | RazorpayX reports: Payroll summary month-wise, Employee-wise salary report, Deduction report, Reimbursement report, Compliance report, Employee master report, Cost center report allocation per cost_center sum gross/net, Variance report vs last month %, payroll cost analysis, audit trails for finance investor reporting, year-end tax, downloadable CSV/XLSX/JSON, Recharts line/bar/pie TPV equivalent payroll cost, headcount | **Medium** |
| **Integrations** | Ledger M3/M4 | + Accounting integration ledger per cost_center expense accounts, payroll journal auto-posted to accounting system (QuickBooks/Zoho Books) via webhook, Attendance integration API Zoho People/Spreadsheet CSV import, HRMS 25+ partners, Bank SFTP for payout file, email SMTP for payslip | **Medium** |

---

## 4. Proposed Comprehensive Architecture — Beyond RazorpayX, Ethiopia-Native

### 4.1 Data Model Enhancements (New Tables — Add Migration 0013_payroll_comprehensive.up.sql)

```sql
-- Departments, Designations, Grades, Branches
create table payroll_departments (id text primary key, merchant_id text not null references merchants(id), name text not null, code text, cost_center text);
create table payroll_designations (id text primary key, merchant_id text not null, title text not null, level int);
create table payroll_grades (id text primary key, merchant_id text not null, name text e.g., G1-G10, min_salary numeric(20,8), max_salary numeric(20,8));
create table payroll_branches (id text primary key, merchant_id text not null, name text, region text, city text);

-- Salary Structure Template
create table payroll_salary_structures (
  id text primary key, merchant_id text not null, name text not null e.g., Fixed CTC 500k, description text,
  ctc_annual numeric(20,8), ctc_monthly numeric(20,8) generated,
  effective_from date, status active/archived,
  created_at timestamptz default now()
);
create table payroll_structure_components (
  id text primary key, structure_id text not null references payroll_salary_structures(id) on delete cascade,
  component_type text check (type in ('earning','deduction','employer_contribution','reimbursement')),
  code text e.g., BASIC, HRA/HOUSING, TRANSPORT, FUEL, SPECIAL_ALLOW, OT, COMMISSION, BONUS, MEDICAL, LTA,
  name text, name_am text,
  calculation_type text check (type in ('fixed','percentage_of_basic','percentage_of_ctc','formula')), 
  amount numeric(20,8), percentage numeric(5,2), formula text e.g., "CTC*0.4" or "BASIC*0.1",
  is_taxable boolean default true,
  is_part_of_gross boolean default true, -- if false then employer contribution not in gross
  is_optional boolean default false,
  order_no int, -- display order in payslip
  meta jsonb default '{}' -- e.g., {tax_exempt_limit: 1000}
);

-- Employee Extended
alter table employees add column department_id text references payroll_departments(id);
alter table employees add column designation_id text references payroll_designations(id);
alter table employees add column grade_id text references payroll_grades(id);
alter table employees add column branch_id text references payroll_branches(id);
alter table employees add column reporting_manager_id text references employees(id);
alter table employees add column ctc_annual numeric(20,8);
alter table employees add column salary_structure_id text references payroll_salary_structures(id);
alter table employees add column probation_end_date date;
alter table employees add column confirmation_status text check (status in ('probation','confirmed','notice','terminated')) default 'probation';
alter table employees add column employment_status text check (status in ('active','on_hold','notice_period','terminated','retired')) default 'active';
alter table employees add column documents jsonb default '[]' -- [{type:contract, file_key, file_hash, status}]
alter table employees add column bank_account_name text;
alter table employees add column date_of_joining date not null default current_date;

-- Salary Revisions History
create table payroll_salary_revisions (
  id text primary key, merchant_id text not null, employee_id text not null references employees(id),
  old_base numeric(20,8), new_base numeric(20,8), old_ctc numeric(20,8), new_ctc numeric(20,8),
  old_structure_id text, new_structure_id text,
  effective_from date not null, reason text, approved_by text, status pending/approved,
  arrear_amount numeric(20,8) default 0, -- auto calc new-old * pending months
  created_at timestamptz default now()
);

-- Attendance & Leave Inputs per Run
create table payroll_attendance_inputs (
  id text primary key, run_id text not null references payroll_runs(id) on delete cascade,
  employee_id text not null references employees(id),
  paid_days int not null default 30, lop_days int not null default 0, total_days int not null default 30,
  ot_weekday_hours numeric(10,2) default 0, ot_weekend_hours numeric(10,2) default 0,
  ot_holiday_hours numeric(10,2) default 0, ot_night_hours numeric(10,2) default 0,
  leave_taken jsonb default '{}' e.g., {"annual":2,"sick":1},
  present_days int,
  created_at timestamptz default now(),
  unique (run_id, employee_id)
);

-- Variable Pay Inputs
create table payroll_variable_inputs (
  id text primary key, run_id text not null references payroll_runs(id) on delete cascade,
  employee_id text not null references employees(id), component_code text e.g., COMMISSION BONUS,
  amount numeric(20,8) not null, is_taxable boolean default true, description text, created_by text
);

-- Loans & Advances
create table payroll_loans (
  id text primary key, merchant_id text, employee_id text not null references employees(id),
  loan_type text check (type in ('personal','salary_advance','housing','other')), principal numeric(20,8),
  interest_rate numeric(5,2) default 0, tenure_months int, emi_amount numeric(20,8),
  total_paid numeric(20,8) default 0, outstanding numeric(20,8), status pending/approved/active/closed/rejected,
  disbursed_at timestamptz, created_at timestamptz default now()
);
create table payroll_loan_repayments (
  id text primary key, loan_id text references payroll_loans(id), run_id text references payroll_runs(id),
  amount numeric(20,8), principal_component numeric(20,8), interest_component numeric(20,8),
  created_at timestamptz default now()
);

-- Reimbursements & Claims enhanced
-- Already payroll_claims exists, enhance with approval flow chain, expense category

-- Payroll Run enhanced to include cost center allocation
alter table payroll_runs add column pay_calendar_id text;
alter table payroll_runs add column cutoff_date date;
alter table payroll_runs add column disbursal_date date;
alter table payroll_runs add column payroll_monthly_data jsonb default '{}' -- {total_paid_days, total_lop, etc}
alter table payroll_runs add column variance_report jsonb -- {vs_last_month_percent: 5.2%}

-- Enhanced Items with breakdown JSON for payslip
alter table payroll_items add column earnings_breakdown jsonb default '[]' -- [{code:BASIC, name:Basic, amount:20000, taxable:true}, ...]
alter table payroll_items add column deductions_breakdown jsonb default '[]'
alter table payroll_items add column employer_contributions_breakdown jsonb default '[]' -- pension 11%
alter table payroll_items add column ytd jsonb default '{}' -- {ytd_gross, ytd_tax, ytd_net}
alter table payroll_items add column paid_days int default 30;
alter table payroll_items add column lop_days int default 0;
alter table payroll_items add column proration_factor numeric(5,4) default 1.0;

-- Compliance Reports generated
create table payroll_compliance_reports (
  id text primary key, merchant_id text not null, period_month int, period_year int,
  report_type text check (type in ('pension_contribution','erca_withholding','annual_tax_certificate','pension_challan','bank_disbursal_file')),
  file_key text -- MinIO presigned 15m
  file_hash text,
  status generated/paid/filed,
  metadata jsonb,
  created_at timestamptz default now()
);

-- F&F Settlement
create table payroll_final_settlements (
  id text primary key, merchant_id text, employee_id text not null references employees(id),
  resignation_date date, last_working_date date, notice_period_days int, notice_served_days int,
  leave_encashment_days numeric(10,2), leave_encashment_amount numeric(20,8),
  severance_amount numeric(20,8), gratuity_amount numeric(20,8), -- per ET labour law
  outstanding_loans numeric(20,8), outstanding_advances numeric(20,8),
  total_payable numeric(20,8), total_deductions numeric(20,8), net_payable numeric(20,8),
  status draft/pending_approval/approved/paid,
  clearance_checklist jsonb default '[]', -- [{item:laptop, status:done}]
  created_at timestamptz default now()
);

-- Employee Portal Magic Links
create table payroll_employee_portal_access (
  id text primary key, merchant_id text, employee_id text references employees(id),
  magic_token_hash text, expires_at timestamptz, last_accessed_at timestamptz,
  created_at timestamptz default now()
);
```

### 4.2 Service Layer — Formula Engine + Proration + Arrears + Auto Compliance

**Formula Engine:** 
- Parse formula `CTC*0.4` or `BASIC*0.1`. Implement simple expression evaluator O(n) token map `map[string]decimal` {"CTC": ctc, "BASIC": basic}. No evil `eval`. Use `github.com/Knetic/govaluate` with decimal override or custom parser for security.

**CalculatePayroll V2 Algorithm:**
```
for each active employee:
  structure = GetStructure(employee.salary_structure_id)
  components = ListComponents(structure_id) sorted order_no
  earnings = {}
  for comp in components where type=earning:
    if calculation_type=fixed => amount
    if percentage_of_basic => basic * percentage
    if percentage_of_ctc => ctc_monthly * percentage
    if formula => evaluate(formula, {BASIC, CTC_MONTHLY})
  Apply proration:
    attendance = GetAttendance(run_id, employee_id)
    paid_days = attendance.paid_days
    total_days = attendance.total_days
    proration_factor = paid_days / total_days
    for each earning where prorate=true: earning.amount *= proration_factor
    gross = sum(earnings where is_part_of_gross)
  Add variable inputs:
    varInputs = ListVariable(run_id, employee_id) -> add to gross if taxable else separate
  Add OT:
    ot_amount = weekday_hours*basic/208*1.25 + weekend*1.5 + holiday*2.0 + night*1.3
    gross += ot_amount
  Add allowances from variable inputs
  Calculate pensionEmp 7% * (gross - non_pensionable?) per law gross includes all? For ET pension base is gross - OT? Need rule — make configurable pension_applicable_gross = gross - (OT + Bonus non-pensionable)
  taxable = gross - pensionEmp - tax_exempt_allowances (e.g., medical up to 1000 non-taxable) - reimbursements non-taxable
  incomeTax = CalculateTax binary search O(log n) taxable
  Deductions: incomeTax + pensionEmp + loanEMIs + advances + custom deductions + LOP recovery
  Net = gross - deductions + reimbursements taxable? Actually reimbursement non-taxable added after tax
  Break breakdown JSON earnings_breakdown, deductions_breakdown, employer_contributions [pension 11%], ytd calc cumulative sum year to date
  Total aggregates
```

**Arrears:** When salary revision effective_from in past, calculate arrear = (new_base - old_base) * months_pending + impact on pension/tax diff. Store as variable input type ARREAR taxable.

**Loan EMI Auto:** ListLoans active employee outstanding>0 -> emi_amount deduction per run until outstanding 0, create loan_repayment record linked to run_id.

**Compliance Automation:**
- After run completed, generate reports:
  - Pension CSV `pension_no, employee_name, basic? gross? pensionable salary, employee_contribution 7%, employer_contribution 11%, total, month` + file_key MinIO + hash
  - ERCA withholding monthly CSV `TIN, employee_name, period, gross, pension, taxable_income, income_tax, net, cost_center` - per ERCA format, include employee TIN 10-digit validation
  - Annual tax certificate per employee `GenerateAnnualTaxCertificate` PDF YTD gross taxable tax pension net
  - Bank disbursal file ISO20022 pain.001 XML per bank_code CBE/Awash/Dashen with <PmtInf> batch booking, employee bank hash masked but actual account encrypted retrieval via crypto.Decrypt AES-GCM using CONNECTOR_ENCRYPTION_KEY, file hash, batch ref

**Disbursal Atomic:** `DisburseRun` Tx:
```
BEGIN;
  SELECT run FOR UPDATE;
  IF run total_net > merchant balance => error insufficient_balance
  Create ledger_books payroll_run if not exists (idempotency ON CONFLICT DO NOTHING)
  Journal M4 Dr expense:salary totalGross Cr payroll_payable totalNet Cr et_income_tax totalTax Cr pension_payable totalPension ValidateBalanced
  Insert ledger_entries
  Create payout_batch book_id = payroll_run book? Actually separate payout_batch book per batch, but per spec per run book + payout batch second journal
  For each item net_pay>0: create payout beneficiary employee bank method=bank_ips status queued
  Second journal Dr payroll_payable totalNet Cr clearing:bank totalNet book payout_batch book
  Update payroll_runs status processing disbursal_date now
  Insert outbox payroll.disbursed
COMMIT;
  Async worker: process payout batch -> connector Bank IPS Initialize Verify success -> update payout_items paid + update payroll_items status paid
```

### 4.3 Employee Lifecycle & F&F

- **Onboarding:** Employee creation wizard 4 steps: Basic Info (code, name EN/AM, email, phone, DOB, gender, nationality ET, Fayda front/back <2MB selfie OTP verification same as merchant owner flow re-used), Employment (DOJ, probation end, department, designation, grade, branch, reporting manager, employment type), Salary (CTC, structure selector, breakdown live preview calculation, bank account name must match employee name fuzzy Levenshtein <3 + bank letter), Documents (contract, TIN certificate, Fayda, bank letter). Fayda verification via `POST /v1/fayda/verify/init` with employee_id.

- **Probation → Confirmation:** Scheduled job daily check probation_end_date <= today and status probation -> send notification to manager for confirmation.

- **Resignation & F&F:** Employee submits resignation date LWD. HR calculates leave encashment Annual balance * per_day `per_day = gross/30` per ET standard. Severance per Labour Proclamation No.1156/2019 Art 39-44: if unlawful termination, severance = 30 days wage * years of service first 3 years + later? Or mutual separation gratuity? Need configurable. Outstanding loans/advances deducted. Net payable. Approval dual finance+hr. Disburse single payout via off-cycle payroll run type bonus/fnf.

### 4.4 Self-Service Portal (Beyond RazorpayX)

- **Auth:** Magic link JWT 24h HMAC SHA256 signed `employee_id+merchant_id+expiry` emailed to employee email + WhatsApp integration optional (share_plus). No password needed, similar to RazorpayX WhatsApp payslips.

- **Dashboard:** Employee sees YTD gross/tax/net, upcoming payroll date, leave balance annual/sick, loan outstanding.

- **Payslips:** List month-wise, view outstanding modern template glassmorphic, download PDF password protected `DOB DDMMYYYY + last4` or `EMP001`, email, WhatsApp share, QR verification link.

- **Claims:** Create reimbursement expense/medical/travel amount description receipt upload MinIO presigned POST 15m file_key <5MB pdf/jpg/png, status pending approved rejected paid, timeline.

- **Loans:** Request loan type personal/salary_advance amount reason, EMI preview, status.

- **Bank Update:** Request bank account change, needs approval, bank account name must match employee name.

- **Documents:** View contract, payslips, tax certificates.

### 4.5 UI/UX Outstanding — Mercury/Linear Inspiration + Ethiopian

- **Tokens:** ET Green #0B6E4F, gold #EAB308, neutral zinc, radius xl 24, shadow soft, motion ease [0.22,1,0.36,1]
- **Payroll Runs Page:** Top 3 KPIs Gross/Tax/Net with sparkline Recharts, tabs Employees/Runs/Settings/Reports. Employees table avatar Fayda badge verified face_score 0.92, bank masked CBE ****, cost_center chips, department chips, grade, status. Runs table pipeline visual stepper draft→calculating→pending_approval→approved→processing→completed + progress bar % employees calculated. Run detail sticky footer totals, row expand breakdown chart pie deductions, approval flow dual avatar finance+admin confetti Lottie, payslip drawer glassmorphic.
- **Create Run Wizard:** Step1 Period Month/Year Type regular/off_cycle/bonus/adjustment Pay calendar cutoff date disbursal date Step2 Select employees active only + filters department branch cost_center + search + select all 500 employees <2s p99, Step3 Variable inputs CSV upload papaparse preview validation icons + OT hours per employee, Step4 Review totals gross/net/tax/pension variance vs last month % badge red/green, Step5 Disburse.
- **Salary Structure Builder:** Drag-drop components earnings/deductions/employer contributions, formula editor live preview CTC 500k annual breakdown monthly, validation basic 40-50% of CTC, order_no display in payslip.
- **Payslip PDF:** Modern template logo company, QR verification, earnings table, deductions table, employer contributions 11% pension, YTD box, net pay large font 24 bold ET Green, footer note "This is computer generated payslip no signature required, verified via QR", bilingual Amharic.
- **Reports:** Compliance center Perplexity-like citations? Actually payroll reports compliance pension ERCA with download buttons, Recharts cost center allocation pie.

### 4.6 AI Differentiation vs RazorpayX (ApexPay Moat)

- **AI Payroll Assist via Swarm:** Goal "Run payroll for July, add 10k bonus for Sales, deduct 500 penalty for EMP003" -> swarm planner decomposes: 1. get_tpv 2. list_employees cost_center Sales 3. create variable inputs bonus 10k 4. calculate_payroll_draft 5. needs_confirmation modal breakdown >100k confirmation true. Much beyond RazorpayX manual.

- **RAG Compliance:** Ask "What is ET pension rate for private employees?" -> answer "Employee 7% employer 11% per Private Organization Employees Pension Proclamation No.1268/2022 [1] score 0.92" with citations mandatory + no hallucination guard 0.65, Amharic/English. Also labour law "When is overtime rate 2.0?" -> "Holiday OT 2.0 per Labour Proclamation 1156/2019 Art 90"

- **Anomaly Detection:** AI critic checks payroll variance >20% vs last month flag, PEP count, dual approval risk, ghost employee detection (same bank account hash duplicate).

- **Payroll Chat:** In merchant-web command center chat glassmorphic "Show me payroll cost for Sales cost center last 3 months" -> tool get_payroll_cost_summaries Recharts.

---

## 5. API Enhancements Needed

```
POST /v1/salary_structures — Create CTC template with components
GET /v1/salary_structures/{id}
POST /v1/employees — enhanced payload department_id designation_id grade_id branch_id salary_structure_id ctc_annual bank_name bank_account_name reporting_manager_id probation_end_date documents[]
POST /v1/employees/bulk — CSV import 500 employees <2s
POST /v1/employees/{id}/fayda/verify/init|confirm
POST /v1/salary_revisions — Create revision effective_from future/past arrear auto calc
POST /v1/payroll_runs — enhanced pay_calendar_id cutoff_date disbursal_date
POST /v1/payroll_runs/{id}/attendance/bulk — CSV paid_days lop_days ot_hours
POST /v1/payroll_runs/{id}/variable_inputs/bulk — CSV component_code amount
POST /v1/payroll_runs/{id}/calculate — v2 formula engine proration arrears loan EMI
POST /v1/payroll_runs/{id}/hold/{employee_id} — Hold salary
POST /v1/payroll_runs/{id}/approve
POST /v1/payroll_runs/{id}/disburse — atomic M4 + payout batch + bank file generation pain.001
GET /v1/payroll_runs/{id}/payslips/{employee_id}/pdf — modern template QR + YTD
GET /v1/payroll_runs/{id}/payslips/bulk/zip — download all ZIP
GET /v1/payroll_reports/pension?month=7&year=2026 — CSV pension contribution
GET /v1/payroll_reports/erca_withholding?month=7&year=2026 — CSV ERCA
GET /v1/payroll_reports/cost_center?month=7&year=2026 — breakdown per cost_center
POST /v1/loans — request loan
POST /v1/claims — reimbursement
POST /v1/final_settlements — F&F
POST /v1/employee_portal/magic_link — generate JWT 24h for employee email
GET /v1/employee_portal/me — employee self-service auth via Bearer employee JWT
```

---

## 6. Implementation Roadmap — 4 Weeks to RazorpayX-Grade + Beyond

**Week 1 — Foundation Data Model + Salary Structure Engine**
- Migration 0013 payroll comprehensive tables listed above
- Seed departments, designations, grades, branches 14? Actually per merchant custom
- Salary structure CRUD handler service validator calculation_type fixed/percentage/formula evaluator secure no eval
- Employee extended fields + bulk CSV papaparse + Fayda verification reuse
- Salary revision history + arrear calc

**Week 2 — Payroll Engine V2 + Attendance + Variable Inputs + Loans**
- Refactor CalculateRun to V2 formula engine + proration (paid_days/total_days) + OT calc + variable inputs + loan EMI auto
- Attendance bulk import CSV + LOP logic
- Variable inputs bulk
- Loans and claims tables + approval flow
- Bank file generation pain.001 XML + MT103 CSV per bank config CBE/Awash etc

**Week 3 — Compliance + Disbursal Atomic + Reports + F&F**
- DisburseRun atomic Tx ledger M4 + payout batch second journal + merchant balance check
- Payout batch process worker Bank IPS connector real mock success via Telebirr/CBE/Bank
- Compliance reports pension CSV ERCA withholding CSV annual tax certificate PDF YTD
- F&F settlement service leave encashment per_day gross/30 + severance formula configurable + clearance checklist
- Payslip PDF v2 outstanding modern template QR verification JWT + YTD + bilingual + email distribution via webhook?

**Week 4 — Self-Service Portal + UI Gold + AI Assist**
- Employee portal magic link JWT auth + dashboard + payslips list + claims create + loan request + docs view
- Merchant-web UI: salary structure builder drag-drop, payroll run wizard 5 steps, employees table Fayda badge, run detail 8 cols sticky footer totals, payslip drawer glassmorphic, reports page Recharts cost center pie, compliance center download buttons
- Mobile Flutter: employee portal screens? Or merchant approvals payroll run approve biometric confetti
- Swarm payroll assist tool calculate_payroll_draft now uses real formula engine + confirmation >100k modal outstanding
- RAG ingests Ethiopian labour law Proclamation 1156/2019 + pension proclamation + ERCA directive for citations
- k6 bench payroll calc 500 employees <2s p99

---

## 7. My Opinion Before Coding

**You are 60% there for payroll core calc, 20% for comprehensive HR.** To be like RazorpayX you need salary structure template + attendance LOP + variable bulk + reimbursements/loans + self-service + compliance filing + F&F + reports. This is not trivial — RazorpayX took years.

**But you have advantage:** Ledger per run book + payout batch + Ethiopian ET tax binary search + Fayda + cost_center allocation is already better architecturally than RazorpayX India-only. If we build formula engine + proration + bank file + compliance CSV + self-service portal, you will surpass RazorpayX for Ethiopia because no one has ET pension 7/11 automated + ERCA CSV + Fayda verified payroll.

**Recommendation:** Don't just clone RazorpayX. Build **RazorpayX + Gusto + Remote.com for Ethiopia** with AI moat. 4 weeks focused will make payroll outstanding and powerful.

**What to prioritize first?**
- Salary structure builder (most critical gap)
- Attendance + proration + variable bulk
- Disbursal atomic + bank file + compliance reports pension/ERCA
- Self-service portal magic link
- F&F settlement

If you approve, I will start with **Migration 0013 + Salary Structure Engine + Employee Extended** then Payroll V2 engine.

Should I start coding this full payroll OS now?

