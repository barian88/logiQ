import 'package:flutter/material.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:gap/gap.dart';

class ErrorDistribution extends StatelessWidget {
  const ErrorDistribution({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Error Distribution',
          style: Theme.of(
            context,
          ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
        const Gap(15),
        PieChartListview(),
      ],
    );
  }
}
