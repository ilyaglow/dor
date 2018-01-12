package dor

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// downloadURL downloads file to a temporary file
func downloadURL(url string, descr string) (filename string, error error) {
	zf, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("dor-%s-", descr))
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		zf.Close()
		os.Remove(zf.Name())
		return "", err
	}
	defer resp.Body.Close()
	defer zf.Close()

	_, err = io.Copy(zf, resp.Body)
	if err != nil {
		zf.Close()
		os.Remove(zf.Name())
		return "", err
	}

	return zf.Name(), nil
}

// zipContent unpacks zip file in memory with one file inside
func zipContent(path string) (zipFile *zip.ReadCloser, f *io.ReadCloser, error error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}

	rc, err := r.File[0].Open()
	if err != nil {
		r.Close()
		return nil, nil, err
	}

	return r, &rc, nil
}

// mapFromURLZip translates a common formatted CSV:
//	rank,domain
// that is packed in ZIP and hosted on a specified url
func mapFromURLZip(url string, desc string) (LookupMap, error) {
	n, err := downloadURL(url, desc)
	if err != nil {
		return nil, err
	}

	z, c, err := zipContent(n)
	if err != nil {
		return nil, err
	}

	defer (*c).Close()
	defer z.Close()
	defer os.Remove(n)

	m := make(LookupMap)
	scanner := bufio.NewScanner(*c)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		pos, d := parts[0], parts[1]
		pint, _ := strconv.ParseInt(pos, 10, 32)
		puint := uint(pint)
		m[d] = puint
	}

	return m, nil
}
