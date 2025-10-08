import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/utils/utils.dart';

class QuizStat extends ConsumerWidget {
  const QuizStat({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);
    final theme = Theme.of(context);

    return quizState.when(
      data: (state) {
        final String completionTime = TimeFormatterUtil.formatCompletionTime(
          state.quiz.completionTime,
        );
        final int correctQuestionsNum = state.quiz.correctQuestionsNum;
        final int totalQuestions = state.quiz.questions.length;
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text('$correctQuestionsNum/$totalQuestions Correct', style: theme.textTheme.bodySmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant
            ),),
            Text('Completion time: $completionTime', style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant
            ),),
          ],
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(
        child: Text('Error: $error'),
      ),
    );
  }
}
