package myapp

import (
	"fmt"
	"github.com/Machiel/slugify"
	"log"
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
		m.name = m.name[:cleanName[1]]
	}

	serie := regexp.MustCompile(`(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,3}))|(?:e\d{1,2})|(episode-(\d{2,3})-?)|((\d{1,2})-(\d{1,2})$)|((saison|season)-(\d{1,2})-episode-(\d{1,2}))`)
	match := serie.FindAllStringSubmatch(m.name, -1)

	m.extractResolution()

	// len(match) > 0 => serie
	// len(match) == 0 => film
	if len(match) > 0 {
		if len(match) > 1 {
			match = match[len(match)-1:]
		}
		for _, v := range match {
			if v[7] != "" && v[8] != "" {
				m.season = formatSaisonNumberOuEpisode(v[7], 's')
				m.episode = formatSaisonNumberOuEpisode(v[8], 'e')
			} else if v[11] != "" && v[12] != "" {
				m.season = formatSaisonNumberOuEpisode(v[11], 's')
				m.episode = formatSaisonNumberOuEpisode(v[12], 'e')
			} else if v[5] != "" {
				m.season = ""
				m.episode = v[5]
			} else if v[2] != "" && v[3] != "" {
				m.season = formatSaisonNumberOuEpisode(v[2], 's')
				m.episode = formatSaisonNumberOuEpisode(v[3], 'e')
			} else {
				m.episode = formatSaisonNumberOuEpisode(v[0], 'e')
				m.season = "s00"
			}
			break
		}
		m.serieNumber = fmt.Sprintf("%s%s", m.season, m.episode)
		if len(serie.FindStringIndex(m.name)) > 0 {
			findIndex := serie.FindAllIndex([]byte(m.name), -1)
			m.serieName = slugify.Slugify(m.complete[:findIndex[len(findIndex)-1][0]-1])
			m.extractYear(m.serieName)
			m.name = m.serieName + "-" + m.serieNumber
		}
	} else {
		if len(video.FindStringIndex(m.name)) > 0 {
			m.name = m.name[:video.FindStringIndex(m.name)[1]]
		}
		m.extractYear(m.completeSlug)

	}
	m.formatageFinal()
}

func (m *myFile) formatageFinal() {
	m.complete = ""

	formatFile := strings.Split(constants.FormatFile(), ",")
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
		}
		if len(m.complete) > 0 && !strings.HasSuffix(m.complete, "-") {
			m.complete += "-"
		}
	}

	if strings.HasSuffix(m.complete, "-") {
		m.complete = m.complete[:len(m.complete)-1]
	}

	separator := strings.TrimSpace(formatFile[0])
	if separator != "-" {
		if separator == "s" {
			separator = " "
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
	}
}

func (m *myFile) extractYear(str string) {
	yearReg := regexp.MustCompile(`(?mi)(19|20)\d{2}`)
	startAndEnd := yearReg.FindStringIndex(str)
	if startAndEnd != nil && len(startAndEnd) > 0 {

		m.year, err = strconv.Atoi(m.name[startAndEnd[0]:startAndEnd[1]])
		if err != nil {
			logger.L(logger.Red, "%s", err)
		}
		if m.serieName != "" {
			m.serieName = m.serieName[:yearReg.FindStringIndex(m.name)[0]]
			m.name = m.serieName
			log.Println(m.serieName, m.name)
		} else {
			m.name = m.name[:yearReg.FindStringIndex(m.name)[0]]
		}
		if strings.HasSuffix(m.name, "-") {
			m.name = m.name[:len(m.name)-1]
		}
		if strings.HasSuffix(m.serieName, "-") {
			m.serieName = m.serieName[:len(m.serieName)-1]
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
