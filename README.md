## simpleMessenger
Here is simple messenger api and app.
It supports 
- creating chats 
- joining chats
- writing messages to chat
- updating history (either unread messages or full history).

Chats are for multiple participants.
The app itself works like interactive menu manager: 
You choose options and fall down the next level of menu.

## How does it work?

Server/Client communication is implemented via RPC connection.
Client connects to server and posts some messages. Messages are written in Kafka's topics.
Clients themselves are Kafka's producers and consumers simultaneously.

Each user chooses id as well as for each chat its id is generated. Identifiers are kept into RedisDB to guarantee no identical ids have been chosen.

## Flags and Menu navigation Guide
When you just launched messenger app, you have two options: `exit` and `login`.
`login` flags are:
- `-user WhoIam` specifies your username that will be used in chats.
- `-ip 127.0.0.1` is the server you connect to.

After you've logged in you can

- `join` existing chat. `-chatId` specifies chat you join to.
- `create` new chat. It automatically joins you to the created chat. You will be given its id.
- `exit` which means you sign out and go back to previous menu level.

After you've got into chat you can:

- `postmsg`. After hitting `Enter` program reads your input until `\n`.
- `update` - Uploads unread messages into specified by `-file saveHere.txt` flag file.
- `hist` Uploads all chat's history into specified by `-file saveHere.txt` flag file.
- `exit`. You leave chat.



