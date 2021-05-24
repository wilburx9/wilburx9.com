package common

import (
	"errors"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

var (
	logrusLevelsToSentryLevels = map[log.Level]sentry.Level{
		log.PanicLevel: sentry.LevelFatal,
		log.FatalLevel: sentry.LevelFatal,
		log.ErrorLevel: sentry.LevelError,
		log.WarnLevel:  sentry.LevelWarning,
		log.InfoLevel:  sentry.LevelInfo,
		log.DebugLevel: sentry.LevelDebug,
		log.TraceLevel: sentry.LevelDebug,
	}
)

// SentryLogrusHook is a Sentry hook for Logrus.
// Original source: https://gist.github.com/HakShak/a5a92e21545206cb185dea54cd9974b5.
// TODO: Remove this when https://github.com/getsentry/sentry-go/issues/43 is implemented
type SentryLogrusHook struct {
	levels []log.Level
}

// NewSentryLogrusHook instantiates SentryLogrusHook
func NewSentryLogrusHook(levels []log.Level) SentryLogrusHook {
	return SentryLogrusHook{levels: levels}
}

// Levels returns supported logging levels
func (hook *SentryLogrusHook) Levels() []log.Level {
	return hook.levels
}

// Fire logs to Sentry
func (hook *SentryLogrusHook) Fire(entry *log.Entry) error {
	if Config.IsDebug() {
		return nil
	}
	var exception error

	if err, ok := entry.Data[log.ErrorKey].(error); ok && err != nil {
		exception = err
	} else {
		// Make a new error with the log message if there is no error provided
		// because stacktraces are neat
		exception = errors.New(entry.Message)
	}

	tags, hasTags := entry.Data["tags"].(map[string]string)

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.AddEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			event.Message = entry.Message
			return event
		})

		scope.SetLevel(logrusLevelsToSentryLevels[entry.Level])

		if hasTags {
			scope.SetTags(tags)
			delete(entry.Data, "tags")  // Remove ugly map rendering
			scope.SetExtras(entry.Data) // Set the extras in Sentry without the redundant tag data
			for k, v := range tags {    // Add the tags in a sane way back to Logrus
				entry.Data[k] = v
			}
		} else {
			scope.SetExtras(entry.Data)
		}

		sentry.CaptureException(exception)
	})

	return nil
}
