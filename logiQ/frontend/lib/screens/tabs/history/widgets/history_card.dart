import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:frontend/models/models.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import '../../../../themes/themes.dart';
import 'package:frontend/utils/utils.dart';

class HistoryCard extends StatelessWidget {
  const HistoryCard({super.key, required this.quizItem, required this.index});

  final Quiz quizItem;
  final int index;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    final purpleGradient = AppGradients.historyCardPurpleGradient(theme);
    final yellowGradient = AppGradients.historyCardYellowGradient(theme);

    final cardBorderColor =
        index % 2 == 0
            ? theme.colorScheme.primary.withAlpha(128)
            : theme.colorScheme.secondary.withAlpha(128);
    final iconColor =
        index % 2 == 0
            ? theme.colorScheme.onPrimary
            : theme.colorScheme.onSecondary;
    final backgroundGradient = index % 2 == 0 ? purpleGradient : yellowGradient;

    final FaIcon icon =
        (() {
          switch (quizItem.type) {
            case QuizType.randomTasks:
              return FaIcon(
                FontAwesomeIcons.diceThree,
                size: 80,
                color: iconColor,
              );
            case QuizType.topicPractice:
              return FaIcon(
                FontAwesomeIcons.bullseye,
                size: 78,
                color: iconColor,
              );
            case QuizType.byDifficulty:
              return FaIcon(
                FontAwesomeIcons.bookOpen,
                size: 65,
                color: iconColor,
              );
            case QuizType.customQuiz:
              return FaIcon(
                FontAwesomeIcons.pencil,
                size: 68,
                color: iconColor,
              );
            default:
              return FaIcon(
                FontAwesomeIcons.diceThree,
                size: 80,
                color: iconColor,
              );
          }
        })();

    return GestureDetector(
      onTap: () {
        context.go('/history/quiz-review/${quizItem.id}/online');
      },
      child: Container(
        width: double.infinity,
        height: 120,
        clipBehavior: Clip.antiAlias,
        decoration: BoxDecoration(
          border: Border.all(color: cardBorderColor),
          borderRadius: AppRadii.medium,
        ),
        child: Row(
          children: [
            Container(
              height: 120,
              width: 120,
              decoration: BoxDecoration(
                gradient: backgroundGradient,
                borderRadius: AppRadii.medium,
              ),
              child: Center(child: icon),
            ),
            Gap(35),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    quizItem.type.displayName,
                    style: theme.textTheme.bodyLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Gap(6),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        '${quizItem.correctQuestionsNum}/10 Correct',
                        style: TextStyle(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                      Gap(5),
                      Text('·'),
                      Gap(5),
                      Text(
                        TimeFormatterUtil.formatCompletionTime(
                          quizItem.completionTime,
                        ),
                        style: TextStyle(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ),
                  Gap(8),
                  Text(
                    TimeFormatterUtil.formatDetailedTime(quizItem.completedAt),
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
