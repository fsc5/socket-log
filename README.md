# socket-log
This is a simple Keylogger Server which takes keystrokes over websockets.

It has two endpoints:
+ /ws
This is the websocket endpoint. It will accept every valid websocket request. The hole text of an message will be parsed as an keystroke
+ /get
You can get an json representation of the keystrokes over this endpoint. The keystrokes are organized by time.

**Pls dont use this project in any illegal way!**

I use this project just for learning purposes. So every suggestion for improvement is appreciated ;)

(I'm also not a native English speaker so grammar or spelling corrections are also appreciated)
