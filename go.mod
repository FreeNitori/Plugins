module git.randomchars.net/FreeNitori/Plugins

go 1.15

require (
	git.randomchars.net/FreeNitori/EmbedUtil v1.0.1
	git.randomchars.net/FreeNitori/FreeNitori v1.12.11
	git.randomchars.net/FreeNitori/Log v1.0.0
	git.randomchars.net/FreeNitori/Multiplexer v1.0.7
	github.com/bwmarrin/discordgo v0.23.2
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376 // indirect
)

replace git.randomchars.net/FreeNitori/FreeNitori v1.12.11 => ../
