import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  final Dio dio = Dio(BaseOptions(baseUrl: 'https://api.apexpay.et/v1', connectTimeout: const Duration(seconds: 10)));
  final _storage = const FlutterSecureStorage();

  ApiClient() {
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (opts, handler) async {
        final token = await _storage.read(key: 'sk_test');
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

  Future<Map<String,dynamic>> post(String path, Map<String,dynamic> data) async {
    final res = await dio.post(path, data: data);
    return res.data as Map<String,dynamic>;
  }

  Future<Map<String,dynamic>> get(String path) async {
    final res = await dio.get(path);
    return res.data as Map<String,dynamic>;
  }
}
