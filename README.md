```
Currently there are four labs:

- old_server
    - gofiber as back-end and currently no frontend
- 0b1
    - fyne as both front-end and back-end.
    - I will probably abandon it though.
- 0b10
    - go-app as both front-end and back-end
    - it builds wasm and has the html + css flexibility
    - builds to web browsers
- 0b11
    - back-end          : gofiber
    - template engine   : django
    - front-end         : htmx
    - websocket         : yes
- 0b100
    - sse
    - I don't know how to make it work with htmx and I don't want to try React out.
- 0b101
    - Client shall ping the server for changes and if there is a change in state, it shall render it into the view.
    - This won't work with htmx. So React is the way.
    - For this to work smoothly, I will try to declare endpoints to the server.
