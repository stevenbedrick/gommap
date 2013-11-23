package main

import (
	"os/exec"
	"fmt"
	"strings"
	"bufio"
	"io"
	"path"
	"regexp"
	"errors"
	"net/http"
	"metamap"
	"metamap/mini" // metamap/mini
	"encoding/xml"
	"encoding/json"
	"log"
	"time"
)

const MetamapHomeDir = "/Users/steven/Desktop/Hacking/mm/mm11/public_mm"
const MetamapCmd = "bin/metamap11v2"
const MetamapArgs = "--XMLf --silent"

type MetamapInstance struct {
	TextInput chan string
	MappedOutput chan metamap.MMOs
	Control chan bool
}

func (m *MetamapInstance) Cleanup() {
	close(m.TextInput)
	close(m.MappedOutput)
	close(m.Control)	
}

// spawn a new MetaMap slave process; returns struct with 
// i/o and control channels
func SpawnMetamap() *MetamapInstance {
	in_channel := make(chan string)
	res_channel := make(chan metamap.MMOs)
	done_channel := make(chan bool)
	
	// start up external metamap process
	go handleMetamap(in_channel, res_channel, done_channel)
	
	return &MetamapInstance{in_channel, res_channel, done_channel}
	
}

func readToEOM(from *bufio.Reader, eom *regexp.Regexp) (string, error) {
	result := ""

	for {
		str, err := from.ReadString('\n')
		if err == io.EOF {
			fmt.Println("At eof, breaking")
			break
		} else if err != nil {
			fmt.Printf("Some sort of error reading a line: %s", err.Error())
			return "", errors.New("Couldn't read a line!")
		} else if eom.MatchString(str) {
			result += str
			// fmt.Println("At end of output")
			break
		}
		result += str
	}

	return result, nil
}



func handleMetamap(text_to_map chan string, mapped_mmos chan metamap.MMOs, done chan bool) {

	cmd := exec.Command(path.Join(MetamapHomeDir, MetamapCmd), MetamapArgs)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error with stdoutpipe!: %s", err.Error())
	}
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("Error with stdinpipe!: %s", err.Error())
	}
	
	buf_reader := bufio.NewReader(stdout)
	buf_writer := bufio.NewWriter(stdin)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error with starting!: %s", err.Error())
		return
	}

	// gotta get past the first line of boilerplate output...
	_, err = buf_reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading first line: %s\n", err.Error())
		return
	}
	
	// "prime the pump"
	buf_writer.WriteString("monkeys\n\n")
	buf_writer.Flush()
	
	// now there's another line ready to read, plus a bunch of crap
	_, err = buf_reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading second line: %s\n", err.Error())
		return
	}
	
	// how do we know when we're at the end of message?
	eom_regex, _ := regexp.Compile("</MMOs>")
	
	// before we can start looking for input to handle, we have to deal with the remnants of our "pump priming"
	result, err := readToEOM(buf_reader, eom_regex)
	
	// ok, now we can start looking for input

	for {
		select {
		case more_input := <-text_to_map:
			// fmt.Println("got input: ---->", more_input, "<-----")
			start_time := time.Now()
			buf_writer.WriteString(more_input + "\n\n")
			buf_writer.Flush()
			result, err = readToEOM(buf_reader, eom_regex)
			decoded := &metamap.MMOs{}
			xml.NewDecoder(strings.NewReader(result)).Decode(decoded)			
			decoded.ParseTime = time.Since(start_time)
			decoded.RawXML = result
			mapped_mmos <- *decoded
		case <-done:
			fmt.Println("done, time to kill")
			// clean up
			fmt.Println("Trying to kill...")
			err = cmd.Process.Kill()
			if err != nil {
				fmt.Printf("Error killing process: %s", err.Error())
			}
			fmt.Printf("Killed!")
			break
		}
	}

}





// handle web service calls- invoke metamap, etc.
func handler(w http.ResponseWriter, r *http.Request) {
	str_to_map := r.FormValue("str")
	str_to_map = strings.TrimSpace(str_to_map)
	
	// TODO: figure out a way to handle situations with multiple newlines in str_to_map; e.g. "sentence one.\n\n sentence two."
	// Right now, MetaMap is processing these as two separate inputs, and since we only try and read one output, we're having problems.
	
	if len(str_to_map) == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(w, "{}")
		return
	} else {
		fmt.Println("About to try and map: ", str_to_map)
		// get an instance from the pool
		mm_instance := <- instance_pool
		
		mm_instance.TextInput <- str_to_map
		temp_res := <- mm_instance.MappedOutput

		// put our instance back on the pool
		instance_pool <- mm_instance
		
		if r.FormValue("format") == "xml" {
			// just send the xml
			w.Header().Set("Content-Type", "text/xml")
			fmt.Fprintln(w, temp_res.RawXML)
		} else if r.FormValue("format") == "json" {
			// turn into smaller, more compact JSON representation
			smaller := mini.FromFullMMO(&temp_res)
			json_bytes, _ := json.Marshal(smaller)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(json_bytes))			
		} else {
			// for now, send XML
			w.Header().Set("Content-Type", "text/xml")
			fmt.Fprintln(w, temp_res.RawXML)
		}
		fmt.Println("Elapsed time: ", temp_res.ParseTime)
	}
}

// set up basic logging of requests
func Log(handler http.Handler) http.Handler { 
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL) 
		log.Printf("instance pool size: %d", len(instance_pool))
handler.ServeHTTP(w, r) 
    }) 
} 

var instance_pool chan MetamapInstance
const MAX_INSTANCES = 5

func main() {
	
	instance_pool = make(chan MetamapInstance, 5)
	for i := 0; i < MAX_INSTANCES; i++ {
		fmt.Println("Setting up instance", i)
		instance_pool <- *SpawnMetamap()
	}


	// we use StripPrefix so that /tmpfiles/somefile will access /tmp/somefile
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
	http.Handle("/ui/", http.StripPrefix("/ui", http.FileServer(http.Dir("./mm_ui"))))
	http.HandleFunc("/parsed", handler)
    http.ListenAndServe(":8080", Log(http.DefaultServeMux)) 
	
	// clean up instances in pool
	for len(instance_pool) > 0 {
		// drain an instance and tell it to clean up
		this_instance := <-instance_pool
		this_instance.Cleanup()
	}
	
}