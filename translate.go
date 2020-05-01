package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/language"
)

type dataInput struct {
	Words        []string `json:"words"`
	TargetLang   string   `json:"target"`
	Romanization bool     `json:"romanization"`
}

type result struct {
	English      string
	Translation  string
	Romanization string
}

func translateWord(ctx context.Context, client translate.Client, word string, targetLang string) string {
	target, err := language.Parse(targetLang)
	translations, err := client.Translate(ctx, []string{word}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	return translations[0].Text
}

func processWordList(data dataInput) []result {
	ctx := context.Background()
	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Tone
	out := []result{}

	for _, word := range data.Words {
		r := result{}

		lang, err := client.DetectLanguage(ctx, []string{word})
		if err != nil {
			log.Fatalf("DetectLanguage: %v", err)
		}
		if len(lang) == 0 || len(lang[0]) == 0 {
			log.Fatalf("DetectLanguage return value empty")
		}

		baseLang := lang[0][0].Language.String()

		if baseLang == "en" {
			r.English = word
			r.Translation = translateWord(ctx, *client, word, data.TargetLang)
			r.Romanization = strings.Join(pinyin.LazyPinyin(r.Translation, pinyinArgs), " ")
		} else if baseLang == data.TargetLang {
			r.English = strings.ToLower(translateWord(ctx, *client, word, "en"))
			r.Translation = word
			r.Romanization = strings.Join(pinyin.LazyPinyin(word, pinyinArgs), " ")
		} else {
			log.Printf("Unsupported language")
			r.Translation = word
			r.Romanization = strings.Join(pinyin.LazyPinyin(r.Translation, pinyinArgs), " ")
			r.English = word
		}

		out = append(out, r)
	}

	return out
}

// Translate gets translation from api service
func Translate(w http.ResponseWriter, r *http.Request) {
	data := dataInput{
		TargetLang:   "zh-CN",
		Romanization: true,
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}
	if len(data.Words) == 0 {
		log.Fatal("Data is empty")
	}

	out := processWordList(data)

	for _, o := range out {
		if data.Romanization {
			fmt.Fprintf(w, "%v, %v, %v\n", o.English, o.Translation, o.Romanization)
		} else {
			fmt.Fprintf(w, "%v, %v\n", o.English, o.Translation)
		}
	}

	// csvWriter := csv.NewWriter(os.Stdout)

	// for _, o := range out {
	// 	if err := csvWriter.Write([]string{o.English, o.Translation, o.Romanization}); err != nil {
	// 		log.Fatalln("error writing record to csv:", err)
	// 	}
	// }

	// csvWriter.Flush()
}
