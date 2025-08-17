package service

import "github.com/fathimasithara01/tradeverse/repository"

type CopyService struct{ Repo *repository.CopyRepository }

func NewCopyService(repo *repository.CopyRepository) *CopyService { return &CopyService{Repo: repo} }

func (s *CopyService) GetCopyStatus(followerID, masterID uint) (bool, error) {
	session, err := s.Repo.GetSessionStatus(followerID, masterID)
	if err != nil {
		return false, nil
	}
	return session.IsActive, nil
}

func (s *CopyService) StartCopying(followerID, masterID uint) error {
	return s.Repo.StartOrUpdateSession(followerID, masterID)
}

func (s *CopyService) StopCopying(followerID, masterID uint) error {
	return s.Repo.StopSession(followerID, masterID)
}
