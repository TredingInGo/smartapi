package smartstream

import (
	"fmt"
	"github.com/TredingInGo/smartapi"
	"github.com/TredingInGo/smartapi/models"
	"log"
	"testing"
	"time"
)

var (
	client      *WebSocket
	currTime    = time.Now()
	baseTime    = time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 9, 0, 0, 0, time.Local)
	chForCandle = make(chan *models.LTPInfo, 100)
)

const (
	feedToken = "eyJhbGciOiJIUzUxMiJ9.eyJ1c2VybmFtZSI6IlA1MTI4NDc5OSIsImlhdCI6MTY4NjM4OTExNCwiZXhwIjoxNjg2NDc1NTE0fQ.R6RPdP736vbf62WibcVDJYFKkXPH3nKNhaKKGL2ybZ1O1c3sgf99bg62z6g4CH-X-G937u1pgKFnpEN3PLyyJQ"
)

func TestSmartStream(t *testing.T) {
	client = New("P51284799", feedToken)
	client.callbacks.onConnected = onConnected
	client.callbacks.onSnapquote = onSnapquote
	client.callbacks.onLTP = onLTP
	go makeCandle(chForCandle, 60*5)
	client.Connect()

}

func onConnected() {
	log.Printf("connected")
	err := client.Subscribe(models.LTP, []models.TokenInfo{models.TokenInfo{ExchangeType: models.NSECM, Token: "2885"}})
	if err != nil {
		log.Printf("error while subscribing")
	}
}

func onSnapquote(snapquote models.SnapQuote) {
	log.Printf("%d", snapquote.BestFiveSell[0])
}

func onLTP(ltpInfo models.LTPInfo) {
	log.Println(ltpInfo)
	chForCandle <- &ltpInfo
}

//func Test_testOnLtp(t *testing.T) {
//	ch := getDummyLtp()
//	chForCandle := make(chan *models.LTPInfo, 100)
//
//	go makeCandle(chForCandle)
//	for data := range ch {
//		log.Println(data)
//		chForCandle <- &data
//	}
//
//	close(chForCandle)
//	time.Sleep(time.Minute * 50)
//}

func makeCandle(ch <-chan *models.LTPInfo, duration int) {
	//candleDuration := 5
	candles := make([]*smartapigo.CandleResponse, 0)
	//t, err := time.Parse("15:04:05", baseTime)
	//if err != nil {
	//	fmt.Println("Error parsing time:", err)
	//	return
	//}

	//formattedBaseTime := t.Format("15:04:05")
	lastSegStart := time.Time{}

	for data := range ch {
		epochSeconds := int64(data.ExchangeFeedTimeEpochMillis) / 1000
		dataTimeFormatted := time.Unix(epochSeconds, 0)
		if len(candles) == 0 {
			candles = append(candles, &smartapigo.CandleResponse{
				Timestamp: dataTimeFormatted,
				Open:      float64(data.LastTradedPrice) / 100,
				High:      float64(data.LastTradedPrice) / 100,
				Low:       float64(data.LastTradedPrice) / 100,
				Close:     float64(data.LastTradedPrice) / 100,
				Volume:    0,
			})
			tempTime := dataTimeFormatted.Sub(baseTime)
			fmt.Println("temp time", tempTime)
			tempTimeInSec := tempTime.Seconds()
			thresHoldTime := (int(tempTimeInSec)) / (duration)
			thresHoldTime++
			lastSegStart = baseTime.Add(time.Duration(thresHoldTime*duration) * time.Second)

		} else {
			if lastSegStart.After(dataTimeFormatted) {
				lastData := candles[len(candles)-1]
				ltp := float64(data.LastTradedPrice) / 100
				if lastData.Low > ltp {
					lastData.Low = ltp
				}

				if lastData.High < ltp {
					lastData.High = ltp
				}

				lastData.Close = ltp

				candles[len(candles)-1] = lastData
			} else {
				fmt.Println(candles[len(candles)-1])
				candles = append(candles, &smartapigo.CandleResponse{
					Timestamp: dataTimeFormatted,
					Open:      float64(data.LastTradedPrice) / 100,
					High:      float64(data.LastTradedPrice) / 100,
					Low:       float64(data.LastTradedPrice) / 100,
					Close:     float64(data.LastTradedPrice) / 100,
					Volume:    0,
				})

				lastSegStart = lastSegStart.Add(time.Duration(duration) * time.Second)
			}
		}
	}

	for _, data := range candles {
		fmt.Println(data)
	}
}

//func getDummyLtp() <-chan models.LTPInfo {
//	ch := make(chan models.LTPInfo)
//	price := []int64{100, -300, 400, 600}
//	sleepTime := []int{500, 100, 800, 1000}
//	count := 0
//
//	stop := make(chan bool)
//
//	go func() {
//		time.Sleep(time.Second * 30)
//		stop <- true
//	}()
//
//	go func() {
//		for {
//			select {
//			case <-stop:
//				close(ch)
//				return
//			default:
//				ch <- models.LTPInfo{
//					TokenInfo:                   models.TokenInfo{ExchangeType: 1, Token: "2885"},
//					SequenceNumber:              0,
//					ExchangeFeedTimeEpochMillis: uint64(time.Now().UnixMilli()),
//					LastTradedPrice:             uint64(245500 + price[count%4]),
//				}
//
//				time.Sleep(time.Millisecond * time.Duration(sleepTime[count%4]))
//				count++
//
//				if count > 500 {
//					count = 0
//				}
//			}
//		}
//	}()
//
//	return ch
//}
