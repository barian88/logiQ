import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'pie_chart_item.dart';
import 'pie_chart_legend.dart';

class ErrorDistributionChart extends StatefulWidget {

  final List<PieChartItem> pieChartItems;

  const ErrorDistributionChart({super.key, required this.pieChartItems});

  @override
  State<ErrorDistributionChart> createState() => _ErrorDistributionChartState();
}

class _ErrorDistributionChartState extends State<ErrorDistributionChart> {
  int touchedIndex = -1;

  @override
  Widget build(BuildContext context) {

    return SizedBox(
      height: 200,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          PieChartLegend(pieChartItems: widget.pieChartItems),
          SizedBox(
            width: 200,
            child: PieChart(
              PieChartData(
                pieTouchData: PieTouchData(touchCallback: touchedCallback),
                startDegreeOffset: 180,
                sectionsSpace: 2,
                centerSpaceRadius: 0,
                sections: showingSections(widget.pieChartItems),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void touchedCallback(FlTouchEvent event, pieTouchResponse) {
    setState(() {
      if (!event.isInterestedForInteractions ||
          pieTouchResponse == null ||
          pieTouchResponse.touchedSection == null) {
        touchedIndex = -1;
        return;
      }
      touchedIndex = pieTouchResponse.touchedSection!.touchedSectionIndex;
    });
  }

  List<PieChartSectionData> showingSections(List<PieChartItem> pieChartItems) {
    return List.generate(pieChartItems.length, (i) {
      final isTouched = i == touchedIndex;
      final item = pieChartItems[i];
      final titlePositionPercentageOffset = 0.55;
      switch (i) {
        case 0:
          return PieChartSectionData(
            color: item.color,
            value: item.value,
            title: isTouched ? '${item.value}%' : '',
            radius: item.radius,
            titlePositionPercentageOffset: titlePositionPercentageOffset,
          );
        case 1:
          return PieChartSectionData(
            color: item.color,
            value: item.value,
            title: isTouched ? '${item.value}%' : '',
            radius: item.radius,
            titlePositionPercentageOffset: titlePositionPercentageOffset,
          );
        case 2:
          return PieChartSectionData(
            color: item.color,
            value: item.value,
            title: isTouched ? '${item.value}%' : '',
            radius: item.radius,
            titlePositionPercentageOffset: titlePositionPercentageOffset,
          );
        default:
          throw Error();
      }
    });
  }
}
