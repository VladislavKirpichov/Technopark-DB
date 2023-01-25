package queries

const (
	SortFlat       string = "flat"
	SortTree       string = "tree"
	SortParentTree string = "parent_tree"
)

var (
	DescSincePostQuery = map[string]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND id < $2 ORDER BY id DESC LIMIT $3",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) ORDER BY path DESC LIMIT $3",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
      AND parent IS NULL
      AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
    ORDER BY path[1] DESC
    LIMIT $3
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path[1] DESC, path[2:]`,
	}
	AscSincePostQuery = map[string]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND id > $2 ORDER BY id LIMIT $3",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) " +
			"ORDER BY path LIMIT $3",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
      AND parent IS NULL
      AND path[1] > (SELECT path[1] FROM posts WHERE id = $2)
    ORDER BY path[1]
    LIMIT $3
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path`,
	}
	DescNoSincePostQuery = map[string]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY id DESC LIMIT $2",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY path DESC LIMIT $2",
		SortParentTree: `
WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
    ORDER BY path[1] DESC
    LIMIT $2
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts WHERE thread = $1 AND path[1] IN (SELECT * FROM roots) ORDER BY path[1] DESC, path[2:]`,
	}
	AscNoSincePostQuery = map[string]string{
		SortFlat: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY id LIMIT $2",
		SortTree: "SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message " +
			"FROM posts WHERE thread = $1 ORDER BY path LIMIT $2\n",
		SortParentTree: `WITH roots AS (
    SELECT DISTINCT path[1]
    FROM posts
    WHERE thread = $1
    ORDER BY path[1]
    LIMIT $2
)
SELECT id, COALESCE(parent, 0), author, forum, thread, created, isEdited, message
FROM posts
WHERE thread = $1 AND path[1] IN (SELECT * FROM roots)
ORDER BY path`,
	}
)
