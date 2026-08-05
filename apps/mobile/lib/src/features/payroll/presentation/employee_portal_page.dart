import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';

class EmployeePortalPage extends StatelessWidget {
  const EmployeePortalPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Employee Portal • የሰራተኛ መግቢያ • Self-Service')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // YTD glass card
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: const LinearGradient(colors: [AppColors.primary, AppColors.primaryLight]),
              borderRadius: BorderRadius.circular(24),
            ),
            child: const Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('YTD • የዓመቱ አጠቃላይ • Abebe Kebede • EMP001', style: TextStyle(color: Colors.white70)),
              SizedBox(height: 8),
              Text('Gross ETB 140k • Tax 12k • Net 98k', style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
              SizedBox(height: 8),
              Text('July 2026 regular • Gross 21,250 Net 16,800 • Paid days 25/30 Factor 0.8333 • OT 5h weekday 1.25x • Pension Emp 7% 1,400 Emplr 11% 2,200', style: TextStyle(color: Colors.white70, fontSize: 11)),
            ]),
          ),
          const SizedBox(height: 20),

          const Text('Payslips • QR verified • Bilingual EN/AM', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          ...List.generate(3, (i) => Card(
            child: ListTile(
              leading: const CircleAvatar(child: Icon(Icons.picture_as_pdf)),
              title: Text('Payslip ${["July 2026", "June 2026", "May 2026"][i]} • ETB ${[16800, 14250, 14000][i]}'),
              subtitle: Text('Gross ${[21250, 19000, 18500][i]} Tax ${[1800, 1600, 1500][i]} • Pension 7%/11% • YTD Gross ${[140000, 118750, 99750][i]} • QR verified'),
              trailing: const Icon(Icons.qr_code),
              onTap: () {},
            ),
          )),
          const SizedBox(height: 20),

          const Text('Loans & Advances • EMI auto deduction', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          Card(
            child: ListTile(
              leading: const CircleAvatar(backgroundColor: AppColors.warning, child: Icon(Icons.payments, color: Colors.white)),
              title: const Text('Salary Advance ETB 20,000 • EMI 5,000 • Outstanding 15,000'),
              subtitle: const Text('Tenure 4mo • 1/4 paid • Next due Aug 2026 • Auto deduction per payroll run O(k)'),
              trailing: Chip(label: Text('Active'), backgroundColor: AppColors.success.withOpacity(0.15)),
            ),
          ),
          const SizedBox(height: 20),

          const Text('Claims & Reimbursements • Expense/Medical/Travel • Receipt MinIO <5MB', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          Card(
            child: ListTile(
              leading: const CircleAvatar(child: Icon(Icons.receipt)),
              title: const Text('Travel ETB 2,000 • travel_receipt.pdf'),
              subtitle: const Text('Status pending → approved → paid via next payroll reimbursement non-taxable • File_key MinIO presigned 15m hash'),
              trailing: const Icon(Icons.pending_actions, color: AppColors.warning),
            ),
          ),
          const SizedBox(height: 20),

          const Text('Documents • Contract TIN Fayda Bank Letter • Vault', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          Row(children: [
            Expanded(child: Card(child: Padding(padding: EdgeInsets.all(12), child: Column(children: [Icon(Icons.description), SizedBox(height: 4), Text('Contract', style: TextStyle(fontSize: 11)), Text('Verified ✓', style: TextStyle(fontSize: 10, color: AppColors.success))])))),
            Expanded(child: Card(child: Padding(padding: EdgeInsets.all(12), child: Column(children: [Icon(Icons.badge), SizedBox(height: 4), Text('Fayda Front', style: TextStyle(fontSize: 11)), Text('Verified 0.92 ✓', style: TextStyle(fontSize: 10, color: AppColors.success))])))),
            Expanded(child: Card(child: Padding(padding: EdgeInsets.all(12), child: Column(children: [Icon(Icons.account_balance), SizedBox(height: 4), Text('Bank Letter', style: TextStyle(fontSize: 11)), Text('CBE ****1234 ✓', style: TextStyle(fontSize: 10, color: AppColors.success))])))),
            Expanded(child: Card(child: Padding(padding: EdgeInsets.all(12), child: Column(children: [Icon(Icons.receipt_long), SizedBox(height: 4), Text('TIN Cert', style: TextStyle(fontSize: 11)), Text('0098765432 ✓', style: TextStyle(fontSize: 10, color: AppColors.success))])))),
          ]),
          const SizedBox(height: 20),

          Card(
            color: AppColors.primary.withOpacity(0.05),
            child: const Padding(
              padding: EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('How to verify payslip QR • Outstanding Modern', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12)),
                SizedBox(height: 8),
                Text('1. Open ApexPay merchant app → Scan QR → /qr/scan overlay rounded 260 corner brackets pulse green + vibration', style: TextStyle(fontSize: 11)),
                Text('2. QR contains runId + employeeCode + netPay hash signed JWT HMAC SHA256 + expiry 24h', style: TextStyle(fontSize: 11)),
                Text('3. Verify via https://apexpay.et/verify/payslip/{runId}/{employeeCode} → shows gross/tax/net breakdown + YTD + ledger M4 balanced ✓', style: TextStyle(fontSize: 11)),
                SizedBox(height: 8),
                Text('Password protected PDF: DOB DDMM + last4 • Bilingual EN/AM • Lottie confetti 3s + haptics • WhatsApp share', style: TextStyle(fontSize: 11, color: AppColors.primary)),
              ]),
            ),
          ),
          const SizedBox(height: 80),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        label: const Text('Request Loan • EMI preview • magic link 24h'),
        icon: const Icon(Icons.request_quote),
      ),
    );
  }
}
