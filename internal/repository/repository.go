package repository

import (
	"Epic55/go_currency_app2/internal/models"
	"Epic55/go_currency_app2/internal/service"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "currency"
)

func NewRepository(ConnectionString string) *Repository {
	ConnectionString1 := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", ConnectionString1)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.SetMaxOpenConns(39)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil
	}

	return &Repository{
		Db: db,
	}
}

func (r *Repository) InsertData(rates models.Rates, formattedDate string) {
	savedItemCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			fmt.Printf("Failed to convert float: %s", err)
			continue
		}
		//startTime := time.Now()

		queryStmt := `INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES ($1, $2, $3, $4)`
		rows, err := r.Db.QueryContext(ctx, queryStmt, item.Title, item.Code, value, formattedDate)
		if err != nil {
			fmt.Println("Failed to insert in the db:", err)
		} else {
			savedItemCount++
			fmt.Println("Item saved", "count", savedItemCount)
		}
		defer rows.Close()
		//duration := time.Since(startTime).Seconds()

	}
	fmt.Println("Items saved:", "All", savedItemCount)
}

func (r *Repository) GetData(ctx context.Context, formattedDate, code string) ([]models.DBItem, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = `SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = $1`
		params = []interface{}{formattedDate}
	} else {
		query = `SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = $1 AND CODE = $2`
		params = []interface{}{formattedDate, code}
	}

	//startTime := time.Now()

	rows, err := r.Db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DBItem
	for rows.Next() {
		var item models.DBItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if len(results) == 0 {
		fmt.Println("No data found with these parameters")
	}

	return results, nil
}

func (r *Repository) DeleteData(ctx context.Context, formattedDate, code string) (int64, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = `DELETE FROM R_CURRENCY WHERE A_DATE = $1`
		params = []interface{}{formattedDate}
	} else {
		query = `DELETE FROM R_CURRENCY WHERE A_DATE = $1 AND CODE = $2`
		params = []interface{}{formattedDate, code}
	}

	//startTime := time.Now()

	result, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	//
	if rowsAffected == 0 {
		fmt.Println("No data deleted with these params")
	}
	return rowsAffected, nil
}

func (r *Repository) scheduler(ctx context.Context, formattedDate string, rates models.Rates) error {
	var count int
	err := r.Db.QueryRowContext(ctx, `SELECT COUNT(*) FROM R_CURRENCY WHERE A_DATE = $1`, formattedDate).Scan(&count)
	if err != nil {
		return err
	}

	for _, item := range rates.Items {
		value, errr := strconv.ParseFloat(item.Value, 64)
		if errr != nil {
			fmt.Println("Failed to convert float:", errr)
			continue
		}
		if count > 0 {
			_, err = r.Db.ExecContext(ctx, `UPDATE R_CURRENCY SET TITLE = $1, VALUE = $2, A_DATE = NOW() WHERE A_DATE = $3 AND CODE = $4`, item.Title, value, formattedDate, item.Code)
			if err != nil {
				return err
			}
		} else {
			_, err = r.Db.ExecContext(ctx, `INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES ($1, $2, $3, $4)`, item.Title, item.Code, value, formattedDate)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (r *Repository) HourTick(date, formattedDate string, ctx context.Context, APIURL string) {

	var service = service.NewService()

	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		err := r.scheduler(ctx, formattedDate, *service.GetData1(ctx, date, APIURL))
		if err != nil {
			fmt.Println("Can't update the date:", err)
		}
	}
}
