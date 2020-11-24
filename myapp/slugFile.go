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

	cleanName := video.FindStringIndex(m.name)
	if len(cleanName) > 0 {
		m.name = m.name[:cleanName[0]]
	}

	serie := regexp.MustCompile(`(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,2}))|(?:e\d{1,2})|(episode-(\d{2,3})-)|((\d{1,2})-(\d{1,2})$)|((saison|season)-(\d{1,2})-episode-(\d{1,2}))`)
	match := serie.FindAllStringSubmatch(m.name, -1)

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
				m.season = "s00"
				m.episode = "e" + v[5]
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
		}
	}
	m.complete = slugify.Slugify(m.name) + ext
	m.serieName = slugify.Slugify(m.serieName)
}

func formatSaisonNumberOuEpisode(num string, seasonOrEpisode byte) string {
	if len(num) == 2 && (num[0] == 's' || num[0] == 'e') {
		num = fmt.Sprintf("%s0%s", string(num[0]), string(num[1]))
	} else if len(num) == 1 {
		num = fmt.Sprintf("%s0%s", string(seasonOrEpisode), string(num[0]))
	} else if num[0] != seasonOrEpisode {
		num = fmt.Sprintf("%s%s", string(seasonOrEpisode), num)
	}
	return num
}
