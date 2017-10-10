package controllers

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Machiel/slugify"
)

// slugFile :
func slugFile(file string) (name, serieName, serieNumberReturn string, year int) {
	var str []string
	ext := filepath.Ext(file)
	file = strings.ToLower(file[:len(file)-len(ext)])
	file = slugify.Slugify(file)

	video := regexp.MustCompile(`^dvdrip$|^dvd-r$|^bluray$|^bdrip$|^brrip$|^cam$|^ts$|^tc$|^vcd$|^md$|^ld$|^r[0-9]{1}$|^xvid$|^divx$|^scr$|^dvdscr$|^repack$|^hdlight$`)
	yearReg := regexp.MustCompile(`^[0-9]{4}$`)

	serie := regexp.MustCompile(`(?mi)[s]{1}[0-9]{1,2}[e]{1}[0-9]{1,3}`)
	serieEpisode := regexp.MustCompile(`episode`)
	serieNumber := regexp.MustCompile(`[0-9]{1,3}`)

	str = strings.Split(file, "-")

	var ifSerie = false

	if serie.MatchString(file) {
		ifSerie = true
	}

	var oldValue string
	for k, v := range str {
		// v = helpers.RemoveDateInName(v)
		v = removeLangInName(v)
		if k > 0 {
			name += " "
		}

		// Si on trouve l'année avant la saison et son épisode on passe à l'élement suivant
		var itsSerie bool
		if ifSerie && yearReg.MatchString(v) {
			if k < len(str)-1 {
				itsSerie = true
			}
		} else {
			name += v
		}

		if yearReg.MatchString(v) && !itsSerie {
			year, _ = strconv.Atoi(v)
			name = name[:len(name)-len(v)]
			break
		} else if yearReg.MatchString(v) && itsSerie {
			name += "-" + str[k+1]
			serieName = name[:len(name)-len(strings.Join(serie.FindAllString(name, -1), ""))] + v
			serieNumberReturn = serie.FindAllString(file, -1)[0]
			year, _ = strconv.Atoi(v)
			break
		} else if serie.MatchString(v) {
			serieName = name[:len(name)-len(strings.Join(serie.FindAllString(name, -1), ""))]
			serieNumberReturn = serie.FindAllString(file, -1)[0]
			break
		} else if serieEpisode.MatchString(oldValue) && serieNumber.MatchString(v) {
			serieName = name[:len(name)-len(oldValue)-len(v)-1]
			serieNumberReturn = serieNumber.FindAllString(file, -1)[0]
			break
		} else if video.MatchString(v) {
			break
		}
		oldValue = v
	}
	return slugify.Slugify(name) + ext, slugify.Slugify(serieName), serieNumberReturn, year
}

func slugSerieSeasonEpisode(serieNumber string) (season, episode int) {
	serie := regexp.MustCompile(`[s]{1}[0-9]{1,2}[e]{1}[0-9]{1,2}`)
	seasonNumber := regexp.MustCompile(`[s]{1}[0-9]{1,2}`)
	episodeNumber := regexp.MustCompile(`[e]{1}[0-9]{1,2}`)
	episodeNumber2 := regexp.MustCompile(`[0-9]{1,3}`)
	if serie.MatchString(serieNumber) {
		season, _ = strconv.Atoi(strings.Join(seasonNumber.FindAllString(serieNumber, -1), "")[1:])
		episode, _ = strconv.Atoi(strings.Join(episodeNumber.FindAllString(serieNumber, -1), "")[1:])
		return season, episode
	} else if episodeNumber2.MatchString(serieNumber) {
		episode, _ = strconv.Atoi(strings.Join(episodeNumber2.FindAllString(serieNumber, -1), ""))
		return 1, episode

	}
	return 1, 0
}

// removeLangInName :
func removeLangInName(s string) string {
	reg := regexp.MustCompile(`^fr$|^en$|^ru$|^us$`)
	if reg.MatchString(s) {
		return ""
	}
	return s
}
