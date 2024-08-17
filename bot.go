package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/abadojack/whatlanggo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TranslationRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type TranslationResponse struct {
	Data string `json:"data"`
}

var (
	allowedGroups map[int64]bool
	allowedUsers  map[int64]bool
	ignoreLangs   map[string]bool
	targetLang    string
	botToken      string
	apiURL        string
)

func main() {
	// Parse command-line flags
	flag.StringVar(&botToken, "token", "", "Telegram Bot Token")
	flag.StringVar(&targetLang, "target", "ZH", "Target language for translation")
	flag.StringVar(&apiURL, "api", "http://127.0.0.1:1188/translate", "API URL for translation service")
	ignoreLangsFlag := flag.String("ignore", "ZH", "Comma-separated list of languages to ignore")
	allowedGroupsFlag := flag.String("groups", "", "Comma-separated list of allowed group IDs")
	allowedUsersFlag := flag.String("users", "", "Comma-separated list of allowed user IDs")
	flag.Parse()

	// Use environment variables if flags are not set
	if botToken == "" {
		botToken = os.Getenv("BOT_TOKEN")
	}
	if targetLang == "ZH" {
		targetLang = getEnv("TARGET_LANG", "ZH")
	}
	if apiURL == "http://127.0.0.1:1188/translate" {
		apiURL = getEnv("API_URL", "http://127.0.0.1:1188/translate")
	}
	if *ignoreLangsFlag == "ZH" {
		*ignoreLangsFlag = getEnv("IGNORE_LANGS", "ZH")
	}
	if *allowedGroupsFlag == "" {
		*allowedGroupsFlag = os.Getenv("ALLOWED_GROUPS")
	}
	if *allowedUsersFlag == "" {
		*allowedUsersFlag = os.Getenv("ALLOWED_USERS")
	}

	// Initialize configuration
	initConfig(*ignoreLangsFlag, *allowedGroupsFlag, *allowedUsersFlag)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Using API URL: %s", apiURL)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go handleMessage(bot, update.Message)
	}
}

func initConfig(ignoreLangsStr, allowedGroupsStr, allowedUsersStr string) {
	ignoreLangs = make(map[string]bool)
	for _, lang := range strings.Split(ignoreLangsStr, ",") {
		ignoreLangs[strings.TrimSpace(strings.ToUpper(lang))] = true
	}

	allowedGroups = make(map[int64]bool)
	for _, groupID := range strings.Split(allowedGroupsStr, ",") {
		if id, err := strconv.ParseInt(strings.TrimSpace(groupID), 10, 64); err == nil {
			allowedGroups[id] = true
		}
	}

	allowedUsers = make(map[int64]bool)
	for _, userID := range strings.Split(allowedUsersStr, ",") {
		if id, err := strconv.ParseInt(strings.TrimSpace(userID), 10, 64); err == nil {
			allowedUsers[id] = true
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if !isAllowed(message) {
		log.Printf("Message from non-whitelisted source: %d", message.Chat.ID)
		return
	}

	if !shouldIgnore(message.Text) {
		translatedText, err := translate(message.Text)
		if err != nil {
			log.Printf("Translation error: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, translatedText)
		msg.ReplyToMessageID = message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func isAllowed(message *tgbotapi.Message) bool {
	if len(allowedGroups) == 0 && len(allowedUsers) == 0 {
		return true
	}
	if message.Chat.Type == "private" {
		return allowedUsers[message.From.ID]
	} else {
		return allowedGroups[message.Chat.ID]
	}
}

func shouldIgnore(text string) bool {
	lang := whatlanggo.DetectLang(text)
	sourceLang := strings.ToUpper(lang.Iso6391())
	return ignoreLangs[sourceLang]
}

func translate(text string) (string, error) {
	requestBody, err := json.Marshal(TranslationRequest{
		Text:       text,
		SourceLang: "auto",
		TargetLang: targetLang,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var translationResp TranslationResponse
	err = json.Unmarshal(body, &translationResp)
	if err != nil {
		return "", err
	}

	return translationResp.Data, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
