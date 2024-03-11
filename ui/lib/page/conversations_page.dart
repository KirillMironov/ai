import 'package:ai/model/conversation.dart';
import 'package:ai/router.dart';
import 'package:ai/service/conversations_service.dart';
import 'package:ai/widget/custom_future_builder.dart';
import 'package:ai/widget/message_item.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';

class ConversationsPage extends StatefulWidget {
  final ConversationsService conversationsService;
  final String? conversationID;

  const ConversationsPage({
    super.key,
    required this.conversationsService,
    this.conversationID,
  });

  @override
  State<ConversationsPage> createState() => _ConversationsPageState();
}

class _ConversationsPageState extends State<ConversationsPage> {
  final buttonColor = const Color.fromRGBO(31, 31, 31, 1.0);

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        return constraints.maxWidth >= 600 ? _buildDesktop() : _buildMobile();
      },
    );
  }

  Widget _buildDesktop() {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(10.0),
          child: Row(
            children: [
              Expanded(flex: 1, child: _buildConversations()),
              Expanded(flex: 4, child: _buildMessages()),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMobile() {
    return Scaffold(
      appBar: AppBar(
        title: const Text('AI'),
      ),
      drawer: Drawer(
        child: Padding(
          padding: const EdgeInsets.all(5.0),
          child: _buildConversations(),
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(10.0),
          child: _buildMessages(),
        ),
      ),
    );
  }

  Widget _buildConversations() {
    return Column(
      children: [
        RoundedButton(
          onTap: () => context.goRoute(Routes.conversations),
          color: buttonColor,
          child: const Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Text(
                  'Start Conversation',
                  style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              Icon(
                Icons.edit_note,
                color: Colors.white,
              )
            ],
          ),
        ),
        const SizedBox(height: 10.0),
        Expanded(
          child: CustomFutureBuilder<List<Conversation>>(
            future: widget.conversationsService.listConversations(0, 50),
            builder: (conversations) {
              return conversations.isEmpty
                  ? const Center(child: Text('No conversations'))
                  : ListView.builder(
                      padding: const EdgeInsets.only(right: 15.0),
                      itemCount: conversations.length,
                      itemBuilder: (context, index) {
                        final conversation = conversations[index];
                        return Padding(
                          padding: const EdgeInsets.symmetric(vertical: 2.0),
                          child: RoundedButton(
                            onTap: () => context.goRouteID(Routes.conversationByID, conversation.id),
                            child: Text(
                              conversation.title,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        );
                      },
                    );
            },
          ),
        ),
        const SizedBox(height: 10.0),
        Align(
          alignment: Alignment.bottomCenter,
          child: RoundedButton(
            onTap: () {},
            color: buttonColor,
            child: const Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    'John Doe',
                    style: TextStyle(color: Colors.white),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Icon(
                  Icons.person,
                  color: Colors.white,
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildMessages() {
    return Column(
      children: [
        Expanded(
          child: widget.conversationID == null
              ? const Center(child: Text('How can I help you today?'))
              : CustomFutureBuilder<Conversation>(
                  future: widget.conversationsService.getConversation(widget.conversationID!),
                  builder: (conversation) {
                    return conversation.messages.isEmpty
                        ? const Center(child: Text('No messages'))
                        : ListView.builder(
                            padding: const EdgeInsets.only(right: 15.0),
                            itemCount: conversation.messages.length,
                            itemBuilder: (context, index) {
                              final message = conversation.messages[index];
                              return MessageItem(
                                role: message.role,
                                content: message.content,
                              );
                            },
                          );
                  },
                ),
        ),
        TextField(
          decoration: InputDecoration(
            hintText: 'Type a message...',
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12.0),
            ),
            suffixIcon: Padding(
              padding: const EdgeInsets.only(right: 10.0),
              child: IconButton(
                icon: const Icon(Icons.send),
                onPressed: () {},
              ),
            ),
          ),
        ),
      ],
    );
  }
}
