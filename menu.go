package main

import (
	"PolAIn/internal/api"
	"log"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (a *App) getMenu() *menu.Menu {
	models := api.GetModels()
	modelItems := make([]*menu.MenuItem, len(models))
	for i, model := range models {
		modelItems[i] = &menu.MenuItem{
			Label: model.Name,
			Type:  menu.RadioType,
			Click: func(current *menu.CallbackData) {
				log.Println("Clicked on model:", current.MenuItem.Label)
				currentModel = model.Name
			},
		}
		modelItems[i].SetChecked(currentModel == model.Name)
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
	modelMenu := &menu.MenuItem{
		Label:   a.Translate("menu.models"),
		Role:    menu.WindowMenuRole,
		Type:    menu.TextType,
		SubMenu: menu.NewMenuFromItems(modelItems[0], modelItems[1:]...),
	}
	return menu.NewMenuFromItems(filemenu, modelMenu)
}
