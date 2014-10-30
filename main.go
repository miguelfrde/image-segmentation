package main

import (
	"fmt"
	"github.com/miguelfrde/image-segmentation/graph"
	"github.com/miguelfrde/image-segmentation/segmentation"
	"html/template"
	"image"
	"image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const RANDOM_STR_SIZE = 25

var templates = template.Must(template.ParseGlob("web/templates/*"))
var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")

func randomString() string {
	chars := make([]byte, RANDOM_STR_SIZE)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func loadImageFromFile(filename string) image.Image {
	f, _ := os.Open(filename)
	defer f.Close()
	img, _, _ := image.Decode(f)
	return img
}

func createFileInFS(file multipart.File, extension string) (string, error) {
	defer file.Close()

	newfilename := randomString()
	imgfile, err := os.Create("tmp/" + newfilename + extension)
	if err != nil {
		return "", err
	}

	defer imgfile.Close()

	_, err = io.Copy(imgfile, file)
	if err != nil {
		return "", err
	}
	return newfilename, nil
}

/* Handlers */

func mainHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "main", nil)
}

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	extension := filepath.Ext(header.Filename)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	filename, err := createFileInFS(file, extension)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	sigma, err := strconv.ParseFloat(r.FormValue("sigma"), 64)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	graphType := graph.KINGSGRAPH
	if r.FormValue("graph") == "2" {
		graphType = graph.GRIDGRAPH
	}

	weightfn := segmentation.NNWeight
	if r.FormValue("weightfn") == "2" {
		weightfn = segmentation.IntensityDifference
	}

	fmt.Println("Segmenting requested image:", header.Filename, "as", filename)
	img := loadImageFromFile("tmp/" + filename + extension)
	segmenter := segmentation.New(img, graphType, weightfn)
	if r.FormValue("color") == "on" {
		segmenter.SetRandomColors(true)
	}

	if algorithm := r.FormValue("algorithm"); algorithm == "1" {
		fmt.Println("Using GBS")
		k, err := strconv.ParseFloat(r.FormValue("k"), 64)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		minSize, err := strconv.ParseFloat(r.FormValue("minsize"), 32)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		segmenter.SegmentGBS(sigma, k, int(minSize))
	} else {
		fmt.Println("Using HMSF")
		minWeight, err := strconv.ParseFloat(r.FormValue("minweight"), 64)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		segmenter.SegmentHMSF(sigma, minWeight)
	}

	toimg, _ := os.Create("tmp/new_" + filename + ".png")
	defer toimg.Close()
	png.Encode(toimg, segmenter.GetResultImage())
	fmt.Fprintln(w, filename, extension)
}

func servePublicFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for public file:", r.URL.Path[1:])
	http.ServeFile(w, r, "web/public/"+r.URL.Path[1:])
}

func serveTmpFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for temporal file:", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/tmp/", serveTmpFile)
	http.HandleFunc("/segment", segmentHandler)

	/* Static files */
	for _, dir := range []string{"css", "img", "js", "components"} {
		http.HandleFunc("/"+dir+"/", servePublicFile)
	}

	fmt.Println("Listening on port " + os.Getenv("PORT") + "...")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
