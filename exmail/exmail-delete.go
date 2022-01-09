/**
    * @Author lirui
    * @Date 2022/1/9 19:58
**/
package main

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("", ""); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("其他文件夹/prdalert", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	log.Println("一共 封邮件", mbox.Messages)

	// 标记删除

	for {

		from := uint32(1)
		to := mbox.Messages
		if mbox.Messages > 3 {
			// We're using unsigned integers here, only subtract if the result is > 0
			from = mbox.Messages - 3
		}

		seqset := new(imap.SeqSet)
		//seqset.AddNum(mbox.Messages)
		seqset.AddRange(from, to)

		// First mark the message as deleted
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}
		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Fatal(err)
		}

		// Then delete it
		if err := c.Expunge(nil); err != nil {
			log.Fatal(err)
		}

		log.Println("Last message has been deleted")

		log.Println("一共 封邮件", mbox.Messages)

	}
}
