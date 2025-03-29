package biz

import (
	"context"
	"testing"
)

type MockCookieClient struct {
}

func (m *MockCookieClient) GetCookie(ctx context.Context, stuID string) (string, error) {
	return "JSESSIONID=9810066688E7B6741FEB5D18A82D6488", nil
}

func TestFreeClassroomBiz_crawFreeClassroom(t *testing.T) {
	cli := new(MockCookieClient)
	fcb := &FreeClassroomBiz{
		cookieCli: cli,
	}
	res, err := fcb.crawFreeClassroom(context.Background(), "2024", "2", "testID", 6, 2, []int{1, 2}, "71")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
