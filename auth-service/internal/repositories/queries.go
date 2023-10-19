package repositories

const (
	insertSessionData string = `
	INSERT INTO 
		session(sid,user_id,email,refresh_token,expired_at, created_at, updated_at) 
	VALUES 
		($1,$2,$3,$4,$5,$6) 
	RETURNING sid`

	selectSessionData = `
	SELECT 
		sid, user_id, email, is_valid 
	FROM 
		session 
	WHERE sid = $1 
	LIMIT 1`
)
