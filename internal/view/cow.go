package view

import (
	"context"
	"fmt"
	"strings"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/k9s/internal/model"
	"github.com/derailed/k9s/internal/ui"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// Cow represents a bomb viewer
type Cow struct {
	*tview.TextView

	actions ui.KeyActions
	app     *App
	says    string
}

// NewCow returns a have a cow viewer.
func NewCow(app *App, says string) *Cow {
	return &Cow{
		TextView: tview.NewTextView(),
		app:      app,
		actions:  make(ui.KeyActions),
		says:     says,
	}
}

// Init initializes the viewer.
func (c *Cow) Init(_ context.Context) error {
	c.SetBorder(true)
	c.SetScrollable(true).SetWrap(true).SetRegions(true)
	c.SetDynamicColors(true)
	c.SetHighlightColor(tcell.ColorOrange)
	c.SetTitleColor(tcell.ColorAqua)
	c.SetInputCapture(c.keyboard)
	c.SetBorderPadding(0, 0, 1, 1)
	c.updateTitle()
	c.SetTextAlign(tview.AlignCenter)

	c.app.Styles.AddListener(c)
	c.StylesChanged(c.app.Styles)

	c.bindKeys()
	c.SetInputCapture(c.keyboard)
	c.talk()

	return nil
}

func (c *Cow) talk() {
	says := c.says
	if len(says) == 0 {
		says = "Nothing to report here. Please move along..."
	}
	c.SetText(cowTalk(says))
}

func cowTalk(says string) string {
	buff := make([]string, 0, len(cow)+3)
	buff = append(buff, " "+strings.Repeat("─", len(says)+8))
	buff = append(buff, fmt.Sprintf("< [red::b]Ruroh? %s[-::-] >", says))
	buff = append(buff, " "+strings.Repeat("─", len(says)+8))
	spacer := strings.Repeat(" ", len(says)/2-8)
	for _, s := range cow {
		buff = append(buff, "[red::b]"+spacer+s)
	}
	return strings.Join(buff, "\n")
}

func (c *Cow) bindKeys() {
	c.actions.Set(ui.KeyActions{
		tcell.KeyEscape: ui.NewKeyAction("Back", c.resetCmd, false),
		ui.KeyQ: ui.NewKeyAction("Back", c.resetCmd, false),
	})
}

func (c *Cow) keyboard(evt *tcell.EventKey) *tcell.EventKey {
	if a, ok := c.actions[ui.AsKey(evt)]; ok {
		return a.Action(evt)
	}

	return evt
}

// StylesChanged notifies the skin changes.
func (c *Cow) StylesChanged(s *config.Styles) {
	c.SetBackgroundColor(c.app.Styles.BgColor())
	c.SetTextColor(c.app.Styles.FgColor())
	c.SetBorderFocusColor(c.app.Styles.Frame().Border.FocusColor.Color())
}

func (c *Cow) resetCmd(evt *tcell.EventKey) *tcell.EventKey {
	return c.app.PrevCmd(evt)
}

// Actions returns menu actions
func (c *Cow) Actions() ui.KeyActions {
	return c.actions
}

// Name returns the component name.
func (c *Cow) Name() string { return "cow" }

// Start starts the view updater.
func (c *Cow) Start() {}

// Stop terminates the updater.
func (c *Cow) Stop() {
	c.app.Styles.RemoveListener(c)
}

// Hints returns menu hints.
func (c *Cow) Hints() model.MenuHints {
	return c.actions.Hints()
}

// ExtraHints returns additional hints.
func (c *Cow) ExtraHints() map[string]string {
	return nil
}

func (c *Cow) updateTitle() {
	c.SetTitle(" Error ")
}

var cow = []string{
	`\   ^__^            `,
	` \  (oo)\_______    `,
	`    (__)\       )\/\`,
	`        ||----w |   `,
	`        ||     ||   `,
}
