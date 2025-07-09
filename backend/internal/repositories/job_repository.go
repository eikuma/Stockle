package repositories

import (
	"github.com/private/Stockle/backend/internal/models"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *models.JobQueue) error
	Update(job *models.JobQueue) error
	GetByID(id string) (*models.JobQueue, error)
	GetNextJob() (*models.JobQueue, error)
	GetPendingJobs() ([]*models.JobQueue, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		db: db,
	}
}

func (r *jobRepository) Create(job *models.JobQueue) error {
	return r.db.Create(job).Error
}

func (r *jobRepository) Update(job *models.JobQueue) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) GetByID(id string) (*models.JobQueue, error) {
	var job models.JobQueue
	err := r.db.Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetNextJob() (*models.JobQueue, error) {
	var job models.JobQueue
	err := r.db.Where("status = ?", models.JobStatusPending).
		Order("priority ASC, created_at ASC").
		First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetPendingJobs() ([]*models.JobQueue, error) {
	var jobs []*models.JobQueue
	err := r.db.Where("status = ?", models.JobStatusPending).
		Order("priority ASC, created_at ASC").
		Find(&jobs).Error
	return jobs, err
}