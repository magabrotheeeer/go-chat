CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user1_id VARCHAR(50) NOT NULL,
    user2_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user1_id, user2_id)
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    author_id VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chats_user1 ON chats(user1_id);
CREATE INDEX idx_chats_user2 ON chats(user2_id);
CREATE INDEX idx_messages_chat_time ON messages(chat_id, created_at DESC);
