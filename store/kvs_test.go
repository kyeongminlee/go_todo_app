package store

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go_todo_app/entity"
	"go_todo_app/testutil"
	"testing"
	"time"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	client := testutil.OpenRedisForTest(t)

	sut := &KVS{Client: client}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	t.Cleanup(func() {
		client.Del(ctx, key)
	})

	//if err := sut.Save(ctx, key, uid); err != nil {
	//	t.Errorf("want no error, got %v", err)
	//}

	err := sut.Save(ctx, key, uid)
	assert.NoError(t, err, "want no error")
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	client := testutil.OpenRedisForTest(t)
	sut := &KVS{Client: client}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		client.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			client.Del(ctx, key)
		})

		got, err := sut.Load(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, uid, got)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Save_notFound"
		ctx := context.Background()

		//got, err := sut.Load(ctx, key)
		//if err == nil || !errors.Is(err, ErrNotFound) {
		//	//t.Errorf("want %v, but got %v(value = %d)", ErrNotFound, err, got)
		//}
		_, err := sut.Load(ctx, key)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
