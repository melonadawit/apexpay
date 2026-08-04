import 'package:firebase_messaging/firebase_messaging.dart';
import './api_client.dart';

/// FCM service — push_devices table FCM token unique, platform android/ios/web per DATABASE
/// Best practice: token refresh listener, background handler, foreground notification

class FCMService {
  final ApiClient apiClient;
  final FirebaseMessaging _fcm = FirebaseMessaging.instance;

  FCMService(this.apiClient);

  Future<void> init() async {
    // Request permission iOS
    await _fcm.requestPermission(alert: true, badge: true, sound: true);

    // Get token
    final token = await _fcm.getToken();
    if (token != null) {
      await _registerToken(token, 'android'); // platform detect via Platform
    }

    // Token refresh
    _fcm.onTokenRefresh.listen((newToken) async {
      await _registerToken(newToken, 'android');
    });

    // Foreground message
    FirebaseMessaging.onMessage.listen((RemoteMessage msg) {
      // Show local notification mock for outstanding UI
      print('FCM foreground: ${msg.notification?.title} • ${msg.data}');
      // In real, use flutter_local_notifications to show
    });

    // Background handler already registered in main.dart via FirebaseMessaging.onBackgroundMessage
  }

  Future<void> _registerToken(String token, String platform) async {
    try {
      await apiClient.post('/devices/register', {
        'fcm_token': token,
        'platform': platform,
        'device_info': {'model': 'Pixel 7', 'os': 'Android 14'},
      });
      print('FCM token registered: ${token.substring(0,20)}... • push_devices unique');
    } catch (e) {
      print('FCM register failed: $e — will retry on next sync');
    }
  }

  Future<void> subscribeTopics() async {
    await _fcm.subscribeToTopic('payments_succeeded');
    await _fcm.subscribeToTopic('payouts_pending_approval');
    await _fcm.subscribeToTopic('payroll_runs_pending');
  }
}

// Background handler must be top-level per Firebase spec
@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  print('FCM background: ${message.data}');
}
