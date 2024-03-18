import 'dart:async';
import 'dart:convert';
import 'package:ai/model/token.dart';
import 'package:ai/storage/token_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPreferencesTokenStorage implements TokenStorage {
  final _tokenKey = 'token';

  late SharedPreferences _prefs;

  SharedPreferencesTokenStorage() {
    _initPrefs();
  }

  Future<void> _initPrefs() async {
    _prefs = await SharedPreferences.getInstance();
  }

  @override
  void saveToken(Token token) async {
    final json = jsonEncode(token);
    await _prefs.setString(_tokenKey, json);
  }

  @override
  Token? getToken() {
    final json = _prefs.getString(_tokenKey);
    return json != null ? Token.fromJson(jsonDecode(json)) : null;
  }
}
