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
import { AlertCircleIcon, CheckCircle2Icon, Plus } from 'lucide-react'
import { UseData } from '@/components/providers/data-provider'
import { formatDegreeTemplateOptions, formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { issueDownloadDegreeZip, searchDegreeTemplateByFaculty } from '@/lib/api/digital-degree'
import { showMessage, showNotification } from '@/lib/utils/common'
import CommonSelect from '../../common-select'
import { Label } from '@/components/ui/label'
import useSWRMutation from 'swr/mutation'
import { OptionType } from '@/types/common'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { getSignDegreeConfig } from '@/lib/utils/handle-storage'
import { verifyDigitalSignature } from '@/lib/utils/handle-vgca'
import { unzipAndSaveClient } from '@/lib/utils/jszip'

interface Props {
  facultyId: string
  certificateType: string
  course: string
}

const IssueDegreeDialog: React.FC<Props> = (props) => {
  const [openSignDialog, setOpenSignDialog] = useState(false)
  const [selectDegreeTemplateId, setSelectDegreeTemplateId] = useState<string>('')
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)
  const signDegreeConfig = getSignDegreeConfig()
  const [issueLoading, setIssueLoading] = useState(false)

  const findLabel = (id: string, options: OptionType[]) => {
    return options?.find((o: OptionType) => o.value === id)?.label
  }

  useEffect(() => {
    if (openSignDialog) {
      queryDegreeTemplatesByFaculty.trigger()
    } else {
      setSelectDegreeTemplateId('')
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.facultyId, openSignDialog])

  const queryDegreeTemplatesByFaculty = useSWRMutation(
    props.facultyId ? 'degree-templates-by-faculty' + props.facultyId : undefined,
    async () => {
      const res = await searchDegreeTemplateByFaculty(props.facultyId)
      return res.data
    },
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy danh sách mẫu bằng số')
      }
    }
  )

  const handleIssueDegreeClick = async () => {
    const matchDegreeTemplate = queryDegreeTemplatesByFaculty.data.find(
      (item: any) => item.id === selectDegreeTemplateId
    )

    if (matchDegreeTemplate.signatureOfUni === '') {
      showMessage('Chưa có chữ ký số cho mẫu của trường văn bằng này')
      return
    }

    await verifyDigitalSignature(matchDegreeTemplate?.signatureOfUni, matchDegreeTemplate?.hash_template)
    // *@*
    // if (!successVerify) {
    //   showMessage('Xác minh chữ ký không thành công')
    //   return
    // }

    try {
      setIssueLoading(true)
      const dirHandle = await (window as any).showDirectoryPicker()

      const perm = await dirHandle.requestPermission({ mode: 'readwrite' })
      if (perm !== 'granted') {
        throw new Error('Bạn đã từ chối quyền truy cập thư mục')
      }

      const blob = await issueDownloadDegreeZip(props.facultyId, selectDegreeTemplateId)

      await unzipAndSaveClient(blob, dirHandle)

      showNotification('success', 'Cấp bằng số thành công')
      setOpenSignDialog(false)
    } catch (err: any) {
      console.log(err)
      showNotification('error', err.message || 'Lỗi khi cấp bằng số')
    } finally {
      setIssueLoading(false)
    }
  }

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
        <Button>
          <Plus />
          <span className='hidden md:block'>Cấp bằng</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cấp văn bằng số</DialogTitle>
          <DialogDescription>
            Cấp văn bằng số sẽ <span className='font-semibold'>xác minh chữ ký</span>,xác minh thành công sẽ tải xuống{' '}
            <span className='font-semibold'>thư mục văn bằng.</span>
          </DialogDescription>
        </DialogHeader>

        {signDegreeConfig?.verifyService !== '' ? (
          props.facultyId ? (
            <Alert variant={'success'}>
              <CheckCircle2Icon />
              <AlertTitle>Sẵn sàng cấp bằng số</AlertTitle>
              <AlertDescription>
                <ul className='list-inside list-disc'>
                  <li>Chuyên ngành: {findLabel(props.facultyId, facultyOptions)}</li>
                  {props.certificateType && <li>Loại bằng: {props.certificateType}</li>}
                  {props.course && <li>Khóa học: {props.course}</li>}
                </ul>
              </AlertDescription>
            </Alert>
          ) : (
            <Alert variant={'destructive'}>
              <AlertCircleIcon />
              <AlertTitle>Cảnh báo</AlertTitle>
              <AlertDescription>
                Vui lòng chọn chuyên ngành trong <span className='font-semibold'>phần tìm kiếm</span> để tiến hành cấp
                bằng số.
              </AlertDescription>
            </Alert>
          )
        ) : (
          <Alert variant={'destructive'}>
            <AlertCircleIcon />
            <AlertTitle>Cảnh báo</AlertTitle>
            <AlertDescription>
              Vui lòng cấu hình ký số cho <span className='font-semibold'>link server ký số</span>
            </AlertDescription>
          </Alert>
        )}

        <Label>Chọn mẫu văn bằng</Label>
        <CommonSelect
          value={selectDegreeTemplateId}
          options={formatDegreeTemplateOptions(queryDegreeTemplatesByFaculty.data || [])}
          handleSelect={setSelectDegreeTemplateId}
          placeholder='Chọn mẫu văn bằng số'
        />

        <DialogFooter>
          <DialogClose asChild>
            <Button variant='outline'>Hủy bỏ</Button>
          </DialogClose>
          <Button
            type='submit'
            isLoading={issueLoading}
            onClick={handleIssueDegreeClick}
            disabled={selectDegreeTemplateId === '' || signDegreeConfig.verifyService === ''}
          >
            Cấp bằng
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default IssueDegreeDialog
