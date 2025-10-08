import 'dart:convert';
import 'dart:typed_data';

import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart' as http;
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../api_service.dart';
import '../app_exception.dart';

part 'upload_api.g.dart';

@riverpod
UploadApi uploadApi(Ref ref) {
  return UploadApi(ref.read(apiServiceProvider));
}

class UploadApi {
  UploadApi(this._apiService);

  final ApiService _apiService;

  Future<String> uploadAvatar(
    Uint8List bytes, {
    required String blobName,
    String contentType = 'image/jpeg',
  }) async {
    final encodedBlobName = Uri.encodeComponent(blobName);
    // 获取预签名URL
    final rawBody = await _apiService.getString(
      '/user/presign-upload?blob_name=$encodedBlobName',
    );
    // 从响应中提取SAS URL
    final sasUrl = _extractSasUrl(rawBody);
    // 使用预签名URL上传文件
    final uploadResponse = await http.put(
      Uri.parse(sasUrl),
      headers: {
        'x-ms-blob-type': 'BlockBlob',
        'Content-Type': contentType,
        'Content-Length': bytes.length.toString(),
      },
      body: bytes,
    );

    if (uploadResponse.statusCode != 201 &&
        uploadResponse.statusCode != 200) {
      throw ServerException(
        'Failed to upload avatar',
        uploadResponse.statusCode,
        originalError: uploadResponse.body,
      );
    }

    return sasUrl.split('?').first;
  }

  String _extractSasUrl(String rawBody) {
    final body = rawBody.trim();
    if (body.isEmpty) {
      throw AppException('empty response body');
    }

    try {
      final decoded = jsonDecode(body);
      if (decoded is String) {
        return decoded;
      }
      if (decoded is Map<String, dynamic>) {
        for (final key in const ['url', 'sasUrl', 'sas', 'data']) {
          final value = decoded[key];
          if (value is String && value.isNotEmpty) {
            return value;
          }
        }
      }
    } catch (e, stack) {
      if (body.startsWith('http')) {
        return body;
      }
    }

    throw AppException('fail to extract sasUrl from response: $body');
  }
}
