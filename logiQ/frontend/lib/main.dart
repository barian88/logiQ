import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'pods/pods.dart';
import 'themes/themes.dart';

void main() {
  runApp(ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerStatefulWidget {
  const MyApp({super.key});

  @override
  ConsumerState<MyApp> createState() => _MyAppState();
}

class _MyAppState extends ConsumerState<MyApp> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) async {
    if (state == AppLifecycleState.paused || state == AppLifecycleState.inactive) {
      // Save the current state of the app when it goes into the background
    }
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routeNotifierProvider);
    final themeMode = ref.watch(themeModeNotifierProvider);

    return MaterialApp.router(
      debugShowCheckedModeBanner: false, // 关闭右上角 Debug 横幅
      routerConfig: router,
      title: 'flutter_application',
      themeMode: themeMode,
      theme: lightTheme,
      darkTheme: darkTheme,
    );
  }
}
