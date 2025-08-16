import { Button } from '@/components/ui/button'
import { uploadDigitalDegreesMinio } from '@/lib/api/digital-degree'
import { showMessage, showNotification } from '@/lib/utils/common'
import { ensurePermission, zipDirectoryHandleToFile } from '@/lib/utils/jszip'
import { FolderUp } from 'lucide-react'
import useSWRMutation from 'swr/mutation'

export const HashUploadButton = () => {
  const mutateuploadDigitalDegreesMinio = useSWRMutation(
    'upload-degree-to-minio',
    (_key, { arg }: { arg: FormData }) => uploadDigitalDegreesMinio(arg),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi mã hóa và đẩy lên Minio')
      },
      onSuccess: () => {
        showNotification('success', 'Mã hóa và đẩy lên Minio thành công')
      }
    }
  )

  const pickZipAndUpload = async () => {
    try {
      const dirHandle = await (window as any).showDirectoryPicker()
      const granted = await ensurePermission(dirHandle, 'read')
      if (!granted) {
        showMessage('Bạn đã từ chối quyền đọc thư mục')
        return
      }
      const zipFile = await zipDirectoryHandleToFile(dirHandle, 'ediplomas.zip')

      const formData = new FormData()
      formData.append('file', zipFile, 'ediplomas.zip')

      await mutateuploadDigitalDegreesMinio.trigger(formData)
    } catch (err: any) {
      if (err?.name === 'AbortError') return
    }
  }

  return (
    <Button
      isLoading={mutateuploadDigitalDegreesMinio.isMutating}
      onClick={pickZipAndUpload}
      title='Mã hóa & lưu lên Minio'
    >
      <FolderUp />
      <span className='hidden md:block'>Mã & lưu</span>
    </Button>
  )
}
