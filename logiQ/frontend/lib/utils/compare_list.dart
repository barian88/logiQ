bool areListsEqual(List<int> a, List<int> b) {
  if (a.length != b.length) {
    return false;
  }

  // 用 Set 构建集合
  final setA = a.toSet();

  // 检查 b 的每个元素是否在 setA 中
  for (var v in b) {
    if (!setA.contains(v)) {
      return false;
    }
  }

  return true;
}