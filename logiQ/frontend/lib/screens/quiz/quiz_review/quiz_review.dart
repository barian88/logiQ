import 'dart:ffi';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/utils/error_handler.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/pods/pods.dart';
import 'package:gap/gap.dart';
import 'widgets/widgets.dart';

class QuizReview extends ConsumerStatefulWidget {
  const QuizReview({super.key, required this.quizId, required this.mode});
  final String quizId;
  final String mode;

  @override
  ConsumerState<QuizReview> createState() => _QuizReviewState();
}

class _QuizReviewState extends ConsumerState<QuizReview> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      _initQuizReview();
    });
  }

  Future<void> _initQuizReview() async {
    final quizNotifier = ref.read(quizNotifierProvider.notifier);
    quizNotifier.setMode(widget.mode);
    await quizNotifier.loadQuizReview(widget.quizId);
  }

  @override
  Widget build(BuildContext context) {

    final quizState = ref.watch(quizNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: Logo(),
        centerTitle: true,
        actions: [ThemeModeSwitch(), const Gap(16)],
      ),
      body: quizState.when(
        data: (state) {
          if (state.quiz.questions.isEmpty) {
            return const Center(
              child: Text('No quiz data available'),
            );
          }
          
          final currentQuizQuestion = state.quiz.questions[state.currentQuestionIndex];
          final questionType = currentQuizQuestion.question.type;
          final userAnswerIndex = currentQuizQuestion.userAnswerIndex;
          final correctAnswerIndex = currentQuizQuestion.question.correctAnswerIndex;
          final isCorrect = currentQuizQuestion.isCorrect;

          return BaseContainer(
            isScrollable: false,
            child: Column(
              children: [
                QuizStat(),
                const Gap(12),
                QuizProgress(),
                const Gap(30),
                Expanded(flex: 4, child: QuestionArea()),
                Divider(height: 0),
                const Gap(10),
                AnswerReview(correctAnswerIndex: correctAnswerIndex, userAnswerIndex: userAnswerIndex, isCorrect: isCorrect, questionType: questionType),
                const Gap(10),
                Expanded(
                  flex: 5,
                  child: Column(
                    children: [OptionArea(), Spacer(), OperationArea(), Gap(35)],
                  ),
                ),
              ],
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error,
          onRetry: _initQuizReview,
          context: context
        ),
      ),
    );
  }
}
