import 'package:frontend/pods/user/user_pod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/material.dart';
import '../../screens/screens.dart';

part 'route_pod.g.dart';

@riverpod
class RouteNotifier extends _$RouteNotifier {
  final GlobalKey<NavigatorState> _rootNavigatorKey =
      GlobalKey<NavigatorState>();
  final GlobalKey<NavigatorState> _shellNavigatorKey =
      GlobalKey<NavigatorState>();

  GoRouter? _router;

  @override
  GoRouter build() {
    // 监听用户状态，一有变化就通知 router 刷新一次
    ref.listen<AsyncValue<UserState>>(userNotifierProvider, (_, __) {
      _router?.refresh();
    });

    return _router ??= GoRouter(
      navigatorKey: _rootNavigatorKey,
      initialLocation: '/home',
      // initialLocation: '/guide-detail/1',
      redirect: (context, state) {
        final userAsync = ref.read(userNotifierProvider);

        // 还在加载 SharedPreferences，先不跳转
        if (userAsync.isLoading) return null;

        final userState = userAsync.valueOrNull;
        final isLoggedIn = userState?.isLoggedIn ?? false;

        // 未登录, 除了登录/注册等页面都重定向到登录
        final location = state.matchedLocation;
        final isAuthRoute = location == '/login' ||
            location.startsWith('/register') ||
            location.startsWith('/verification') ||
            location.startsWith('/change-password');

        if (!isLoggedIn && !isAuthRoute) return '/login';

        return null;
      },
      routes: [
        ShellRoute(
          navigatorKey: _shellNavigatorKey,
          builder: (context, state, child) {
            return TabScaffold(child: child);
          },
          routes: [
            GoRoute(
              path: '/home',
              pageBuilder:
                  (context, state) => NoTransitionPage(child: const Home()),
              routes: [
                GoRoute(
                  parentNavigatorKey: _rootNavigatorKey,
                  path: '/new-quiz/:mode',
                  builder: (context, state) {
                    final mode = state.pathParameters['mode'] ?? 'normal';
                    return NewQuiz(mode: mode);
                  },
                ),
              ],
            ),
            GoRoute(
              path: '/history',
              pageBuilder:
                  (context, state) => NoTransitionPage(child: const History()),
              routes: [
                GoRoute(
                  parentNavigatorKey: _rootNavigatorKey,
                  path: '/quiz-review/:id/:mode',
                  builder: (context, state) {
                    final id = state.pathParameters['id'] ?? '';
                    final mode = state.pathParameters['mode'] ?? 'normal';
                    return QuizReview(quizId: id, mode: mode);
                  },
                ),
              ],
            ),
            GoRoute(
              path: '/guide',
              pageBuilder:
                  (context, state) => NoTransitionPage(child: const Guide()),
              routes: [
                GoRoute(
                  parentNavigatorKey: _rootNavigatorKey,
                  path: '/detail/:id',
                  pageBuilder: (context, state) {
                    final id = int.tryParse(state.pathParameters['id'] ?? '') ?? 0;

                    return MaterialPage(child: GuideDetail(id: id));
                  },
                ),
              ]
            ),
            GoRoute(
              path: '/user',
              pageBuilder:
                  (context, state) => NoTransitionPage(child: const User()),
              routes: [
                GoRoute(
                  parentNavigatorKey: _rootNavigatorKey,
                  path: '/update-profile',
                  pageBuilder:
                      (context, state) => MaterialPage(child: const UpdateProfile()),
                ),
              ]
            ),
          ],
        ),

        // Authentication Routes
        GoRoute(
          path: '/login',
          pageBuilder: (context, state) => MaterialPage(child: const Login()),
        ),
        GoRoute(
          path: '/register',
          pageBuilder:
              (context, state) => MaterialPage(child: const Register()),
        ),
        GoRoute(
          path: '/verification/:parentPage/:email',
          pageBuilder: (context, state) {
            final parentPage = state.pathParameters['parentPage'] ?? '';
            final email = state.pathParameters['email'] ?? '';
            return MaterialPage(
              child: Verification(email: email, parentPage: parentPage),
            );
          },
        ),
        // 修改密码 login和user都可以访问，独立路由
        GoRoute(
          path: '/change-password/:temporaryToken',
          pageBuilder: (context, state) {
            final temporaryToken = state.pathParameters['temporaryToken'] ?? '';
            return MaterialPage(
              child: ChangePassword(temporaryToken: temporaryToken),
            );
          },
        ),
      ],
    );
  }
}
