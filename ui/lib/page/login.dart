import 'package:ai/router.dart';
import 'package:ai/service/authenticator.dart';
import 'package:ai/storage/token.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class LoginPage extends StatelessWidget {
  final AuthenticatorService authenticatorService;
  final TokenStorage tokenStorage;

  LoginPage({
    super.key,
    required this.authenticatorService,
    required this.tokenStorage,
  });

  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final width = screenWidth > 600 ? 400.0 : screenWidth;

    return SafeArea(
      child: Scaffold(
        body: Padding(
          padding: const EdgeInsets.all(10.0),
          child: Align(
            alignment: Alignment.center,
            child: SizedBox(
              width: width,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  _textField('Username', false, _usernameController),
                  const SizedBox(height: 15.0),
                  _textField('Password', true, _passwordController),
                  const SizedBox(height: 15.0),
                  Row(
                    children: [
                      _button(context, 'Sign In', _signIn),
                      const SizedBox(width: 10.0),
                      _button(context, 'Sign Up', _signUp),
                    ],
                  )
                ],
              ),
            ),
          ),
        )
      ),
    );
  }

  Widget _textField(String hint, bool obscureText, TextEditingController controller) {
    return TextField(
      controller: controller,
      obscureText: obscureText,
      decoration: InputDecoration(
        hintText: hint,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12.0),
        ),
      ),
    );
  }

  Widget _button(BuildContext context, String text, void Function(BuildContext) onTap) {
    return Expanded(
      child: RoundedButton(
        color: const Color.fromRGBO(31, 31, 31, 1.0),
        onTap: () => onTap(context),
        child: Center(
          child: Text(
            text,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ),
    );
  }

  void _signIn(BuildContext context) async {
    try {
      final token = await authenticatorService.signIn(
        _usernameController.text,
        _passwordController.text
      );
      tokenStorage.saveToken(token);
    } catch(e) {
      if (!context.mounted) return;
      ScaffoldMessenger.of(context).showMaterialBanner(
        MaterialBanner(
            content: Text(e.toString()),
            actions: [
              TextButton(
                onPressed: () => ScaffoldMessenger.of(context).clearMaterialBanners(),
                child: const Text('DISMISS'),
              ),
            ],
        )
      );
      return;
    }
    if (!context.mounted) return;
    context.pushNamed(Routes.conversations.name);
  }

  void _signUp(BuildContext context) async {
    try {
      final token = await authenticatorService.signUp(
          _usernameController.text,
          _passwordController.text
      );
      tokenStorage.saveToken(token);
    } catch(e) {
      if (!context.mounted) return;
      ScaffoldMessenger.of(context).showMaterialBanner(
          MaterialBanner(
            content: Text(e.toString()),
            actions: [
              TextButton(
                onPressed: () => ScaffoldMessenger.of(context).clearMaterialBanners(),
                child: const Text('DISMISS'),
              ),
            ],
          )
      );
    }
  }
}
