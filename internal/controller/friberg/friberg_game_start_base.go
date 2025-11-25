package friberg

import (
	"context"

	"friberg/api/friberg/game"
	imongo "friberg/internal/library/mongo"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ControllerGame) StartBase(ctx context.Context, req *game.StartBaseReq) (res *game.StartBaseRes, err error) {
	r := g.RequestFromCtx(ctx)
	matchStage := bson.D{
		{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}},
	}
	projectStage := ProjectStage
	pipeline := mongo.Pipeline{
		matchStage,
		projectStage,
	}

	cursor, err := imongo.TableSubject().Coll().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 直接解码为 subjectinfo
	var result subjectinfo
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		return nil, gerror.New("no document found")
	}
	result.Frequency = 10
	// 存入 session/Redis
	r.Session.Set("subject", result)
	return &game.StartBaseRes{}, nil
}
