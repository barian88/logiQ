import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../pods/pods.dart';

class ThemeModeSwitch extends ConsumerWidget {
  const ThemeModeSwitch({
    super.key,
  });


  @override
  Widget build(BuildContext context, WidgetRef ref) {

    final themeMode = ref.watch(themeModeNotifierProvider);
    final themeNotifier = ref.read(themeModeNotifierProvider.notifier);

    final WidgetStateProperty<Icon> thumbIcon = WidgetStateProperty<Icon>.fromMap(
      <WidgetStatesConstraint, Icon>{
        WidgetState.selected: Icon(Icons.dark_mode, color: Theme.of(context).colorScheme.primary,),
        WidgetState.any: Icon(Icons.light_mode),
      },
    );

    return Switch(
      thumbIcon: thumbIcon,
      value: themeMode == ThemeMode.dark,
      onChanged: (value) {
        themeNotifier.toggleThemeMode(value);
      },
    );
  }
}
