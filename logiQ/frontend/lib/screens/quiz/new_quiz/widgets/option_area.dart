import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/utils/error_handler.dart';
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
        return _buildOptionArea(context, currentQuizQuestion);
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error));
  }

  Widget _buildOptionArea(BuildContext context, QuizQuestion currentQuizQuestion) {
    final theme = Theme.of(context);

    final choiceColors = [
      theme.colorScheme.blue,
      theme.colorScheme.secondary,
      theme.colorScheme.green,
      theme.colorScheme.red,
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        final maxWidth = constraints.maxWidth;
        if(currentQuizQuestion.question.type == QuestionType.trueFalse){
          return Column(
            children: [
              ChoiceOptionsContainer(
                option: 'True',
                color: choiceColors[2],
                maxWidth: maxWidth,
              ),
              Gap(20),
              ChoiceOptionsContainer(
                option: 'False',
                color: choiceColors[3],
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
                    color: choiceColors[0],
                    maxWidth: maxWidth,
                  ),
                  Gap(20),
                  ChoiceOptionsContainer(
                    option: 'C',
                    color: choiceColors[2],
                    maxWidth: maxWidth,
                  ),
                ],
              ),
              Gap(20),
              Column(
                children: [
                  ChoiceOptionsContainer(
                    option: 'B',
                    color: choiceColors[1],
                    maxWidth: maxWidth,
                  ),
                  Gap(20),
                  ChoiceOptionsContainer(
                    option: 'D',
                    color: choiceColors[3],
                    maxWidth: maxWidth,
                  ),
                ],
              ),
            ],
          );
        }
      },
    );
  }
}

class ChoiceOptionsContainer extends ConsumerWidget {
  const ChoiceOptionsContainer({
    super.key,
    required this.option,
    required this.maxWidth,
    required this.color,
  });

  final String option;
  final double maxWidth;
  final Color color;


  @override
  Widget build(BuildContext context, WidgetRef ref) {

   final quizState = ref.watch(quizNotifierProvider);
   final quizNotifier = ref.read(quizNotifierProvider.notifier);

   return quizState.when(
     data: (state) {
       final currentQuizQuestion = state.quiz.questions[state.currentQuestionIndex];
       final userAnswerIndex = currentQuizQuestion.userAnswerIndex;
       // 初始化时检查是否已选择该选项
       final currentOptionIndex = _getCurrentIndex();
       final isSelected = userAnswerIndex.contains(currentOptionIndex);

       return SizedBox(
      width: currentQuizQuestion.question.type == QuestionType.trueFalse ? maxWidth : (maxWidth - 20) / 2,
      child: AspectRatio(
        aspectRatio: currentQuizQuestion.question.type == QuestionType.trueFalse ? 4 : 1.6,
        child: InkWell(
          onTap: () {
            quizNotifier.setUserAnswerIndex(currentOptionIndex);},
          borderRadius: AppRadii.medium,
          child: Container(
            decoration: BoxDecoration(
              color: isSelected
                  ? color.withAlpha(140) // 选中状态的颜色变化
                  : color,
              borderRadius: AppRadii.medium,
              border: isSelected
                  ? Border.all(color: Colors.grey.withAlpha(128), width: 5)
                  : null,
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
      ),
    );
     },
     loading: () => const SizedBox.shrink(),
     error: (error, stackTrace) => const SizedBox.shrink(),
   );
  }

  int _getCurrentIndex() {
    switch (option){
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