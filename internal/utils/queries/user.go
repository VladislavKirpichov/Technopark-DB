package queries

var (
	UserQuery = map[string]string{
		"Get":               `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1`,
		"Create":            `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`,
		"GetUsersByUserNOE": `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2`,
		"Update": `UPDATE users SET fullname = COALESCE(NULLIF($1, ''), fullname), 
			about = COALESCE(NULLIF($2, ''), about), 
			email = COALESCE(NULLIF($3, ''), email) WHERE nickname = $4 
			RETURNING nickname, fullname, about, email`,
	}
)
