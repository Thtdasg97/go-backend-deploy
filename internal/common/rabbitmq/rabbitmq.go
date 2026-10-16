package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-backend/internal/common/env"
	"log"
	"math/rand"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

func NewRabbitMQ(env *env.Env) *RabbitMQ {
	// connect tới rabbitMQ
	conn, err := amqp.DialConfig(env.RabbitMQURL, amqp.Config{
		Properties: amqp.Table{
			"connection_name": "go-backend",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connect RabbitMQ - OK")

	return &RabbitMQ{
		Conn: conn,
	}
}

// tạo correctionID ngẫu nhiên
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func queueDeclareMain(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name - tạo queue với tên giống với event
		true,      // durability - có tồn tại queue hay không khi rmq restart
		false,     // delete when unused -
		false,     // exclusive - tính năng độc quyền
		false,     // noWait có đợi queue tạo thành công hay không?
		nil,       // arguments
	)

}

// queue tạm (queue_temp)

func queueDeclareTemp(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name - tạo queue với tên giống với event
		true,      // durability - có tồn tại queue hay không khi rmq restart
		true,      // delete when unused - queue tạm nên tự động xóa
		false,     // exclusive - tính năng độc quyền
		false,     // noWait có đợi queue tạo thành công hay không?
		nil,       // arguments
	)

}

// /
func publishNotReply(ctx context.Context, ch *amqp.Channel, queueName string, corrId string, body []byte) error {
	return ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json", // trong Golang, nó sẽ hiểu json là byte?
			CorrelationId: corrId,
			DeliveryMode:  amqp.Persistent, // lỡ rapid sập -> vẫn giữ lại message đang nằm trong queue
			// ReplyTo:       q.Name, bắn rồi quên lun thì không cần dòng này
			Body: body,
		})
}

func publishReply(ctx context.Context, ch *amqp.Channel, queueName string, replyTo string, corrId string, body []byte) error {
	return ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json", // trong Golang, nó sẽ hiểu json là byte?
			CorrelationId: corrId,
			DeliveryMode:  amqp.Persistent, // lỡ rapid sập -> vẫn giữ lại message đang nằm trong queue
			ReplyTo:       replyTo,         // bắn rồi quên lun thì không cần dòng này
			Body:          body,
		})
}

type Reply struct {
	ErrorString string
	Data        json.RawMessage
}

// xây dựng func bắn đi trước

func (r *RabbitMQ) Send(ctx context.Context, queueName string, payload any) (err error) {

	// dùng Marshal vì muốn trả byte
	body, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := r.Conn.Channel()

	if err != nil {
		fmt.Println(err)
		return // tự động return biến được đặt ở trên.
	}

	defer ch.Close()

	// tạo cái queue chính
	// nếu đã có queue, nó sẽ kiểm tra đúng các setting hay không, nếu không đúng => báo lỗi, nếu đúng thì đi tiếp.

	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	corrId := randomString(32)

	// bắn vào main queue đã tạo ở trên

	publishNotReply(ctx, ch, queueMain.Name, corrId, body)

	if err != nil {
		fmt.Println(err)
		return // tự động return biến được đặt ở trên.
	}

	return
}

// function lắng nghe
// chạy bằng go routine, để không block main
// khi start BE lên là chạy luôn
func (r *RabbitMQ) On(ctx context.Context, queueName string, handler func(context.Context, []byte) error) (err error) {
	ch, err := r.Conn.Channel()

	if err != nil {
		fmt.Println(err)
		return
	}

	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// consume
	msgs, err := ch.Consume(
		queueMain.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	// sử dụng Go routine để không bị lock

	go func() {
		defer ch.Close() // ý đồ để đây là gì, ông thầy nói mình chả hiểu lắm

		for d := range msgs {

			// chạy handler
			// lưu db lâu quá để timeout vô
			err := handler(ctx, d.Body)
			if err != nil {
				d.Nack(false, false)
				continue
			}

			// lỗi chạy NACK, không lỗi chạy ACK

			d.Ack(false)

		}
	}()

	return
}

// tạo queue phụ và consume
func (r *RabbitMQ) Request(ctx context.Context, queueName string, payload any, result any) (err error) {

	// dùng Marshal vì muốn trả byte
	body, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := r.Conn.Channel()

	if err != nil {
		fmt.Println(err)
		return // tự động return biến được đặt ở trên.
	}

	// ===============Block scope - temp====================
	queueTemp, err := queueDeclareTemp(ch, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	// consume
	msgs, err := ch.Consume(
		queueTemp.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	// ===============Block scope====================

	// ===============Block scope - main====================
	// tạo cái queue chính
	// nếu đã có queue, nó sẽ kiểm tra đúng các setting hay không, nếu không đúng => báo lỗi, nếu đúng thì đi tiếp.

	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	corrId := randomString(32)

	// bắn vào main queue đã tạo ở trên

	err = publishReply(ctx, ch, queueMain.Name, queueTemp.Name, corrId, body)

	if err != nil {
		fmt.Println(err)
		return // tự động return biến được đặt ở trên.
	}
	// ====================================

	// setting timeout

	for {
		select {
		// done nghĩa là nó đang bị timeout -> tạo lỗi và return
		case <-ctx.Done():
			err = errors.New("Request timeout cmnr -_- ")
			fmt.Println(err)
			return

		// Trường hợp lấy được dữ liệu
		case d := <-msgs:
			if d.CorrelationId != corrId {
				continue
			}

			var bodyReply Reply
			err = json.Unmarshal(d.Body, &bodyReply)
			if err != nil {
				fmt.Println(err)
				return
			}

			if bodyReply.ErrorString != "" {

				err = errors.New(bodyReply.ErrorString)
				fmt.Println(err)
				return

			}

			err = json.Unmarshal(bodyReply.Data, result)
			if err != nil {
				fmt.Println(err)
				return
			}
			return

		}
	}
}

func (r *RabbitMQ) OnReply(ctx context.Context, queueName string, handler func(context.Context, []byte) (any, error)) (err error) {
	ch, err := r.Conn.Channel()

	if err != nil {
		fmt.Println(err)
		return
	}

	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// consume
	msgs, err := ch.Consume(
		queueMain.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	// sử dụng Go routine để không bị lock

	go func() {
		defer ch.Close() // ý đồ để đây là gì, ông thầy nói mình chả hiểu lắm

		for d := range msgs {

			// tạo hàm trả lỗi cho các logic ở dưới
			// tạo hàm anomymous -> gắn vô biến

			replyError := func(replyError error) {
				d.Nack(false, false)

				bodyReply := Reply{
					ErrorString: replyError.Error(),
					Data:        nil,
				}

				body, err := json.Marshal(bodyReply)

				if err != nil {
					fmt.Println(err)
					return
				}

				err = publishNotReply(ctx, ch, d.ReplyTo, d.CorrelationId, body)
				if err != nil {
					fmt.Println(err)
					return
				}

			}

			// chạy handler
			// lưu db lâu quá để timeout vô

			result, err := handler(ctx, d.Body)
			// xử lý lỗi cho error
			if err != nil {
				replyError(err)
				continue
			}

			// xử lý data
			data, err := json.Marshal(result)
			if err != nil {
				replyError(err)
				continue
			}

			// xử lý format
			bodyReply := Reply{
				ErrorString: "",
				Data:        data,
			}

			body, err := json.Marshal(bodyReply)
			if err != nil {
				replyError(err)
				continue
			}

			err = publishNotReply(ctx, ch, d.ReplyTo, d.CorrelationId, body)
			if err != nil {
				d.Nack(false, false)
				continue
			}

			// lỗi chạy NACK, không lỗi chạy ACK

			d.Ack(false)

		}
	}()

	return
}
