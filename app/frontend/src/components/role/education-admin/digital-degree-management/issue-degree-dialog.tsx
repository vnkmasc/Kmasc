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
import { KeyRound } from 'lucide-react'
import { UseData } from '@/components/providers/data-provider'
import { formatDegreeTemplateOptions, formatFacultyOptionsByID } from '@/lib/utils/format-api'
import useSWR from 'swr'
import { issueDownloadDegreeZip, searchDegreeTemplateByFaculty } from '@/lib/api/degree'
import { showNotification } from '@/lib/utils/common'
import CommonSelect from '../common-select'
import { Label } from '@/components/ui/label'
import useSWRMutation from 'swr/mutation'
import { OptionType } from '@/types/common'

const IssueDegreeDialog: React.FC = () => {
  const [openSignDialog, setOpenSignDialog] = useState(false)
  const [selectFacultyId, setSelectFacultyId] = useState<string>('')
  const [selectDegreeTemplateId, setSelectDegreeTemplateId] = useState<string>('')
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)
  useEffect(() => {
    setSelectDegreeTemplateId('')
  }, [selectFacultyId])

  const findLabel = (id: string, options: OptionType[]) => {
    return options?.find((o: OptionType) => o.value === id)?.label
  }

  const queryDegreeTemplatesByFaculty = useSWR(
    selectFacultyId === '' ? undefined : 'degree-templates-by-faculty' + selectFacultyId,
    async () => {
      const res = await searchDegreeTemplateByFaculty(selectFacultyId)
      return formatDegreeTemplateOptions(res.data)
    },
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy danh sách mẫu bằng số')
      }
    }
  )

  const mutateIssueDigitalDegree = useSWRMutation(
    'issue-digital-degree-faculty',
    () =>
      issueDownloadDegreeZip(
        selectFacultyId,
        selectDegreeTemplateId,
        `VBS-${findLabel(selectFacultyId, facultyOptions)}-${findLabel(selectDegreeTemplateId, queryDegreeTemplatesByFaculty.data ?? [])}.zip`
      ),
    {
      onSuccess: () => {
        showNotification('success', 'Cấp bằng số cho chuyên ngành thành công')
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
        <Button>
          <KeyRound />
          <span className='hidden md:block'>Cấp bằng</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cấp bằng số</DialogTitle>
          <DialogDescription>
            Cấp bằng số cho các chuyên ngành và tự động tải xuống tệp <span className='font-bold'>.zip</span>
          </DialogDescription>
        </DialogHeader>
        <Label>Chọn chuyên ngành</Label>
        <CommonSelect
          value={selectFacultyId}
          options={facultyOptions}
          handleSelect={setSelectFacultyId}
          placeholder='Chọn chuyên ngành'
        />

        <Label>Chọn mẫu văn bằng</Label>
        <CommonSelect
          value={selectDegreeTemplateId}
          options={queryDegreeTemplatesByFaculty.data || []}
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
            Cấp bằng
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default IssueDegreeDialog
