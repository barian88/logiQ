import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:frontend/models/models.dart';
import 'package:frontend/repositories/repositories.dart';
import 'package:frontend/apis/apis.dart';
import 'dart:async';

part 'quiz_pod.g.dart';

// 保持provider活跃，防止状态丢失。
// 主要是离线模式下，用户切换页面时状态不会丢失，从而能够实现Review
// 但是在线模式下，用户会重新加载测验，可以覆盖掉之前的状态
@Riverpod(keepAlive: true)
class QuizNotifier extends _$QuizNotifier {
  Timer? _timer;

  @override
  AsyncValue<QuizState> build() {
    ref.keepAlive();
    // 在build方法中注册dispose回调
    ref.onDispose(() {
      _timer?.cancel();
    });

    return const AsyncValue.data(QuizState.initial());
  }

// 设置测验模式 - online 或 offline
  void setMode(String mode) {
    final quizMode = mode == 'offline' ? QuizMode.offline : QuizMode.online;
    final currentState = state.asData?.value ?? const QuizState.initial();
    state = AsyncValue.data(currentState.copyWith(mode: quizMode));
  }

  // 获取新的测验 - mode normal 和 mode offline
  Future<void> loadNewQuiz({required String type, String? category, String? difficulty}) async {

    try {
      final repository = ref.read(quizRepositoryProvider);
      final currentMode = state.value?.mode ?? QuizMode.online;
      state = const AsyncValue.loading(); // 要在获取mode之后再设置loading状态，否则会覆盖mode
      Quiz quiz;
      // 如果是online模式，从服务器获取题目
      if (currentMode == QuizMode.online) {
        quiz = await repository.createNewQuiz(
          type: type,
          category: category,
          difficulty: difficulty,
        );
      } else {
        quiz = await repository.loadOfflineQuiz();
      }

      // 初始化QuizState
      state = AsyncValue.data(QuizState.loaded(quiz, mode: currentMode));
      // 开始计时
      startTimer();
    } catch (error, stackTrace) {
      state = AsyncValue.error(error, stackTrace);
    }
  }

  // 获取测验回顾
  Future<void> loadQuizReview(String quizId) async {
    final currentMode = state.value?.mode ?? QuizMode.online;

    try {
      final repository = ref.read(quizRepositoryProvider);
      // online模式下从服务器获取测验
      if(currentMode == QuizMode.online) {
        // 确认是在线模式后再设置loading状态，否则会覆盖state的数据，缓存丢失，无法实现离线模式的Review
        state = const AsyncValue.loading();

        final quiz = await repository.getQuiz(quizId);
        state = AsyncValue.data(QuizState.review(quiz, mode: currentMode));
      }
      // offline模式下从本地获取测验，不使用loading刷新状态：保留数据和mode
      else {
        Quiz quiz = state.value?.quiz ?? await repository.loadOfflineQuiz();
        // 拿到isCorrect的数据
        quiz = await repository.getLocalQuizReview(quiz);
        state = AsyncValue.data(QuizState.review(quiz, mode: currentMode));
      }

    } catch (error, stackTrace) {
      state = AsyncValue.error(error, stackTrace);
    }
  }

  // 确认提交测验
  Future<SubmitQuizResult> submitQuiz() async {
    stopTimer();
    try{
      final currentMode = state.value?.mode ?? QuizMode.online;
      final quiz = state.value?.quiz;
      // 在线模式下提交测验
      if(currentMode == QuizMode.online){
        // 构建quizRequest
        SubmitQuizRequest quizRequest;
        if (quiz == null) {
          return SubmitQuizResult.error('Quiz data is not available');
        } else {
          quizRequest = SubmitQuizRequest(type: quiz.type.name, questions: quiz.questions, completionTime: quiz.completionTime);
        }
        final quizRepository = ref.read(quizRepositoryProvider);
        final quizRes = await quizRepository.submitQuiz(quizRequest);
        return SubmitQuizResult.success(quizRes.id);
      }
      // 离线模式下不提交，直接返回成功。判断正确性在loadReview中进行
      return SubmitQuizResult.success(quiz?.id ?? '');

    }catch(e){
      return SubmitQuizResult.error('Failed to submit quiz: ${e.toString()}');
    }
  }

  // 下一题
  void nextQuestion() {
    state.whenData((quizState) {
      if (quizState.currentQuestionIndex <
          quizState.quiz.questions.length - 1) {
        state = AsyncValue.data(
          quizState.copyWith(
            currentQuestionIndex: quizState.currentQuestionIndex + 1,
          ),
        );
      }
    });
  }

  // 上一题
  void previousQuestion() {
    state.whenData((quizState) {
      if (quizState.currentQuestionIndex > 0) {
        state = AsyncValue.data(
          quizState.copyWith(
            currentQuestionIndex: quizState.currentQuestionIndex - 1,
          ),
        );
      }
    });
  }

  // 设置用户答案
  void setUserAnswerIndex(int optionIndex) {
    state.whenData((quizState) {
      final currentQuizQuestion =
          quizState.quiz.questions[quizState.currentQuestionIndex];
      List<int> userAnswerIndex = List.from(
        currentQuizQuestion.userAnswerIndex,
      );

      if (currentQuizQuestion.question.type == QuestionType.singleChoice ||
          currentQuizQuestion.question.type == QuestionType.trueFalse) {
        userAnswerIndex = [optionIndex];
      } else if (currentQuizQuestion.question.type ==
          QuestionType.multipleChoice) {
        if (currentQuizQuestion.userAnswerIndex.contains(optionIndex)) {
          userAnswerIndex.remove(optionIndex);
        } else {
          userAnswerIndex.add(optionIndex);
        }
      }

      final updatedQuizQuestion = currentQuizQuestion.copyWith(
        userAnswerIndex: userAnswerIndex,
      );

      final updatedQuestions = [...quizState.quiz.questions];
      updatedQuestions[quizState.currentQuestionIndex] = updatedQuizQuestion;

      final updatedQuiz = quizState.quiz.copyWith(questions: updatedQuestions);

      state = AsyncValue.data(quizState.copyWith(quiz: updatedQuiz));
    });
  }

  // 开始计时器
  void startTimer() {
    _timer?.cancel();

    state.whenData((quizState) {
      final updatedQuiz = quizState.quiz.copyWith(completionTime: 0);
      state = AsyncValue.data(quizState.copyWith(quiz: updatedQuiz));

      _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
        state.whenData((currentState) {
          final updatedQuiz = currentState.quiz.copyWith(
            completionTime: currentState.quiz.completionTime + 1,
          );
          state = AsyncValue.data(currentState.copyWith(quiz: updatedQuiz));
        });
      });
    });
  }

  // 停止计时器
  void stopTimer() {
    _timer?.cancel();
  }


  // 重置测验
  void resetQuiz() {
    _timer?.cancel();
    state = const AsyncValue.data(QuizState.initial());
  }
}

// 测验状态类
class QuizState {
  final Quiz quiz;
  final int currentQuestionIndex;
  final QuizStatus status;
  final QuizMode mode;

  const QuizState({
    required this.quiz,
    this.currentQuestionIndex = 0,
    this.status = QuizStatus.initial,
    this.mode = QuizMode.online,
  });

  const QuizState.initial()
    : this(
        quiz: const Quiz(id: '', type: QuizType.randomTasks, questions: []),
        status: QuizStatus.initial,
        mode: QuizMode.online,
      );

  const QuizState.loaded(Quiz quiz, {QuizMode mode = QuizMode.online})
    : this(quiz: quiz, status: QuizStatus.active, mode: mode);

  const QuizState.review(Quiz quiz, {QuizMode mode = QuizMode.online})
    : this(quiz: quiz, status: QuizStatus.review, mode: mode);

  QuizState copyWith({
    Quiz? quiz,
    int? currentQuestionIndex,
    QuizStatus? status,
    QuizMode? mode,
  }) {
    return QuizState(
      quiz: quiz ?? this.quiz,
      currentQuestionIndex: currentQuestionIndex ?? this.currentQuestionIndex,
      status: status ?? this.status,
      mode: mode ?? this.mode,
    );
  }
}

enum QuizStatus { initial, active, review }

enum QuizMode { online, offline }

// submit结果类
class SubmitQuizResult {
  final bool isSuccess;
  final String? quizId;
  final String? errorMessage;

  const SubmitQuizResult({required this.isSuccess, this.quizId, this.errorMessage});

  factory SubmitQuizResult.success(String quizId) {
    return SubmitQuizResult(isSuccess: true, quizId: quizId);
  }

  factory SubmitQuizResult.error(String message) {
    return SubmitQuizResult(isSuccess: false, errorMessage: message);
  }
}
