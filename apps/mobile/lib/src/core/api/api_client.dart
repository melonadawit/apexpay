import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

// Base URL is configurable at build time so the app can target a local/dev/prod API:
//   flutter run --dart-define=APEXPAY_API_URL=http://10.0.2.2:8080/v1
const String _defaultBase = 'https://api.apexpay.et/v1';
const String apiBaseUrl = String.fromEnvironment('APEXPAY_API_URL', defaultValue: _defaultBase);

class ApiClient {
  final Dio dio = Dio(BaseOptions(baseUrl: apiBaseUrl, connectTimeout: const Duration(seconds: 10)));
  final _storage = const FlutterSecureStorage();

  ApiClient() {
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (opts, handler) async {
        // Prefer the dashboard session token if present, else the merchant API key.
        final token = await _storage.read(key: 'apexpay_session') ??
            await _storage.read(key: 'sk_test');
        if (token != null) opts.headers['Authorization'] = 'Bearer $token';
        opts.headers['X-Request-Id'] = DateTime.now().millisecondsSinceEpoch.toString();
        return handler.next(opts);
      },
      onError: (e, handler) {
        // Log but don't expose FIN/account in logs
        return handler.next(e);
      },
    ));
  }

  Future<Map<String, dynamic>> post(String path, Map<String, dynamic> data) async {
    final res = await dio.post(path, data: data);
    return res.data as Map<String, dynamic>;
  }

  Future<Map<String, dynamic>> get(String path) async {
    final res = await dio.get(path);
    return res.data as Map<String, dynamic>;
  }

  // Stores a session/API token securely for subsequent requests.
  Future<void> setAuthToken(String token, {String key = 'apexpay_session'}) async {
    await _storage.write(key: key, value: token);
  }

  Future<void> clearAuth() async {
    await _storage.delete(key: 'apexpay_session');
    await _storage.delete(key: 'sk_test');
  }
}

