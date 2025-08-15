package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// RobotRule represents a rule in the robots.txt file
type RobotsRules struct {
	Host       string   // Hostname of the site
	Disallow   []string // List of disallowed paths
	CrawlDelay float64  // Optional crawl delay
}

// gets the robots.txt file for a given domain as a string and returns a populated struct containing the rules
func robots(urlStr string) RobotsRules {
	var rules RobotsRules
	rules.CrawlDelay = 0.1 // Default crawl delay
	robotsData := getData(urlStr, &rules)
	if robotsData == nil {
		fmt.Println("failed to get data")
		return RobotsRules{}
	}
	robotsData = append(robotsData, '\n')

	// parse the robots.txt file
	scanner := bytes.NewBuffer(robotsData)
	userAgentRegex := regexp.MustCompile(`User-agent: (\S+)`)
	disallowRegex := regexp.MustCompile(`Disallow: (\S*)`)
	crawlDelayRegex := regexp.MustCompile(`Crawl-delay: (\S+)`)
	userAgentAll := false // Track if we're in the "User-agent: *" section

	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			break // EOF reached
		}
		line = strings.TrimSpace(line) // Remove extra spaces
		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle "User-agent"
		if match := userAgentRegex.FindStringSubmatch(line); match != nil {
			userAgentAll = (match[1] == "*") // Only process rules under "User-agent: *"
			continue
		}

		// If we're in the "User-agent: *" section, process rules
		if userAgentAll {
			switch {
			case disallowRegex.MatchString(line):
				rules.Disallow = append(rules.Disallow, disallowRegex.FindStringSubmatch(line)[1])
			case crawlDelayRegex.MatchString(line):
				if delay, err := strconv.ParseFloat(crawlDelayRegex.FindStringSubmatch(line)[1], 64); err == nil {
					rules.CrawlDelay = delay
				} else {
					fmt.Println("Error parsing crawl delay:", err)
				}
			}
		}
	}
	fmt.Println("Rules: ", rules)
	return rules
}

// helper method
func getData(urlStr string, rules *RobotsRules) []byte {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL: ", err)
		return nil
	}

	rules.Host = urlObj.Host
	resp, err := http.Get(urlObj.Scheme + "://" + urlObj.Host + "/robots.txt")

	if err != nil {
		fmt.Println("Error getting data: ", err)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return nil
	}
	return body
}

func isAllowed(urlStr string, rules RobotsRules) bool {
	path := strings.TrimPrefix(urlStr, "http://"+rules.Host+"/")
	for _, disallow := range rules.Disallow {
		pattern := "^" + strings.ReplaceAll(disallow, `\*`, `.*`) + "$"

		matched, err := regexp.MatchString(pattern, path)
		if err != nil {
			fmt.Println("Regex error:", err)
			continue
		}
		if matched {
			return false
		}
	}
	return true
}
