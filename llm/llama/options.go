package llama

type Option func(llama *Llama)

func WithServerHost(host string) Option {
	return func(llama *Llama) {
		llama.serverHost = host
	}
}

func WithServerPort(port int) Option {
	return func(llama *Llama) {
		llama.serverPort = port
	}
}

func WithContextSize(size int) Option {
	return func(llama *Llama) {
		llama.contextSize = size
	}
}

func WithNumSlots(num int) Option {
	return func(llama *Llama) {
		llama.numSlots = num
	}
}

func WithNumThreads(num int) Option {
	return func(llama *Llama) {
		llama.numThreads = num
	}
}

func WithMmap(mmap bool) Option {
	return func(llama *Llama) {
		llama.mmap = mmap
	}
}

func WithCacheChatPrompt(cacheChatPrompt bool) Option {
	return func(llama *Llama) {
		llama.cacheChatPrompt = cacheChatPrompt
	}
}

func WithSystemPrompt(systemPrompt string) Option {
	return func(llama *Llama) {
		llama.systemPrompt = systemPrompt
	}
}

func WithChatTemplate(chatTemplate string) Option {
	return func(llama *Llama) {
		llama.chatTemplate = chatTemplate
	}
}
