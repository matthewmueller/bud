package sqs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/matthewmueller/bud/job"
	"github.com/matthewmueller/bud/log"
)

// New SQS queue
func New(client sqsiface.SQSAPI, log log.Log, queueUrl string) *Queue {
	return &Queue{client, log, queueUrl}
}

type Queue struct {
	client   sqsiface.SQSAPI
	log      log.Log
	queueUrl string
}

var _ job.Queue = (*Queue)(nil)

// Push multiple messages to SQS
func (q *Queue) Push(ctx context.Context, messages ...*job.Message) error {
	// TODO: parallelize sending batches
	batches := chunkJobs(messages, 10)
	for _, batch := range batches {
		if err := q.pushBatch(ctx, batch); err != nil {
			return err
		}
	}
	return nil
}

func (q *Queue) pushBatch(ctx context.Context, batch []*job.Message) error {
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(batch))
	for i, job := range batch {
		entries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(job.ID),
			MessageBody: aws.String(string(job.Payload)),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"SentTimestamp": {
					DataType:    aws.String("Number"),
					StringValue: aws.String(strconv.FormatInt(job.Timestamp.Unix(), 10)),
				},
			},
		}
	}
	result, err := q.client.SendMessageBatchWithContext(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(q.queueUrl),
		Entries:  entries,
	})
	if err != nil {
		return fmt.Errorf("sqs: error pushing jobs: %w", err)
	} else if len(result.Failed) > 0 {
		return fmt.Errorf("sqs: error pushing jobs: %v", result.Failed)
	}
	return nil
}

func chunkJobs(messages []*job.Message, chunkSize int) [][]*job.Message {
	var chunks [][]*job.Message
	size := len(messages)
	for i := 0; i < size; i += chunkSize {
		end := i + chunkSize
		if end > size {
			end = size
		}
		chunks = append(chunks, messages[i:end])
	}
	return chunks
}

// // Work the queue
// func (q *Queue) Work(ctx context.Context) error {
// 	err := job.Loop(ctx, q.log, q.pull)
// 	if err != nil {
// 		return fmt.Errorf("sqs: error working queue: %w", err)
// 	}
// 	return nil
// }

func (q *Queue) Size(ctx context.Context) (int64, error) {
	out, err := q.client.GetQueueAttributesWithContext(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(q.queueUrl),
		AttributeNames: []*string{aws.String("ApproximateNumberOfMessages")},
	})
	if err != nil {
		return 0, fmt.Errorf("sqs: error getting size: %w", err)
	}
	count := out.Attributes["ApproximateNumberOfMessages"]
	return strconv.ParseInt(aws.StringValue(count), 10, 64)
}

// Pull a job from the queue and process it
func (q *Queue) Pull(ctx context.Context, handle func(ctx context.Context, msg *job.Message) error) error {
	msg, err := q.popMessage(ctx)
	if err != nil {
		return err
	}
	timeStamp, err := getTimestamp(msg)
	if err != nil {
		return err
	}
	message := &job.Message{
		ID:        *msg.MessageId,
		Payload:   []byte(*msg.Body),
		Timestamp: timeStamp,
	}
	if err := handle(ctx, message); err != nil {
		return fmt.Errorf("sqs: unable to handle message %s: %w", message.ID, err)
	}
	return q.deleteMessage(ctx, msg)
}

func (q *Queue) Worker(handler job.Handler) job.Worker {
	return job.Looper(q, handler)
}

// Pop a message off the queue
func (q *Queue) popMessage(ctx context.Context) (*sqs.Message, error) {
	out, err := q.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(q.queueUrl),
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(20),
		MessageAttributeNames: []*string{aws.String("All")},
	})
	if err != nil {
		return nil, fmt.Errorf("sqs: error receiving message: %w", err)
	}
	if len(out.Messages) == 0 {
		return nil, job.ErrEmptyQueue
	}
	msg := out.Messages[0]
	return msg, nil
}

// Delete a message from the queue
func (q *Queue) deleteMessage(ctx context.Context, msg *sqs.Message) error {
	_, err := q.client.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueUrl),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("sqs: unable to delete message %q. %w", *msg.MessageId, err)
	}
	return nil
}

func getTimestamp(msg *sqs.Message) (t time.Time, err error) {
	val, ok := msg.MessageAttributes["SentTimestamp"]
	if !ok {
		return t, nil
	} else if val.DataType == nil || *val.DataType != "Number" || val.StringValue == nil {
		return t, nil
	}
	n, err := strconv.Atoi(*val.StringValue)
	if err != nil {
		return t, fmt.Errorf("sqs: error parsing timestamp: %w", err)
	}
	t = time.Unix(int64(n), 0)
	return t, nil
}
