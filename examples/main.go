package main

import (
	goerrors "errors"
	"fmt"
	"strings"

	"github.com/bsido/go-errors/errors"
)

func main() {
	fmt.Print(errors.New("test"))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("test").Wrap(errors.New("wrapped plain error")))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("test").
		Wrap(errors.New("wrapped").
			Wrap(errors.New("wrapped plain error"))))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("test").
		Wrap(errors.New("wrapped")))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("test").
		Cause(goerrors.New("cause of the 1st error\nsecond line")).
		Wrap(errors.New("wrapped").
			Cause(goerrors.New("cause of the 2nd error\nsecond line"))))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("test notes").
		Cause(goerrors.New("cause \nsecond line")).
		Note("this is because \nthis and this").
		Note("also \nthat").
		Help("do this \nbecause of reasons").
		Help("also \ndo that"))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("with code").
		Cause(goerrors.New("cause \nsecond line")).
		Code(404))
	fmt.Print("\n----\n")

	err := errors.New("with code")
	fmt.Print(errors.New("test").
		HelpIf("do this \nbecause the condition is true", helpIfErrorContainsSt(err, "with code")))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("'Foo' is not an iterator").
		Code(277).
		Causef(`src/main.rs:4:16

     for foo in Foo {}
                ^^^ 'Foo' is not an iterator
`).Note("maybe try calling '.iter()' or a similar method").
		Help("the trait 'std::iter::Iterator' is not implemented for 'Foo'").
		Note("required by 'std::iter::IntoIterator::into_iter'").
		Wrap(errors.New("'&str' is not an iterator").
			Code(277).
			Causef(`src/main.rs:5:16

	 for foo in "" {}
				^^ '&str' is not an iterator
`).Help("call '.chars()' or '.bytes() on '&str'").
			Help("the trait 'std::iter::Iterator' is not implemented for '&str'").
			Note("required by 'std::iter::IntoIterator::into_iter'")))
	fmt.Print("\n----\n")

	fmt.Print(errors.New("this is the main message").
		Note("this is the 1st note\nthis is the second line of the first note").
		Note("this is the 2nd note\nthis is the second line of the second note"))

	fmt.Print("\n----\n")

	other := errors.New("second").
		Note("this might be happening because")

	fmt.Print(errors.New("first").
		Help("do this!").
		Wrap(other))
	fmt.Print("\n----\n")

	if err := callFunction(); err != nil {
		// print the error
		fmt.Print(err)
	}
}

func callFunction() error {
	err := callThirdParty()
	if err != nil {
		// 1. set a reasonable message
		// 2. optional: add additional details to the message as the cause
		// 3. wrap the error from coming from the function
		// 4. optional: add help and/or note messages
		return errors.New("failed to execute our own function").
			// optional: otherwise the cause will be empty if we wrap the error
			// you may leave it empty if the above message is sufficient
			Causef("some additional details about the error").
			// we wrap the error because this is the second level of the call stack
			Wrap(err).
			Help("do this to fix the error") // optional
	}

	return nil
}

func callThirdParty() error {
	// call a third party library
	err := goerrors.New("something went wrong")
	if err != nil {
		// 1. set a reasonable message
		// 2. set the error from the third party library as the cause
		return errors.New("failed to execute third party library").
			// we set the error as the cause because we know that this is not our error
			Cause(err)
	}

	return nil
}

func helpIfErrorContainsSt(err error, contains string) func() bool {
	return func() bool {
		return strings.Contains(err.Error(), contains)
	}
}
