import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/models/models.dart';
import 'package:go_router/go_router.dart';
import '../../../../themes/themes.dart';

class QuizButtonContainer extends ConsumerWidget {
  final String title;
  final String emoji;
  final int index;

  const QuizButtonContainer({
    super.key,
    required this.title,
    required this.emoji,
    required this.index,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final gradient =
        index % 2 == 0
            ? AppGradients.buttonPurpleGradient(theme)
            : AppGradients.buttonYellowGradient(theme);

    return Material(
      elevation: 1,
      surfaceTintColor:
          theme.brightness == Brightness.dark
              ? Theme.of(context).colorScheme.surfaceTint
              : Theme.of(context).colorScheme.surface,
      borderRadius: AppRadii.medium,
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () {
          switch (index) {
            case 0:
              // 使用了两种传参方式，一种是路径参数（必须的mode），一种是查询参数
              context.push('/home/new-quiz/normal?type=${QuizType.randomTasks.name}');
              break;
            case 1:
              _showPicker(context, QuizType.topicPractice);
              break;
            case 2:
              _showPicker(context, QuizType.byDifficulty);
              break;
            case 3:
              _showPicker(context, QuizType.customQuiz);
              break;
            // offline mode
            case 4:
              _showDialog(context);
              break;
          }
        },
        child: Column(
          children: [
            Ink(
              height: 100,
              width: 140,
              decoration: BoxDecoration(gradient: gradient),
              child: Center(child: Text(emoji, style: TextStyle(fontSize: 46))),
            ),
            SizedBox(height: 40, child: Center(child: Text(title))),
          ],
        ),
      ),
    );
  }

  void _showPicker(BuildContext context, QuizType quizType) {
    final topicList = [
      QuestionCategory.truthTable.displayName.toString(),
      QuestionCategory.equivalence.displayName.toString(),
      QuestionCategory.inference.displayName.toString(),
    ];
    final difficultyList = [
      QuestionDifficulty.easy.displayName.toString(),
      QuestionDifficulty.medium.displayName.toString(),
      QuestionDifficulty.hard.displayName.toString(),
    ];

    // 根据不同的type显示不同的选项
    int selectedIndexTopic = 0;
    int selectedIndexDifficulty = 0;

    final Widget picker;

    switch (quizType) {
      case QuizType.topicPractice:
        picker = CupertinoPicker(
          itemExtent: 45,
          onSelectedItemChanged: (index) {
            selectedIndexTopic = index;
          },
          children: topicList.map((e) => Center(child: Text((e)))).toList(),
        );
        break;
      case QuizType.byDifficulty:
        picker = CupertinoPicker(
          itemExtent: 45,
          onSelectedItemChanged: (index) {
            selectedIndexDifficulty = index;
          },
          children: difficultyList.map((e) => Center(child: Text((e)))).toList(),
        );
        break;
      case QuizType.customQuiz:
        picker = Row(
          children: [
            Expanded(child: CupertinoPicker(
              itemExtent: 45,
              onSelectedItemChanged: (index) {
                selectedIndexTopic = index;
              },
              children: topicList.map((e) => Center(child: Text((e)))).toList(),
            )),
            Expanded(child: CupertinoPicker(
              itemExtent: 45,
              onSelectedItemChanged: (index) {
                selectedIndexDifficulty = index;
              },
              children: difficultyList.map((e) => Center(child: Text((e)))).toList(),
            ))
          ],
        );
        break;
      default:
        picker = Center(child: Text('Unsupported quiz type'));
        break;
    }


    showModalBottomSheet(
      context: context,
      builder: (_) => SizedBox(
        height: 250,
        child: Column(
          children: [
            CupertinoButton(
              child: Text('OK'),
              onPressed: () {
                switch (quizType) {
                  case QuizType.topicPractice:
                    context.push('/home/new-quiz/normal?type=${QuizType.topicPractice.name}&category=${topicList[selectedIndexTopic]}');
                    break;
                    case QuizType.byDifficulty:
                    context.push('/home/new-quiz/normal?type=${QuizType.byDifficulty.name}&difficulty=${difficultyList[selectedIndexDifficulty]}');
                    break;
                  case QuizType.customQuiz:
                    context.push('/home/new-quiz/normal?type=${QuizType.customQuiz.name}&category=${topicList[selectedIndexTopic]}&difficulty=${difficultyList[selectedIndexDifficulty]}');
                    break;
                  default:
                    context.push('/home/new-quiz/normal?type=${QuizType.randomTasks.name}');
                    break;
                }
                // 关闭底部弹出框
                Navigator.pop(context);
              },
            ),
            Expanded(
              child: picker,
            ),
          ],
        ),
      ),
    );
  }

  void _showDialog(BuildContext context) {
    showCupertinoDialog(
      context: context,
      builder: (dialogContext) => CupertinoAlertDialog(
        title: const Text('Offline Practice'),
        content: Text(
          'You are about to start an offline quiz. Questions come from the local question bank and the results will not be recorded.',
        ),
        actions: [
          CupertinoDialogAction(
            isDefaultAction: true,
            onPressed: () {
              Navigator.of(dialogContext).pop();
              context.push('/home/new-quiz/offline?type=${QuizType.randomTasks.name}');
            },
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
}
