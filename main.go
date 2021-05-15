package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	_ "os"
	"os/exec"
	_ "reflect"
	"regexp"
	"strconv"
	"strings"
)

var TIME = 20

func main() {

	// generate csv file
	fmt.Print(fmt.Sprintf("GENERATE POWERTOP REPORT IN %v SECONDES\n",TIME))
	_, err := exec.Command("sudo", "powertop", "--csv",fmt.Sprintf("--time=%v",TIME)).Output()
	if err != nil {
		log.Fatal(err)
	}
	// filter powertop.csv
	// get only powerful electric headset
	outFilter, err := exec.Command("sed", "-n", "/Usage;Wakeups\\/s;GPU ops\\/s;/,/^$/p", "powertop.csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(string(outFilter)))
	r.Comma = ';'
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// split without header
	for _, line := range records[1:] {
		// if my value is empty
		if line[7] == "" {
			continue
		}
		//fmt.Println(line)
		fields := map[string]interface{}{
			"Usage(um/s)": findAndConvertPrefix(line[0])*math.Pow10(6),
			"Wakeups/s":   findAndConvertPrefix(line[1]),
			"GPU_ops/s": line[2],
			"Disk_Io/s": line[3],
			"GFX_wakeups/s": line[4],
			"Category": line[5],
			"Description": line[6],
			"Power(mW)": findAndConvertPrefix(line[7])*math.Pow10(3),
		}
		fmt.Println(fields)
	}
}

func findAndConvertPrefix(data string) float64 {
	reg_result := regexp.MustCompile(`([0-9]+\.?[0-9]*)[\s]*([a-zA-Z \/]+)?`).FindSubmatch([]byte(data))
	if len(reg_result) == 0 {
		return 0
	}
	val, err := strconv.ParseFloat(string(reg_result[1]), 64)
	if err != nil {
		fmt.Println(err)
	}
	if len(reg_result) > 2 {
		format := string(reg_result[2])
		if strings.ContainsAny(format, "n") {
			val *= math.Pow10(-9)
		}
		if strings.ContainsAny(format, "u") {
			val *= math.Pow10(-6)
		}
		if strings.ContainsAny(format, "m") {
			val *= math.Pow10(-3)
		}
		if strings.ContainsAny(format, "k") {
			val *= math.Pow10(3)
		}
	}
	return val
}
