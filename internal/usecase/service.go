package usecase

import (
	"github.com/v.kirpichov/db_tp/internal/utils/errors"
	"github.com/v.kirpichov/db_tp/internal/models"
	"github.com/v.kirpichov/db_tp/internal/repository"
)

type ServiceUsecase struct {
	serviceRepository repository.ServiceR
}

func NewServiceUsecase(serviceRepository repository.ServiceR) *ServiceUsecase {
	return &ServiceUsecase{serviceRepository: serviceRepository}
}

func (usecase *ServiceUsecase) Clear() (err error) {
	err = usecase.serviceRepository.Clear()
	if err != nil {
		err = errors.ServerInternal
		return
	}
	return
}

func (usecase *ServiceUsecase) Status() (status *models.ForumStatus, err error) {
	status, err = usecase.serviceRepository.Status()
	if err != nil {
		err = errors.ServerInternal
		return
	}
	return
}
