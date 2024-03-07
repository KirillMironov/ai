import 'dart:async';
import 'package:ai/model/conversation.dart';
import 'package:ai/model/message.dart';
import 'package:ai/model/role.dart';
import 'package:ai/service/conversations_service.dart';
import 'package:ai/service/grpc_service.dart';
import 'package:ai/api/ai.pbgrpc.dart' as api;

final class GrpcConversationsService extends GrpcService implements ConversationsService {
  GrpcConversationsService(super.host, super.port, super.webPort, super.secure);

  @override
  Future<List<Conversation>> listConversations(int offset, int limit) async {
    final channel = createChannel();
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.listConversations(api.ListConversationsRequest()
      ..limit = limit
      ..offset = offset
      );

      return response.conversations.map((e) => Conversation(
        e.id,
        e.title,
        List.empty(),
        e.createdAt.toDateTime(),
        e.updatedAt.toDateTime()
      )).toList();
    } catch (e) {
      throw handleException(e, 'failed to list conversations');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<Conversation> getConversation(String id) async {
    final channel = createChannel();
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.getConversation(api.GetConversationRequest(id: id));
      final conversation = response.conversation;

      return Conversation(
          conversation.id,
          conversation.title,
          response.messages.map((e) => Message(
              e.id,
              e.role.toRole(),
              e.content,
              e.createdAt.toDateTime(),
              e.updatedAt.toDateTime(),
          )).toList(),
          conversation.createdAt.toDateTime(),
          conversation.updatedAt.toDateTime()
      );
    } catch (e) {
      throw handleException(e, 'failed to get conversation');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Future<Message> sendMessage(String conversationId, String content) async {
    final channel = createChannel();
    final client = api.ConversationsClient(channel);

    try {
      final response = await client.sendMessage(api.SendMessageRequest()
        ..conversationId = conversationId
        ..content = content
      );
      final message = response.message;

      return Message(
          message.id,
          message.role.toRole(),
          message.content,
          message.createdAt.toDateTime(),
          message.updatedAt.toDateTime(),
      );
    } catch (e) {
      throw handleException(e, 'failed to send message');
    } finally {
      await channel.shutdown();
    }
  }

  @override
  Stream<Message> sendMessageStream(String conversationId, String content) {
    final channel = createChannel();
    final client = api.ConversationsClient(channel);

    try {
      final request = api.SendMessageStreamRequest()
        ..conversationId = conversationId
        ..content = content;

      return client.sendMessageStream(request).map((event) => Message(
          event.message.id,
          event.message.role.toRole(),
          event.message.content,
          event.message.createdAt.toDateTime(),
          event.message.updatedAt.toDateTime(),
      ));
    } catch (e) {
      throw handleException(e, 'failed to send message stream');
    }
  }
}

extension APIRoleToRole on api.Role {
  Role toRole() {
    switch (this) {
      case api.Role.ROLE_ASSISTANT:
        return Role.assistant;
      case api.Role.ROLE_USER:
        return Role.user;
      default:
        throw Exception('unexpected role: $this');
    }
  }
}