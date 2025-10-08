import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/utils/utils.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/pods/pods.dart';

class Register extends ConsumerStatefulWidget {
  const Register({super.key});

  @override
  ConsumerState<Register> createState() => _RegisterState();
}

class _RegisterState extends ConsumerState<Register> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();

  bool _isPasswordVisible = false;
  bool _isConfirmPasswordVisible = false;

  @override
  void dispose() {
    _usernameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _handleRegister() async {
    if (_formKey.currentState!.validate()) {
      final registerNotifier = ref.read(registerNotifierProvider.notifier);
      final username = _usernameController.text;
      final email = _emailController.text;
      final password = _passwordController.text;

      final result = await registerNotifier.sendRegisterRequest(username, email, password);

      if (mounted) {
        if (result.isSuccess) {
          context.push('/verification/register/$email');
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
      appBar: AppBar(title: const Logo()),
      body: BaseContainer(
        isScrollable: true,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Gap(10),
            const Text(
              "Sign Up",
              style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
            ),
            const Gap(10),
            Text(
              "Create an account to continue",
              style: TextStyle(color: theme.colorScheme.onSurfaceVariant),
            ),
            const Gap(35),
            Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Username
                  Text('Username', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _usernameController,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a username';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      enabledBorder: UnderlineInputBorder(borderSide: BorderSide(width: 1, color: theme.colorScheme.primary)),
                      focusedBorder: UnderlineInputBorder(borderSide: BorderSide(width: 2, color: theme.colorScheme.primary)),
                      hintText: 'your username',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Email
                  Text('Email', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _emailController,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter an email';
                      }
                      if (!EmailUtil.isValidEmail(value)) {
                        return 'Please enter a valid email address';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      enabledBorder: UnderlineInputBorder(borderSide: BorderSide(width: 1, color: theme.colorScheme.primary)),
                      focusedBorder: UnderlineInputBorder(borderSide: BorderSide(width: 2, color: theme.colorScheme.primary)),
                      hintText: 'your email',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Password
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
                      hintText: 'your password',
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

                  // Confirm Password
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
                      hintText: 'your password again',
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
            const Gap(40),
            FilledButton(
              onPressed: _handleRegister,
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
                shape: RoundedRectangleBorder(borderRadius: AppRadii.medium),
              ),
              child: Text(
                "Continue",
                style: theme.textTheme.titleMedium?.copyWith(
                  color: Colors.white,
                ),
              ),
            ),
            const Gap(30),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  "Already have an account?",
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const Gap(5),
                InkWell(
                  onTap: () {
                    context.go("/login");
                  },
                  child: Text(
                    "Sign In",
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.primary,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}