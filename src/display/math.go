package display

import "fmt"

type single struct {
	value uint8
	last  uint8
}

type double struct {
	value int16
	last  int16
}

func (d *single) from(r *Display) {
	b := r.next()
	d.value = uint8(b)
}

func (d *single) changed() bool {
	if d.last != d.value {
		d.last = d.value
		return true
	}
	return false
}

func (d *single) get() uint8 { d.last = d.value; return d.value }

func (d *double) from(r *Display) {
	ub := r.next()
	lb := r.next()
	w := int(ub)<<8 | int(lb)
	d.value = int16(w + 32768)
}

func (d *single) String() string {
	return fmt.Sprintf("%d", d.value)
}

func (d *double) changed() bool {
	if d.last != d.value {
		d.last = d.value
		return true
	}
	return false
}

func (d *double) String() string {
	return fmt.Sprintf("%d", d.value)
}

func (d *double) get() int16 { d.last = d.value; return d.value }
