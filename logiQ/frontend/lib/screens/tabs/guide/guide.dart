import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:frontend/themes/themes.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import 'widgets/widgets.dart';
import '../../../widgets/widgets.dart';
import '../../../data/data.dart';


class Guide extends StatelessWidget {
  const Guide({super.key});

  @override
  Widget build(BuildContext context) {

    final guideModules = allGuideModules;

    final displayList = List.generate(guideModules.length, (index){
      final item = guideModules[index];
      return Column(
        children: [
          InkWell(
              onTap: (){
                context.push('/guide/detail/${item.id}');
              },
              child: GuideDisplay(guideModule: item)),
          Gap(30),
        ],
      );
    });
    final top = Container(
      decoration:BoxDecoration(
        gradient: AppGradients.guideCardGradient(Theme.of(context)),
        borderRadius: AppRadii.medium,
      ),
      width: double.infinity,
      height: 180,
      child: Center(
        child: Column(
          mainAxisSize:  MainAxisSize.min,
          children: [
            Icon(FontAwesomeIcons.solidStar , size: 40, color: Colors.grey[200],),
            Gap(20),
            Text('LogiQ Features', style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: Colors.grey[200],
              fontWeight: FontWeight.bold,
            )),
            Gap(30),
          ],
        ),
      ),
    );
    return BaseContainer(child: Column(children: [top, ...displayList],));
  }
}
