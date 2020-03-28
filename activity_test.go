package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestProjectActivityService_List(t *testing.T) {
	projectIDOrKey := "TEST"
	want := struct {
		spath string
	}{
		spath: "projects/" + projectIDOrKey + "/activities",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	}
	pas := backlog.ExportNewProjectActivityService(cm)
	pas.List(projectIDOrKey)
}

func TestProjectActivityService_List_projectIDOrKeyIsEmpty(t *testing.T) {
	projectIDOrKey := ""
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			t.Error("clientMethod.Get must never be called")
			return nil, errors.New("error")
		},
	}
	pas := backlog.ExportNewProjectActivityService(cm)
	pas.List(projectIDOrKey)
}

func TestSpaceActivityService_List(t *testing.T) {
	want := struct {
		spath string
	}{
		spath: "space/activities",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	}
	sas := backlog.ExportNewSpaceActivityService(cm)
	sas.List()
}

func TestUserActivityService_List(t *testing.T) {
	id := 1234
	want := struct {
		spath string
	}{
		spath: "users/" + strconv.Itoa(id) + "/activities",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			return nil, errors.New("error")
		},
	}
	uas := backlog.ExportNewUserActivityService(cm)
	uas.List(id)
}

func TestUserActivityService_List_invaliedID(t *testing.T) {
	id := 0
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			t.Error("clientMethod.Get must never be called")
			return nil, errors.New("error")
		},
	}
	uas := backlog.ExportNewUserActivityService(cm)
	uas.List(id)
}

func TestBaseActivityService_GetList(t *testing.T) {
	aos := &backlog.ActivityOptionService{}
	type want struct {
		activityTypeId []string
		minId          []string
		maxId          []string
		count          []string
		order          []string
	}
	cases := map[string]struct {
		options   []backlog.ActivityOption
		wantError bool
		want      want
	}{
		"NoOption": {
			options:   []backlog.ActivityOption{},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithActivityTypeIDs_[0]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{0}),
			},
			wantError: true,
			want:      want{},
		},
		"WithActivityTypeIDs_[1]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{1}),
			},
			wantError: false,
			want: want{
				activityTypeId: []string{"1"},
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithActivityTypeIDs_[26]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{26}),
			},
			wantError: false,
			want: want{
				activityTypeId: []string{"26"},
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithActivityTypeIDs_[27]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{27}),
			},
			wantError: true,
			want:      want{},
		},
		"WithActivityTypeIDs_[1...26]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
					14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
				}),
			},
			wantError: false,
			want: want{
				activityTypeId: []string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13",
					"14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26",
				},
				minId: nil,
				maxId: nil,
				count: nil,
				order: nil,
			},
		},
		"WithActivityTypeIDs_[0,1]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{0, 1}),
			},
			wantError: true,
			want:      want{},
		},
		"WithActivityTypeIDs_empty": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{}),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithActivityTypeIDs_[1,1]": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{1, 1}),
			},
			wantError: false,
			want: want{
				activityTypeId: []string{"1", "1"},
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithMinID_0": {
			options: []backlog.ActivityOption{
				aos.WithMinID(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithMinID_1": {
			options: []backlog.ActivityOption{
				aos.WithMinID(1),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          []string{"1"},
				maxId:          nil,
				count:          nil,
				order:          nil,
			},
		},
		"WithMaxID_0": {
			options: []backlog.ActivityOption{
				aos.WithMaxID(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithMaxID_1": {
			options: []backlog.ActivityOption{
				aos.WithMaxID(1),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          []string{"1"},
				count:          nil,
				order:          nil,
			},
		},
		"WithCount_0": {
			options: []backlog.ActivityOption{
				aos.WithCount(0),
			},
			wantError: true,
			want:      want{},
		},
		"WithCount_1": {
			options: []backlog.ActivityOption{
				aos.WithCount(1),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          []string{"1"},
				order:          nil,
			},
		},
		"WithCount_100": {
			options: []backlog.ActivityOption{
				aos.WithCount(100),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          []string{"100"},
				order:          nil,
			},
		},
		"WithCount_101": {
			options: []backlog.ActivityOption{
				aos.WithCount(101),
			},
			wantError: true,
			want:      want{},
		},
		"WithOrder_asc": {
			options: []backlog.ActivityOption{
				aos.WithOrder(backlog.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          []string{backlog.OrderAsc},
			},
		},
		"WithOrder_desc": {
			options: []backlog.ActivityOption{
				aos.WithOrder(backlog.OrderDesc),
			},
			wantError: false,
			want: want{
				activityTypeId: nil,
				minId:          nil,
				maxId:          nil,
				count:          nil,
				order:          []string{backlog.OrderDesc},
			},
		},
		"WithOrder_empty": {
			options: []backlog.ActivityOption{
				aos.WithOrder(""),
			},
			wantError: true,
			want:      want{},
		},
		"WithOrder_invalied": {
			options: []backlog.ActivityOption{
				aos.WithOrder("test"),
			},
			wantError: true,
			want:      want{},
		},
		"MultipleOptions": {
			options: []backlog.ActivityOption{
				aos.WithActivityTypeIDs([]int{1, 2}),
				aos.WithMinID(1),
				aos.WithMaxID(100),
				aos.WithCount(20),
				aos.WithOrder(backlog.OrderAsc),
			},
			wantError: false,
			want: want{
				activityTypeId: []string{"1", "2"},
				minId:          []string{"1"},
				maxId:          []string{"100"},
				count:          []string{"20"},
				order:          []string{backlog.OrderAsc},
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/activity.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					// Check options.
					v := *params.ExportURLValues()
					assert.Equal(t, tc.want.activityTypeId, v["activityTypeId[]"])
					assert.Equal(t, tc.want.minId, v["minId"])
					assert.Equal(t, tc.want.maxId, v["maxId"])
					assert.Equal(t, tc.want.count, v["count"])
					assert.Equal(t, tc.want.order, v["order"])

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			}
			ps := backlog.ExportNewSpaceActivityService(cm)

			if _, err := ps.List(tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}