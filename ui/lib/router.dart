import 'package:ai/page/conversations_page.dart';
import 'package:ai/page/login_page.dart';
import 'package:ai/service/authenticator_service.dart';
import 'package:ai/service/conversations_service.dart';
import 'package:ai/storage/token_storage.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
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
  final AuthenticatorService authenticatorService;
  final ConversationsService conversationsService;

  Router(this.tokenStorage, this.authenticatorService, this.conversationsService);

  GoRouter router() => GoRouter(
        routes: [
          Route(
            route: Routes.conversations,
            childBuilder: (state) => ConversationsPage(conversationsService: conversationsService),
          ),
          Route(
            route: Routes.login,
            childBuilder: (state) => LoginPage(authenticatorService: authenticatorService, tokenStorage: tokenStorage),
          ),
          Route(
            route: Routes.conversationByID,
            childBuilder: (state) {
              return ConversationsPage(
                conversationsService: conversationsService,
                conversationID: state.pathParameterID(),
              );
            },
          ),
        ],
        redirect: (context, state) {
          try {
            final token = tokenStorage.getToken();
            return token == null || JwtDecoder.isExpired(token.jwt)
                ? Routes.login.path
                : state.matchedLocation == Routes.login.path
                    ? Routes.conversations.path
                    : null;
          } catch (_) {
            return Routes.login.path;
          }
        },
      );
}

class Route extends GoRoute {
  final Routes route;
  final Function(GoRouterState) childBuilder;

  Route({required this.route, required this.childBuilder})
      : super(
          name: route.name,
          path: route.path,
          pageBuilder: (context, state) => CustomTransitionPage(
            key: UniqueKey(),
            child: childBuilder(state),
            transitionsBuilder: (context, animation, secondaryAnimation, child) {
              return FadeTransition(
                opacity: CurveTween(curve: Curves.easeInOutCirc).animate(animation),
                child: child,
              );
            },
          ),
        );
}

/// Use url path strategy only on web to successfully build
/// on other platforms using conditional import
void setUrlStrategy() {
  plugins.usePathUrlStrategy();
}

void usePathUrlStrategy() {}

extension GoRouterHelper on BuildContext {
  void goRoute(Routes route) => GoRouter.of(this).goNamed(route.name);

  void goRouteID(Routes route, String id) => GoRouter.of(this).goNamed(route.name, pathParameters: {'id': id});
}

extension GoRouterStateHelper on GoRouterState {
  String? pathParameterID() => pathParameters['id'];
}
