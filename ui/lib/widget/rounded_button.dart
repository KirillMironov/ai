import 'package:flutter/material.dart';

class RoundedButton extends StatelessWidget {
  final Widget child;
  final VoidCallback onTap;
  final Color? color;
  final BorderRadius borderRadius;
  final EdgeInsetsGeometry padding;

  RoundedButton({
    super.key,
    required this.child,
    required this.onTap,
    this.color,
    BorderRadius? borderRadius,
    EdgeInsetsGeometry? padding,
  }) :
    borderRadius = borderRadius ?? BorderRadius.circular(12.0),
    padding = padding ?? const EdgeInsets.all(15.0);

  @override
  Widget build(BuildContext context) {
    return Material(
      color: color,
      borderRadius: borderRadius,
      child: InkWell(
        borderRadius: borderRadius,
        onTap: onTap,
        child: Padding(
          padding: padding,
          child: child,
        ),
      ),
    );
  }
}
