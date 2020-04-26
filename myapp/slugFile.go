package myapp

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Machiel/slugify"
)

func (m *myFile) slugFile() {
	ext := filepath.Ext(m.complete)
	m.name = strings.ToLower(m.complete[:len(m.complete)-len(ext)])
	m.name = slugify.Slugify(m.name)

	var err error

	video := regexp.MustCompile(`(?mi)-(french|dvdrip|multi|vostfr|dvd-r|bluray|bdrip|brrip|cam|ts|tc|vcd|md|ld|r[0-9]|xvid|divx|scr|dvdscr|repack|hdlight|720p|480p|1080p|2160p|uhd)`)
	yearReg := regexp.MustCompile(`(?mi)-[0-9]{4}`)

	serie := regexp.MustCompile(`(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,2}))|(?:e\d{1,2})|(episode-(\d{2,3})-)`)
	match := serie.FindAllStringSubmatch(m.name, -1)

	if len(match) > 0 {
		for _, v := range match {
			if v[2] == "" && v[4] == "" {
				m.episode = formatSaisonNumberOuEpisode(v[0], 'e')
			} else if v[5] != "" {
				m.season = "s00"
				m.episode = "e" + v[5]
			} else {
				m.season = formatSaisonNumberOuEpisode(v[2], 's')
				m.episode = formatSaisonNumberOuEpisode(v[3], 'e')
			}
			break
		}
		m.serieNumber = fmt.Sprintf("%s%s", m.season, m.episode)
		if len(serie.FindStringIndex(m.name)) > 0 {
			m.serieName = m.complete[:serie.FindStringIndex(m.name)[0]-1]
			if len(yearReg.FindStringIndex(m.serieName)) > 0 {
				m.year, err = strconv.Atoi(yearReg.FindString(m.serieName)[1:])
				if err != nil {
					log.Println(err)
				}
				m.name = m.serieName[:yearReg.FindStringIndex(m.serieName)[0]]
				m.name = m.name + "-" + m.serieNumber
			} else {
				m.name = m.serieName + "-" + m.serieNumber
			}
		}
	} else {
		if len(video.FindStringIndex(m.name)) > 0 {
			m.name = m.name[:video.FindStringIndex(m.name)[0]]
			if len(yearReg.FindStringIndex(m.name)) > 0 {
				m.year, err = strconv.Atoi(yearReg.FindString(m.name)[1:])
				if err != nil {
					log.Println(err)
				}
				if len(yearReg.FindStringIndex(m.name)) > 0 {
					m.name = m.name[:yearReg.FindStringIndex(m.name)[0]]
				}
			}
		} else {
			if len(yearReg.FindStringIndex(m.name)) > 0 {
				m.year, err = strconv.Atoi(yearReg.FindString(m.name)[1:])
				if err != nil {
					log.Println(err)
				}
				m.name = m.name[:yearReg.FindStringIndex(m.name)[0]]
			}
			//else {
			//	m.complete = m.c
			//}
		}
	}
	m.complete = slugify.Slugify(m.name) + ext
	m.serieName = slugify.Slugify(m.serieName)
}

func formatSaisonNumberOuEpisode(num string, seasonOrEpisode byte) string {
	if len(num) == 2 && (num[0] == 's' || num[0] == 'e') {
		num = fmt.Sprintf("%s0%s", string(num[0]), string(num[1]))
	} else if num[0] != seasonOrEpisode {
		num = fmt.Sprintf("%s%s", string(seasonOrEpisode), num)
	}
	return num
}

// slugFile :
//func slugFile(file string) (name, serieName, serieNumberReturn string, year int) {
//	var str []string
//	ext := filepath.Ext(file)
//	file = strings.ToLower(file[:len(file)-len(ext)])
//	file = slugify.Slugify(file)
//
//	video := regexp.MustCompile(
//		`^french$|^dvdrip$|^multi$|^vostfr$|
//				^dvd-r$|^bluray$|^bdrip$|
//				^brrip$|^cam$|^ts$|^tc$|^vcd$|
//				^md$|^ld$|^r[0-9]$|
//				^xvid$|^divx$|^scr$|
//				^dvdscr$|^repack$|
//				^multi$|^hdlight$|^720p$|^480p$|^1080p$|^2160p$|^uhd$`)
//	yearReg := regexp.MustCompile(`^[0-9]{4}$`)
//
//	serie := regexp.MustCompile(`(?mi)[s][0-9]{1,2}[e][0-9]{1,3}`)
//	serieEpisode := regexp.MustCompile(`episode`)
//	serieNumber := regexp.MustCompile(`^[0-9]{2,3}$`)
//
//	str = strings.Split(file, "-")
//
//	var ifSerie = false
//
//	if serie.MatchString(file) {
//		ifSerie = true
//	}
//
//	//if _, s, _ := slugSerieSeasonEpisode(file); s != 0{
//	//	ifSerie = true
//	//}
//
//	var oldValue string
//	var oldName string
//
//	for k, v := range str {
//		// v = helpers.RemoveDateInName(v)
//		v = removeLangInName(v)
//		if k > 0 {
//			name += " "
//		}
//
//		// Si on trouve l'année avant la saison et son épisode on passe à l'élement suivant
//		var itsSerie bool
//		if ifSerie && yearReg.MatchString(v) {
//			if k < len(str)-1 {
//				itsSerie = true
//			}
//		} else {
//			name += v
//		}
//
//		if yearReg.MatchString(v) && !itsSerie {
//			year, _ = strconv.Atoi(v)
//			name = name[:len(name)-len(v)]
//			break
//		} else if yearReg.MatchString(v) && itsSerie {
//			name += "-" + str[k+1]
//			serieName = name[:len(name)-len(strings.Join(serie.FindAllString(name, -1), ""))] + v
//			serieNumberReturn, _, _ = slugSerieSeasonEpisode(serie.FindAllString(file, -1)[0])
//			year, _ = strconv.Atoi(v)
//			name = name[:len(name)-len(serie.FindAllString(file, -1)[0])] + serieNumberReturn
//			break
//		} else if serie.MatchString(v) && itsSerie{
//			serieName = name[:len(name)-len(strings.Join(serie.FindAllString(name, -1), ""))]
//			serieNumberReturn, _, _ = slugSerieSeasonEpisode(v)
//			name = serieName + " " + serieNumberReturn
//			break
//		} else if serieEpisode.MatchString(oldValue) && serieNumber.MatchString(v) && itsSerie{
//			serieName = name[:len(name)-len(oldValue)-len(v)-1]
//			serieNumberReturn, _, _ = slugSerieSeasonEpisode(v)
//			name = serieName + " " + serieNumberReturn
//			break
//		} else if serieNumber.MatchString(v) && itsSerie {
//			var toto = false
//			for _, value := range str {
//				if serie.MatchString(value) {
//					toto = true
//					break
//				}
//			}
//			if toto {
//				continue
//			}
//			serieNumberReturn, _, _ = slugSerieSeasonEpisode(v)
//			serieName = name[:len(name)-len(v)-1]
//			name = serieName + " " + serieNumberReturn
//			break
//		} else if video.MatchString(v) {
//			name = oldName
//			year, _ = searchYear(str, yearReg)
//			break
//		}
//		oldName = name
//		oldValue = v
//	}
//
//	return slugify.Slugify(name) + ext, slugify.Slugify(serieName), serieNumberReturn, year
//}

//func searchYear(str []string, reg *regexp.Regexp) (int, error) {
//	for _, val := range str {
//		if reg.MatchString(val) {
//			return strconv.Atoi(val)
//		}
//	}
//	return 0, nil
//}

//func (m *myFile) slugSerieSeasonEpisode() {
//	serie := regexp.MustCompile(`[s][0-9]{1,2}[e][0-9]{1,2}`)
//	seasonNumber := regexp.MustCompile(`[s][0-9]{1,2}`)
//	episodeNumber := regexp.MustCompile(`[e][0-9]{1,2}`)
//	episodeNumber2 := regexp.MustCompile(`[0-9]{2,3}`)
//	log.Println("serienumber : ", m.serieNumber)
//	if serie.MatchString(m.serieNumber) {
//		m.season, _ = strconv.Atoi(strings.Join(seasonNumber.FindAllString(m.serieNumber, -1), "")[1:])
//		m.episode, _ = strconv.Atoi(strings.Join(episodeNumber.FindAllString(m.serieNumber, -1), "")[1:])
//		m.serieNumber = checkIfTwoNumberToSeasonOrEpisode(m.season, m.episode)
//		return
//	} else if episodeNumber2.MatchString(m.serieNumber) {
//		m.season = 0
//		m.episode, _ = strconv.Atoi(m.serieNumber)
//		m.serieNumber = checkIfTwoNumberToSeasonOrEpisode(m.season, m.episode)
//		return
//	}
//	m.serieNumber = ""
//	m.season = 0
//	m.episode = 0
//}

// removeLangInName :
//func removeLangInName(s string) string {
//	reg := regexp.MustCompile(`^fr$|^en$|^ru$|^us$`)
//	if reg.MatchString(s) {
//		return ""
//	}
//	return s
//}

func checkIfTwoNumberToSeasonOrEpisode(season, episode int) string {
	strSeason := strconv.Itoa(season)
	strEpisode := strconv.Itoa(episode)
	var str string
	if len(strSeason) == 1 {
		str = "s0" + strSeason
	} else {
		str = "s" + strSeason
	}
	if len(strEpisode) == 1 {
		str = str + "e0" + strEpisode
	} else {
		str = str + "e" + strEpisode
	}
	return str
}
