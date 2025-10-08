import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class BottomNavBar extends StatelessWidget {
  const BottomNavBar({super.key});

  static const _tabs = [
    '/home',
    '/history',
    '/guide',
    '/user'
  ];

  @override
  Widget build(BuildContext context) {

    final location = GoRouterState.of(context).uri.toString();
    int currentIndex = _tabs.indexWhere((t) => location.startsWith(t));
    if (currentIndex == -1) currentIndex = 0;

    return NavigationBar(
      onDestinationSelected: (int index) {
        context.go(_tabs[index]);
      },
      // indicatorColor: Theme.of(context).colorScheme,
      selectedIndex: currentIndex,
      destinations: <Widget>[
        NavigationDestination(
          selectedIcon: Icon(Icons.home),
          icon: Icon(Icons.home_outlined),
          label: 'Home',
        ),
        NavigationDestination(
          selectedIcon: Icon(Icons.list_alt_outlined),
          icon: Icon(Icons.list_alt),
          label: 'History',
        ),
        NavigationDestination(
          selectedIcon: Icon(Icons.tips_and_updates),
          icon: Icon(Icons.tips_and_updates_outlined),
          label: 'Guide',
        ),
        NavigationDestination(
          selectedIcon: Icon(Icons.person),
          icon: Icon(Icons.person_outline),
          label: 'User',
        ),
      ],
    );
  }
}
