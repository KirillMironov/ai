//
//  Generated code. Do not modify.
//  source: ai.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class Role extends $pb.ProtobufEnum {
  static const Role ROLE_UNSPECIFIED = Role._(0, _omitEnumNames ? '' : 'ROLE_UNSPECIFIED');
  static const Role ROLE_ASSISTANT = Role._(1, _omitEnumNames ? '' : 'ROLE_ASSISTANT');
  static const Role ROLE_USER = Role._(2, _omitEnumNames ? '' : 'ROLE_USER');

  static const $core.List<Role> values = <Role> [
    ROLE_UNSPECIFIED,
    ROLE_ASSISTANT,
    ROLE_USER,
  ];

  static final $core.Map<$core.int, Role> _byValue = $pb.ProtobufEnum.initByValue(values);
  static Role? valueOf($core.int value) => _byValue[value];

  const Role._($core.int v, $core.String n) : super(v, n);
}


const _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
