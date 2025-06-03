package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"douyin/mylog" // Your logger
)

// ConfigMail holds SMTP server configuration.
type ConfigMail struct {
	Host               string `mapstructure:"host" json:"host" yaml:"host"`
	Port               int    `mapstructure:"port" json:"port" yaml:"port"`
	Username           string `mapstructure:"username" json:"username" yaml:"username"`
	Password           string `mapstructure:"password" json:"password" yaml:"password"` // Consider secure handling
	From               string `mapstructure:"from" json:"from" yaml:"from"`             // Email address
	UseTLS             bool   `mapstructure:"use_tls" json:"use_tls" yaml:"use_tls"`    // Whether to use STARTTLS (direct TLS, not STARTTLS)
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify" json:"insecure_skip_verify" yaml:"insecure_skip_verify"` // For dev/test with self-signed certs
}

// Client is an email sending client.
type Client struct {
	cfg ConfigMail
}

// NewClient creates a new email client.
func NewClient(cfg ConfigMail) *Client {
	if cfg.Port == 0 {
		if cfg.UseTLS {
			cfg.Port = 465 // Default SMTPS port
		} else {
			cfg.Port = 587 // Default SMTP submission port (often uses STARTTLS via smtp.SendMail)
		}
	}
	return &Client{cfg: cfg}
}

// Send sends an email.
// 'body' is expected to be plain text. For HTML, set Content-Type header appropriately in msg.
func (c *Client) Send(to []string, subject, body string) error {
	if c.cfg.Host == "" || c.cfg.From == "" {
		mylog.Error("SMTP host or from address is not configured. Email not sent.")
		return fmt.Errorf("SMTP host or from address not configured")
	}

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)

	// Construct the email message
	// Adding basic headers for compatibility
	// Ensure CRLF line endings for email headers and body separation
	header := make(map[string]string)
	header["From"] = c.cfg.From
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"
	// Add MIME-Version and Date for good measure, though often handled by MTAs
	// header["MIME-Version"] = "1.0"
	// header["Date"] = time.Now().Format(time.RFC1123Z)


	var messageBuilder strings.Builder
	for k, v := range header {
		messageBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	messageBuilder.WriteString("\r\n") // End of headers
	messageBuilder.WriteString(body)

	msg := []byte(messageBuilder.String())

	if c.cfg.UseTLS { // This usually means SMTPS (TLS from the start)
		tlsconfig := &tls.Config{
			InsecureSkipVerify: c.cfg.InsecureSkipVerify, // Use with caution!
			ServerName:         c.cfg.Host,
		}

		// Dial connection
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			mylog.Errorf("Failed to dial TLS for email to %v: %v", to, err)
			return err
		}

		// Create new SMTP client
		client, err := smtp.NewClient(conn, c.cfg.Host)
		if err != nil {
			mylog.Errorf("Failed to create SMTP client with TLS for email to %v: %v", to, err)
			conn.Close() // Close connection if NewClient fails
			return err
		}
		defer client.Close()

		// Authenticate
		if c.cfg.Password != "" { // Only Auth if password is provided
		    if err = client.Auth(auth); err != nil {
			    mylog.Errorf("SMTP Auth error with TLS for email to %v: %v", to, err)
			    return err
		    }
        }

		// Set From
		if err = client.Mail(c.cfg.From); err != nil {
			mylog.Errorf("SMTP Mail (from) error with TLS for email to %v: %v", to, err)
			return err
		}
		// Set To
		for _, recipient := range to {
			if err = client.Rcpt(recipient); err != nil {
				mylog.Errorf("SMTP Rcpt (to %s) error with TLS for email to %v: %v", recipient, to, err)
				return err // Return on first recipient error
			}
		}
		// Get Data writer
		w, err := client.Data()
		if err != nil {
			mylog.Errorf("SMTP Data command error with TLS for email to %v: %v", to, err)
			return err
		}
		// Write message
		_, err = w.Write(msg)
		if err != nil {
			mylog.Errorf("Error writing email body with TLS for email to %v: %v", to, err)
			w.Close() // Close writer on error
			return err
		}
		// Close writer
		err = w.Close()
		if err != nil {
			mylog.Errorf("Error closing data writer with TLS for email to %v: %v", to, err)
			return err
		}
		// Quit
		client.Quit() // Best effort quit

	} else {
		// Standard smtp.SendMail (often uses port 25 or 587, might use STARTTLS implicitly if server supports)
		err := smtp.SendMail(addr, auth, c.cfg.From, to, msg)
		if err != nil {
			mylog.Errorf("Failed to send email to %v via SendMail: %v", to, err)
			return err
		}
	}

	mylog.Infof("Email sent to %v, subject: %s", to, subject)
	return nil
}
