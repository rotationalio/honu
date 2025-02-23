package opts

type TxOptions struct {
	// Begin a read-only transaction that cannot modify the database and can be
	// executed concurrently with other read-only transactions.
	ReadOnly bool
}

type ReadOptions struct {
	// Allow the inclusion of tombstone records (otherwise they are skipped or not found).
	Tombstones bool
}

type WriteOptions struct {
	// Require that the object exists before performing a deletion; otherwise an error
	// is returned (by default, an object that doesn't exist when deleted is ignored).
	CheckDelete bool

	// If the object already exists, return an error to prevent overwriting it.
	NoOverwrite bool

	// Require the object to exist before performing an update; otherwise an error
	// is returned (by default, an object that doesn't exist when updated is created).
	CheckUpdate bool
}

//===========================================================================
// TxOptions Getters and Defaults
//===========================================================================

func (to *TxOptions) GetReadOnly() bool {
	if to == nil {
		return false
	}
	return to.ReadOnly
}

//===========================================================================
// ReadOptions Getters and Defaults
//===========================================================================

func (ro *ReadOptions) GetTombstones() bool {
	if ro == nil {
		return false
	}
	return ro.Tombstones
}

//===========================================================================
// WriteOptions Getters and Defaults
//===========================================================================

func (wo *WriteOptions) GetCheckDelete() bool {
	if wo == nil {
		return false
	}
	return wo.CheckDelete
}

func (wo *WriteOptions) GetNoOverwrite() bool {
	if wo == nil {
		return false
	}
	return wo.NoOverwrite
}

func (wo *WriteOptions) GetCheckUpdate() bool {
	if wo == nil {
		return false
	}
	return wo.CheckUpdate
}
