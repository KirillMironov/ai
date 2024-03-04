enum Role {
  assistant,
  user;

  @override
  String toString() {
    switch (this) {
      case Role.assistant:
        return 'Assistant';
      case Role.user:
        return 'User';
      default:
        return name;
    }
  }
}
