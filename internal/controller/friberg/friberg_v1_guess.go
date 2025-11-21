package friberg

import (
	"context"
	"log"

	v1 "friberg/api/friberg/v1"
	imongo "friberg/internal/library/mongo"
	"friberg/internal/model/iostruct"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ControllerV1) Guess(ctx context.Context, req *v1.GuessReq) (res *v1.GuessRes, err error) {
	g.Log().Debug(ctx, "guess req:", req)
	r := g.RequestFromCtx(ctx)
	sessionData, err := r.Session.Data()
	if err != nil {
		return nil, err
	}
	answer, err := parseSubjectInfoFromRedis(sessionData["subject"])
	if err != nil {
		return nil, err
	}
	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, gerror.New("invalid id format")
	}

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: objectId},
		}},
	}
	limitStage := bson.D{
		{Key: "$limit", Value: 1},
	}

	projectStage := ProjectStage
	pipeline := mongo.Pipeline{
		matchStage,
		limitStage,
		projectStage,
	}
	cursor, err := imongo.TableSubject.Collection(ctx).Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("failed to aggregate: %v", err)
	}
	var guess subjectinfo
	if cursor.Next(ctx) {
		if err := cursor.Decode(&guess); err != nil {
			return nil, err
		}
	} else {
		return nil, gerror.New("no document found")
	}
	res = &v1.GuessRes{}
	if guess.ID == answer.ID {
		res.Success = true
		return res, nil
	} else {
		res.Success = false
		res.Result.Name.Value = guess.Name
		res.Result.ID.Value = guess.ID
	}

	res.Result.ReleaseDate.Value = guess.ReleaseDate
	guessTime := gtime.NewFromStr(guess.ReleaseDate)
	answerTime := gtime.NewFromStr(answer.ReleaseDate)

	// 计算天数差
	guessTimeDiff := guessTime.Sub(answerTime).Hours() / 24

	switch {
	case guessTimeDiff == 0:
		res.Result.ReleaseDate.Value = "status_equal"
	case guessTimeDiff > 365:
		res.Result.ReleaseDate.Value = "status_large_late"
	case guessTimeDiff < -365:
		res.Result.ReleaseDate.Value = "status_large_early"
	case guessTimeDiff > 0:
		res.Result.ReleaseDate.Value = "status_early"
	case guessTimeDiff < 0:
		res.Result.ReleaseDate.Value = "status_late"
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid guess time diff")
	}

	for platform := range guess.Platforms {
		tmp := iostruct.GuessResult{
			Value: platform,
		}
		if answer.Platforms[platform] {
			tmp.Status = "status_equal"
		} else {
			tmp.Status = "status_missing"
		}
		res.Result.Platforms = append(res.Result.Platforms, tmp)
	}

	for tag := range guess.Tags {
		tmp := iostruct.GuessResult{
			Value: tag,
		}
		if answer.Tags[tag] {
			tmp.Status = "status_equal"
		} else {
			tmp.Status = "status_missing"
		}
		res.Result.Tags = append(res.Result.Tags, tmp)
	}

	return res, nil
}
