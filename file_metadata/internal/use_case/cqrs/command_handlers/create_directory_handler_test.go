package command_handlers

import (
	"elex_storage/file_metadata/internal/test"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"testing"
)

func TestCreateDirectoryHandler(t *testing.T) {
	cmd := commands.CreateDirectoryCommand{Name: "test dir", ParentId: nil}
	app, logger, directoryMetadataRepository, _ := test.InjectMock(t)
	handler := NewCreateDirectoryHandler(logger, directoryMetadataRepository)
	if err, _ := handler.Handle(cmd); err != nil {
		t.Fatal(err.Error())
	}
	dirs, err := directoryMetadataRepository.GetDirectories()
	if err != nil {
		t.Fatal(err.Error())
	}
	if dirs == nil || len(*dirs) < 1 {
		t.Errorf("expected one or more directories, got %d", len(*dirs))
	}
	if (*dirs)[1].Name != "test dir" {
		t.Errorf("expected directory name 'test dir', got '%s'", (*dirs)[0].Name)
	}
	app.Done()
}
