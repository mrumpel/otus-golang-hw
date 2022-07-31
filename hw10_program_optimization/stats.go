package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	s := bufio.NewScanner(r)
	result := make(DomainStat)

	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, fmt.Errorf("error in regexp compiling: %w", err)
	}

	for s.Scan() {
		var user User
		if err = easyjson.Unmarshal(s.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("error in json unmarshalling: %w", err)
		}

		if re.MatchString(user.Email) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}

	if s.Err() != nil {
		return nil, fmt.Errorf("error in scanning data: %w", s.Err())
	}

	return result, nil
}
