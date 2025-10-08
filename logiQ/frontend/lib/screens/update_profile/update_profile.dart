import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/themes/colors.dart';
import 'package:frontend/utils/utils.dart';
import 'package:frontend/widgets/containers/base_container.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';

class UpdateProfile extends ConsumerStatefulWidget {
  const UpdateProfile({super.key});

  @override
  ConsumerState<UpdateProfile> createState() => _UpdateProfileState();
}

class _UpdateProfileState extends ConsumerState<UpdateProfile> {
  final _usernameController = TextEditingController();

  @override
  void initState() {
    super.initState();
    final user = ref.read(userNotifierProvider).value?.user;
    if (user != null) {
      _usernameController.text = user.username;
    }
  }

  @override
  void dispose() {
    _usernameController.dispose();
    super.dispose();
  }

  Future<void> _updateProfile() async {
    final newUsername = _usernameController.text;
    final userNotifier = ref.read(userNotifierProvider.notifier);

    final result = await userNotifier.updateProfile(newUsername, "");

    if (mounted) {
      final theme = Theme.of(context);
      if (result.isSuccess) {
        ToastHelper.showSuccess(theme, "Username updated successfully!");
        context.go('/user'); // back to user page
      } else {
        // Show error message
        ErrorHandler.showErrorToast(context, result);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text('Change Username', style: theme.textTheme.bodyLarge),
        centerTitle: true,
        actions: [
          TextButton(
            onPressed: _updateProfile,
            child: Text(
              'Save',
              style: TextStyle(color: theme.colorScheme.primary, fontSize: 16),
            ),
          ),
          const Gap(16),
        ],
      ),
      body: BaseContainer(child: Form(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'New Username',
              style: theme.textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: AppColors.grey1,
              ),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: _usernameController,
              decoration: InputDecoration(
                enabledBorder: UnderlineInputBorder(
                  borderSide: BorderSide(width: 1, color: theme.colorScheme.primary),
                ),
                focusedBorder: UnderlineInputBorder(
                  borderSide: BorderSide(width: 2, color: theme.colorScheme.primary),
                ),
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter a username';
                }
                if (value.length > 20) {
                  return 'Username cannot exceed 20 characters';
                }
                return null;
              },
            ),
            const SizedBox(height: 8),
            Text(
              'Username must be at most 20 characters long.',
              style: theme.textTheme.bodySmall?.copyWith(
                color: AppColors.grey1,
              ),
            ),
          ],
        ),
      )),
    );
  }
}