package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekchills/go-resume-screener/database"
	"github.com/ekchills/go-resume-screener/services"
	"github.com/gin-gonic/gin"
)

type ResumeController struct{}

type shortlistReq struct {ResumeId string `json:"resumeId"`}


func (r *ResumeController) GetAllResumes(context *gin.Context) {
	resumesService := &services.ResumeService{Db: database.DB, Context: context}
	allResumes, err := resumesService.ListAllResumes()
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to fetch resumes"})
		return
	}
	context.JSON(200, gin.H{"resumes": allResumes})
}

func (r *ResumeController) ShortListResume(context *gin.Context) {

	var b shortlistReq

	err := context.ShouldBindJSON(&b)
	if  err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		fmt.Println("Error in binding JSON:", err)
		return
	}

	convertedId, err := strconv.Atoi(b.ResumeId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume ID"})
		return
	}
	resumesService := &services.ResumeService{Db: database.DB, Context: context}
	err = resumesService.ShortListResume(uint(convertedId))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shortlist resume"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Resume shortlisted successfully"})
}

func (r *ResumeController) GetResumeDetail(c *gin.Context) {
	resumeId := c.Param("id")
	convertedId, err := strconv.Atoi(resumeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume ID"})
		return
	}
	if resumeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume ID is required"})
		return
	}
	resumeService := &services.ResumeService{Db: database.DB, Context: c}
	resume, err := resumeService.GetResumeById(uint(convertedId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resume details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"resume": resume})
}