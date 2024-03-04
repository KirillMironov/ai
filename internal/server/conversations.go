package server

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/KirillMironov/ai/internal/api/ai"
	"github.com/KirillMironov/ai/internal/model"
)

type conversationsService interface {
	ListConversations(ctx context.Context, token string, offset, limit int) ([]model.Conversation, error)
	GetConversation(ctx context.Context, token string, id string) (model.Conversation, error)
	SendMessage(ctx context.Context, request model.SendMessageRequest) (model.Message, error)
	SendMessageStream(ctx context.Context, request model.SendMessageRequest, onChunk func(model.Message) error) error
}

type Conversations struct {
	service conversationsService
	api.UnimplementedConversationsServer
}

func NewConversations(service conversationsService) Conversations {
	return Conversations{service: service}
}

func (c Conversations) ListConversations(ctx context.Context, request *api.ListConversationsRequest) (*api.ListConversationsResponse, error) {
	token, err := c.tokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	conversations, err := c.service.ListConversations(ctx, token, int(request.GetOffset()), int(request.GetLimit()))
	if err != nil {
		slog.Error("failed to list conversations", err)
		return nil, err
	}

	return &api.ListConversationsResponse{Conversations: conversationsToAPI(conversations)}, nil
}

func (c Conversations) GetConversation(ctx context.Context, request *api.GetConversationRequest) (*api.GetConversationResponse, error) {
	token, err := c.tokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	conversation, err := c.service.GetConversation(ctx, token, request.GetId())
	if err != nil {
		slog.Error("failed to get conversation", err)
		return nil, err
	}

	response := &api.GetConversationResponse{
		Conversation: conversationToAPI(conversation),
		Messages:     messagesToAPI(conversation.Messages),
	}

	return response, nil
}

func (c Conversations) SendMessage(ctx context.Context, request *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	token, err := c.tokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	req := model.SendMessageRequest{
		Token:          token,
		ConversationID: request.GetConversationId(),
		Role:           model.RoleUser,
		Content:        request.GetContent(),
	}

	message, err := c.service.SendMessage(ctx, req)
	if err != nil {
		slog.Error("failed to send message", err)
		return nil, err
	}

	return &api.SendMessageResponse{Message: conversationMessageToAPI(message)}, nil
}

func (c Conversations) SendMessageStream(request *api.SendMessageStreamRequest, stream api.Conversations_SendMessageStreamServer) error {
	token, err := c.tokenFromHeader(stream.Context())
	if err != nil {
		return err
	}

	req := model.SendMessageRequest{
		Token:          token,
		ConversationID: request.GetConversationId(),
		Role:           model.RoleUser,
		Content:        request.GetContent(),
	}

	onChunk := func(message model.Message) error {
		return stream.Send(&api.SendMessageStreamResponse{Message: conversationMessageToAPI(message)})
	}

	if err = c.service.SendMessageStream(stream.Context(), req, onChunk); err != nil {
		slog.Error("failed to send message stream", err)
		return err
	}

	return nil
}

func (c Conversations) tokenFromHeader(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.InvalidArgument, "failed to get request metadata")
	}

	token, ok := md["jwt"]
	if !ok || len(token) == 0 {
		return "", status.Error(codes.InvalidArgument, "jwt token not found in metadata")
	}

	return token[0], nil
}

func conversationsToAPI(conversations []model.Conversation) []*api.Conversation {
	apiConversations := make([]*api.Conversation, 0, len(conversations))

	for _, conversation := range conversations {
		apiConversations = append(apiConversations, conversationToAPI(conversation))
	}

	return apiConversations
}

func messagesToAPI(messages []model.Message) []*api.Message {
	apiMessages := make([]*api.Message, 0, len(messages))

	for _, message := range messages {
		apiMessages = append(apiMessages, conversationMessageToAPI(message))
	}

	return apiMessages
}

func conversationToAPI(conversation model.Conversation) *api.Conversation {
	return &api.Conversation{
		Id:        conversation.ID,
		Title:     conversation.Title,
		CreatedAt: timestamppb.New(conversation.CreatedAt),
		UpdatedAt: timestamppb.New(conversation.UpdatedAt),
	}
}

func conversationMessageToAPI(message model.Message) *api.Message {
	return &api.Message{
		Id:        message.ID,
		Role:      roleToAPI(message.Role),
		Content:   message.Content,
		CreatedAt: timestamppb.New(message.CreatedAt),
		UpdatedAt: timestamppb.New(message.UpdatedAt),
	}
}

func roleToAPI(role model.Role) api.Role {
	switch role {
	case model.RoleAssistant:
		return api.Role_ROLE_ASSISTANT
	case model.RoleUser:
		return api.Role_ROLE_USER
	default:
		return api.Role_ROLE_UNSPECIFIED
	}
}
