package counter

type Unsigned uint

func NewUnsigned() *Unsigned {
	u := Unsigned(0)
	return &u
}

func (u *Unsigned) Inc() {
	*u++
}

func (u *Unsigned) Get() uint {
	return uint(*u)
}
