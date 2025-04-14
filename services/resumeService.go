package services

import (
	"fmt"
	"strings"

	"github.com/ekchills/go-resume-screener/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResumeService struct {
	Db *gorm.DB
	Context *gin.Context
}


func (r *ResumeService) Save(resumeData *models.AnalyzedResume) error {
	userId, ok := r.Context.Get("userId")
	fmt.Println("User ID from context:", userId)
	if(!ok) {
		return nil
	}
	if err := r.Db.Create(&models.Resume{
		Email: resumeData.Email,
		Name:  resumeData.Name,
		Phone: resumeData.Phone,
		Skills: strings.Join(resumeData.Skills, ","),
		Education: strings.Join(resumeData.Education, ","),
		Experience: strings.Join(resumeData.Experience, ","),
		UserID: userId.(uint),
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ResumeService) ListAllResumes() (*[]models.Resume, error) {
	userId, ok := r.Context.Get("userId")
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}
	var resumes []models.Resume
	if err := r.Db.Where("user_id = ?", userId).Find(&resumes).Error; err != nil {
		return nil, err
	}
	return &resumes, nil
}

func (rr *ResumeService) ShortListResume(resumeId uint) error {
	var resume models.Resume
	if err := rr.Db.First(&resume, resumeId).Error; err != nil {
		return err
	}
	if err := rr.Db.Model(&resume).Update("Shortlisted", true).Error; err != nil {
		return err	
	}
	return nil
}

func (r *ResumeService) GetResumeById(resumeId uint) (*models.Resume, error) {
	userId , ok := r.Context.Get("userId")
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}
	var resume models.Resume
	if err := r.Db.Where("id = ? AND user_id = ?", resumeId, userId).First(&resume).Error; err != nil {
		return nil, err
	}
	
	return &resume, nil
}