import 'package:ai/page/conversations_page.dart';
import 'package:ai/page/login_page.dart';
import 'package:ai/router.dart';
import 'package:ai/service/grpc_authenticator_service.dart';
import 'package:ai/service/grpc_conversations_service.dart';
import 'package:ai/storage/shared_preferences_token_storage.dart';
import 'package:flutter/material.dart' hide Router;
import 'package:go_router/go_router.dart';

void main() {
  // --dart-define args
  const aiHost = String.fromEnvironment('AI_HOST', defaultValue: 'localhost');
  const aiPort = int.fromEnvironment('AI_PORT', defaultValue: 8080);
  const aiWebPort = int.fromEnvironment('AI_WEB_PORT', defaultValue: 9090);
  const aiSecure = bool.fromEnvironment('AI_SECURE', defaultValue: false);

  // routing
  WidgetsFlutterBinding.ensureInitialized();
  GoRouter.optionURLReflectsImperativeAPIs = true;
  setUrlStrategy();

  // di
  final tokenStorage = SharedPreferencesTokenStorage();
  final authenticatorService = GrpcAuthenticatorService(aiHost, aiPort, aiWebPort, aiSecure);
  final conversationsService = GrpcConversationsService(aiHost, aiPort, aiWebPort, aiSecure, tokenStorage);
  final loginPage = LoginPage(authenticatorService: authenticatorService, tokenStorage: tokenStorage);
  final conversationsPage = ConversationsPage(conversationsService: conversationsService);
  final router = Router(tokenStorage, conversationsService, loginPage, conversationsPage);

  runApp(App(router: router.router()));
}

class App extends StatelessWidget {
  final GoRouter router;

  const App({
    super.key,
    required this.router,
  });

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'AI',
      theme: ThemeData.dark(useMaterial3: true),
      debugShowCheckedModeBanner: false,
      routerConfig: router,
    );
  }
}
