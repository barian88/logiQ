class TimeFormatterUtil{
  static String getFormattedTime(int time) {
    if (time == 0) {
      return '00:00:00';
    }

    final hours = time ~/ 3600;
    final minutes = (time % 3600) ~/ 60;
    final seconds = time % 60;

    final hoursStr = hours.toString().padLeft(2, '0');
    final minutesStr = minutes.toString().padLeft(2, '0');
    final secondsStr = seconds.toString().padLeft(2, '0');

    return '$hoursStr:$minutesStr:$secondsStr';
  }

  // 格式化完成时间用于前端显示 (相对时间)
  static String formatRelativeTime(DateTime? dateTime) {
    if (dateTime == null) return 'unknown';
    
    final now = DateTime.now();
    final difference = now.difference(dateTime);
    
    if (difference.inDays > 0) {
      return '${difference.inDays} Days Ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours} Hours Ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes} Minutes Ago';
    } else {
      return 'Just Now';
    }
  }
  
  // 详细时间格式 (YYYY-MM-DD HH:mm)
  static String formatDetailedTime(DateTime? dateTime) {
    if (dateTime == null) return 'unknown';
    return '${dateTime.day.toString().padLeft(2, '0')}-${dateTime.month.toString().padLeft(2, '0')}-${dateTime.year} · ${formatRelativeTime(dateTime)}';
  }
  
  // 格式化答题用时 (分钟+秒)
  static String formatCompletionTime(int timeInSeconds) {
    if (timeInSeconds == 0) return '0 s';

    final minutes = timeInSeconds ~/ 60;
    final seconds = timeInSeconds % 60;

    if (minutes > 0) {
      return '${minutes} min ${seconds} sec';
    } else {
      return '${seconds} sec';
    }
  }

}