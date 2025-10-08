import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/models/models.dart';
import 'package:frontend/themes/themes.dart';
import 'package:gap/gap.dart';
import 'pie_chart_item.dart';
import 'error_distribution_chart.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/utils/error_handler.dart';

class PieChartListview extends ConsumerStatefulWidget {
  const PieChartListview({super.key});

  @override
  ConsumerState<PieChartListview> createState() => _PieChartListviewState();
}

class _PieChartListviewState extends ConsumerState<PieChartListview> {

  final PageController _controller = PageController();
  int _currentPage = 0;

  @override
  Widget build(BuildContext context) {
    final statsState = ref.watch(userStatsNotifierProvider);

    return statsState.when(
      loading: () => const SizedBox(
        height: 200,
        child: Center(child: CircularProgressIndicator()),
      ),
      error: (error, stackTrace) => SizedBox(
        height: 200,
        child: ErrorHandler.buildErrorWidget(
          error,
          onRetry: () => ref.read(userStatsNotifierProvider.notifier).refresh(),
        ),
      ),
      data: (stats) {
        final errorDistribution = stats?.errorDistribution;

        if (errorDistribution == null) {
          return const SizedBox(
            height: 200,
            child: Center(child: Text('No error distribution data available')),
          );
        }

        final colorScheme = Theme.of(context).colorScheme;
        final listByCategory = errorDistribution.dataByCategory;
        final listByDifficulty = errorDistribution.dataByDifficulty;

        // 容器的高度定义为了200，根据此定义半径
        final pieChartItemsByCategory = <PieChartItem>[
          PieChartItem(
            title: QuestionCategoryExtension.getDisplayNameFromName(listByCategory[0].type),
            value: listByCategory[0].value,
            color: colorScheme.primary,
            radius: 90,
          ),
          PieChartItem(
            title: QuestionCategoryExtension.getDisplayNameFromName(listByCategory[1].type),
            value: listByCategory[1].value,
            color: colorScheme.secondary,
            radius: 80,
          ),
          PieChartItem(
            title: QuestionCategoryExtension.getDisplayNameFromName(listByCategory[2].type),
            value: listByCategory[2].value,
            color: colorScheme.blue,
            radius: 75,
          ),
        ];
        final pieChartItemsByDifficulty = <PieChartItem>[
          PieChartItem(
            title: QuestionDifficultyExtension.getDisplayNameFromName(listByDifficulty[0].type),
            value: listByDifficulty[0].value,
            color: colorScheme.green,
            radius: 80,
          ),
          PieChartItem(
            title: QuestionDifficultyExtension.getDisplayNameFromName(listByDifficulty[1].type),
            value: listByDifficulty[1].value,
            color: colorScheme.secondary,
            radius: 90,
          ),
          PieChartItem(
            title: QuestionDifficultyExtension.getDisplayNameFromName(listByDifficulty[2].type),
            value: listByDifficulty[2].value,
            color: colorScheme.red,
            radius: 75,
          ),
        ];

        final pages = <Widget>[
          ErrorDistributionChart(pieChartItems: pieChartItemsByCategory),
          ErrorDistributionChart(pieChartItems: pieChartItemsByDifficulty),
        ];

        return Column(
          children: [
            SizedBox(
              height: 200,
              child: PageView.builder(
                controller: _controller,
                itemCount: pages.length,
                scrollDirection: Axis.horizontal,
                onPageChanged: (index) {
                  setState(() {
                    _currentPage = index;
                  });
                },
                itemBuilder: (context, index) {
                  return Card(elevation: 1, child: pages[index]);
                },
              ),
            ),
            const Gap(10),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(pages.length, (index) {
                final isActive = index == _currentPage;
                return AnimatedContainer(
                  duration: const Duration(milliseconds: 300),
                  margin: const EdgeInsets.symmetric(horizontal: 6),
                  width: isActive ? 10 : 6,
                  height: isActive ? 10 : 6,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color:
                        isActive
                            ? colorScheme.primary
                            : Colors.grey.withAlpha(128),
                  ),
                );
              }),
            ),
          ],
        );
      },
    );
  }
}
