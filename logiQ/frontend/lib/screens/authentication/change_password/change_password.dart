import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/utils/utils.dart';
import 'package:gap/gap.dart';
import 'package:frontend/pods/pods.dart';
import 'package:go_router/go_router.dart';

class ChangePassword extends ConsumerStatefulWidget {
  const ChangePassword({super.key, required this.temporaryToken});

  final String temporaryToken;

  @override
  ConsumerState<ChangePassword> createState() => _ChangePasswordState();
}

class _ChangePasswordState extends ConsumerState<ChangePassword> {
  final _formKey = GlobalKey<FormState>();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  bool _isPasswordVisible = false;
  bool _isConfirmPasswordVisible = false;

  @override
  void dispose() {
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _handleChangePassword() async {
    if (_formKey.currentState!.validate()) {
      final changePasswordNotifier = ref.read(changePasswordNotifierProvider.notifier);
      final newPassword = _passwordController.text;

      final result = await changePasswordNotifier.resetPassword(widget.temporaryToken, newPassword);

      if (mounted) {
        if (result.isSuccess) {
          await ToastHelper.showSuccess(Theme.of(context), 'Password changed successfully, please sign in again');
          context.go('/login');
        } else {
          ErrorHandler.showErrorToast(context, result);
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Logo(),
      ),
      body: BaseContainer(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Gap(30),
            const Text(
              "Change Password 🔑",
              style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
            ),
            const Gap(40),
            Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Set Password', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _passwordController,
                    obscureText: !_isPasswordVisible,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a password';
                      }
                      if (value.length < 6) {
                        return 'Password must be at least 6 characters';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      hintText: 'your new password',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                      enabledBorder: UnderlineInputBorder(borderSide: BorderSide(width: 1, color: theme.colorScheme.primary)),
                      focusedBorder: UnderlineInputBorder(borderSide: BorderSide(width: 2, color: theme.colorScheme.primary)),
                      suffixIcon: IconButton(
                        icon: Icon(_isPasswordVisible ? Icons.visibility : Icons.visibility_off, size: 18),
                        onPressed: () => setState(() => _isPasswordVisible = !_isPasswordVisible),
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text('Confirm Password', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _confirmPasswordController,
                    obscureText: !_isConfirmPasswordVisible,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please confirm your password';
                      }
                      if (value != _passwordController.text) {
                        return 'Passwords do not match';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      hintText: 'your new password again',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                      enabledBorder: UnderlineInputBorder(borderSide: BorderSide(width: 1, color: theme.colorScheme.primary)),
                      focusedBorder: UnderlineInputBorder(borderSide: BorderSide(width: 2, color: theme.colorScheme.primary)),
                      suffixIcon: IconButton(
                        icon: Icon(_isConfirmPasswordVisible ? Icons.visibility : Icons.visibility_off, size: 18),
                        onPressed: () => setState(() => _isConfirmPasswordVisible = !_isConfirmPasswordVisible),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const Gap(60),
            FilledButton(
              onPressed: _handleChangePassword,
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
                shape: RoundedRectangleBorder(borderRadius: AppRadii.medium),
              ),
              child: Text(
                "Confirm",
                style: theme.textTheme.titleMedium?.copyWith(
                  color: Colors.white,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}