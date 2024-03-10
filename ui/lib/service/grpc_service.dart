import 'package:flutter/material.dart';
import 'package:grpc/grpc_or_grpcweb.dart';

base class GrpcService {
  final String host;
  final int port;
  final int webPort;
  final bool secure;

  GrpcService(this.host, this.port, this.webPort, this.secure);

  @protected
  GrpcOrGrpcWebClientChannel createChannel() {
    return GrpcOrGrpcWebClientChannel.toSeparatePorts(
      host: host,
      grpcPort: port,
      grpcTransportSecure: secure,
      grpcWebPort: webPort,
      grpcWebTransportSecure: secure,
    );
  }

  @protected
  Exception handleException(Object e, String message) {
    switch (e) {
      case final GrpcError e:
        return Exception('$message: ${e.message} (code: ${e.code})');
      default:
        return Exception('$message: $e');
    }
  }
}
