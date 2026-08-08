import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'src/core/theme/app_theme.dart';
import 'src/core/router/app_router.dart';
import 'src/core/storage/hive_boxes.dart';
import 'src/core/api/fcm_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Hive for offline draft links - optimal local cache
  await Hive.initFlutter();
  await HiveBoxes.init();

  // Firebase + FCM for push. Firebase.initializeApp() is commented out because no
  // google-services.json is committed; uncomment once Firebase project config is added.
  // Register the top-level background handler BEFORE runApp.
  FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);

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
