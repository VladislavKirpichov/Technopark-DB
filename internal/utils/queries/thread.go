package queries

var (
	ThreadQuery = map[string]string{
		"GetBySlug":        `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE slug = $1`,
		"PostsCreate":      `INSERT INTO posts(parent, author, forum, thread, message, created) VALUES `,
		"CreatePostsTwo":   ` RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
		"VoteByID":         `INSERT INTO votes (nickname, thread, value) VALUES ($1, $2, $3) ON CONFLICT (nickname, thread) DO UPDATE SET value = $3`,
		"CreatePostsBatch": `INSERT INTO posts(parent, author, forum, thread, message, created) VALUES (NULLIF($1, 0), $2, $3, $4, $5, $6) RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
		"VoteBySlug": `INSERT INTO votes (nickname, thread, value) VALUES ($1, (SELECT id FROM threads WHERE slug=$2), $3) 
		ON CONFLICT (nickname, thread) DO UPDATE SET value = $3`,
		"UpdateByID": `UPDATE threads SET title = COALESCE(NULLIF($1, ''), title), 
		message = COALESCE(NULLIF($2, ''), message) WHERE id = $3 
		RETURNING id, COALESCE(slug, ''), author, forum, title, message, created, votes`,
		"GetByID": `SELECT id, COALESCE(slug, ''), author, forum, title, message, created, votes FROM threads WHERE id = $1`,
		"UpdateBySlug": `UPDATE threads SET title = COALESCE(NULLIF($1, ''), title), 
		message = COALESCE(NULLIF($2, ''), message) WHERE slug = $3 
		RETURNING id, slug, author, forum, title, message, created, votes`,
	}
)
