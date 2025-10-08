class User {
  final String id;
  final String username;
  final String email;
  final String profilePictureUrl;

  User({
    required this.id,
    required this.username,
    required this.email,
    required this.profilePictureUrl,
  });

  User copyWith({
    String? id,
    String? username,
    String? email,
    String? profilePictureUrl,
  }) {
    return User(
      id: id ?? this.id,
      username: username ?? this.username,
      email: email ?? this.email,
      profilePictureUrl: profilePictureUrl ?? this.profilePictureUrl,
    );
  }

  // JSON 序列化支持
  factory User.fromJson(Map<String, dynamic> json) => User(
    id: json['_id'] ?? json['id'] ?? '', // 兼容MongoDB的_id字段
    username: json['username'] ?? '',
    email: json['email'] ?? '',
    profilePictureUrl: json['profile_picture_url'] ?? '',
  );

  Map<String, dynamic> toJson() => {
    '_id': id,
    'username': username,
    'email': email,
    'profile_picture_url': profilePictureUrl,
  };
}