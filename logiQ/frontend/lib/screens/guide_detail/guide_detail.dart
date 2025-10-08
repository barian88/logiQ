import 'package:flutter/material.dart';
import 'package:frontend/widgets/widgets.dart';
import 'package:frontend/themes/themes.dart';
import 'package:frontend/data/data.dart';
import 'package:gap/gap.dart';
import 'package:video_player/video_player.dart';

const mediaHeight = 500.0;
const mediaAspectRatio = 1179 / 2556;
const mediaWidth = mediaHeight * mediaAspectRatio;
const adminMediaWidth = 500.0;

class GuideDetail extends StatefulWidget {
  final int id;
  const GuideDetail({super.key, required this.id});

  @override
  State<GuideDetail> createState() => _GuideDetailState();
}

class _GuideDetailState extends State<GuideDetail> {
  late PageController _pageController;
  int _currentPage = 0;
  late GuideModule _currentModule;

  @override
  void initState() {
    super.initState();
    // 根据传入的 id 查找对应的模块
    _currentModule = allGuideModules.firstWhere(
      (module) => module.id == widget.id,
      orElse: () => allGuideModules.first, // 如果未找到，则默认显示第一个模块
    );
    _pageController = PageController(initialPage: _currentPage);
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        title: Text(_currentModule.title), // 显示模块名称
        centerTitle: true,
      ),
      body: BaseContainer(
        isScrollable: false,
        child: Column(
          children: [
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: _currentModule.items.length,
                onPageChanged: (index) {
                  setState(() {
                    _currentPage = index;
                  });
                },
                itemBuilder: (context, index) {
                  final item = _currentModule.items[index];
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _GuideMedia(item: item, id: widget.id),
                      const Gap(24),
                      // 标题
                      Text(
                        item.title,
                        style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const Gap(16),
                      // 描述
                      Text(
                        item.description,
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                    ],
                  );
                },
              ),
            ),
            // 页面指示器
            Padding(
              padding: const EdgeInsets.only(bottom: 16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(_currentModule.items.length, (index) {
                  return Container(
                    margin: const EdgeInsets.symmetric(horizontal: 4.0),
                    width: _currentPage == index ? 12.0 : 8.0,
                    height: 8.0,
                    decoration: BoxDecoration(
                      color: _currentPage == index ? Theme.of(context).colorScheme.primary : Colors.grey,
                      borderRadius: BorderRadius.circular(4.0),
                    ),
                  );
                }),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _GuideMedia extends StatelessWidget {
  const _GuideMedia({required this.item, required this.id});

  final GuideItem item;
  final int id;

  @override
  Widget build(BuildContext context) {

    final theme = Theme.of(context);

    if (item.hasVideo) {
      final videoPath = item.videoAsset;
      if (videoPath != null && videoPath.isNotEmpty) {
        return _GuideVideoPlayer(assetPath: videoPath, id: id,);
      }
      return _GuidePlaceholder(message: 'Video asset not provided');
    }

    // 没视频返回图片组件
    final assetPath = item.imagePath;
    if (assetPath == null || assetPath.isEmpty) {
      return _GuidePlaceholder(message: 'Image asset not provided');
    }

    return Center(
      child: Container(
        width: mediaWidth,
        height: mediaHeight,
        decoration: BoxDecoration(
          borderRadius: AppRadii.medium,
          border: Border.all(color: theme.colorScheme.outline, width: 1.0),
        ),
        child: ClipRRect(
          borderRadius: AppRadii.medium,
          child: Image.asset(
            assetPath,
            fit: BoxFit.cover,
            width: mediaWidth,
            height: mediaHeight,
            errorBuilder: (context, error, stackTrace) => Container(
              width: mediaWidth,
              height: mediaHeight,
              color: Colors.grey[300],
              alignment: Alignment.center,
              child: const Text('Image failed to load'),
            ),
          ),
        ),
      )
    );
  }
}

class _GuidePlaceholder extends StatelessWidget {
  const _GuidePlaceholder({required this.message});

  final String message;

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: AppRadii.medium,
      child: Container(
        width: mediaWidth,
        height: mediaHeight,
        color: Colors.grey[300],
        alignment: Alignment.center,
        child: Text(message),
      ),
    );
  }
}

// 视频播放
class _GuideVideoPlayer extends StatefulWidget {
  const _GuideVideoPlayer({required this.assetPath, required this.id});

  final String assetPath;
  final int id; //用来判读是不是admin Id == 5
  @override
  State<_GuideVideoPlayer> createState() => _GuideVideoPlayerState();
}

class _GuideVideoPlayerState extends State<_GuideVideoPlayer> {
  late VideoPlayerController _controller;
  late Future<void> _initialization;

  @override
  void initState() {
    super.initState();
    _controller = VideoPlayerController.asset(widget.assetPath);
    _initialization = _initialiseController();
  }

  Future<void> _initialiseController() async {
    await _controller.initialize();
    await _controller.setLooping(true);
    await _controller.setVolume(0);
    if (mounted) {
      await _controller.play();
    }
  }

  @override
  void didUpdateWidget(covariant _GuideVideoPlayer oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.assetPath != widget.assetPath) {
      final oldController = _controller;
      _controller = VideoPlayerController.asset(widget.assetPath);
      _initialization = _initialiseController();
      oldController.dispose();
      setState(() {});
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Container(
        width: widget.id == 5 ? adminMediaWidth : mediaWidth,
        height: mediaHeight,
        decoration: BoxDecoration(
          borderRadius: AppRadii.medium,
          border: Border.all(color: theme.colorScheme.outline),
        ),
        child: ClipRRect(
          borderRadius: AppRadii.medium,
          child: FutureBuilder<void>(
            future: _initialization,
            builder: (context, snapshot) {
              if (snapshot.connectionState != ConnectionState.done) {
                return const Center(child: CircularProgressIndicator());
              }

              if (snapshot.hasError) {
                return Container(
                  color: Colors.grey[300],
                  alignment: Alignment.center,
                  child: const Text('Video failed to load'),
                );
              }

              final videoValue = _controller.value;
              if (videoValue.size.isEmpty) {
                return Container(
                  color: Colors.grey[300],
                  alignment: Alignment.center,
                  child: const Text('Unable to play video'),
                );
              }
              return FittedBox(
                fit: BoxFit.contain,
                child: SizedBox(
                  width: videoValue.size.width,
                  height: videoValue.size.height,
                  child: VideoPlayer(_controller),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}
