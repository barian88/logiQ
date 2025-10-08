import 'package:flutter/material.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import '../../widgets/widgets.dart';
import 'user/widgets/widgets.dart';

class TabScaffold extends StatelessWidget {
  final Widget child;

  const TabScaffold({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    final Widget title;
    final List<Widget> actions;

    final location = GoRouterState.of(context).uri.toString();
    if (location == '/user') {
      title = ThemeModeSwitch();
      actions = [LogoutButton(), Gap(16)];

    } else {
      title = const Logo();
      actions = [ThemeModeSwitch(), Gap(16)];
    }

    return Scaffold(
      appBar: AppBar(title: title, actions: actions),
      bottomNavigationBar: BottomNavBar(),
      body: child,
    );
  }
}
