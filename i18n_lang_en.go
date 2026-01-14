//
// FilePath    : video-trim\i18n_lang_en.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : 英文翻译
//

package main

var langEN = map[string]string{
	KeyLanguageName:          "English",
	KeyTitle:                 "Video Trimmer",
	KeyHeaderUpload:          "Upload videos (trim head/tail seconds)",
	KeyChooseVideo:           "Choose videos",
	KeyUploadButton:          "Upload & Process",
	KeyUploadingText:         "Uploading and processing...",
	KeyClearButton:           "Clear uploaded and output files",
	KeyConfirmClear:          "Clear all uploaded and output files? This cannot be undone.",
	KeySelectAtLeastOne:      "Please select at least one video file before uploading",
	KeyFileTooLargePrefix:    "File \"",
	KeyFileTooLargeSuffix:    "\" exceeds allowed size ",
	KeyFileTooLargeEnd:       ", please reduce file size and retry.",
	KeyHeadLabel:             "Head trim seconds (editable)",
	KeyTailLabel:             "Tail trim seconds (editable, default 0)",
	KeyHint:                  "After processing, you'll be redirected to the download page; ensure browser and server are on the same LAN.",
	KeyProcessedTitle:        "Processed, click to download:",
	KeyDownloadAll:           "Download All",
	KeyDownload:              "Download",
	KeyReturnUpload:          "Return to Upload",
	KeyRemove:                "Remove",
	KeyAlertNoTrim:           "Head and tail trims are both 0, no processing needed",
	KeyNoProcessedFilesHint:  "No files were successfully processed, please check source files or FFmpeg logs.",
	KeyRequestBodyTooLarge:   "File too large, maximum allowed upload size is %s. Please reduce file size and retry.",
	KeyRequestParseError:     "Request body too large or unable to parse form",
	KeyCannotReadFile:        "Unable to read file %s",
	KeyFileEmptyOrUnreadable: "File %s is empty or unreadable",
	KeyNotSupportedVideo:     "File %s is not a supported video format (magic number check failed)",
}
