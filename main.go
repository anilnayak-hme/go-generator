package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ObjectiveLog struct {
	UserId        int       `db:"user_id"`
	Date          time.Time `db:"date"`
	Order         int       `db:"order"`
	ObjectiveId   int       `db:"objective_id"`
	PlanListId    int       `db:"plan_list_id"`
	CompletedAt   time.Time `db:"completed_at"`
	LastUpdatedAt time.Time `db:"last_updated_at"`
	StartTime     time.Time `db:"start_time"`
	Precondition  int       `db:"precondition"`
	Point         int       `db:"point"`
	IsDeleted     bool      `db:"is_deleted"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const maxCount = 2000

func generateData(index, batch int, wg *sync.WaitGroup) {
	defer wg.Done()
	objectiveLogs := []ObjectiveLog{}
	fmt.Println("Started", index)
	db, err := sqlx.Connect("mysql", "root:evolution@(localhost:3306)/navi")
	defer db.Close()
	checkErr(err)
	for i := 0; i < maxCount; i++ {
		o := ObjectiveLog{}
		faker.FakeData(&o)
		objectiveLogs = append(objectiveLogs, o)
	}
	res, err := db.NamedExec(`INSERT INTO objective_log 
		(user_id, date, order_n, objective_id, plan_list_id, completed_at, last_updated_at, start_time, point, is_deleted)
		Values (:user_id, :date, :order, :objective_id, :plan_list_id, :completed_at, :last_updated_at, :start_time, :point, :is_deleted)`, objectiveLogs,
	)
	checkErr(err)
	fmt.Println(res)
}

func main() {
	total := 100000000
	threads := 8
	batches := total / maxCount / threads
	fmt.Println("total batches", batches)
	for j := 0; j < batches; j++ {
		fmt.Println("starting batch", j)
		wg := sync.WaitGroup{}
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go generateData(i, j, &wg)
		}
		wg.Wait()
	}
}
