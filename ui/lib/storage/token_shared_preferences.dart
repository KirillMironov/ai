import 'dart:async';
import 'package:ai/storage/token.dart';
import 'package:shared_preferences/shared_preferences.dart';

class TokenStorageSharedPreferences implements TokenStorage {
  late SharedPreferences _prefs;
  final _jwtKey = 'jwt';

  TokenStorageSharedPreferences() {
    _initPrefs();
  }

  Future<void> _initPrefs() async {
    _prefs = await SharedPreferences.getInstance();
  }

  @override
  void saveToken(String token) async {
    await _prefs.setString(_jwtKey, token);
  }

  @override
  String? getToken() {
    return _prefs.getString(_jwtKey);
  }
}
