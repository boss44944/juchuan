package main

import "time"

type History struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Direction string    `json:"direction"`
	Content   string    `json:"content"`
	Filename  string    `json:"filename"`
	Filepath  string    `json:"filepath"`
	Mime      string    `json:"mime"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
