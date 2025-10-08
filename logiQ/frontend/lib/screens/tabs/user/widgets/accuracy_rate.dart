import 'package:flutter/material.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:gap/gap.dart';

class AccuracyRate extends StatelessWidget {
  const AccuracyRate({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Accuracy Rate',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
        const Gap(30),
        AccuracyRateChart()
      ],
    );
  }
}
