package session

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Session struct {
	History []string
}

func (s *Session) AddMessage(msg string) {
	s.History = append(s.History, msg)
	if len(s.History) > 10 { // keep last 10 exchanges
		s.History = s.History[len(s.History)-10:]
	}
}

func (s *Session) PromptLoop(handle func(string) string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			fmt.Println("👋 bye!")
			return
		}
		if input == "" {
			continue
		}
		s.AddMessage("User: " + input)
		resp := handle(input)
		s.AddMessage("AI: " + resp)
		fmt.Println(resp)
	}
}
