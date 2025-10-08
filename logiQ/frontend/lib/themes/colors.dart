import 'package:flutter/material.dart';

class AppColors {
  static const lightPrimaryColor = Color(0xFF6949FF);
  static const lightSecondaryColor = Color(0xFFFFC107);

  static const darkPrimaryColor = Color(0xFF4730C5);
  static const darkSecondaryColor = Color(0xFFBF9101);

  // 不随主题变化的颜色，直接从AppColors中使用
  static const lightPurple = Color(0xFF8F7DFC);
  static const lightYellow = Color(0xFFFFE187);

  static const bluePurple = Color(0xFF95A4FC);

  static const grey1 = Color(0xFF9E9E9E);
  static const grey1WithOp = Color(0xBF9E9E9E);

  // 跟随主题变化的颜色，注册到主题中，提前判断
  static const red1 = Color(0xFFF85555);
  static const red2 = Color(0xFFB63E3E);

  static const green1 = Color(0xFF12D18E);
  static const green2 = Color(0xFF069F69);

  static const blue1 = Color(0xFF377AFF);
  static const blue2 = Color(0xFF1F5DB7);
}


