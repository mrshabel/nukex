package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
 _   _ _    _ _  ______   __
| \ | | |  | | |/ |  __\ \/ /
|  \| | |  | | ' /| |__ \  / 
| . ' | |  | |  < |  __|/  \ 
| |\  | |__| | . \| |__/ /\ \
|_| \_|\____/|_|\_\____/_/\_\
`

// showBanner displays the app banner
func showBanner() {
	// Clear screen
	fmt.Println()

	banner := titleStyle.Render(logo)
	subtitle := dimStyle.Render(`Clean up your "node_modules" directories with ease!`)
	credits := dimStyle.Render("@mrshabel made it")

	fmt.Println(banner)
	fmt.Println(subtitle)
	fmt.Println(credits)
	fmt.Println()
}

// showResults renders the found directories in their shortened form
func showResults(dirs []DirInfo) {
	fmt.Println(successStyle.Render(fmt.Sprintf("🎉 Found %d node_modules directories:", len(dirs))))
	fmt.Println()

	for i, dir := range dirs {
		displayPath := shortenPath(dir.Path)

		// format it as: 1 - ../../
		line := fmt.Sprintf("%d - %s %s",
			i+1,
			pathStyle.Render(displayPath),
			sizeStyle.Render(fmt.Sprintf("(%s)", dir.FormattedSize)))

		fmt.Println(itemStyle.Render(line))
	}
	fmt.Println()
}

// showCompletionMessage renders completion message after cleaning the directories
func showCompletionMessage(cleanedDirs []DirInfo) {
	totalSize := calculateTotalSize(cleanedDirs)

	completion := fmt.Sprintf(`
	🎉 Cleanup completed successfully!

	✨ Deleted %d directories
	💾 Freed up %s of space

	🚀 You're no longer in the node modules black hole!
`, len(cleanedDirs), totalSize)

	fmt.Println(successStyle.Render(completion))
	fmt.Println(dimStyle.Render("Thanks for using Nukex!"))
}

// selectDirectoriesToDelete renders a multiselect form for choosing directories to delete
func selectDirectoriesToDelete(dirs []DirInfo) ([]DirInfo, error) {
	// compose options with the directory index as its value
	var options []huh.Option[int]
	var selected []int

	for idx, dir := range dirs {
		// append label and index
		label := fmt.Sprintf("%s (%s)", dir.Path, dir.FormattedSize)
		options = append(options, huh.Option[int]{Key: label, Value: idx})
	}

	// build form
	form := huh.NewForm(huh.NewGroup(huh.NewMultiSelect[int]().Options(
		options...,
	).Title("📂 Select directories to delete:").Value(&selected)))

	err := form.Run()
	if err != nil {
		return nil, err
	}

	// get the selected directories
	var selectedDirs []DirInfo
	for _, idx := range selected {
		selectedDirs = append(selectedDirs, dirs[idx])
	}

	return selectedDirs, nil
}

// confirmDeletion renders a confirmation prompt for the user
func confirmDeletion(selectedDirs []DirInfo) (bool, error) {
	// print selected directories information
	fmt.Println()
	fmt.Println(warningStyle.Render(fmt.Sprintf("⚠️  You've selected %d directories for deletion:", len(selectedDirs))))
	fmt.Println()

	for idx, dir := range selectedDirs {
		// 1 - path (size)
		info := fmt.Sprintf(" %d - %s %s", idx+1, pathStyle.Render(dir.Path), sizeStyle.Render(fmt.Sprintf("(%s)", dir.FormattedSize)))

		fmt.Println(info)
	}

	totalSize := calculateTotalSize(selectedDirs)
	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("💾 Total space to be freed: %s", totalSize)))
	fmt.Println()

	var confirm bool

	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Proceed with deletion?").
		Description("This action cannot be undone").Value(&confirm)))
	if err := form.Run(); err != nil {
		return confirm, err
	}
	return confirm, nil
}

// deleteDirectories removes all the selected directories
func deleteDirectories(selectedDirs []DirInfo) {
	// print header
	fmt.Println()
	fmt.Println(headerStyle.Render("🗑️	Deleting directories..."))
	fmt.Println()

	// delete directories and its children
	for _, dir := range selectedDirs {
		if err := os.RemoveAll(dir.Path); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to delete %s: %v", dir.Path, err)))
			continue
		}

		fmt.Println(successStyle.Render(fmt.Sprintf("✅ Deleted %s", dir.Path)))
	}

}

// styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#01FAC6")).
			Bold(true).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(1, 1)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginLeft(2)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	sizeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAA00")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)
