import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../themes/colors.dart';
import 'package:frontend/pods/pods.dart';
import 'package:frontend/utils/error_handler.dart';

class AccuracyRateChart extends ConsumerStatefulWidget {
  const AccuracyRateChart({super.key});

  @override
  ConsumerState<AccuracyRateChart> createState() => _AccuracyRateChartState();
}

class _AccuracyRateChartState extends ConsumerState<AccuracyRateChart> {
  
  List<FlSpot> _generateSpots(List<dynamic> accuracyData) {
    if (accuracyData.isEmpty) return [];
    
    List<FlSpot> spots = [];
    for (int i = 0; i < accuracyData.length; i++) {
      final item = accuracyData[i];
      spots.add(FlSpot(i.toDouble(), item.value)); // 直接使用百分比(0-1)范围值
    }
    return spots;
  }

  @override
  Widget build(BuildContext context) {
    final statsState = ref.watch(userStatsNotifierProvider);

    return statsState.when(
      loading: () => const SizedBox(
        height: 270,
        child: Center(child: CircularProgressIndicator()),
      ),
      error: (error, stackTrace) => SizedBox(
        height: 270,
        child: ErrorHandler.buildErrorWidget(
          error,
          onRetry: () => ref
              .read(userStatsNotifierProvider.notifier)
              .refresh(),
        ),
      ),
      data: (stats) {
        final accuracyRate = stats?.accuracyRate;

        if (accuracyRate == null || accuracyRate.data.isEmpty) {
          return const SizedBox(
            height: 270,
            child: Center(child: Text('No accuracy rate data available')),
          );
        }

        return SizedBox(
          width: double.infinity,
          height: 270,
          child: LineChart(mainData(accuracyRate.data)),
        );
      },
    );
  }

  Widget bottomTitleWidgets(double value, TitleMeta meta, List<dynamic> accuracyData) {
    final style = Theme.of(context).textTheme.bodyMedium?.copyWith(
      fontWeight: FontWeight.bold,
      color: AppColors.grey1,
    );
    
    final index = value.toInt();
    final isLast = value == meta.max;
    
    Widget text;
    if (index >= 0 && index < accuracyData.length) {
      if (isLast) {
        text = Text('Latest', style: style);
      } else if (index == 0) {
        text = Text('', style: style);
      } else {
        text = Text('', style: style);
      }
    } else {
      text = Text('', style: style);
    }

    return SideTitleWidget(
      meta: meta,
      child: Padding(
        padding: EdgeInsets.only(right: isLast ? 48 : 0), // 给最后一个标签加右 padding
        child: text,
      ),
    );
  }

  Widget leftTitleWidgets(double value, TitleMeta meta) {
    final style = Theme.of(context).textTheme.bodyMedium?.copyWith(
      fontWeight: FontWeight.bold,
      color: AppColors.grey1,
    );
    String text;
    if (value == 0.0) {
      text = '0%';
    } else if (value == 0.25) {
      text = '25%';
    } else if (value == 0.5) {
      text = '50%';
    } else if (value == 0.75) {
      text = '75%';
    } else if (value == 1.0) {
      text = '100%';
    } else {
      return Container();
    }

    return Text(text, style: style, textAlign: TextAlign.left);
  }

  LineChartData mainData(List<dynamic> accuracyData) {
    final spots = _generateSpots(accuracyData);
    final maxX = (accuracyData.length - 1).toDouble();
    final maxY = 1.0;
    
    return LineChartData(
      // baselineY: 0,
      // baselineX: 0,
      gridData: FlGridData(
        show: true,
        drawVerticalLine: false,
        horizontalInterval: 0.25,
        getDrawingHorizontalLine: (value) {
          return FlLine(
            color: Theme.of(context).colorScheme.onSurfaceVariant.withValues(alpha: 0.3),
            strokeWidth: 1,
          );
        },
      ),
      titlesData: FlTitlesData(
        show: true,
        rightTitles: const AxisTitles(
          sideTitles: SideTitles(showTitles: false),
        ),
        topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
        bottomTitles: AxisTitles(
          sideTitles: SideTitles(
            showTitles: true,
            reservedSize: 30,
            interval: 1,
            getTitlesWidget: (value, meta) => bottomTitleWidgets(value, meta, accuracyData),
          ),
        ),
        leftTitles: AxisTitles(
          sideTitles: SideTitles(
            showTitles: true,
            interval: 0.25,
            getTitlesWidget: leftTitleWidgets,
            reservedSize: 42,
          ),
        ),
      ),
      borderData: FlBorderData(
        show: true,
        border: Border.symmetric(
          horizontal: BorderSide(color: AppColors.grey1, width: 1.5),
        ),
      ),
      minX: 0,
      maxX: maxX,
      minY: 0,
      maxY: maxY,
      lineBarsData: [
        LineChartBarData(
          spots: spots,
          isCurved: true,
          preventCurveOverShooting: true, // 防止曲线过冲，不会超出0-1的范围
          color: AppColors.bluePurple,
          barWidth: 5,
          isStrokeCapRound: true,
          dotData: FlDotData(
            show: true,
            getDotPainter: (_, _, _, index) {
              final isLast = (index == spots.length - 1);

              if (isLast) {
                return FlDotCirclePainter(
                  radius: 3,
                  color: Theme.of(context).colorScheme.secondary,
                  strokeWidth: 3.5,
                  strokeColor: AppColors.bluePurple,
                );
              } else {
                return FlDotCirclePainter(radius: 0);
              }
            },
          ),
          belowBarData: BarAreaData(
            show: true,
            gradient: LinearGradient(
              colors: [
                AppColors.bluePurple.withValues(alpha: 0.5),
                AppColors.bluePurple.withValues(alpha: 0.2),
                AppColors.bluePurple.withValues(alpha: 0.0),
              ],
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
            ),
          ),
        ),
      ],
    );
  }
}
