package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/russross/blackfriday/v2"
	"go_day06/pkg/entities/articles"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) error
}

type PostUseCase interface {
	CreateArticle(ctx context.Context, title, content string) error
	GetArticles(ctx context.Context, limit, offset int) ([]articles.Article, int, error)
	GetArticleByID(ctx context.Context, id int) (*articles.Article, error)
}

type Adapter struct {
	auth AuthUseCase
	post PostUseCase
}

func New(a AuthUseCase, p PostUseCase) *Adapter {
	return &Adapter{auth: a, post: p}
}

func StartServer(port int, a *Adapter) error {
	route := chi.NewMux()
	rl := NewRateLimiter(100, 1*time.Second)

	fileServer := http.FileServer(http.Dir("static"))
	route.Use(rl.RateLimiting)

	route.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	route.Get("/admin", a.LoginPage)
	route.Post("/admin", a.Login)
	route.Get("/admin/post", a.AuthMiddleware(a.PostPage))
	route.Post("/admin/post", a.PostArticle)
	route.Get("/article/{id}", a.ArticlePage)
	route.Get("/", a.MainPage)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), route))

	return nil
}

func (a *Adapter) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *Adapter) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" || a.auth.VerifyToken(r.Context(), cookie.Value) != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Adapter) Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	password := r.FormValue("password")

	token, err := a.auth.Login(r.Context(), user, password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/admin",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/admin/post", http.StatusSeeOther)
		return
	}
}

func (a *Adapter) PostPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../templates/postPage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func truncate(s string, length int) string {
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}

func substr(s string, start int, end int) string {
	if len(s) < start {
		return ""
	}
	if len(s) < end {
		return s[start:]
	}
	return s[start:end]
}

func (a *Adapter) MainPage(w http.ResponseWriter, r *http.Request) {
	//var articles []articles.Article
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * 3

	posts, totalPosts, err := a.post.GetArticles(r.Context(), 3, offset)
	pageData := struct {
		Articles    []articles.Article
		CurrentPage int
		TotalPages  int
	}{
		Articles:    posts,
		CurrentPage: page,
		TotalPages:  (totalPosts + 2) / 3,
	}

	funcMap := template.FuncMap{
		"truncate": truncate,
		"add":      add,
		"sub":      sub,
	}
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("../../templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *Adapter) PostArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	err := a.post.CreateArticle(r.Context(), title, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Adapter) ArticlePage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/article/"):])
	if err != nil || id < 1 {
		http.Error(w, "Page not found", http.StatusBadRequest)
		return
	}
	article, err := a.post.GetArticleByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseFiles("../../templates/article.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlContent := blackfriday.Run([]byte(article.Content))

	pageData := struct {
		Title   string
		Content template.HTML
	}{
		Title:   article.Title,
		Content: template.HTML(htmlContent),
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
