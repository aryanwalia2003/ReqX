// Ek chhota row-editor helper — headers, cookies jaise "list of key/value
// pairs with a stable id for React keys" UI teeno jagah (SendRequestPanel,
// AuthFieldsEditor, CollectionEditorPanel) same shape chahiye tha.

export interface KeyValueRow {
  id: number
  key: string
  value: string
}

let nextRowId = 1

export function emptyRow(): KeyValueRow {
  return { id: nextRowId++, key: '', value: '' }
}

export function rowsFrom(record: Record<string, string> | undefined): KeyValueRow[] {
  const entries = Object.entries(record ?? {})
  return entries.length > 0
    ? entries.map(([key, value]) => ({ id: nextRowId++, key, value }))
    : [emptyRow()]
}

export function rowSetters(setRows: (updater: (rows: KeyValueRow[]) => KeyValueRow[]) => void) {
  return {
    update(id: number, patch: Partial<KeyValueRow>) {
      setRows((rows) => rows.map((row) => (row.id === id ? { ...row, ...patch } : row)))
    },
    remove(id: number) {
      setRows((rows) => (rows.length > 1 ? rows.filter((row) => row.id !== id) : rows))
    },
  }
}

export function rowsToRecord(rows: KeyValueRow[]): Record<string, string> | undefined {
  const record: Record<string, string> = {}
  for (const row of rows) {
    if (row.key.trim()) record[row.key.trim()] = row.value
  }
  return Object.keys(record).length > 0 ? record : undefined
}
