import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/models.dart';
import '../api_service.dart';

part 'user_api.g.dart';

@riverpod
UserApi userApi(Ref ref) {
  return UserApi(ref.read(apiServiceProvider));
}

class UserApi {
  final ApiService _apiService;

  UserApi(this._apiService);

  Future<User> getCurrentUser() async {
    final response = await _apiService.get('/users/profile');
    return User.fromJson(response);
  }

  Future<UserStats> getUserStats() async {
    final response = await _apiService.get('/user-stats');
    return UserStats.fromJson(response);
  }

  Future<User> updateProfile(String username, String profilePictureUrl) async {
    // 后端做了空字段忽略处理
    final response = await _apiService.post('/user/update', body: {
      'username': username,
      'profile_picture_url': profilePictureUrl,
    });
    return User.fromJson(response);
  }

}


