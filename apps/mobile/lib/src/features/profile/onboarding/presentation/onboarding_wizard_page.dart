import 'package:flutter/material.dart';
import '../../../../core/theme/app_theme.dart';

// Outstanding 6-step wizard mirroring web wizard but mobile-optimized with Fayda front/back camera
class OnboardingWizardPage extends StatefulWidget {
  const OnboardingWizardPage({super.key});
  @override State<OnboardingWizardPage> createState() => _OnboardingWizardPageState();
}

class _OnboardingWizardPageState extends State<OnboardingWizardPage> {
  int _step = 0;
  final PageController _pc = PageController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Onboarding • ምዝገባ'), leading: IconButton(icon: const Icon(Icons.close), onPressed: ()=> Navigator.pop(context))),
      body: Column(
        children: [
          // Progress donut + stepper outstanding
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(children: List.generate(6, (i)=> Expanded(child: Container(height: 6, margin: const EdgeInsets.symmetric(horizontal: 3), decoration: BoxDecoration(color: i<=_step ? AppColors.primary : Colors.black12, borderRadius: BorderRadius.circular(3)))))),
          ),
          Expanded(
            child: PageView(
              controller: _pc,
              onPageChanged: (i)=> setState(()=>_step=i),
              children: [
                _stepBusiness(),
                _stepOwnersFayda(),
                _stepBank(),
                _stepDocs(),
                _stepCompliance(),
                _stepReview(),
              ],
            ),
          ),
          Padding(padding: const EdgeInsets.all(16), child: Row(children: [
            if (_step>0) Expanded(child: OutlinedButton(onPressed: ()=> _pc.previousPage(duration: const Duration(milliseconds: 300), curve: Curves.easeOut), child: const Text('Back'))),
            const SizedBox(width: 12),
            Expanded(child: ElevatedButton(onPressed: ()=> _pc.nextPage(duration: const Duration(milliseconds: 300), curve: Curves.easeOut), child: Text(_step==5?'Submit • አስገባ':'Next • ቀጣይ'))),
          ])),
        ],
      ),
    );
  }

  Widget _stepBusiness() => ListView(padding: const EdgeInsets.all(20), children: [
    const Text('Business Info • የንግድ መረጃ', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
    const SizedBox(height:16),
    const TextField(decoration: InputDecoration(labelText: 'Legal Name • ህጋዊ ስም')),
    const SizedBox(height:12),
    const TextField(decoration: InputDecoration(labelText: 'TIN • 10 digits')),
    const SizedBox(height:12),
    DropdownButtonFormField(items: const [DropdownMenuItem(value:'plc',child: Text('PLC')), DropdownMenuItem(value:'sole',child: Text('Sole Prop'))], onChanged: (v){}, decoration: const InputDecoration(labelText: 'Business Type')),
  ]);

  Widget _stepOwnersFayda() => ListView(padding: const EdgeInsets.all(20), children: [
    const Text('Owners & Fayda ID Verification', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
    const SizedBox(height:8),
    const Text('Capture front/back of Fayda card + selfie + OTP — per NBE + id.gov.et', style: TextStyle(color: Colors.black54)),
    const SizedBox(height:16),
    Card(child: ListTile(leading: const Icon(Icons.badge, color: AppColors.primary), title: const Text('Abebe Kebede • Owner 100%'), subtitle: const Text('Fayda • FIN ****1234 • Verified 0.92'), trailing: const Icon(Icons.check_circle, color: Colors.green))),
    const SizedBox(height:12),
    _faydaCaptureCard('Fayda Front • de face', Icons.credit_card, 'Capture front image • የፊት ገጽታ'),
    const SizedBox(height:12),
    _faydaCaptureCard('Fayda Back • dos', Icons.credit_card_outlined, 'Capture back image • የኋላ ገጽታ'),
    const SizedBox(height:12),
    _faydaCaptureCard('Selfie • ራስ ፎቶ', Icons.face, 'Liveness selfie • የፊት ፎቶ'),
    const SizedBox(height:12),
    const TextField(decoration: InputDecoration(labelText: 'FIN 12-digit or FAN 16 alias • ፋይዳ ቁጥር', prefixIcon: Icon(Icons.numbers)), keyboardType: TextInputType.number),
    const SizedBox(height:12),
    const TextField(decoration: InputDecoration(labelText: 'OTP 6-digit • የተላከ ኮድ', prefixIcon: Icon(Icons.sms)), keyboardType: TextInputType.number),
  ]);

  Widget _faydaCaptureCard(String title, IconData icon, String action) => GlassCard(child: Row(children: [
    Icon(icon, color: AppColors.primary, size: 32),
    const SizedBox(width:12),
    Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [Text(title, style: const TextStyle(fontWeight: FontWeight.w600)), Text(action, style: const TextStyle(fontSize:12,color: Colors.black54))])),
    const Icon(Icons.camera_alt, color: AppColors.primary),
  ]));

  Widget _stepBank() => ListView(padding: const EdgeInsets.all(20), children: const [
    Text('Bank Account • የባንክ ሂሳብ', style: TextStyle(fontSize:20,fontWeight: FontWeight.bold)),
    SizedBox(height:16),
    TextField(decoration: InputDecoration(labelText: 'Bank — e.g., Commercial Bank Ethiopia', prefixIcon: Icon(Icons.account_balance))),
    SizedBox(height:12),
    TextField(decoration: InputDecoration(labelText: 'Account Number', prefixIcon: Icon(Icons.numbers))),
    SizedBox(height:12),
    TextField(decoration: InputDecoration(labelText: 'Account Name must match Legal Name')),
  ]);

  Widget _stepDocs() => ListView(padding: const EdgeInsets.all(20), children: [
    const Text('Documents Vault • ሰነዶች', style: TextStyle(fontSize:20,fontWeight: FontWeight.bold)),
    const LinearProgressIndicator(value: 0.6, color: AppColors.primary),
    const SizedBox(height:12),
    ...['Company Registration • የኩባንያ ምዝገባ','TIN Certificate • ቲን ሰርተፍኬት','Business License • ንግድ ፈቃድ','Bank Letter • የባንክ ደብዳቤ','Fayda Front • ፋይዳ ፊት','Fayda Back • ፋይዳ ጀርባ'].map((t)=> Card(child: ListTile(leading: const Icon(Icons.picture_as_pdf, color: Colors.red), title: Text(t), trailing: const Icon(Icons.cloud_upload, color: AppColors.primary)))),
  ]);

  Widget _stepCompliance() => ListView(padding: const EdgeInsets.all(20), children: [
    const Text('Compliance Preview • ተገዢነት', style: TextStyle(fontSize:20,fontWeight: FontWeight.bold)),
    const SizedBox(height:12),
    const GlassCard(child: Row(children: [Icon(Icons.verified_user, color: Colors.green, size:32), SizedBox(width:12), Expanded(child: Text('Risk: Medium • 42/100\nFayda: Verified 0.92 • OTP • consent logged\nTIN: Valid • Bank: Verified • Restricted Industry: Pass\nWebsite: refund ✓ privacy ✓ terms ✓'))])),
  ]);

  Widget _stepReview() => ListView(padding: const EdgeInsets.all(20), children: const [
    Text('Review & Submit', style: TextStyle(fontSize:20,fontWeight: FontWeight.bold)),
    SizedBox(height:12),
    Text('By submitting, you confirm business info true per NBE ONPS/02/2020 directive and consent Fayda verification via id.gov.et with OTP consent. Files encrypted at rest, FIN hashed, front/back images <2MB.'),
    SizedBox(height:20),
    Text('After submit, compliance team reviews in Kanban board. Dual approval if risk high or TPV>1M ETB. Test keys immediately, live after approval.', style: TextStyle(color: Colors.black54)),
  ]);
}
