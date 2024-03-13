package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	api "github.com/KirillMironov/ai/internal/api/llm"
	"github.com/KirillMironov/ai/internal/model"
)

const titleLength = 10

type (
	authenticatorService interface {
		Authenticate(token string) (tokenPayload model.TokenPayload, err error)
	}

	conversationsStorage interface {
		SaveConversation(ctx context.Context, conversation model.Conversation) error
		GetConversationsByUserID(ctx context.Context, userID string, offset, limit int) ([]model.Conversation, error)
		GetConversationByID(ctx context.Context, id string) (conversation model.Conversation, exists bool, err error)
		DeleteConversation(ctx context.Context, id string) error
	}
)

type Conversations struct {
	authenticatorService authenticatorService
	conversationsStorage conversationsStorage
	llmClient            api.LLMClient
}

func NewConversations(
	authenticatorService authenticatorService,
	conversationsStorage conversationsStorage,
	llmClient api.LLMClient,
) Conversations {
	return Conversations{
		authenticatorService: authenticatorService,
		conversationsStorage: conversationsStorage,
		llmClient:            llmClient,
	}
}

func (c Conversations) ListConversations(ctx context.Context, token string, offset, limit int) ([]model.Conversation, error) {
	userID, err := c.authenticate(token)
	if err != nil {
		return nil, err
	}

	conversations, err := c.conversationsStorage.GetConversationsByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("get conversations by user id: %w", err)
	}

	return conversations, nil
}

func (c Conversations) GetConversation(ctx context.Context, token string, id string) (model.Conversation, error) {
	userID, err := c.authenticate(token)
	if err != nil {
		return model.Conversation{}, err
	}

	conversation, exists, err := c.conversationsStorage.GetConversationByID(ctx, id)
	switch {
	case err != nil:
		return model.Conversation{}, fmt.Errorf("get conversation by id: %w", err)
	case !exists:
		return model.Conversation{}, fmt.Errorf("conversation with id '%s' not found", id)
	case conversation.UserID != userID:
		return model.Conversation{}, fmt.Errorf("user id mismatch")
	}

	return conversation, nil
}

func (c Conversations) DeleteConversation(ctx context.Context, token string, id string) error {
	if _, err := c.GetConversation(ctx, token, id); err != nil {
		return err
	}

	if err := c.conversationsStorage.DeleteConversation(ctx, id); err != nil {
		return fmt.Errorf("delete conversation with id '%s': %w", id, err)
	}

	return nil
}

func (c Conversations) SendMessage(ctx context.Context, request model.SendMessageRequest) (model.Message, error) {
	if request.Content == "" {
		return model.Message{}, errors.New("empty content")
	}

	userID, err := c.authenticate(request.Token)
	if err != nil {
		return model.Message{}, err
	}

	conversation, err := c.getOrCreateConversation(ctx, userID, request.ConversationID, request.Content)
	if err != nil {
		return model.Message{}, err
	}

	conversation.Messages = append(conversation.Messages, newMessage(model.RoleUser, request.Content))

	response, err := c.llmClient.ChatCompletion(ctx, &api.ChatCompletionRequest{Messages: messagesToAPI(conversation.Messages)})
	if err != nil {
		return model.Message{}, fmt.Errorf("send message to LLM: %w", err)
	}

	message := newMessage(model.RoleAssistant, response.GetMessage().GetContent())
	conversation.Messages = append(conversation.Messages, message)
	conversation.UpdatedAt = time.Now()

	if err = c.conversationsStorage.SaveConversation(ctx, conversation); err != nil {
		return model.Message{}, fmt.Errorf("save message: %w", err)
	}

	return message, nil
}

func (c Conversations) SendMessageStream(ctx context.Context, request model.SendMessageRequest, onChunk func(model.Message) error) error {
	if request.Content == "" {
		return errors.New("empty content")
	}

	userID, err := c.authenticate(request.Token)
	if err != nil {
		return err
	}

	conversation, err := c.getOrCreateConversation(ctx, userID, request.ConversationID, request.Content)
	if err != nil {
		return err
	}

	conversation.Messages = append(conversation.Messages, newMessage(model.RoleUser, request.Content))

	stream, err := c.llmClient.ChatCompletionStream(ctx, &api.ChatCompletionStreamRequest{Messages: messagesToAPI(conversation.Messages)})
	if err != nil {
		return fmt.Errorf("send message to LLM: %w", err)
	}

	message := newMessage(model.RoleAssistant, "")

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		content := response.GetMessage().GetContent()
		message.Content += content
		chunkMessage := message
		chunkMessage.Content = content
		if err = onChunk(chunkMessage); err != nil {
			return err
		}
	}

	conversation.Messages = append(conversation.Messages, message)
	conversation.UpdatedAt = time.Now()

	if err = c.conversationsStorage.SaveConversation(ctx, conversation); err != nil {
		return fmt.Errorf("save conversation: %w", err)
	}

	return nil
}

func (c Conversations) authenticate(token string) (userID string, err error) {
	tokenPayload, err := c.authenticatorService.Authenticate(token)
	if err != nil {
		return "", fmt.Errorf("authenticate by token: %w", err)
	}

	return tokenPayload.UserID, nil
}

func (c Conversations) getOrCreateConversation(ctx context.Context, userID, conversationID, content string) (conversation model.Conversation, err error) {
	if conversationID == "" {
		contentRunes := []rune(content)
		title := string(contentRunes[:min(titleLength, len(contentRunes))])
		now := time.Now()
		conversation = model.Conversation{
			ID:        uuid.NewString(),
			UserID:    userID,
			Title:     title,
			CreatedAt: now,
			UpdatedAt: now,
		}
	} else {
		var exists bool
		conversation, exists, err = c.conversationsStorage.GetConversationByID(ctx, conversationID)
		switch {
		case err != nil:
			return model.Conversation{}, fmt.Errorf("get conversation by id: %w", err)
		case !exists:
			return model.Conversation{}, fmt.Errorf("conversation with id '%s' not found", conversationID)
		case conversation.UserID != userID:
			return model.Conversation{}, fmt.Errorf("user id mismatch")
		}
	}

	return conversation, nil
}

func newMessage(role model.Role, content string) model.Message {
	now := time.Now()
	return model.Message{
		ID:        uuid.NewString(),
		Role:      role,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func messagesToAPI(messages []model.Message) []*api.Message {
	apiMessages := make([]*api.Message, 0, len(messages))

	for _, message := range messages {
		apiMessages = append(apiMessages, &api.Message{
			Role:    roleToAPI(message.Role),
			Content: message.Content,
		})
	}

	return apiMessages
}

func roleToAPI(role model.Role) api.Role {
	switch role {
	case model.RoleAssistant:
		return api.Role_ROLE_LLM
	case model.RoleUser:
		return api.Role_ROLE_USER
	default:
		return api.Role_ROLE_UNSPECIFIED
	}
}
