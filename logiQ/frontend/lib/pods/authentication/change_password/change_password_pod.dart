import 'package:frontend/apis/apis.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:frontend/repositories/repositories.dart';
part 'change_password_pod.g.dart';

@riverpod
class ChangePasswordNotifier extends _$ChangePasswordNotifier {
  @override
  ChangePasswordNotifierModel build() {
    return const ChangePasswordNotifierModel();
  }

  Future<ResetPasswordResult> resetPassword(String temporaryToken, String newPassword) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      final resetPasswordRequest = ResetPasswordRequest(
        temporaryToken: temporaryToken,
        newPassword: newPassword,
      );

      await authRepository.resetPassword(resetPasswordRequest);
      state = state.copyWith(isLoading: false);
      return ResetPasswordResult.success();
    } catch (e) {
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }

      state = state.copyWith(isLoading: false, errorMessage: errorMessage);
      return ResetPasswordResult.error(errorMessage);
    }
  }
}

class ChangePasswordNotifierModel {
  final bool isLoading;
  final String? errorMessage;

  const ChangePasswordNotifierModel({
    this.isLoading = false,
    this.errorMessage,
  });

  ChangePasswordNotifierModel copyWith({
    bool? isLoading,
    String? errorMessage,
  }) {
    return ChangePasswordNotifierModel(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

class ResetPasswordResult {
  final bool isSuccess;
  final String? errorMessage;

  const ResetPasswordResult._({required this.isSuccess, this.errorMessage});

  factory ResetPasswordResult.success() =>
      const ResetPasswordResult._(isSuccess: true);
  factory ResetPasswordResult.error(String message) =>
      ResetPasswordResult._(isSuccess: false, errorMessage: message);
}