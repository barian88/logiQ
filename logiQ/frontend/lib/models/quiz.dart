import 'quiz_question.dart';

enum QuizType{
  randomTasks,
  topicPractice,
  byDifficulty,
  customQuiz,
}

extension QuizTypeExtension on QuizType {
  String get displayName {
    switch (this) {
      case QuizType.randomTasks:
        return 'Random Tasks';
      case QuizType.topicPractice:
        return 'Topic Practice';
      case QuizType.byDifficulty:
        return 'By Difficulty';
      case QuizType.customQuiz:
        return 'Custom Quiz';
    }
  }

  static String getDisplayNameFromName(String name) {
    try {
      final type = QuizType.values.firstWhere((e) => e.name == name);
      return type.displayName;
    } catch (e) {
      return name; // 如果找不到匹配的枚举值，返回原始名称
    }
  }

  static String getNameFromDisplayName(String displayName) {
    try {
      final type = QuizType.values.firstWhere((e) => e.displayName == displayName);
      return type.name;
    } catch (e) {
      return displayName; // 如果找不到匹配的枚举值，返回原始显示名称
    }
  }
}

class Quiz{

  final String id;
  final List<QuizQuestion> questions;
  final QuizType type;
  final int correctQuestionsNum;
  final int completionTime;
  final DateTime? completedAt;


  const Quiz({
    required this.id,
    required this.questions,
    required this.type,
    this.correctQuestionsNum = 0,
    this.completionTime = 0,
    this.completedAt,
  });

  Quiz copyWith({
    String? id,
    List<QuizQuestion>? questions,
    QuizType? type,
    int? correctQuestionsNum,
    int? completionTime,
    DateTime? completedAt,
  }) {
    return Quiz(
      id: id ?? this.id,
      questions: questions ?? this.questions,
      type: type ?? this.type,
      correctQuestionsNum: correctQuestionsNum ?? this.correctQuestionsNum,
      completionTime: completionTime ?? this.completionTime,
      completedAt: completedAt ?? this.completedAt,
    );
  }

  // 便利方法
  int get totalQuestionsNum => questions.length;
  
  // JSON 序列化支持
  factory Quiz.fromJson(Map<String, dynamic> json) => Quiz(
    id: json['_id'] ?? json['id'] ?? '', // 兼容MongoDB的_id字段
    questions: (json['questions'] as List<dynamic>?)
        ?.map((item) => QuizQuestion.fromJson(item))
        .toList() ?? [],
    type: QuizType.values.firstWhere(
      (e) => e.name == json['type'],
      orElse: () => QuizType.randomTasks,
    ),
    correctQuestionsNum: json['correct_questions_num'] ?? 0,
    completionTime: json['completion_time'] ?? 0,
    completedAt: json['completed_at'] != null 
        ? DateTime.parse(json['completed_at']) 
        : null,
  );

  Map<String, dynamic> toJson() => {
    '_id': id,
    'questions': questions.map((q) => q.toJson()).toList(),
    'type': type.name,
    'correct_questions_num': correctQuestionsNum,
    'completion_time': completionTime,
    if (completedAt != null) 'completed_at': completedAt!.toIso8601String(),
  };

}