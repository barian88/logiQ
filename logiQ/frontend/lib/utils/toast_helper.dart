import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:frontend/themes/themes.dart';
import '../themes/light_dark_theme.dart';

class ToastHelper {
  static Future<void> showToast({
    required String message,
    required Color backgroundColor,
  }) async {
    await Fluttertoast.showToast(
      timeInSecForIosWeb: 2,
      msg: message,
      toastLength: Toast.LENGTH_LONG,
      gravity: ToastGravity.BOTTOM,
      backgroundColor: backgroundColor,
      textColor: Colors.white,
      fontSize: 14.0,
    );
  }

  static Future<void> showError(ThemeData theme, String message) async {
    await showToast(
      message: message,
      backgroundColor: theme.colorScheme.red,
    );
  }

  static Future<void> showSuccess(ThemeData theme, String message) async {
    await showToast(
      message: message,
      backgroundColor: AppColors.grey1WithOp
    );
  }

  static Future<void> showWarning(ThemeData theme, String message) async {
    await showToast(
      message: message,
      backgroundColor: theme.colorScheme.secondary,
    );
  }
}