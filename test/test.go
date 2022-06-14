package test

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/dnlo/struct2csv"
)

// create table objective_log (
// 	id int AUTO_INCREMENT,
// 	user_id int,
// 	date DATETIME,
// 	order_n int,
// 	objective_id int,
// 	plan_list_id int,
// 	completed_at DATETIME,
// 	last_updated_at DATETIME,
// 	start_time datetime,
// 	point int,
// 	is_deleted boolean,
// 	PRIMARY KEY(id)
// );

type ObjectiveLog struct {
	UserId        int
	Date          time.Time
	Order         int
	ObjectiveId   int
	PlanListId    int
	CompletedAt   time.Time
	LastUpdatedAt time.Time
	StartTime     time.Time
	Precondition  int
	Point         int
	IsDeleted     bool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const maxCount = 250000

func generateData(index, batch int, wg *sync.WaitGroup) {
	defer wg.Done()
	objectiveLogs := []ObjectiveLog{}
	fmt.Println("Started", index)
	for i := 0; i < maxCount; i++ {
		o := ObjectiveLog{}
		faker.FakeData(&o)
		objectiveLogs = append(objectiveLogs, o)
	}
	filepath := "./out/" + strconv.Itoa(batch) + "_" + strconv.Itoa(index) + ".csv"
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	checkErr(err)
	defer file.Close()
	w := struct2csv.NewWriter(file)
	w.WriteStructs(objectiveLogs)
	checkErr(err)
	fmt.Println("Done", index)
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
