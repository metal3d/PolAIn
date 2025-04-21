package main

import (
	"PolAIn/internal/api"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var currentModel *ModelPresentation

type ModelPresentation struct {
	*api.ModelDefinition
}

func (mp *ModelPresentation) GetLabel() string {
	icon := "🔞"
	var label string
	if mp.Uncensorded {
		label = icon + mp.Name + "\n" + mp.Description
	} else {
		label = mp.Name + "\n" + mp.Description
	}
	return label
}

var modelList = []*ModelPresentation{}

func init() {
	models := api.GetModels()
	for _, model := range models {
		modelList = append(modelList, &ModelPresentation{&model})
	}
	currentModel = modelList[0]
}

func (a *App) getMenu() *menu.Menu {
	var modelMenu *menu.MenuItem

	modelItems := make([]*menu.MenuItem, len(modelList))
	for i, model := range modelList {
		modelItems[i] = &menu.MenuItem{
			Label: model.GetLabel(),
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
				Label: a.Translate("menu.conversation.new"),
				Type:  menu.TextType,
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
				Label: a.Translate("menu.help.title"),
				Role:  menu.WindowMenuRole,
				Type:  menu.TextType,
				Click: func(_ *menu.CallbackData) {
					runtime.EventsEmit(a.ctx, "show-help")
				},
			},
		),
	}
	return menu.NewMenuFromItems(filemenu, modelMenu, helpmenu)
}
