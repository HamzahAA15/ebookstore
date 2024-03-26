package postgresql_test

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository/postgresql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBookRepository_GetBooks(t *testing.T) {
	books := []model.Book{
		{
			ID:         1,
			Title:      "test",
			Author:     "test",
			Price:      100,
			CategoryID: 1,
		},
	}

	type fields struct {
		data []model.Book
		err  error
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Book
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: books,
				err:  nil,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    books,
			wantErr: false,
		},
		{
			name: "SelectContext error",
			fields: fields{
				data: nil,
				err:  errors.New("error"),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			testDB := postgresql.NewBookRepository(sqlxDB)

			query := "SELECT id, title, author, price, category_id FROM books WHERE deleted_at IS NULL"

			mockExpectQuery := m.ExpectQuery(query)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "title", "author", "price", "category_id"})
				rowData := books
				row.AddRow(rowData[0].ID, rowData[0].Title, rowData[0].Author, rowData[0].Price, rowData[0].CategoryID)

				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetBooks(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestBookRepository_GetBookByID(t *testing.T) {
	book := model.Book{
		ID:         1,
		Title:      "test",
		Author:     "test",
		Price:      100,
		CategoryID: 1,
	}

	type fields struct {
		data model.Book
		err  error
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Book
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: book,
				err:  nil,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    book,
			wantErr: false,
		},
		{
			name: "SelectContext error",
			fields: fields{
				data: model.Book{},
				err:  errors.New("error"),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    model.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			testDB := postgresql.NewBookRepository(sqlxDB)

			query := "SELECT id, title, author, price, category_id FROM books WHERE id = $1 AND deleted_at IS NULL"

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.id)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "title", "author", "price", "category_id"})
				rowData := book
				row.AddRow(rowData.ID, rowData.Title, rowData.Author, rowData.Price, rowData.CategoryID)

				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetBookByID(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBookRepository_GetCategoryByID(t *testing.T) {
	cat := model.Category{
		ID:   1,
		Name: "test",
	}

	type fields struct {
		data model.Category
		err  error
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Category
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: cat,
				err:  nil,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    cat,
			wantErr: false,
		},
		{
			name: "SelectContext error",
			fields: fields{
				data: model.Category{},
				err:  errors.New("error"),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    model.Category{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			testDB := postgresql.NewBookRepository(sqlxDB)

			query := "SELECT id, name FROM categories WHERE id = $1"

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.id)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "name"}).AddRow(cat.ID, cat.Name)
				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetCategoryByID(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)

		})
	}
}
