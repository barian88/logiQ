import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gap/gap.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/models/models.dart';

import '../../../../utils/utils.dart';

class QuestionArea extends ConsumerWidget {
  const QuestionArea({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);
    final theme = Theme.of(context);

    return quizState.when(
      data: (state) {
        final currentQuizQuestion =
            state.quiz.questions[state.currentQuestionIndex];
        return _buildQuestionArea(context, theme, currentQuizQuestion);
      },
      loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error)
    );
  }

  Widget _buildQuestionArea(
    BuildContext context,
    ThemeData theme,
    QuizQuestion currentQuizQuestion,
  ) {
    // This should be replaced with the actual options data
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            "Question",
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          Row(
            children: [
              Chip(
                label: Text(
                  currentQuizQuestion.question.type.displayName,
                  style: theme.textTheme.bodySmall,
                ),
                padding: EdgeInsets.all(0),
              ),
              Gap(8),
              Chip(
                label: Text(
                  currentQuizQuestion.question.category.displayName,
                  style: theme.textTheme.bodySmall,
                ),
                padding: EdgeInsets.all(0),
              ),
              Gap(8),
              Chip(
                label: Text(
                  currentQuizQuestion.question.difficulty.displayName,
                  style: theme.textTheme.bodySmall,
                ),
                padding: EdgeInsets.all(0),
              ),
            ],
          ),
          Gap(8),
          Text(
            currentQuizQuestion.question.questionText,
            style: theme.textTheme.bodyLarge?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          Gap(16),
          ...getOptionWidgets(currentQuizQuestion.question.options, theme),
        ],
      ),
    );
  }
}

List<Widget> getOptionWidgets(List<String> options, ThemeData theme) {
  final tags = ['A', 'B', 'C', 'D'];
  final optionWidgets = <Widget>[];
  for (int i = 0; i < options.length; i++) {
    final option = options[i];
    optionWidgets.add(
      Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '${tags[i]}. ',
            style: theme.textTheme.bodyLarge?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const Gap(8),
          Expanded(
            child: Text(
              option,
              style: theme.textTheme.bodyLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
              softWrap: true,
            ),
          ),
        ],
      ),
    );
    if (i != options.length - 1) {
      optionWidgets.add(Gap(8)); // Add space between options
    }
  }
  return optionWidgets;
}
