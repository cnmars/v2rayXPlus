package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"gitlab.com/xiayesuifeng/v2rayxplus/conf"
	core2 "gitlab.com/xiayesuifeng/v2rayxplus/core"
	"gitlab.com/xiayesuifeng/v2rayxplus/styles"
	"io/ioutil"
	"strings"
)

type ConfigList struct {
	widgets.QWidget

	vboxLayout *widgets.QVBoxLayout

	buttonGroup *widgets.QButtonGroup

	addButton *widgets.QPushButton

	_ func() `constructor:"init"`

	_ func(name string) `signal:"configChange"`
	_ func(name string) `signal:"editConfig"`
	_ func(name string) `signal:"removeConfig"`
}

func (ptr *ConfigList) init() {
	ptr.vboxLayout = widgets.NewQVBoxLayout2(ptr)
	ptr.vboxLayout.SetSpacing(0)

	ptr.buttonGroup = widgets.NewQButtonGroup(ptr)

	ptr.addButton = widgets.NewQPushButton(ptr)
	ptr.addButton.SetFixedSize2(45, 45)
	ptr.addButton.SetStyleSheet(styles.GetStyleSheet(styles.AddButton))

	infos, err := ioutil.ReadDir(conf.V2rayConfigPath)
	if err == nil {
		for _, info := range infos {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
				name := strings.Split(info.Name(), ".json")[0]
				tmp := NewConfigListItem2(name, ptr)
				tmp.ConnectEditConfig(ptr.EditConfig)
				tmp.ConnectRemoveConfig(ptr.RemoveConfig)
				ptr.vboxLayout.AddWidget(tmp, 0, core.Qt__AlignCenter)
				ptr.buttonGroup.AddButton(tmp, 0)
			}
		}
	}

	if len(ptr.buttonGroup.Buttons()) > 0 {
		ptr.buttonGroup.Buttons()[0].SetChecked(true)
	}

	ptr.vboxLayout.AddSpacing(30)
	ptr.vboxLayout.AddWidget(ptr.addButton, 0, core.Qt__AlignHCenter)

	ptr.SetLayout(ptr.vboxLayout)

	ptr.initConnect()
}

func (ptr *ConfigList) initConnect() {
	ptr.buttonGroup.ConnectButtonClicked(func(button *widgets.QAbstractButton) {
		ptr.ConfigChange(button.Text())
	})

	ptr.addButton.ConnectClicked(func(checked bool) {
		v2ray := conf.NewV2rayConfig()
		name, path := core2.GetConfigName()
		if err := v2ray.Save(path); err != nil {
			widgets.QMessageBox_Warning(ptr, "错误", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}

		item := NewConfigListItem2(name, ptr)
		item.ConnectEditConfig(ptr.EditConfig)
		item.ConnectRemoveConfig(ptr.RemoveConfig)
		ptr.vboxLayout.InsertWidget(len(ptr.buttonGroup.Buttons()), item, 0, core.Qt__AlignCenter)
		ptr.buttonGroup.AddButton(item, 0)

		item.EditConfig(name)
	})
}
