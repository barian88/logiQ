import 'question.dart';

class QuizQuestion {
  final Question question;
  final List<int> userAnswerIndex;
  final bool isCorrect;

  const QuizQuestion({
    required this.question,
    this.userAnswerIndex = const [],
    this.isCorrect = false,
  });

  QuizQuestion copyWith({
    Question? question,
    List<int>? userAnswerIndex,
    bool? isCorrect,
  }) {
    return QuizQuestion(
      question: question ?? this.question,
      userAnswerIndex: userAnswerIndex ?? this.userAnswerIndex,
      isCorrect: isCorrect ?? this.isCorrect,
    );
  }

  // JSON 序列化支持
  factory QuizQuestion.fromJson(Map<String, dynamic> json) => QuizQuestion(
    question: Question.fromJson(json['question'] ?? {}),
    userAnswerIndex: List<int>.from(json['user_answer_index'] ?? []),
    isCorrect: json['is_correct'] ?? false,
  );

  Map<String, dynamic> toJson() => {
    'question': question.toJson(),
    'user_answer_index': userAnswerIndex,
    'is_correct': isCorrect,
  };

  // 便利方法
  bool get isAnswered => userAnswerIndex.isNotEmpty;

  // 获取问题ID的便利方法
  String get questionId => question.id;
  
  // 获取问题文本的便利方法
  String get questionText => question.questionText;
  
  // 获取选项的便利方法
  List<String> get options => question.options;
  
  // 获取正确答案的便利方法
  List<int> get correctAnswerIndex => question.correctAnswerIndex;
}