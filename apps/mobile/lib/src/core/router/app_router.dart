import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../features/auth/presentation/login_page.dart';
import '../../features/dashboard/presentation/dashboard_page.dart';
import '../../features/links/presentation/create_link_sheet.dart';
import '../../features/qr/presentation/qr_scanner_page.dart';
import '../../features/approvals/presentation/approvals_page.dart';
import '../../features/profile/onboarding/presentation/onboarding_wizard_page.dart';
import '../../features/payroll/presentation/payroll_runs_page.dart';
import '../../features/payroll/presentation/payroll_run_detail_page.dart';
import '../../features/payroll/presentation/employee_portal_page.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/login',
    routes: [
      GoRoute(path: '/login', builder: (c,s) => const LoginPage()),
      GoRoute(path: '/onboarding', builder: (c,s) => const OnboardingWizardPage()),
      ShellRoute(
        builder: (context, state, child) => ScaffoldWithNav(child: child),
        routes: [
          GoRoute(path: '/dashboard', builder: (c,s) => const DashboardPage()),
          GoRoute(path: '/approvals', builder: (c,s) => const ApprovalsPage()),
          GoRoute(path: '/payroll', builder: (c,s) => const PayrollRunsPage()),
          GoRoute(path: '/employee/portal', builder: (c,s) => const EmployeePortalPage()),
        ],
      ),
      GoRoute(path: '/links/create', builder: (c,s) => const CreateLinkSheet()),
      GoRoute(path: '/qr/scan', builder: (c,s) => const QrScannerPage()),
      GoRoute(path: '/payroll/:runId', builder: (c,s) => PayrollRunDetailPage(runId: s.pathParameters['runId']!)),
      GoRoute(path: '/employee/payroll/:runId', builder: (c,s) => PayrollRunDetailPage(runId: s.pathParameters['runId']!)),
    ],
  );
});

class ScaffoldWithNav extends StatelessWidget {
  final Widget child;
  const ScaffoldWithNav({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(icon: Icon(Icons.dashboard_outlined), label: 'Dashboard'),
          NavigationDestination(icon: Icon(Icons.check_circle_outline), label: 'Approvals'),
          NavigationDestination(icon: Icon(Icons.groups_outlined), label: 'Payroll • ደሞዝ'),
          NavigationDestination(icon: Icon(Icons.person_outline), label: 'Me • Portal'),
        ],
        onDestinationSelected: (i) {
          if (i==0) context.go('/dashboard');
          else if (i==1) context.go('/approvals');
          else if (i==2) context.go('/payroll');
          else context.go('/employee/portal');
        },
      ),
    );
  }
}
