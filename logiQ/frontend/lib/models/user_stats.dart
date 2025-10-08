class UserStats {
  final String id;
  final Performance performance;
  final AccuracyRate accuracyRate;
  final ErrorDistribution errorDistribution;

  const UserStats({
    required this.id,
    required this.performance,
    required this.accuracyRate,
    required this.errorDistribution,
  });

  factory UserStats.fromJson(Map<String, dynamic> json) {
    return UserStats(
      id: json['_id'],
      performance: Performance.fromJson(json['performance']),
      accuracyRate: AccuracyRate.fromJson(json['accuracy_rate']),
      errorDistribution: ErrorDistribution.fromJson(json['error_distribution']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      '_id': id,
      'performance': performance.toJson(),
      'accuracy_rate': accuracyRate.toJson(),
      'error_distribution': errorDistribution.toJson(),
    };
  }
}

class Performance {
  final int taskNum;
  final int score;
  final double avgTime;

  const Performance({
    required this.taskNum,
    required this.score,
    required this.avgTime,
  });

  factory Performance.fromJson(Map<String, dynamic> json) {
    return Performance(
      taskNum: json['task_num'],
      score: json['score'],
      avgTime: (json['avg_time'] as num).toDouble(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'task_num': taskNum,
      'score': score,
      'avg_time': avgTime,
    };
  }
}

class AccuracyRate {
  final List<AccuracyRateItem> data;

  const AccuracyRate({
    required this.data,
  });

  factory AccuracyRate.fromJson(Map<String, dynamic> json) {
    return AccuracyRate(
      data: (json['data'] as List)
          .map((item) => AccuracyRateItem.fromJson(item))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'data': data.map((item) => item.toJson()).toList(),
    };
  }
}

class AccuracyRateItem {
  final DateTime date;
  final double value;

  const AccuracyRateItem({
    required this.date,
    required this.value,
  });

  factory AccuracyRateItem.fromJson(Map<String, dynamic> json) {
    return AccuracyRateItem(
      date: DateTime.parse(json['date']),
      value: (json['value'] as num).toDouble(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'date': date.toIso8601String(),
      'value': value,
    };
  }
}

class ErrorDistribution {
  final List<ErrorDistributionItem> dataByCategory;
  final List<ErrorDistributionItem> dataByDifficulty;

  const ErrorDistribution({
    required this.dataByCategory,
    required this.dataByDifficulty,
  });

  factory ErrorDistribution.fromJson(Map<String, dynamic> json) {
    return ErrorDistribution(
      dataByCategory: (json['data_by_category'] as List)
          .map((item) => ErrorDistributionItem.fromJson(item))
          .toList(),
      dataByDifficulty: (json['data_by_difficulty'] as List)
          .map((item) => ErrorDistributionItem.fromJson(item))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'data_by_category': dataByCategory.map((item) => item.toJson()).toList(),
      'data_by_difficulty': dataByDifficulty.map((item) => item.toJson()).toList(),
    };
  }
}

class ErrorDistributionItem {
  final String type;
  final double value;

  const ErrorDistributionItem({
    required this.type,
    required this.value,
  });

  factory ErrorDistributionItem.fromJson(Map<String, dynamic> json) {
    return ErrorDistributionItem(
      type: json['type'],
      value: (json['value'] as num).toDouble(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'type': type,
      'value': value,
    };
  }
}