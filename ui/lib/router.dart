import 'package:ai/page/conversations.dart';
import 'package:ai/page/login.dart';
import 'package:ai/storage/token.dart';
import 'package:go_router/go_router.dart';

enum Routes {
  conversations,
  login;

  String get name {
    switch (this) {
      case Routes.conversations:
        return 'conversations';
      case Routes.login:
        return 'login';
    }
  }

  String get path {
    switch (this) {
      case Routes.conversations:
        return '/';
      case Routes.login:
        return '/login';
    }
  }
}

class Router {
  final TokenStorage tokenStorage;
  final LoginPage loginPage;
  final ConversationsPage conversationsPage;

  Router(this.tokenStorage, this.loginPage, this.conversationsPage);

  GoRouter router() => GoRouter(
    initialLocation: Routes.login.path,
    routes: [
      GoRoute(
        name: Routes.conversations.name,
        path: Routes.conversations.path,
        builder: (context, state) => conversationsPage,
      ),
      GoRoute(
        name: Routes.login.name,
        path: Routes.login.path,
        builder: (context, state) => loginPage,
      )
    ],
    redirect: (context, state) {
      final conversationsLocation = state.namedLocation(Routes.conversations.name);
      final loginLocation = state.namedLocation(Routes.login.name);

      try {
        return tokenStorage.getToken() != null
            ? conversationsLocation
            : loginLocation;
      } catch(_) {}

      return null;
    },
  );
}
