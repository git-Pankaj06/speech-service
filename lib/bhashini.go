package bhashini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AudioContent struct {
	AudioContent string `json:"audioContent"`
}

type Language struct {
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage,omitempty"`
}

type Config struct {
	Language     Language `json:"language"`
	ServiceId    string   `json:"serviceId"`
	AudioFormat  string   `json:"audioFormat,omitempty"`
	SamplingRate int      `json:"samplingRate,omitempty"`
}

type Asr struct {
	TaskType string `json:"taskType"`
	Config   Config `json:"config"`
}

type InputData struct {
	Audio []AudioContent `json:"audio"`
}

type Payload struct {
	PipelineTasks []Asr     `json:"pipelineTasks"`
	InputData     InputData `json:"inputData"`
}

type RespBody struct {
	PipelineResponse []struct {
		TaskType string
		Config   string
		Output   []struct {
			Source string
			Target string `json:"target,omitempty"`
		} `json:"output"`
	} `json:"pipelineResponse"`
}

func RecognizeAndTranslate(audio string) (string, error) {
	payload, err := getPayload(audio)

	if err != nil {
		return "", err
	}

	auth_key := os.Getenv("BHASHINI_AUTHORIZATION")
	base_url := os.Getenv("BHASHINI_URL")

	req, err := http.NewRequest("POST", base_url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error in making post request")
		return "", err
	}

	req.Header.Set("Authorization", auth_key)

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer res.Body.Close()

	var respBody RespBody

	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		return "", err
	}

	if len(respBody.PipelineResponse) < 1 {
		return "", fmt.Errorf("No response")
	}

	if len(respBody.PipelineResponse[1].Output) == 0 {
		return "", fmt.Errorf("No response")
	}

	return respBody.PipelineResponse[1].Output[0].Target, nil
}

func getPayload(audio string) ([]byte, error) {
	payload := Payload{
		PipelineTasks: []Asr{
			{
				TaskType: "asr",
				Config: Config{
					Language: Language{
						SourceLanguage: "hi",
					},
					ServiceId:    "ai4bharat/conformer-hi-gpu--t4",
					AudioFormat:  "flac",
					SamplingRate: 16000,
				},
			},
			{
				TaskType: "translation",
				Config: Config{
					Language: Language{
						SourceLanguage: "hi",
						TargetLanguage: "en",
					},
					ServiceId: "ai4bharat/indictrans-v2-all-gpu--t4",
				},
			},
		},
		InputData: InputData{
			Audio: []AudioContent{
				{AudioContent: audio},
			},
		},
	}

	return json.Marshal(payload)
}
