package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	cfgpkg "github.com/phishguard/phishguard/internal/config"
)

// DeliverMessage is the payload the API publishes and the worker consumes.
type DeliverMessage struct {
	RecipientID string `json:"recipient_id"`
	CampaignID  string `json:"campaign_id"`
}

// Validate rejects empty or whitespace-only identifiers (poison to the worker).
func (m DeliverMessage) Validate() error {
	if strings.TrimSpace(m.RecipientID) == "" {
		return fmt.Errorf("recipient_id is required")
	}
	if strings.TrimSpace(m.CampaignID) == "" {
		return fmt.Errorf("campaign_id is required")
	}
	return nil
}

type Client struct {
	sqs      *sqs.Client
	queueURL string
}

func New(ctx context.Context, cfg cfgpkg.Config) (*Client, error) {
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKey, cfg.AWSSecretKey, "",
		)),
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	var sqsOpts []func(*sqs.Options)
	endpoint := cfg.SQSEndpoint
	if endpoint == "" {
		endpoint = cfg.AWSEndpointURL
	}
	if endpoint != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	return &Client{
		sqs:      sqs.NewFromConfig(awsCfg, sqsOpts...),
		queueURL: cfg.SQSQueueURL,
	}, nil
}

func (c *Client) PublishDeliver(ctx context.Context, msg DeliverMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (c *Client) Receive(ctx context.Context, max int32, waitSeconds int32) ([]Received, error) {
	out, err := c.sqs.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: max,
		WaitTimeSeconds:     waitSeconds,
		VisibilityTimeout:   30,
	})
	if err != nil {
		return nil, err
	}

	var msgs []Received
	for _, m := range out.Messages {
		var payload DeliverMessage
		if err := json.Unmarshal([]byte(aws.ToString(m.Body)), &payload); err != nil {
			// Delete poison messages that cannot be parsed so they do not block the queue.
			_ = c.Delete(ctx, aws.ToString(m.ReceiptHandle))
			continue
		}
		msgs = append(msgs, Received{
			Payload:       payload,
			ReceiptHandle: aws.ToString(m.ReceiptHandle),
		})
	}
	return msgs, nil
}

func (c *Client) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

type Received struct {
	Payload       DeliverMessage
	ReceiptHandle string
}

func ParseBody(body string) (DeliverMessage, error) {
	var msg DeliverMessage
	err := json.Unmarshal([]byte(body), &msg)
	return msg, err
}
