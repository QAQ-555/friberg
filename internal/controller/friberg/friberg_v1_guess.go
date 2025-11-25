package friberg

import (
	"context"
	v1 "friberg/api/friberg/v1"
	imongo "friberg/internal/library/mongo"
	"friberg/internal/model/iostruct"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ControllerV1) Guess(ctx context.Context, req *v1.GuessReq) (res *v1.GuessRes, err error) {
	g.Log().Debug(ctx, "guess req:", req)
	res = &v1.GuessRes{}

	r := g.RequestFromCtx(ctx)

	// =============================
	// 获取玩家 session 数据
	// =============================
	sessionData, err := r.Session.Data()
	if err != nil {
		return nil, err
	}

	answer := subjectinfo{}
	err = gconv.Struct(sessionData["subject"], &answer)
	if err != nil {
		return nil, err
	}
	g.Dump("session subject:", answer)

	// =============================
	// 更新 Frequency 并处理 Session 生命周期
	// =============================
	answer.Frequency = answer.Frequency - 1
	res.Frequency = answer.Frequency

	if answer.Frequency <= 0 {
		r.Session.Close() // 最后一次机会或猜对时释放
	} else {
		r.Session.Set("subject", answer)
	}

	// =============================
	// Mongo 查询真实 subject
	// =============================
	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, gerror.New("invalid id format")
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
		bson.D{{Key: "$limit", Value: 1}},
		ProjectStage,
	}

	cursor, err := imongo.TableSubject().Coll().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var guess subjectinfo
	if cursor.Next(ctx) {
		if err := cursor.Decode(&guess); err != nil {
			return nil, err
		}
	} else {
		return nil, gerror.New("no document found")
	}

	// =============================
	// 判断猜测是否成功
	// =============================
	if guess.ID == answer.ID {
		res.Success = true
		r.Session.Close()
		return res, nil
	}
	res.Success = false

	// =============================
	// 填充结果 Response
	// =============================
	res.Result.ID.Value = guess.ID
	res.Result.Name.Value = guess.Name

	// ReleaseDate 差值计算
	guessTime := gtime.NewFromStr(guess.ReleaseDate)
	answerTime := gtime.NewFromStr(answer.ReleaseDate)
	daysDiff := guessTime.Sub(answerTime).Hours() / 24

	switch {
	case daysDiff == 0:
		res.Result.ReleaseDate.Value = "status_equal"
	case daysDiff > 365:
		res.Result.ReleaseDate.Value = "status_large_late"
	case daysDiff < -365:
		res.Result.ReleaseDate.Value = "status_large_early"
	case daysDiff > 0:
		res.Result.ReleaseDate.Value = "status_early"
	case daysDiff < 0:
		res.Result.ReleaseDate.Value = "status_late"
	}

	// 平台和标签对比
	for platform := range guess.Platforms {
		status := "status_missing"
		if answer.Platforms[platform] {
			status = "status_equal"
		}
		res.Result.Platforms = append(res.Result.Platforms, iostruct.GuessResult{
			Value:  platform,
			Status: status,
		})
	}

	for tag := range guess.Tags {
		status := "status_missing"
		if answer.Tags[tag] {
			status = "status_equal"
		}
		res.Result.Tags = append(res.Result.Tags, iostruct.GuessResult{
			Value:  tag,
			Status: status,
		})
	}

	// =============================
	//  返回结果
	// =============================
	return res, nil
}
