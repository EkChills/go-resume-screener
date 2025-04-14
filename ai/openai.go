package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ekchills/go-resume-screener/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct{}

func (o *OpenAIClient) NewClient() *openai.Client {
	apiKey := os.Getenv("OPENAI_API_KEY") // This gets the actual key from the environment
	fmt.Println("API Key:", apiKey)
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &client
}

func (o *OpenAIClient) AnalyzeResume(resume string) (*models.AnalyzedResume, error) {
	client := o.NewClient()

	completion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.DeveloperMessage(`You are an intelligent resume analysis assistant. Your task is to extract structured information from the resume text provided by the user.

				Return the data strictly in the following JSON format:

				{
				"name": "",
				"email": "",
				"phone": "",
				"skills": ["", ""],
				"education": ["", ""],
				"experience": ["", ""]
				}

				Guidelines:
				- "name": Extract the full name of the resume owner.
				- "email": Extract the most likely email address found.
				- "phone": Extract the phone number, preferably in international format if available.
				- "skills": List 4–8 key technical and soft skills mentioned.
				- "education": Summarize educational qualifications in short form (e.g., "BSc in Computer Science", "MSc in Data Science").
				- "experience": Summarize past work experiences like: ["Company ABC for 3 years", "Company XYZ for 2 years"].

				Important:
				- Do not include any extra explanation, text, or markdown—only return the final JSON object.
				- Be accurate and concise. Do not guess if the information is not clearly provided.
				- if the information is not available, return an empty string or empty array depending on the format specified for that field.`),
				openai.UserMessage(resume),
			},
			Model: openai.ChatModelGPT3_5Turbo,
		},
	)

	if err != nil {
		return nil, err
	}

	var analyzedResume models.AnalyzedResume

	err = json.Unmarshal([]byte(completion.Choices[0].Message.Content), &analyzedResume)
	if err != nil  {
		return nil, err
	}

	return &analyzedResume, nil
}
