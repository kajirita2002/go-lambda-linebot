package main

import "linebot/gurunavi"

func TextRestaurants(g *gurunavi.GurunaviResponseBody) string {
	var t string
	for _, r := range g.Rest {
		t += r.Name + "\n" + r.URL + "\n"
	}
	return t
}
