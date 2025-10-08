import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/utils/utils.dart';
import 'package:frontend/pods/pods.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';

class OperationArea extends ConsumerWidget {
  const OperationArea({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);
    final quizNotifier = ref.read(quizNotifierProvider.notifier);

    final theme = Theme.of(context);

    return quizState.when(
      data: (state) {
        final isFirst = state.currentQuestionIndex == 0;
        final isLast =
            state.currentQuestionIndex == state.quiz.questions.length - 1;
        final backStatus = isFirst ? null : quizNotifier.previousQuestion;
        final nextStatus =
            isLast
                ? () => _handleSubmit(context, ref)
                : quizNotifier.nextQuestion;

        return _buildOperationArea(
          theme,
          isFirst,
          isLast,
          backStatus,
          nextStatus,
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stackTrace) => Center(child: Text('Error: $error')),
    );
  }

  Widget _buildOperationArea(
    ThemeData theme,
    bool isFirst,
    bool isLast,
    VoidCallback? backStatus,
    VoidCallback? nextStatus,
  ) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        ElevatedButton(
          onPressed: backStatus,
          style: ElevatedButton.styleFrom(
            backgroundColor: theme.colorScheme.secondary, // Change the color to red
            minimumSize: Size(10, 50),
          ),
          child: Icon(Icons.arrow_back_ios_new, color: Colors.white),
        ),
        Gap(30),
        Expanded(
          child: ElevatedButton(
            onPressed: nextStatus,
            style: ElevatedButton.styleFrom(
              backgroundColor: theme.colorScheme.primary,
              overlayColor: theme.colorScheme.onPrimary.withOpacity(0.1),
              // Change the color to red
              minimumSize: Size(10, 50),
            ),
            child: Text(
              isLast ? 'Submit' : 'Next',
              style: theme.textTheme.bodyLarge?.copyWith(
                color: Colors.white, // Ensure text is readable
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
      ],
    );
  }

  void _handleSubmit(BuildContext context, WidgetRef ref) async {
    final quizNotifier = ref.read(quizNotifierProvider.notifier);
    final quizState = ref.watch(quizNotifierProvider);

    showCupertinoDialog(
      context: context,
      builder: (BuildContext context) {
        return CupertinoAlertDialog(
          title: const Text('Submit'),
          content: Column(
            children: [
              const Text('Are you sure to submit the quiz?'),
              if (_getUnansweredQuestionsMessage(ref).isNotEmpty) 
                Text(_getUnansweredQuestionsMessage(ref)),
            ],
          ),
          actions: [
            CupertinoDialogAction(
              child: const Text('Cancel'),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
            CupertinoDialogAction(
              isDestructiveAction: true,
              child: const Text('Yes'),
              onPressed: () async{
                Navigator.of(context).pop();
                // 确认提交
                final result = await quizNotifier.submitQuiz();
                if(result.isSuccess){
                //   弹出对话框
                  showCupertinoDialog(
                    context: context,
                    builder: (context) => CupertinoAlertDialog(
                      content: Text("Submit Successfully!"),
                      actions: [
                        CupertinoDialogAction(
                          child: Text("View Result", style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: Theme.of(context).colorScheme.primary
                          ),),
                          onPressed: () {
                            Navigator.pop(context);
                            final mode = quizState.value?.mode ?? QuizMode.online;
                            context.go('/history/quiz-review/${result.quizId}/${mode.name}');
                          },
                        ),
                      ],
                    ),
                  );
                }else{
                  ErrorHandler.showErrorToast(context, result);
                }
              },
            ),
          ],
        );
      },
    );
  }

  String _getUnansweredQuestionsMessage(WidgetRef ref) {
    final quizState = ref.watch(quizNotifierProvider);
    final questions = quizState.value?.quiz.questions ?? [];
    final unansweredQuestions = <int>[];
    for (int i = 0; i < questions.length; i++) {
      if (questions[i].userAnswerIndex.isEmpty) {
        unansweredQuestions.add(i + 1);
      }
    }
    return unansweredQuestions.isNotEmpty
        ? 'Unanswered question(s): ${unansweredQuestions.join(', ')}'
        : '';
  }
}
