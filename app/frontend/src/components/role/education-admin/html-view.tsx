import { Loader2 } from 'lucide-react'

interface Props {
  html?: string
  loading?: boolean
  // sandbox để an toàn khi render HTML từ server
  sandbox?: boolean
}

const HtmlView: React.FC<Props> = ({ html, loading = false, sandbox = true }) => {
  if (!html && !loading) {
    return (
      <div className='h-full w-full'>
        <p className='text-center text-red-500'>Không có nội dung HTML</p>
      </div>
    )
  }

  if (loading) {
    return (
      <div className='flex h-full w-full items-center justify-center'>
        <Loader2 className='h-4 w-4 animate-spin' />
        <p className='text-center text-sm text-gray-500'>Đang tải tài liệu HTML...</p>
      </div>
    )
  }

  const getHTMLDoc = () => {
    if (!html) return ''
    // Nếu đã là tài liệu đầy đủ thì dùng luôn, nếu không thì wrap lại
    const hasHtml = /<html[\s>]/i.test(html)
    if (hasHtml) return html
    return `<!DOCTYPE html><html><head><meta charset="utf-8"></head><body>${html}</body></html>`
  }

  const iframeSandbox = sandbox ? 'allow-same-origin allow-popups allow-forms allow-scripts' : undefined

  return (
    <div className={'h-full min-h-[700px] w-full overflow-x-scroll'}>
      <iframe className='h-full min-h-[700px] w-full' srcDoc={getHTMLDoc()} sandbox={iframeSandbox} />
    </div>
  )
}

export default HtmlView
