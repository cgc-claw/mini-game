package physics

type AABB struct {
	X, Y          float64
	Width, Height float64
}

func NewAABB(x, y, w, h float64) *AABB {
	return &AABB{X: x, Y: y, Width: w, Height: h}
}

func (a *AABB) Intersects(b *AABB) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

func (a *AABB) MinX() float64  { return a.X }
func (a *AABB) MaxX() float64  { return a.X + a.Width }
func (a *AABB) MinY() float64  { return a.Y }
func (a *AABB) MaxY() float64  { return a.Y + a.Height }

func (a *AABB) Top() float64   { return a.Y }
func (a *AABB) Bottom() float64 { return a.Y + a.Height }
func (a *AABB) Left() float64   { return a.X }
func (a *AABB) Right() float64  { return a.X + a.Width }

const (
	Gravity       = 0.5
	JumpForce     = -12
	MoveSpeed     = 5
	TerminalVelocity = 15
)

func ApplyGravity(vy float64) float64 {
	vy += Gravity
	if vy > TerminalVelocity {
		vy = TerminalVelocity
	}
	return vy
}
