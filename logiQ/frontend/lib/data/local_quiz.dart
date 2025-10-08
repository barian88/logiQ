import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart' show rootBundle;
import '../models/models.dart';

Future<Quiz> loadQuizFromAssets() async {
  final path = 'assets/questions/questions.json';

  // 读取文件内容
  final String jsonString = await rootBundle.loadString(path);

  // 解析 JSON 数据
  final List<dynamic> questionList = jsonDecode(jsonString) as List<dynamic>;

  final List<Question> questions = questionList.map((raw) {
    final map = Map<String, dynamic>.from(raw as Map<String, dynamic>);

    // 处理 MongoDB 风格的 _id 字段
    final dynamic rawId = map['_id'];
    if (rawId is Map && rawId['\$oid'] is String) {
      final oid = rawId['\$oid'] as String;
      map['id'] = oid;
      map['_id'] = oid;
    } else if (rawId is String) {
      map['id'] = rawId;
      map['_id'] = rawId;
    } else {
      map['_id'] = map['id'];
    }

    return Question.fromJson(map);
  }).toList();

  final List<QuizQuestion> quizQuestions = questions
      .map((question) => QuizQuestion(question: question))
      .toList();
  debugPrint('Loaded ${quizQuestions.length} questions from $path');

  return Quiz(
    id: 'local_quiz',
    questions: quizQuestions,
    type: QuizType.randomTasks,
  );
}
