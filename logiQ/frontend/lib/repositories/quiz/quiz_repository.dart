import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/models.dart';
import '../../apis/apis.dart';
import '../../data/local_quiz.dart';
import '../../utils/utils.dart';

part 'quiz_repository.g.dart';

@riverpod
QuizRepository quizRepository(Ref ref) {
  return QuizRepositoryImpl(ref.read(quizApiProvider));
}

abstract class QuizRepository {
  
  // 创建新的测验（用于开始新测验）
  Future<Quiz> createNewQuiz({
    required String type,
    String? category,
    String? difficulty,
  });

  // 加载离线测验（用于离线模式）
  Future<Quiz> loadOfflineQuiz();

  // 提交测验
  Future<Quiz> submitQuiz(SubmitQuizRequest request);

  // 获取所有测验历史记录（用于history页面）
  Future<List<Quiz>> getAllQuizHistory();

  // 获取单个测验（用于回顾）
  Future<Quiz> getQuiz(String id);

  // 获取本地测验回顾（用于离线模式下的回顾, 主要是判断题目的正确性）
  Future<Quiz> getLocalQuizReview(Quiz quiz);

}

class QuizRepositoryImpl implements QuizRepository {
  final QuizApi _quizApi;

  QuizRepositoryImpl(this._quizApi);

  
  @override
  Future<Quiz> createNewQuiz({
    required String type,
    String? category,
    String? difficulty,
  }) async {
    return await _quizApi.createNewQuiz(
        type: type,
        category: category,
        difficulty: difficulty,
      );
  }

  @override
  Future<Quiz> loadOfflineQuiz() async {
    return loadQuizFromAssets();
  }

  @override
  Future<Quiz> submitQuiz(SubmitQuizRequest request) async {
    return _quizApi.submitQuiz(request);
  }

  @override
  Future<List<Quiz>> getAllQuizHistory() async {
    return _quizApi.getAllQuizHistory();
  }

  @override
  Future<Quiz> getQuiz(String id) async {
    return _quizApi.getQuiz(id);
  }

  @override
  Future<Quiz> getLocalQuizReview(Quiz quiz) async {
    // 主要是判断题目的正确性
    final questions = quiz.questions;
    var correctCount = 0;
    final newQuestions = List<QuizQuestion>.from(questions);
    for (var i = 0; i < questions.length; i++) {
      final q = questions[i];
      final isCorrect = areListsEqual(q.userAnswerIndex, q.correctAnswerIndex);
      newQuestions[i] = q.copyWith(isCorrect: isCorrect);
      if (isCorrect) {
        correctCount++;
      }
    }
    quiz = quiz.copyWith(correctQuestionsNum: correctCount, questions: newQuestions);
    return quiz;
  }
}