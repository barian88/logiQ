import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gap/gap.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/utils/error_handler.dart';

class Performance extends ConsumerWidget {
  const Performance({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {

    final statsState = ref.watch(userStatsNotifierProvider);

    return statsState.when(
      loading: () => const SizedBox(
        height: 120,
        child: Center(child: CircularProgressIndicator()),
      ),
      // 为了了和data状态下的标题对齐，这里也加上标题
      error: (error, stackTrace) => Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Performance',
            style: Theme.of(
              context,
            ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
          ),
          ErrorHandler.buildErrorWidget(
            error,
            onRetry: () => ref.read(userStatsNotifierProvider.notifier).refresh(),
            context: context,
          )
        ],
      ),

      data: (stats) {
        final performance = stats?.performance;
        if (performance == null) {
          return const SizedBox(
            height: 120,
            child: Center(child: Text('No performance data available')),
          );
        }

        final theme = Theme.of(context);

        final divider = Container(
          width: 1,
          height: 35,
          color: theme.colorScheme.primary.withAlpha(128),
          margin: const EdgeInsets.symmetric(horizontal: 10),
        );

        final performanceItems = [
          Column(
            children: [
              Text(
                '${performance.taskNum}',
                style: const TextStyle(fontSize: 19, fontWeight: FontWeight.bold),
              ),
              const Text('Tasks'),
            ],
          ),
          divider,
          Column(
            children: [
              Text(
                '${performance.score}',
                style: const TextStyle(fontSize: 19, fontWeight: FontWeight.bold),
              ),
              const Text('Score'),
            ],
          ),
          divider,
          Column(
            children: [
              Text(
                '${performance.avgTime.toStringAsFixed(1)}s',
                style: const TextStyle(fontSize: 19, fontWeight: FontWeight.bold),
              ),
              const Text('Avg. Time'),
            ],
          ),
        ];

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
              ],
            ),
            const Gap(15),
            LayoutBuilder(
              builder: (context, constraints) {
                final width = constraints.maxWidth;
                return SizedBox(
                  width: width,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: performanceItems,
                  ),
                );
              },
            ),
          ],
        );
      },
    );
  }
}
