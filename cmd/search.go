package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	// Kullanıcıların sohbet geçmişini saklamak için bir harita ve kilit oluşturuyoruz
	userHistories = make(map[string][]string)
	mu            sync.Mutex
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Args:  cobra.MinimumNArgs(1), // Minimum 1 arg required
	Run: func(cmd *cobra.Command, args []string) {
		userID := "exampleUserID" // Kullanıcı kimliği, burada sabit olarak verilmiş durumda
		response := GetResponse(userID, args)
		fmt.Println("Generated Response:", response)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func GetResponse(userID string, args []string) []string {
	// Sohbet geçmişine erişim için kilit kullanıyoruz
	mu.Lock()
	defer mu.Unlock()

	// Gelen argümanları bir cümle haline getiriyoruz
	userArgs := strings.Join(args[0:], " ")

	// Kullanıcının daha önceki konuşmalarını alıyoruz
	history, exists := userHistories[userID]
	if !exists {
		history = []string{}
	}

	// Gelen yeni soruyu geçmişe ekle
	history = append(history, "User: "+userArgs)

	// Sohbet geçmişini tek bir bağlam haline getir
	contextText := strings.Join(history, "\n")

	// Google AI client'ını başlat
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyBKN5QfkEn7YJhCVhKFI0WjG__iVaNKUsM"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Modeli seç ve kullanıcının sorusunu işleme koy
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(contextText))
	if err != nil {
		log.Fatal(err)
	}

	// İlk adayı alıyoruz ve `Part` içeriğini elde ediyoruz
	if len(resp.Candidates) > 0 {
		firstPart := resp.Candidates[0].Content.Parts[0]

		// Type assertion kullanarak `firstPart`'ın `Text` tipi olup olmadığını kontrol edelim
		if text, ok := firstPart.(genai.Text); ok {
			generatedText := string(text)
			fmt.Println("Generated Text:", generatedText)

			// Yanıtı da geçmişe ekle
			history = append(history, "AI: "+generatedText)

			// Güncellenmiş sohbet geçmişini kaydet
			userHistories[userID] = history

			return []string{generatedText}
		}

		// Eğer `Part` tipi `Text` değilse
		log.Println("Unsupported Part type. Only Text type is supported.")
		return []string{"Desteklenmeyen içerik türü, sadece metin destekleniyor."}
	}

	// Eğer hiçbir aday yoksa, bu durumu da işleyelim
	return []string{"Bir şeyler ters gitti, lütfen tekrar deneyin."}
}

func HandleAPIRequest(userID string, query string) []string {
	return GetResponse(userID, []string{query})
}
