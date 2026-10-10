import FileIcon from '~icons/lucide/file'
import FileArchiveIcon from '~icons/lucide/file-archive'
import FileCodeIcon from '~icons/lucide/file-code'
import FileImageIcon from '~icons/lucide/file-image'
import FileMusicIcon from '~icons/lucide/file-music'
import FileSpreadsheetIcon from '~icons/lucide/file-spreadsheet'
import FileTextIcon from '~icons/lucide/file-text'
import FileVideoIcon from '~icons/lucide/file-video'
import FolderIcon from '~icons/lucide/folder'

export const formatBytes = (bytes: number, decimals = 0): string => {
  if (bytes === 0) return '--'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
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
