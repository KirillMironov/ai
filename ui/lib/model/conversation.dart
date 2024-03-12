import 'package:ai/model/message.dart';

class Conversation {
  final String id;
  final String title;
  final List<Message> messages;
  final DateTime createdAt;
  final DateTime updatedAt;

  Conversation(this.id, this.title, this.messages, this.createdAt, this.updatedAt);
}

class ConversationDescription {
  final String id;
  final String title;
  final DateTime createdAt;
  final DateTime updatedAt;

  ConversationDescription(this.id, this.title, this.createdAt, this.updatedAt);
}
