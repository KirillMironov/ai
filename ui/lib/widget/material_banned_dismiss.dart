import 'package:flutter/material.dart';

class MaterialBannerDismiss extends MaterialBanner {
  final BuildContext context;
  final String text;

  MaterialBannerDismiss(this.context, this.text, {super.key})
      : super(
          content: SelectableText(text),
          actions: [
            TextButton(
              onPressed: () => ScaffoldMessenger.of(context).clearMaterialBanners(),
              child: const Text('DISMISS'),
            ),
          ],
        );

  void show() {
    final messenger = ScaffoldMessenger.of(context);
    messenger.clearMaterialBanners();
    messenger.showMaterialBanner(this);
  }
}
