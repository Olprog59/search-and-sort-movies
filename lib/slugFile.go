package lib

import (
	"bytes"
	"context"
	"errors"
	"github.com/Machiel/slugify"
	"github.com/sam-docker/media-organizer/logger"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (m *myFile) slugFile() error {
	m.ext = filepath.Ext(m.complete)
	m.name = strings.ToLower(m.complete[:len(m.complete)-len(m.ext)])

	// remove : [Kaerizaki-Fansub] - from : [Kaerizaki-Fansub] One Piece 1093 VOSTFR FHD (1920x1080) .mp4
	m.removeFirstBrackets()

	m.name = slugify.Slugify(m.name)
	m.completeSlug = m.name

	video := regexp.MustCompile(`(?mi)-(french|vf|dvdrip|multi|vostfr|subfrench|dvd-r|bluray|bdrip|brrip|cam|ts|tc|vcd|md|ld|r[0-9]|xvid|divx|scr|dvdscr|repack|hdlight|720p|480p|1080p|2160p|uhd|4k|1920x1080)`)
	language := regexp.MustCompile(`(?mi)-(french|multi|vostfr|subfrench|vo)`)

	cleanName := video.FindStringIndex(m.name)
	if len(cleanName) > 0 {
		cleanLanguage := language.FindStringIndex(m.name)
		if len(cleanLanguage) > 0 {
			m.language = m.name[cleanLanguage[0]+1 : cleanLanguage[1]]
		}
		m.name = strings.Replace(m.name, m.language, "", -1)
	}

	m.extractYear(m.name)
	m.extractResolution()

	serie := regexp.MustCompile(`(?mi)((s\d{1,2})(?:\W+)?(e?\d{1,4}))|e\d{1,4}|(episode-(\d{2,4})-?)|((\d{1,2})-(\d{2,4}))|((saison|season)-(\d{1,2})-episode-(\d{1,4}))|((\d{1,2})x(\d{1,4}))`)
	match := serie.FindAllStringSubmatch(m.name, -1)

	// len(match) > 0 => serie
	// len(match) == 0 => film
	if len(match) > 0 {
		for _, v := range match {
			if v[7] != "" && v[8] != "" {
				m.season = formatSaisonNumberOuEpisode(v[7])
				m.episode = formatSaisonNumberOuEpisode(v[8])
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else if v[11] != "" && v[12] != "" {
				m.season = formatSaisonNumberOuEpisode(v[11])
				m.episode = formatSaisonNumberOuEpisode(v[12])
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else if v[5] != "" {
				m.season = 0
				m.episode, _ = strconv.Atoi(v[5])
				m.name = strings.Replace(m.name, "episode-"+v[5], "", -1)
			} else if v[2] != "" && v[3] != "" {
				m.season = formatSaisonNumberOuEpisode(v[2])
				m.episode = formatSaisonNumberOuEpisode(v[3])
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else if v[14] != "" && v[15] != "" {
				m.season = formatSaisonNumberOuEpisode(v[14])
				m.episode = formatSaisonNumberOuEpisode(v[15])
				m.name = strings.Replace(m.name, v[0], "", -1)
			} else {
				m.episode = formatSaisonNumberOuEpisode(v[0])
				m.season = 0
				m.name = strings.Replace(m.name, v[0], "-", -1)
			}
			break
		}

		more := regexp.MustCompile(`(?mi)-{2,}`)
		places := more.FindStringIndex(m.name)
		if len(places) > 0 {
			m.serieName = m.name[:places[0]]
		}
		if m.serieName == "" && strings.HasSuffix(m.name, "-") {
			m.serieName = m.name[:len(m.name)-1]
		}
		m.name = m.serieName

	} else {
		more := regexp.MustCompile(`(?mi)-{2,}`)
		places := more.FindStringIndex(m.name)
		if len(places) > 0 {
			m.name = m.name[:places[0]]
		}
	}
	err := m.formatageFinal()
	if err != nil {
		return err
	}
	return nil
}

func (m *myFile) formatageFinal() error {
	m.complete = ""
	if m.file == "" {
		if m.serieName != "" {
			m.formatageSerie()
		} else {
			m.formatageMovie()
		}
		return nil
	}
	duration, err := getDurationMovie(m.file)
	if err != nil {
		logger.L(logger.Red, "%s", err)
		return err
	}
	sec := duration / 60
	if m.episode > 0 && sec < 70 {
		m.formatageSerie()
	} else if m.episode > 0 && sec >= 70 {
		return errors.New("inconsistency between file name and duration")
	} else if m.episode == 0 {
		if sec > 60 {
			m.formatageMovie()
		} else {
			// TODO : formatage si il y a une incohérence entre le nom du fichier et la durée
			return errors.New("inconsistency between file name and duration")
		}
	}
	return nil
}

func getDurationMovie(fileName string) (float64, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	var outputBuf bytes.Buffer
	var errorBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", fileName)
	cmd.Stdout = &outputBuf
	cmd.Stderr = &errorBuf
	//logger.L(logger.Yellow, "cmd : %v", cmd.String())
	err := cmd.Run()
	if err != nil {
		logger.L(logger.Red, "error : %s", errorBuf.String())
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(outputBuf.String()), 64)
}

func (m *myFile) extractResolution() {
	resolution := regexp.MustCompile(`(?mi)(720p|480p|1080p|2160p|uhd|4k|1920x1080)`)
	tabResolution := resolution.FindStringIndex(m.completeSlug)
	if len(tabResolution) > 0 {
		isSeparator := func() int {
			if m.completeSlug[tabResolution[0]] == '-' {
				return tabResolution[0] + 1
			}
			return tabResolution[0]
		}
		m.resolution = m.completeSlug[isSeparator():tabResolution[1]]
		switch m.resolution {
		case "1920x1080":
			m.name = strings.Replace(m.name, m.resolution, "", -1)
			m.resolution = "1080p"
		}

	}
}

func (m *myFile) extractYear(str string) {
	yearReg := regexp.MustCompile(`(?mi)(\s|-)(19|20)\d{2}(\s|-|$)`)
	startAndEnd := yearReg.FindStringIndex(str)

	if startAndEnd != nil && len(startAndEnd) > 0 {
		year := str[startAndEnd[0]:startAndEnd[1]]
		// logger.L(logger.Yellow, "year : %s", year)
		year = strings.Replace(year, " ", "", -1)
		year = strings.Replace(year, "-", "", -1)
		m.year, err = strconv.Atoi(year)

		if err != nil {
			logger.L(logger.Red, "%s", err)
		}
		if len(yearReg.FindStringIndex(m.name)) > 0 {
			m.name = strings.Replace(m.name, year, "", -1)
			//m.name = m.name[:yearReg.FindStringIndex(m.name)[0]]
		}
		if strings.HasSuffix(m.name, "-") {
			m.name = m.name[:len(m.name)-1]
		}
	}
}

func formatSaisonNumberOuEpisode(num string) int {
	if num[0] == 's' || num[0] == 'e' {
		atoi, err := strconv.Atoi(num[1:])
		if err != nil {
			return 0
		}
		return atoi
	}
	atoi, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return atoi
}
