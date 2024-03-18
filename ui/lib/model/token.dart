import 'dart:convert';
import 'package:jwt_decoder/jwt_decoder.dart';

class Token {
  final String jwt;
  final String userID;
  final String username;

  Token._(this.jwt, this.userID, this.username);

  factory Token.fromJWT(String jwt) {
    final data = JwtDecoder.decode(jwt)['data'];
    final decoded = utf8.decode(base64.decode(data));
    final tokenPayload = jsonDecode(decoded);
    return Token._(jwt, tokenPayload['user_id'], tokenPayload['username']);
  }

  Token.fromJson(Map<String, dynamic> json)
      : jwt = json['jwt'] as String,
        userID = json['user_id'] as String,
        username = json['username'] as String;

  Map<String, dynamic> toJson() => {
        'jwt': jwt,
        'user_id': userID,
        'username': username,
      };
}
