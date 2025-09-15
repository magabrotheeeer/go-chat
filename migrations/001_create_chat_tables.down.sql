DROP INDEX idx_messages_chat_time;
DROP INDEX idx_chats_user2;
DROP INDEX idx_chats_user1;

DROP TABLE messages;
DROP TABLE chats;

DROP EXTENSION IF EXISTS "uuid-ossp";
