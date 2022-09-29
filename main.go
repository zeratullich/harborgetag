package main

import (
	"flag"
	"fmt"
	"harborgetag/imagetag"
	"harborgetag/tools"
	"log"
)

var username = flag.String("username", "admin", "Harbor auth user name")
var password = flag.String("password", "123456", "Harbor auth user password")
var registry = flag.String("registry", "harbor.myquanyi.com", "Harbor registered mirror address")
var image = flag.String("image", "qy/book-store", "Harbor registered mirror image")
var order = flag.String("order", "Descending-Versions", `Allows the user to alter the ordering of the ImageTags in the build parameter.

	Natural-Ordering ... same Ordering as the tags had in prior versions
	Reverse-Natural-Ordering ... the reversed original ordering
	Descending-Versions ... attempts to pars the tags to a version and order them descending
	Ascending-Versions ... attempts to pars the tags to a version and order them ascending
`)
var verifySSL = flag.Bool("verifySSL", true, "Whether ssl authentication is enabled in harbor")
var filter = flag.String("filter", ".*", `Regular expression to filter image tag e.g. v(\d+\.)*\d+ for tags like v23.3.2`)
var scheme = tools.StringEnumVar("scheme", "https", "URL request prefix, only one of [http, https] can be selected")

func main() {
	flag.Parse()

	imageTag := imagetag.ImageTag{Username: *username, Password: *password, Registry: *registry, Image: *image, VerifySSL: *verifySSL, Order: *order, Filter: *filter}

	if *scheme == "http" {
		imageTag.VerifySSL = false
		imageTag.Scheme = "http://"
	} else {
		imageTag.Scheme = "https://"
	}

	err := imageTag.GetImageTagsFromRegistry()
	if err != nil {
		fmt.Printf("get image tag failed: %v", err)
	} else {
		log.Println("get image tags is:", imageTag.Tags)
		for k := range imageTag.Tags {
			fmt.Printf("%s:%s\n", *image, imageTag.Tags[k])
		}
	}
}
