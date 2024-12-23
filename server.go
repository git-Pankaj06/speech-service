package main

import (
	"encoding/json"
	"log"
	"net/http"

	bhashini "github.com/fundu-games/speech-service/lib"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type AudioBody struct {
	Content string `json:"content"`
}

func main() {
	loadEnv()

	r := mux.NewRouter()
	r.HandleFunc("/recognize", recognize).Methods("POST")
	http.ListenAndServe(":4000", r)
}

func recognize(w http.ResponseWriter, r *http.Request) {
	var audioBody AudioBody

	err := json.NewDecoder(r.Body).Decode(&audioBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	content := audioBody.Content

	result, err := bhashini.RecognizeAndTranslate(content)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write([]byte(result))
}

// func google(content string) (string, error) {
// 	ctx := context.Background()
// 	client, err := speech.NewClient(ctx)

// 	if err != nil {
// 		fmt.Print(err)

// 		return "", err
// 	}

// 	defer client.Close()

// 	config := speechpb.RecognitionConfig{
// 		Encoding:                 speechpb.RecognitionConfig_FLAC,
// 		SampleRateHertz:          16000,
// 		LanguageCode:             "en-IN",
// 		AlternativeLanguageCodes: []string{"hi-IN"},
// 	}

// 	audio := speechpb.RecognitionAudio{
// 		AudioSource: &speechpb.RecognitionAudio_Content{
// 			Content: []byte(content),
// 		},
// 	}
// 	req := &speechpb.RecognizeRequest{
// 		Config: &config,
// 		Audio:  &audio,
// 	}

// 	resp, err := client.Recognize(ctx, req)

// 	if err != nil {
// 		fmt.Print(err)

// 		return "", err
// 	}

// 	results := resp.Results
// 	var result string

// 	if len(results) > 0 {
// 		alternatives := results[0].Alternatives

// 		if len(alternatives) > 0 {
// 			result = alternatives[0].Transcript
// 		}
// 	}

// 	return result, nil
// }

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}
