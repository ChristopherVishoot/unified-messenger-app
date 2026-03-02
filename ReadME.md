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
