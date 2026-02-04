package levels

import "github.com/PawelZabc/ProjektZespolowy/internal/game/physics"

var BasicLevel = []Change{
	{Value: 20, Axis: physics.DirX},
	{Value: 5, Axis: physics.DirZ},
	{Value: 5, Axis: physics.DirX},

	{Value: -5, Axis: physics.DirZ},
	{Value: 2, Axis: physics.DirX},

	{Value: -5, Axis: physics.DirZ},
	{Value: -2, Axis: physics.DirX},
	{Value: -7, Axis: physics.DirZ},
	{Value: 10, Axis: physics.DirX},
	{Value: 7, Axis: physics.DirZ},
	{Value: -3, Axis: physics.DirX},
	{Value: 5, Axis: physics.DirZ},

	{Value: 3, Axis: physics.DirX},
	{Value: 30, Axis: physics.DirZ},
	{Value: -4, Axis: physics.DirX},

	// corridors1 beg

	{Value: 5, Axis: physics.DirZ},
	{Value: -7, Axis: physics.DirX},

	{Value: -3, Axis: physics.DirX, Skip: true},
	{Value: -4, Axis: physics.DirX},
	{Value: 4, Axis: physics.DirZ},
	{Value: 4, Axis: physics.DirX},
	{Value: -4, Axis: physics.DirZ},
	{Value: 3, Axis: physics.DirX, Skip: true},

	{Value: 7, Axis: physics.DirZ},
	{Value: -10, Axis: physics.DirX},
	{Value: -7, Axis: physics.DirZ},

	{Value: -4, Axis: physics.DirX},
	{Value: -3, Axis: physics.DirZ, Skip: true},
	{Value: 4, Axis: physics.DirX},

	{Value: -7, Axis: physics.DirZ},
	{Value: 3, Axis: physics.DirX},
	{Value: 7, Axis: physics.DirZ},
	{Value: 4, Axis: physics.DirX},
	{Value: -7, Axis: physics.DirZ},
	{Value: 3, Axis: physics.DirX},
	{Value: 7, Axis: physics.DirZ},
	{Value: 4, Axis: physics.DirX},
	{Value: -2, Axis: physics.DirZ},

	{Value: -3, Axis: physics.DirX},

	// corridors1 end

	{Value: -20, Axis: physics.DirZ},
	{Value: -5, Axis: physics.DirX},

	{Value: 10, Axis: physics.DirZ},
	{Value: -8, Axis: physics.DirX},

	// black room beg
	{Value: 4, Axis: physics.DirZ},
	{Value: 1, Axis: physics.DirX},
	{Value: 6, Axis: physics.DirZ},
	{Value: -6, Axis: physics.DirX},
	{Value: -6, Axis: physics.DirZ},
	{Value: 1, Axis: physics.DirX},
	{Value: -4, Axis: physics.DirZ},
	// black room end

	{Value: -8, Axis: physics.DirX},
	{Value: -10, Axis: physics.DirZ},
	{Value: -5, Axis: physics.DirX},
	{Value: 20, Axis: physics.DirZ},
	{Value: -4, Axis: physics.DirX},

	// corridors 2 beg

	{Value: 2, Axis: physics.DirZ},
	{Value: 5, Axis: physics.DirX},

	{Value: 3, Axis: physics.DirX, Skip: true},
	{Value: -4, Axis: physics.DirZ},
	{Value: 4, Axis: physics.DirX},
	{Value: 4, Axis: physics.DirZ},
	{Value: -4, Axis: physics.DirX},
	{Value: -3, Axis: physics.DirX, Skip: true},

	{Value: -11, Axis: physics.DirZ},
	{Value: 10, Axis: physics.DirX},
	{Value: 11, Axis: physics.DirZ},

	{Value: 4, Axis: physics.DirX},
	{Value: 3, Axis: physics.DirZ, Skip: true},
	{Value: -4, Axis: physics.DirX},

	{Value: 7, Axis: physics.DirZ},
	{Value: -3, Axis: physics.DirX},
	{Value: -7, Axis: physics.DirZ},
	{Value: -4, Axis: physics.DirX},
	{Value: 7, Axis: physics.DirZ},
	{Value: -3, Axis: physics.DirX},
	{Value: -7, Axis: physics.DirZ},

	{Value: -8, Axis: physics.DirX},
	{Value: -5, Axis: physics.DirZ},

	//corridors2 end

	{Value: -3, Axis: physics.DirX},
	{Value: -30, Axis: physics.DirZ},

	// purple room beg
	{Value: 2, Axis: physics.DirX},
	{Value: -9, Axis: physics.DirZ},
	{Value: 26, Axis: physics.DirX},
	{Value: -3, Axis: physics.DirZ},
	{Value: 10, Axis: physics.DirX},
	{Value: 11, Axis: physics.DirZ},
	{Value: -10, Axis: physics.DirX},
	{Value: -3, Axis: physics.DirZ},
	{Value: -21, Axis: physics.DirX},
	{Value: 4, Axis: physics.DirZ},
	{Value: 3, Axis: physics.DirX},
	// purple room end

	{Value: 5, Axis: physics.DirZ},
	{Value: 5, Axis: physics.DirX},
	{Value: -5, Axis: physics.DirZ},
}
