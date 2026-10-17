package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"chanlearning/pkg/domain/models"
)

func ProcessText(ctx context.Context, body []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var task models.TextTask
	if err := json.Unmarshal(body, &task); err != nil {
		return fmt.Errorf("decode text task: %w", err)
	}
	if task.ID == "" || strings.TrimSpace(task.Text) == "" {
		return errors.New("text task must include an ID and non-empty text")
	}

	log.Printf("processed task id=%s result=%q", task.ID, strings.ToUpper(task.Text))
	return nil
}
