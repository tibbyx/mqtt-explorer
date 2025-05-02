Build or Run (and build):

```bash
$ make build
$ make run
```


What the thing can do:

Connect to ip:port/topic
- listen to messages from every client
- send messages to every client
- ... that's it.


With this basics, we can do a lot.
- A chat system (like revolt or discord)
How? Well we have the ID. The rest we can store into the database.
The topics can be... like revolt/discord servers.


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
topic-list---+current-topic-messages-----clients------+
| main       | Currently at main         | main.c     |
| zion       +---------------------------+ virus      |
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
