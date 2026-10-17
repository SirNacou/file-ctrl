import FileIcon from '~icons/lucide/file'
import FileArchiveIcon from '~icons/lucide/file-archive'
import FileCodeIcon from '~icons/lucide/file-code'
import FileImageIcon from '~icons/lucide/file-image'
import FileMusicIcon from '~icons/lucide/file-music'
import FileSpreadsheetIcon from '~icons/lucide/file-spreadsheet'
import FileTextIcon from '~icons/lucide/file-text'
import FileVideoIcon from '~icons/lucide/file-video'
import FolderIcon from '~icons/lucide/folder'

export const formatBytes = (bytes: number | bigint, decimals = 0): string => {
  const k = 1024n
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
  const dm = Math.max(0, decimals)

  if (typeof bytes === 'bigint') {
    let value = bytes
    let i = 0

    while (value >= k && i < sizes.length - 1) {
      value /= k
      i++
    }

    return `${value} ${sizes[i]}`
  }

  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), sizes.length - 1)

  return `${Number((bytes / Math.pow(1024, i)).toFixed(dm))} ${sizes[i]}`
}

export function getFileIcon(name: string, isDir: boolean) {
  if (isDir) return FolderIcon

  const ext = name.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'webp':
    case 'svg':
      return FileImageIcon
    case 'mp4':
    case 'mkv':
    case 'webm':
      return FileVideoIcon
    case 'mp3':
    case 'wav':
    case 'ogg':
      return FileMusicIcon
    case 'zip':
    case 'tar':
    case 'gz':
    case 'rar':
      return FileArchiveIcon
    case 'pdf':
    case 'txt':
    case 'md':
      return FileTextIcon
    case 'csv':
    case 'xlsx':
      return FileSpreadsheetIcon
    case 'go':
    case 'ts':
    case 'tsx':
    case 'js':
    case 'json':
    case 'yaml':
    case 'toml':
      return FileCodeIcon
    default:
      return FileIcon
  }
}
