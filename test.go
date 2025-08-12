package main

import (
	"NewCsTeamServer/profile"
	"fmt"
	"os"
)

func main() {
	data, err := os.Open("jquery-c2.4.5.profile")
	if err != nil {
		fmt.Println("[info]", "Not Find Profile!")
		return
	}
	profile.ProfileConfig = profile.GetProfile(data)
	fmt.Println(profile.ProfileConfig)

}
