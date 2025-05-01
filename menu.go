package main

import (
	"PolAIn/internal/api"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var currentModel *ModelPresentation
var (
	textIcon      = ""
	adultIcon     = "🔞"
	eyeIcon       = "👁️"
	reasoningIcon = "🧠"
	emptyIcon     = "＿"
)

type LabelParts struct {
	Icons string
	Text  string
}

type ModelPresentation struct {
	*api.ModelDefinition
}

func (mp *ModelPresentation) getLabelParts() LabelParts {
	icons := textIcon
	if mp.Vision {
		icons += eyeIcon
	} else {
		icons += emptyIcon
	}
	if mp.Reasoning {
		icons += reasoningIcon
	} else {
		icons += emptyIcon
	}
	if mp.Uncensorded {
		icons += adultIcon
	} else {
		icons += emptyIcon
	}

	text := fmt.Sprintf("%s (%s, provider: %s)", mp.Name, mp.Description, mp.Provider)

	return LabelParts{
		Icons: icons,
		Text:  text,
	}
}

func (mp *ModelPresentation) getLabel() string {
	label := mp.getLabelParts()
	s := fmt.Sprintf("%s %s", label.Icons, label.Text)
	return s
}

var modelList = []*ModelPresentation{}

func init() {
	models := api.GetModels()

	for _, model := range models {
		modelList = append(modelList, &ModelPresentation{&model})
		if !model.Uncensorded && currentModel == nil {
			currentModel = modelList[len(modelList)-1]
		}
	}
}

func (a *App) getMenu() *menu.Menu {
	var modelMenu *menu.MenuItem

	modelItems := make([]*menu.MenuItem, len(modelList))
	for i, model := range modelList {
		modelItems[i] = &menu.MenuItem{
			Label: model.getLabel(),
			Type:  menu.RadioType,
			Click: func(current *menu.CallbackData) {
				// find the index of the menu
				if model.Uncensorded {
					runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
						Type:    runtime.InfoDialog,
						Title:   a.Translate("model.alert.uncensored.title"),
						Message: a.Translate("model.alert.uncensored.message"),
					})
				}
				currentModel = model
				runtime.EventsEmit(a.ctx, "selected-model", currentModel)
			},
		}
		modelItems[i].SetChecked(currentModel.Name == model.Name)
	}

	filemenu := &menu.MenuItem{
		Label: a.Translate("menu.conversation"),
		Role:  menu.WindowMenuRole,
		Type:  menu.TextType,
		SubMenu: menu.NewMenuFromItems(
			&menu.MenuItem{
				Label:       a.Translate("menu.conversation.new"),
				Accelerator: keys.CmdOrCtrl("n"),
				Type:        menu.TextType,
				Click: func(_ *menu.CallbackData) {
					a.NewConversation()
				},
			},
		),
	}
	modelMenu = &menu.MenuItem{
		Label:   a.Translate("menu.models"),
		Role:    menu.WindowMenuRole,
		Type:    menu.TextType,
		SubMenu: menu.NewMenuFromItems(modelItems[0], modelItems[1:]...),
	}

	helpmenu := &menu.MenuItem{
		Label: a.Translate("menu.help.title"),
		Role:  menu.WindowMenuRole,
		Type:  menu.TextType,
		SubMenu: menu.NewMenuFromItems(
			&menu.MenuItem{
				Label:       a.Translate("menu.help.title"),
				Accelerator: keys.CmdOrCtrl("h"),
				Role:        menu.WindowMenuRole,
				Type:        menu.TextType,
				Click: func(_ *menu.CallbackData) {
					runtime.EventsEmit(a.ctx, "show-help")
				},
			},
		),
	}
	return menu.NewMenuFromItems(filemenu, modelMenu, helpmenu)
}
