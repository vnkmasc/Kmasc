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
import { useState } from 'react'
import { KeyRound } from 'lucide-react'
import { UseData } from '@/components/providers/data-provider'
import { formatDegreeTemplateOptions, formatFacultyOptionsByID } from '@/lib/utils/format-api'
import useSWR from 'swr'
import { issueDownloadDegreeZip, searchDegreeTemplateByFaculty } from '@/lib/api/digital-degree'
import { showNotification } from '@/lib/utils/common'
import CommonSelect from '../../common-select'
import { Label } from '@/components/ui/label'
import useSWRMutation from 'swr/mutation'
import { OptionType } from '@/types/common'

interface Props {
  facultyId: string
  certificateType: string
  course: string
}

const SignDegreeDialog: React.FC<Props> = (props) => {
  const [openSignDialog, setOpenSignDialog] = useState(false)
  const [selectDegreeTemplateId, setSelectDegreeTemplateId] = useState<string>('')
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)

  const findLabel = (id: string, options: OptionType[]) => {
    return options?.find((option: OptionType) => option.value === id)?.label
  }

  const queryDegreeTemplatesByFaculty = useSWR(
    props.facultyId ? 'degree-templates-by-faculty' + props.facultyId : undefined,
    async () => {
      const res = await searchDegreeTemplateByFaculty(props.facultyId)
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
        props.facultyId,
        selectDegreeTemplateId,
        `VBS-${findLabel(props.facultyId, facultyOptions)}-${findLabel(selectDegreeTemplateId, queryDegreeTemplatesByFaculty.data ?? [])}.zip`
      ),
    {
      onSuccess: () => {
        showNotification('success', 'Ký bằng số cho chuyên ngành thành công')
        setOpenSignDialog(false)
      }
    }
  )

  return (
    <Dialog
      open={openSignDialog}
      onOpenChange={(open) => {
        if (!open) {
          setSelectDegreeTemplateId('')
        }
        setOpenSignDialog(open)
      }}
    >
      <DialogTrigger asChild>
        <Button variant={'secondary'}>
          <KeyRound />
          <span className='hidden md:block'>Ký số</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Ký văn bằng số</DialogTitle>
          <DialogDescription>
            Ký bằng số cho các chuyên ngành và tự động tải xuống tệp <span className='font-semibold'>.zip</span>
          </DialogDescription>
        </DialogHeader>

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
            onClick={() => {
              if (selectDegreeTemplateId === '') {
                showNotification('error', 'Vui lòng chọn mẫu văn bằng số')
                return
              }
              mutateIssueDigitalDegree.trigger()
            }}
          >
            Ký số
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default SignDegreeDialog
