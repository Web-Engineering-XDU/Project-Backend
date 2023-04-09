package interfaces

import "context"

type Agent interface {
	Run(ctx context.Context , e Event)
}