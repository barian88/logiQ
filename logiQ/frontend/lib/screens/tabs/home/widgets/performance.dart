import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:frontend/themes/radii.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:gap/gap.dart';

class Performance extends StatelessWidget {
  const Performance({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Performance',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            InkWell(
              onTap: (){
                context.go('/user');
              },
              splashColor: theme.colorScheme.primary.withAlpha(50),
              borderRadius: AppRadii.small,
              child: Text(
                'View All',
                style: theme.textTheme.bodyMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.primary,
                ),
              ),
            ),
          ],
        ),
        const Gap(24),
        AccuracyRateChart()
      ],
    );
  }
}
