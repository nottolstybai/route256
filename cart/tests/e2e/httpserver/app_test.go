package httpserver

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteItem(t *testing.T) {
	type AddRequest struct {
		Count uint16 `json:"count"`
	}

	type wantData struct {
		statusCodeOnAdd    int
		statusCodeOnDelete int
	}

	type inputData struct {
		AddRequestURL    string
		AddRequestBody   AddRequest
		DeleteRequestURL string
	}

	client := MakeClient()
	serverObj := MakeServer()

	router := http.NewServeMux()
	router.HandleFunc("POST /user/{user_id}/cart/{sku_id}", serverObj.AddItem)
	router.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", serverObj.DeleteItem)

	testServer := httptest.NewServer(router)

	testCases := []struct {
		name  string
		input inputData
		want  wantData
	}{
		{
			name: "Add valid item and delete item",
			input: inputData{
				AddRequestURL:    testServer.URL + "/user/31337/cart/773297411",
				AddRequestBody:   AddRequest{Count: 10},
				DeleteRequestURL: testServer.URL + "/user/31337/cart/2958025",
			},
			want: wantData{
				statusCodeOnAdd:    http.StatusOK,
				statusCodeOnDelete: http.StatusNoContent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// first send add_item_request
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.input.AddRequestBody)
			if err != nil {
				t.Fatalf("Error encoding request object %v", err)
			}

			addReq, err := http.NewRequest(http.MethodPost, tc.input.AddRequestURL, &buf)
			if err != nil {
				t.Fatalf("Error creating request object %v", err)
			}
			defer addReq.Body.Close()

			resp, err := client.Do(addReq)
			if err != nil {
				t.Fatalf("Error sending request=%s: %v", tc.input.AddRequestURL, err)
			}

			if resp.StatusCode != tc.want.statusCodeOnAdd {
				t.Fatalf("expected status code %d but got %d", tc.want.statusCodeOnAdd, resp.StatusCode)
			}

			//then send delete_item_request
			delReq, err := http.NewRequest(http.MethodDelete, tc.input.DeleteRequestURL, nil)
			if err != nil {
				t.Fatalf("Error creating request object %v", err)
			}

			resp, err = client.Do(delReq)
			if err != nil {
				t.Fatalf("Error sending request=%s: %v", tc.input.AddRequestURL, err)
			}

			if resp.StatusCode != tc.want.statusCodeOnDelete {
				t.Fatalf("expected status code %d but got %d", tc.want.statusCodeOnDelete, resp.StatusCode)
			}
		})
	}
}

func TestListItems(t *testing.T) {
	type (
		AddRequest struct {
			Count uint16 `json:"count"`
		}
		Item struct {
			SkuID int64  `json:"sku_id"`
			Name  string `json:"name"`
			Count uint16 `json:"count"`
			Price uint32 `json:"price"`
		}
		ListItemsResponse struct {
			Items      []Item `json:"items"`
			TotalPrice uint32 `json:"total_price"`
		}
	)

	type wantData struct {
		statusCodeOnAdd  int
		statusCodeOnList int
		response         ListItemsResponse
	}

	type inputData struct {
		AddRequestURL  string
		AddRequestBody AddRequest
		ListRequestURL string
	}

	client := MakeClient()
	serverObj := MakeServer()

	router := http.NewServeMux()
	router.HandleFunc("POST /user/{user_id}/cart/{sku_id}", serverObj.AddItem)
	router.HandleFunc("GET /user/{user_id}/cart", serverObj.ListItems)

	testServer := httptest.NewServer(router)

	testCases := []struct {
		name  string
		input inputData
		want  wantData
	}{
		{
			name: "Add valid item and list item",
			input: inputData{
				AddRequestURL:  testServer.URL + "/user/31337/cart/773297411",
				AddRequestBody: AddRequest{Count: 10},
				ListRequestURL: testServer.URL + "/user/31337/cart",
			},
			want: wantData{
				statusCodeOnAdd:  http.StatusOK,
				statusCodeOnList: http.StatusOK,
				response: ListItemsResponse{
					Items: []Item{
						{
							SkuID: 773297411,
							Count: 10,
							Name:  "Кроссовки Nike JORDAN",
							Price: 2202,
						},
					},
					TotalPrice: 22020,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// first send add_item_request
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.input.AddRequestBody)
			if err != nil {
				t.Fatalf("Error encoding request object %v", err)
			}

			addReq, err := http.NewRequest(http.MethodPost, tc.input.AddRequestURL, &buf)
			if err != nil {
				t.Fatalf("Error creating request object %v", err)
			}
			defer addReq.Body.Close()

			resp, err := client.Do(addReq)
			if err != nil {
				t.Fatalf("Error sending request=%s: %v", tc.input.AddRequestURL, err)
			}

			if resp.StatusCode != tc.want.statusCodeOnAdd {
				t.Fatalf("expected status code %d but got %d", tc.want.statusCodeOnAdd, resp.StatusCode)
			}

			//then send list_items_request
			listReq, err := http.NewRequest(http.MethodGet, tc.input.ListRequestURL, nil)
			if err != nil {
				t.Fatalf("Error creating request object %v", err)
			}

			resp, err = client.Do(listReq)
			if err != nil {
				t.Fatalf("Error sending request=%s: %v", tc.input.AddRequestURL, err)
			}

			if resp.StatusCode != tc.want.statusCodeOnList {
				t.Fatalf("expected status code %d but got %d", tc.want.statusCodeOnList, resp.StatusCode)
			}

			var actualResponse ListItemsResponse
			if err := json.NewDecoder(resp.Body).Decode(&actualResponse); err != nil {
				t.Fatalf("Error decoding response %v", err)
			}
			defer resp.Body.Close()

			require.Equal(t, tc.want.response, actualResponse)
		})
	}
}
