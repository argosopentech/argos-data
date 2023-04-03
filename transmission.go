package main

import (
	"fmt"
	"os/exec"
)

// Requires transmission-cli to be installed
// sudo apt install transmission-cli

func TransmissionShowMagnetLink(torrentFile string) string {
	// Example shell command
	// transmission-show -m Argos-Translate-LibreTranslate-2022-04-30.zip.torrent

	// Create command
	cmd := exec.Command("transmission-show", "-m", torrentFile)

	// Run command
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Convert output to string
	magnetLink := string(out)

	return magnetLink
}

func Transmission(link string) {
	// Example shell command
	// transmission-cli magnet:?xt=urn:btih...

	// Create command
	cmd := exec.Command("transmission-cli", link)

	// Run command
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Get magnet link
	var torrentFile string = "Argos-Translate-LibreTranslate-2022-04-30.zip.torrent"
	var magnetLink string = TransmissionShowMagnetLink(torrentFile)

	// Download torrent
	Transmission(magnetLink)
}
