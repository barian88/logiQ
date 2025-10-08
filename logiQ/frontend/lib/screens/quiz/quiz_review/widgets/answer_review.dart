import 'package:flutter/material.dart';
import 'package:frontend/models/models.dart';
import 'package:frontend/themes/themes.dart';

class AnswerReview extends StatelessWidget {
  const AnswerReview({
    super.key,
    required this.correctAnswerIndex,
    required this.userAnswerIndex,
    required this.isCorrect,
    required this.questionType,
  });

  final List<int> correctAnswerIndex;
  final List<int> userAnswerIndex;
  final bool isCorrect;
  final QuestionType questionType;

  String _convertIndexesToLabels(List<int> indexes) {
    if (indexes.isEmpty) {
      return 'None';
    }

    if (questionType == QuestionType.trueFalse) {
      const labels = ['True', 'False'];
      return indexes.map((i) => labels[i]).join(', ');
    } else {
      const letters = ['A', 'B', 'C', 'D'];
      return indexes.map((i) => letters[i]).join(', ');
    }
  }


  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          'Correct Answer: ${_convertIndexesToLabels(correctAnswerIndex)}',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.green,
          ),
        ),
        Text(
          'Your Answer: ${_convertIndexesToLabels(userAnswerIndex)}',
          style: theme.textTheme.bodySmall?.copyWith(
            color: isCorrect ? theme.colorScheme.green : theme.colorScheme.red,
          ),
        ),
      ],
    );
  }
}
