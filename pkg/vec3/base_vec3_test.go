package vec3

import (
	"testing"
)

func TestEquals(t *testing.T) {
	v1 := Vec3{3, 3, 3}
	v2 := Vec3{3, 3, 3}
	v3 := Vec3{-5, 0, 2}

	if !v1.Eq(v2) {
		t.Errorf("Expected %v == %v. Got %v != %v", v1, v2, v1, v2)
		t.Fail()
	}
	if v1.Eq(v3) {
		t.Errorf("Expected %v != %v. Got %v == %v", v1, v3, v1, v3)
		t.Fail()
	}

}

func TestPlus(t *testing.T) {

	v1 := Vec3{3, 3, 3}
	v2 := v1.Plus(2)

	if !(v2.X == 5 && v2.Y == 5 && v2.Z == 5) {
		t.Fail()
	}

}
