package app

type Handlers struct {
	ListDirectory           *ListDirectoryHandler
	ListUnregisteredSources *ListUnregisteredSourcesHandler
	ListStorageSources      *ListStorageSourcesHandler
	RegisterStorageSource   *RegisterStorageSourceHandler
	RemoveStorageSource     *RemoveStorageSourceHandler
	RenameSource            *RenameStorageHandler
	CreateFolder            *CreateFolderHandler
	RenameFile              *RenameFileHandler
	DeleteFiles             *DeleteFilesHandler
	OpenFile                *OpenFileHandler
	PrepareArchive          *PrepareArchiveHandler
	ValidateUpload          *ValidateUploadHandler
	UploadFile              *UploadFileHandler
}
