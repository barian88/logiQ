import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../repositories/repositories.dart';
import '../../../apis/apis.dart';
import '../../../apis/app_exception.dart';

part 'register_pod.g.dart';

@riverpod
class RegisterNotifier extends _$RegisterNotifier {
  @override
  RegisterNotifierModel build() {
    return const RegisterNotifierModel();
  }

  Future<RegisterResult> sendRegisterRequest(String username, String email, String password) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      final registerRequest = RegisterRequest(
        email: email,
        password: password,
        username: username,
      );

      await authRepository.registerRequest(registerRequest);

      state = state.copyWith(isLoading: false);
      return RegisterResult.success();

    } catch (e) {
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }

      state = state.copyWith(isLoading: false, errorMessage: errorMessage);
      return RegisterResult.error(errorMessage);
    }
  }
}

class RegisterNotifierModel {
  final bool isLoading;
  final String? errorMessage;

  const RegisterNotifierModel({
    this.isLoading = false,
    this.errorMessage,
  });

  RegisterNotifierModel copyWith({
    bool? isLoading,
    String? errorMessage,
  }) {
    return RegisterNotifierModel(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}

class RegisterResult {
  final bool isSuccess;
  final String? errorMessage;

  const RegisterResult._({
    required this.isSuccess,
    this.errorMessage,
  });

  factory RegisterResult.success() => const RegisterResult._(isSuccess: true);
  factory RegisterResult.error(String message) => RegisterResult._(
        isSuccess: false,
        errorMessage: message,
      );
}
