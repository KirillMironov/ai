import 'package:ai/api/ai.pbgrpc.dart' as api;
import 'package:ai/service/authenticator_service.dart';
import 'package:ai/service/grpc_service.dart';

final class GrpcAuthenticatorService extends GrpcService implements AuthenticatorService {
  GrpcAuthenticatorService(super.host, super.webHost, super.port, super.webPort, super.secure);

  @override
  Future<String> signUp(String username, String password) async {
    final channel = createChannel();
    final client = api.AuthenticatorClient(channel);

    try {
      final response = await client.signUp(api.SignUpRequest(
        username: username,
        password: password,
      ));

      return response.token;
    } catch (e) {
      throw handleException(e, 'failed to sign up');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<String> signIn(String username, String password) async {
    final channel = createChannel();
    final client = api.AuthenticatorClient(channel);

    try {
      final response = await client.signIn(api.SignInRequest(
        username: username,
        password: password,
      ));

      return response.token;
    } catch (e) {
      throw handleException(e, 'failed to sign in');
    } finally {
      await channel.shutdown();
    }
  }
}
