import 'package:ai/model/conversation.dart';
import 'package:ai/model/message.dart';

abstract interface class ConversationsService {
  Future<List<ConversationDescription>> listConversations(int offset, int limit);
  Future<List<Message>> getMessagesByConversationID(String conversationID);
  Future<void> deleteConversationByID(String conversationID);
  Future<Message> sendMessage(String conversationID, String content);
  Stream<Message> sendMessageStream(String conversationID, String content);
}
