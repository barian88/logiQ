import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/utils/error_handler.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/pods/pods.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import 'widgets/widgets.dart';
import 'package:frontend/models/models.dart';

class NewQuiz extends ConsumerStatefulWidget {
  const NewQuiz({super.key, required this.mode});

  final String mode;

  @override
  ConsumerState<NewQuiz> createState() => _NewQuizState();
}

class _NewQuizState extends ConsumerState<NewQuiz> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      _initQuiz();
    });
  }

  _initQuiz() async {
    final quizNotifier = ref.read(quizNotifierProvider.notifier);
    quizNotifier.setMode(widget.mode);
    final (type, category, difficulty) = _parseQueryParameters();
    await quizNotifier.loadNewQuiz(
      type: type,
      category: category,
      difficulty: difficulty,
    );

  }

  @override
  Widget build(BuildContext context) {
    final quizState = ref.watch(quizNotifierProvider);

    return quizState.when(
      loading:
          () =>
              _buildQuizScreen(context, const Center(child: CircularProgressIndicator())),
      error:
          (error, stackTrace) => _buildQuizScreen(context,
              Padding(
                padding: const EdgeInsets.fromLTRB(0, 0, 0, 200),
                child: ErrorHandler.buildErrorWidget(error,
                onRetry: _initQuiz, context: context),
              )),
      data: (state) {
        //  api返回成功，但server没有返回题目
        if (state.quiz.questions.isEmpty) {
          return _buildQuizScreen(context, Center(child: Text('No quiz data available')));
        }
        // 正常加载题目
        return _buildQuizScreen(context, Column(
          children: [
            QuizProgress(),
            const Gap(30),
            Expanded(flex: 4, child: QuestionArea()),
            Divider(height: 50),
            Expanded(
              flex: 5,
              child: Column(
                children: [OptionArea(), Spacer(), OperationArea(), Gap(35)],
              ),
            ),
          ],
        ));
      },
    );
  }

  Widget _buildQuizScreen(BuildContext context, Widget content) {
    return Scaffold(
      appBar: AppBar(
        title: QuizTimer(),
        centerTitle: true,
        actions: [ThemeModeSwitch(), const Gap(16)],
      ),
      body: BaseContainer(
        isScrollable: false,
        child: content,
      ),
    );
  }

  (String, String?, String?) _parseQueryParameters() {
    // 获取查询参数
    final state = GoRouterState.of(context);
    // 传入的typeParam是name
    final typeParam = state.uri.queryParameters['type'];
    // 传入的categoryParam和difficultyParam是displayName
    final categoryParam = state.uri.queryParameters['category'];
    final difficultyParam = state.uri.queryParameters['difficulty'];

    // 得到对应的enum.name
    final String type = typeParam ?? QuizType.randomTasks.name;
    final String? category;
    final String? difficulty;

    category =
        (categoryParam != null)
            ? QuestionCategoryExtension.getNameFromDisplayName(categoryParam)
            : null;
    difficulty =
        (difficultyParam != null)
            ? QuestionDifficultyExtension.getNameFromDisplayName(
              difficultyParam,
            )
            : null;

    return (type, category, difficulty);
  }
}
