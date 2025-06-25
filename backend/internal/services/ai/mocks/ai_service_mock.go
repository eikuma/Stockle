package mocks

import (
	"context"
	"errors"
	"time"
)

type MockAIService struct {
	ShouldFail        bool
	ShouldTimeout     bool
	ShouldReturnEmpty bool
	CallCount         int
	LastInput         string
}

type SummarizeResponse struct {
	Summary  string `json:"summary"`
	Category string `json:"category"`
}

func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

func (m *MockAIService) SummarizeContent(ctx context.Context, content string) (*SummarizeResponse, error) {
	m.CallCount++
	m.LastInput = content

	if m.ShouldTimeout {
		time.Sleep(10 * time.Second)
		return nil, context.DeadlineExceeded
	}

	if m.ShouldFail {
		return nil, errors.New("AI service temporarily unavailable")
	}

	if m.ShouldReturnEmpty {
		return &SummarizeResponse{
			Summary:  "",
			Category: "",
		}, nil
	}

	return &SummarizeResponse{
		Summary:  "This is a mock summary of the provided content. The article discusses important topics relevant to the user's interests.",
		Category: "Technology",
	}, nil
}

func (m *MockAIService) CategorizContent(ctx context.Context, content string) (string, error) {
	m.CallCount++
	m.LastInput = content

	if m.ShouldTimeout {
		time.Sleep(10 * time.Second)
		return "", context.DeadlineExceeded
	}

	if m.ShouldFail {
		return "", errors.New("AI categorization service temporarily unavailable")
	}

	if m.ShouldReturnEmpty {
		return "", nil
	}

	return "Technology", nil
}

func (m *MockAIService) Reset() {
	m.ShouldFail = false
	m.ShouldTimeout = false
	m.ShouldReturnEmpty = false
	m.CallCount = 0
	m.LastInput = ""
}

type MockGroqService struct {
	*MockAIService
}

func NewMockGroqService() *MockGroqService {
	return &MockGroqService{
		MockAIService: NewMockAIService(),
	}
}

type MockClaudeService struct {
	*MockAIService
}

func NewMockClaudeService() *MockClaudeService {
	return &MockClaudeService{
		MockAIService: NewMockAIService(),
	}
}

type MockRateLimiter struct {
	ShouldLimit bool
	CallCount   int
}

func NewMockRateLimiter() *MockRateLimiter {
	return &MockRateLimiter{}
}

func (m *MockRateLimiter) Allow() bool {
	m.CallCount++
	return !m.ShouldLimit
}

func (m *MockRateLimiter) Reset() {
	m.ShouldLimit = false
	m.CallCount = 0
}