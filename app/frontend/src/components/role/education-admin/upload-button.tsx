'use client'
import { Button } from '@/components/ui/button'
import { useRef, forwardRef, useImperativeHandle } from 'react'
import { UploadIcon } from 'lucide-react'

interface Props {
  // eslint-disable-next-line no-unused-vars
  handleUpload: (file: FormData) => void
  loading: boolean | false
  title?: string
  icon?: React.ReactNode
}

export interface UploadButtonRef {
  triggerUpload: () => void
}

const UploadButton = forwardRef<UploadButtonRef, Props>((props, ref) => {
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleButtonClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  useImperativeHandle(ref, () => ({
    triggerUpload: handleButtonClick
  }))

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (!files || files.length === 0) return

    try {
      // Process each file individually
      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        const formData = new FormData()
        formData.append('file', file)
        props.handleUpload(formData)
      }
    } catch (error: any) {
      console.error('Upload failed:', error)
    } finally {
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  return (
    <>
      <input
        ref={fileInputRef}
        type='file'
        accept='.xlsx, .xls, .csv, .pdf'
        onChange={handleFileChange}
        className='hidden'
        multiple
      />
      <Button
        variant='outline'
        onClick={handleButtonClick}
        isLoading={props.loading}
        title='Có hỗ trợ tải nhiều tệp cùng lúc'
      >
        {props.icon || <UploadIcon />}
        <span className='hidden sm:block'>{props.title || 'Tải tệp lên'}</span>
      </Button>
    </>
  )
})

UploadButton.displayName = 'UploadButton'

export default UploadButton
