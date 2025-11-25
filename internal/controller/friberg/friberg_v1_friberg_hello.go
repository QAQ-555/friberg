package friberg

import (
	"context"
	"fmt"
	"log"

	"github.com/gogf/gf/v2/frame/g"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	v1 "friberg/api/friberg/v1"
	imongo "friberg/internal/library/mongo"
)

func (c *ControllerV1) FribergHello(ctx context.Context, req *v1.FribergHelloReq) (res *v1.FribergHelloRes, err error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: "^portal"},
				{Key: "$options", Value: "i"},
			}},
		}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: 1},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "name", Value: 1},
		}},
	}

	// Print the aggregation results
	pipeline := mongo.Pipeline{
		matchStage,
		limitStage,
		projectStage,
	}
	cursor, err := imongo.TableSubject().Coll().Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("failed to aggregate: %v", err)
	}
	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("failed to decode results: %v", err)
	}

	// Print the aggregation results.
	for _, result := range results {
		g.DumpWithType(result)
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
		g.DumpWithType(res)
	}
	return nil, nil
}
