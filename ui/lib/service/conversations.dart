import '../model/conversation.dart';
import '../model/message.dart';

abstract interface class ConversationsService {
  Future<List<Conversation>> listConversations(int offset, int limit);
  Future<Conversation> getConversation(String id);
  Future<Message> sendMessage(String conversationID, String content);
  Stream<Message> sendMessageStream(String conversationID, String content);
}
