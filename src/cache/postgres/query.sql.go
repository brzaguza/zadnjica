// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addEngineRank = `-- name: AddEngineRank :exec
INSERT INTO engine_ranks (
  engine_name, engine_rank, engine_page, engine_on_page_rank, result_id
) VALUES (
  $1, $2, $3, $4, $5
)
`

type AddEngineRankParams struct {
	EngineName       string
	EngineRank       int64
	EnginePage       int64
	EngineOnPageRank int64
	ResultID         int64
}

func (q *Queries) AddEngineRank(ctx context.Context, arg AddEngineRankParams) error {
	_, err := q.db.Exec(ctx, addEngineRank,
		arg.EngineName,
		arg.EngineRank,
		arg.EnginePage,
		arg.EngineOnPageRank,
		arg.ResultID,
	)
	return err
}

const addImageResult = `-- name: AddImageResult :exec
INSERT INTO image_results (
  image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, result_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
`

type AddImageResultParams struct {
	ImageOriginalHeight  int64
	ImageOriginalWidth   int64
	ImageThumbnailHeight int64
	ImageThumbnailWidth  int64
	ImageThumbnailUrl    string
	ImageSource          string
	ImageSourceUrl       string
	ResultID             int64
}

func (q *Queries) AddImageResult(ctx context.Context, arg AddImageResultParams) error {
	_, err := q.db.Exec(ctx, addImageResult,
		arg.ImageOriginalHeight,
		arg.ImageOriginalWidth,
		arg.ImageThumbnailHeight,
		arg.ImageThumbnailWidth,
		arg.ImageThumbnailUrl,
		arg.ImageSource,
		arg.ImageSourceUrl,
		arg.ResultID,
	)
	return err
}

const addResult = `-- name: AddResult :one
INSERT INTO results (
  query, url, rank, score, title, description
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id
`

type AddResultParams struct {
	Query       string
	Url         string
	Rank        int64
	Score       float64
	Title       string
	Description string
}

func (q *Queries) AddResult(ctx context.Context, arg AddResultParams) (int64, error) {
	row := q.db.QueryRow(ctx, addResult,
		arg.Query,
		arg.Url,
		arg.Rank,
		arg.Score,
		arg.Title,
		arg.Description,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteAllResultsOlderThanXDays = `-- name: DeleteAllResultsOlderThanXDays :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 day' * $1
`

func (q *Queries) DeleteAllResultsOlderThanXDays(ctx context.Context, dollar_1 interface{}) error {
	_, err := q.db.Exec(ctx, deleteAllResultsOlderThanXDays, dollar_1)
	return err
}

const deleteAllResultsOlderThanXHours = `-- name: DeleteAllResultsOlderThanXHours :exec
DELETE FROM results
WHERE created_at < NOW() - INTERVAL '1 hour' * $1
`

func (q *Queries) DeleteAllResultsOlderThanXHours(ctx context.Context, dollar_1 interface{}) error {
	_, err := q.db.Exec(ctx, deleteAllResultsOlderThanXHours, dollar_1)
	return err
}

const deleteAllResultsWithQuery = `-- name: DeleteAllResultsWithQuery :exec
DELETE FROM results
WHERE query = $1
`

func (q *Queries) DeleteAllResultsWithQuery(ctx context.Context, query string) error {
	_, err := q.db.Exec(ctx, deleteAllResultsWithQuery, query)
	return err
}

const getImageResultsByQueryAndEngineWithEngineRanks = `-- name: GetImageResultsByQueryAndEngineWithEngineRanks :many
SELECT results.id, query, url, rank, score, title, description, created_at, image_results.id, image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, image_results.result_id, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, engine_ranks.result_id FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetImageResultsByQueryAndEngineWithEngineRanksParams struct {
	Query      string
	EngineName string
}

type GetImageResultsByQueryAndEngineWithEngineRanksRow struct {
	ID                   int64
	Query                string
	Url                  string
	Rank                 int64
	Score                float64
	Title                string
	Description          string
	CreatedAt            pgtype.Timestamp
	ID_2                 int64
	ImageOriginalHeight  int64
	ImageOriginalWidth   int64
	ImageThumbnailHeight int64
	ImageThumbnailWidth  int64
	ImageThumbnailUrl    string
	ImageSource          string
	ImageSourceUrl       string
	ResultID             int64
	ID_3                 int64
	EngineName           string
	EngineRank           int64
	EnginePage           int64
	EngineOnPageRank     int64
	ResultID_2           int64
}

func (q *Queries) GetImageResultsByQueryAndEngineWithEngineRanks(ctx context.Context, arg GetImageResultsByQueryAndEngineWithEngineRanksParams) ([]GetImageResultsByQueryAndEngineWithEngineRanksRow, error) {
	rows, err := q.db.Query(ctx, getImageResultsByQueryAndEngineWithEngineRanks, arg.Query, arg.EngineName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetImageResultsByQueryAndEngineWithEngineRanksRow
	for rows.Next() {
		var i GetImageResultsByQueryAndEngineWithEngineRanksRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.ImageOriginalHeight,
			&i.ImageOriginalWidth,
			&i.ImageThumbnailHeight,
			&i.ImageThumbnailWidth,
			&i.ImageThumbnailUrl,
			&i.ImageSource,
			&i.ImageSourceUrl,
			&i.ResultID,
			&i.ID_3,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes = `-- name: GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT results.id, query, url, rank, score, title, description, created_at, image_results.id, image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, image_results.result_id, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, engine_ranks.result_id FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2 AND created_at > NOW() - INTERVAL '1 minute' * $3
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesParams struct {
	Query      string
	EngineName string
	Column3    interface{}
}

type GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow struct {
	ID                   int64
	Query                string
	Url                  string
	Rank                 int64
	Score                float64
	Title                string
	Description          string
	CreatedAt            pgtype.Timestamp
	ID_2                 int64
	ImageOriginalHeight  int64
	ImageOriginalWidth   int64
	ImageThumbnailHeight int64
	ImageThumbnailWidth  int64
	ImageThumbnailUrl    string
	ImageSource          string
	ImageSourceUrl       string
	ResultID             int64
	ID_3                 int64
	EngineName           string
	EngineRank           int64
	EnginePage           int64
	EngineOnPageRank     int64
	ResultID_2           int64
}

func (q *Queries) GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes(ctx context.Context, arg GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesParams) ([]GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow, error) {
	rows, err := q.db.Query(ctx, getImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes, arg.Query, arg.EngineName, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow
	for rows.Next() {
		var i GetImageResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.ImageOriginalHeight,
			&i.ImageOriginalWidth,
			&i.ImageThumbnailHeight,
			&i.ImageThumbnailWidth,
			&i.ImageThumbnailUrl,
			&i.ImageSource,
			&i.ImageSourceUrl,
			&i.ResultID,
			&i.ID_3,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImageResultsByQueryWithEngineRanks = `-- name: GetImageResultsByQueryWithEngineRanks :many
SELECT results.id, query, url, rank, score, title, description, created_at, image_results.id, image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, image_results.result_id, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, engine_ranks.result_id FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetImageResultsByQueryWithEngineRanksRow struct {
	ID                   int64
	Query                string
	Url                  string
	Rank                 int64
	Score                float64
	Title                string
	Description          string
	CreatedAt            pgtype.Timestamp
	ID_2                 int64
	ImageOriginalHeight  int64
	ImageOriginalWidth   int64
	ImageThumbnailHeight int64
	ImageThumbnailWidth  int64
	ImageThumbnailUrl    string
	ImageSource          string
	ImageSourceUrl       string
	ResultID             int64
	ID_3                 int64
	EngineName           string
	EngineRank           int64
	EnginePage           int64
	EngineOnPageRank     int64
	ResultID_2           int64
}

func (q *Queries) GetImageResultsByQueryWithEngineRanks(ctx context.Context, query string) ([]GetImageResultsByQueryWithEngineRanksRow, error) {
	rows, err := q.db.Query(ctx, getImageResultsByQueryWithEngineRanks, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetImageResultsByQueryWithEngineRanksRow
	for rows.Next() {
		var i GetImageResultsByQueryWithEngineRanksRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.ImageOriginalHeight,
			&i.ImageOriginalWidth,
			&i.ImageThumbnailHeight,
			&i.ImageThumbnailWidth,
			&i.ImageThumbnailUrl,
			&i.ImageSource,
			&i.ImageSourceUrl,
			&i.ResultID,
			&i.ID_3,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImageResultsByQueryWithEngineRanksNotOlderThanXminutes = `-- name: GetImageResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT results.id, query, url, rank, score, title, description, created_at, image_results.id, image_original_height, image_original_width, image_thumbnail_height, image_thumbnail_width, image_thumbnail_url, image_source, image_source_url, image_results.result_id, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, engine_ranks.result_id FROM results
JOIN image_results ON results.id = image_results.result_id
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND created_at > NOW() - INTERVAL '1 minute' * $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesParams struct {
	Query   string
	Column2 interface{}
}

type GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesRow struct {
	ID                   int64
	Query                string
	Url                  string
	Rank                 int64
	Score                float64
	Title                string
	Description          string
	CreatedAt            pgtype.Timestamp
	ID_2                 int64
	ImageOriginalHeight  int64
	ImageOriginalWidth   int64
	ImageThumbnailHeight int64
	ImageThumbnailWidth  int64
	ImageThumbnailUrl    string
	ImageSource          string
	ImageSourceUrl       string
	ResultID             int64
	ID_3                 int64
	EngineName           string
	EngineRank           int64
	EnginePage           int64
	EngineOnPageRank     int64
	ResultID_2           int64
}

func (q *Queries) GetImageResultsByQueryWithEngineRanksNotOlderThanXminutes(ctx context.Context, arg GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesParams) ([]GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesRow, error) {
	rows, err := q.db.Query(ctx, getImageResultsByQueryWithEngineRanksNotOlderThanXminutes, arg.Query, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesRow
	for rows.Next() {
		var i GetImageResultsByQueryWithEngineRanksNotOlderThanXminutesRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.ImageOriginalHeight,
			&i.ImageOriginalWidth,
			&i.ImageThumbnailHeight,
			&i.ImageThumbnailWidth,
			&i.ImageThumbnailUrl,
			&i.ImageSource,
			&i.ImageSourceUrl,
			&i.ResultID,
			&i.ID_3,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultsByQueryAndEngineWithEngineRanks = `-- name: GetResultsByQueryAndEngineWithEngineRanks :many
SELECT results.id, query, url, rank, score, title, description, created_at, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, result_id FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetResultsByQueryAndEngineWithEngineRanksParams struct {
	Query      string
	EngineName string
}

type GetResultsByQueryAndEngineWithEngineRanksRow struct {
	ID               int64
	Query            string
	Url              string
	Rank             int64
	Score            float64
	Title            string
	Description      string
	CreatedAt        pgtype.Timestamp
	ID_2             int64
	EngineName       string
	EngineRank       int64
	EnginePage       int64
	EngineOnPageRank int64
	ResultID         int64
}

func (q *Queries) GetResultsByQueryAndEngineWithEngineRanks(ctx context.Context, arg GetResultsByQueryAndEngineWithEngineRanksParams) ([]GetResultsByQueryAndEngineWithEngineRanksRow, error) {
	rows, err := q.db.Query(ctx, getResultsByQueryAndEngineWithEngineRanks, arg.Query, arg.EngineName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetResultsByQueryAndEngineWithEngineRanksRow
	for rows.Next() {
		var i GetResultsByQueryAndEngineWithEngineRanksRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes = `-- name: GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes :many
SELECT results.id, query, url, rank, score, title, description, created_at, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, result_id FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2 AND created_at > NOW() - INTERVAL '1 minute' * $3
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesParams struct {
	Query      string
	EngineName string
	Column3    interface{}
}

type GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow struct {
	ID               int64
	Query            string
	Url              string
	Rank             int64
	Score            float64
	Title            string
	Description      string
	CreatedAt        pgtype.Timestamp
	ID_2             int64
	EngineName       string
	EngineRank       int64
	EnginePage       int64
	EngineOnPageRank int64
	ResultID         int64
}

func (q *Queries) GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes(ctx context.Context, arg GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesParams) ([]GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow, error) {
	rows, err := q.db.Query(ctx, getResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutes, arg.Query, arg.EngineName, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow
	for rows.Next() {
		var i GetResultsByQueryAndEngineWithEngineRanksNotOlderThanXminutesRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultsByQueryWithEngineRanks = `-- name: GetResultsByQueryWithEngineRanks :many
SELECT results.id, query, url, rank, score, title, description, created_at, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, result_id FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetResultsByQueryWithEngineRanksRow struct {
	ID               int64
	Query            string
	Url              string
	Rank             int64
	Score            float64
	Title            string
	Description      string
	CreatedAt        pgtype.Timestamp
	ID_2             int64
	EngineName       string
	EngineRank       int64
	EnginePage       int64
	EngineOnPageRank int64
	ResultID         int64
}

func (q *Queries) GetResultsByQueryWithEngineRanks(ctx context.Context, query string) ([]GetResultsByQueryWithEngineRanksRow, error) {
	rows, err := q.db.Query(ctx, getResultsByQueryWithEngineRanks, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetResultsByQueryWithEngineRanksRow
	for rows.Next() {
		var i GetResultsByQueryWithEngineRanksRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultsByQueryWithEngineRanksNotOlderThanXminutes = `-- name: GetResultsByQueryWithEngineRanksNotOlderThanXminutes :many
SELECT results.id, query, url, rank, score, title, description, created_at, engine_ranks.id, engine_name, engine_rank, engine_page, engine_on_page_rank, result_id FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND created_at > NOW() - INTERVAL '1 minute' * $2
ORDER BY results.rank ASC, engine_ranks.engine_rank ASC
`

type GetResultsByQueryWithEngineRanksNotOlderThanXminutesParams struct {
	Query   string
	Column2 interface{}
}

type GetResultsByQueryWithEngineRanksNotOlderThanXminutesRow struct {
	ID               int64
	Query            string
	Url              string
	Rank             int64
	Score            float64
	Title            string
	Description      string
	CreatedAt        pgtype.Timestamp
	ID_2             int64
	EngineName       string
	EngineRank       int64
	EnginePage       int64
	EngineOnPageRank int64
	ResultID         int64
}

func (q *Queries) GetResultsByQueryWithEngineRanksNotOlderThanXminutes(ctx context.Context, arg GetResultsByQueryWithEngineRanksNotOlderThanXminutesParams) ([]GetResultsByQueryWithEngineRanksNotOlderThanXminutesRow, error) {
	rows, err := q.db.Query(ctx, getResultsByQueryWithEngineRanksNotOlderThanXminutes, arg.Query, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetResultsByQueryWithEngineRanksNotOlderThanXminutesRow
	for rows.Next() {
		var i GetResultsByQueryWithEngineRanksNotOlderThanXminutesRow
		if err := rows.Scan(
			&i.ID,
			&i.Query,
			&i.Url,
			&i.Rank,
			&i.Score,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.ID_2,
			&i.EngineName,
			&i.EngineRank,
			&i.EnginePage,
			&i.EngineOnPageRank,
			&i.ResultID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultsTTLByQuery = `-- name: GetResultsTTLByQuery :one
SELECT created_at FROM results
WHERE query = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetResultsTTLByQuery(ctx context.Context, query string) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, getResultsTTLByQuery, query)
	var created_at pgtype.Timestamp
	err := row.Scan(&created_at)
	return created_at, err
}

const getResultsTTLByQueryAndEngine = `-- name: GetResultsTTLByQueryAndEngine :one
SELECT created_at FROM results
JOIN engine_ranks ON results.id = engine_ranks.result_id
WHERE query = $1 AND engine_ranks.engine_name = $2
ORDER BY created_at DESC
LIMIT 1
`

type GetResultsTTLByQueryAndEngineParams struct {
	Query      string
	EngineName string
}

func (q *Queries) GetResultsTTLByQueryAndEngine(ctx context.Context, arg GetResultsTTLByQueryAndEngineParams) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, getResultsTTLByQueryAndEngine, arg.Query, arg.EngineName)
	var created_at pgtype.Timestamp
	err := row.Scan(&created_at)
	return created_at, err
}