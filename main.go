package main

import (
	"bytes"
	"fmt"
	"io" // io.Copy для копирования байтов
	"log"
	"os"
	"os/exec"
	"path/filepath" // безопасная сборка путей
)

// Функция для копирования файла из src в dst с схоранением прав доступа
// copyFile: копирует обычный файл src -> dst, создаёт родительскую папку и сохраняет права
func copyFile(src, dst string, perm os.FileMode) error {
	// гарантируем, что папка для dst есть
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %q: %w", filepath.Dir(dst), err)
	}
	// открываем исходный файл на чтение
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src %q: %w", src, err)
	}
	defer in.Close()

	// создаём/перезаписываем файл назначения
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm.Perm())
	if err != nil {
		return fmt.Errorf("create dst %q: %w", dst, err)
	}
	defer out.Close()

	// копируем содержимое «потоком»
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %q -> %q: %w", src, dst, err)
	}
	return nil
}

// CopySelected: копирует заданные относительные пути (папки/файлы) из текущей директории в destRoot
// Пример: CopySelected("/home/user/.config", []string{"hypr", "waybar"})
func CopySelected(destRoot string, paths []string) error {
	for _, rel := range paths {
		srcPath := filepath.Clean(rel) // путь источника (из текущей папки)
		info, err := os.Lstat(srcPath) // Lstat — чтобы различать симлинки
		if err != nil {
			return fmt.Errorf("stat %q: %w", srcPath, err)
		}

		dstPath := filepath.Join(destRoot, rel) // куда класть внутри destRoot

		// Если это каталог — обходим рекурсивно и копируем содержимое
		if info.IsDir() {
			err = filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// относительный путь внутри каталога
				subRel, err := filepath.Rel(srcPath, path)
				if err != nil {
					return err
				}
				target := filepath.Join(dstPath, subRel)

				switch mode := fi.Mode(); {
				case fi.IsDir():
					// создаём каталог (если уже есть — ок)
					return os.MkdirAll(target, mode.Perm())
				case mode&os.ModeSymlink != 0 || !mode.IsRegular():
					// для простоты пропускаем симлинки/спецфайлы
					return nil
				default:
					// обычный файл
					return copyFile(path, target, fi.Mode())
				}
			})
			if err != nil {
				return err
			}
			continue
		}

		// Если это файл (и не симлинк) — копируем как файл
		if info.Mode()&os.ModeSymlink == 0 && info.Mode().IsRegular() {
			if err := copyFile(srcPath, dstPath, info.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	var neededPackages bool = true

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
	destRoot := filepath.Join(os.Getenv("HOME"), ".config")

	// ЧТО копировать из текущей директории проекта
	toCopy := []string{"hypr", "waybar"} // твоя структура: папки с конфигами

	if err := CopySelected(destRoot, toCopy); err != nil {
		log.Fatal("Копирование не удалось: ", err)
	}
	fmt.Println("Готово: скопированы конфиги в:", destRoot)
}
