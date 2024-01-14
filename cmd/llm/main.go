package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KirillMironov/ai/llm"
	"github.com/KirillMironov/ai/llm/llama"
)

func main() {
	ctx := context.Background()

	llamaLLM := llama.New("/bin/llama/server", "/models/llama-2-7b-chat.Q4_K_M.gguf")

	if err := llamaLLM.Start(ctx); err != nil {
		log.Fatalf("failed to start llama LLM: %v", err)
	}
	defer llamaLLM.Close(ctx)

	http.HandleFunc("/completion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Prompt string `json:"prompt"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}

		response, err := llamaLLM.Completion(r.Context(), llm.CompletionRequest{Prompt: request.Prompt})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
