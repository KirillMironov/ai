import 'package:ai/router.dart';
import 'package:ai/service/authenticator.dart';
import 'package:ai/storage/token.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class LoginPage extends StatefulWidget {
  final AuthenticatorService authenticatorService;
  final TokenStorage tokenStorage;

  const LoginPage({
    super.key,
    required this.authenticatorService,
    required this.tokenStorage,
  });

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
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
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    _textFormField('Username', false, _usernameController, 'Empty username'),
                    const SizedBox(height: 15.0),
                    _textFormField('Password', true, _passwordController, 'Empty password'),
                    const SizedBox(height: 15.0),
                    Row(
                      children: [
                        _button('Sign In', () => _login(widget.authenticatorService.signIn)),
                        const SizedBox(width: 10.0),
                        _button('Sign Up', () => _login(widget.authenticatorService.signUp)),
                      ],
                    )
                  ],
                ),
              ),
            ),
          ),
        )
      ),
    );
  }

  Widget _textFormField(String hint, bool obscureText, TextEditingController controller, String validationError) {
    return TextFormField(
      controller: controller,
      obscureText: obscureText,
      decoration: InputDecoration(
        hintText: hint,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12.0),
        ),
      ),
      validator: (value) {
        if (value == null || value.isEmpty) {
          return validationError;
        }
        return null;
      },
    );
  }

  Widget _button(String text, VoidCallback onTap) {
    return Expanded(
      child: RoundedButton(
        color: const Color.fromRGBO(31, 31, 31, 1.0),
        onTap: onTap,
        child: Center(
          child: Text(
            text,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ),
    );
  }

  void _login(Future<String> Function(String username, String password) loginAction) async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    try {
      final token = await loginAction(_usernameController.text, _passwordController.text);
      widget.tokenStorage.saveToken(token);
    } catch(e) {
      if (!mounted) return;
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

    if (!mounted) return;
    context.pushNamed(Routes.conversations.name);
  }
}
