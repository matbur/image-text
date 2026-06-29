package version

// Commit is the git SHA, set at build time via -ldflags.
var Commit = "dev"

func init() {
	if len(Commit) > 7 {
		Commit = Commit[:7]
	}
}
