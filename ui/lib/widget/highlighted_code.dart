import 'package:flutter/material.dart';
import 'package:flutter_highlight/themes/tomorrow-night.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:highlight/highlight.dart' show Highlight, Node;
import 'package:highlight/languages/1c.dart';
import 'package:highlight/languages/ada.dart';
import 'package:highlight/languages/arduino.dart';
import 'package:highlight/languages/armasm.dart';
import 'package:highlight/languages/asciidoc.dart';
import 'package:highlight/languages/aspectj.dart';
import 'package:highlight/languages/avrasm.dart';
import 'package:highlight/languages/awk.dart';
import 'package:highlight/languages/bash.dart';
import 'package:highlight/languages/basic.dart';
import 'package:highlight/languages/brainfuck.dart';
import 'package:highlight/languages/capnproto.dart';
import 'package:highlight/languages/clean.dart';
import 'package:highlight/languages/clojure-repl.dart';
import 'package:highlight/languages/clojure.dart';
import 'package:highlight/languages/cmake.dart';
import 'package:highlight/languages/coffeescript.dart';
import 'package:highlight/languages/cpp.dart';
import 'package:highlight/languages/crystal.dart';
import 'package:highlight/languages/css.dart';
import 'package:highlight/languages/d.dart';
import 'package:highlight/languages/dart.dart';
import 'package:highlight/languages/delphi.dart';
import 'package:highlight/languages/diff.dart';
import 'package:highlight/languages/django.dart';
import 'package:highlight/languages/dns.dart';
import 'package:highlight/languages/dockerfile.dart';
import 'package:highlight/languages/dos.dart';
import 'package:highlight/languages/elixir.dart';
import 'package:highlight/languages/elm.dart';
import 'package:highlight/languages/erlang-repl.dart';
import 'package:highlight/languages/erlang.dart';
import 'package:highlight/languages/excel.dart';
import 'package:highlight/languages/fortran.dart';
import 'package:highlight/languages/fsharp.dart';
import 'package:highlight/languages/glsl.dart';
import 'package:highlight/languages/gml.dart';
import 'package:highlight/languages/go.dart';
import 'package:highlight/languages/golo.dart';
import 'package:highlight/languages/gradle.dart';
import 'package:highlight/languages/groovy.dart';
import 'package:highlight/languages/haskell.dart';
import 'package:highlight/languages/haxe.dart';
import 'package:highlight/languages/http.dart';
import 'package:highlight/languages/ini.dart';
import 'package:highlight/languages/isbl.dart';
import 'package:highlight/languages/java.dart';
import 'package:highlight/languages/javascript.dart';
import 'package:highlight/languages/json.dart';
import 'package:highlight/languages/julia-repl.dart';
import 'package:highlight/languages/julia.dart';
import 'package:highlight/languages/kotlin.dart';
import 'package:highlight/languages/less.dart';
import 'package:highlight/languages/lisp.dart';
import 'package:highlight/languages/llvm.dart';
import 'package:highlight/languages/lua.dart';
import 'package:highlight/languages/makefile.dart';
import 'package:highlight/languages/markdown.dart';
import 'package:highlight/languages/matlab.dart';
import 'package:highlight/languages/mipsasm.dart';
import 'package:highlight/languages/nginx.dart';
import 'package:highlight/languages/nix.dart';
import 'package:highlight/languages/objectivec.dart';
import 'package:highlight/languages/ocaml.dart';
import 'package:highlight/languages/perl.dart';
import 'package:highlight/languages/pgsql.dart';
import 'package:highlight/languages/php.dart';
import 'package:highlight/languages/plaintext.dart';
import 'package:highlight/languages/powershell.dart';
import 'package:highlight/languages/prolog.dart';
import 'package:highlight/languages/protobuf.dart';
import 'package:highlight/languages/puppet.dart';
import 'package:highlight/languages/purebasic.dart';
import 'package:highlight/languages/python.dart';
import 'package:highlight/languages/r.dart';
import 'package:highlight/languages/routeros.dart';
import 'package:highlight/languages/ruby.dart';
import 'package:highlight/languages/rust.dart';
import 'package:highlight/languages/scala.dart';
import 'package:highlight/languages/scss.dart';
import 'package:highlight/languages/shell.dart';
import 'package:highlight/languages/smalltalk.dart';
import 'package:highlight/languages/sql.dart';
import 'package:highlight/languages/swift.dart';
import 'package:highlight/languages/thrift.dart';
import 'package:highlight/languages/typescript.dart';
import 'package:highlight/languages/vala.dart';
import 'package:highlight/languages/vbnet.dart';
import 'package:highlight/languages/vbscript.dart';
import 'package:highlight/languages/verilog.dart';
import 'package:highlight/languages/vhdl.dart';
import 'package:highlight/languages/vim.dart';
import 'package:highlight/languages/x86asm.dart';
import 'package:highlight/languages/xml.dart';
import 'package:highlight/languages/yaml.dart';

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
    final highlight = Highlight()..registerLanguages(highlightedLanguages);

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

  final highlightedLanguages = {
    '1c': lang1C,
    'ada': ada,
    'arduino': arduino,
    'armasm': armasm,
    'asciidoc': asciidoc,
    'aspectj': aspectj,
    'avrasm': avrasm,
    'awk': awk,
    'bash': bash,
    'basic': basic,
    'brainfuck': brainfuck,
    'capnproto': capnproto,
    'clean': clean,
    'clojure-repl': clojureRepl,
    'clojure': clojure,
    'cmake': cmake,
    'coffeescript': coffeescript,
    'cpp': cpp,
    'crystal': crystal,
    'css': css,
    'd': d,
    'dart': dart,
    'delphi': delphi,
    'diff': diff,
    'django': django,
    'dns': dns,
    'dockerfile': dockerfile,
    'dos': dos,
    'elixir': elixir,
    'elm': elm,
    'erlang-repl': erlangRepl,
    'erlang': erlang,
    'excel': excel,
    'fortran': fortran,
    'fsharp': fsharp,
    'glsl': glsl,
    'gml': gml,
    'go': go,
    'golo': golo,
    'gradle': gradle,
    'groovy': groovy,
    'haskell': haskell,
    'haxe': haxe,
    'http': http,
    'ini': ini,
    'isbl': isbl,
    'java': java,
    'javascript': javascript,
    'json': json,
    'julia-repl': juliaRepl,
    'julia': julia,
    'kotlin': kotlin,
    'less': less,
    'lisp': lisp,
    'llvm': llvm,
    'lua': lua,
    'makefile': makefile,
    'markdown': markdown,
    'matlab': matlab,
    'mipsasm': mipsasm,
    'nginx': nginx,
    'nix': nix,
    'objectivec': objectivec,
    'ocaml': ocaml,
    'perl': perl,
    'pgsql': pgsql,
    'php': php,
    'plaintext': plaintext,
    'powershell': powershell,
    'prolog': prolog,
    'protobuf': protobuf,
    'puppet': puppet,
    'purebasic': purebasic,
    'python': python,
    'r': r,
    'routeros': routeros,
    'ruby': ruby,
    'rust': rust,
    'scala': scala,
    'scss': scss,
    'shell': shell,
    'smalltalk': smalltalk,
    'sql': sql,
    'swift': swift,
    'thrift': thrift,
    'typescript': typescript,
    'vala': vala,
    'vbnet': vbnet,
    'vbscript': vbscript,
    'verilog': verilog,
    'vhdl': vhdl,
    'vim': vim,
    'x86asm': x86Asm,
    'xml': xml,
    'yaml': yaml,
  };
}
