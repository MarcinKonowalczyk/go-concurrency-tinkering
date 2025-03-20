package log

import (
	"context"
	"fmt"
	_log "log"
	"math/rand"
	"net/http"

	"golang.org/x/exp/constraints"
)

type ID[T constraints.Integer] interface {
	ID() T
}

type baseid struct {
	id int64
}

func (i baseid) ID() int64 {
	return i.id
}

// Note: This is a painful way of doing this. Ideally the tests would then check if we do actually cover all the
// cases of `constraints.Integer`. Also, probably much simpler to just use `int64` everywhere.
// But this is a good exercise in generics.
func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case int16:
		return int64(v), true
	case int8:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	case uintptr:
		return int64(v), true

	case ID[int64]:
		return v.ID(), true
	case ID[int32]:
		return int64(v.ID()), true
	case ID[int16]:
		return int64(v.ID()), true
	case ID[int8]:
		return int64(v.ID()), true
	case ID[uint64]:
		return int64(v.ID()), true
	case ID[uint32]:
		return int64(v.ID()), true
	case ID[uint16]:
		return int64(v.ID()), true
	case ID[uint8]:
		return int64(v.ID()), true
	case ID[int]:
		return int64(v.ID()), true
	case ID[uint]:
		return int64(v.ID()), true
	case ID[uintptr]:
		return int64(v.ID()), true

	default:
		return 0, false
	}
}

func idFromContext(ctx context.Context) (ID[int64], bool) {
	value := ctx.Value("id")
	if value == nil {
		// no value
		return nil, false
	}
	id, ok := toInt64(value)
	if ok {
		return baseid{id: id}, true
	}

	return nil, false
}


func Println(ctx context.Context, v ...any) {
	id, ok := idFromContext(ctx)
	if ok {
		v = append([]any{fmt.Sprintf("[%d]", id.ID())}, v...)
	}
	_log.Println(v...)
}

func Fatal(ctx context.Context, v ...any) {
	id, ok := idFromContext(ctx)
	if ok {
		v = append([]any{fmt.Sprintf("[%d]", id.ID())}, v...)
	}
	_log.Fatal(v...)	
}

func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add request ID to the context
		ctx := r.Context();
		id := rand.Int63() % 1000000;
		ctx = context.WithValue(ctx, "id", id)
		f(w, r.WithContext(ctx))
	}
}