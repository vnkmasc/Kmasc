'use client'
import useSWRMutation from 'swr/mutation'
import { Button } from '@/components/ui/button'
import { FileTextIcon, DownloadIcon } from 'lucide-react'
import { showNotification } from '@/lib/utils/common'

interface Props {
  action: () => Promise<any>
}

const CertificateBlankButton: React.FC<Props> = (props) => {
  const mutateViewFile = useSWRMutation(`certificate-file-view`, props.action, {
    onSuccess: (data) => {
      const url = URL.createObjectURL(data)

      window.open(url, '_blank')

      setTimeout(() => {
        URL.revokeObjectURL(url)
      }, 100)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi xem tệp')
    }
  })
  const mutateDownloadFile = useSWRMutation(`certificate-file-download`, props.action, {
    onSuccess: (data) => {
      const url = URL.createObjectURL(data)

      const link = document.createElement('a')
      link.href = url
      link.download = `certificate.pdf`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      setTimeout(() => {
        URL.revokeObjectURL(url)
      }, 100)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi tải xuống tệp')
    }
  })

  return (
    <div className='flex items-center gap-2'>
      <Button onClick={() => mutateViewFile.trigger()} variant={'outline'}>
        <FileTextIcon />
        <span className='hidden md:block'>Xem tệp</span>
      </Button>

      <Button onClick={() => mutateDownloadFile.trigger()}>
        <DownloadIcon />
        <span className='hidden md:block'>Tải xuống</span>
      </Button>
    </div>
  )
}

export default CertificateBlankButton
