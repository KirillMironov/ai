import 'package:ai/widget/highlighted_code.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:markdown/markdown.dart' as md;

class HighlightedMarkdown extends StatelessWidget {
  final String content;

  const HighlightedMarkdown(this.content, {super.key});

  @override
  Widget build(BuildContext context) {
    return SelectionArea(
      child: MarkdownBody(
        data: content,
        builders: {'code': CodeMarkdownElementBuilder()},
      ),
    );
  }
}

class CodeMarkdownElementBuilder extends MarkdownElementBuilder {
  @override
  Widget? visitElementAfter(md.Element element, TextStyle? preferredStyle) {
    final scrollController = ScrollController();

    return element.textContent.contains('\n')
        ? LayoutBuilder(
            builder: (context, constraints) {
              return Scrollbar(
                controller: scrollController,
                trackVisibility: true,
                child: SingleChildScrollView(
                  controller: scrollController,
                  scrollDirection: Axis.horizontal,
                  child: ConstrainedBox(
                    constraints: BoxConstraints(minWidth: constraints.maxWidth),
                    child: HighlightedCode(
                      code: element.textContent.trimRight(),
                      padding: const EdgeInsets.all(10.0),
                    ),
                  ),
                ),
              );
            },
          )
        : Container(
            color: Colors.grey.shade800,
            padding: const EdgeInsets.all(1.0),
            child: Text(
              element.textContent.trimRight(),
              style: GoogleFonts.jetBrainsMono(),
            ),
          );
  }
}
