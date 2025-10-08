import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/apis/apis.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/pods/pods.dart';
import 'package:gap/gap.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import '../../../../utils/utils.dart';

class UserInfoCard extends ConsumerWidget {
  const UserInfoCard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userState = ref.watch(userNotifierProvider);
    final userNotifier = ref.read(userNotifierProvider.notifier);
    final user = userState.value?.user;

    final uploadApi = ref.read(uploadApiProvider);

    if (user == null) {
      return const SizedBox.shrink();
    }

    final theme = Theme.of(context);
    final textColor = Colors.white;

    handleUpdateProfilePicture() async {
      // 使用 image_picker 选择图片
      final picker = ImagePicker();
      final XFile? file = await picker.pickImage(source: ImageSource.gallery);

      // 用户选择了图片
      if (file != null) {
        // 显示加载指示器
        showDialog(
          context: context,
          barrierDismissible: false,
          builder: (context) => const Center(child: CircularProgressIndicator()),
        );
        try {
          // 读取图片文件为字节
          final bytes = await file.readAsBytes();
          final blobName = 'avatars/${user.id}_${DateTime.now().millisecondsSinceEpoch}';
          final suffix = file.name.split('.').last;
          // 上传图片并获取公开 URL
          final publicUrl = await uploadApi.uploadAvatar(
            bytes,
            blobName: blobName,
            contentType: 'image/$suffix',
          );

          // 更新用户资料
          final result = await userNotifier.updateProfile(
            user.username,
            publicUrl,
          );

          if (!context.mounted) return;

          if (result.isSuccess) {
            ToastHelper.showSuccess(theme, 'Profile updated successfully!');
          } else {
            ToastHelper.showError(theme, result.message);
          }
        } catch (e) {
          if (context.mounted) {
            ToastHelper.showError(
              theme,
              ErrorHandler.getErrorMessage(e),
            );
          }
        } finally {
          if (context.mounted) {
            Navigator.of(context, rootNavigator: true).pop(); // 关掉 loading
          }
        }

      }

    }

    return Stack(
      clipBehavior: Clip.none,
      children: [
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
            gradient: AppGradients.cardPurpleGradient(theme),
            borderRadius: AppRadii.medium,
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(16, 35, 16, 30),
            child: Column(
              children: [
                InkWell(
                  onTap: () {
                    context.push('/user/update-profile');
                  },
                  child: Row(
                    children: [
                      Text('Username', style: TextStyle(color: textColor)),
                      const Spacer(),
                      Text(
                        user.username,
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: textColor,
                        ),
                      ),
                      const Gap(5),
                      Icon(Icons.arrow_forward_ios, color: textColor, size: 16),
                    ],
                  ),
                ),
                const Gap(16),
                Row(
                  children: [
                    Text('Email', style: TextStyle(color: textColor)),
                    const Spacer(),
                    Text(user.email, style: TextStyle(color: textColor)),
                  ],
                ),
                const Gap(16),
                InkWell(
                  onTap: () {
                    final email = user.email;

                    context.push('/verification/user/$email');
                  },
                  child: Row(
                    children: [
                      Text('Change Password', style: TextStyle(color: textColor)),
                      const Spacer(),
                      Icon(Icons.arrow_forward_ios, color: textColor, size: 16),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        Align(
          alignment: Alignment.topCenter,
          child: Transform.translate(
            offset: const Offset(0, -35),
            child: CircleAvatar(
              radius: 37,
              backgroundColor: AppColors.grey1,
              child: InkWell(
                onTap: handleUpdateProfilePicture,
                child: CircleAvatar(
                  radius: 35,
                  backgroundImage: NetworkImage(user.profilePictureUrl),
                  onBackgroundImageError: (error, stackTrace) {
                    // 处理图片加载失败
                  },
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
