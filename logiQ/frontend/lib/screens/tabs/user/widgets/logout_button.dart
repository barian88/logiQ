import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/pods/pods.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/utils/utils.dart';

class LogoutButton extends ConsumerWidget {
  const LogoutButton({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {

    final userNotifier = ref.watch(userNotifierProvider.notifier);
    final theme = Theme.of(context);

    return TextButton.icon(
      onPressed: () {
        userNotifier.logout();
        ToastHelper.showSuccess(theme, 'Logged out successfully');
        context.go('/login');
      },
      icon: Icon(
        Icons.logout,
        color: theme.colorScheme.red,
      ),
      label: Text(
        'Logout',
        style: TextStyle(
          color: theme.colorScheme.red,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
