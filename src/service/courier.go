package service

import (
	"math"
	"time"
	"yandex-team.ru/bstask/models"
)

func (s *Service) CreateCouriers(courierList []models.CreateCourierDto) ([]models.CourierDto, error) {
	return s.Repository.SaveCouriers(courierList)
}

func (s *Service) GetCouriers(limit int32, offset int32) ([]models.CourierDto, error) {
	return s.Repository.GetCouriers(limit, offset)
}

func (s *Service) GetCourierById(id int64) (models.CourierDto, error) {
	return s.Repository.GetCourierById(id)
}

func (s *Service) GetCourierMetaInfoById(id int64, startDate string, endDate string) (models.CourierMetaInfoDto, error) {
	count, sum, err := s.Repository.GetCourierSumAndCount(id, startDate, endDate)
	if err != nil {
		return models.CourierMetaInfoDto{}, err
	}

	courier, err := s.Repository.GetCourierById(id)
	if err != nil {
		return models.CourierMetaInfoDto{}, err
	}

	courierMetaInfo := models.CourierMetaInfoDto{
		CourierId:    id,
		CourierType:  courier.CourierType,
		WorkingHours: courier.WorkingHours,
		Regions:      courier.Regions}

	if count == 0 {
		return courierMetaInfo, nil
	}

	startTime, _ := time.Parse(models.Date, startDate)
	endTime, _ := time.Parse(models.Date, endDate)

	rating := int32(math.Round(
		float64(count) / endTime.Sub(startTime).Hours() * float64(models.GetCourierRatingCoefficient(&courier)),
	))
	earnings := sum * models.GetCourierEarningsCoefficient(&courier)

	courierMetaInfo.Rating = &rating
	courierMetaInfo.Earnings = &earnings

	return courierMetaInfo, nil
}
