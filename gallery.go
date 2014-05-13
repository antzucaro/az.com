package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

import _ "image/jpeg"

const MEDIA_URL string = "http://media.antzucaro.com"

var dir string

func main() {
	check_args()

	images := getImagesByExt(".jpg")
	if len(images) > 0 {
		dir = create_dest_dir()
	}

	for _, image := range images {
		processImage(image)
	}
}

func getImagesByExt(targetExt string) (images []string) {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext == targetExt {
			images = append(images, path)
		}
		return nil
	})
	return
}

func check_args() {
	// Gather arguments
	if len(os.Args) < 2 {
		fmt.Println("Missing argument: directory name")
		os.Exit(1)
	}
}

func create_dest_dir() (dir string) {
	now := time.Now()
	dir = fmt.Sprintf("%d/%d/%s", now.Year(), now.Month(), os.Args[1])

	m_dir := fmt.Sprintf("%s/m", dir)
	l_dir := fmt.Sprintf("%s/l", dir)

	if err := os.MkdirAll(m_dir, os.FileMode(0755)); err != nil {
		log.Fatal("Issue creating medium directory!")
	}
	if err := os.MkdirAll(l_dir, os.FileMode(0755)); err != nil {
		log.Fatal("Issue creating large directory!")
	}
	return
}

func parse_fn(fn string) (bn string, ext string) {
	// Extract a basename and extension from the current file
	ext = filepath.Ext(fn)
	bn = fn[:len(fn)-len(ext)]
	return
}

func processImage(path string) error {
	// basename, extension of the file
	bn, ext := parse_fn(path)

	// calculate how much imagemagick is going to resize by
	m_option, l_option := calc_resize_options(path)

	// actually do the resizing
	m_fn, l_fn := resize(path, bn, ext, m_option, l_option)

	// prompt for caption on STDIN
	caption := get_caption(path)

	// spit out the HTML
	generate_html(m_fn, l_fn, caption)

	return nil
}

func calc_resize_options(fn string) (m_option string, l_option string) {

	C_L_SIZE := 1000
	C_M_SIZE := 650
	PANO_L_SIZE := 1600
	PANO_M_SIZE := 650

	width, height := get_width_height(fn)

	if width > height {
		// landscape
		if width > 4000 {
			l_option = fmt.Sprintf("%dx", PANO_L_SIZE)
			m_option = fmt.Sprintf("%dx", PANO_M_SIZE)
		} else {
			l_option = fmt.Sprintf("%dx", C_L_SIZE)
			m_option = fmt.Sprintf("%dx", C_M_SIZE)
		}
	} else {
		// portrait
		if height > 4000 {
			l_option = fmt.Sprintf("x%d", PANO_L_SIZE)
			m_option = fmt.Sprintf("x%d", PANO_M_SIZE)
		} else {
			l_option = fmt.Sprintf("x%d", C_L_SIZE)
			m_option = fmt.Sprintf("x%d", C_M_SIZE)
		}
	}
	return
}

func resize(fn string, bn string, ext string, m_option string, l_option string) (m_fn string, l_fn string) {
	l_fn = fmt.Sprintf("%s/l/%s_l%s", dir, bn, ext)
	m_fn = fmt.Sprintf("%s/m/%s_m%s", dir, bn, ext)
	go exec.Command("convert", "-resize", l_option, fn, l_fn).Run()
	go exec.Command("convert", "-resize", m_option, "-quality", "65", fn, m_fn).Run()
	return
}

func get_caption(fn string) (caption string) {
	stdin := bufio.NewReader(os.Stdin)

	go exec.Command("/usr/bin/feh", "-.", fn).Run()

	fmt.Printf("Enter caption: ")
	caption, _ = stdin.ReadString('\n')
	caption = strings.TrimSpace(caption)

	return
}

func generate_html(m_fn string, l_fn string, caption string) {
	MEDIA_URL := "/uploads/"

	fn := fmt.Sprintf("%s.html", os.Args[1])

	var of *os.File
	of, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		of, err = os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0755)
	}

	defer of.Close()

    /*
	fmt.Fprintf(of, "<div class='wp-caption aligncenter'>\n")
	fmt.Fprintf(of, "  <a href=\"%s%s\" title=\"%s\">\n", MEDIA_URL, l_fn, caption)
	fmt.Fprintf(of, "    <img alt=\"%s\" title=\"%s\" src=\"%s%s\">\n", caption, caption, MEDIA_URL, m_fn)
	fmt.Fprintf(of, "  </a>\n")
	fmt.Fprintf(of, "    <p class='wp-caption-text'>%s</p>\n", caption)
	fmt.Fprintf(of, "</div>\n\n")
    */

	fmt.Fprintln(of, "{{% polaroid")
	fmt.Fprintf(of, "   \"%s%s\"\n", MEDIA_URL, m_fn)
	fmt.Fprintf(of, "   \"%s%s\"\n", MEDIA_URL, l_fn)
	fmt.Fprintf(of, "   \"%s\"\n", caption)
	fmt.Fprintln(of, "%}}\n")
}

func get_width_height(fn string) (width int, height int) {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	width, height = bounds.Size().X, bounds.Size().Y

	return
}
