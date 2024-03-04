abstract interface class AuthenticatorService {
  Future<String> signUp(String username, password);
  Future<String> signIn(String username, password);
}
