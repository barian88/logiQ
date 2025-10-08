import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/pods/pods.dart';

class OtpInputArea extends ConsumerWidget {
  const OtpInputArea({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {

    final verificationNotifier = ref.read(verificationNotifierProvider.notifier);

    final theme = Theme.of(context);

    return PinCodeTextField(
      appContext: context,
      length: 4,
      obscureText: false,
      // animationType: AnimationType.none,
      keyboardType: TextInputType.number,
      pinTheme: PinTheme(
        shape: PinCodeFieldShape.box,
        borderRadius: AppRadii.small,
        fieldHeight: 60,
        fieldWidth: 60,
        inactiveColor: Colors.grey.withAlpha(128),
        activeColor: theme.colorScheme.primary,
        selectedColor: theme.colorScheme.primary,
      ),
      onChanged: (value) {
        // print("Current input: $value");
        verificationNotifier.update(verificationCode: value);
      },

    );
  }
}
