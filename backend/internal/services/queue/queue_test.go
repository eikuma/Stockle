package queue

import (
	"context"
	"testing"
	"time"

	"github.com/private/Stockle/backend/internal/services/queue/mocks"
)

func TestQueueService_EnqueueJob_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{
		"article_id": "123",
		"content":    "Test article content",
		"user_id":    "user_123",
	}
	
	job, err := queueService.EnqueueJob(ctx, "summarize_article", payload)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if job.ID == "" {
		t.Error("Expected non-empty job ID")
	}
	
	if job.Type != "summarize_article" {
		t.Errorf("Expected job type 'summarize_article', got: %s", job.Type)
	}
	
	if job.Status != mocks.JobStatusPending {
		t.Errorf("Expected job status '%s', got: %s", mocks.JobStatusPending, job.Status)
	}
	
	if job.Payload["article_id"] != "123" {
		t.Errorf("Expected article_id '123', got: %v", job.Payload["article_id"])
	}
}

func TestQueueService_EnqueueJob_Failure(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	queueService.ShouldFailEnqueue = true
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	
	_, err := queueService.EnqueueJob(ctx, "test_job", payload)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	expectedErrorMsg := "failed to enqueue job"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMsg, err.Error())
	}
}

func TestQueueService_DequeueJob_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	enqueuedJob, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	dequeuedJob, err := queueService.DequeueJob(ctx)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if dequeuedJob == nil {
		t.Fatal("Expected job, got nil")
	}
	
	if dequeuedJob.ID != enqueuedJob.ID {
		t.Errorf("Expected job ID '%s', got: '%s'", enqueuedJob.ID, dequeuedJob.ID)
	}
	
	if dequeuedJob.Status != mocks.JobStatusProcessing {
		t.Errorf("Expected job status '%s', got: %s", mocks.JobStatusProcessing, dequeuedJob.Status)
	}
}

func TestQueueService_DequeueJob_EmptyQueue(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	job, err := queueService.DequeueJob(ctx)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if job != nil {
		t.Error("Expected nil job from empty queue")
	}
}

func TestQueueService_DequeueJob_Failure(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	queueService.ShouldFailDequeue = true
	ctx := context.Background()
	
	_, err := queueService.DequeueJob(ctx)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	expectedErrorMsg := "failed to dequeue job"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMsg, err.Error())
	}
}

func TestQueueService_CompleteJob_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	job, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	err := queueService.CompleteJob(ctx, job.ID)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	completedJob, _ := queueService.GetJobStatus(ctx, job.ID)
	
	if completedJob.Status != mocks.JobStatusCompleted {
		t.Errorf("Expected job status '%s', got: %s", mocks.JobStatusCompleted, completedJob.Status)
	}
	
	if completedJob.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set")
	}
}

func TestQueueService_FailJob_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	job, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	errorMsg := "Processing failed due to network error"
	err := queueService.FailJob(ctx, job.ID, errorMsg)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	failedJob, _ := queueService.GetJobStatus(ctx, job.ID)
	
	if failedJob.Retries != 1 {
		t.Errorf("Expected retries to be 1, got: %d", failedJob.Retries)
	}
	
	if failedJob.Status != mocks.JobStatusPending {
		t.Errorf("Expected job status '%s' after first failure, got: %s", mocks.JobStatusPending, failedJob.Status)
	}
}

func TestQueueService_FailJob_MaxRetriesReached(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	job, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	errorMsg := "Persistent failure"
	
	for i := 0; i < 3; i++ {
		queueService.FailJob(ctx, job.ID, errorMsg)
	}
	
	finalJob, _ := queueService.GetJobStatus(ctx, job.ID)
	
	if finalJob.Status != mocks.JobStatusFailed {
		t.Errorf("Expected job status '%s', got: %s", mocks.JobStatusFailed, finalJob.Status)
	}
	
	if finalJob.Error != errorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", errorMsg, finalJob.Error)
	}
	
	if finalJob.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set for failed job")
	}
}

func TestQueueService_GetJobStatus_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	originalJob, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	retrievedJob, err := queueService.GetJobStatus(ctx, originalJob.ID)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if retrievedJob.ID != originalJob.ID {
		t.Errorf("Expected job ID '%s', got: '%s'", originalJob.ID, retrievedJob.ID)
	}
	
	if retrievedJob.Type != originalJob.Type {
		t.Errorf("Expected job type '%s', got: '%s'", originalJob.Type, retrievedJob.Type)
	}
}

func TestQueueService_GetJobStatus_NotFound(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	_, err := queueService.GetJobStatus(ctx, "nonexistent-job-id")
	
	if err == nil {
		t.Error("Expected error for nonexistent job, got nil")
	}
	
	expectedErrorMsg := "job not found"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got: '%s'", expectedErrorMsg, err.Error())
	}
}

func TestQueueService_ProcessJobsAsync_Success(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	payload := map[string]interface{}{"test": "data"}
	job, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	processed := false
	processor := func(j *mocks.Job) error {
		if j.ID == job.ID {
			processed = true
		}
		return nil
	}
	
	queueService.ProcessJobsAsync(ctx, processor)
	
	time.Sleep(500 * time.Millisecond)
	
	if !processed {
		t.Error("Expected job to be processed")
	}
	
	processedJob, _ := queueService.GetJobStatus(ctx, job.ID)
	if processedJob.Status != mocks.JobStatusCompleted {
		t.Errorf("Expected job status '%s', got: %s", mocks.JobStatusCompleted, processedJob.Status)
	}
}

func TestQueueService_ProcessJobsAsync_ProcessorError(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	payload := map[string]interface{}{"test": "data"}
	job, _ := queueService.EnqueueJob(ctx, "test_job", payload)
	
	processor := func(j *mocks.Job) error {
		return errors.New("processing error")
	}
	
	queueService.ProcessJobsAsync(ctx, processor)
	
	time.Sleep(500 * time.Millisecond)
	
	processedJob, _ := queueService.GetJobStatus(ctx, job.ID)
	if processedJob.Retries != 1 {
		t.Errorf("Expected job retries to be 1, got: %d", processedJob.Retries)
	}
}

func TestQueueService_GetQueueLength(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	initialLength, _ := queueService.GetQueueLength(ctx)
	if initialLength != 0 {
		t.Errorf("Expected initial queue length 0, got: %d", initialLength)
	}
	
	payload := map[string]interface{}{"test": "data"}
	queueService.EnqueueJob(ctx, "test_job_1", payload)
	queueService.EnqueueJob(ctx, "test_job_2", payload)
	
	length, err := queueService.GetQueueLength(ctx)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if length != 2 {
		t.Errorf("Expected queue length 2, got: %d", length)
	}
	
	queueService.DequeueJob(ctx)
	
	lengthAfterDequeue, _ := queueService.GetQueueLength(ctx)
	if lengthAfterDequeue != 1 {
		t.Errorf("Expected queue length 1 after dequeue, got: %d", lengthAfterDequeue)
	}
}

func TestQueueService_ConcurrentOperations(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	numJobs := 10
	payload := map[string]interface{}{"concurrent": "test"}
	
	for i := 0; i < numJobs; i++ {
		go func(i int) {
			queueService.EnqueueJob(ctx, "concurrent_job", payload)
		}(i)
	}
	
	time.Sleep(100 * time.Millisecond)
	
	length, _ := queueService.GetQueueLength(ctx)
	if length != numJobs {
		t.Errorf("Expected queue length %d, got: %d", numJobs, length)
	}
}

func TestQueueService_Reset(t *testing.T) {
	queueService := mocks.NewMockQueueService()
	ctx := context.Background()
	
	payload := map[string]interface{}{"test": "data"}
	queueService.EnqueueJob(ctx, "test_job", payload)
	queueService.ShouldFailEnqueue = true
	
	initialLength, _ := queueService.GetQueueLength(ctx)
	if initialLength == 0 {
		t.Error("Expected queue to have jobs before reset")
	}
	
	queueService.Reset()
	
	lengthAfterReset, _ := queueService.GetQueueLength(ctx)
	if lengthAfterReset != 0 {
		t.Errorf("Expected queue length 0 after reset, got: %d", lengthAfterReset)
	}
	
	if queueService.ShouldFailEnqueue {
		t.Error("Expected ShouldFailEnqueue to be false after reset")
	}
	
	if queueService.CallCount != 0 {
		t.Errorf("Expected CallCount to be 0 after reset, got: %d", queueService.CallCount)
	}
}