import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/utils/utils.dart';
import 'package:frontend/pods/pods.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';

class Login extends ConsumerStatefulWidget {
  const Login({super.key});

  @override
  ConsumerState<Login> createState() => _LoginState();
}

class _LoginState extends ConsumerState<Login> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _passwordVisible = false;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void handleLogin() async {
    if (_formKey.currentState!.validate()) {
      final loginNotifier = ref.read(loginNotifierProvider.notifier);
      final email = _emailController.text;
      final password = _passwordController.text;

      final result = await loginNotifier.login(email, password);

      if (mounted) {
        if (result.isSuccess) {
          context.go('/home');
        } else {
          ErrorHandler.showErrorToast(context, result);
        }
      }
    }
  }

  void handleForgotPassword() async {
    final email = _emailController.text;
    if (email.isEmpty) {
      await ToastHelper.showWarning(Theme.of(context), "Please enter your email address first");
      return;
    }
    if (!EmailUtil.isValidEmail(email)) {
      await ToastHelper.showWarning(Theme.of(context), "Please enter a valid email address");
      return;
    }
    context.push('/verification/login/$email');
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Logo(), centerTitle: false),
      body: BaseContainer(
        isScrollable: true,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Gap(40),
            const Text(
              "Hello there 👋",
              style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
            ),
            const Gap(10),
            Text(
              "Enter your email and password to sign in ",
              style: TextStyle(color: theme.colorScheme.onSurfaceVariant),
            ),
            const Gap(50),
            Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Email',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _emailController,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter your email';
                      }
                      if (!EmailUtil.isValidEmail(value)) {
                        return 'Please enter a valid email address';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      enabledBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: theme.colorScheme.primary),
                      ),
                      focusedBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 2, color: theme.colorScheme.primary),
                      ),
                      hintText: 'your email',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Password',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextFormField(
                    controller: _passwordController,
                    obscureText: !_passwordVisible,
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter your password';
                      }
                      return null;
                    },
                    decoration: InputDecoration(
                      hintText: 'your password',
                      hintStyle: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                      enabledBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 1, color: theme.colorScheme.primary),
                      ),
                      focusedBorder: UnderlineInputBorder(
                        borderSide: BorderSide(width: 2, color: theme.colorScheme.primary),
                      ),
                      suffixIcon: IconButton(
                        icon: Icon(
                          _passwordVisible ? Icons.visibility : Icons.visibility_off,
                          size: 18,
                        ),
                        onPressed: () {
                          setState(() {
                            _passwordVisible = !_passwordVisible;
                          });
                        },
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const Gap(20),
            Row(
              children: [
                const Spacer(),
                InkWell(
                  onTap: handleForgotPassword,
                  child: Text(
                    "Forgot Password ?",
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.primary,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            const Gap(40),
            FilledButton(
              onPressed: handleLogin,
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
                shape: RoundedRectangleBorder(borderRadius: AppRadii.medium),
              ),
              child: Text(
                "SIGN IN",
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
                  "Don't have an account? ",
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const Gap(5),
                InkWell(
                  onTap: () {
                    context.push("/register");
                  },
                  child: Text(
                    "Sign Up",
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