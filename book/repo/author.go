package repo

import (
	"database/sql"

	"github.com/reeechart/booql/book/models"
)

var (
	authorGetAll  = "GetAllAuthorsQuery"
	authorGetById = "GetAuthorByIdQuery"
	authorInsert  = "InsertAuthorQuery"
	authorUpdate  = "UpdateAuthorQuery"
	authorQueries = map[string]string{
		authorGetAll:  `SELECT * FROM public.author`,
		authorGetById: `SELECT * FROM public.author WHERE id = $1`,
		authorInsert:  `INSERT INTO public.author (name) VALUES ($1) RETURNING id, name`,
		authorUpdate:  `UPDATE public.author SET name = $2 WHERE id = $1 RETURNING id, name`,
	}
)

type AuthorRepo struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func NewAuthorRepo(db *sql.DB) *AuthorRepo {
	stmts := make(map[string]*sql.Stmt, len(authorQueries))
	for queryKey, query := range authorQueries {
		stmt, err := db.Prepare(query)
		if err != nil {
			panic(err)
		}
		stmts[queryKey] = stmt
	}

	return &AuthorRepo{db: db, stmts: stmts}
}

func (repo *AuthorRepo) ListAuthors() ([]models.Author, error) {
	authors := []models.Author{}
	rows, err := repo.stmts[authorGetAll].Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		author := models.Author{}
		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (repo *AuthorRepo) GetAuthorById(id int32) (*models.Author, error) {
	rows, err := repo.stmts[authorGetById].Query(id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	author := models.Author{}
	err = rows.Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (repo *AuthorRepo) AddAuthor(name string) (*models.Author, error) {
	rows, err := repo.stmts[authorInsert].Query(name)
	if err != nil {
		return nil, err
	}

	rows.Next()
	author := models.Author{}
	err = rows.Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (repo *AuthorRepo) UpdateAuthor(id int32, name string) (*models.Author, error) {
	rows, err := repo.stmts[authorUpdate].Query(id, name)
	if err != nil {
		return nil, err
	}

	rows.Next()
	author := models.Author{}
	err = rows.Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}
