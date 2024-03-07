import 'package:ai/page/conversations_page.dart';
import 'package:ai/page/login_page.dart';
import 'package:ai/storage/token_storage.dart';
import 'package:go_router/go_router.dart';
import 'package:ai/router.dart'
  if (dart.library.html) 'package:flutter_web_plugins/flutter_web_plugins.dart'
  as plugins;

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

/// Use url path strategy only on web to successfully build
/// on other platforms using conditional import
void setUrlStrategy() {
  plugins.usePathUrlStrategy();
}

void usePathUrlStrategy() {}
