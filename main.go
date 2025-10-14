package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	// для io.Copy при копировании файлов
	// безопасная сборка путей (учитывает разделители)
)

// Функция для копирования файла из src в dst с схоранением прав доступа
func copyFile(src, dst string) error {

	//Открываем исходный файл
	// := - короткое объявление переменных
	// os.Open возвращает два значения of.File and error
	in, err := os.Open(src)

	// Если ошибка
	if err != nil {
		return fmt.Errorf("open src err: ", err)
	}
	// откладыает is.Close() до выхода из функции
	// Ставим defer после проверки err,чтобы in не была nill
	defer in.Close()

	return nil
}

func main() {

	var neededPackages bool = false

	if neededPackages == true {
		pkgList := []string{
			"hyprland",
			"xdg-desktop-portal-hyprland",
			"xorg-xwayland",
			"waybar",
			"hyprpaper",
			"hypridle",
			"hyprlock",
			"foot",
			"wl-clipboard",
			"grim",
			"slurp",
			"gamescope",
			"steam",
			"xdg-desktop-portal",
			"qt5-wayland",
			"qt6-wayland",
			"pipewire",
			"wireplumber",
			"pipewire-alsa",
			"pipewire-pulse",
			"polkit-gnome",
			"gvfs",
			"noto-fonts",
			"noto-fonts-emoji",
			"ttf-dejavu",
			"sddm",
			"qt6-virtualkeyboard",
			"brightnessctl",
			"wofi",
			"dolphin",
		}

		for _, pkg := range pkgList {

			hyprcmd := exec.Command("pacman", "--needed", "--noconfirm", "-S", pkg)

			hyprInstall, errInstall := hyprcmd.CombinedOutput()

			if errInstall != nil {

				log.Fatal("Ошибки установки", string(bytes.TrimSpace(hyprInstall)), errInstall)
			} else {

				fmt.Println(string(bytes.TrimSpace(hyprInstall)))
			}
		}

	}

}
