import 'package:ai/page/conversations_page.dart';
import 'package:ai/page/login_page.dart';
import 'package:ai/router.dart';
import 'package:ai/service/grpc_authenticator_service.dart';
import 'package:ai/storage/shared_preferences_token_storage.dart';
import 'package:flutter/material.dart' hide Router;
import 'package:go_router/go_router.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  GoRouter.optionURLReflectsImperativeAPIs = true;
  setUrlStrategy();
  final tokenStorage = SharedPreferencesTokenStorage();
  final authenticatorService = GrpcAuthenticatorService('localhost', 8080, 9090, false);
  final loginPage = LoginPage(authenticatorService: authenticatorService, tokenStorage: tokenStorage);
  const conversationsPage = ConversationsPage();
  final router = Router(tokenStorage, loginPage, conversationsPage);
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
