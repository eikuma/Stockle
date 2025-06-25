package ai

import (
	"context"
	"testing"
	"time"

	"github.com/private/Stockle/backend/internal/services/ai/mocks"
)

func TestAIService_SummarizeContent_Success(t *testing.T) {
	mockService := mocks.NewMockAIService()
	ctx := context.Background()
	
	content := "This is a test article about artificial intelligence and machine learning technologies."
	
	result, err := mockService.SummarizeContent(ctx, content)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if result.Summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	if result.Category == "" {
		t.Error("Expected non-empty category")
	}
	
	if mockService.CallCount != 1 {
		t.Errorf("Expected 1 call, got: %d", mockService.CallCount)
	}
	
	if mockService.LastInput != content {
		t.Errorf("Expected input '%s', got: '%s'", content, mockService.LastInput)
	}
}

func TestAIService_SummarizeContent_Timeout(t *testing.T) {
	mockService := mocks.NewMockAIService()
	mockService.ShouldTimeout = true
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	content := "Test content"
	
	_, err := mockService.SummarizeContent(ctx, content)
	
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got: %v", err)
	}
}

func TestAIService_SummarizeContent_Failure(t *testing.T) {
	mockService := mocks.NewMockAIService()
	mockService.ShouldFail = true
	ctx := context.Background()
	
	content := "Test content"
	
	_, err := mockService.SummarizeContent(ctx, content)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	expectedErrorMsg := "AI service temporarily unavailable"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMsg, err.Error())
	}
}

func TestAIService_SummarizeContent_EmptyResponse(t *testing.T) {
	mockService := mocks.NewMockAIService()
	mockService.ShouldReturnEmpty = true
	ctx := context.Background()
	
	content := "Test content"
	
	result, err := mockService.SummarizeContent(ctx, content)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if result.Summary != "" {
		t.Errorf("Expected empty summary, got: '%s'", result.Summary)
	}
	
	if result.Category != "" {
		t.Errorf("Expected empty category, got: '%s'", result.Category)
	}
}

func TestAIService_CategorizContent_Success(t *testing.T) {
	mockService := mocks.NewMockAIService()
	ctx := context.Background()
	
	content := "This is a test article about technology trends and innovations."
	
	category, err := mockService.CategorizContent(ctx, content)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if category == "" {
		t.Error("Expected non-empty category")
	}
	
	expectedCategory := "Technology"
	if category != expectedCategory {
		t.Errorf("Expected category '%s', got: '%s'", expectedCategory, category)
	}
}

func TestAIService_CategorizContent_Failure(t *testing.T) {
	mockService := mocks.NewMockAIService()
	mockService.ShouldFail = true
	ctx := context.Background()
	
	content := "Test content"
	
	_, err := mockService.CategorizContent(ctx, content)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	expectedErrorMsg := "AI categorization service temporarily unavailable"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMsg, err.Error())
	}
}

func TestGroqService_FallbackToClaudeOnFailure(t *testing.T) {
	groqService := mocks.NewMockGroqService()
	claudeService := mocks.NewMockClaudeService()
	
	groqService.ShouldFail = true
	
	ctx := context.Background()
	content := "Test content for fallback scenario"
	
	_, err := groqService.SummarizeContent(ctx, content)
	if err == nil {
		t.Error("Expected Groq service to fail")
	}
	
	result, err := claudeService.SummarizeContent(ctx, content)
	if err != nil {
		t.Fatalf("Expected Claude service to succeed as fallback, got error: %v", err)
	}
	
	if result.Summary == "" {
		t.Error("Expected non-empty summary from Claude fallback")
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	limiter := mocks.NewMockRateLimiter()
	
	if !limiter.Allow() {
		t.Error("Expected rate limiter to allow request")
	}
	
	if limiter.CallCount != 1 {
		t.Errorf("Expected 1 call, got: %d", limiter.CallCount)
	}
}

func TestRateLimiter_Deny(t *testing.T) {
	limiter := mocks.NewMockRateLimiter()
	limiter.ShouldLimit = true
	
	if limiter.Allow() {
		t.Error("Expected rate limiter to deny request")
	}
	
	if limiter.CallCount != 1 {
		t.Errorf("Expected 1 call, got: %d", limiter.CallCount)
	}
}

func TestRateLimiter_MultipleRequests(t *testing.T) {
	limiter := mocks.NewMockRateLimiter()
	
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}
	
	if limiter.CallCount != 5 {
		t.Errorf("Expected 5 calls, got: %d", limiter.CallCount)
	}
	
	limiter.ShouldLimit = true
	
	if limiter.Allow() {
		t.Error("Expected request to be denied after rate limit exceeded")
	}
}

func TestServiceReset(t *testing.T) {
	mockService := mocks.NewMockAIService()
	
	mockService.ShouldFail = true
	mockService.ShouldTimeout = true
	mockService.CallCount = 5
	mockService.LastInput = "previous input"
	
	mockService.Reset()
	
	if mockService.ShouldFail {
		t.Error("Expected ShouldFail to be false after reset")
	}
	
	if mockService.ShouldTimeout {
		t.Error("Expected ShouldTimeout to be false after reset")
	}
	
	if mockService.CallCount != 0 {
		t.Errorf("Expected CallCount to be 0 after reset, got: %d", mockService.CallCount)
	}
	
	if mockService.LastInput != "" {
		t.Errorf("Expected LastInput to be empty after reset, got: '%s'", mockService.LastInput)
	}
}