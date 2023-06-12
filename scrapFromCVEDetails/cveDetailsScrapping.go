package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
	"github.com/tealeg/xlsx"
)

type Vulnerability struct {
	CVE                string
	CWE                string
	Summary            string
	Description        string
	Exploits           string
	VulnerabilityTypes string
	PublishDate        string
	UpdateDate         string
	Score              string
	GainedAccessLevel  string
	AccessComplexity   string
	Authentication     string
	Confidentiality    string
	Integrity          string
	Availability       string
}

func main() {
	// Base URL of the CVE details page
	baseURL := "https://www.cvedetails.com/vulnerability-list.php?vendor_id=0&product_id=0&version_id=0"

	// Prompt user for start and end page inputs
	var startPage, endPage int
	fmt.Print("Enter start page number: ")
	fmt.Scanln(&startPage)
	fmt.Print("Enter end page number: ")
	fmt.Scanln(&endPage)

	// Prompt user for export format
	var exportFormat string
	fmt.Print("Enter export format (csv or xlsx): ")
	fmt.Scanln(&exportFormat)

	// Create a slice to store the vulnerabilities
	vulnerabilities := []Vulnerability{}

	// Initialize the progress bar
	bar := pb.StartNew(endPage - startPage + 1)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop through the pages
	for page := startPage; page <= endPage; page++ {
		// Increment the wait group counter
		wg.Add(1)

		// Launch a goroutine to scrape the page
		go func(page int) {
			defer wg.Done()

			// Construct the URL with page query
			url := fmt.Sprintf("%s&page=%d", baseURL, page)

			// Fetch the HTML content from the web
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// Create a new document and load the HTML content
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Find the total number of vulnerabilities on the current page
			totalVulns := doc.Find("tr.srrowns").Length()

			// Find the rows with class 'srrowns'
			doc.Find("tr.srrowns").Each(func(i int, rowSelection *goquery.Selection) {
				// Create a new vulnerability object
				vuln := Vulnerability{}

				// Get the cells in the row
				rowSelection.Find("td").Each(func(j int, cellSelection *goquery.Selection) {
					cellText := strings.TrimSpace(cellSelection.Text())

					// Extract the relevant data based on the cell's index
					switch j {
					case 1:
						vuln.CVE = cellSelection.Find("a").Text()
					case 2:
						vuln.CWE = cellSelection.Find("a").Text()
					case 3:
						vuln.Summary = cellText
					case 4:
						vuln.Exploits = cellText
					case 5:
						vuln.VulnerabilityTypes = cellText
					case 6:
						vuln.PublishDate = cellText
					case 7:
						vuln.UpdateDate = cellText
					case 8:
						vuln.Score = cellText
					case 9:
						vuln.GainedAccessLevel = cellText
					case 10:
						vuln.AccessComplexity = cellText
					case 11:
						vuln.Authentication = cellText
					case 12:
						vuln.Confidentiality = cellText
					case 13:
						vuln.Integrity = cellText
					case 14:
						vuln.Availability = cellText
					}
				})

				// Find the next row and get the description
				description := rowSelection.Next().Find("td.cvesummarylong").Text()
				vuln.Description = strings.TrimSpace(description)

				// Append the vulnerability to the slice
				vulnerabilities = append(vulnerabilities, vuln)
			})

			// Increment the progress bar
			bar.Add(totalVulns)
		}(page)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Finish the progress bar
	bar.Finish()

	// Print the extracted data
	for _, vuln := range vulnerabilities {
		fmt.Printf("CVE: %s\n", vuln.CVE)
		fmt.Printf("CWE: %s\n", vuln.CWE)
		fmt.Printf("Summary: %s\n", vuln.Summary)
		fmt.Printf("Description: %s\n", vuln.Description)
		fmt.Printf("Exploits: %s\n", vuln.Exploits)
		fmt.Printf("Vulnerability Types: %s\n", vuln.VulnerabilityTypes)
		fmt.Printf("Publish Date: %s\n", vuln.PublishDate)
		fmt.Printf("Update Date: %s\n", vuln.UpdateDate)
		fmt.Printf("Score: %s\n", vuln.Score)
		fmt.Printf("Gained Access Level: %s\n", vuln.GainedAccessLevel)
		fmt.Printf("Access Complexity: %s\n", vuln.AccessComplexity)
		fmt.Printf("Authentication: %s\n", vuln.Authentication)
		fmt.Printf("Confidentiality Impact: %s\n", vuln.Confidentiality)
		fmt.Printf("Integrity Impact: %s\n", vuln.Integrity)
		fmt.Printf("Availability Impact: %s\n", vuln.Availability)
		fmt.Println("--------------")
	}

	var fileName string
	fmt.Print("Enter the file name (default: vulnerabilities.csv[.xlsx]): ")
	_, _ = fmt.Scanln(&fileName)

	// Prompt user for export format
	fmt.Print("Enter export format (csv or xlsx): ")
	_, _ = fmt.Scanln(&exportFormat)

	// Export the data based on the chosen format
	switch strings.ToLower(exportFormat) {
	case "xlsx", "excel", "x":
		if len(fileName) > 0 {
			exportCSV(vulnerabilities, fileName)
		} else {
			exportXLSX(vulnerabilities, "vulnerabilities.xlsx")
		}
	default:
		if len(fileName) > 0 {
			exportCSV(vulnerabilities, fileName)
		} else {
			exportCSV(vulnerabilities, "vulnerabilities.csv")
		}
	}
}

func exportCSV(vulnerabilities []Vulnerability, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the header
	file.WriteString("CVE,CWE,Summary,Description,Exploits,Vulnerability Types,Publish Date,Update Date,Score,Gained Access Level,Access Complexity,Authentication,Confidentiality Impact,Integrity Impact,Availability Impact\n")

	// Write each vulnerability to the file
	for _, vuln := range vulnerabilities {
		file.WriteString(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			vuln.CVE, vuln.CWE, vuln.Summary, vuln.Description, vuln.Exploits, vuln.VulnerabilityTypes,
			vuln.PublishDate, vuln.UpdateDate, vuln.Score, vuln.GainedAccessLevel, vuln.AccessComplexity,
			vuln.Authentication, vuln.Confidentiality, vuln.Integrity, vuln.Availability))
	}

	fmt.Printf("Data exported to %s\n", filename)
}

func exportXLSX(vulnerabilities []Vulnerability, filename string) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Vulnerabilities")
	if err != nil {
		log.Fatal(err)
	}

	// Write the header
	row := sheet.AddRow()
	row.AddCell().Value = "CVE"
	row.AddCell().Value = "CWE"
	row.AddCell().Value = "Summary"
	row.AddCell().Value = "Description"
	row.AddCell().Value = "Exploits"
	row.AddCell().Value = "Vulnerability Types"
	row.AddCell().Value = "Publish Date"
	row.AddCell().Value = "Update Date"
	row.AddCell().Value = "Score"
	row.AddCell().Value = "Gained Access Level"
	row.AddCell().Value = "Access Complexity"
	row.AddCell().Value = "Authentication"
	row.AddCell().Value = "Confidentiality Impact"
	row.AddCell().Value = "Integrity Impact"
	row.AddCell().Value = "Availability Impact"

	// Write each vulnerability to the sheet
	for _, vuln := range vulnerabilities {
		row = sheet.AddRow()
		row.AddCell().Value = vuln.CVE
		row.AddCell().Value = vuln.CWE
		row.AddCell().Value = vuln.Summary
		row.AddCell().Value = vuln.Description
		row.AddCell().Value = vuln.Exploits
		row.AddCell().Value = vuln.VulnerabilityTypes
		row.AddCell().Value = vuln.PublishDate
		row.AddCell().Value = vuln.UpdateDate
		row.AddCell().Value = vuln.Score
		row.AddCell().Value = vuln.GainedAccessLevel
		row.AddCell().Value = vuln.AccessComplexity
		row.AddCell().Value = vuln.Authentication
		row.AddCell().Value = vuln.Confidentiality
		row.AddCell().Value = vuln.Integrity
		row.AddCell().Value = vuln.Availability
	}

	err = file.Save(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Data exported to %s\n", filename)
}
