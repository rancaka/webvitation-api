package main

import (
	"fmt"

	"github.com/gosimple/slug"
)

func main() {
	asd := []string{
		"Iin",
		"Susi dan Imran",
		"Elisabet Bunga Eria",
		"Rida Demsy",
		"Yolanda",
	}

	// fmt.Println("update invitees set invitation_time_from = '19:00', invitation_time_to = '20:00' where slug in (")
	fmt.Println("insert into invitees(name, slug, event_id) values")
	for _, t := range asd {
		sl := slug.Make(t)
		// fmt.Printf(`'%s',`, sl)
		fmt.Printf(`('%s', '%s', 'a7219998-ba5d-11eb-ae4c-42010ab8005f'),`, t, sl)
	}

	fmt.Print("\n\n")

	for _, t := range asd {
		sl := slug.Make(t)
		fmt.Printf("https://tamuindonesia.com/gizela-dan-fajri-wedding/to/%s\n", sl)
	}
}
