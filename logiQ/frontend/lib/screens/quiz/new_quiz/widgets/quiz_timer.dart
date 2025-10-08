import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/utils/utils.dart';


class QuizTimer extends ConsumerStatefulWidget {
  const QuizTimer ({super.key});
  @override
  ConsumerState<QuizTimer> createState() => _QuizTimerState();
}

class _QuizTimerState extends ConsumerState<QuizTimer> {
  @override
  void initState(){
    super.initState();
    // 计时器现在在 loadNewQuiz 完成后自动启动
  }

@override
  Widget build(BuildContext context) {

    final quizState = ref.watch(quizNotifierProvider);
    final theme = Theme.of(context);

    return quizState.when(
      data: (state) {
        final time = state.quiz.completionTime;
        return Text(TimeFormatterUtil.getFormattedTime(time), 
          style: theme.textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.bold
          ));
      },
      loading: () => Text('--:--', 
        style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.bold
        )),
      error: (error, stackTrace) => Text('--:--', 
        style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.bold
        )),
    );
  }
}
