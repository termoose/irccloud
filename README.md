# IRCCloud Terminal Client

Do you miss the R(etro) in IRC? Still have dreams about BitchX? Who doesn't!
For all of this to work you need an IRCCloud account. If you're not familiar with [IRCCloud](http://irccloud.com), then I think you should check it out!

![Live demo!](preview.gif)

## Navigation

- <kbd>Ctrl</kbd>+<kbd>Space</kbd>: Select channel
- <kbd>Ctrl</kbd>+<kbd>b</kbd>: Switch to channel with most recent activity
- <kbd>Home</kbd>, <kbd>Ctrl</kbd>+<kbd>a</kbd>: Move to the beginning of the line
- <kbd>End</kbd>, <kbd>Ctrl</kbd>+<kbd>e</kbd>: Move to the end of the line
- <kbd>Ctrl</kbd>+<kbd>k</kbd>: Delete from the cursor to the end of the line
- <kbd>Ctrl</kbd>+<kbd>w</kbd>: Delete the last word before the cursor
- <kbd>Ctrl</kbd>+<kbd>u</kbd>: Delete the entire line

## Installing

```bash
go get -u github.com/termoose/irccloud
```

Make sure `$GOPATH/bin` is added to the `$PATH` variable.

```bash
irccloud
```

On the first run a config file will be generated for you in `~/.config/irccloud`, update it with your IRCCloud username and password.

### Homebrew

Soon

### Flatpak (Linux)

Probably soon

## TODO
- ~Get text input working, return = send text to channel~
- ~Make sure channel join opens a new buffer~
- ~Properly quit on Control+C so it doesn't mess up the terminal~
- ~New query window opened for new private messages~
- ~Show latest active channels, sorted list of oldest to newest activity~
- Don't quit when we lose websocket connection, display "Reconnecting" dialogue and retry until Ctrl+C
- Create a loading bar based on `numbuffers` in `oob_include`
- Make auto-completion on nick names in the current channel buffer
- Add list of trigger names to config, include own nick by default?
- Special highlight for channels with trigger events in event bar
- ~Show timestamps on chat messages~
- ~Support nick change messages~
- ~Show non-message data in chat like join/leave events~
- Find out how to change channel view with alt+Fn or something similar
- Show operator/voice status in channel member list
- Get and show fancy ANSI splash art, BitchX style!
- Should we handle multiple servers? Does irccloud do this? How? (it kinda works)
- Add support for IRCCloud file upload
- Get image preview in terminal, ascii-render fallback in unsupported terminal
- Support IRCCloud's pastebin service, query users with very large text inputs if they want to use pastebin service