import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_theme.dart';

class PayrollRunsPage extends StatelessWidget {
  const PayrollRunsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final runs = [
      {"ref": "July2026_Regular", "period": "07/2026", "type": "regular", "status": "pending_approval", "gross": "200000", "net": "150000", "count": 10, "variance": "+5.2%"},
      {"ref": "June2026_Regular", "period": "06/2026", "type": "regular", "status": "completed", "gross": "190000", "net": "142500", "count": 10, "variance": "-2.1%"},
    ];

    return Scaffold(
      appBar: AppBar(title: const Text('Payroll Runs • ደሞዝ ሩጫዎች')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // KPI card outstanding
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: const LinearGradient(colors: [AppColors.primary, AppColors.primaryLight]),
              borderRadius: BorderRadius.circular(24),
            ),
            child: const Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Total Net July • የተጣራ', style: TextStyle(color: Colors.white70)),
                SizedBox(height: 8),
                Text('ETB 150,000', style: TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.bold)),
                SizedBox(height: 8),
                Text('10 employees • Ledger M4 balanced • Bank pain.001 generated', style: TextStyle(color: Colors.white70, fontSize: 12)),
              ],
            ),
          ),
          const SizedBox(height: 20),
          const Text('Payroll Runs • Status pipeline visual stepper', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          ...runs.map((r) => Card(
            child: ListTile(
              leading: CircleAvatar(
                backgroundColor: r["status"]=="completed" ? AppColors.success.withValues(alpha: 0.15) : AppColors.warning.withValues(alpha: 0.15),
                child: Icon(r["status"]=="completed" ? Icons.check : Icons.pending_actions, color: r["status"]=="completed" ? AppColors.success : AppColors.warning),
              ),
              title: Text("${r["ref"]} • ${r["period"]}"),
              subtitle: Text("Gross ${r["gross"]} Net ${r["net"]} • ${r["count"]} emps • Variance ${r["variance"]} • Type ${r["type"]}"),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => context.push('/payroll/${r["ref"]}'),
            ),
          )),
          const SizedBox(height: 24),
          const Text('Quick Actions', style: TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          Row(children: [
            Expanded(child: ElevatedButton.icon(onPressed: () {}, icon: const Icon(Icons.add), label: const Text('Create Run'))),
            const SizedBox(width: 12),
            Expanded(child: OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.upload_file), label: const Text('Import Attendance'))),
          ]),
          const SizedBox(height: 12),
          OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.picture_as_pdf), label: const Text('Pension CSV + ERCA CSV + Bank pain.001')),
          const SizedBox(height: 24),
          Card(
            color: AppColors.primary.withValues(alpha: 0.05),
            child: const Padding(
              padding: EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('Ledger M4 per run book outstanding', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12)),
                SizedBox(height: 8),
                Text('Dr expense:salary 200k + Dr expense:pension_employer 22k', style: TextStyle(fontSize: 11, fontFamily: 'monospace')),
                Text('Cr payroll_payable 150k Cr et_income_tax_payable 20k Cr pension_payable 36k balanced ValidateBalanced', style: TextStyle(fontSize: 11, fontFamily: 'monospace')),
                SizedBox(height: 8),
                Text('Second journal disburse: Dr payroll_payable 150k Cr clearing:bank 150k via payout batch pain.001 XML', style: TextStyle(fontSize: 11, fontFamily: 'monospace')),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.push('/payroll/create'),
        label: const Text('Create Run • 5 steps wizard'),
        icon: const Icon(Icons.play_arrow),
      ),
    );
  }
}
