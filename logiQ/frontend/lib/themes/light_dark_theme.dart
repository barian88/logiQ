import 'package:flutter/material.dart';
import 'colors.dart';

final ThemeData lightTheme = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
      seedColor: AppColors.lightPrimaryColor,
      brightness: Brightness.light,
    ).copyWith(
      primary: AppColors.lightPrimaryColor,
      secondary: AppColors.lightSecondaryColor,
    ));

final ThemeData darkTheme = ThemeData(
  useMaterial3: true,
  colorScheme: ColorScheme.fromSeed(
    seedColor: AppColors.darkPrimaryColor,
    brightness: Brightness.dark,
  ).copyWith(
    primary: AppColors.darkPrimaryColor,
    secondary: AppColors.darkSecondaryColor,
  ),

);


extension CustomColorScheme on ColorScheme {
  Color get red =>
      brightness == Brightness.light
          ? AppColors.red1
          : AppColors.red2;
  Color get green =>
      brightness == Brightness.light
          ? AppColors.green1
          : AppColors.green2;
  Color get blue =>
      brightness == Brightness.light
          ? AppColors.blue1
          : AppColors.blue2;
}
