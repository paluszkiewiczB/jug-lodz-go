package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	whoErr    = true
	sentErr   = true
	structErr = true
)

func main() {
	name, err := whoAmI()
	if err != nil {
		log.Printf("could not check who I am, %s", err)
		// in most cases, you'd see an early 'return' here
		// return
	}

	hello(name)

	name, err = whoAreYou()
	// theoretically, you can check for specific error values with `err == ErrWhoAreYou`, but it's not recommended
	// it will not work with wrapped errors
	if errors.Is(err, ErrWhoAreYou) {
		log.Printf("type: %T, err: %s", err, err)
		log.Printf("could not check who you are, %s", err)
	} else {
		log.Printf("unsupported erorr for whoAreYou, exiting")
		return
	}

	name, err = whoIsIt()
	if err != nil {
		log.Printf("type: %T, err: %s", err, err.Error())
	}

	ce := ComplexErr{}
	if errors.As(err, &ce) {
		log.Printf("could not check who it is, %s", err)
		log.Printf("since the error occured at: %s, we could retry at: %s", ce.At, ce.At.Add(5*time.Second))
	}
}

// whoAmI returns the name or the error.
// Error is similar to the `new Exception()' in Java. You can check if it appeared, but it's hard to check for specific values.
func whoAmI() (string, error) {
	if whoErr {
		// alternatively, you can use fmt.Errorf("try again after: %d seconds", 5) - it allows message formatting
		return "", errors.New("try again later")
	}

	return "World", nil
}

func whoAreYou() (string, error) {
	if sentErr {
		return "", fmt.Errorf("wrapping mocno: %w", ErrWhoAreYou)
	}

	return "World", nil
}

func whoIsIt() (string, error) {
	if structErr {
		customErr := ComplexErr{
			Message: "try again later",
			At:      time.Now(),
			User:    "system", // read from the context/environment, whatever
			Reason:  "flag structErr is toggled",
		}
		return "", fmt.Errorf("hehe error wrapping: %w", customErr)
	}

	return "World", nil

}

func hello(who string) {
	log.Printf("Hello, %q!", who)
}

// ErrWhoAreYou is a sentinel error - it is returned as-is without any additional information.
// It is used to describe specific errors, usually when you cannot start or proceed with the operation and don't want to rely on the string returned from Error() method.
// Can you imagine doing the error check in Java based on exception.getMessage()?
// Nice articles:
// - https://www.digitalocean.com/community/tutorials/handling-errors-in-go
// - https://www.digitalocean.com/community/tutorials/how-to-add-extra-information-to-errors-in-go#handling-specific-errors-using-sentinel-errors
var ErrWhoAreYou = errors.New("failed to check who are you")

type ComplexErr struct {
	Message string
	At      time.Time
	User    string
	Reason  string
	Cause   error
}

func (e ComplexErr) Error() string {
	sb := strings.Builder{}
	sb.WriteString(e.Message)
	sb.WriteString(" at: ")
	sb.WriteString(e.At.Format(time.RFC3339))
	sb.WriteString(" by user: '")
	sb.WriteString(e.User)
	sb.WriteString("' because: ")
	sb.WriteString(e.Reason)
	return sb.String()
}

func (e ComplexErr) Unwrap() error {
	return e.Cause
}
