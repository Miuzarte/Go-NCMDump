package main

import (
	"Go-NCMDump/NCMDump"
	"Go-NCMDump/SimpleLogFormatter"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DIRECTMODE = iota
	CONFIGMODE
)

var (
	exeDir string
	mode   int
)

var (
	corePlaintext = []byte(
		"bro u kn" + "ow, go i" +
			"s the be" + "st langu" +
			"age in t" + "he world")
	metaPlaintext = []byte(
		"dlrow eh" + "t ni ega" +
			"ugnal ts" + "eb eht s" +
			"i og ,wo" + "nk u orb")
	coreCiphertext = []byte{
		0x8B, 0x02, 0x74, 0x51, 0x6E, 0x05, 0x3F, 0xDE,
		0x29, 0xC6, 0xAE, 0xE4, 0xF2, 0x80, 0x87, 0x67,
		0x28, 0x53, 0x5E, 0xF4, 0x5E, 0xD7, 0xB7, 0x4F,
		0xBD, 0x26, 0xB5, 0xC5, 0x07, 0xB4, 0xC2, 0x14,
		0x81, 0xF9, 0x36, 0x15, 0x44, 0x53, 0xD3, 0x94,
		0x41, 0x94, 0x18, 0x27, 0xFC, 0x76, 0x58, 0xC3,
	}
	metaCiphertext = []byte{
		0x62, 0xDF, 0x8C, 0x62, 0x53, 0x9B, 0x75, 0x25,
		0x64, 0x7E, 0x30, 0x07, 0xB0, 0xDE, 0x00, 0x7E,
		0x73, 0x9E, 0xC7, 0xE1, 0x32, 0x23, 0x15, 0x2A,
		0xFD, 0x65, 0xBD, 0x1B, 0xFA, 0x45, 0xB1, 0xDD,
		0xA1, 0x1B, 0x67, 0x23, 0x6B, 0xC5, 0xD1, 0x06,
		0xE8, 0x90, 0x6C, 0x11, 0x04, 0x90, 0x21, 0xA3,
	}
)

var (
	formatter = &SimpleLogFormatter.LogFormat{}
	config    *Config
	dumper    = NCMDump.New()
	wg        sync.WaitGroup
)

func INIT() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(formatter)

	exePath, err := os.Executable()
	if err != nil {
		log.Panic("failed to get executable path: ", err)
	}
	exeDir = filepath.Dir(exePath)

	if len(os.Args) > 1 {
		mode = DIRECTMODE
	} else {
		mode = CONFIGMODE
	}

	config = config.init()
	coreOk, metaOk := checkKeys(config.CoreKey, config.MetaKey)
	if !coreOk && !metaOk {
		secondsToExitLog("keys seem to be incorrect")
	} else if !coreOk {
		secondsToExitLog("core key seems to be incorrect")
	} else if !metaOk {
		secondsToExitLog("meta key seems to be incorrect")
	}

	dumper.SetKeys(
		config.CoreKey,
		config.MetaKey).
		SetCoverOutput(
			config.CoverOutput,
			config.CoverEmbed,
			config.HighDefinitionCover)
}

func main() {
	INIT()

	var NCMList []string
	switch mode {
	case DIRECTMODE:
		for _, path := range os.Args[1:] {
			fs, err := os.Stat(path)
			if err != nil && !fs.IsDir() {
				NCMList = append(NCMList, path)
			} else if fs.IsDir() {
				filepath.Walk(path,
					func(path string, info os.FileInfo, err error) error {
						if err != nil {
							log.Fatal("failed to walk: ", err)
						}
						if !info.IsDir() && strings.HasSuffix(info.Name(), ".ncm") {
							NCMList = append(NCMList, path)
						}
						return nil
					})
			}

		}
		NCMList = os.Args[1:]
	case CONFIGMODE:
		_ = filepath.Walk(config.InputDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal("failed to walk: ", err)
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".ncm") {
					NCMList = append(NCMList, path)
				}
				return nil
			})
	}

	if config.MultiThread {
		multiThreadDump(NCMList)
	} else {
		singleThreadDump(NCMList)
	}

	if formatter.HasError {
		fmt.Println("Press ENTER to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	}
	if formatter.HasWarn {
		fmt.Println("exit in 5 seconds...")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

}

func multiThreadDump(NCMList []string) {
	wg.Add(len(NCMList))

	sem := make(chan struct{}, runtime.NumCPU())

	for _, path := range NCMList {
		go func() {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()

			err := dumper.DumpFile(path)
			if err != nil {
				log.Error(err)
			}
		}()
	}

	wg.Wait()
}

func singleThreadDump(NCMList []string) {
	wg.Add(len(NCMList))

	for _, path := range NCMList {
		go func() {
			defer func() {
				wg.Done()
			}()

			err := dumper.DumpFile(path)
			if err != nil {
				log.Error(err)
			}
		}()
	}

	wg.Wait()
}

func secondsToExitLog(msg ...interface{}) {
	log.Error(msg...)
	log.Warn("exit in 5 seconds, press ENTER to continue.")
	timer := time.AfterFunc(5*time.Second, func() {
		log.Info("time up")
		os.Exit(0)
	})
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	timer.Stop()
}
