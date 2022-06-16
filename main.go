package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type link struct {
	short string
	to    string
}

type App struct {
	links      []link
	httpIP     string
	httpPort   int
	httpSubdir string
	debug      bool
	testConfig bool
}

func fixuri(uri string) string {
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	return uri
}

func (a *App) getlinks(path string) error {
	f, err := os.Open(path)

	var e error

	if err != nil {
		return err
	}
	defer f.Close()

	var links []link

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() != "" {
			data := strings.Fields(s.Text())
			if len(data) != 2 {
				e = fmt.Errorf("invalid syntax %s", data)
			} else {
				l := link{short: data[0], to: data[1]}
				links = append(links, l)
			}
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	a.links = links
	return e
}

func (a App) httpGetLink(w http.ResponseWriter, r *http.Request) {
	uri := strings.Replace(r.RequestURI, a.httpSubdir, "", 1)
	log.WithFields(log.Fields{
		"uri": uri,
	}).Debugf("HTTP> got request")
	for _, l := range a.links {
		if l.short == uri {
			log.Debug("CORE> short found")
			http.Redirect(w, r, l.to, http.StatusFound)
			return
		}
	}
	log.Debug("CORE> short not found")
	w.WriteHeader(404)
	fmt.Fprintf(w, "404\n")
}

func main() {
	linksPath := flag.String("links", "links.txt", "Path to the links file")
	webPort := flag.Int("http-port", 8080, "HTTP Listen port")
	webIP := flag.String("http-host", "127.0.0.1", "HTTP Listen IP")
	webDir := flag.String("http-dir", "/", "If proxied in subfolder")
	debugState := flag.Bool("debug", false, "Enable debug logs")
	testConfig := flag.Bool("test", false, "Test links file for syntax")

	flag.Parse()

	app := App{
		debug:      *debugState,
		httpIP:     *webIP,
		httpPort:   *webPort,
		httpSubdir: fixuri(*webDir),
		testConfig: *testConfig,
	}

	if app.debug {
		log.SetLevel(log.DebugLevel)
	}

	err := app.getlinks(*linksPath)
	if err != nil {
		eFields := log.Fields{
			"error": err,
		}
		if !app.testConfig {
			log.WithFields(eFields).Error("CORE> Failed to load short link")
		} else {
			log.WithFields(eFields).Fatal("Failed to load short link")
		}
	}

	log.Infof("CORE> Loaded %v urls", len(app.links))

	log.Infof("HTTP> Listening on http://%s:%v", app.httpIP, app.httpPort)

	http.HandleFunc("/", app.httpGetLink)
	err = http.ListenAndServe(fmt.Sprintf("%s:%v", app.httpIP, app.httpPort), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("HTTP:")
	}
}
