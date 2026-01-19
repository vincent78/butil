package console

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/peterh/liner"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/utils/fileUtil"
)

var (
	// u: unlock, s: signXX, sendXX, n: newAccount, i: importXX
	passwordRegexp = regexp.MustCompile(`personal.[nusi]`)
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

var (
	Newline  = []byte{'\n'}
	Space    = []byte{' '}
	CmdSplit = " "
)

// HistoryFile is the file within the data directory to store input scrollback.
const HistoryFile = "history"

// DefaultPrompt is the default prompt line prefix to use for user input querying.
const DefaultPrompt = "> "

type Console struct {
	prompt   string       // Input prompt prefix string
	prompter UserPrompter // Input prompter to allow interactive user feedback
	histPath string       // Absolute path to the console scrollback history
	history  []string     // Scroll history maintained by the console
	printer  io.Writer    // Output writer to serialize any display strings to
	client   *Client      // the client of console,include the handle func
	config   Config

	interactiveStopped chan struct{}
	stopInteractiveCh  chan struct{}
	signalReceived     chan struct{}
	stopped            chan struct{}
	wg                 sync.WaitGroup
	stopOnce           sync.Once
}

func New(config Config) (*Console, error) {
	// Handle unset config values gracefully
	if config.Prompter == nil {
		config.Prompter = Stdin
	}
	if config.Prompt == "" {
		config.Prompt = DefaultPrompt
	}
	if config.Printer == nil {
		config.Printer = colorable.NewColorableStdout()
	}

	// Initialize the console and return
	console := &Console{
		client:             NewClient(config.Printer),
		prompt:             config.Prompt,
		prompter:           config.Prompter,
		printer:            config.Printer,
		histPath:           filepath.Join(config.DataDir, HistoryFile),
		config:             config,
		interactiveStopped: make(chan struct{}),
		stopInteractiveCh:  make(chan struct{}),
		signalReceived:     make(chan struct{}, 1),
		stopped:            make(chan struct{}),
	}

	if config.DataDir != "" {
		if err := fileUtil.MkDir(config.DataDir); err != nil {
			return nil, err
		}
	}

	if err := console.init(); err != nil {
		return nil, err
	}

	console.wg.Add(1)
	go console.interruptHandler()

	return console, nil
}

func (c *Console) RegisterHandler(api, name string, handler ApiHandleFunc) {
	c.client.RegisterApiHandler(api, name, handler)
}

func (c *Console) doHandler(apiName string, args ...string) model.RespModel {
	return c.client.DoHandler(apiName, args...)
}

func (c *Console) Print(msg string, opt ...any) {
	s := fmt.Sprintf("%v%v%v", string(Newline), msg, string(Newline))
	fmt.Fprintf(c.client.printer, s, opt...)
	fmt.Fprintf(c.client.printer, "%v%v", string(Newline), c.prompt)
}

// init retrieves the available APIs from the remote RPC provider and initializes
// the console's JavaScript namespaces based on the exposed modules.
func (c *Console) init() error {
	c.initConsoleObject()

	// Configure the input prompter for history and tab completion.
	if c.prompter != nil {
		if content, err := os.ReadFile(c.histPath); err != nil {
			c.prompter.SetHistory(nil)
		} else {
			c.history = strings.Split(string(content), string(Newline))
			c.prompter.SetHistory(c.history)
		}
		c.prompter.SetWordCompleter(c.AutoCompleteInput)
	}
	return nil
}

func (c *Console) initConsoleObject() {
}

//var defaultAPIs = map[string]string{"eth": "1.0", "net": "1.0", "debug": "1.0"}

func (c *Console) clearHistory() {
	c.history = nil
	c.prompter.ClearHistory()
	if err := os.Remove(c.histPath); err != nil {
		logger1.Error("can't delete history file: %v", err)
	} else {
		logger1.Debug("history file deleted")
	}
}

// consoleOutput is an override for the console.log and console.error methods to
// stream the output into the configured output stream instead of stdout.
//func (c *Console) consoleOutput(call goja.FunctionCall) goja.Value {
//	var output []string
//	for _, argument := range call.Arguments {
//		output = append(output, fmt.Sprintf("%v", argument))
//	}
//	fmt.Fprintln(c.printer, strings.Join(output, " "))
//	return goja.Null()
//}

// AutoCompleteInput is a pre-assembled word completer to be used by the user
// input prompter to provide hints to the user about the methods available.
func (c *Console) AutoCompleteInput(line string, pos int) (string, []string, string) {
	// No completions can be provided for empty inputs
	if len(line) == 0 || pos == 0 {
		return "", nil, ""
	}
	// Chunk data to relevant part for autocompletion
	// E.g. in case of nested lines eth.getBalance(eth.coinb<tab><tab>
	start := pos - 1
	for ; start > 0; start-- {
		// Skip all methods and namespaces (i.e. including the dot)
		if line[start] == '.' || (line[start] >= 'a' && line[start] <= 'z') || (line[start] >= 'A' && line[start] <= 'Z') {
			continue
		}
		// Handle web3 in a special way (i.e. other numbers aren't auto completed)
		if start >= 3 && line[start-3:start] == "web3" {
			start -= 3
			continue
		}
		// We've hit an unexpected character, autocomplete form here
		start++
		break
	}
	//return line[:start], c.jsre.CompleteKeywords(line[start:pos]), line[pos:]
	return line[:start], []string{}, line[pos:]
}

// Welcome show summary of current Geth instance and some metadata about the
// console's available modules.
func (c *Console) Welcome() {
	message := "Welcome to the console!\n"
	message += "To help, press the help() \n"
	message += "To exit, press ctrl-d or type exit"
	fmt.Fprintln(c.printer, message)
}

// Evaluate executes code and pretty prints the result to the specified output
// stream.
func (c *Console) Evaluate(statement string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(c.printer, "[native] error: %v\n", r)
			fmt.Fprintf(c.printer, "[native] cmd: %v\n", statement)
		}
	}()

	if len(statement) <= 0 {
		return
	}

	// 去掉最后一个换行符
	if strings.HasSuffix(statement, string(Newline)) {
		statement = statement[:len(statement)-1]
	}

	inputs := strings.Split(statement, CmdSplit)
	apiCmd := inputs[0]
	params := inputs[1:]

	resp := c.doHandler(apiCmd, params...)
	if resp.Code == model.Success {
		if resp.Data != nil {
			fmt.Fprintf(c.printer, "%v \n", resp.Data)
		}
	} else {
		fmt.Fprintf(c.printer, "%v - %v \n", resp.Code, resp.Message)
	}
	time.Sleep(500 * time.Microsecond)
	// Avoid exiting Interactive when jsre was interrupted by SIGINT.
	c.clearSignalReceived()
}

// interruptHandler runs in its own goroutine and waits for signals.
// When a signal is received, it interrupts the JS interpreter.
func (c *Console) interruptHandler() {
	defer c.wg.Done()

	// During Interactive, liner inhibits the signal while it is prompting for
	// input. However, the signal will be received while evaluating JS.
	//
	// On unsupported terminals, SIGINT can also happen while prompting.
	// Unfortunately, it is not possible to abort the prompt in this case and
	// the c.readLines goroutine leaks.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	defer signal.Stop(sig)

	for {
		select {
		case <-sig:
			c.setSignalReceived()
			//c.jsre.Interrupt(errors.New("interrupted"))
		case <-c.stopInteractiveCh:
			close(c.interactiveStopped)
			//c.jsre.Interrupt(errors.New("interrupted"))
		case <-c.stopped:
			return
		}
	}
}

func (c *Console) setSignalReceived() {
	select {
	case c.signalReceived <- struct{}{}:
	default:
	}
}

func (c *Console) clearSignalReceived() {
	select {
	case <-c.signalReceived:
	default:
	}
}

// StopInteractive causes Interactive to return as soon as possible.
func (c *Console) StopInteractive() {
	select {
	case c.stopInteractiveCh <- struct{}{}:
	case <-c.stopped:
	}
}

// Interactive starts an interactive user session, where input is prompted from
// the configured user prompter.
func (c *Console) Interactive() {
	var (
		prompt      = c.prompt             // the current prompt line (used for multi-line inputs)
		indents     = 0                    // the current number of input indents (used for multi-line inputs)
		input       = ""                   // the current user input
		inputLine   = make(chan string, 1) // receives user input
		inputErr    = make(chan error, 1)  // receives liner errors
		requestLine = make(chan string)    // requests a line of input
	)

	defer func() {
		c.writeHistory()
	}()

	// The line reader runs in a separate goroutine.
	go c.readLines(inputLine, inputErr, requestLine)
	defer close(requestLine)

	for {
		// Send the next prompt, triggering an input read.
		requestLine <- prompt

		select {
		case <-c.interactiveStopped:
			fmt.Fprintln(c.printer, "down, exiting console")
			return

		case <-c.signalReceived:
			// SIGINT received while prompting for input -> unsupported terminal.
			// I'm not sure if the best choice would be to leave the console running here.
			// Bash keeps running in this case. node.js does not.
			fmt.Fprintln(c.printer, "caught interrupt, exiting")
			return

		case err := <-inputErr:
			if err == liner.ErrPromptAborted {
				// When prompting for multi-line input, the first Ctrl-C resets
				// the multi-line state.
				prompt, indents, input = c.prompt, 0, ""
				continue
			}
			return

		case line := <-inputLine:
			// User input was returned by the prompter, handle special cases.
			if indents <= 0 && exit.MatchString(line) {
				return
			}
			if onlyWhitespace.MatchString(line) {
				continue
			}
			// Append the line to the input and check for multi-line interpretation.
			input += line + "\n"
			indents = countIndents(input)
			if indents <= 0 {
				prompt = c.prompt
			} else {
				prompt = strings.Repeat(".", indents*3) + " "
			}
			// If all the needed lines are present, save the command and run it.
			if indents <= 0 {
				if len(input) > 0 && input[0] != ' ' && !passwordRegexp.MatchString(input) {
					if command := strings.TrimSpace(input); len(c.history) == 0 || command != c.history[len(c.history)-1] {
						c.history = append(c.history, command)
						if c.prompter != nil {
							c.prompter.AppendHistory(command)
						}
					}
				}
				c.Evaluate(input)
				input = ""

			}
		}
	}
}

// readLines runs in its own goroutine, prompting for input.
func (c *Console) readLines(input chan<- string, errc chan<- error, prompt <-chan string) {
	for p := range prompt {
		line, err := c.prompter.PromptInput(p)
		if err != nil {
			errc <- err
		} else {
			input <- line
		}
	}
}

// countIndents returns the number of indentations for the given input.
// In case of invalid input such as var a = } the result can be negative.
func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '   // keep track of the string open char to allow var str = "I'm ....";
		charEscaped = false // keep track if the previous char was the '\' char, allow var str = "abc\"def";
	)

	for _, c := range input {
		switch c {
		case '\\':
			// indicate next char as escaped when in string and previous char isn't escaping this backslash
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c { // end string
				inString = false
			} else if !inString && !charEscaped { // begin string
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString { // ignore brackets when in string, allow var str = "a{"; without indenting
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}

	return indents
}

// Stop cleans up the console and terminates the runtime environment.
func (c *Console) Stop(graceful bool) error {
	c.stopOnce.Do(func() {
		// Stop the interrupt handler.
		close(c.stopped)
		c.wg.Wait()
	})
	return nil
}

func (c *Console) writeHistory() error {
	if err := os.WriteFile(c.histPath, []byte(strings.Join(c.history, "\n")), 0600); err != nil {
		return err
	}
	return os.Chmod(c.histPath, 0600) // Force 0600, even if it was different previously
}
