// ApexPay merchant app smoke test.
//
// Verifies the design system brand colors and the Material 3 color scheme
// derived from the ApexPay brand seed. Kept google-font/Hive/Firebase-free so
// it runs headless in CI without network or platform plugins.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:apexpay_merchant/src/core/theme/app_theme.dart';

void main() {
  test('ApexPay brand palette matches the ET green + Abyssinia gold identity', () {
    expect(AppColors.primary, const Color(0xFF0B6E4F));
    expect(AppColors.accentGold, const Color(0xFFEAB308));
    expect(AppColors.success, const Color(0xFF10B981));
    expect(AppColors.error, const Color(0xFFEF4444));
  });

  test('light color scheme derives from the brand seed', () {
    final scheme = ColorScheme.fromSeed(seedColor: AppColors.primary, brightness: Brightness.light);
    expect(scheme.primary, isNot(Colors.black));
    expect(scheme.brightness, Brightness.light);
  });

  testWidgets('ApexPay brand primary color renders a MaterialApp', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        theme: ThemeData(colorScheme: ColorScheme.fromSeed(seedColor: AppColors.primary)),
        home: const Scaffold(body: Center(child: Text('ApexPay'))),
      ),
    );
    expect(find.text('ApexPay'), findsOneWidget);
  });
}
