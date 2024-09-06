package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadArticles() ([]Article, error) {
	var articles []Article
	files, err := os.ReadDir("data")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			content, err := os.ReadFile(filepath.Join("data", file.Name()))
			if err != nil {
				return nil, err
			}

			var article Article
			err = json.Unmarshal(content, &article)
			if err != nil {
				return nil, err
			}
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func LoadArticleByID(id string) (Article, error) {
	files, err := os.ReadDir("data")
	if err != nil {
		return Article{}, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			content, err := os.ReadFile(filepath.Join("data", file.Name()))
			if err != nil {
				return Article{}, err
			}

			var article Article
			err = json.Unmarshal(content, &article)
			if err != nil {
				return Article{}, err
			}

			if article.ID == id {
				return article, nil
			}
		}
	}

	return Article{}, os.ErrNotExist
}
