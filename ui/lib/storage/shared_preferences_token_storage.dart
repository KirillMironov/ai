import 'dart:async';
import 'package:ai/storage/token_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPreferencesTokenStorage implements TokenStorage {
  late SharedPreferences _prefs;
  final _jwtKey = 'jwt';

  SharedPreferencesTokenStorage() {
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
