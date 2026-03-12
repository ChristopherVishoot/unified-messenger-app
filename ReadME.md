# Problem Statement
I currently use WhatsApp, Signal, and Telegram to keep in touch with people. It would be nice to have the following -

- A unified UI to view messages
- Able to reply to messages
- An AI agent to schedule calls (zoom meetings, google meets, etc)
- Able to add a `messenger app` to the unified UI

* WhatsApp (for the prototype)
* Signal
* Telegram

Storage -

The messages will be stored for a maximum of 30 days by default, and the number of days can be updated by the user admin.

Instead of using Postgres, we use a Redis Cache (with RediSearch), storing messages only in memory and not to disk.

We can still perform advanced search queries on the data using RediSearch.

Data lifecycle:
- Auto-deleted after retention period (30 days default)
- Lost on server/pod restart (no persistence)
- Users can call `clear_my_data` to immediately delete all their messages
- Data encrypted in transit (TLS) between services

Auth & Authorization:
- User auth: OAuth2 Proxy (Google/GitHub login)
- Agent auth: Scoped API key per user (stored in Kubernetes Secret)
- Service-to-service: Kubernetes ServiceAccount tokens

MCP tool restrictions:
- Agent has **read-only** access to messages
- `search_messages` — scoped to authenticated user's data only
- `get_conversation` — scoped to authenticated user's data only
- No write/delete permissions for the agent

AI Agent:
- Local LLM via Ollama (`llama3` or `mistral`)
- Runs on `http://localhost:11434`
- Connects to MCP server for read-only message access
- Handles scheduling requests based on message context

## Redis Data Model

Redis is a key-value store. Each message is stored as a Hash with a UUID as the key.

**Key Format:** `message:{uuid}`

**Value (Hash fields):**

| Field | Description |
|-------|-------------|
| `user_id` | Owner of the message (for scoping queries to authenticated user) |
| `conversation_id` | Groups messages into conversations (for `get_conversation` tool) |
| `sender_id` | Who sent the message |
| `sender_name` | Display name for UI |
| `platform` | `whatsapp`, `signal`, or `telegram` |
| `direction` | `inbound` or `outbound` |
| `message` | Actual message content |
| `timestamp` | When the message was inserted into the cache |
| `message_id` | Original message ID from the messenger platform (if available) |
| `hash` | Hash value of the inserted object (for deduplication) |

**TTL:** Each key has a TTL of 30 days (2592000 seconds) by default, auto-deleted by Redis.

## RediSearch Index

RediSearch is a module that adds secondary indexing and full-text search on top of Redis data.

**Index Definition:**
```
FT.CREATE idx:messages ON HASH PREFIX 1 message:
  SCHEMA 
    user_id TAG
    conversation_id TAG
    sender_id TAG
    platform TAG
    direction TAG
    message TEXT
    timestamp NUMERIC SORTABLE
```

**Example Queries:**
- Search user's messages: `FT.SEARCH idx:messages "@user_id:{user123} @message:zoom" SORTBY timestamp DESC`
- Get conversation: `FT.SEARCH idx:messages "@user_id:{user123} @conversation_id:{conv456}" SORTBY timestamp ASC`
