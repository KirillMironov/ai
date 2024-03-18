import 'dart:async';
import 'dart:convert';
import 'package:ai/model/user.dart';
import 'package:ai/storage/user_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPreferencesUserStorage implements UserStorage {
  final _tokenKey = 'token';

  late SharedPreferences _prefs;

  SharedPreferencesUserStorage() {
    _initPrefs();
  }

  Future<void> _initPrefs() async {
    _prefs = await SharedPreferences.getInstance();
  }

  @override
  void saveUser(User user) async {
    final json = jsonEncode(user);
    await _prefs.setString(_tokenKey, json);
  }

  @override
  User? getUser() {
    final json = _prefs.getString(_tokenKey);
    return json != null ? User.fromJson(jsonDecode(json)) : null;
  }
}
