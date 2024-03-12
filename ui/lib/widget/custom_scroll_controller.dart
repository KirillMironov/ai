import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';

class CustomScrollController extends ScrollController {
  final Duration duration;
  final Curve curve;

  CustomScrollController({
    Duration? duration,
    Curve? curve,
  })  : duration = duration ?? const Duration(milliseconds: 200),
        curve = curve ?? Curves.fastOutSlowIn;

  void scrollDown() {
    if (positions.isEmpty) return;
    SchedulerBinding.instance.addPostFrameCallback((_) {
      animateTo(
        position.maxScrollExtent,
        duration: duration,
        curve: curve,
      );
    });
  }
}
