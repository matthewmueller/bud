package sqs_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/livebud/buddy/job"
	"github.com/livebud/buddy/job/sqs"
	"github.com/livebud/buddy/log"
	"github.com/matryer/is"
)

func TestReal(t *testing.T) {
	if os.Getenv("AWS_PROFILE") == "" {
		t.Skip("AWS_PROFILE not set")
	} else if os.Getenv("AWS_REGION") == "" {
		t.Skip("AWS_REGION not set")
	} else if os.Getenv("QUEUE_URL") == "" {
		t.Skip("QUEUE_URL not set")
	}
	is := is.New(t)
	ctx := context.Background()
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	log := log.Default()
	sqs := sqs.New(awssqs.New(session), log, os.Getenv("QUEUE_URL"))
	err := sqs.Push(ctx, &job.Message{
		ID:        "1",
		Payload:   []byte(`{"url":"https://example.com"}`),
		Timestamp: time.Now(),
	})
	is.NoErr(err)
	err = sqs.Push(ctx, &job.Message{
		ID:        "2",
		Payload:   []byte(`{"url":"https://example2.com"}`),
		Timestamp: time.Now(),
	})
	is.NoErr(err)

	size, err := sqs.Size(ctx)
	is.NoErr(err)
	fmt.Println("queue size", size)
	err = sqs.Pull(ctx, func(ctx context.Context, msg *job.Message) error {
		fmt.Println("processing job:", msg.ID)
		// is.Equal(job.ID, "1")
		// is.Equal(string(job.Payload), `{"url":"https://example.com"}`)
		fmt.Println(msg.Timestamp.Format(time.Kitchen))
		// is.Equal(job.Timestamp, int64(1682305589))
		// time.Sleep(1 * time.Second)
		return nil
	})
	is.NoErr(err)
	err = sqs.Pull(ctx, func(ctx context.Context, msg *job.Message) error {
		fmt.Println("processing job:", msg.ID)
		// is.Equal(job.ID, "1")
		// is.Equal(string(job.Payload), `{"url":"https://example.com"}`)
		fmt.Println(msg.Timestamp.Format(time.Kitchen))
		// is.Equal(job.Timestamp, int64(1682305589))
		// time.Sleep(1 * time.Second)
		return nil
	})
	is.NoErr(err)
	// err = client.Dequeue(ctx, &job)
	// is.NoErr(err)
	// is.Equal(job.ID, 2)
	// is.Equal(job.URL, "https://example2.com")
	// is.Equal(job.Timestamp, int64(1682305590))
	// time.Sleep(1 * time.Second)
	// err = client.Dequeue(ctx, &job)
	// is.True(err != nil)
	// is.True(errors.Is(err, job.ErrEmptyQueue))
}

type mockClient struct {
	sqsiface.SQSAPI
	MockReceiveMessageWithContext     func(ctx aws.Context, in *awssqs.ReceiveMessageInput, options ...request.Option) (*awssqs.ReceiveMessageOutput, error)
	MockSendMessageBatchWithContext   func(ctx aws.Context, in *awssqs.SendMessageBatchInput, options ...request.Option) (*awssqs.SendMessageBatchOutput, error)
	MockDeleteMessageWithContext      func(ctx aws.Context, in *awssqs.DeleteMessageInput, options ...request.Option) (*awssqs.DeleteMessageOutput, error)
	MockGetQueueAttributesWithContext func(ctx aws.Context, in *awssqs.GetQueueAttributesInput, options ...request.Option) (*awssqs.GetQueueAttributesOutput, error)
}

func (m *mockClient) ReceiveMessageWithContext(ctx aws.Context, in *awssqs.ReceiveMessageInput, options ...request.Option) (*awssqs.ReceiveMessageOutput, error) {
	return m.MockReceiveMessageWithContext(ctx, in, options...)
}

func (m *mockClient) SendMessageBatchWithContext(ctx aws.Context, in *awssqs.SendMessageBatchInput, options ...request.Option) (*awssqs.SendMessageBatchOutput, error) {
	return m.MockSendMessageBatchWithContext(ctx, in, options...)
}

func (m *mockClient) DeleteMessageWithContext(ctx aws.Context, in *awssqs.DeleteMessageInput, options ...request.Option) (*awssqs.DeleteMessageOutput, error) {
	return m.MockDeleteMessageWithContext(ctx, in, options...)
}

func (m *mockClient) GetQueueAttributesWithContext(ctx aws.Context, in *awssqs.GetQueueAttributesInput, options ...request.Option) (*awssqs.GetQueueAttributesOutput, error) {
	return m.MockGetQueueAttributesWithContext(ctx, in, options...)
}

func TestProcessor(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	url := "https://sqs.us-west-2.amazonaws.com/036813706318/test-queue"
	buf := []*awssqs.SendMessageBatchRequestEntry{}
	sendCount := 0
	recvCount := 0
	mock := &mockClient{
		MockSendMessageBatchWithContext: func(_ aws.Context, in *awssqs.SendMessageBatchInput, _ ...request.Option) (*awssqs.SendMessageBatchOutput, error) {
			sendCount++
			is.True(in.QueueUrl != nil)
			is.Equal(*in.QueueUrl, url)
			is.Equal(len(in.Entries), 1)
			is.True(in.Entries[0].MessageBody != nil)
			buf = append(buf, in.Entries[0])
			return &awssqs.SendMessageBatchOutput{}, nil
		},
		MockReceiveMessageWithContext: func(_ aws.Context, in *awssqs.ReceiveMessageInput, _ ...request.Option) (*awssqs.ReceiveMessageOutput, error) {
			recvCount++
			is.True(in.QueueUrl != nil)
			is.Equal(*in.QueueUrl, url)
			is.Equal(in.MaxNumberOfMessages, aws.Int64(1))
			if len(buf) == 0 {
				return &awssqs.ReceiveMessageOutput{}, nil
			}
			msg := buf[0]
			return &awssqs.ReceiveMessageOutput{
				Messages: []*awssqs.Message{
					&awssqs.Message{
						MessageId:         msg.Id,
						Body:              msg.MessageBody,
						ReceiptHandle:     msg.Id,
						MessageAttributes: msg.MessageAttributes,
					},
				},
			}, nil
		},
		MockDeleteMessageWithContext: func(_ aws.Context, in *awssqs.DeleteMessageInput, _ ...request.Option) (*awssqs.DeleteMessageOutput, error) {
			is.True(in.QueueUrl != nil)
			is.Equal(*in.QueueUrl, url)
			is.True(in.ReceiptHandle != nil)
			for i, msg := range buf {
				if *in.ReceiptHandle == *msg.Id {
					buf = append(buf[:i], buf[i+1:]...)
					break
				}
			}
			return &awssqs.DeleteMessageOutput{}, nil
		},
	}
	queue := sqs.New(mock, log.Default(), url)
	err := queue.Push(ctx, &job.Message{
		ID:        "1",
		Payload:   []byte(`{"url":"https://example.com"}`),
		Timestamp: time.Now(),
	})
	is.NoErr(err)
	err = queue.Push(ctx, &job.Message{
		ID:        "2",
		Payload:   []byte(`{"url":"https://example2.com"}`),
		Timestamp: time.Now(),
	})
	is.NoErr(err)
	err = job.Work(ctx, queue.Worker(&Handler{log.Default(), 0}))
	is.NoErr(err)
	// worker.New(log.Default(), queue)
	// err := sqsQueue.Push(ctx, &queue.Message{
	// 	ID:        "2",
	// 	Payload:   []byte(`{"url":"https://example2.com"}`),
	// 	Timestamp: time.Now(),
	// })
	// is.NoErr(err)
	// processor := &Processor{
	// 	Queue: sqsQueue,
	// }

	// err := client.Enqueue(ctx, Job{
	// 	ID:        1,
	// 	URL:       "https://example.com",
	// 	Timestamp: 1682305589,
	// })
	// is.NoErr(err)
	// err = client.Enqueue(ctx, Job{
	// 	ID:        2,
	// 	URL:       "https://example2.com",
	// 	Timestamp: 1682305590,
	// })
	// is.NoErr(err)
	// var job Job
	// err = client.Dequeue(ctx, &job)
	// is.NoErr(err)
	// is.Equal(job.ID, 1)
	// is.Equal(job.URL, "https://example.com")
	// is.Equal(job.Timestamp, int64(1682305589))
	// err = client.Dequeue(ctx, &job)
	// is.NoErr(err)
	// is.Equal(job.ID, 2)
	// is.Equal(job.URL, "https://example2.com")
	// is.Equal(job.Timestamp, int64(1682305590))
	// is.Equal(sendCount, 2)
	// is.Equal(recvCount, 2)
	// err = client.Dequeue(ctx, &job)
	// is.True(err != nil)
	// is.True(errors.Is(err, queue.ErrEmptyQueue))

}

// // GHtoDBSync processor
// type Handler struct {
// 	Queue queue.Queue
// }

// var _ queue.Handler = (*Handler)(nil)

// func (p *Handler) Push(ctx context.Context, id string, job *Job) error {
// 	payload, err := json.Marshal(job)
// 	if err != nil {
// 		return err
// 	}
// 	return p.Queue.Push(ctx, &queue.Message{
// 		ID:        id,
// 		Payload:   payload,
// 		Timestamp: time.Now(),
// 	})
// }

// func Work(ctx context.Context, queue queue.Queue, fn func(ctx context.Context, msg *queue.Message) error) error {
// 	return queue.Pull(context.Background(), fn)
// }

// // Handle the next job
// func (p *Handler) Handle(ctx context.Context) error {
// 	return p.Queue.Pull(ctx, p.handle)

// }

// func New(log log.Log, queue job.Queue) *Handler {
// 	return &Handler{
// 		log:   log,
// 		queue: queue,
// 	}
// }

type Handler struct {
	log      log.Log
	attempts int
}

var _ job.Handler = (*Handler)(nil)

type Job struct {
	URL string
}

func (h *Handler) Handle(ctx context.Context, msg *job.Message) error {
	var job Job
	if err := json.Unmarshal(msg.Payload, &job); err != nil {
		return err
	}
	h.log.Infof("gh-to-db: handling job %s with url %s", msg.ID, job.URL)
	if h.attempts < 2 {
		h.attempts++
		panic(errors.New("oh noz!"))
	}
	// return fmt.Errorf("unable to handle message")
	// if err := h.handle(ctx, msg); err != nil {
	// 	h.log.Error("gh-to-db: error handling %s: %w", msg.ID, err)
	// }
	return nil
}

func (h *Handler) HandleError(ctx context.Context, err error) error {
	if errors.Is(err, job.ErrEmptyQueue) {
		return err
	}
	var je job.Error
	if errors.As(err, &je) {
		h.log.Errorf("gh-to-db: error handling job %s. %s", je.ID, je)
	}
	return nil
}
