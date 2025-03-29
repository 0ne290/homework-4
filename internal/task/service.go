package task

import (
	"context"
	"github.com/google/uuid"
	"homework-4/internal"
)

type Service interface {
	Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error)
	GetAll(ctx context.Context) (*GetAllResponse, error)
	GetByUuid(ctx context.Context, taskUuid uuid.UUID) (*GetByUuidResponse, error)
	Update(ctx context.Context, taskUuid uuid.UUID) (*UpdateResponse, error)
	Delete(ctx context.Context, taskUuid uuid.UUID) (*DeleteResponse, error)
}

type RealService struct {
	unitOfWork   UnitOfWork
	timeProvider internal.TimeProvider
	uuidProvider internal.UuidProvider
}

func NewRealService(unitOfWork UnitOfWork, timeProvider internal.TimeProvider, uuidProvider internal.UuidProvider) *RealService {
	return &RealService{unitOfWork, timeProvider, uuidProvider}
}

func (s *RealService) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	timeNow := s.timeProvider.Now()
	task := newTask(s.uuidProvider.Random(), request.Title, request.Description, timeNow, timeNow)

	err = repository.Create(ctx, task)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{}, nil
}

func (s *RealService) GetAll(ctx context.Context) (*GetAllResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := repository.GetAll(ctx)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &GetAllResponse{tasks}, nil
}

func (s *RealService) GetByUuid(ctx context.Context, taskUuid uuid.UUID) (*GetByUuidResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	task, err := repository.GetByUuid(ctx, taskUuid)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &GetByUuidResponse{task}, nil
}

func (s *RealService) Update(ctx context.Context, taskUuid uuid.UUID) (*UpdateResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	task, err := repository.GetByUuid(ctx, taskUuid)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = task.Update(s.timeProvider.Now())
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = repository.Update(ctx, task)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{task.Status}, nil
}

func (s *RealService) Delete(ctx context.Context, taskUuid uuid.UUID) (*DeleteResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	err = repository.Delete(ctx, taskUuid)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{"task deleted"}, nil
}
