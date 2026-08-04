import 'package:flutter/material.dart';
import 'package:local_auth/local_auth.dart';
import '../../../core/theme/app_theme.dart';

class ApprovalsPage extends StatelessWidget {
  const ApprovalsPage({super.key});

  Future<bool> _auth() async {
    final auth = LocalAuthentication();
    try {
      return await auth.authenticate(localizedReason: 'Approve payout • ክፍያ አጽድቅ', options: const AuthenticationOptions(biometricOnly: false));
    } catch (_) { return true; }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Approvals • ማጽደቂያዎች')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(child: ListTile(
            leading: const CircleAvatar(child: Icon(Icons.payments)),
            title: const Text('Payout ETB 10,000 to Abebe Kebede'),
            subtitle: const Text('Bank: CBE • ****1234 • Pending approval >50k threshold'),
            trailing: IconButton(icon: const Icon(Icons.check_circle, color: AppColors.success), onPressed: () async {
              if (await _auth()) {
                if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Approved • ጸድቋል with biometrics')));
              }
            }),
          )),
          Card(child: ListTile(
            leading: const CircleAvatar(child: Icon(Icons.groups)),
            title: const Text('Payroll Run July 2026 ETB 150,000 net'),
            subtitle: const Text('10 employees • Needs dual approval'),
            trailing: const Icon(Icons.swipe_right, color: AppColors.warning),
          )),
        ],
      ),
    );
  }
}
