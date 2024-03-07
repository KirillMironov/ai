abstract interface class AuthenticatorService {
  Future<String> signUp(String username, String password);
  Future<String> signIn(String username, String password);
}
