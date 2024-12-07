-- name: AddSong :one
insert into songs (
    name,
    release_date,
    text,
    link,
    group_name
) values (
          $1, $2, $3, $4, $5
         )
returning *;

-- name: DeleteSong :exec
delete from songs
where id = $1;

-- name: UpdateSong :one
update songs set
                 name = $2,
                 release_date = $3,
                 text = $4,
                 link = $5,
                 group_name = $6
where id = $1
returning *;

-- name: GetSong :one
SELECT release_date, text, link
FROM songs
where group_name = $1 and name = $2;

-- name: ListSongs :many
SELECT *
FROM songs
LIMIT $1
OFFSET $2;

-- name: GetSongLyrics :many
SELECT
    array_to_string(
            (string_to_array(text, E'\n\n'))[ (($2 - 1) * $3 + 1) : (($2 - 1) * $3 + $3) ],
            E'\n\n'
    ) AS paginated_verses
FROM songs
WHERE id = $1;