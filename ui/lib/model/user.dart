import 'dart:convert';
import 'package:jwt_decoder/jwt_decoder.dart';

class User {
  final String id;
  final String username;
  final String jwt;

  User._(this.id, this.username, this.jwt);

  factory User.fromJWT(String jwt) {
    final data = JwtDecoder.decode(jwt)['data'];
    final decoded = utf8.decode(base64.decode(data));
    final tokenPayload = jsonDecode(decoded);
    return User._(tokenPayload['user_id'], tokenPayload['username'], jwt);
  }

  User.fromJson(Map<String, dynamic> json)
      : id = json['id'] as String,
        username = json['username'] as String,
        jwt = json['jwt'] as String;

  Map<String, dynamic> toJson() => {
        'id': id,
        'username': username,
        'jwt': jwt,
      };
}
