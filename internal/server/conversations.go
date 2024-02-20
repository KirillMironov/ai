package server

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/ai/internal/api/ai"
	"github.com/KirillMironov/ai/internal/model"
)

type conversationsService interface {
	ListConversations(ctx context.Context, token string, offset, limit int) ([]model.Conversation, error)
	GetConversation(ctx context.Context, token string, id string) (model.Conversation, []model.Message, error)
	SendMessage(ctx context.Context, token string, request model.SendMessageRequest) (model.Message, error)
	SendMessageStream(ctx context.Context, token string, request model.SendMessageRequest, onChunk func(model.Message) error) error
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

	conversations, err := c.service.ListConversations(ctx, token, int(request.Offset), int(request.Limit))
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

	conversation, messages, err := c.service.GetConversation(ctx, token, request.Id)
	if err != nil {
		slog.Error("failed to get conversation", err)
		return nil, err
	}

	response := &api.GetConversationResponse{
		Conversation: conversationToAPI(conversation),
		Messages:     messagesToAPI(messages),
	}

	return response, nil
}

func (c Conversations) SendMessage(ctx context.Context, request *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	token, err := c.tokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	role, err := roleFromAPI(request.Role)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	req := model.SendMessageRequest{
		ConversationID: request.ConversationId,
		Role:           role,
		Content:        request.Content,
	}

	message, err := c.service.SendMessage(ctx, token, req)
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

	role, err := roleFromAPI(request.Role)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	req := model.SendMessageRequest{
		ConversationID: request.ConversationId,
		Role:           role,
		Content:        request.Content,
	}

	onChunk := func(message model.Message) error {
		return stream.Send(&api.SendMessageStreamResponse{Message: conversationMessageToAPI(message)})
	}

	if err = c.service.SendMessageStream(stream.Context(), token, req, onChunk); err != nil {
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
		CreatedAt: conversation.CreatedAt.Unix(),
		UpdatedAt: conversation.UpdatedAt.Unix(),
	}
}

func conversationMessageToAPI(message model.Message) *api.Message {
	return &api.Message{
		Id:      message.ID,
		Role:    message.Role.String(),
		Content: message.Content,
	}
}

func roleFromAPI(apiRole string) (model.Role, error) {
	var role model.Role
	switch apiRole {
	case "assistant":
		role = model.RoleAssistant
	case "user":
		role = model.RoleUser
	default:
		return role, fmt.Errorf("unexpected role: '%s'", apiRole)
	}
	return role, nil
}
