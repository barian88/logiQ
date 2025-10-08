import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/models.dart';
import '../../apis/apis.dart';

part 'auth_repository.g.dart';

@riverpod
AuthRepository authRepository(Ref ref) {
  return AuthRepositoryImpl(ref.read(authApiProvider));
}

abstract class AuthRepository {
  Future<AuthResponse> login(LoginRequest request);
  Future<void> registerRequest(RegisterRequest request);
  Future<VerifyCodeResponse> verifyCode(VerifyCodeRequest request);
  Future<User> completeRegistration(CompleteRegistrationRequest request);
  Future<void> logout();
  Future<void> sendVerificationCode(String email);
  Future<void> resetPassword(ResetPasswordRequest request);
}

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi _authApi;

  AuthRepositoryImpl(this._authApi);

  @override
  Future<AuthResponse> login(LoginRequest request) async {
    return _authApi.login(request);
  }

  @override
  Future<void> registerRequest(RegisterRequest request) async {
    return _authApi.registerRequest(request);
  }

  @override
  Future<VerifyCodeResponse> verifyCode(VerifyCodeRequest request) async {
    return _authApi.verifyCode(request);
  }

  @override
  Future<User> completeRegistration(CompleteRegistrationRequest request) async {
    return _authApi.completeRegistration(request);
  }

  @override
  Future<void> logout() async {
    return _authApi.logout();
  }


  @override
  Future<void> sendVerificationCode(String email) async {
    return _authApi.sendVerificationCode(email);
  }

  @override
  Future<void> resetPassword(ResetPasswordRequest request){
    return _authApi.resetPassword(request);
  }
}