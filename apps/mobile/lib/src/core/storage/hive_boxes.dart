import 'package:hive_flutter/hive_flutter.dart';

class HiveBoxes {
  static late Box draftLinksBox;
  static late Box offlineQueueBox;

  static Future init() async {
    draftLinksBox = await Hive.openBox('draft_links'); // offline draft links
    offlineQueueBox = await Hive.openBox('offline_queue');
  }
}
