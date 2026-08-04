import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

class QrScannerPage extends StatefulWidget {
  const QrScannerPage({super.key});
  @override State<QrScannerPage> createState() => _QrScannerPageState();
}

class _QrScannerPageState extends State<QrScannerPage> {
  final controller = MobileScannerController();
  String? result;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Scan QR • QR ስካን')),
      body: Stack(
        children: [
          MobileScanner(controller: controller, onDetect: (barcodes){
            final code = barcodes.barcodes.first.rawValue;
            if (code!=null) setState(()=> result=code);
          }),
          // Outstanding overlay with rounded square guides
          Center(
            child: Container(
              width: 260, height: 260,
              decoration: BoxDecoration(border: Border.all(color: Colors.white, width: 3), borderRadius: BorderRadius.circular(24)),
              child: Stack(children: [
                // corner brackets animated
                Positioned(top:0,left:0,child: Container(width:30,height:30,decoration: BoxDecoration(border: Border(left: BorderSide(color: Colors.greenAccent,width: 4), top: BorderSide(color: Colors.greenAccent,width:4))))),
              ]),
            ),
          ),
          if (result!=null)
            Positioned(bottom:20,left:20,right:20,child: Card(child: Padding(padding: const EdgeInsets.all(16), child: Column(children: [
              const Icon(Icons.check_circle, color: Colors.green, size: 48),
              const SizedBox(height:8),
              Text('Scanned • ተቃኝቷል:\n$result', textAlign: TextAlign.center),
              const SizedBox(height:12),
              ElevatedButton(onPressed: ()=> Navigator.pop(context), child: const Text('Done • ጨርስ')),
            ])))),
        ],
      ),
    );
  }
}
