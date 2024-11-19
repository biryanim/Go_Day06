package articles

import "time"

type Article struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}
