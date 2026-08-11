import 'package:flutter/material.dart';
import 'package:share_plus/share_plus.dart';
import '../../../core/theme/app_theme.dart';

class CreateLinkSheet extends StatefulWidget {
  const CreateLinkSheet({super.key});
  @override State<CreateLinkSheet> createState() => _CreateLinkSheetState();
}

class _CreateLinkSheetState extends State<CreateLinkSheet> {
  final _amountCtrl = TextEditingController(text: '500');
  final _descCtrl = TextEditingController(text: 'Tutoring • አስተማሪ');
  String? _linkUrl;
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create Link • ሊንክ ፍጠር')),
      body: DraggableScrollableSheet(
        initialChildSize: 0.9,
        builder: (context, controller) => ListView(
          controller: controller,
          padding: const EdgeInsets.all(20),
          children: [
            // Amount chips for outstanding UX
            const Text('Amount • መጠን ETB'),
            const SizedBox(height: 12),
            Wrap(spacing: 8, children: [100,500,1000,5000].map((a) => ChoiceChip(label: Text('ETB $a'), selected: _amountCtrl.text==a.toString(), onSelected: (s){ setState(()=> _amountCtrl.text=a.toString()); })).toList()),
            const SizedBox(height: 16),
            TextField(controller: _amountCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Custom amount', prefixText: 'ETB ')),
            const SizedBox(height: 16),
            TextField(controller: _descCtrl, decoration: const InputDecoration(labelText: 'Description • መግለጫ')),
            const SizedBox(height: 24),
            if (_linkUrl==null)
              ElevatedButton(
                onPressed: _loading ? null : () async {
                  setState(()=> _loading=true);
                  await Future.delayed(const Duration(seconds: 1));
                  setState(() {
                    _linkUrl='https://checkout.apexpay.et/c/abc123';
                    _loading=false;
                  });
                },
                child: _loading ? const CircularProgressIndicator(color: Colors.white) : const Text('Generate Link • ሊንክ አመንጭ'),
              ),
            if (_linkUrl!=null) ...[
              const SizedBox(height: 20),
              GlassCard(child: Column(children: [
                const Icon(Icons.qr_code_2, size: 120, color: AppColors.primary),
                const SizedBox(height: 12),
                Text(_linkUrl!, style: const TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 12),
                Row(children: [
                  Expanded(child: ElevatedButton.icon(onPressed: ()=> Share.share('Pay ETB ${_amountCtrl.text} for ${_descCtrl.text}: $_linkUrl'), icon: const Icon(Icons.share), label: const Text('Share via Telegram/WhatsApp'))),
                ]),
              ])),
            ]
          ],
        ),
      ),
    );
  }
}
