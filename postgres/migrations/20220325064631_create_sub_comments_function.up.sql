CREATE OR REPLACE FUNCTION comments (p_id int)
RETURNS jsonb[] AS $$
BEGIN
RETURN (
    SELECT array_agg(row_to_json(t)) FROM
        (
            SELECT c.*, comments(c.id) AS comments, u.email, u.name, u.profile_picture
            FROM comment AS c
            LEFT JOIN user_account AS u ON u.id = c.user_id
            WHERE c.parent_id = p_id
        ) t
);
END;
$$ LANGUAGE plpgsql
