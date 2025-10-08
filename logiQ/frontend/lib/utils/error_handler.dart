import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../apis/apis.dart';
import 'toast_helper.dart';

class ErrorHandler {
  static String getErrorMessage(Object error) {
    if (error is AppException) {
      return error.message;
    }
    
    if (error is NetworkException) {
      return 'Network connection failed. Please check your network.';
    }
    
    if (error is ServerException) {
      switch (error.statusCode) {
        case 404:
          return 'Requested resource not found';
        case 500:
          return 'Server error. Please try again later.';
        case 503:
          return 'Service unavailable. Please try again later.';
        default:
          return error.message;
      }
    }
    
    if (error is ValidationException) {
      return error.message;
    }
    
    if (error is AuthException) {
      return error.message;
    }
    
    return 'Unknown error. Please try again later.';
  }

  static void showErrorToast(BuildContext context, Object error) {
    ToastHelper.showError(Theme.of(context), getErrorMessage(error));
  }

  static Widget buildErrorWidget(Object error, {VoidCallback? onRetry, BuildContext? context}) {
    final message = getErrorMessage(error);
    
    return Builder(
      builder: (context) {
        final colorScheme = Theme.of(context).colorScheme;
        
        return Center(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  Icons.error_outline,
                  size: 64,
                  color: colorScheme.onSurface.withValues(alpha: 0.5),
                ),
                const SizedBox(height: 16),
                Text(
                  message,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 16,
                    color: colorScheme.onSurface.withValues(alpha: 0.7),
                  ),
                ),
                if (onRetry != null) ...[
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: onRetry,
                    child: const Text('Retry'),
                  ),
                ],
              ],
            ),
          ),
        );
      },
    );
  }
}

// 扩展方法，方便在Widget中使用
extension AsyncValueErrorHandling<T> on AsyncValue<T> {
  Widget buildErrorWidget({VoidCallback? onRetry}) {
    return when(
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stackTrace) => ErrorHandler.buildErrorWidget(error, onRetry: onRetry),
      data: (data) => const SizedBox.shrink(),
    );
  }
}
