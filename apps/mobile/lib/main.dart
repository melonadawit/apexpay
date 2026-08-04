import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'src/core/theme/app_theme.dart';
import 'src/core/router/app_router.dart';
import 'src/core/storage/hive_boxes.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Hive for offline draft links - optimal local cache
  await Hive.initFlutter();
  await HiveBoxes.init();
  
  // Firebase would be init here for FCM
  // await Firebase.initializeApp();

  runApp(const ProviderScope(child: ApexPayApp()));
}

class ApexPayApp extends ConsumerWidget {
  const ApexPayApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    return MaterialApp.router(
      title: 'ApexPay Merchant',
      theme: AppTheme.light(),
      darkTheme: AppTheme.dark(),
      themeMode: ThemeMode.system,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
