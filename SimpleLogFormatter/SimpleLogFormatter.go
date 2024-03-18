package SimpleLogFormatter

import (
	"bytes"

	"github.com/Miuzarte/ANSIFmt"
	log "github.com/sirupsen/logrus"
)

type LogFormat struct {
	RecordLog         func(content []byte)
	HasError, HasWarn bool
}

func (f *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(logLevelBanner[entry.Level])
	buf.WriteString(entry.Time.Format("[01/02|15:04:05] "))
	buf.WriteString(entry.Message)
	buf.WriteString("\n")

	if entry.Level == log.ErrorLevel {
		f.HasError = true
	}

	if entry.Level == log.WarnLevel {
		f.HasWarn = true
	}

	if f.RecordLog != nil {
		f.RecordLog(buf.Bytes())
	}
	return buf.Bytes(), nil
}

var (
	logLevelBanner = map[log.Level]string{
		log.TraceLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightBlue).
			Sprint("[TRAC]"),
		log.DebugLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightGreen).
			Sprint("[DEBU]"),
		log.InfoLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightWhite).
			Sprint("[INFO]"),
		log.WarnLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightYellow).
			Sprint("[WARN]"),
		log.ErrorLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed).
			Sprint("[ERRO]"),
		log.FatalLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed, ANSIFmt.Style.SlowBlink).
			Sprint("[FATA]"),
		log.PanicLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed, ANSIFmt.Style.SlowBlink, ANSIFmt.Style.Invert).
			Sprint("[PANI]"),
	}
)
