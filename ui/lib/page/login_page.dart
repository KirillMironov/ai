import 'package:ai/model/user.dart';
import 'package:ai/router.dart';
import 'package:ai/service/authenticator_service.dart';
import 'package:ai/storage/user_storage.dart';
import 'package:ai/widget/material_banner_dismiss.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';

class LoginPage extends StatelessWidget {
  final AuthenticatorService authenticatorService;
  final UserStorage userStorage;

  LoginPage({
    super.key,
    required this.authenticatorService,
    required this.userStorage,
  });

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
                        _button('Sign In', () => _login(context, authenticatorService.signIn)),
                        const SizedBox(width: 10.0),
                        _button('Sign Up', () => _login(context, authenticatorService.signUp)),
                      ],
                    )
                  ],
                ),
              ),
            ),
          ),
        ),
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

  void _login(BuildContext context, Future<String> Function(String username, String password) loginAction) async {
    if (!_formKey.currentState!.validate()) return;

    try {
      final token = await loginAction(_usernameController.text, _passwordController.text);
      userStorage.saveUser(User.fromJWT(token));
    } catch (e) {
      if (!context.mounted) return;
      MaterialBannerDismiss(context, e.toString()).show();
      return;
    }

    if (!context.mounted) return;
    ScaffoldMessenger.of(context).removeCurrentMaterialBanner();
    context.goRoute(Routes.conversations);
  }
}
