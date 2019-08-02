package utils

import "testing"

func TestCountDigits(t *testing.T) {

	if val := CountDigits(10); val != 2 {
		t.Errorf("Input=10. The digits should be 2. Get %v.", val)
	}

	if val := CountDigits(0); val != 0 {
		t.Errorf("Input=0. The digits should be 0. Get %v.", val)
	}

	if val := CountDigits(1000000000); val != 10 {
		t.Errorf("Input=1000000000. The digits should be 10. Get %v.", val)
	}
}
