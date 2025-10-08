import 'package:flutter/material.dart';
import 'package:gap/gap.dart';

import 'pie_chart_item.dart';

class PieChartLegend extends StatelessWidget {
  const PieChartLegend({super.key, required this.pieChartItems});

  final List<PieChartItem> pieChartItems;

  @override
  Widget build(BuildContext context) {

    final theme = Theme.of(context);

    final List<Widget> legendList = [];

    for (int i=0; i < pieChartItems.length; i++) {
      final item = pieChartItems[i];
      final legend = Row(
        children: [
          Container(
            width: 15,
            height: 15,
            decoration: BoxDecoration(
              color: item.color,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 8),
          Text(item.title, style: theme.textTheme.bodySmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant
          ),),
        ],
      );
      legendList.add(legend);
      if (i != pieChartItems.length - 1) {
        legendList.add(Gap(15)); // 添加间隔
      }
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisAlignment: MainAxisAlignment.center,
      children: legendList,
    );
  }

}
