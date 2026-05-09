package internal

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// huhKeyMap returns a keymap with Esc bound to quit so every huh form
// shows "esc quit" in the bottom help row and actually exits on Esc.
func huhKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	)
	return km
}

// selectReleases shows a multi-select checkbox prompt with all releases pre-selected.
// Space to toggle, A to toggle all, Enter to confirm, Esc to quit.
func selectReleases(title string) ([]string, error) {
	selected := []string{releaseNBPlus, releaseMockGW}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Description("Space to toggle · A to toggle all · Enter to confirm").
				Options(
					huh.NewOption("payments-nbplus", releaseNBPlus).Selected(true),
					huh.NewOption("mock-go", releaseMockGW).Selected(true),
				).
				Value(&selected),
		),
	).WithKeyMap(huhKeyMap()).Run()

	if err == huh.ErrUserAborted {
		return nil, nil
	}
	return selected, err
}
