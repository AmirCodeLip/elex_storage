package command_handlers

import (
	"elex_storage/file_metadata/internal/test"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/shared_kernel/event_models"
	"testing"

	"github.com/google/uuid"
)

func TestFileCreatedHandler(t *testing.T) {
	fileEntity := event_models.FileCreated{
		Id:    uuid.New(),
		Name:  "test_file.jpg",
		Drive: "cloud_1",
		Size:  200,
	}
	cmd := commands.FileCreatedCommand{FileEntity: fileEntity}
	app, logger, directoryMetadataRepositoryd, fileMetadataRepository := test.InjectMock(t)
	handler := NewFileCreatedHandler(logger, fileMetadataRepository, directoryMetadataRepositoryd)
	if err := handler.Handle(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if err := handler.Handle(cmd); err != nil {
		t.Fatal(err.Error())
	}
	files, err := fileMetadataRepository.GetFiles()
	if err != nil {
		t.Fatal(err.Error())
	}
	if files == nil || len(*files) < 1 {
		t.Errorf("expected one or more files, got %d", len(*files))
	}
	if (*files)[1].Name != "test_file" {
		t.Errorf("expected file name 'test_file', got '%s'", (*files)[0].Name)
	}
	app.Done()
}
