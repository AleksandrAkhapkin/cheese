package logger

import (
	"bytes"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
	"mime/multipart"
	"net/http"
	urler "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	chatID  = ""
	infoMsg = "INFO"
	errMsg  = "ERROR"
	token   = ""
	sashaID = ""
)

func NewLogger(t *config.Telegram) error {

	chatID = t.ChatID
	token = t.TelegramToken
	sashaID = t.SashaID

	return nil
}

var logger = zerolog.New(os.Stdout)
var Debug = true

func CheckDebug() {
	if Debug {
		infoMsg += "-DEBUG"
		errMsg += "-DEBUG"
	}
}

func LogError(err error) {

	logger.Err(err).Send()
	SendError(err)
}

func LogInfo(msg string) {
	logger.Info().Msg(msg)
	SendMessage(msg)
}

func LogOSForSasha(msg string) {
	logger.Info().Msg(msg)
	SendOSForSassha(msg)
}

func LogFatal(err error) {
	t := fmt.Sprintf("[%s]", time.Now().Format("2006-01-02T15:04:05"))
	err = errors.Wrap(err, t)
	SendError(err)
	logger.Fatal().Err(err).Send()

}
func SendError(err error) {

	url := makeURLSendMessage(errMsg, urler.QueryEscape(err.Error()))
	if err := send(url); err != nil {
		logger.Err(err).Send()
	}
}

func SendMessage(msg string) {
	url := makeURLSendMessage(infoMsg, msg)
	if err := send(url); err != nil {
		logger.Err(err).Send()
	}
}

func SendOSForSassha(msg string) {
	url := makeURLSendOS(infoMsg, msg)
	if err := send(url); err != nil {
		logger.Err(err).Send()
	}
}

func makeURLSendMessage(typeMsg, text string) string {

	text = fmt.Sprintf("%s [%s]: %s", typeMsg, time.Now().Format("2006-01-02T15:04:05"), text)
	str := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		token, chatID, text)
	return strings.ReplaceAll(str, " ", "+")
}

func makeURLSendOS(typeMsg, text string) string {

	text = fmt.Sprintf("%s [%s]: %s", typeMsg, time.Now().Format("2006-01-02T15:04:05"), text)
	str := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		token, sashaID, text)
	return strings.ReplaceAll(str, " ", "+")
}

func send(urlForSend string) error {
	req, err := http.NewRequest(http.MethodPost, urlForSend, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			logger.Err(err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d", res.StatusCode)
	}
	return nil
}

func Postphoto(directory string, userID int, billID int) {

	urlForPhoto := "https://api.telegram.org/bot" + token + "/" + "sendPhoto?chat_id=" + chatID
	directory = fmt.Sprintf("%s/%d/bill_id_%d", directory, userID, billID)
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	file, err := os.Open(directory)
	defer func() {
		err := file.Close()
		if err != nil {
			LogError(errors.Wrap(err, "err with file.Close in Postphoto"))
		}
	}()

	form, err := writer.CreateFormFile("photo", filepath.Base(directory))
	_, err = io.Copy(form, file)
	if err != nil {
		LogError(errors.Wrap(err, "err with io.Copy in Postphoto"))
		return
	}

	_ = writer.WriteField("caption", fmt.Sprintf("BILL_ID: %d\nuser_id: %d", billID, userID))
	err = writer.Close()
	if err != nil {
		LogError(errors.Wrap(err, "err with writer.Close in Postphoto"))
	}

	_, err = http.Post(urlForPhoto, writer.FormDataContentType(), buf)
	if err != nil {
		LogError(errors.Wrap(err, "err with http.Post in Postphoto"))
	}
}
