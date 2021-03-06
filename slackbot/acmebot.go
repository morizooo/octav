package slackbot

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat/go-cloud-acmeagent"
	"github.com/lestrrat/go-pdebug"
	"github.com/nlopes/slack"
)

type rtmctx struct {
	RTM     *slack.RTM
	Message *slack.MessageEvent
	Sync    bool
}

func (ctx *rtmctx) Reply(txt string) {
	ctx.RTM.SendMessage(ctx.RTM.NewOutgoingMessage(txt, ctx.Message.Channel))
}

func StartRTM(done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if pdebug.Enabled {
		g := pdebug.Marker("StartRTM")
		defer g.End()
	}

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for loop := true; loop; {
		select {
		case msg := <-rtm.IncomingEvents:
			if err := handleMessage(rtm, msg); err != nil {
				if pdebug.Enabled {
					pdebug.Printf("handleMessage: %s", err)
				}
				loop = false
			}
		case <-done:
			loop = false
		}
	}
}

type SlackLink struct {
	Text string
	URL  string
}

func parseSlackLink(s string) (*SlackLink, error) {
	if len(s) == 0 || s[0] != '<' {
		return nil, errors.New("not a link")
	}
	sl := &SlackLink{}
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '|':
			sl.Text = s[1:i]
		case '>':
			if l := len(sl.Text); l > 0 {
				sl.URL = sl.Text
				sl.Text = s[len(sl.Text)+2 : i]
			} else {
				sl.Text = s[1:i]
			}
			return sl, nil
		}
	}

	return nil, errors.New("not a link")
}

var spacesRx = regexp.MustCompile(`\s+`)

func handleMessage(rtm *slack.RTM, msg slack.RTMEvent) (err error) {
	switch msg.Data.(type) {
	case *slack.RTMError:
		return msg.Data.(*slack.RTMError)
	case *slack.InvalidAuthEvent:
		return errors.New("invalid auth")
	case *slack.MessageEvent:
		sm := msg.Data.(*slack.MessageEvent)

		cmd := spacesRx.Split(strings.TrimSpace(sm.Text), -1)
		if len(cmd) < 3 {
			return nil
		}

		sl, err := parseSlackLink(cmd[0])
		if err != nil || sl.Text != "@"+slackUser {
			return nil
		}

		if cmd[1] != "acme" {
			return nil
		}
		ctx := rtmctx{RTM: rtm, Message: sm}
		handleLetsEncryptCmd(&ctx, cmd[2:])
	}

	return nil
}

func handleLetsEncryptCmd(ctx *rtmctx, cmd []string) {
	switch cmd[0] {
	case "help":
		handleHelpCmd(ctx)
	case "authz":
		if len(cmd) < 2 {
			return
		}
		handleAuthzCmd(ctx, cmd[1:])
	case "cert":
		if len(cmd) < 3 {
			return
		}
		handleCertCmd(ctx, cmd[1:])
	}
}

func handleHelpCmd(ctx *rtmctx) {
	ctx.Reply(`usage: acme [cert|authz] [subcmds...]

acme cert issue <domain>
acme cert delete <domain>
acme cert upload <domain>
acme authz request <domain>
acme authz delete <domain>
`)
}

func handleAuthzCmd(ctx *rtmctx, cmd []string) {
	if len(cmd) < 2 {
		return
	}

	sl, err := parseSlackLink(cmd[1])
	if err != nil {
		return
	}
	domain := sl.Text
	switch cmd[0] {
	case "request":
		handleAuthzRequestCmd(ctx, domain)
	case "delete":
		handleAuthzDeleteCmd(ctx, domain)
	case "show":
		handleAuthzShowCmd(ctx, domain)
	default:
		ctx.Reply("Usage: `acme authz [request|delete|show] <domain>`")
	}
}

func handleAuthzDeleteCmd(ctx *rtmctx, domain string) {
	if err := acmeStateStore.DeleteAuthorization(domain); err != nil {
		ctx.Reply(":exclamation: Deleting authorization failed: " + err.Error())
		return
	}
	ctx.Reply(":tada: Deleted authorization")
}

func handleAuthzRequestCmd(ctx *rtmctx, domain string) {
	ctx.Reply(":white_check_mark: Authorizing *" + domain + "*")

	var authz acmeagent.Authorization
	if err := acmeStateStore.LoadAuthorization(domain, &authz); err != nil {
		ctx.Reply(":white_check_mark: Authorization for domain not found in storage.")
	} else {
		if authz.IsExpired() {
			ctx.Reply(":exclamation: Authorization expired, going to run authorization again")
		} else {
			ctx.Reply(":exclamation: Authorization already exists. Run `acme cert` to issue certificates for this domain")
			return
		}
	}

	ctx.Reply(":white_check_mark: Running authorization (this may take a few minutes)")

	cb := func() {
		if err := acmeAgent.AuthorizeForDomain(domain); err != nil {
			ctx.Reply(":exclamation: Authorization failed: " + err.Error())
			return
		}
		ctx.Reply(":tada: Authorization for domain *" + domain + "* complete")
	}

	if ctx.Sync {
		cb()
	} else {
		go cb()
	}
}

func handleAuthzShowCmd(ctx *rtmctx, domain string) {
	var authz acmeagent.Authorization
	if err := acmeStateStore.LoadAuthorization(domain, &authz); err != nil {
		ctx.Reply(":white_check_mark: Authorization for domain not found in storage.")
		return
	}

	buf, _ := json.MarshalIndent(authz, "", "  ")
	ctx.Reply("```\n" + string(buf) + "\n```")
}

func handleCertCmd(ctx *rtmctx, cmd []string) {
	switch len(cmd) {
	case 0, 1:
		ctx.Reply("Usage: `acme cert [issue|delete|upload] <domain>`")
		return
	default:
	}

	sl, err := parseSlackLink(cmd[1])
	if err != nil {
		return
	}
	domain := sl.Text

	switch cmd[0] {
	case "issue":
		handleCertIssueCmd(ctx, domain)
	case "delete":
		handleCertDeleteCmd(ctx, domain)
	case "upload":
		handleCertUploadCmd(ctx, domain)
	default:
		ctx.Reply("Usage: `acme cert [issue|delete|upload] <domain>`")
	}
}

func handleCertDeleteCmd(ctx *rtmctx, domain string) {
	ctx.Reply(":white_check_mark: Deleting certificates for *" + domain + "*")
	if err := acmeStateStore.DeleteCert(domain); err != nil {
		ctx.Reply(":exclamation: Failed to delete certificates: " + err.Error())
	} else {
		ctx.Reply(":tada: Deleted certificates")
	}
}

func handleCertIssueCmd(ctx *rtmctx, domain string) {
	ctx.Reply(":white_check_mark: Issueing certificates for *" + domain + "*")
	cert, err := acmeStateStore.LoadCert(domain)
	if err != nil {
		ctx.Reply(":white_check_mark: Certificates for domain not found in storage.")
	} else {
		if time.Now().After(cert.NotAfter) {
			ctx.Reply(":exclamation: Certificate expired, going to issue it again")
		} else {
			ctx.Reply(":exclamation: Certificate already exists. Run `acme upload` to upload the certificate")
			return
		}
	}

	// run handleAuthzCmd to make sure that the authorization is there
	osync := ctx.Sync
	ctx.Sync = true
	handleAuthzRequestCmd(ctx, domain)
	ctx.Sync = osync

	ctx.Reply(":white_check_mark: Fetching certificates")
	// Do this in a goroutine so we don't block from doing other things
	go func() {
		if err := acmeAgent.IssueCertificate(domain, nil, false); err != nil {
			ctx.Reply(":exclamation: Failed to fetch certificates: " + err.Error())
			return
		}
		ctx.Reply(":tada: Issueing certificates for domain *" + domain + "* complete")
	}()
}

func handleCertUploadCmd(ctx *rtmctx, domain string) {
	ctx.Reply(":white_check_mark: Uploading certificates for *" + domain + "*")

	name, err := acmeAgent.UploadCertificate(domain)
	if err != nil {
		ctx.Reply(":exclamation: Failed to upload certificates: " + err.Error())
		return
	}

	ctx.Reply(":tada: Certificates uploaded as *" + name + "*")
}
