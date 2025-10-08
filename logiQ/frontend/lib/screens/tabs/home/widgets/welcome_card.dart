import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../themes/themes.dart';
import 'package:frontend/pods/pods.dart';

class WelcomeCard extends ConsumerWidget {
  const WelcomeCard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userState = ref.watch(userNotifierProvider);
    final user = userState.value?.user;
    
    if (user == null) {
      return const SizedBox.shrink();
    }

    final theme = Theme.of(context);
    final purpleGradient = AppGradients.cardPurpleGradient(theme);

    return Container(
      width: double.infinity,
      height: 120,
      decoration: BoxDecoration(
        gradient: purpleGradient,
        borderRadius: AppRadii.medium,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                "Hi, ${user.username} 👋",
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                "Time to unlock some logic!",
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          CircleAvatar(
              radius: 35,
              backgroundImage: NetworkImage(user.profilePictureUrl)),
        ],
      ),
    );
  }
}
