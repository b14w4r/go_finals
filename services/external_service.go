package services

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "errors"

    "gorm.io/gorm"
    "gopkg.in/gomail.v2"
)

type ExternalService interface {
    GetKeyRateCBR() (float64, error)
    SendEmail(to, subject, body string) error
    ComputeHMAC(data string, secret string) string
}

type externalService struct {
    smtpHost string
    smtpPort string
    smtpUser string
    smtpPass string
}

func NewExternalService(smtpHost, smtpPort, smtpUser, smtpPass, pgpPublicKey, pgpPrivateKey string) ExternalService {
    return &externalService{smtpHost: smtpHost, smtpPort: smtpPort, smtpUser: smtpUser, smtpPass: smtpPass}
}

func (s *externalService) GetKeyRateCBR() (float64, error) {
    // TODO: implement SOAP client to CBR
    return 0, nil
}

func (s *externalService) SendEmail(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.smtpUser)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    port, err := strconv.Atoi(s.smtpPort)
    if err != nil {
        return err
    }
    d := gomail.NewDialer(s.smtpHost, port, s.smtpUser, s.smtpPass)
    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}

func (s *externalService) ComputeHMAC(data string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}
