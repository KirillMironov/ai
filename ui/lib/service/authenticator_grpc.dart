import 'package:ai/service/authenticator.dart';
import 'package:ai/api/ai.pbgrpc.dart' as api;
import 'package:grpc/grpc.dart';

class AuthenticatorServiceGRPC implements AuthenticatorService {
  final String host;
  final int port;

  AuthenticatorServiceGRPC(this.host, this.port);

  @override
  Future<String> signUp(String username, String password) async {
    final channel = ClientChannel(host, port: port, options: const ChannelOptions(credentials: ChannelCredentials.insecure()));
    final client = api.AuthenticatorClient(channel);

    try {
      final response = await client.signUp(api.SignUpRequest()
        ..username = username
        ..password = password);

      return response.token;
    } catch (e) {
      throw Exception('failed to sign up: $e');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<String> signIn(String username, String password) async {
    final channel = ClientChannel(host, port: port, options: const ChannelOptions(credentials: ChannelCredentials.insecure()));
    final client = api.AuthenticatorClient(channel);

    try {
      final response = await client.signIn(api.SignInRequest()
        ..username = username
        ..password = password);

      return response.token;
    } catch (e) {
      throw Exception('failed to sign in: $e');
    } finally {
      await channel.shutdown();
    }
  }
}
