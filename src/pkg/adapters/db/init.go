package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go_day06/pkg/config"
	"go_day06/pkg/entities/articles"
	"log"
)

type Adapter struct {
	db *pgx.Conn
}

func New(cfg *config.AppConfig) *Adapter {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")
	//defer conn.Close(context.Background())
	for _, sql := range cfg.SQLCommands {
		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			log.Fatalf("error executing SQL command: %v\n", err)
		}
	}
	return &Adapter{db: conn}
}

func Close(a *Adapter) {
	a.db.Close(context.Background())
}

func (a *Adapter) PostArticle(ctx context.Context, title string, content string) error {
	_, err := a.db.Exec(ctx, "INSERT INTO articles(title, content) VALUES($1, $2);", title, content)
	return err
}

func (a *Adapter) GetArticles(ctx context.Context, limit, offset int) ([]articles.Article, int, error) {
	rows, err := a.db.Query(ctx, "SELECT id, title, content, created_at FROM articles ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []articles.Article
	for rows.Next() {
		article := articles.Article{}
		if err = rows.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, article)
	}

	var totalArticles int
	err = a.db.QueryRow(ctx, "SELECT count(*) FROM articles").Scan(&totalArticles)
	if err != nil {
		return nil, 0, err
	}
	return posts, totalArticles, nil
}

func (a *Adapter) GetArticleById(ctx context.Context, id int) (*articles.Article, error) {
	rows, err := a.db.Query(ctx, "SELECT id, title, content, created_at FROM articles WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	article := articles.Article{}
	for rows.Next() {
		if err = rows.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt); err != nil {
			return nil, err
		}
	}
	return &article, nil
}
