package core

import "context"

func (r *mutationResolver) PerformManualMovement(ctx context.Context, vec ManualMovementPositionVector) (*bool, error) {

	err := r.App.ManualMovementService.MoveRelative(&PositionVector{
		X: vec.X,
		Y: vec.Y,
		Z: vec.Z,
		E: vec.E,
	})
	return nil, err
}
