package queries

var (
	PostQuery = map[string]string{
		"Get": `SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message FROM posts WHERE id = $1`,
		"Update": `UPDATE posts SET message = COALESCE(NULLIF($1, ''), message), 
		isEdited = CASE WHEN (isEdited = TRUE OR (isEdited = FALSE AND NULLIF($1, '') IS NOT NULL AND NULLIF($1, '') <> message)) 
		THEN TRUE ELSE FALSE END WHERE id = $2 
		RETURNING id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message`,
	}
)
