package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"linebot/gurunavi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Webhook struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// botは*Client型
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	// エラー処理(envがうまく入ったか)
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			// サーバー側のエラーを返す
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message:":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
		}, nil
	}
	// requestのログ
	log.Print(request.Headers)
	log.Print(request.Body)

	if !validateSignature(os.Getenv("LINE_CHANNEL_SECRET"), request.Headers["X-Line-Signature"], []byte(request.Body)) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", linebot.ErrInvalidSignature.Error()),
		}, nil
	}

	// リクエストのボディ部を格納するための構造体を定義
	var webhook Webhook

	// jsonを構造体に変換 &webhookにjsonが変換されて入る
	if err := json.Unmarshal([]byte(request.Body), &webhook); err != nil {
		log.Print(err)
		// クライアントのエラーを返す
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
		}, nil
	}

	// webhookで集めたイベントの配列のポインタをここで処理する
	for _, event := range webhook.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			// そのインターフェースを実装する型によって処理を分岐させることができます。
			switch m := event.Message.(type) {
			case *linebot.TextMessage:
				switch request.Path {
				case "/parrot":
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(m.Text)).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
						}, nil
					}
				case "/restaurants":
					g, err := gurunavi.SearchRestaurants(m.Text)
					if err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
						}, nil
					}
					var sm linebot.SendingMessage
					switch {
					case g.Error != nil:
						t := g.Error[0].Message
						sm = linebot.NewTextMessage(t)
					default:
						f := FlexRestaurants(g)
						sm = linebot.NewFlexMessage("飲食店検索結果", f)
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, sm).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
						}, err
					}
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func validateSignature(channelSecret string, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	hash := hmac.New(sha256.New, []byte(channelSecret))
	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}

func main() {
	lambda.Start(handler)
}
