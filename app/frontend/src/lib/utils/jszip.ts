import JSZip from 'jszip'

export type ZipEntry = {
  path: string
  file: File | Blob
}

const MIME_ZIP = 'application/zip'

export async function ensurePermission(handle: any, mode: 'read' | 'readwrite' = 'read'): Promise<boolean> {
  try {
    const status = await handle?.queryPermission?.({ mode })
    if (status === 'granted') return true
  } catch {
    console.error('Error querying permission')
  }
  try {
    const granted = await handle?.requestPermission?.({ mode })
    return granted === 'granted'
  } catch {
    return false
  }
}

async function ensureReadWritePermission(dirHandle: FileSystemDirectoryHandle): Promise<boolean> {
  return ensurePermission(dirHandle, 'readwrite')
}

export function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = filename
  document.body.appendChild(anchor)
  anchor.click()
  anchor.remove()
  setTimeout(() => URL.revokeObjectURL(url), 100)
}

async function ensureDirectory(
  root: FileSystemDirectoryHandle,
  pathSegments: string[]
): Promise<FileSystemDirectoryHandle> {
  let current = root
  for (const segment of pathSegments) {
    if (!segment) continue
    current = await current.getDirectoryHandle(segment, { create: true })
  }
  return current
}

async function writeFileToDirectory(dirHandle: FileSystemDirectoryHandle, filePath: string, data: Blob) {
  const segments = filePath.split('/')
  const filename = segments.pop() as string
  const targetDir = await ensureDirectory(dirHandle, segments)
  const fileHandle = await targetDir.getFileHandle(filename, { create: true })
  const writable = await fileHandle.createWritable()
  await writable.write(data)
  await writable.close()
}

export async function unzipAndSaveClient(zipBlob: Blob, dirHandle: FileSystemDirectoryHandle) {
  const hasPerm = await ensureReadWritePermission(dirHandle)
  if (!hasPerm) {
    throw new Error('Quyền ghi vào thư mục đã bị từ chối')
  }
  const zip = await JSZip.loadAsync(zipBlob)
  for (const path of Object.keys(zip.files)) {
    const entry = zip.files[path]
    if (entry.dir) continue
    const fileData = await entry.async('blob')
    await writeFileToDirectory(dirHandle, path, fileData)
  }
}

export function toZipEntriesFromFileList(fileList: FileList): ZipEntry[] {
  const entries: ZipEntry[] = []
  for (let i = 0; i < fileList.length; i++) {
    const file = fileList[i]
    const path = (file as any).webkitRelativePath || file.name
    entries.push({ path, file })
  }
  return entries
}

async function addDirectoryHandleToZip(zip: JSZip, dirHandle: FileSystemDirectoryHandle, basePath: string) {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  for await (const [name, handle] of dirHandle.entries()) {
    const currentPath = basePath ? `${basePath}/${name}` : name
    if ((handle as any).kind === 'directory') {
      await addDirectoryHandleToZip(zip, handle as FileSystemDirectoryHandle, currentPath)
    } else {
      const file = await (handle as FileSystemFileHandle).getFile()
      zip.file(currentPath, file)
    }
  }
}

export async function zipDirectoryHandleToFile(
  dirHandle: FileSystemDirectoryHandle,
  zipName = 'archive.zip'
): Promise<File> {
  const zip = new JSZip()
  await addDirectoryHandleToZip(zip, dirHandle, '')
  const blob = await zip.generateAsync({ type: 'blob' })
  return new File([blob], zipName, { type: MIME_ZIP })
}

export async function zipFilesWithPaths(files: ZipEntry[], zipName = 'archive.zip'): Promise<File> {
  const zip = new JSZip()
  for (const { path, file } of files) {
    zip.file(path, file)
  }
  const blob = await zip.generateAsync({ type: 'blob' })
  return new File([blob], zipName, { type: MIME_ZIP })
}
