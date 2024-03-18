import 'package:ai/model/user.dart';

abstract interface class UserStorage {
  void saveUser(User user);
  User? getUser();
}
