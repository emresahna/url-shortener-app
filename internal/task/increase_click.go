package task

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"strings"
)

func (ts *task) IncreaseClickTask(ctx context.Context, t *asynq.Task) error {
	keys, _ := ts.rc.Keys(ctx, "clicks:*").Result()

	for _, key := range keys {
		shortCode := strings.Split(key, ":")
		originalUrl, _ := ts.db.GetURLByCode(ctx, shortCode[1])
		fmt.Println(originalUrl)
	}

	return nil
}
