import 'package:flutter/material.dart';

class BaseContainer extends StatelessWidget {
  final Widget child;
  final bool isScrollable;

  const BaseContainer({
    super.key,
    required this.child,
    this.isScrollable = true,
  });

  @override
  Widget build(BuildContext context) {
    if (isScrollable) {
      return SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 12.0),
          child: child,
        ),
      );
    } else {
      return Padding(
        padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 12.0),
        child: child,
      );
    }
  }
}
