package postHandler

import (
	"context"
	"go_day06/pkg/entities/articles"
)

type storage interface {
	PostArticle(ctx context.Context, title string, content string) error
	GetArticles(ctx context.Context, limit, offset int) ([]articles.Article, int, error)
	GetArticleById(ctx context.Context, id int) (*articles.Article, error)
}

type PostHandler struct {
	s storage
}

func New(s storage) *PostHandler {
	return &PostHandler{
		s: s,
	}
}

func (p *PostHandler) CreateArticle(ctx context.Context, title, content string) error {
	return p.s.PostArticle(ctx, title, content)
}

func (p *PostHandler) GetArticles(ctx context.Context, limit, offset int) ([]articles.Article, int, error) {
	return p.s.GetArticles(ctx, limit, offset)
}

func (p *PostHandler) GetArticleByID(ctx context.Context, id int) (*articles.Article, error) {
	return p.s.GetArticleById(ctx, id)
}
