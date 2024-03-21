import 'package:ai/widget/code_block.dart';
import 'package:flutter/material.dart';

class HighlightedText extends StatelessWidget {
  final String content;

  const HighlightedText({super.key, required this.content});

  @override
  Widget build(BuildContext context) {
    return SelectionArea(
      child: Text.rich(
        TextSpan(children: _contentToSpans(content)),
      ),
    );
  }

  List<InlineSpan> _contentToSpans(String content) {
    final regex = RegExp(r'```(?:\w+\n)?([^`]+)```|`([^`\n]+)`|([^`]+)');
    final matches = regex.allMatches(content);

    List<InlineSpan> spans = [];

    for (Match match in matches) {
      final code = match.group(1);
      final inlineCode = match.group(2);
      final text = match.group(3);

      if (code != null) {
        spans.add(
          WidgetSpan(
            child: CodeBlock(
              code: code.trim(),
              padding: const EdgeInsets.all(10.0),
            ),
          ),
        );
      } else if (inlineCode != null) {
        spans.add(
          TextSpan(
            text: inlineCode,
            style: TextStyle(backgroundColor: Colors.grey.shade800),
          ),
        );
      } else if (text != null) {
        spans.add(TextSpan(text: text));
      }
    }

    return spans;
  }
}
