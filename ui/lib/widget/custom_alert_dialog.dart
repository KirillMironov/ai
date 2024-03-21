import 'package:flutter/material.dart';

class CustomAlertDialog extends StatelessWidget {
  final String title;
  final VoidCallback? onCancel;
  final VoidCallback? onContinue;

  const CustomAlertDialog({
    super.key,
    required this.title,
    this.onCancel,
    this.onContinue,
  });

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(title),
      actions: [
        TextButton(
          onPressed: onCancel != null ? onCancel! : () => Navigator.of(context).pop(),
          child: const Text('Cancel'),
        ),
        TextButton(
          onPressed: onContinue,
          child: const Text('Continue'),
        )
      ],
    );
  }
}
