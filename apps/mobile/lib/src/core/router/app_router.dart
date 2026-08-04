import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../features/auth/presentation/login_page.dart';
import '../../features/dashboard/presentation/dashboard_page.dart';
import '../../features/links/presentation/create_link_sheet.dart';
import '../../features/qr/presentation/qr_scanner_page.dart';
import '../../features/approvals/presentation/approvals_page.dart';
import '../../features/profile/onboarding/presentation/onboarding_wizard_page.dart';

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
        ],
      ),
      GoRoute(path: '/links/create', builder: (c,s) => const CreateLinkSheet()),
      GoRoute(path: '/qr/scan', builder: (c,s) => const QrScannerPage()),
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
        ],
        onDestinationSelected: (i) {
          if (i==0) context.go('/dashboard');
          else context.go('/approvals');
        },
      ),
    );
  }
}
