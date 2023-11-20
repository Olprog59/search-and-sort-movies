package myapp

import (
	"fmt"
	"github.com/Machiel/slugify"
	"path/filepath"
	"regexp"
	"search-and-sort-movies/myapp/constants"
	"search-and-sort-movies/myapp/logger"
	"strconv"
	"strings"
)

func (m *myFile) slugFile() {
	m.ext = filepath.Ext(m.complete)
	m.name = strings.ToLower(m.complete[:len(m.complete)-len(m.ext)])
	m.name = slugify.Slugify(m.name)
	m.completeSlug = m.name

	video := regexp.MustCompile(`(?mi)-(french|vf|dvdrip|multi|vostfr|dvd-r|bluray|bdrip|brrip|cam|ts|tc|vcd|md|ld|r[0-9]|xvid|divx|scr|dvdscr|repack|hdlight|720p|480p|1080p|2160p|uhd|4k)`)

	cleanName := video.FindStringIndex(m.name)
	if len(cleanName) > 0 {
		m.language = m.name[cleanName[0]+1 : cleanName[1]]
		m.name = strings.Replace(m.name, m.language, "", -1)
	}

	m.extractYear(m.name)
	m.extractResolution()

	serie := regexp.MustCompile(`(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,4}))|e\d{1,4}|(episode-(\d{2,4})-?)|((\d{1,2})-(\d{2,4}))|((saison|season)-(\d{1,2})-episode-(\d{1,4}))`)
	match := serie.FindAllStringSubmatch(m.name, -1)

	// len(match) > 0 => serie
	// len(match) == 0 => film
	if len(match) > 0 {
		for _, v := range match {
			if v[7] != "" && v[8] != "" {
				m.season = formatSaisonNumberOuEpisode(v[7], 's')
				m.episode = formatSaisonNumberOuEpisode(v[8], 'e')
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else if v[11] != "" && v[12] != "" {
				m.season = formatSaisonNumberOuEpisode(v[11], 's')
				m.episode = formatSaisonNumberOuEpisode(v[12], 'e')
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else if v[5] != "" {
				m.season = "s00"
				m.episode = "e" + v[5]
				m.name = strings.Replace(m.name, "episode-"+v[5], "", -1)
			} else if v[2] != "" && v[3] != "" {
				m.season = formatSaisonNumberOuEpisode(v[2], 's')
				m.episode = formatSaisonNumberOuEpisode(v[3], 'e')
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else {
				m.episode = formatSaisonNumberOuEpisode(v[0], 'e')
				m.season = "s00"
				m.name = strings.Replace(m.name, v[0], "-", -1)
			}
			break
		}
		m.serieNumber = fmt.Sprintf("%s%s", m.season, m.episode)
		more := regexp.MustCompile(`(?mi)-{2,}`)
		places := more.FindStringIndex(m.name)
		if len(places) > 0 {
			m.serieName = m.name[:places[0]]
		}
		m.name = m.serieName + "-" + m.serieNumber
		//}
	} else {
		more := regexp.MustCompile(`(?mi)-{2,}`)
		places := more.FindStringIndex(m.name)
		if len(places) > 0 {
			m.name = m.name[:places[0]]
		}
	}
	m.formatageFinal()
}

func (m *myFile) formatageFinal() {
	m.complete = ""

	formatFile := strings.Split(constants.FormatFile(), ",")
	isNameInFormatFile := false
	for _, v := range formatFile[1:] {
		v = strings.TrimSpace(v)
		if m.resolution != "" && v == "resolution" {
			m.complete += m.resolution
		}
		if m.year > 0 && v == "year" {
			m.complete += fmt.Sprintf("%d", m.year)
		}
		if v == "name" {
			m.complete += slugify.Slugify(m.name)
			isNameInFormatFile = true
		}
		if len(m.complete) > 0 && !strings.HasSuffix(m.complete, "-") {
			m.complete += "-"
		}
	}

	if !isNameInFormatFile {
		m.complete += slugify.Slugify(m.name)
	}

	if strings.HasSuffix(m.complete, "-") {
		m.complete = m.complete[:len(m.complete)-1]
	}

	separator := strings.TrimSpace(formatFile[0])
	if separator != "-" {
		if separator == "s" {
			separator = " "
		}
		if separator == "" {
			separator = "-"
		}
		m.complete = strings.ReplaceAll(m.complete, "-", separator)
	}
	fmt.Printf("%#v\n", m)
	m.complete += m.ext
	m.serieName = slugify.Slugify(m.serieName)
}

func (m *myFile) extractResolution() {
	resolution := regexp.MustCompile(`(?mi)(720p|480p|1080p|2160p|uhd|4k)`)
	tabResolution := resolution.FindStringIndex(m.completeSlug)
	if len(tabResolution) > 0 {
		isSeparator := func() int {
			if m.completeSlug[tabResolution[0]] == '-' {
				return tabResolution[0] + 1
			}
			return tabResolution[0]
		}
		m.resolution = m.completeSlug[isSeparator():tabResolution[1]]
		m.name = strings.Replace(m.name, m.resolution, "", -1)
	}
}

func (m *myFile) extractYear(str string) {
	yearReg := regexp.MustCompile(`(?mi)(19|20)\d{2}`)
	startAndEnd := yearReg.FindStringIndex(str)
	if startAndEnd != nil && len(startAndEnd) > 0 {
		m.year, err = strconv.Atoi(str[startAndEnd[0]:startAndEnd[1]])

		if err != nil {
			logger.L(logger.Red, "%s", err)
		}
		if len(yearReg.FindStringIndex(m.name)) > 0 {
			m.name = strings.Replace(m.name, str[startAndEnd[0]:startAndEnd[1]], "", -1)
			//m.name = m.name[:yearReg.FindStringIndex(m.name)[0]]
		}
		if strings.HasSuffix(m.name, "-") {
			m.name = m.name[:len(m.name)-1]
		}
	}
}

func formatSaisonNumberOuEpisode(num string, seasonOrEpisode byte) string {
	if len(num) == 2 && (num[0] == 's' || num[0] == 'e') {
		num = fmt.Sprintf("%s0%s", string(num[0]), string(num[1]))
	} else if len(num) == 4 && (num[0] == 's' || num[0] == 'e') {
		if string(num[1]) == "0" {
			num = fmt.Sprintf("%s%s", string(num[0]), string(num[2:]))
		} else {
			num = fmt.Sprintf("%s%s", string(num[0]), string(num[1:]))
		}
	} else if len(num) == 1 {
		num = fmt.Sprintf("%s0%s", string(seasonOrEpisode), string(num[0]))
	} else if num[0] != seasonOrEpisode {
		num = fmt.Sprintf("%s%s", string(seasonOrEpisode), num)
	}
	return num
}
