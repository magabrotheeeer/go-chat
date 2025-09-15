package domain

import "time"

type Chat struct {
	ID        string    `json:"id"`
	User1ID   string    `json:"user1_id"`
	User2ID   string    `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateChatRequest struct {
	User1ID string `json:"user1_id" binding:"required"`
	User2ID string `json:"user2_id" binding:"required"`
}
