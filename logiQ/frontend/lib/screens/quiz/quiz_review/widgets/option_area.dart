import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gap/gap.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/models/models.dart';
import 'package:frontend/pods/pods.dart';


class OptionArea extends ConsumerWidget {
  const OptionArea({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);

    return quizState.when(
      data: (state) {
        final currentQuizQuestion = state.quiz.questions[state.currentQuestionIndex];
        
        return LayoutBuilder(
          builder: (context, constraints) {
            final maxWidth = constraints.maxWidth;
            if(currentQuizQuestion.question.type == QuestionType.trueFalse){
              return Column(
                children: [
                  ChoiceOptionsContainer(
                    option: 'True',
                    maxWidth: maxWidth,
                  ),
                  Gap(20),
                  ChoiceOptionsContainer(
                    option: 'False',
                    maxWidth: maxWidth,
                  ),
                ],
              );
            }
            else{
              return Row(
                children: [
                  Column(
                    children: [
                      ChoiceOptionsContainer(
                        option: 'A',
                        maxWidth: maxWidth,
                      ),
                      Gap(20),
                      ChoiceOptionsContainer(
                        option: 'C',
                        maxWidth: maxWidth,
                      ),
                    ],
                  ),
                  Gap(20),
                  Column(
                    children: [
                      ChoiceOptionsContainer(
                        option: 'B',
                        maxWidth: maxWidth,
                      ),
                      Gap(20),
                      ChoiceOptionsContainer(
                        option: 'D',
                        maxWidth: maxWidth,
                      ),
                    ],
                  ),
                ],
              );
            }
          },
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(
        child: Text('Error: $error'),
      ),
    );
  }
}

class ChoiceOptionsContainer extends ConsumerWidget {
  const ChoiceOptionsContainer({
    super.key,
    required this.option,
    required this.maxWidth,
  });

  final String option;
  final double maxWidth;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);

    return quizState.when(
      data: (state) {
        final currentQuizQuestion = state.quiz.questions[state.currentQuestionIndex];
        final userAnswerIndex = currentQuizQuestion.userAnswerIndex;
        final correctAnswerIndex = currentQuizQuestion.question.correctAnswerIndex;

        final currentOptionIndex = _getCurrentIndex();

        // 判断当前选项颜色
        final theme = Theme.of(context);
        Color backgroundColor;
        if (correctAnswerIndex.contains(currentOptionIndex)) {
          backgroundColor = theme.colorScheme.green;
        } else if (userAnswerIndex.contains(currentOptionIndex)) {
          backgroundColor = theme.colorScheme.red;
        } else {
          backgroundColor = Colors.grey;
        }

        return SizedBox(
          width: currentQuizQuestion.question.type == QuestionType.trueFalse
              ? maxWidth
              : (maxWidth - 20) / 2,
          child: AspectRatio(
            aspectRatio: currentQuizQuestion.question.type == QuestionType.trueFalse ? 4 : 1.6,
            child: Container(
              decoration: BoxDecoration(
                color: backgroundColor,
                borderRadius: AppRadii.medium,
                border: Border.all(color: Colors.grey.withAlpha(64), width: 1),
              ),
              child: Center(
                child: Text(
                  option,
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
          ),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(
        child: Text('Error: $error'),
      ),
    );
  }

  int _getCurrentIndex() {
    switch (option) {
      case 'A':
        return 0;
      case 'B':
        return 1;
      case 'C':
        return 2;
      case 'D':
        return 3;
      case 'True':
        return 0;
      case 'False':
        return 1;
      default:
        return -1;
    }
  }
}


