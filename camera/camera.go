package camera

type Camera struct {
	X, Y       float64
	ScreenW    int
	ScreenH    int
	TargetX    float64
	TargetY    float64
	Smoothing  float64
}

func New(w, h int) *Camera {
	return &Camera{
		X:        0,
		Y:        0,
		ScreenW:  w,
		ScreenH:  h,
		Smoothing: 0.1,
	}
}

func (c *Camera) Follow(targetX, targetY float64) {
	c.TargetX = targetX - float64(c.ScreenW)/3
	c.TargetY = targetY - float64(c.ScreenH)/2
	
	if c.TargetX < 0 {
		c.TargetX = 0
	}
	
	c.X += (c.TargetX - c.X) * c.Smoothing
	c.Y += (c.TargetY - c.Y) * c.Smoothing
}

func (c *Camera) WorldToScreen(x, y float64) (float64, float64) {
	return x - c.X, y - c.Y
}

func (c *Camera) Reset() {
	c.X = 0
	c.Y = 0
}
