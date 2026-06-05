package command_handlers

import (
	"elex_storage/file_storage/internal/test"
	"elex_storage/file_storage/internal/use_case/cqrs/commands"
	"fmt"

	"testing"
)

func TestSaveFileHandler(t *testing.T) {
	app, logger, cfg := test.InjectMock(t)
	txtFile := []byte("this is simple text file for test the project")
	cmd := commands.SaveFileCommand{Name: "test_save_file.txt", Data: &txtFile}
	NewSaveFileHandler(logger, cfg, nil, nil, nil)

	fmt.Println(cmd.Name)
	app.Done()
}
