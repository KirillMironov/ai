import 'package:ai/model/token.dart';

abstract interface class TokenStorage {
  void saveToken(Token token);
  Token? getToken();
}
