import 'package:ai/model/user.dart';

abstract interface class UserStorage {
  void saveUser(User user);
  void deleteUser();
  User? getUser();
}
