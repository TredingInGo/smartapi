package smartapigo

import (
	"encoding/json"
	"net/http"
	"time"
)

// CandleParams represents parameters for getting CandleData.
type CandleParams struct {
	Exchange    string `json:"exchange"`
	SymbolToken string `json:"symboltoken"`
	Interval    string `json:"interval"`
	FromDate    string `json:"fromdate"`
	ToDate      string `json:"todate"`
}

type CandleResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int       `json:"volume"`
}

func (c *Client) GetCandleData(candleParam CandleParams) ([]CandleResponse, error) {
	var (
		e          [][]interface{}
		candleData *json.RawMessage
	)

	params := structToMap(candleParam, "json")
	err := c.doEnvelope(http.MethodPost, URICandleData, params, nil, &candleData, true)

	// if no data found for given date
	err = json.Unmarshal(*candleData, &e)
	if err != nil {
		return nil, err
	}

	return getData(e), err
}

func getData(e [][]interface{}) []CandleResponse {
	var resp []CandleResponse
	for _, k := range e {
		var entity CandleResponse
		entity.Timestamp, _ = time.Parse(time.RFC3339, k[0].(string))
		entity.Open = k[1].(float64)
		entity.High = k[2].(float64)
		entity.Low = k[3].(float64)
		entity.Close = k[4].(float64)
		entity.Volume = int(k[5].(float64))

		resp = append(resp, entity)
	}

	return resp
}
