import 'package:flutter/material.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
part 'theme_mode_pod.g.dart';


@riverpod
class ThemeModeNotifier extends _$ThemeModeNotifier {
  @override
  ThemeMode build() {
    return ThemeMode.light;
  }

  void toggleThemeMode(bool value) {
    state = value ? ThemeMode.dark : ThemeMode.light;
  }
}
