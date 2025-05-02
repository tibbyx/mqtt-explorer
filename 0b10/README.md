Build or Run (and build):

```bash
$ make build
$ make run
```


What mqtt can do:

Connect to ip:port/topic
- listen to messages from every client
- send messages to every client
- ... that's it.


What mqtt cannot do:
- Differenciate the clients
- ...
- Or can it? (Vsause music kicks in)
- The user who will use will simply write the messages, but we shall append the clientId to the messages..!
- Hah! Crazy! Insane even!


With this basics, we can do a lot.
- A chat system (like revolt or discord)
- ...
- How? Well we have the ID (In the messages).
- We can use database to display previously entered messages in topics by ids.
- The topics can be... like revolt/discord servers!


The design:

# Base With middle
```
+-----------------------------------------------------+
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                          .                          |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
+-----------------------------------------------------+
```
# Initial
```
+-----------------------------------------------------+
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                  form------------+                  |
|                  |               |                  |
|                  |       .       |                  |
|                  |               |                  |
|                  +---------------+                  |
|                                                     |
|                                                     |
|                                                     |
|                                                     |
+-----------------------------------------------------+
```
# Wrong Creds
```
+-----------------------------------------------------+
|                                                     |
|                                                     |
|                                                     |
|                                                     |
|                  red-form--------+                  |
|                  |               |                  |
|                  |       .       |                  |
|                  |               |                  |
|                  explanation-----+                  |
|                  |               |                  |
|                  +---------------+                  |
|                                                     |
|                                                     |
+-----------------------------------------------------+
```
# Entry
```
topic-list---+current-topic--------------clients------+
| main       | Currently at main         | main.c     |
| zion       messages--------------------+ virus      |
|            |                           | cogmind    |
|            |                           | soldier    |
|            |                           | slayer     |
|            |                           |            |
|            |C|main.c 11:33             |            |
|            |Playing the real resource  |            |
|            |management game...         |            |
|            |v|virus 11:34              |            |
|            |Zooming...                 |            |
|            type-messages---------------+            |
|            |> Longing for the surface  |            |
+------------+---------------------------+------------+
```
# If no topics
```
topic-list---connect-to-topic------------clients------+
|            |>                          | cogmind    |
|            +---------------------------+            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
|            |             .             |            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
|            |                           |            |
+------------+---------------------------+------------+
```
