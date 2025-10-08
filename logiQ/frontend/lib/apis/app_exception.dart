class AppException implements Exception {
  final String message;
  final String? code;
  final dynamic originalError;

  AppException(this.message, {this.code, this.originalError});

  @override
  String toString() => 'AppException: $message';
}

class NetworkException extends AppException {
  NetworkException(String message, {String? code, dynamic originalError})
      : super(message, code: code, originalError: originalError);
}

class ServerException extends AppException {
  final int statusCode;
  
  ServerException(String message, this.statusCode, {String? code, dynamic originalError})
      : super(message, code: code, originalError: originalError);
}

class ValidationException extends AppException {
  final Map<String, String> errors;
  
  ValidationException(String message, this.errors, {String? code, dynamic originalError})
      : super(message, code: code, originalError: originalError);
}

class AuthException extends AppException {
  AuthException(String message, {String? code, dynamic originalError})
      : super(message, code: code, originalError: originalError);
}

class CacheException extends AppException {
  CacheException(String message, {String? code, dynamic originalError})
      : super(message, code: code, originalError: originalError);
}