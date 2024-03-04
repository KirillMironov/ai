import 'package:ai/model/conversation.dart';
import 'package:ai/model/message.dart';
import 'package:ai/model/role.dart';
import 'package:ai/widget/message_item.dart';
import 'package:ai/widget/rounded_button.dart';
import 'package:flutter/material.dart';

class ConversationsPage extends StatefulWidget {
  const ConversationsPage({super.key});

  @override
  State<ConversationsPage> createState() => _ConversationsPageState();
}

class _ConversationsPageState extends State<ConversationsPage> {
  final buttonColor = const Color.fromRGBO(31, 31, 31, 1.0);

  List<Conversation> conversations = [
    Conversation("85a004ba-9006-4b69-86fe-513ce477ead3", "user_id", "What is th", DateTime.timestamp(), DateTime.timestamp()),
    Conversation("85a004ba-9006-4b69-86fe-513ce477ead4", "user_id", "What is th", DateTime.timestamp(), DateTime.timestamp()),
  ];

  List<Message> messages = [
    Message("51f1e357-d7b6-48a2-b783-19d3bf7f3e6d", Role.user, "What is the most populated city in Colombia?", DateTime.timestamp(), DateTime.timestamp()),
    Message("234709b5-1ecc-431f-90e4-8bae379fb00b", Role.assistant, "The most populated city in Colombia is Bogotá, with a population of approximately 8.8 million people according to the National Administrative Department of Statistics (DANE) in 2020. Bogotá is not only the capital and largest city of Colombia but also one of the most important cultural, economic, and political centers in Latin America.", DateTime.timestamp(), DateTime.timestamp()),
    Message("2c58e1e5-44e1-4cf4-bfbe-7f8597f53036", Role.user, "Thanks! I\u0027m planning a trip there and wanted to know more about the local attractions.", DateTime.timestamp(), DateTime.timestamp()),
    Message("bb380b1d-796a-4763-8c0a-6fb2c59fc9e5", Role.assistant, "Bogotá has a lot to offer for tourists! Some of the most popular attractions include:\n1. La Candelaria: This is the historic center of Bogotá, where you can find colonial-style buildings, churches, and museums.\n2. Plaza Bolivar: Located in the heart of La Candelaria, this square is surrounded by important buildings such as the Cathedral of Bogotá and the Colombian Supreme Court.\n3. Paloquemao Market: This market is a great place to try local food and drinks, buy souvenirs, and experience the city\u0027s vibrant culture.\n4. Monserrate: This mountain is located just outside of Bogotá and offers stunning views of the city and surrounding landscape. You can take a funicular or cable car to the top.\n5. Museo del Oro: The Gold Museum is one of the most important museums in Colombia, featuring an extensive collection of pre-Columbian gold artifacts from various indigenous cultures.\n6. Parque 93: This park is a popular spot for locals and tourists alike, offering a peaceful escape from the hustle and bustle of the city. You can also find several restaurants and cafes here.\n7. Teatro Colón: This opera house is one of the most iconic buildings in Bogotá and hosts various cultural events throughout the year.\n8. La Macarena: This neighborhood is known for its trendy bars, restaurants, and nightlife. It\u0027s a great place to experience the city\u0027s vibrant nightlife.\nThese are just a few of the many attractions that Bogotá has to offer. Depending on your interests, you may also want to visit other nearby sites such as the Salt Cathedral in Zipaquirá or the coffee region of Risaralda.\n", DateTime.timestamp(), DateTime.timestamp()),
  ];

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        return constraints.maxWidth >= 600 ? _buildDesktop() : _buildMobile();
      }
    );
  }

  Widget _buildDesktop() {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(10.0),
          child: Row(
            children: [
              Expanded(
                  flex: 1,
                  child: _buildConversations()
              ),
              Expanded(
                  flex: 4,
                  child: _buildMessages()
              ),
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
          onTap: () {},
          color: buttonColor,
          child: const Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Text(
                  'Start Conversation',
                  style: TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold
                  ),
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
          child: ListView.builder(
            padding: const EdgeInsets.only(right: 15.0),
            itemCount: conversations.length,
            itemBuilder: (context, index) {
              final conversation = conversations[index];
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 2.0),
                child: RoundedButton(
                  onTap: () {},
                  child: Text(
                      conversation.title,
                      overflow: TextOverflow.ellipsis,
                  ),
                ),
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
      ]
    );
  }

  Widget _buildMessages() {
    return Column(
      children: [
        Expanded(
          child: ListView.builder(
            itemCount: messages.length,
            itemBuilder: (context, index) {
              final message = messages[index];
              return MessageItem(
                role: message.role,
                content: message.content,
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
