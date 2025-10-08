import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../repositories/repositories.dart';
import '../../../apis/apis.dart';
import '../../../apis/app_exception.dart';

part 'verification_pod.g.dart';

@riverpod
class VerificationNotifier extends _$VerificationNotifier {
  final _verificationNotifierModel = VerificationNotifierModel();

  @override
  VerificationNotifierModel build() {
    return _verificationNotifierModel;
  }

  void update({
    String? verificationCode,
    bool? isLoading,
    String? errorMessage,
  }) {
    state = state.copyWith(
      verificationCode: verificationCode,
      isLoading: isLoading,
      errorMessage: errorMessage,
    );
  }

  // 发送验证码
  Future<VerificationResult> sendVerificationCode(String email) async {
    update(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      await authRepository.sendVerificationCode(email);
      
      update(isLoading: false);
      return VerificationResult.success();
      
    } catch (e) {
      // 如果是我们的自定义异常，直接使用其message
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }
      
      update(isLoading: false, errorMessage: errorMessage);
      return VerificationResult.error(errorMessage);
    }
  }

  // 注册流程（包含验证code和完成注册两部分）
  Future<VerificationResult> verifyRegistration(String email) async {
    if (state.verificationCode.isEmpty) {
      final error = "Please enter verification code";
      update(errorMessage: error);
      return VerificationResult.error(error);
    }

    update(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      
      // 第一步：验证验证码，获取临时token
      final verifyRequest = VerifyCodeRequest(
        email: email,
        verificationCode: state.verificationCode,
        purpose: 'registration',
      );
      
      final verifyResponse = await authRepository.verifyCode(verifyRequest);
      
      // 第二步：使用临时token完成注册
      final completeRequest = CompleteRegistrationRequest(
        temporaryToken: verifyResponse.temporaryToken,
      );
      
      await authRepository.completeRegistration(completeRequest);
      
      update(isLoading: false);
      return VerificationResult.success();
      
    } catch (e) {
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }
      
      update(isLoading: false, errorMessage: errorMessage);
      return VerificationResult.error(errorMessage);
    }
  }

  // 密码重置流程 （只是验证code，返回成功验证的token）
  Future<VerificationResult> verifyPasswordReset(String email) async {
    if (state.verificationCode.isEmpty) {
      final error = "Please enter verification code";
      update(errorMessage: error);
      return VerificationResult.error(error);
    }

    update(isLoading: true, errorMessage: null);

    try {
      final authRepository = ref.read(authRepositoryProvider);
      
      // 验证验证码，获取临时token（用于密码重置）
      final verifyRequest = VerifyCodeRequest(
        email: email,
        verificationCode: state.verificationCode,
        purpose: 'password_reset',
      );
      
      final verifyResponse = await authRepository.verifyCode(verifyRequest);
      
      update(isLoading: false);
      return VerificationResult.successWithToken(verifyResponse.temporaryToken);
      
    } catch (e) {
      String errorMessage;
      if (e is AppException) {
        errorMessage = e.message;
      } else {
        errorMessage = e.toString();
      }
      
      update(isLoading: false, errorMessage: errorMessage);
      return VerificationResult.error(errorMessage);
    }
  }


}

class VerificationNotifierModel {
  final String verificationCode;
  final bool isLoading;
  final String? errorMessage;

  VerificationNotifierModel({
    this.verificationCode = '',
    this.isLoading = false,
    this.errorMessage,
  });

  VerificationNotifierModel copyWith({
    String? verificationCode,
    bool? isLoading,
    String? errorMessage,
  }) {
    return VerificationNotifierModel(
      verificationCode: verificationCode ?? this.verificationCode,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}

class VerificationResult {
  final bool isSuccess;
  final String? errorMessage;
  final String? temporaryToken;

  const VerificationResult._({
    required this.isSuccess,
    this.errorMessage,
    this.temporaryToken,
  });

  factory VerificationResult.success() => const VerificationResult._(isSuccess: true);
  factory VerificationResult.successWithToken(String token) => VerificationResult._(
    isSuccess: true,
    temporaryToken: token,
  );
  factory VerificationResult.error(String message) => VerificationResult._(
    isSuccess: false,
    errorMessage: message,
  );
}