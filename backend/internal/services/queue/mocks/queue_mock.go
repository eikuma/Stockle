package mocks

import (
	"context"
	"errors"
	"sync"
	"time"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Status      JobStatus              `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Retries     int                    `json:"retries"`
	MaxRetries  int                    `json:"max_retries"`
}

type MockQueueService struct {
	mu                sync.RWMutex
	jobs              map[string]*Job
	ShouldFailEnqueue bool
	ShouldFailDequeue bool
	ProcessingDelay   time.Duration
	CallCount         int
}

func NewMockQueueService() *MockQueueService {
	return &MockQueueService{
		jobs: make(map[string]*Job),
	}
}

func (m *MockQueueService) EnqueueJob(ctx context.Context, jobType string, payload map[string]interface{}) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.CallCount++
	
	if m.ShouldFailEnqueue {
		return nil, errors.New("failed to enqueue job")
	}
	
	job := &Job{
		ID:         generateJobID(),
		Type:       jobType,
		Payload:    payload,
		Status:     JobStatusPending,
		CreatedAt:  time.Now(),
		Retries:    0,
		MaxRetries: 3,
	}
	
	m.jobs[job.ID] = job
	return job, nil
}

func (m *MockQueueService) DequeueJob(ctx context.Context) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.CallCount++
	
	if m.ShouldFailDequeue {
		return nil, errors.New("failed to dequeue job")
	}
	
	for _, job := range m.jobs {
		if job.Status == JobStatusPending {
			job.Status = JobStatusProcessing
			return job, nil
		}
	}
	
	return nil, nil
}

func (m *MockQueueService) CompleteJob(ctx context.Context, jobID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	job, exists := m.jobs[jobID]
	if !exists {
		return errors.New("job not found")
	}
	
	job.Status = JobStatusCompleted
	now := time.Now()
	job.ProcessedAt = &now
	
	return nil
}

func (m *MockQueueService) FailJob(ctx context.Context, jobID string, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	job, exists := m.jobs[jobID]
	if !exists {
		return errors.New("job not found")
	}
	
	job.Retries++
	if job.Retries >= job.MaxRetries {
		job.Status = JobStatusFailed
		job.Error = errorMsg
		now := time.Now()
		job.ProcessedAt = &now
	} else {
		job.Status = JobStatusPending
	}
	
	return nil
}

func (m *MockQueueService) GetJobStatus(ctx context.Context, jobID string) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	job, exists := m.jobs[jobID]
	if !exists {
		return nil, errors.New("job not found")
	}
	
	return job, nil
}

func (m *MockQueueService) ProcessJobsAsync(ctx context.Context, processor func(*Job) error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				job, err := m.DequeueJob(ctx)
				if err != nil || job == nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				
				if m.ProcessingDelay > 0 {
					time.Sleep(m.ProcessingDelay)
				}
				
				if err := processor(job); err != nil {
					m.FailJob(ctx, job.ID, err.Error())
				} else {
					m.CompleteJob(ctx, job.ID)
				}
			}
		}
	}()
}

func (m *MockQueueService) GetQueueLength(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	count := 0
	for _, job := range m.jobs {
		if job.Status == JobStatusPending {
			count++
		}
	}
	
	return count, nil
}

func (m *MockQueueService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.jobs = make(map[string]*Job)
	m.ShouldFailEnqueue = false
	m.ShouldFailDequeue = false
	m.ProcessingDelay = 0
	m.CallCount = 0
}

func generateJobID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}