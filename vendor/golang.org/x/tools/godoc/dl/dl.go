// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package dl implements a simple downloads frontend server.
//
// It accepts HTTP POST requests to create a new download metadata entity, and
// lists entities with sorting and filtering.
// It is designed to run only on the instance of godoc that serves golang.org.
package dl

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	"google.golang.org/cloud/compute/metadata"
)

const (
	gcsBaseURL    = "https://storage.googleapis.com/golang/"
	cacheKey      = "download_list_3" // increment if listTemplateData changes
	cacheDuration = time.Hour
)

var builderKey string

func init() {
	builderKey, _ = metadata.ProjectAttributeValue("builder-key")
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.Handle("/dl", http.RedirectHandler("/dl/", http.StatusFound))
	mux.HandleFunc("/dl/", getHandler) // also serves listHandler
	mux.HandleFunc("/dl/upload", uploadHandler)
	mux.HandleFunc("/dl/init", initHandler)
}

type File struct {
	Filename string
	OS       string
	Arch     string
	Version  string
	Checksum string `datastore:",noindex"`
	Size     int64  `datastore:",noindex"`
	Kind     string // "archive", "installer", "source"
	Uploaded time.Time
}

func (f File) PrettyOS() string {
	if f.OS == "darwin" {
		switch {
		case strings.Contains(f.Filename, "osx10.8"):
			return "OS X 10.8+"
		case strings.Contains(f.Filename, "osx10.6"):
			return "OS X 10.6+"
		}
	}
	return pretty(f.OS)
}

func (f File) PrettySize() string {
	const mb = 1 << 20
	if f.Size == 0 {
		return ""
	}
	if f.Size < mb {
		// All Go releases are >1mb, but handle this case anyway.
		return fmt.Sprintf("%v bytes", f.Size)
	}
	return fmt.Sprintf("%.0fMB", float64(f.Size)/mb)
}

func (f File) Highlight() bool {
	switch {
	case f.Kind == "source":
		return true
	case f.Arch == "amd64" && f.OS == "linux":
		return true
	case f.Arch == "amd64" && f.Kind == "installer":
		switch f.OS {
		case "windows":
			return true
		case "darwin":
			if !strings.Contains(f.Filename, "osx10.6") {
				return true
			}
		}
	}
	return false
}

func (f File) URL() string {
	return gcsBaseURL + f.Filename
}

type Release struct {
	Version string
	Stable  bool
	Files   []File
	Visible bool // show files on page load
}

type Feature struct {
	// The File field will be filled in by the first stable File
	// whose name matches the given fileRE.
	File
	fileRE *regexp.Regexp

	Platform     string // "Microsoft Windows", "Mac OS X", "Linux"
	Requirements string // "Windows XP and above, 64-bit Intel Processor"
}

// featuredFiles lists the platforms and files to be featured
// at the top of the downloads page.
var featuredFiles = []Feature{
	{
		Platform:     "Microsoft Windows",
		Requirements: "Windows XP or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.windows-amd64\.msi$`),
	},
	{
		Platform:     "Apple OS X",
		Requirements: "OS X 10.8 or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.darwin-amd64(-osx10\.8)?\.pkg$`),
	},
	{
		Platform:     "Linux",
		Requirements: "Linux 2.6.23 or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.linux-amd64\.tar\.gz$`),
	},
	{
		Platform: "Source",
		fileRE:   regexp.MustCompile(`\.src\.tar\.gz$`),
	},
}

// data to send to the template; increment cacheKey if you change this.
type listTemplateData struct {
	Featured         []Feature
	Stable, Unstable []Release
	LoginURL         string
}

var (
	listTemplate  = template.Must(template.New("").Funcs(templateFuncs).Parse(templateHTML))
	templateFuncs = template.FuncMap{"pretty": pretty}
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var (
		c = appengine.NewContext(r)
		d listTemplateData
	)
	if _, err := memcache.Gob.Get(c, cacheKey, &d); err != nil {
		if err == memcache.ErrCacheMiss {
			log.Debugf(c, "cache miss")
		} else {
			log.Errorf(c, "cache get error: %v", err)
		}

		var fs []File
		_, err := datastore.NewQuery("File").Ancestor(rootKey(c)).GetAll(c, &fs)
		if err != nil {
			log.Errorf(c, "error listing: %v", err)
			return
		}
		d.Stable, d.Unstable = filesToReleases(fs)
		if len(d.Stable) > 0 {
			d.Featured = filesToFeatured(d.Stable[0].Files)
		}

		d.LoginURL, _ = user.LoginURL(c, "/dl")
		if user.Current(c) != nil {
			d.LoginURL, _ = user.LogoutURL(c, "/dl")
		}

		item := &memcache.Item{Key: cacheKey, Object: &d, Expiration: cacheDuration}
		if err := memcache.Gob.Set(c, item); err != nil {
			log.Errorf(c, "cache set error: %v", err)
		}
	}
	if err := listTemplate.ExecuteTemplate(w, "root", d); err != nil {
		log.Errorf(c, "error executing template: %v", err)
	}
}

func filesToFeatured(fs []File) (featured []Feature) {
	for _, feature := range featuredFiles {
		for _, file := range fs {
			if feature.fileRE.MatchString(file.Filename) {
				feature.File = file
				featured = append(featured, feature)
				break
			}
		}
	}
	return
}

func filesToReleases(fs []File) (stable, unstable []Release) {
	sort.Sort(fileOrder(fs))

	var r *Release
	var stableMaj, stableMin int
	add := func() {
		if r == nil {
			return
		}
		if r.Stable {
			if len(stable) == 0 {
				// Display files for latest stable release.
				stableMaj, stableMin, _ = parseVersion(r.Version)
				r.Visible = len(stable) == 0
			}
			stable = append(stable, *r)
			return
		}
		if len(unstable) != 0 {
			// Only show one (latest) unstable version.
			return
		}
		maj, min, _ := parseVersion(r.Version)
		if maj < stableMaj || maj == stableMaj && min <= stableMin {
			// Display unstable version only if newer than the
			// latest stable release.
			return
		}
		r.Visible = true
		unstable = append(unstable, *r)
	}
	for _, f := range fs {
		if r == nil || f.Version != r.Version {
			add()
			r = &Release{
				Version: f.Version,
				Stable:  isStable(f.Version),
			}
		}
		r.Files = append(r.Files, f)
	}
	add()
	return
}

// isStable reports whether the version string v is a stable version.
func isStable(v string) bool {
	return !strings.Contains(v, "beta") && !strings.Contains(v, "rc")
}

type fileOrder []File

func (s fileOrder) Len() int      { return len(s) }
func (s fileOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s fileOrder) Less(i, j int) bool {
	a, b := s[i], s[j]
	if av, bv := a.Version, b.Version; av != bv {
		return versionLess(av, bv)
	}
	if a.OS != b.OS {
		return a.OS < b.OS
	}
	if a.Arch != b.Arch {
		return a.Arch < b.Arch
	}
	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}
	return a.Filename < b.Filename
}

func versionLess(a, b string) bool {
	// Put stable releases first.
	if isStable(a) != isStable(b) {
		return isStable(a)
	}
	maja, mina, ta := parseVersion(a)
	majb, minb, tb := parseVersion(b)
	if maja == majb {
		if mina == minb {
			return ta >= tb
		}
		return mina >= minb
	}
	return maja >= majb
}

func parseVersion(v string) (maj, min int, tail string) {
	if i := strings.Index(v, "beta"); i > 0 {
		tail = v[i:]
		v = v[:i]
	}
	if i := strings.Index(v, "rc"); i > 0 {
		tail = v[i:]
		v = v[:i]
	}
	p := strings.Split(strings.TrimPrefix(v, "go1."), ".")
	maj, _ = strconv.Atoi(p[0])
	if len(p) < 2 {
		return
	}
	min, _ = strconv.Atoi(p[1])
	return
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	c := appengine.NewContext(r)

	// Authenticate using a user token (same as gomote).
	user := r.FormValue("user")
	if !validUser(user) {
		http.Error(w, "bad user", http.StatusForbidden)
		return
	}
	if builderKey == "" {
		http.Error(w, "no builder-key found in project metadata", http.StatusInternalServerError)
		return
	}
	if r.FormValue("key") != userKey(c, user) {
		http.Error(w, "bad key", http.StatusForbidden)
		return
	}

	var f File
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		log.Errorf(c, "error decoding upload JSON: %v", err)
		http.Error(w, "Something broke", http.StatusInternalServerError)
		return
	}
	if f.Filename == "" {
		http.Error(w, "Must provide Filename", http.StatusBadRequest)
		return
	}
	if f.Uploaded.IsZero() {
		f.Uploaded = time.Now()
	}
	k := datastore.NewKey(c, "File", f.Filename, 0, rootKey(c))
	if _, err := datastore.Put(c, k, &f); err != nil {
		log.Errorf(c, "putting File entity: %v", err)
		http.Error(w, "could not put File entity", http.StatusInternalServerError)
		return
	}
	if err := memcache.Delete(c, cacheKey); err != nil {
		log.Errorf(c, "cache delete error: %v", err)
	}
	io.WriteString(w, "OK")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/dl/")
	if name == "" {
		listHandler(w, r)
		return
	}
	if !fileRe.MatchString(name) {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, gcsBaseURL+name, http.StatusFound)
}

func validUser(user string) bool {
	switch user {
	case "adg", "bradfitz", "cbro":
		return true
	}
	return false
}

func userKey(c context.Context, user string) string {
	h := hmac.New(md5.New, []byte(builderKey))
	h.Write([]byte("user-" + user))
	return fmt.Sprintf("%x", h.Sum(nil))
}

var fileRe = regexp.MustCompile(`^go[0-9a-z.]+\.[0-9a-z.-]+\.(tar\.gz|pkg|msi|zip)$`)

func initHandler(w http.ResponseWriter, r *http.Request) {
	var fileRoot struct {
		Root string
	}
	c := appengine.NewContext(r)
	k := rootKey(c)
	err := datastore.RunInTransaction(c, func(c context.Context) error {
		err := datastore.Get(c, k, &fileRoot)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		_, err = datastore.Put(c, k, &fileRoot)
		return err
	}, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	io.WriteString(w, "OK")
}

// rootKey is the ancestor of all File entities.
func rootKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "FileRoot", "root", 0, nil)
}

// pretty returns a human-readable version of the given OS, Arch, or Kind.
func pretty(s string) string {
	t, ok := prettyStrings[s]
	if !ok {
		return s
	}
	return t
}

var prettyStrings = map[string]string{
	"darwin":  "OS X",
	"freebsd": "FreeBSD",
	"linux":   "Linux",
	"windows": "Windows",

	"386":   "32-bit",
	"amd64": "64-bit",

	"archive":   "Archive",
	"installer": "Installer",
	"source":    "Source",
}
