package i18n

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"douyin/mylog" // Your logger
)

var (
   bundle         *i18n.Bundle
   supportedLangs []language.Tag
   defaultLang    language.Tag
)

// InitI18n initializes the i18n bundle with message files from localesPath.
func InitI18n(dLang language.Tag, sLangs []language.Tag, localesPath string) error {
	if len(sLangs) == 0 {
		return fmt.Errorf("no supported languages provided for i18n initialization")
	}
	defaultLang = dLang
	supportedLangs = sLangs
	bundle = i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := filepath.Glob(filepath.Join(localesPath, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to glob locale files in %s: %w", localesPath, err)
	}
	if len(files) == 0 {
		mylog.Warnf("No locale files found in %s. i18n might not work as expected.", localesPath)
	}

	for _, file := range files {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			mylog.Errorf("Failed to load message file %s: %v", file, err)
			// Continue loading other files
		} else {
			mylog.Infof("Loaded message file: %s", file)
		}
	}
	return nil
}

// GetLocalizer creates a localizer based on Gin context (Accept-Language header and 'lang' query param).
// It negotiates the language against the list of supported languages.
func GetLocalizer(c *gin.Context) *i18n.Localizer {
   langQuery := c.Query("lang")
   acceptLang := c.GetHeader("Accept-Language")

   desiredTags := []language.Tag{}
   if langQuery != "" {
	   tag, err := language.Parse(langQuery)
	   if err == nil {
		   desiredTags = append(desiredTags, tag)
	   }
   }

   if acceptLang != "" {
	   tags, _, err := language.ParseAcceptLanguage(acceptLang)
	   if err == nil {
		   desiredTags = append(desiredTags, tags...)
	   }
   }

   // Fallback to default if no specific language preference is found or if bundle is nil
   if bundle == nil { // Should not happen if InitI18n was called successfully
        mylog.Error("i18n Bundle is not initialized. Falling back to a placeholder localizer.")
        // Create a dummy bundle and localizer to prevent panic, though localization will not work.
        dummyBundle := i18n.NewBundle(language.English)
        return i18n.NewLocalizer(dummyBundle, language.English.String())
   }

   if len(desiredTags) == 0 {
	   desiredTags = append(desiredTags, defaultLang)
   }

   // Create a matcher with supported languages
   matcher := language.NewMatcher(supportedLangs)
   // Negotiate the best language
   tag, _, _ := matcher.Match(desiredTags...) // The chosen language tag

   return i18n.NewLocalizer(bundle, tag.String())
}

// MustLocalize is a helper to simplify localization calls from handlers.
// It retrieves the localizer from Gin context.
func MustLocalize(c *gin.Context, messageID string, templateData ...map[string]interface{}) string {
    // Ensure bundle is initialized before trying to use it.
    if bundle == nil {
        mylog.Error("i18n Bundle is not initialized. Cannot localize message.")
        // Return a formatted ID or placeholder if bundle is nil
        return formatIDWithData(messageID, templateData...)
    }

    loc, ok := c.Get("localizer")
    var localizerInstance *i18n.Localizer
    if !ok || loc == nil {
        mylog.Warn("Localizer not found in context. Creating a default one for this request.")
        // Fallback to default language localizer if not set in context
        localizerInstance = i18n.NewLocalizer(bundle, defaultLang.String())
    } else {
        localizerInstance, ok = loc.(*i18n.Localizer)
        if !ok {
            mylog.Error("Localizer in context is of wrong type. Using default.")
            localizerInstance = i18n.NewLocalizer(bundle, defaultLang.String())
        }
    }


	config := &i18n.LocalizeConfig{MessageID: messageID}
	if len(templateData) > 0 && templateData[0] != nil {
		config.TemplateData = templateData[0]
	}

	localizedMsg, err := localizerInstance.Localize(config)
	if err != nil {
		mylog.Warnf("Failed to localize messageID '%s': %v. Returning ID itself or formatted ID.", messageID, err)
		return formatIDWithData(messageID, templateData...)
	}
	return localizedMsg
}

// formatIDWithData is an internal helper to format a message ID with template data,
// used as a fallback when localization fails or bundle is nil.
func formatIDWithData(messageID string, templateData ...map[string]interface{}) string {
    if len(templateData) > 0 && templateData[0] != nil {
        formattedID := messageID
        for k, v := range templateData[0] {
            formattedID = strings.ReplaceAll(formattedID, "{"+k+"}", fmt.Sprintf("%v",v))
        }
        return formattedID
    }
    return messageID
}
