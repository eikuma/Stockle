package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eikuma/stockle/backend/internal/models"
	"github.com/eikuma/stockle/backend/internal/repositories"
	"github.com/google/uuid"
)

type JobService struct {
	jobRepo     repositories.JobRepository
	articleRepo repositories.ArticleRepository
	aiService   *AIService
}

type JobPayload struct {
	ArticleID string                 `json:"article_id"`
	JobType   string                 `json:"job_type"`
	Options   map[string]interface{} `json:"options"`
}

func NewJobService(jobRepo repositories.JobRepository, articleRepo repositories.ArticleRepository, aiService *AIService) *JobService {
	return &JobService{
		jobRepo:     jobRepo,
		articleRepo: articleRepo,
		aiService:   aiService,
	}
}

func (s *JobService) EnqueueSummaryJob(articleID string, priority int) error {
	payload := JobPayload{
		ArticleID: articleID,
		JobType:   "summarize",
		Options: map[string]interface{}{
			"summary_type": "medium",
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %w", err)
	}

	job := &models.JobQueue{
		ID:         uuid.New().String(),
		JobType:    "summarize",
		Priority:   priority,
		Status:     "pending",
		Payload:    string(payloadJSON),
		MaxRetries: 3,
	}

	return s.jobRepo.Create(job)
}

func (s *JobService) ProcessJob(ctx context.Context, job *models.JobQueue) error {
	var payload JobPayload
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal job payload: %w", err)
	}

	switch job.JobType {
	case "summarize":
		return s.processSummaryJob(ctx, job, &payload)
	default:
		return fmt.Errorf("unknown job type: %s", job.JobType)
	}
}

func (s *JobService) processSummaryJob(ctx context.Context, job *models.JobQueue, payload *JobPayload) error {
	// 記事の取得
	article, err := s.articleRepo.GetByID(payload.ArticleID)
	if err != nil {
		return fmt.Errorf("failed to get article: %w", err)
	}

	// 既に要約が存在する場合はスキップ
	if article.Summary != nil && *article.Summary != "" {
		return nil
	}

	// 要約生成リクエストの作成
	summaryType := "medium"
	if st, ok := payload.Options["summary_type"].(string); ok {
		summaryType = st
	}

	req := &SummaryRequest{
		Content:     *article.Content,
		Title:       article.Title,
		URL:         article.URL,
		Language:    article.Language,
		SummaryType: summaryType,
	}

	// 要約生成
	summary, err := s.aiService.GenerateSummary(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// 記事の更新
	article.Summary = &summary.Summary
	article.SummaryGenerationStatus = "completed"
	article.SummaryGeneratedAt = &summary.GeneratedAt
	article.SummaryModelVersion = &summary.ModelVersion

	if err := s.articleRepo.Update(article); err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	log.Printf("Summary generated for article %s using %s", article.ID, summary.Provider)
	return nil
}

func (s *JobService) StartWorker(ctx context.Context, workerID int) {
	log.Printf("Starting worker %d", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", workerID)
			return
		default:
			job, err := s.jobRepo.GetNextJob()
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if job == nil {
				time.Sleep(5 * time.Second)
				continue
			}

			s.processJobWithRetry(ctx, job)
		}
	}
}

func (s *JobService) processJobWithRetry(ctx context.Context, job *models.JobQueue) {
	job.Status = "processing"
	job.StartedAt = timePtr(time.Now())
	s.jobRepo.Update(job)

	err := s.ProcessJob(ctx, job)

	if err != nil {
		log.Printf("Job %s failed: %v", job.ID, err)

		job.RetryCount++
		job.ErrorMessage = stringPtr(err.Error())

		if job.RetryCount >= job.MaxRetries {
			job.Status = "failed"
		} else {
			job.Status = "pending"
		}
	} else {
		job.Status = "completed"
		job.CompletedAt = timePtr(time.Now())
	}

	s.jobRepo.Update(job)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func stringPtr(s string) *string {
	return &s
}