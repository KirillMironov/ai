import 'package:ai/model/role.dart';

class Message {
  final String id;
  final Role role;
  String content;
  final DateTime createdAt;
  final DateTime updatedAt;

  Message(this.id, this.role, this.content, this.createdAt, this.updatedAt);
}
