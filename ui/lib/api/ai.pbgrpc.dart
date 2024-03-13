//
//  Generated code. Do not modify.
//  source: ai.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'ai.pb.dart' as $0;
import 'google/protobuf/empty.pb.dart' as $1;

export 'ai.pb.dart';

@$pb.GrpcServiceName('ai.Authenticator')
class AuthenticatorClient extends $grpc.Client {
  static final _$signUp = $grpc.ClientMethod<$0.SignUpRequest, $0.SignUpResponse>(
      '/ai.Authenticator/SignUp',
      ($0.SignUpRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.SignUpResponse.fromBuffer(value));
  static final _$signIn = $grpc.ClientMethod<$0.SignInRequest, $0.SignInResponse>(
      '/ai.Authenticator/SignIn',
      ($0.SignInRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.SignInResponse.fromBuffer(value));

  AuthenticatorClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.SignUpResponse> signUp($0.SignUpRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$signUp, request, options: options);
  }

  $grpc.ResponseFuture<$0.SignInResponse> signIn($0.SignInRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$signIn, request, options: options);
  }
}

@$pb.GrpcServiceName('ai.Authenticator')
abstract class AuthenticatorServiceBase extends $grpc.Service {
  $core.String get $name => 'ai.Authenticator';

  AuthenticatorServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SignUpRequest, $0.SignUpResponse>(
        'SignUp',
        signUp_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SignUpRequest.fromBuffer(value),
        ($0.SignUpResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SignInRequest, $0.SignInResponse>(
        'SignIn',
        signIn_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SignInRequest.fromBuffer(value),
        ($0.SignInResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SignUpResponse> signUp_Pre($grpc.ServiceCall call, $async.Future<$0.SignUpRequest> request) async {
    return signUp(call, await request);
  }

  $async.Future<$0.SignInResponse> signIn_Pre($grpc.ServiceCall call, $async.Future<$0.SignInRequest> request) async {
    return signIn(call, await request);
  }

  $async.Future<$0.SignUpResponse> signUp($grpc.ServiceCall call, $0.SignUpRequest request);
  $async.Future<$0.SignInResponse> signIn($grpc.ServiceCall call, $0.SignInRequest request);
}
@$pb.GrpcServiceName('ai.Conversations')
class ConversationsClient extends $grpc.Client {
  static final _$listConversations = $grpc.ClientMethod<$0.ListConversationsRequest, $0.ListConversationsResponse>(
      '/ai.Conversations/ListConversations',
      ($0.ListConversationsRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.ListConversationsResponse.fromBuffer(value));
  static final _$getConversation = $grpc.ClientMethod<$0.GetConversationRequest, $0.GetConversationResponse>(
      '/ai.Conversations/GetConversation',
      ($0.GetConversationRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.GetConversationResponse.fromBuffer(value));
  static final _$deleteConversation = $grpc.ClientMethod<$0.DeleteConversationRequest, $1.Empty>(
      '/ai.Conversations/DeleteConversation',
      ($0.DeleteConversationRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $1.Empty.fromBuffer(value));
  static final _$sendMessage = $grpc.ClientMethod<$0.SendMessageRequest, $0.SendMessageResponse>(
      '/ai.Conversations/SendMessage',
      ($0.SendMessageRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.SendMessageResponse.fromBuffer(value));
  static final _$sendMessageStream = $grpc.ClientMethod<$0.SendMessageStreamRequest, $0.SendMessageStreamResponse>(
      '/ai.Conversations/SendMessageStream',
      ($0.SendMessageStreamRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.SendMessageStreamResponse.fromBuffer(value));

  ConversationsClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options,
        interceptors: interceptors);

  $grpc.ResponseFuture<$0.ListConversationsResponse> listConversations($0.ListConversationsRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$listConversations, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetConversationResponse> getConversation($0.GetConversationRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$getConversation, request, options: options);
  }

  $grpc.ResponseFuture<$1.Empty> deleteConversation($0.DeleteConversationRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$deleteConversation, request, options: options);
  }

  $grpc.ResponseFuture<$0.SendMessageResponse> sendMessage($0.SendMessageRequest request, {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$sendMessage, request, options: options);
  }

  $grpc.ResponseStream<$0.SendMessageStreamResponse> sendMessageStream($0.SendMessageStreamRequest request, {$grpc.CallOptions? options}) {
    return $createStreamingCall(_$sendMessageStream, $async.Stream.fromIterable([request]), options: options);
  }
}

@$pb.GrpcServiceName('ai.Conversations')
abstract class ConversationsServiceBase extends $grpc.Service {
  $core.String get $name => 'ai.Conversations';

  ConversationsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ListConversationsRequest, $0.ListConversationsResponse>(
        'ListConversations',
        listConversations_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListConversationsRequest.fromBuffer(value),
        ($0.ListConversationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConversationRequest, $0.GetConversationResponse>(
        'GetConversation',
        getConversation_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetConversationRequest.fromBuffer(value),
        ($0.GetConversationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteConversationRequest, $1.Empty>(
        'DeleteConversation',
        deleteConversation_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteConversationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendMessageRequest, $0.SendMessageResponse>(
        'SendMessage',
        sendMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SendMessageRequest.fromBuffer(value),
        ($0.SendMessageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendMessageStreamRequest, $0.SendMessageStreamResponse>(
        'SendMessageStream',
        sendMessageStream_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.SendMessageStreamRequest.fromBuffer(value),
        ($0.SendMessageStreamResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListConversationsResponse> listConversations_Pre($grpc.ServiceCall call, $async.Future<$0.ListConversationsRequest> request) async {
    return listConversations(call, await request);
  }

  $async.Future<$0.GetConversationResponse> getConversation_Pre($grpc.ServiceCall call, $async.Future<$0.GetConversationRequest> request) async {
    return getConversation(call, await request);
  }

  $async.Future<$1.Empty> deleteConversation_Pre($grpc.ServiceCall call, $async.Future<$0.DeleteConversationRequest> request) async {
    return deleteConversation(call, await request);
  }

  $async.Future<$0.SendMessageResponse> sendMessage_Pre($grpc.ServiceCall call, $async.Future<$0.SendMessageRequest> request) async {
    return sendMessage(call, await request);
  }

  $async.Stream<$0.SendMessageStreamResponse> sendMessageStream_Pre($grpc.ServiceCall call, $async.Future<$0.SendMessageStreamRequest> request) async* {
    yield* sendMessageStream(call, await request);
  }

  $async.Future<$0.ListConversationsResponse> listConversations($grpc.ServiceCall call, $0.ListConversationsRequest request);
  $async.Future<$0.GetConversationResponse> getConversation($grpc.ServiceCall call, $0.GetConversationRequest request);
  $async.Future<$1.Empty> deleteConversation($grpc.ServiceCall call, $0.DeleteConversationRequest request);
  $async.Future<$0.SendMessageResponse> sendMessage($grpc.ServiceCall call, $0.SendMessageRequest request);
  $async.Stream<$0.SendMessageStreamResponse> sendMessageStream($grpc.ServiceCall call, $0.SendMessageStreamRequest request);
}
