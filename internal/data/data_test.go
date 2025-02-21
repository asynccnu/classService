package data

import (
	"context"
	"fmt"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/asynccnu/classService/internal/model"
	"testing"
)

var dt *Data

func TestMain(m *testing.M) {
	cli, err := NewEsClient(&conf.Data{Es: &conf.Data_ES{
		Url:      "http://127.0.0.1:9200",
		Setsniff: false,
		Username: "elastic",
		Password: "12345678",
	}})
	if err != nil {
		panic(fmt.Sprintf("failed to create elasticsearch client: %v", err))
	}
	dt = &Data{cli: cli}
	m.Run()
}

func TestData_AddClassInfo(t *testing.T) {
	t.Run("add a class_info", func(t *testing.T) {
		err := dt.AddClassInfo(context.Background(), model.ClassInfo{
			ID:           "Class:haha:2024:1:1:1-2:cc:somewhere1:65535",
			Day:          1,
			Teacher:      "cc",
			Where:        "somewhere1",
			ClassWhen:    "1-2",
			WeekDuration: "1-16周",
			Classname:    "haha",
			Credit:       1,
			Weeks:        65535,
			Semester:     "1",
			Year:         "2024",
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("add same class_infos", func(t *testing.T) {
		err := dt.AddClassInfo(context.Background(), model.ClassInfo{
			ID:           "Class:haha123:2024:1:1:1-2:cchh:somewhere2:65535",
			Day:          1,
			Teacher:      "cchh",
			Where:        "somewhere2",
			ClassWhen:    "1-2",
			WeekDuration: "1-16周",
			Classname:    "haha123",
			Credit:       1,
			Weeks:        65535,
			Semester:     "1",
			Year:         "2024",
		})
		if err != nil {
			t.Error(err)
		}
		err = dt.AddClassInfo(context.Background(), model.ClassInfo{
			ID:           "Class:haha123:2024:1:1:1-2:cchh:somewhere2:65535",
			Day:          1,
			Teacher:      "cchh",
			Where:        "somewhere2",
			ClassWhen:    "1-2",
			WeekDuration: "1-16周",
			Classname:    "haha123",
			Credit:       1,
			Weeks:        65535,
			Semester:     "1",
			Year:         "2024",
		})
		if err != nil {
			t.Error(err)
		}
	})

}

func TestData_RemoveClassInfo(t *testing.T) {
	err := dt.AddClassInfo(context.Background(), model.ClassInfo{
		ID:           "Class:haha:2023:1:1:1-2:cc:somewhere2:65535",
		Day:          1,
		Teacher:      "cc",
		Where:        "somewhere2",
		ClassWhen:    "1-2",
		WeekDuration: "1-16周",
		Classname:    "haha",
		Credit:       1,
		Weeks:        65535,
		Semester:     "1",
		Year:         "2023",
	})
	if err != nil {
		t.Error(err)
	}
	dt.RemoveClassInfo(context.Background(), "2023", "1")
}

func TestData_SearchClassInfo(t *testing.T) {
	t.Run("search by teacher", func(t *testing.T) {
		res, err := dt.SearchClassInfo(context.Background(), "cc", "2024", "1")
		if err != nil {
			t.Error(err)
		}
		for _, info := range res {
			fmt.Println(info)
		}
	})
	t.Run("search by classname", func(t *testing.T) {
		res, err := dt.SearchClassInfo(context.Background(), "haha", "2024", "1")
		if err != nil {
			t.Error(err)
		}
		for _, info := range res {
			fmt.Println(info)
		}
	})
}
