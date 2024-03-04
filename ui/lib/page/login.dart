import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';

class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

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
                  _textField('Username'),
                  const SizedBox(height: 15.0),
                  _textField('Password'),
                  const SizedBox(height: 15.0),
                  Row(
                    children: [
                      _button('Sign In'),
                      const SizedBox(width: 10.0),
                      _button('Sign Up'),
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

  Widget _textField(String hint) {
    return TextField(
      decoration: InputDecoration(
        hintText: hint,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12.0),
        ),
      ),
    );
  }

  Widget _button(String text) {
    return Expanded(
      child: RoundedButton(
        color: const Color.fromRGBO(31, 31, 31, 1.0),
        child: Center(
          child: Text(
            text,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        onTap: () {},
      ),
    );
  }
}
