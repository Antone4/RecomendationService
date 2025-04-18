package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recommendation-service/internal/database"

	"github.com/gin-gonic/gin"
)

type Clip struct {
	ID        string            `json:"clip_id"`
	Title     string            `json:"title,omitempty"`
	VideoPath string            `json:"video_path,omitempty"`
	ThumbPath string            `json:"thumb_path,omitempty"`
	Subs      map[string]string `json:"subs"`
	CreatedAt string            `json:"created_at,omitempty"`
	UpdatedAt string            `json:"updated_at,omitempty"`
}

type SearchResult struct {
	Results []struct {
		Source struct {
			MovieID   string `json:"movie_id"`
			ClipID    string `json:"clip_id"`
			Subtitles []struct {
				Text         string `json:"text"`
				PathToRusSub string `json:"path_to_rus_sub"`
				PathToEngSub string `json:"path_to_eng_sub"`
			} `json:"subtitles"`
		} `json:"_source"`
	} `json:"results"`
}

func RecommendHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Получить слова пользователя
	var userWords []database.UserWord
	database.DB.Preload("Word").Where("user_id = ?", userID).Find(&userWords)

	knownWords := map[string]bool{}
	for _, uw := range userWords {
		if uw.KnowledgeLevel >= 3 {
			knownWords[uw.Word.Word] = true
		}
	}

	// Получить популярные слова других пользователей
	var commonWords []string
	database.DB.
		Model(&database.Word{}).
		Select("words.word").
		Joins("JOIN user_words ON words.id = user_words.word_id").
		Where("user_words.knowledge_level >= 3").
		Group("words.word").
		Having("COUNT(user_words.user_id) >= 2"). // можно варьировать
		Find(&commonWords)

	// Исключить известные слова
	targetWords := []string{}
	for _, word := range commonWords {
		if !knownWords[word] {
			targetWords = append(targetWords, word)
		}
		if len(targetWords) >= 10 {
			break
		}
	}

	// Поиск клипов по словам
	uniqueClips := map[string]Clip{}
	for _, word := range targetWords {
		clips, err := searchClips(word)
		if err != nil {
			continue
		}
		for _, clip := range clips {
			if len(uniqueClips) >= 20 {
				break
			}
			if _, ok := uniqueClips[clip.ID]; !ok {
				uniqueClips[clip.ID] = clip
			}
		}
		if len(uniqueClips) >= 20 {
			break
		}
	}

	// Ответ
	result := []Clip{}
	for _, clip := range uniqueClips {
		result = append(result, clip)
	}
	c.JSON(http.StatusOK, result)
}

func searchClips(word string) ([]Clip, error) {
	url := fmt.Sprintf("http://localhost:8000/search/?query=%s", word)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}

	var clips []Clip
	for _, r := range sr.Results {
		clipID := r.Source.ClipID
		if clipID == "" {
			continue
		}

		var eng, rus string
		for _, sub := range r.Source.Subtitles {
			if sub.PathToEngSub != "" {
				eng = sub.PathToEngSub
			}
			if sub.PathToRusSub != "" {
				rus = sub.PathToRusSub
			}
		}

		clips = append(clips, Clip{
			ID:        clipID,
			VideoPath: "",
			ThumbPath: "",
			Subs: map[string]string{
				"eng": eng,
				"rus": rus,
			},
		})
	}

	return clips, nil
}
