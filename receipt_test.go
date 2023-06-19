package main

import "testing"

func TestCalculatePoints(t *testing.T) {

	// Tests morning-receipt.json
	// 9 characters, 5 points for 2 items, 1 point for item description ceil(1.40 * 0.2) = 1
	// Total : 15 points

	got := calculatePoints(Receipt{"Walgreens", "2022-01-02", "08:13", "2.65", []ItemsList{{"Pepsi - 12-oz", "1.25"}, {"Dasani", "1.40"}}})
	want := 15

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	// Tests simple-receipt.json
	// 6 characters, 25 points 1.25 mod 0.25 = 0
	// Total : 31 points

	got1 := calculatePoints(Receipt{"Target", "2022-01-02", "13:13", "1.25", []ItemsList{{"Pepsi - 12-0z", "1.25"}}})
	want1 := 31

	if got1 != want1 {
		t.Errorf("got %q, wanted %q", got1, want1)
	}
}
