import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../repositories/repositories.dart';
import '../../../apis/apis.dart';
import '../../../apis/app_exception.dart';
import '../../user/user_pod.dart';

part 'login_pod.g.dart';

@riverpod
class LoginNotifier extends _$LoginNotifier {
  @override
  LoginNotifierModel build() {
    return const LoginNotifierModel();
  }

  Future<LoginResult> login(String email, String password) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      final userNotifier = ref.read(userNotifierProvider.notifier);

      final loginRequest = LoginRequest(email: email, password: password);
      final authResponse = await authRepository.login(loginRequest);

      // Save user info and token to UserPod
      await userNotifier.login(authResponse.token, authResponse.user);

      state = state.copyWith(isLoading: false);
      return LoginResult.success();

    } catch (e) {
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }

      state = state.copyWith(isLoading: false, errorMessage: errorMessage);
      return LoginResult.error(errorMessage);
    }
  }
}

class LoginNotifierModel {
  final bool isLoading;
  final String? errorMessage;

  const LoginNotifierModel({
    this.isLoading = false,
    this.errorMessage,
  });

  LoginNotifierModel copyWith({
    bool? isLoading,
    String? errorMessage,
  }) {
    return LoginNotifierModel(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}

class LoginResult {
  final bool isSuccess;
  final String? errorMessage;

  const LoginResult._({
    required this.isSuccess,
    this.errorMessage,
  });

  factory LoginResult.success() => const LoginResult._(isSuccess: true);
  factory LoginResult.error(String message) => LoginResult._(
        isSuccess: false,
        errorMessage: message,
      );
}
