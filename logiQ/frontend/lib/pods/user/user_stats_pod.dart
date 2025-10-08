import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'package:frontend/models/models.dart';
import 'package:frontend/pods/user/user_pod.dart';
import 'package:frontend/repositories/repositories.dart';

part 'user_stats_pod.g.dart';

@riverpod
class UserStatsNotifier extends _$UserStatsNotifier {
  @override
  Future<UserStats?> build() async {
    // 缓存统计数据，避免频繁重建
    ref.keepAlive();

    final userState = ref.watch(userNotifierProvider).valueOrNull;
    // 如果用户未登录，返回null。
    // 否则token恢复还没完成，就获取用户统计数据，得到401错误，会被跳转到登录页
    if (userState == null || !userState.isLoggedIn) {
      return null;
    }

    return _fetchUserStats();
  }

  Future<UserStats> _fetchUserStats() {
    final repository = ref.read(userRepositoryProvider);
    return repository.getUserStats();
  }

  Future<void> refresh() async {
    final userState = ref.read(userNotifierProvider).valueOrNull;
    if (userState == null || !userState.isLoggedIn) {
      state = const AsyncValue.data(null);
      return;
    }

    state = const AsyncValue.loading();
    state = await AsyncValue.guard(_fetchUserStats);
  }
}
