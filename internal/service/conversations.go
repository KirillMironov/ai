package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/KirillMironov/ai/internal/api/llm"
	"github.com/KirillMironov/ai/internal/model"
)

type (
	authenticatorService interface {
		Authenticate(token string) (tokenPayload model.TokenPayload, err error)
	}

	conversationsStorage interface {
		SaveConversation(ctx context.Context, conversation model.Conversation) error
		GetConversationsByUserID(ctx context.Context, userID string, offset, limit int) ([]model.Conversation, error)
		GetConversationByID(ctx context.Context, id string) (conversation model.Conversation, exists bool, err error)
	}

	messagesStorage interface {
		SaveMessage(ctx context.Context, conversationID string, message model.Message) error
		GetMessagesByConversationID(ctx context.Context, conversationID string) ([]model.Message, error)
	}
)

type Conversations struct {
	authenticatorService authenticatorService
	conversationsStorage conversationsStorage
	messagesStorage      messagesStorage
	llmClient            api.LLMClient
}

func NewConversations(
	authenticatorService authenticatorService,
	conversationsStorage conversationsStorage,
	messagesStorage messagesStorage,
	llmClient api.LLMClient,
) Conversations {
	return Conversations{
		authenticatorService: authenticatorService,
		conversationsStorage: conversationsStorage,
		messagesStorage:      messagesStorage,
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

func (c Conversations) GetConversation(ctx context.Context, token string, id string) (model.Conversation, []model.Message, error) {
	userID, err := c.authenticate(token)
	if err != nil {
		return model.Conversation{}, nil, err
	}

	conversation, exists, err := c.conversationsStorage.GetConversationByID(ctx, id)
	if err != nil {
		return model.Conversation{}, nil, fmt.Errorf("get conversation by id: %w", err)
	}
	if !exists {
		return model.Conversation{}, nil, fmt.Errorf("conversation with id '%s' not found", id)
	}
	if conversation.UserID != userID {
		return model.Conversation{}, nil, fmt.Errorf("user id mismatch")
	}

	messages, err := c.messagesStorage.GetMessagesByConversationID(ctx, conversation.ID)
	if err != nil {
		return model.Conversation{}, nil, fmt.Errorf("get messages by conversation id: %w", err)
	}

	return conversation, messages, nil
}

func (c Conversations) SendMessage(ctx context.Context, token string, request model.SendMessageRequest) (model.Message, error) {
	userID, err := c.authenticate(token)
	if err != nil {
		return model.Message{}, err
	}

	if request.Content == "" {
		return model.Message{}, errors.New("empty content")
	}

	var conversation model.Conversation

	if request.ConversationID == "" {
		contentRunes := []rune(request.Content)
		title := string(contentRunes[:min(10, len(contentRunes))])
		now := time.Now()

		conversation = model.Conversation{
			ID:        uuid.NewString(),
			UserID:    userID,
			Title:     title,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err = c.conversationsStorage.SaveConversation(ctx, conversation); err != nil {
			return model.Message{}, fmt.Errorf("save conversation: %w", err)
		}
	} else {
		var exists bool
		conversation, exists, err = c.conversationsStorage.GetConversationByID(ctx, request.ConversationID)
		if err != nil {
			return model.Message{}, fmt.Errorf("get conversation by id: %w", err)
		}
		if !exists {
			return model.Message{}, fmt.Errorf("conversation with id '%s' not found", request.ConversationID)
		}
		if conversation.UserID != userID {
			return model.Message{}, fmt.Errorf("user id mismatch")
		}
	}

	messages, err := c.messagesStorage.GetMessagesByConversationID(ctx, conversation.ID)
	if err != nil {
		return model.Message{}, fmt.Errorf("get messages by conversation id: %w", err)
	}

	messages = append(messages, model.Message{
		Role:    request.Role,
		Content: request.Content,
	})

	response, err := c.llmClient.ChatCompletion(ctx, &api.ChatCompletionRequest{Messages: messagesToAPI(messages)})
	if err != nil {
		return model.Message{}, fmt.Errorf("send message to LLM: %w", err)
	}

	message := messageFromAPI(response.Message)

	message.ID = uuid.NewString()
	if err = c.messagesStorage.SaveMessage(ctx, conversation.ID, message); err != nil {
		return model.Message{}, fmt.Errorf("save message: %w", err)
	}

	conversation.UpdatedAt = time.Now()
	if err = c.conversationsStorage.SaveConversation(ctx, conversation); err != nil {
		return model.Message{}, fmt.Errorf("save conversation: %w", err)
	}

	return message, nil
}

func (c Conversations) SendMessageStream(ctx context.Context, token string, request model.SendMessageRequest, onChunk func(model.Message) error) error {
	userID, err := c.authenticate(token)
	if err != nil {
		return err
	}

	if request.Content == "" {
		return errors.New("empty content")
	}

	var conversation model.Conversation

	if request.ConversationID == "" {
		contentRunes := []rune(request.Content)
		title := string(contentRunes[:min(10, len(contentRunes))])
		now := time.Now()

		conversation = model.Conversation{
			ID:        uuid.NewString(),
			UserID:    userID,
			Title:     title,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err = c.conversationsStorage.SaveConversation(ctx, conversation); err != nil {
			return fmt.Errorf("save conversation: %w", err)
		}
	} else {
		var exists bool
		conversation, exists, err = c.conversationsStorage.GetConversationByID(ctx, request.ConversationID)
		if err != nil {
			return fmt.Errorf("get conversation by id: %w", err)
		}
		if !exists {
			return fmt.Errorf("conversation with id '%s' not found", request.ConversationID)
		}
		if conversation.UserID != userID {
			return fmt.Errorf("user id mismatch")
		}
	}

	messages, err := c.messagesStorage.GetMessagesByConversationID(ctx, conversation.ID)
	if err != nil {
		return fmt.Errorf("get messages by conversation id: %w", err)
	}

	messages = append(messages, model.Message{
		Role:    request.Role,
		Content: request.Content,
	})

	stream, err := c.llmClient.ChatCompletionStream(ctx, &api.ChatCompletionStreamRequest{Messages: messagesToAPI(messages)})
	if err != nil {
		return fmt.Errorf("send message to LLM: %w", err)
	}

	var message model.Message

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		msg := messageFromAPI(response.Message)
		if message.Role == 0 {
			message.Role = msg.Role
		}
		message.Content += msg.Content
		if err = onChunk(messageFromAPI(response.Message)); err != nil {
			return err
		}
	}

	message.ID = uuid.NewString()
	if err = c.messagesStorage.SaveMessage(ctx, conversation.ID, message); err != nil {
		return fmt.Errorf("save message: %w", err)
	}

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

func messageFromAPI(apiMessage *api.Message) model.Message {
	return model.Message{
		Role:    roleFromAPI(apiMessage.Role),
		Content: apiMessage.Content,
	}
}

func roleFromAPI(apiRole api.Role) model.Role {
	switch apiRole {
	case api.Role_ROLE_LLM:
		return model.RoleAssistant
	case api.Role_ROLE_USER:
		return model.RoleUser
	default:
		return 0
	}
}
