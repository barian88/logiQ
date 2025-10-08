import 'package:flutter/material.dart';
import 'package:frontend/data/data.dart';
import 'package:gap/gap.dart';
import '../../../../themes/themes.dart';

class GuideDisplay extends StatelessWidget {
  const GuideDisplay({super.key, required this.guideModule});
  final GuideModule guideModule;
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    final rowHeight = 115.0 + 2 * 6; // Image height + vertical padding

    return LayoutBuilder(builder: (context,constraints){
      return SizedBox(
        height: rowHeight,
        width: constraints.maxWidth,
        child: Row(
          children: [
            DecoratedBox(
              decoration: BoxDecoration(
                color: theme.colorScheme.surfaceContainer,
                borderRadius: AppRadii.medium,
              ),
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 15, vertical: 6),
                child: guideModule.id == 5 ? Image.asset(guideModule.imagePath,height: 120, width: 120, ) :Image.asset(guideModule.imagePath, height: 120),
              ),
            ),
            Gap(20),
            Expanded(
              child: LayoutBuilder(
                builder: (context, constraint) {
                  // Calculate the width based on the available right space
                  final width = constraint.maxWidth;
                  return SizedBox(
                    height: rowHeight,
                    width: width,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Spacer(),
                        Row(
                          children: [
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    guideModule.title,
                                    style: theme.textTheme.bodyLarge?.copyWith(
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                  Gap(5),
                                  Text(
                                    guideModule.description,
                                    style: theme.textTheme.bodySmall?.copyWith(
                                      color: theme.colorScheme.onSurfaceVariant,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            Gap(10),
                            Icon(
                              Icons.arrow_forward_ios,
                              size: 16,
                              color: theme.colorScheme.primary,
                            ),
                          ],
                        ),
                        Spacer(),
                        Container(
                          alignment: Alignment.bottomCenter,
                          width: width,
                          height: 1,
                          color: theme.colorScheme.onSurfaceVariant.withValues(
                            alpha: 0.3,
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      );
    });


  }
}
