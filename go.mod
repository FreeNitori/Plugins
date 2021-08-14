module git.randomchars.net/freenitori/plugins

go 1.15

require (
	git.randomchars.net/freenitori/embedutil v1.0.2
	git.randomchars.net/freenitori/freenitori/v2 v2.0.0
	git.randomchars.net/freenitori/log v1.0.0
	git.randomchars.net/freenitori/multiplexer v1.0.12
	github.com/bwmarrin/discordgo v0.23.2
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	layeh.com/gopus v0.0.0-20210501142526-1ee02d434e32
)

replace git.randomchars.net/freenitori/freenitori/v2 v2.0.0 => ../freenitori
