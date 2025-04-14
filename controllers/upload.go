package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/ekchills/go-resume-screener/ai"
	"github.com/ekchills/go-resume-screener/database"
	"github.com/ekchills/go-resume-screener/services"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
)

type UploadController struct{}

func (u *UploadController) UploadResume(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	tempDir, err := os.MkdirTemp("", "temp-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary directory.",
		})
		return
	}
	for _, file := range files {
		if file.Header.Get("Content-Type") != "application/pdf" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid file type. Only PDF files are allowed.",
			})
			return
		}
		err := c.SaveUploadedFile(file, tempDir+"/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file.",
			})
			return
		}
		f, r, err := pdf.Open(tempDir + "/" + file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file."})
			return 
		}
		defer f.Close()
	
		var buf bytes.Buffer
		b, err := r.GetPlainText()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file."})
			return 
		}
		buf.ReadFrom(b)
		aiClient := &ai.OpenAIClient{}
		content, err := aiClient.AnalyzeResume(buf.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze resume."+ err.Error()})
			return 
		}
		rs := &services.ResumeService{Db: database.DB, Context: c}
		err = rs.Save(content)
		fmt.Println("Resume Analysis Result: ", content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"+ err.Error()})
			return 
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"msg": "File uploaded successfully."})
	
}