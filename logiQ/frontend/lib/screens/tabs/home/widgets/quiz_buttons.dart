import 'package:flutter/material.dart';
import 'package:gap/gap.dart';
import 'quiz_button_container.dart';

class QuizButtons extends StatelessWidget {
  const QuizButtons({super.key});

  static const buttonList = [
    {'title': 'Random Tasks', 'emoji': '🎲'},
    {'title': 'Topic Practice', 'emoji': '🎯'},
    {'title': 'By Difficulty', 'emoji': '📚'},
    {'title': 'Custom Quiz	', 'emoji': '✏️'},
    {'title': 'Offline Practice', 'emoji': '🗂️'},
  ];

  @override
  Widget build(BuildContext context) {
    List<Widget> quizButtons = List.generate(buttonList.length, (index) {
      final item = buttonList[index];
      return Row(
        children: [
          QuizButtonContainer(
            title: item['title']!,
            emoji: item['emoji']!,
            index: index,
          ),
          const Gap(10),
        ],
      );
    });

    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Quiz',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            Icon(
              Icons.arrow_forward_ios,
              size: 16,
              color: Theme.of(context).colorScheme.primary,
            ),
          ],
        ),
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 10.0),
            child: Row(children: quizButtons),
          ),
        ),
      ],
    );
  }
}
