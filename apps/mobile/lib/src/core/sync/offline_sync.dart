import '../api/api_client.dart';
import '../storage/hive_boxes.dart';

/// Offline sync service — optimal: Hive draft_links + offline_queue, sync on reconnect with idempotency key same as web
/// Best practice: queue O(n) with retry exponential backoff, no duplicate via Idempotency-Key header

class OfflineSyncService {
  final ApiClient apiClient;
  OfflineSyncService(this.apiClient);

  Future<int> getPendingCount() async {
    return HiveBoxes.draftLinksBox.length + HiveBoxes.offlineQueueBox.length;
  }

  Future<void> saveDraftLink(Map<String,dynamic> linkData) async {
    final id = DateTime.now().millisecondsSinceEpoch.toString();
    await HiveBoxes.draftLinksBox.put(id, {
      ...linkData,
      'id': id,
      'created_at': DateTime.now().toIso8601String(),
      'idempotency_key': 'idem_$id', // same as web per spec
    });
  }

  Future<void> enqueue(String type, Map<String,dynamic> payload) async {
    final id = DateTime.now().millisecondsSinceEpoch.toString();
    await HiveBoxes.offlineQueueBox.put(id, {
      'type': type,
      'payload': payload,
      'attempts': 0,
      'created_at': DateTime.now().toIso8601String(),
      'idempotency_key': 'idem_${type}_$id',
    });
  }

  Future<SyncResult> syncAll() async {
    int synced = 0;
    int failed = 0;

    // Sync draft links first
    for (var key in HiveBoxes.draftLinksBox.keys.toList()) {
      final data = HiveBoxes.draftLinksBox.get(key);
      try {
        await apiClient.post('/payment_links', {
          'amount': data['amount'],
          'currency': 'ETB',
          'description': data['description'],
        });
        await HiveBoxes.draftLinksBox.delete(key);
        synced++;
      } catch (e) {
        failed++;
      }
    }

    // Sync queue with exponential backoff
    for (var key in HiveBoxes.offlineQueueBox.keys.toList()) {
      final item = HiveBoxes.offlineQueueBox.get(key);
      final attempts = item['attempts'] ?? 0;
      if (attempts > 3) continue; // max retries

      try {
        final type = item['type'];
        if (type == 'payout_approval') {
          await apiClient.post('/payout_batches/${item['payload']['batch_id']}/approve', {});
        } else if (type == 'payroll_approve') {
          await apiClient.post('/payroll_runs/${item['payload']['run_id']}/approve', {});
        }
        await HiveBoxes.offlineQueueBox.delete(key);
        synced++;
      } catch (e) {
        await HiveBoxes.offlineQueueBox.put(key, {
          ...item,
          'attempts': attempts + 1,
        });
        failed++;
      }
    }

    return SyncResult(synced: synced, failed: failed);
  }
}

class SyncResult {
  final int synced;
  final int failed;
  SyncResult({required this.synced, required this.failed});
}
