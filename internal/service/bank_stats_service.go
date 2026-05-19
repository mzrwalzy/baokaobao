package service

import (
	"baokaobao/internal/repository"
)

type BankStatsService struct {
	repo *repository.Repository
}

func NewBankStatsService(repo *repository.Repository) *BankStatsService {
	return &BankStatsService{repo: repo}
}

func (s *BankStatsService) GetBankStatsList(page, pageSize int, sortBy, sortOrder string) ([]repository.BankStat, int64, error) {
	return s.repo.GetBankStatsList(page, pageSize, sortBy, sortOrder)
}

func (s *BankStatsService) GetBankStats(bankID int64) (*repository.BankDetailStat, error) {
	return s.repo.GetBankStats(bankID)
}
