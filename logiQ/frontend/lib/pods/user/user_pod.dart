import 'dart:developer';

import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';
import '../../models/models.dart';
import 'package:frontend/repositories/repositories.dart';

part 'user_pod.g.dart';

@riverpod
class UserNotifier extends _$UserNotifier {
  @override
  AsyncValue<UserState> build() {
    // 保持provider活跃，防止状态丢失
    ref.keepAlive();
    // 异步初始化用户状态
    loadUserFromStorage();
    // 加载时返回loading状态
    return const AsyncValue.loading();
  }

  Future<void> login(String token, User user) async {
    // 保存到本地存储
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('auth_token', token);
    await prefs.setString('user_data', jsonEncode(user.toJson()));
    // 更新状态
    state = AsyncValue.data(
      UserState(isLoggedIn: true, token: token, user: user),
    );
  }

  Future<void> logout() async {
    // 清除本地存储
    final prefs = await SharedPreferences.getInstance();
    log('Logging out, clearing stored user data');
    await prefs.remove('auth_token');
    await prefs.remove('user_data');

    // 重置状态
    state = const AsyncValue.data(UserState());
  }

  Future<void> loadUserFromStorage() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('auth_token');
      final userData = prefs.getString('user_data');

      if (token != null && userData != null) {
        final userJson = jsonDecode(userData) as Map<String, dynamic>;
        final user = User.fromJson(userJson);

        state = AsyncValue.data(
          UserState(isLoggedIn: true, token: token, user: user),
        );
      } else {
        // 没有存储的用户数据，设置默认未登录状态
        state = const AsyncValue.data(UserState());
      }
    } catch (error, stackTrace) {
      state = const AsyncValue.data(UserState());
    }
  }

  // 更新用户信息
  Future<UpdateProfileResult> updateProfile(
    String username,
    String profilePictureUrl,
  ) async {
    final currentState = state.value;
    try {
      final repository = ref.read(userRepositoryProvider);
      // 调用更新接口,存入数据库
      final updatedUser = await repository.updateProfile(
        username,
        profilePictureUrl,
      );
      // 更新state
      state = AsyncValue.data(currentState!.copyWith(user: updatedUser));
      // 同步更新本地存储
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('user_data', jsonEncode(updatedUser.toJson()));

      return UpdateProfileResult.success(updatedUser);
    } catch (error, stackTrace) {
      if (currentState != null) {
        state = AsyncValue.data(currentState);
      } else {
        state = const AsyncValue.data(UserState());
      }
      log('Failed to update profile: $error', stackTrace: stackTrace);
      return UpdateProfileResult.error(
        'Failed to update profile: ${error.toString()}',
      );
    }
  }
}

class UserState {
  final bool isLoggedIn;
  final String? token;
  final User? user;

  const UserState({
    this.isLoggedIn = false,
    this.token,
    this.user,
  });

  UserState copyWith({
    bool? isLoggedIn,
    String? token,
    User? user,
  }) {
    return UserState(
      isLoggedIn: isLoggedIn ?? this.isLoggedIn,
      token: token ?? this.token,
      user: user ?? this.user,
    );
  }
}

// 更新用户信息结果
class UpdateProfileResult {
  final bool isSuccess;
  final String message;
  final User? updatedUser;

  const UpdateProfileResult({
    required this.isSuccess,
    this.updatedUser,
    required this.message,
  });

  factory UpdateProfileResult.success(User updatedUser) {
    return UpdateProfileResult(
      isSuccess: true,
      updatedUser: updatedUser,
      message: "Profile updated successfully",
    );
  }

  factory UpdateProfileResult.error(String message) {
    return UpdateProfileResult(isSuccess: false, message: message);
  }
}
