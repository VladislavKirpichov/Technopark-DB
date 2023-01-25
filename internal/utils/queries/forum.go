package queries

var (
	ForumQuery = map[string]string{
		"Create":                `INSERT INTO forums ("user", slug, title) VALUES ((SELECT nickname FROM users WHERE nickname = $3), $1, $2) RETURNING slug, title, "user", posts, threads`,
		"Get":                   `SELECT id, slug, title, "user", posts, threads FROM forums WHERE slug = $1`,
		"GetThreadsDesc":        ` AND created <= $2 ORDER BY created DESC LIMIT $3`,
		"GetThreadsSinceDesc":   ` ORDER BY created DESC LIMIT $2`,
		"GetThreadsNoDesc":      ` AND created >= $2 ORDER BY created LIMIT $3`,
		"GetThreadsSinceNoDesc": ` ORDER BY created LIMIT $2`,
		"GetThreads":            `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE forum = $1`,
		"GetUsers":              `SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users AS fu JOIN users AS u ON fu.nickname = u.nickname WHERE fu.forum = $1 `,
		"GetUsersDesc":          ` ORDER BY u.nickname DESC LIMIT $2`,
		"GetUsersSinceDesc":     ` AND u.nickname < $2 ORDER BY u.nickname DESC LIMIT $3`,
		"GetUsersNoDesc":        ` ORDER BY u.nickname LIMIT $2`,
		"GetUsersSinceNoDesc":   ` AND u.nickname > $2 ORDER BY u.nickname LIMIT $3`,
		"CreateThread": `INSERT INTO threads (slug, author, forum, title, message, created) VALUES (NULLIF($1, ''), (SELECT nickname FROM users WHERE nickname = $2), 
		(SELECT slug FROM forums WHERE slug = $3), $4, $5, $6) RETURNING id, $1, author, forum, title, message, created, votes`,
	}
)
