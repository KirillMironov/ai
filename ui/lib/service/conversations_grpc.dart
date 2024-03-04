import 'dart:async';
import 'package:ai/model/conversation.dart';
import 'package:ai/model/message.dart';
import 'package:ai/model/role.dart';
import 'package:ai/service/conversations.dart';
import 'package:fixnum/fixnum.dart';
import 'package:grpc/grpc.dart';
import 'package:ai/api/ai.pbgrpc.dart' as api;

class ConversationsServiceGRPC implements ConversationsService {
  final String host;
  final int port;

  ConversationsServiceGRPC(this.host, this.port);

  @override
  Future<List<Conversation>> listConversations(int offset, int limit) async {
    final channel = ClientChannel(host, port: port);
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.listConversations(api.ListConversationsRequest(offset: offset, limit: limit));

      return response.conversations.map((e) => Conversation(
        e.id,
        e.title,
        List.empty(),
        e.createdAt.toDateTime(),
        e.updatedAt.toDateTime()
      )).toList();
    } catch (e) {
      throw Exception('failed to list conversations: $e');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<Conversation> getConversation(String id) async {
    final channel = ClientChannel(host, port: port);
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.getConversation(api.GetConversationRequest(id: id));
      final conversation = response.conversation;

      return Conversation(
          conversation.id,
          conversation.title,
          response.messages.map((e) => Message(
              e.id,
              e.role == 'user' ? Role.user : Role.assistant,
              e.content,
              DateTime.now(),
              DateTime.now(),
          )).toList(),
          conversation.createdAt.toDateTime(),
          conversation.updatedAt.toDateTime()
      );
    } catch (e) {
      throw Exception('failed to get conversation: $e');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<Message> sendMessage(String conversationId, String role, String content) async {
    final channel = ClientChannel(host, port: port);
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.sendMessage(api.SendMessageRequest()
        ..conversationId = conversationId
        ..role = role
        ..content = content);
      final message = response.message;

      return Message(
          message.id,
          message.role == 'user' ? Role.user : Role.assistant,
          message.content,
          DateTime.now(),
          DateTime.now(),
      );
    } catch (e) {
      throw Exception('failed to send message: $e');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Stream<Message> sendMessageStream(String conversationId, String role, String content) {
    final channel = ClientChannel(host, port: port);
    final client = api.ConversationsClient(channel);

    try {
      final request = api.SendMessageStreamRequest()
        ..conversationId = conversationId
        ..role = role
        ..content = content;

      return client.sendMessageStream(request).map((event) => Message(
          event.message.id,
          event.message.role == 'user' ? Role.user : Role.assistant,
          event.message.content,
          DateTime.now(),
          DateTime.now()
      ));
    } catch (e) {
      throw Exception('failed to send message stream: $e');
    }
  }
}

extension Int64ToDateTime on Int64 {
  DateTime toDateTime() {
    return DateTime.fromMillisecondsSinceEpoch(1000 * toInt());
  }
}
