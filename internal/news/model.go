package news

import "time"

type Article struct {
	Title    string
	Body     string
	Datetime time.Time
}
