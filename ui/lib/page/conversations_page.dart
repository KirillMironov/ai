import 'package:ai/model/conversation.dart';
import 'package:ai/model/message.dart';
import 'package:ai/model/role.dart';
import 'package:ai/router.dart';
import 'package:ai/service/conversations_service.dart';
import 'package:ai/storage/user_storage.dart';
import 'package:ai/widget/custom_future_builder.dart';
import 'package:ai/widget/custom_scroll_controller.dart';
import 'package:ai/widget/material_banned_dismiss.dart';
import 'package:ai/widget/message_item.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';

class ConversationsPage extends StatefulWidget {
  final ConversationsService conversationsService;
  final UserStorage userStorage;
  final String? conversationID;

  const ConversationsPage({
    super.key,
    required this.conversationsService,
    required this.userStorage,
    this.conversationID,
  });

  @override
  State<ConversationsPage> createState() => _ConversationsPageState();
}

class _ConversationsPageState extends State<ConversationsPage> {
  final _buttonColor = const Color.fromRGBO(31, 31, 31, 1.0);
  final _messagesScrollController = CustomScrollController();
  final _messageInputController = TextEditingController();

  late Future<List<Message>> _messagesFuture;

  List<Message> _messages = List.empty(growable: true);
  bool _isSendButtonEnabled = true;
  String _username = '';

  @override
  void initState() {
    if (widget.conversationID != null) {
      _messagesFuture = widget.conversationsService.getMessagesByConversationID(widget.conversationID!);
    }

    try {
      final user = widget.userStorage.getUser();
      if (user != null) {
        _username = user.username;
      }
    } catch(_) {}

    super.initState();
  }

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
          color: _buttonColor,
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
          child: CustomFutureBuilder<List<ConversationDescription>>(
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
                            padding: const EdgeInsets.all(10.0),
                            onTap: () => context.goRouteID(Routes.conversationByID, conversation.id),
                            child: Row(
                              children: [
                                Expanded(
                                  child: Text(
                                    conversation.title,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                ),
                                const SizedBox(width: 5.0),
                                IconButton(
                                  visualDensity: VisualDensity.compact,
                                  iconSize: 17.0,
                                  icon: const Icon(Icons.delete_outline_sharp),
                                  onPressed: () => {
                                    showDialog(
                                      context: context,
                                      builder: (context) {
                                        return Column(
                                          mainAxisAlignment: MainAxisAlignment.center,
                                          children: [
                                            const Text('Delete conversation?'),
                                            TextButton(
                                              onPressed: () => _deleteConversationByID(conversation.id),
                                              child: const Text('Delete'),
                                            )
                                          ],
                                        );
                                      },
                                    )
                                  },
                                )
                              ],
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
            color: _buttonColor,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    _username,
                    style: const TextStyle(color: Colors.white),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                const Icon(
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
          child: widget.conversationID == null && _messages.isEmpty
              ? const Center(child: Text('How can I help you today?'))
              : CustomFutureBuilder<List<Message>>(
                  future: _messagesFuture,
                  builder: (messages) {
                    _messages = messages;
                    return messages.isEmpty
                        ? const Center(child: Text('No messages'))
                        : ListView.builder(
                            controller: _messagesScrollController,
                            padding: const EdgeInsets.only(right: 15.0),
                            itemCount: messages.length,
                            itemBuilder: (context, index) {
                              final message = messages[index];
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
          controller: _messageInputController,
          minLines: 1,
          maxLines: 10,
          decoration: InputDecoration(
            hintText: 'Type a message...',
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12.0),
            ),
            suffixIcon: Padding(
              padding: const EdgeInsets.only(right: 10.0),
              child: _isSendButtonEnabled
                  ? IconButton(icon: const Icon(Icons.send), onPressed: () => _sendMessageStream())
                  : IconButton(icon: const Icon(Icons.send), onPressed: null, color: Colors.grey.withOpacity(0.5)),
            ),
          ),
        ),
      ],
    );
  }

  void _sendMessageStream() async {
    final content = _messageInputController.text;
    if (content.isEmpty) return;

    setState(() {
      _isSendButtonEnabled = false;
      _messages.add(_createMessage(Role.user, content));
      _messages.add(_createMessage(Role.assistant, ''));
      _messagesFuture = Future(() => _messages);
      _messageInputController.clear();
      _messagesScrollController.scrollDown();
    });

    try {
      await widget.conversationsService.sendMessageStream(widget.conversationID ?? '', content).forEach((e) {
        setState(() {
          _messages.last.content += e.content;
          _messagesScrollController.scrollDown();
        });
      });
    } catch (e) {
      if (!mounted) return;
      MaterialBannerDismiss(context, e.toString()).show();
      return;
    } finally {
      setState(() {
        _isSendButtonEnabled = true;
      });
    }

    final conversationID = await widget.conversationsService.listConversations(0, 1).then((value) => value.first.id);

    if (!mounted) return;
    context.goRouteID(Routes.conversationByID, conversationID);
  }

  void _deleteConversationByID(String id) async {
    try {
      await widget.conversationsService.deleteConversationByID(id);
    } catch (e) {
      if (!mounted) return;
      Navigator.of(context).pop();
      MaterialBannerDismiss(context, e.toString()).show();
      return;
    }

    if (widget.conversationID != null && widget.conversationID == id) {
      if (!mounted) return;
      context.goRoute(Routes.conversations);
    } else {
      setState(() {});
      if (!mounted) return;
      Navigator.of(context).pop();
    }
  }

  Message _createMessage(Role role, String content) => Message('', role, content, DateTime.now(), DateTime.now());
}
