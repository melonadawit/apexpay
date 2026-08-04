import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_theme.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  @override State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailCtrl = TextEditingController(text: 'meron@demo.et');
  final _passCtrl = TextEditingController(text: 'demo123');
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(gradient: LinearGradient(colors: [Color(0xFF0B6E4F), Color(0xFF10A37A)], begin: Alignment.topLeft, end: Alignment.bottomRight)),
        child: SafeArea(
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24),
              child: GlassCard(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const Icon(Icons.account_balance_wallet, size: 64, color: AppColors.primary),
                    const SizedBox(height: 16),
                    Text('ApexPay', style: Theme.of(context).textTheme.displayLarge, textAlign: TextAlign.center),
                    Text('Merchant • Ethiopia', style: Theme.of(context).textTheme.bodyLarge, textAlign: TextAlign.center),
                    const SizedBox(height: 32),
                    TextField(controller: _emailCtrl, decoration: const InputDecoration(labelText: 'Email • ኢሜይል', prefixIcon: Icon(Icons.email_outlined))),
                    const SizedBox(height: 16),
                    TextField(controller: _passCtrl, obscureText: true, decoration: const InputDecoration(labelText: 'Password • የይለፍ ቃል', prefixIcon: Icon(Icons.lock_outline))),
                    const SizedBox(height: 24),
                    ElevatedButton(
                      onPressed: _loading ? null : () async {
                        setState(()=> _loading=true);
                        await Future.delayed(const Duration(seconds: 1));
                        if (mounted) context.go('/dashboard');
                        setState(()=> _loading=false);
                      },
                      child: _loading ? const SizedBox(height:20,width:20,child: CircularProgressIndicator(color: Colors.white,strokeWidth:2)) : const Text('Login • ግባ'),
                    ),
                    const SizedBox(height: 12),
                    TextButton(onPressed: ()=> context.go('/onboarding'), child: const Text('Start onboarding • ምዝገባ ጀምር')),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
