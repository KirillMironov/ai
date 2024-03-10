import 'package:ai/model/conversation.dart';
import 'package:ai/router.dart';
import 'package:ai/service/conversations_service.dart';
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
          child: FutureBuilder<List<Conversation>>(
            future: widget.conversationsService.listConversations(0, 50),
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                return snapshot.data == null || snapshot.data!.isEmpty
                    ? const Center(child: Text('No conversations'))
                    : ListView.builder(
                        padding: const EdgeInsets.only(right: 15.0),
                        itemCount: snapshot.data != null ? snapshot.data!.length : 0,
                        itemBuilder: (context, index) {
                          final conversation = snapshot.data?[index];
                          return Padding(
                            padding: const EdgeInsets.symmetric(vertical: 2.0),
                            child: RoundedButton(
                              onTap: () => context.goRouteID(Routes.conversationByID, conversation!.id),
                              child: Text(
                                conversation != null ? conversation.title : '',
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                          );
                        },
                      );
              } else if (snapshot.hasError) {
                return Center(child: SelectableText(snapshot.error.toString()));
              }
              return const CircularProgressIndicator();
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
              : FutureBuilder<Conversation>(
                  future: widget.conversationsService.getConversation(widget.conversationID!),
                  builder: (context, snapshot) {
                    if (snapshot.hasData) {
                      return snapshot.data == null || snapshot.data!.messages.isEmpty
                          ? const Center(child: Text('No messages'))
                          : ListView.builder(
                              padding: const EdgeInsets.only(right: 15.0),
                              itemCount: snapshot.data != null ? snapshot.data!.messages.length : 0,
                              itemBuilder: (context, index) {
                                final message = snapshot.data!.messages[index];
                                return MessageItem(
                                  role: message.role,
                                  content: message.content,
                                );
                              },
                            );
                    } else if (snapshot.hasError) {
                      return Center(child: SelectableText(snapshot.error.toString()));
                    }
                    return const CircularProgressIndicator();
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
