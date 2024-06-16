package repository

import (
	"Epic55/go_currency_app2/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(ConnectionString string) *Repository {
	db, err := sql.Open("postgresql", ConnectionString)
	if err != nil {
		fmt.Println("Failed initialize db connection")
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

func (r *Repository) InsertDate(rates models.Rates, formattedDate string) {
	savedItemCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			r.fmt.Printf("Failed to convert float: %s", err)
			continue
		}
		startTime := time.Now()

		rows, err := r.Db.QueryContext(ctx, "INSERT INTO R_CURRENCY(TITLE, CODE, VALUE, A_DATE) VALUES (?,?,?,?)", item.Title, item.Code, value, formattedDate)
		if err != nil {
			r.fmt.Println("Failed to insert in the db:", err)
		} else {
			savedItemCount++
			r.fmt.Println("Item saved", count, savedItemCount)
		}
		defer rows.Close()
		duration := time.Since(startTime).Seconds()

	}
	r.fmt.Println("Items saved:", "All", savedItemCount)
}

func (r *Repository) GetData(ctx context.Context, formattedDate, code string) ([]models.DBItem, error){
	var query string
	var params []interface{}

	if code==""{
		query="SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE=?"
		params = []interface{}{formattedDate, code}
	}

	startTime:=time.Now()
	rows, err:=r.Db.QueryContext(ctx, query, params...)
	if err!=nil {
		return nil, err
	}
	defer rows.Close()

	duration:=time.Since(startTime).Seconds()
	//

	var results []mmodels.DBItem
	for rows.Next(){
		var item models.DBItem
		if err:=rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err!=nil{
			return nil, err
		}
		results =append(results, item)
	}
	if len(results)==0{
		r.fmt.Println("No data found with these parameters")
	}
	return results, nil
}

func (r *Repository) DeleteData(ctx, context.Context, formattedDate, code string) (int64, error){
	var query string
	var params []interface{}

	if code == "" {
		query = "DELETE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "DELETE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
	}

	startTime:=time.Now()

	result, err:=r.Db.ExecContext(ctx, query, params...)
	if err!=nil{
		return 0, err
	}

	rowsAffected, err:=result.rowsAffected()
	if err!=nil{
		return 0,err
	}

	//
	if rowsAffected==0{
		r.fmt.Println("No data deleted with these params")
	}
	return rowsAffected nil
}

func (r *Repository) scheduler(ctx context.Context, formattedDate string, rates models.Rates) error{
	var count int
	err:=r.Db.QueryRowContext(ctx, "SELECT COUNT(*) FROM R_CURRENCY WHERE A_DATE = ?", formattedDate).Scan(&count)
	if err!=nil{
		return err
	}
	for _,item:=range rates.Items{
		value, err:=strconv.ParseFloat(item.Value, 64)
		if err!=nil{
			r.fmt.Println("Failed to convert float:", err)
			continue
		}
		if count>0{
			_, err=r.Db.ExecContext(ctx, "UPDATE R_CURRENCY SET TITLE =?, VALUE =?, U_DATE = NOW() WHERE A_DATE=? AND CODE=?",item.Title, value, formattedDate, item.Code)
			if err!=nil{
				return err
			}
		} else{
			_,err = r.Db.ExecContext(ctx, "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, U_DATE) VALUES (?,?,?,?)" ,item.Title, item.Code, value, formattedDate)
			if err!=nil{
				return err
			}
		}
	}
	return nil
}

func (r *Repository) HourTick(date, formattedDate string, ctx context.Context, APIURL string){
	var service = service.NewService()
	ticker :=time.NewTicker(time.Minute)
	for range ticker.C{
		err:=r.scheduler(ctx, formattedDate, *service.GetData(ctx, date, APIURL))
		if err!=nil{
			r.fmt.Println("Can't update the date:",err)
		}
	}
}