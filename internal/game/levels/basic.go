package levels

import "github.com/PawelZabc/ProjektZespolowy/internal/game/physics"

var BasicLevel = []Change{
	{Value: 20, Axis: physics.DirX},
	{Value: 5, Axis: physics.DirZ},
	{Value: 5, Axis: physics.DirX},
	{Value: -5, Axis: physics.DirZ},
	{Value: 10, Axis: physics.DirX},
	{Value: 30, Axis: physics.DirZ},
	{Value: -10, Axis: physics.DirX},
	{Value: -20, Axis: physics.DirZ},
	{Value: -5, Axis: physics.DirX},
	{Value: 10, Axis: physics.DirZ},
	{Value: -20, Axis: physics.DirX},
	{Value: -10, Axis: physics.DirZ},
	{Value: -5, Axis: physics.DirX},
	{Value: 20, Axis: physics.DirZ},
	{Value: -10, Axis: physics.DirX},
	{Value: -30, Axis: physics.DirZ},
	{Value: 10, Axis: physics.DirX},
	{Value: 5, Axis: physics.DirZ},
	{Value: 5, Axis: physics.DirX},
	{Value: -5, Axis: physics.DirZ},
}
