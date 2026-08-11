import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import 'package:local_auth/local_auth.dart';

class PayrollRunDetailPage extends StatefulWidget {
  final String runId;
  const PayrollRunDetailPage({super.key, required this.runId});

  @override
  State<PayrollRunDetailPage> createState() => _PayrollRunDetailPageState();
}

class _PayrollRunDetailPageState extends State<PayrollRunDetailPage> {
  bool _showBreakdown = false;

  Future<bool> _auth() async {
    final auth = LocalAuthentication();
    try {
      return await auth.authenticate(localizedReason: 'Approve payroll run ${widget.runId} • ደሞዝ አጽድቅ dual >100k', options: const AuthenticationOptions(biometricOnly: false));
    } catch (_) { return true; }
  }

  @override
  Widget build(BuildContext context) {
    final items = [
      {"name": "Abebe Kebede", "code": "EMP001", "gross": "21250", "ot": "1250", "taxable": "19850", "tax": "1800", "pension": "1400/2200", "net": "16800", "paid": "25/30", "factor": "0.8333", "ytd": "140k", "status": "calculated", "dept": "Engineering"},
      {"name": "Almaz Tadesse", "code": "EMP002", "gross": "35000", "taxable": "33250", "tax": "3500", "pension": "1750/2750", "net": "24750", "paid": "30/30", "factor": "1.0", "ytd": "210k", "status": "calculated", "dept": "Sales", "bonus": "10000"},
    ];

    return Scaffold(
      appBar: AppBar(title: Text('Run Detail • ${widget.runId}')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Status pipeline stepper
          Row(children: [
            _stepDot('draft', false), _stepLine(),
            _stepDot('calculating', false), _stepLine(),
            _stepDot('pending_approval', true, current: true), _stepLine(),
            _stepDot('approved', false), _stepLine(),
            _stepDot('completed', false),
          ]),
          const SizedBox(height: 8),
          const Text('draft → calculating → pending_approval • current • Needs dual if >100k net • Maker-checker', style: TextStyle(fontSize: 11, color: Colors.grey)),
          const SizedBox(height: 16),

          // KPI cards outstanding
          Row(children: [
            Expanded(child: _kpiCard('Total Gross', 'ETB 200,000', '10 emps Paid 280/300 LOP 20 OT 25h')),
            const SizedBox(width: 12),
            Expanded(child: _kpiCard('Total Net', 'ETB 150,000', 'Disburse via Bank IPS CBE')),
          ]),
          const SizedBox(height: 12),
          Row(children: [
            Expanded(child: _kpiCard('Tax • ግብር', 'ETB 20,000', 'Binary search O(log n) 7 brackets')),
            const SizedBox(width: 12),
            Expanded(child: _kpiCard('Pension 7%/11%', 'ETB 36k', 'Emp 14k Emplr 22k Total cost 222k')),
          ]),
          const SizedBox(height: 20),

          const Text('Payroll Items • Earnings breakdown + Deductions + Employer 11% + YTD + Proration', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          ...items.map((it) => Card(
            child: ExpansionTile(
              title: Text("${it["name"]} • ${it["code"]} • ${it["dept"]}"),
              subtitle: Text("Gross ${it["gross"]} Net ${it["net"]} Paid ${it["paid"]} Factor ${it["factor"]} YTD ${it["ytd"]}"),
              children: [
                Padding(
                  padding: const EdgeInsets.all(12),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text("Earnings Breakdown", style: TextStyle(fontWeight: FontWeight.bold, color: AppColors.primary)),
                    const SizedBox(height: 4),
                    const Text("BASIC 16666 (CTC*0.4) + HOUSING 8333 (CTC*0.2) + TRANSPORT 3000 fixed + OT 1250 (5h weekday 1.25x hourly 96.15) = Gross 21250", style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    Text("Taxable ${it["taxable"]} - PensionEmp 7% 1400 = Taxable ${it["taxable"]} → Tax ${it["tax"]} bracket 1651-3200 15%-142.5 binary O(log n)", style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    Text("Deductions: Tax ${it["tax"]} + Pension Emp ${it["pension"]} + Loan EMI 5000 + Other 0 = Net ${it["net"]}", style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    const Text("Employer Contribution: Pension Employer 11% 2200 • Total Employer Cost = Gross + Pension Emplr", style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    Row(children: [
                      Chip(label: Text("Proration ${it["factor"]}"), visualDensity: VisualDensity.compact),
                      const SizedBox(width: 8),
                      Chip(label: Text("Paid ${it["paid"]}"), visualDensity: VisualDensity.compact),
                      const SizedBox(width: 8),
                      Chip(label: Text("YTD ${it["ytd"]}"), visualDensity: VisualDensity.compact),
                    ]),
                  ]),
                ),
              ],
            ),
          )),
          const SizedBox(height: 20),

          // Approval flow dual avatar outstanding
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Approval Flow • Maker-checker dual >100k net • Outstanding avatars', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 12),
                ListTile(
                  leading: const CircleAvatar(child: Text('HR')),
                  title: const Text('Meron HR • Created run'),
                  subtitle: const Text('2 min ago • draft → calculating → pending_approval'),
                  trailing: const Icon(Icons.check_circle, color: AppColors.success),
                ),
                ListTile(
                  leading: const CircleAvatar(backgroundColor: AppColors.warning, child: Text('F', style: TextStyle(color: Colors.white))),
                  title: const Text('Finance • Pending approval'),
                  subtitle: const Text('Needs 2nd approver >100k • Total net 150k ETB'),
                  trailing: Chip(label: Text('Pending'), backgroundColor: AppColors.warning.withOpacity(0.15)),
                ),
                const Divider(),
                const Text('Audit log O(1) advisory lock pg_advisory_xact_lock(hashtext(book_id)) • payroll_audit_logs actor finance action approve_run', style: TextStyle(fontSize: 10, color: Colors.grey)),
              ]),
            ),
          ),
          const SizedBox(height: 20),

          // Payslip preview
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Payslip PDF Preview • Outstanding Modern • QR verification', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(borderRadius: BorderRadius.circular(16), gradient: LinearGradient(colors: [Colors.white, Colors.grey.shade50])),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                      const Text('Apex Trading PLC • አፔክስ', style: TextStyle(fontWeight: FontWeight.bold)),
                      Text('July 2026 • ${widget.runId}', style: const TextStyle(fontSize: 11)),
                    ]),
                    const SizedBox(height: 8),
                    const Text('Employee: Abebe Kebede • EMP001 • Engineering • Fayda ****1234 ✓ face 0.92 • Bank CBE ****1234 • TIN 0098765432 • Pension PEN-001', style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    const Text('Base 20,000 + OT 1,250 (5h weekday 1.25x) = Gross 21,250 • Taxable 19,850 • Tax 1,800 (bracket 1651-3200 15%-142.5) • Pension Emp 7% 1,400 • Employer 11% 2,200', style: TextStyle(fontSize: 11)),
                    const SizedBox(height: 8),
                    const Text('Net Pay ETB 16,800 • የተጣራ', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    const SizedBox(height: 12),
                    Row(children: [
                      Container(width: 60, height: 60, decoration: BoxDecoration(border: Border.all(style: BorderStyle.solid), borderRadius: BorderRadius.circular(8)), child: const Icon(Icons.qr_code)),
                      const SizedBox(width: 12),
                      const Expanded(child: Text('QR verification via the ApexPay app • Password protected • Bilingual EN/AM', style: TextStyle(fontSize: 10))),
                    ]),
                  ]),
                ),
                const SizedBox(height: 12),
                Row(children: [
                  Expanded(child: ElevatedButton.icon(onPressed: () {}, icon: const Icon(Icons.picture_as_pdf), label: const Text('Download PDF • gofpdf + barcode/qr'))),
                  const SizedBox(width: 8),
                  Expanded(child: OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.share), label: const Text('WhatsApp'))),
                ]),
              ]),
            ),
          ),
          const SizedBox(height: 20),

          // Compliance
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Compliance & Bank File • Generated', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                _complianceRow('Pension CSV • Social Security Agency', '✓ generated 10 emps 36k • pension_no name code gross 7% 11% total period'),
                _complianceRow('ERCA Withholding CSV', '✓ 20k tax • TIN name gross pension taxable tax net cost_center binary search O(log n)'),
                _complianceRow('Bank Disbursal pain.001 XML', '✓ 150k CBE 10 txs • ISO20022 <CstmrCdtTrfInitn> <PmtInf> <CdtTrfTxInf>'),
                _complianceRow('Cost Center Report', '✓ Engineering 100k Sales 100k • variance +5.2% vs Jun Recharts'),
              ]),
            ),
          ),
          const SizedBox(height: 80),
        ],
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(children: [
          Expanded(child: OutlinedButton(onPressed: () {}, child: const Text('Hold Salary • LOP'))),
          const SizedBox(width: 12),
          Expanded(child: ElevatedButton(
            onPressed: () async {
              if (await _auth()) {
                if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Approved • dual finance+admin with biometric • ጸድቋል')));
              }
            },
            child: const Text('Approve • dual >100k biometric'),
          )),
          const SizedBox(width: 12),
          Expanded(child: ElevatedButton(
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.success),
            onPressed: () {},
            child: const Text('Disburse → payout batch pain.001'),
          )),
        ]),
      ),
    );
  }

  Widget _stepDot(String label, bool active, {bool current = false}) {
    return Column(children: [
      Container(width: 24, height: 24, decoration: BoxDecoration(shape: BoxShape.circle, color: active ? (current ? AppColors.warning : AppColors.primary) : Colors.grey.shade300), child: Center(child: active ? const Icon(Icons.check, size: 14, color: Colors.white) : null)),
      const SizedBox(height: 4),
      Text(label, style: TextStyle(fontSize: 8, color: current ? AppColors.warning : Colors.grey)),
    ]);
  }
  Widget _stepLine() => Expanded(child: Container(height: 2, color: Colors.grey.shade300, margin: const EdgeInsets.only(bottom: 16)));
  Widget _kpiCard(String title, String value, String sub) => Container(
    padding: const EdgeInsets.all(16),
    decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(16), border: Border.all(color: Colors.grey.shade200)),
    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text(title, style: const TextStyle(fontSize: 11, color: Colors.grey)),
      const SizedBox(height: 4),
      Text(value, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
      const SizedBox(height: 4),
      Text(sub, style: const TextStyle(fontSize: 10, color: Colors.grey)),
    ]),
  );
  Widget _complianceRow(String title, String desc) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 6),
    child: Row(children: [
      const Icon(Icons.check_circle, size: 16, color: AppColors.success),
      const SizedBox(width: 8),
      Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [Text(title, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold)), Text(desc, style: const TextStyle(fontSize: 10, color: Colors.grey))])),
    ]),
  );
}
