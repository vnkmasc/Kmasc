'use client'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { useEffect, useState } from 'react'
import { UseData } from '@/components/providers/data-provider'
import { formatDegreeTemplateOptions, formatFacultyOptionsByID } from '@/lib/utils/format-api'
import useSWR from 'swr'
import { downloadDegreeZip, searchDegreeTemplateByFaculty } from '@/lib/api/degree'
import { showNotification } from '@/lib/utils/common'
import CommonSelect from '../../common-select'
import { Label } from '@/components/ui/label'
import useSWRMutation from 'swr/mutation'
import { Download } from 'lucide-react'

const DownloadDialog: React.FC = () => {
  const [openSignDialog, setOpenSignDialog] = useState(false)
  const [selectFacultyId, setSelectFacultyId] = useState<string>('')
  const [selectDegreeTemplateId, setSelectDegreeTemplateId] = useState<string>('')

  useEffect(() => {
    setSelectDegreeTemplateId('')
  }, [selectFacultyId])

  const queryDegreeTemplatesByFaculty = useSWR(
    selectFacultyId === '' ? undefined : 'degree-templates-by-faculty-1' + selectFacultyId,
    () => searchDegreeTemplateByFaculty(selectFacultyId),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy danh sách mẫu bằng số')
      }
    }
  )

  const mutateIssueDigitalDegree = useSWRMutation(
    'download-digital-degree-faculty',
    () => downloadDegreeZip(selectFacultyId, selectDegreeTemplateId),
    {
      onSuccess: () => {
        showNotification('success', 'Tải bằng số của chuyên ngành thành công')
        setOpenSignDialog(false)
      }
    }
  )

  return (
    <Dialog
      open={openSignDialog}
      onOpenChange={(open) => {
        if (!open) {
          setSelectFacultyId('')
          setSelectDegreeTemplateId('')
        }
        setOpenSignDialog(open)
      }}
    >
      <DialogTrigger asChild>
        <Button variant={'outline'}>
          <Download />
          <span className='hidden md:block'>Tải xuống</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tải bằng số</DialogTitle>
          <DialogDescription>
            Tải bằng số của chuyên ngành theo dạng <span className='font-bold'>.zip</span>
          </DialogDescription>
        </DialogHeader>
        <Label>Chọn chuyên ngành</Label>
        <CommonSelect
          value={selectFacultyId}
          options={formatFacultyOptionsByID(UseData().facultyList)}
          handleSelect={setSelectFacultyId}
          placeholder='Chọn chuyên ngành'
        />

        <Label>Chọn mẫu văn bằng</Label>
        <CommonSelect
          value={selectDegreeTemplateId}
          options={formatDegreeTemplateOptions(queryDegreeTemplatesByFaculty.data?.data ?? [])}
          handleSelect={setSelectDegreeTemplateId}
          placeholder='Chọn mẫu văn bằng số'
        />

        <DialogFooter>
          <DialogClose asChild>
            <Button variant='outline'>Hủy bỏ</Button>
          </DialogClose>
          <Button
            type='submit'
            isLoading={mutateIssueDigitalDegree.isMutating}
            onClick={() => mutateIssueDigitalDegree.trigger()}
          >
            Tải xuống
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default DownloadDialog
