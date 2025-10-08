import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import '../../../utils/error_handler.dart';
import '../../../widgets/widgets.dart';
import 'widgets/widgets.dart';
import 'package:frontend/pods/pods.dart';

class Home extends ConsumerStatefulWidget {
  const Home({super.key});

  @override
  ConsumerState<Home> createState() => _HomeState();
}

class _HomeState extends ConsumerState<Home> {
  @override
  void initState() {
    super.initState();
    _initializeUserData();
  }

  Future<void> _initializeUserData() async {
    final userNotifier = ref.read(userNotifierProvider.notifier);
    await userNotifier.loadUserFromStorage();
  }

  @override
  Widget build(BuildContext context) {
    final userState = ref.watch(userNotifierProvider);
    
    return BaseContainer(
      child: userState.when(
        data: (userData) => Column(
          children: [
            WelcomeCard(),
            Gap(26),
            QuizButtons(),
            Gap(26),
            Performance(),
          ],
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error,
        onRetry: _initializeUserData, context: context),
      ),
    );
  }
}


