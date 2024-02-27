package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"

	api "github.com/KirillMironov/ai/internal/api/llm"
	"github.com/KirillMironov/ai/internal/mock"
	"github.com/KirillMironov/ai/internal/model"
)

const (
	testConversationID = "conversation_id"
	testContent        = "content"
)

var (
	errAuthenticator = errors.New("authenticator error")
	errLLMClient     = errors.New("llm client error")
)

func TestConversations_ListConversations(t *testing.T) {
	type mocks struct {
		authenticatorService authenticatorService
		conversationsStorage conversationsStorage
		llmClient            api.LLMClient
	}

	tests := []struct {
		name              string
		token             string
		wantErr           bool
		wantConversations []model.Conversation
		mocks             mocks
	}{
		{
			name:              "success",
			token:             testToken,
			wantErr:           false,
			wantConversations: []model.Conversation{{ID: testConversationID, UserID: testUserID}},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationsByUserIDFunc: func(_ context.Context, _ string, _ int, _ int) ([]model.Conversation, error) {
						return []model.Conversation{{ID: testConversationID, UserID: testUserID}}, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name:              "authenticator error",
			token:             testToken,
			wantErr:           true,
			wantConversations: nil,
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{}, errAuthenticator
					},
				},
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name:              "storage error",
			token:             testToken,
			wantErr:           true,
			wantConversations: nil,
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationsByUserIDFunc: func(_ context.Context, _ string, _ int, _ int) ([]model.Conversation, error) {
						return nil, errStorage
					},
				},
				llmClient: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conversationsService := NewConversations(tc.mocks.authenticatorService, tc.mocks.conversationsStorage, tc.mocks.llmClient)

			conversations, err := conversationsService.ListConversations(context.Background(), tc.token, 0, 0)
			if (err != nil) != tc.wantErr {
				t.Fatalf("conversationsService.ListConversations() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := conversations, tc.wantConversations; !cmp.Equal(got, want) {
				t.Errorf("conversations = %v, want %v", got, want)
			}
		})
	}
}

func TestConversations_GetConversation(t *testing.T) {
	type mocks struct {
		authenticatorService authenticatorService
		conversationsStorage conversationsStorage
		llmClient            api.LLMClient
	}

	tests := []struct {
		name             string
		token            string
		conversationID   string
		wantErr          bool
		wantConversation model.Conversation
		mocks            mocks
	}{
		{
			name:             "success",
			token:            testToken,
			conversationID:   testConversationID,
			wantErr:          false,
			wantConversation: model.Conversation{ID: testConversationID, UserID: testUserID},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{ID: testConversationID, UserID: testUserID}, true, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name:             "conversation not found",
			token:            testToken,
			conversationID:   testConversationID,
			wantErr:          true,
			wantConversation: model.Conversation{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name:             "user id mismatch",
			token:            testToken,
			conversationID:   testConversationID,
			wantErr:          true,
			wantConversation: model.Conversation{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{ID: testConversationID, UserID: "user-id-mismatch"}, true, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name:             "authenticator error",
			token:            testToken,
			conversationID:   testConversationID,
			wantErr:          true,
			wantConversation: model.Conversation{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{}, errAuthenticator
					},
				},
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name:             "storage error",
			token:            testToken,
			conversationID:   testConversationID,
			wantErr:          true,
			wantConversation: model.Conversation{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, errStorage
					},
				},
				llmClient: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conversationsService := NewConversations(tc.mocks.authenticatorService, tc.mocks.conversationsStorage, tc.mocks.llmClient)

			conversation, err := conversationsService.GetConversation(context.Background(), tc.token, tc.conversationID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("conversationsService.GetConversation() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := conversation, tc.wantConversation; !cmp.Equal(got, want) {
				t.Errorf("conversation = %q, want %q", got, want)
			}
		})
	}
}

func TestConversations_SendMessage(t *testing.T) {
	type mocks struct {
		authenticatorService authenticatorService
		conversationsStorage conversationsStorage
		llmClient            api.LLMClient
	}

	tests := []struct {
		name        string
		request     model.SendMessageRequest
		wantErr     bool
		wantMessage model.Message
		mocks       mocks
	}{
		{ //nolint:dupl
			name: "role user",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     false,
			wantMessage: model.Message{Role: model.RoleUser, Content: testContent},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionFunc: func(_ context.Context, _ *api.ChatCompletionRequest, _ ...grpc.CallOption) (*api.ChatCompletionResponse, error) {
						return &api.ChatCompletionResponse{Message: &api.Message{Role: api.Role_ROLE_USER, Content: testContent}}, nil
					},
				},
			},
		},
		{ //nolint:dupl
			name: "role assistant",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleAssistant,
				Content:        testContent,
			},
			wantErr:     false,
			wantMessage: model.Message{Role: model.RoleAssistant, Content: testContent},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionFunc: func(_ context.Context, _ *api.ChatCompletionRequest, _ ...grpc.CallOption) (*api.ChatCompletionResponse, error) {
						return &api.ChatCompletionResponse{Message: &api.Message{Role: api.Role_ROLE_LLM, Content: testContent}}, nil
					},
				},
			},
		},
		{
			name:        "empty content",
			request:     model.SendMessageRequest{Content: ""},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: nil,
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name: "conversation not found by id",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name: "authenticator error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{}, errAuthenticator
					},
				},
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name: "storage error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, errStorage
					},
				},
				llmClient: nil,
			},
		},
		{
			name: "llm client error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionFunc: func(_ context.Context, _ *api.ChatCompletionRequest, _ ...grpc.CallOption) (*api.ChatCompletionResponse, error) {
						return &api.ChatCompletionResponse{}, errLLMClient
					},
				},
			},
		},
		{
			name: "save conversation error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return errStorage
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionFunc: func(_ context.Context, _ *api.ChatCompletionRequest, _ ...grpc.CallOption) (*api.ChatCompletionResponse, error) {
						return &api.ChatCompletionResponse{Message: &api.Message{Role: api.Role_ROLE_LLM, Content: testContent}}, nil
					},
				},
			},
		},
		{
			name: "user id mismatch",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantMessage: model.Message{},
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{ID: testConversationID, UserID: "mismatch-user-id"}, true, nil
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionFunc: func(_ context.Context, _ *api.ChatCompletionRequest, _ ...grpc.CallOption) (*api.ChatCompletionResponse, error) {
						return &api.ChatCompletionResponse{Message: &api.Message{Role: api.Role_ROLE_LLM, Content: testContent}}, nil
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conversationsService := NewConversations(tc.mocks.authenticatorService, tc.mocks.conversationsStorage, tc.mocks.llmClient)

			message, err := conversationsService.SendMessage(context.Background(), tc.request)
			if (err != nil) != tc.wantErr {
				t.Fatalf("conversationsService.SendMessage() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := message.Role, tc.wantMessage.Role; got != want {
				t.Errorf("message.Role = %q, want %q", got, want)
			}
			if got, want := message.Content, tc.wantMessage.Content; got != want {
				t.Errorf("message.Content = %q, want %q", got, want)
			}
		})
	}
}

func TestConversations_SendMessageStream(t *testing.T) {
	type mocks struct {
		authenticatorService authenticatorService
		conversationsStorage conversationsStorage
		llmClient            api.LLMClient
	}

	tests := []struct {
		name        string
		request     model.SendMessageRequest
		wantErr     bool
		wantContent string
		mocks       mocks
	}{
		{
			name: "role user",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     false,
			wantContent: testContent,
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: llmClientChatCompletionStream(api.Role_ROLE_USER, testContent),
			},
		},
		{
			name: "role assistant",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleAssistant,
				Content:        testContent,
			},
			wantErr:     false,
			wantContent: testContent,
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: llmClientChatCompletionStream(api.Role_ROLE_LLM, testContent),
			},
		},
		{
			name:        "empty content",
			request:     model.SendMessageRequest{Content: ""},
			wantErr:     true,
			wantContent: "",
			mocks: mocks{
				authenticatorService: nil,
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name: "conversation not found by id",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantContent: "",
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, nil
					},
				},
				llmClient: nil,
			},
		},
		{
			name: "authenticator error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantContent: "",
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{}, errAuthenticator
					},
				},
				conversationsStorage: nil,
				llmClient:            nil,
			},
		},
		{
			name: "storage error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: testConversationID,
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantContent: "",
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					GetConversationByIDFunc: func(_ context.Context, _ string) (model.Conversation, bool, error) {
						return model.Conversation{}, false, errStorage
					},
				},
				llmClient: nil,
			},
		},
		{
			name: "llm client error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantContent: "",
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return nil
					},
				},
				llmClient: &mock.LLMClient{
					ChatCompletionStreamFunc: func(_ context.Context, _ *api.ChatCompletionStreamRequest, _ ...grpc.CallOption) (api.LLM_ChatCompletionStreamClient, error) {
						return nil, errLLMClient
					},
				},
			},
		},
		{
			name: "save conversation error",
			request: model.SendMessageRequest{
				Token:          testToken,
				ConversationID: "",
				Role:           model.RoleUser,
				Content:        testContent,
			},
			wantErr:     true,
			wantContent: testContent,
			mocks: mocks{
				authenticatorService: &mock.AuthenticatorService{
					AuthenticateFunc: func(_ string) (model.TokenPayload, error) {
						return model.TokenPayload{UserID: testUserID, Username: testUsername}, nil
					},
				},
				conversationsStorage: &mock.ConversationsStorage{
					SaveConversationFunc: func(_ context.Context, _ model.Conversation) error {
						return errStorage
					},
				},
				llmClient: llmClientChatCompletionStream(api.Role_ROLE_LLM, testContent),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conversationsService := NewConversations(tc.mocks.authenticatorService, tc.mocks.conversationsStorage, tc.mocks.llmClient)

			var content string
			onChunk := func(message model.Message) error {
				content += message.Content
				return nil
			}

			err := conversationsService.SendMessageStream(context.Background(), tc.request, onChunk)
			if (err != nil) != tc.wantErr {
				t.Fatalf("conversationsService.SendMessageStream() = %v, wantErr %v", err, tc.wantErr)
			}
			if got, want := content, tc.wantContent; got != want {
				t.Errorf("content = %q, want %q", got, want)
			}
		})
	}
}

func llmClientChatCompletionStream(role api.Role, content string) *mock.LLMClient {
	return &mock.LLMClient{
		ChatCompletionStreamFunc: func(_ context.Context, _ *api.ChatCompletionStreamRequest, _ ...grpc.CallOption) (api.LLM_ChatCompletionStreamClient, error) {
			var idx int
			return &mock.LLMChatCompletionStreamClient{
				RecvFunc: func() (*api.ChatCompletionStreamResponse, error) {
					if idx >= len(content) {
						return nil, io.EOF
					}
					resp := &api.ChatCompletionStreamResponse{Message: &api.Message{Role: role, Content: string(content[idx])}}
					idx++
					return resp, nil
				},
			}, nil
		},
	}
}
