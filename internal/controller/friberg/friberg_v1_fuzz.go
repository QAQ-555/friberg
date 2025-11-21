package friberg

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	v1 "friberg/api/friberg/v1"
	imongo "friberg/internal/library/mongo"
)

func (c *ControllerV1) Fuzz(ctx context.Context, req *v1.FuzzReq) (res *v1.FuzzRes, err error) {
	res = &v1.FuzzRes{}
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: "^" + req.Name},
				{Key: "$options", Value: "i"},
			}},
		}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: 6},
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
	cursor, err := imongo.TableSubject.Collection(ctx).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res.GameInfo); err != nil {
		return nil, err
	}
	return res, nil
}
