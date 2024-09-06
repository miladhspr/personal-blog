package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"personal-blog/internal"
	"personal-blog/middlewares"
	"time"
)

func main() {
	// Public Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/article", articleHandler)

	// Private Routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/", adminHandler)
	adminMux.HandleFunc("/admin/new", newArticleHandler)
	adminMux.HandleFunc("/admin/edit", editArticleHandler)
	adminMux.HandleFunc("/admin/delete", deleteArticleHandler)

	mux.Handle("/admin/", middlewares.BasicAuthMiddleware(adminMux))

	fmt.Println("server on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := internal.LoadArticles()
	if err != nil {
		http.Error(w, "unable to load articles", http.StatusInternalServerError)
		return
	}
	temp, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "unable to load template", http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, articles)
	if err != nil {
		http.Error(w, "unable to render template", http.StatusInternalServerError)
		return
	}
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Article ID is required", http.StatusBadRequest)
		return
	}

	article, err := internal.LoadArticleByID(id)
	if err != nil {
		http.Error(w, "Unable to load article", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/article.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, article)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := internal.LoadArticles()
	if err != nil {
		http.Error(w, "Unable to load articles", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, articles)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

func newArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/new_article.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
		return
	}
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := r.FormValue("date")
		article := internal.Article{
			ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:   title,
			Content: content,
			Date:    date,
		}
		fileName := fmt.Sprintf("data/%s.json", article.ID)
		file, err := os.Create(fileName)
		if err != nil {
			http.Error(w, "Unable to save article", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = json.NewEncoder(file).Encode(article)
		if err != nil {
			http.Error(w, "Unable to save article", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Article not found", http.StatusNotFound)
		}
		article, err := internal.LoadArticleByID(id)
		if err != nil {
			http.Error(w, "Unable to load article", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/edit_article.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, article)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
		return
	}
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := r.FormValue("date")

		article := internal.Article{
			ID:      id,
			Title:   title,
			Content: content,
			Date:    date,
		}

		fileName := fmt.Sprintf("data/%d.json", article.ID)
		file, err := os.Create(fileName)
		if err != nil {
			http.Error(w, "Unable to save article", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = json.NewEncoder(file).Encode(article)
		if err != nil {
			http.Error(w, "Unable to save article", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Article ID is required", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("data/%s.json", id)
	err := os.Remove(fileName)
	if err != nil {
		http.Error(w, "Unable to delete article", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
