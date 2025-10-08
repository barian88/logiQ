
class Question {

  final String id;
  final String questionText;
  final List<String> options;
  final List<int> correctAnswerIndex;

  final QuestionType type;
  final QuestionCategory category;
  final QuestionDifficulty difficulty;

  Question({
    required this.id,
    required this.questionText,
    required this.options,
    required this.correctAnswerIndex,
    required this.type,
    required this.category,
    required this.difficulty,
  });

  Question copyWith({
    String? id,
    String? questionText,
    List<String>? options,
    List<int>? correctAnswerIndex,
    QuestionType? type,
    QuestionCategory? category,
    QuestionDifficulty? difficulty,
  }) {
    return Question(
      id: id ?? this.id,
      questionText: questionText ?? this.questionText,
      options: options ?? this.options,
      correctAnswerIndex: correctAnswerIndex ?? this.correctAnswerIndex,
      type: type ?? this.type,
      category: category ?? this.category,
      difficulty: difficulty ?? this.difficulty,
    );
  }

  // JSON 序列化支持
  factory Question.fromJson(Map<String, dynamic> json) => Question(
    id: json['_id'] ?? json['id'] ?? '', // 兼容MongoDB的_id字段
    questionText: json['question_text'] ?? '',
    options: List<String>.from(json['options'] ?? []),
    correctAnswerIndex: List<int>.from(json['correct_answer_index'] ?? []),
    type: QuestionType.values.firstWhere(
      (e) => e.name == json['type'],
      orElse: () => QuestionType.singleChoice,
    ),
    category: QuestionCategory.values.firstWhere(
      (e) => e.name == json['category'],
      orElse: () => QuestionCategory.truthTable,
    ),
    difficulty: QuestionDifficulty.values.firstWhere(
      (e) => e.name == json['difficulty'],
      orElse: () => QuestionDifficulty.easy,
    ),
  );

  Map<String, dynamic> toJson() => {
    '_id': id,  // 兼容MongoDB的_id字段
    'question_text': questionText,
    'options': options,
    'correct_answer_index': correctAnswerIndex,
    'type': type.name,
    'category': category.name,
    'difficulty': difficulty.name,
  };

}

enum QuestionType {
  singleChoice,
  multipleChoice,
  trueFalse,
}

extension QuestionTypeExtension on QuestionType {
  String get displayName {
    switch (this) {
      case QuestionType.singleChoice:
        return 'Single Choice';
      case QuestionType.multipleChoice:
        return 'Multiple Choice';
      case QuestionType.trueFalse:
        return 'True / False';
    }
  }

  static String getDisplayNameFromName(String name) {
    try {
      final type = QuestionType.values.firstWhere((e) => e.name == name);
      return type.displayName;
    } catch (e) {
      return name; // 如果找不到匹配的枚举值，返回原始名称
    }
  }

  static String getNameFromDisplayName(String displayName) {
    try {
      final type = QuestionType.values.firstWhere((e) => e.displayName == displayName);
      return type.name;
    } catch (e) {
      return displayName; // 如果找不到匹配的枚举值，返回原始显示名称
    }
  }
}

enum QuestionCategory {
  truthTable,
  equivalence,
  inference,
}

extension QuestionCategoryExtension on QuestionCategory {
  String get displayName {
    switch (this) {
      case QuestionCategory.truthTable:
        return 'Truth Table';
      case QuestionCategory.equivalence:
        return 'Equivalence';
      case QuestionCategory.inference:
        return 'Inference';
    }
  }

  static String getDisplayNameFromName(String name) {
    try {
      final category = QuestionCategory.values.firstWhere((e) => e.name == name);
      return category.displayName;
    } catch (e) {
      return name; // 如果找不到匹配的枚举值，返回原始名称
    }
  }

  static String getNameFromDisplayName(String displayName) {
    try {
      final category = QuestionCategory.values.firstWhere((e) => e.displayName == displayName);
      return category.name;
    } catch (e) {
      return displayName; // 如果找不到匹配的枚举值，返回原始显示名称
    }
  }
}

enum QuestionDifficulty {
  easy,
  medium,
  hard,
}

extension QuestionDifficultyExtension on QuestionDifficulty {
  String get displayName {
    switch (this) {
      case QuestionDifficulty.easy:
        return 'Easy';
      case QuestionDifficulty.medium:
        return 'Medium';
      case QuestionDifficulty.hard:
        return 'Hard';
    }
  }

  static String getDisplayNameFromName(String name) {
    try {
      final difficulty = QuestionDifficulty.values.firstWhere((e) => e.name == name);
      return difficulty.displayName;
    } catch (e) {
      return name; // 如果找不到匹配的枚举值，返回原始名称
    }
  }

  static String getNameFromDisplayName(String displayName) {
    try {
      final difficulty = QuestionDifficulty.values.firstWhere((e) => e.displayName == displayName);
      return difficulty.name;
    } catch (e) {
      return displayName; // 如果找不到匹配的枚举值，返回原始显示名称
    }
  }
}
