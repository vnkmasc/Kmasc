import { Loader2 } from 'lucide-react'

interface Props {
  url: string | undefined
  loading: boolean
}

const PDFView: React.FC<Props> = (props) => {
  if (!props.url && !props.loading)
    return (
      <div className='h-full w-full'>
        <p className='text-center text-red-500'>Không có file PDF</p>
      </div>
    )
  if (props.loading)
    return (
      <div className='flex h-full w-full items-center justify-center'>
        <Loader2 className='h-4 w-4 animate-spin' />
        <p className='text-center text-sm text-gray-500'>Đang tải file PDF...</p>
      </div>
    )
  const iframUrl = URL.createObjectURL(props.url as unknown as Blob)

  //   useEffect(() => {
  //     return () => {
  //       console.log('component unmount')

  //       URL.revokeObjectURL(iframUrl)
  //     }
  //   }, [iframUrl])

  return (
    <div className='h-full min-h-[500px] w-full'>
      <iframe src={iframUrl} className='h-full w-full' />
    </div>
  )
}

export default PDFView
