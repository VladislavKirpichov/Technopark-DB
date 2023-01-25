package queries

var (
	ServiceQuery = map[string]string{
		"Clear":        `TRUNCATE users, forums, threads, votes, posts, forum_users`,
		"queryUsers":   `SELECT COUNT(*) FROM users`,
		"queryForums":  `SELECT COUNT(*) FROM forums`,
		"queryThreads": `SELECT COUNT(*) FROM threads`,
		"queryPosts":   `SELECT COUNT(*) FROM posts`,
	}
)
