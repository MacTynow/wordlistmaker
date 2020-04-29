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

type result struct {
	English string
	Hanzi   string
	Pinyin  []string
}

func translateWord(ctx context.Context, client translate.Client, word string, targetLang string) string {
	target, err := language.Parse(targetLang)
	translations, err := client.Translate(ctx, []string{word}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	return translations[0].Text
}

// Translate gets translation from api service
func Translate(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Words []string `json:"words"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}
	if len(data.Words) == 0 {
		log.Fatal("Data is empty")
	}

	ctx := context.Background()

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
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
			r.Hanzi = translateWord(ctx, *client, word, "zh")
			r.Pinyin = pinyin.LazyPinyin(r.Hanzi, a)
		} else if baseLang == "zh-CN" || baseLang == "zh-TW" {
			r.Hanzi = word
			r.Pinyin = pinyin.LazyPinyin(word, a)
			r.English = strings.ToLower(translateWord(ctx, *client, word, "en"))
		} else {
			log.Printf("Unsupported language")
			r.Hanzi = word
			r.Pinyin = pinyin.LazyPinyin(r.Hanzi, a)
			r.English = word
		}

		out = append(out, r)
	}

	for _, o := range out {
		fmt.Fprintf(w, "%v, %v, %v\n", o.English, o.Hanzi, strings.Join(o.Pinyin, " "))
	}
}
