package archivedeps

// This package is the sandbox's *copy* of the api an archive library exposes —
// the same mechanic as hashdeps, iodeps and stddeps, for the same reason: the
// sandbox may import nothing but the sandbox, so `archive/zip` may not appear
// inside it. The contract is restated here, and the adapter — which lives
// outside the sandbox — is what fills it.
//
// An archive crosses this boundary whole, as a []byte, and its entries as a
// list of File: no reader, no writer and no type of the concrete library.

// File is one entry of an archive: where it sits in the archive and what it
// holds.
type File struct {
	// Path is the entry's name inside the archive, slash-separated and
	// relative to its root ("data/backofficedb/backoffice-user/size"). It is
	// handed back exactly as the archive spells it: whether it is safe to
	// write anywhere is the caller's to check.
	Path string
	// Content is the entry's whole uncompressed content.
	Content []byte
}

// Contract is the archive library injected whole as the Deps.ArchiveDeps
// field.
type Contract struct {
	// Zip packs files, in the order given, into one zip archive and returns
	// its bytes. Every entry is compressed and carries no modification time,
	// so the same files always pack into the same bytes.
	Zip func(files []File) ([]byte, error)

	// Unzip reads every file of the zip archive in archive, in the order the
	// archive lists them. Directory entries are not reported. The error
	// reports an archive that is not a zip or is damaged, and one whose
	// files add up to more than maxBytes once uncompressed — which is what
	// keeps a small archive from unpacking into all of the memory.
	Unzip func(archive []byte, maxBytes int64) ([]File, error)
}
