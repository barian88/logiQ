import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/utils/error_handler.dart';
import 'package:gap/gap.dart';
import '../../../widgets/widgets.dart';
import '../../../pods/history/history_pod.dart';
import 'widgets/widgets.dart';

class History extends ConsumerStatefulWidget {
  const History({super.key});

  @override
  ConsumerState<History> createState() => _HistoryState();
}

class _HistoryState extends ConsumerState<History> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(historyNotifierProvider.notifier).loadAllHistory();
    });
  }

  @override
  Widget build(BuildContext context) {
    final historyState = ref.watch(historyNotifierProvider);

    return BaseContainer(
      child: historyState.when(
        data: (quizzes) {
          if (quizzes.isEmpty) {
            return Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text('No quiz history available.', style: Theme.of(context).textTheme.bodyLarge),
                const Gap(10),
                Text('Start a new quiz to see your history here!', style: Theme.of(context).textTheme.bodyMedium)
              ],
            );
          }
          
          final cardList = List.generate(quizzes.length, (index) {
            final quiz = quizzes[index];
            return Column(
              children: [
                HistoryCard(quizItem: quiz, index: index),
                const Gap(30),
              ],
            );
          });

          return Column(children: cardList);
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error,
            onRetry: () {
              ref.read(historyNotifierProvider.notifier).loadAllHistory();
            }, context: context
        ),
      ),
    );
  }
}
