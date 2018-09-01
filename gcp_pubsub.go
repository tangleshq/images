package images

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yall.in"
)

type GCPPubSub struct {
	topic, id string
	client    *pubsub.Client
}

func NewGCPPubSub(client *pubsub.Client, topic, subID string) *GCPPubSub {
	return &GCPPubSub{
		topic:  topic,
		id:     subID,
		client: client,
	}
}

func (g *GCPPubSub) Listen(ctx context.Context, callback func(ctx context.Context, sha256 string) (Image, error)) error {
	topic, err := g.client.CreateTopic(ctx, g.topic)
	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.AlreadyExists {
			topic = g.client.Topic(g.topic)
		} else {
			return fmt.Errorf("%T: %+v", err, err)
		}
	}
	sub, err := g.client.CreateSubscription(ctx, g.id, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.AlreadyExists {
			sub = g.client.Subscription(g.id)
		} else {
			return fmt.Errorf("%T: %+v", err, err)
		}
	}
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log := yall.FromContext(ctx)
		log = log.WithField("id", msg.ID).WithField("data", string(msg.Data))
		log.Debug("got message")
		img, err := callback(ctx, string(msg.Data))
		if err != nil {
			log.WithError(err).Error("error processing msg")
			msg.Nack()
			return
		}
		msg.Ack()
		log.WithField("image", img.SHA256).Info("image processed")
	})
	if err != nil {
		return err
	}
	return nil
}
