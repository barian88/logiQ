import 'dart:convert';

import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/models.dart';
import '../api_service.dart';

part 'quiz_api.g.dart';

@riverpod
QuizApi quizApi(Ref ref) {
  return QuizApi(ref.read(apiServiceProvider));
}

class QuizApi {
  final ApiService _apiService;

  QuizApi(this._apiService);

  // 创建新的测验（用于开始新测验）
  Future<Quiz> createNewQuiz({
    required String type,
    String? category,
    String? difficulty,
  }) async {
    final body = <String, String>{};

    body['type'] = type;
    if (category != null) body['category'] = category;
    if (difficulty != null) body['difficulty'] = difficulty;

    final response = await _apiService.post('/quiz/new', body: body);
    return Quiz.fromJson(response);
  }

  // 提交测验
  Future<Quiz> submitQuiz(SubmitQuizRequest request) async {
    final response = await _apiService.post('/quiz/submit', body: request.toJson());
    return Quiz.fromJson(response);
  }

  // 获取所有测验历史记录（用于history页面）
  Future<List<Quiz>> getAllQuizHistory() async {
    final body = await _apiService.getString('/quiz/history');
    final decoded = jsonDecode(body);

    final dynamic data =
    decoded is Map<String, dynamic> ? decoded['data'] ?? [] : decoded;

    if (data is! List) {
      return [];
    }

    return data
        .whereType<Map<String, dynamic>>()
        .map(Quiz.fromJson)
        .toList();
  }

  // 获取单个测验（用于回顾）
  Future<Quiz> getQuiz(String id) async {
    final response = await _apiService.get('/quiz/$id');
    return Quiz.fromJson(response);
  }

}

// 请求和响应模型
class SubmitQuizRequest {

  final String type;
  final List<QuizQuestion> questions;
  final int completionTime;

  SubmitQuizRequest({
    required this.type,
    required this.questions,
    required this.completionTime,
  });

  Map<String, dynamic> toJson() => {
    'type': type,
    'questions': questions.map((q) => q.toJson()).toList(),
    'completion_time': completionTime,
  };
}

