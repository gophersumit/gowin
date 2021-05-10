package notify

import "github.com/martinlindhe/notify"

const appName = "gowin"

func Alert(header, description string) {

	notify.Alert(appName, header, description, "")
}
