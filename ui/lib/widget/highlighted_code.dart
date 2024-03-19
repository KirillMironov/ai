import 'package:flutter/material.dart';
import 'package:flutter_highlight/themes/tomorrow-night.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:highlight/highlight.dart' show highlight, Node;

class HighlightedCode extends StatelessWidget {
  final String code;
  final EdgeInsetsGeometry? padding;
  final TextStyle textStyle;
  final Map<String, TextStyle> theme;

  HighlightedCode({
    super.key,
    required this.code,
    this.padding,
    TextStyle? textStyle,
    Map<String, TextStyle>? theme,
  })  : textStyle = textStyle ?? GoogleFonts.jetBrainsMono(),
        theme = theme ?? tomorrowNightTheme;

  @override
  Widget build(BuildContext context) {
    return Container(
      color: theme['root']?.backgroundColor ?? const Color(0xffffffff),
      padding: padding,
      child: SelectableText.rich(
        TextSpan(
          style: textStyle,
          children: _convert(highlight.parse(code, autoDetection: true).nodes!),
        ),
      ),
    );
  }

  List<TextSpan> _convert(List<Node> nodes) {
    List<TextSpan> spans = [];
    List<List<TextSpan>> stack = [];

    var currentSpans = spans;

    void traverse(Node node) {
      if (node.value != null) {
        currentSpans.add(node.className == null
            ? TextSpan(text: node.value)
            : TextSpan(text: node.value, style: theme[node.className!]));
      } else if (node.children != null) {
        List<TextSpan> tmp = [];
        currentSpans.add(TextSpan(children: tmp, style: theme[node.className!]));
        stack.add(currentSpans);
        currentSpans = tmp;

        for (var n in node.children!) {
          traverse(n);
          if (n == node.children!.last) {
            currentSpans = stack.isEmpty ? spans : stack.removeLast();
          }
        }
      }
    }

    for (var node in nodes) {
      traverse(node);
    }

    return spans;
  }
}
