import 'package:ai/page/conversations_page.dart';
import 'package:ai/page/login_page.dart';
import 'package:ai/service/conversations_service.dart';
import 'package:ai/storage/token_storage.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:ai/router.dart' if (dart.library.html) 'package:flutter_web_plugins/flutter_web_plugins.dart'
    as plugins;

enum Routes {
  conversations,
  login,
  conversationByID;

  String get path {
    switch (this) {
      case Routes.conversations:
        return '/';
      case Routes.login:
        return '/login';
      case Routes.conversationByID:
        return '/c/:id';
    }
  }
}

class Router {
  final TokenStorage tokenStorage;
  final ConversationsService conversationsService;
  final LoginPage loginPage;
  final ConversationsPage conversationsPage;

  Router(
    this.tokenStorage,
    this.conversationsService,
    this.loginPage,
    this.conversationsPage,
  );

  GoRouter router() => GoRouter(
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
          ),
          GoRoute(
            name: Routes.conversationByID.name,
            path: Routes.conversationByID.path,
            builder: (context, state) {
              final conversationID = state.pathParameters['id'];
              return ConversationsPage(
                conversationsService: conversationsService,
                conversationID: conversationID,
              );
            },
          ),
        ],
        redirect: (context, state) {
          try {
            return tokenStorage.getToken() == null ? Routes.login.path : null;
          } catch (_) {
            return Routes.login.path;
          }
        },
      );
}

/// Use url path strategy only on web to successfully build
/// on other platforms using conditional import
void setUrlStrategy() {
  plugins.usePathUrlStrategy();
}

void usePathUrlStrategy() {}

extension GoRouterHelper on BuildContext {
  void goNamedID(String name, String id) => GoRouter.of(this).goNamed(name, pathParameters: {'id': id});
}
