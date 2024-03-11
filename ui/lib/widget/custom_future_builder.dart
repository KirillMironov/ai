import 'package:flutter/material.dart';

class CustomFutureBuilder<T> extends StatelessWidget {
  final Future<T> future;
  final Widget Function(T) builder;

  const CustomFutureBuilder({super.key, required this.future, required this.builder});

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<T>(
      future: future,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          return builder(snapshot.data as T);
        }
        return Center(
          child: snapshot.hasError
              ? SelectableText(snapshot.error.toString())
              : const CircularProgressIndicator(),
        );
      },
    );
  }
}
