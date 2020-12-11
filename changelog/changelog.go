package changelog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

type Changes struct {
	Tag        string
	Added      string
	Changed    string
	Deprecated string
	Fixed      string
	Removed    string
	Security   string
	Time       time.Time
}

func NewChanges(tag string) *Changes {
	return &Changes{
		Tag: tag,
	}
}

func (changes *Changes) WriteTo(w io.Writer) (int64, error) {

	buf := bytes.NewBufferString("")
	fmt.Fprintf(buf, "## [%s]", changes.Tag)
	if changes.Time.Unix() > 0 {
		fmt.Fprintf(buf, " - %s", changes.Time.Format("2006-01-02"))
	}
	buf.WriteString("\n")

	if changes.Added != "" {
		fmt.Fprintf(buf, "### Added\n%s\n\n", changes.Added)
	}

	if changes.Changed != "" {
		fmt.Fprintf(buf, "### Changed\n%s\n\n", changes.Changed)
	}
	if changes.Deprecated != "" {
		fmt.Fprintf(buf, "### Deprecated\n%s\n\n", changes.Deprecated)
	}
	if changes.Fixed != "" {
		fmt.Fprintf(buf, "### Fixed\n%s\n\n", changes.Fixed)
	}
	if changes.Removed != "" {
		fmt.Fprintf(buf, "### Removed\n%s\n\n", changes.Removed)
	}
	if changes.Security != "" {
		fmt.Fprintf(buf, "### Security\n%s\n\n", changes.Security)
	}

	return buf.WriteTo(w)
}

type Reference struct {
	Tag     string
	Raw     string
	From    string
	To      string
	BaseURL string
}

func (ref *Reference) WriteTo(w io.Writer) (int64, error) {
	if ref.To != "" && ref.From != "" {
		n, err := fmt.Fprintf(w, "[%s]: %s/compare/%s...%s\n", ref.Tag, ref.BaseURL, ref.From, ref.To)
		return int64(n), err
	}
	n, err := fmt.Fprintf(w, "%s\n", ref.Raw)
	return int64(n), err
}

type Contents struct {
	Header     string
	Changes    []*Changes
	Unreleased *Changes
	Last       *Changes
	Rest       string
	Refs       []Reference
}

func (c Contents) ChangeLogInfo(version string) string {
	s := ""
	for _, element := range c.Changes {
		if strings.Compare(element.Tag, version) == 0 {
			if element.Added != "" {
				s += fmt.Sprintln("Added")
				s += fmt.Sprintln(element.Added)
				s += fmt.Sprintln()
			}

			if element.Changed != "" {
				s += fmt.Sprintln("Changed")
				s += fmt.Sprintln(element.Changed)
				s += fmt.Sprintln()
			}

			if element.Deprecated != "" {
				s += fmt.Sprintln("Deprecated")
				s += fmt.Sprintln(element.Deprecated)
				s += fmt.Sprintln()
			}

			if element.Fixed != "" {
				s += fmt.Sprintln("Fixed")
				s += fmt.Sprintln(element.Fixed)
				s += fmt.Sprintln()
			}

			if element.Removed != "" {
				s += fmt.Sprintln("Removed")
				s += fmt.Sprintln(element.Removed)
				s += fmt.Sprintln()
			}

			if element.Security != "" {
				s += fmt.Sprintln("Security")
				s += fmt.Sprintln(element.Security)
				s += fmt.Sprintln()
			}
			return s
		}
	}
	return "Versin Not Found"
}

var unreleasedRe = regexp.MustCompile(`(?i)^##\s*\[?(unreleased)\]?\s*$`)
var sectionRe = regexp.MustCompile(`(?i)^###\s(added|changed|deprecated|fixed|removed|security)\s*$`)
var changeRe = regexp.MustCompile(`(?i)^##\s*\[?([0-9.]+)\]?\s*-?\s*([0-9\-]+)?\s*$`)
var changeRefRe = regexp.MustCompile(`(?i)^\[([^\]]+)\]:\s*(.*)/compare/(.*)\.\.\.(.*)$`)
var refRe = regexp.MustCompile(`(?i)^\[([^\]]+)\]:\s*(.*)$`)

func Parse(r io.Reader) (*Contents, error) {
	var contents Contents
	var section string

	header := bytes.NewBufferString("")
	rest := bytes.NewBufferString("")
	items := bytes.NewBufferString("")

	var changes *Changes

	finishSection := func() {
		if section != "" && items.Len() > 0 {
			switch strings.ToLower(section) {
			case "added":
				changes.Added = strings.Trim(items.String(), "\n ")
			case "changed":
				changes.Changed = strings.Trim(items.String(), "\n ")
			case "deprecated":
				changes.Deprecated = strings.Trim(items.String(), "\n ")
			case "fixed":
				changes.Fixed = strings.Trim(items.String(), "\n ")
			case "removed":
				changes.Removed = strings.Trim(items.String(), "\n ")
			case "security":
				changes.Security = strings.Trim(items.String(), "\n ")
			}
			items.Reset()
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		u := unreleasedRe.FindStringSubmatch(line)
		s := sectionRe.FindStringSubmatch(line)
		c := changeRe.FindStringSubmatch(line)
		rf := refRe.FindStringSubmatch(line)
		cr := changeRefRe.FindStringSubmatch(line)

		isHeaderUnreleased := len(u) > 0
		isHeaderChanges := len(c) > 0
		isHeaderSection := len(s) > 0
		isRef := (len(rf) > 0)
		isChangeRef := (len(cr) > 0)
		isHeader := isHeaderUnreleased || isHeaderChanges
		isHeaderOrSubheader := isHeaderUnreleased || isHeaderChanges || isHeaderSection
		isEndOfSection := isHeaderOrSubheader || isRef

		if changes == nil && isHeader {
			contents.Header = header.String()
		}

		if changes != nil && isEndOfSection {
			finishSection()
		}
		if changes != nil && isHeader {
			contents.Changes = append(contents.Changes, changes)
		}

		if changes != nil && isRef {
			contents.Changes = append(contents.Changes, changes)
			changes = nil
		}

		if isRef && !isChangeRef {
			ref := Reference{
				Tag: rf[1],
				Raw: rf[0],
			}
			contents.Refs = append(contents.Refs, ref)
		}

		if isChangeRef {
			ref := Reference{
				Tag:     cr[1],
				Raw:     cr[0],
				BaseURL: cr[2],
				From:    cr[3],
				To:      cr[4],
			}
			contents.Refs = append(contents.Refs, ref)
		}

		if isHeaderUnreleased {
			changes = NewChanges(u[1])
			contents.Unreleased = changes
		}

		if isHeaderChanges {
			changes = NewChanges(c[1])
			t, err := time.Parse("2006-01-02", c[2])
			if err != nil {
				return nil, err
			}
			changes.Time = t
		}

		if len(s) > 0 {
			section = s[1]
		}

		if section != "" && !isHeaderOrSubheader {
			items.WriteString(line)
			items.WriteString("\n")
		}

		if changes != nil && strings.ToLower(changes.Tag) == "unreleased" {
			continue
		}

		if contents.Header == "" && changes == nil {
			header.WriteString(line)
			header.WriteString("\n")
		} else if !isRef {
			rest.WriteString(line)
			rest.WriteString("\n")
		}
	}
	contents.Rest = rest.String()
	if changes != nil {
		finishSection()
	}

	if contents.Unreleased == nil {
		contents.Unreleased = NewChanges("Unreleased")
	}

	return &contents, nil
}

func (contents *Contents) WriteTo(w io.Writer) (int64, error) {
	buf := bytes.NewBufferString("")
	buf.WriteString(contents.Header)
	contents.Unreleased.WriteTo(buf)
	buf.WriteString(contents.Rest)
	for _, ref := range contents.Refs {
		ref.WriteTo(buf)
	}
	return buf.WriteTo(w)
}
