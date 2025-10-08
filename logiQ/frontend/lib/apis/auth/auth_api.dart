import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/models.dart';
import '../api_service.dart';

part 'auth_api.g.dart';

@riverpod
AuthApi authApi(Ref ref) {
  return AuthApi(ref.read(apiServiceProvider));
}

class AuthApi {
  final ApiService _apiService;

  AuthApi(this._apiService);

  Future<AuthResponse> login(LoginRequest request) async {
    final response = await _apiService.post(
      '/auth/login',
      body: request.toJson(),
    );
    return AuthResponse.fromJson(response);
  }

  Future<void> registerRequest(RegisterRequest request) async {
    await _apiService.post('/auth/register-request', body: request.toJson());
  }

  Future<VerifyCodeResponse> verifyCode(VerifyCodeRequest request) async {
    final response = await _apiService.post(
      '/auth/verify-code',
      body: request.toJson(),
    );
    return VerifyCodeResponse.fromJson(response);
  }

  Future<User> completeRegistration(CompleteRegistrationRequest request) async {
    final response = await _apiService.post(
      '/auth/complete-registration',
      body: request.toJson(),
    );
    return User.fromJson(response);
  }

  Future<void> logout() async {
    await _apiService.post('/auth/logout');
  }

  Future<void> sendVerificationCode(String email) async {
    await _apiService.post('/auth/send-verification', body: {'email': email});
  }

  Future<void> resetPassword(ResetPasswordRequest request) async {
    await _apiService.post('/auth/update-password', body: request.toJson());
  }
}

// 请求模型
class LoginRequest {
  final String email;
  final String password;

  LoginRequest({required this.email, required this.password});

  Map<String, dynamic> toJson() => {'email': email, 'password': password};
}

class RegisterRequest {
  final String email;
  final String password;
  final String username;

  RegisterRequest({
    required this.email,
    required this.password,
    required this.username,
  });

  Map<String, dynamic> toJson() => {
    'email': email,
    'password': password,
    'username': username,
  };
}

class VerifyCodeRequest {
  final String email;
  final String verificationCode;
  final String purpose;

  VerifyCodeRequest({
    required this.email,
    required this.verificationCode,
    required this.purpose,
  });

  Map<String, dynamic> toJson() => {
    'email': email,
    'verification_code': verificationCode,
    'purpose': purpose,
  };
}

class CompleteRegistrationRequest {
  final String temporaryToken;

  CompleteRegistrationRequest({required this.temporaryToken});

  Map<String, dynamic> toJson() => {'temporary_token': temporaryToken};
}

class ResetPasswordRequest {
  final String temporaryToken;
  final String newPassword;

  ResetPasswordRequest({
    required this.temporaryToken,
    required this.newPassword,
  });

  Map<String, dynamic> toJson() => {
    'temporary_token': temporaryToken,
    'new_password': newPassword,
  };
}

// 响应模型
class AuthResponse {
  final String token;
  final User user;

  AuthResponse({required this.token, required this.user});

  factory AuthResponse.fromJson(Map<String, dynamic> json) =>
      AuthResponse(token: json['token'], user: User.fromJson(json['user']));
}

class VerifyCodeResponse {
  final bool success;
  final String message;
  final String temporaryToken;
  final int expiresIn;

  VerifyCodeResponse({
    required this.success,
    required this.message,
    required this.temporaryToken,
    required this.expiresIn,
  });

  factory VerifyCodeResponse.fromJson(Map<String, dynamic> json) =>
      VerifyCodeResponse(
        success: json['success'],
        message: json['message'],
        temporaryToken: json['temporary_token'],
        expiresIn: json['expires_in'],
      );
}
