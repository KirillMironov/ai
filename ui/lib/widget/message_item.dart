import 'package:ai/model/role.dart';
import 'package:flutter/material.dart';

class MessageItem extends StatelessWidget {
  final Role role;
  final String content;

  const MessageItem({
    super.key,
    required this.role,
    required this.content,
  });

  @override
  Widget build(BuildContext context) {
    final avatarBackgroundColor = role == Role.user ? Colors.purple : Colors.teal;
    final avatarIcon = role == Role.user ? Icons.person : Icons.assistant;
    final padding = MediaQuery.of(context).size.width < 600
        ? EdgeInsets.zero
        : const EdgeInsets.symmetric(horizontal: 150.0);

    return Container(
      padding: padding,
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              CircleAvatar(
                backgroundColor: avatarBackgroundColor,
                child: Icon(avatarIcon),
              ),
              const SizedBox(width: 8.0),
              Expanded(
                child: Text(
                  role.toString(),
                  style: const TextStyle(fontWeight: FontWeight.bold),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: 4.0),
          Padding(
            padding: const EdgeInsets.only(left: 48.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SelectableText(
                  content,
                  textAlign: TextAlign.left,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
