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
import { findLabel, showMessage, showNotification } from '@/lib/utils/common'
import CommonSelect from '../../common-select'
import { Label } from '@/components/ui/label'
import useSWRMutation from 'swr/mutation'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { getSignDegreeConfig } from '@/lib/utils/handle-storage'
import { verifyDigitalSignature } from '@/lib/utils/handle-vgca'
import { ensurePermission, unzipAndSaveClient } from '@/lib/utils/jszip'

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

    let dirHandle: any
    try {
      dirHandle = await (window as any).showDirectoryPicker()
    } catch (err: any) {
      if (err.name === 'AbortError') {
        return
      }
      showNotification('error', 'Không thể mở dialog chọn thư mục')
      return
    }

    const granted = await ensurePermission(dirHandle, 'readwrite')
    if (!granted) {
      showMessage('Bạn đã từ chối quyền đọc thư mục')
      return
    }

    try {
      setIssueLoading(true)
      const successVerify = await verifyDigitalSignature(
        matchDegreeTemplate?.signatureOfUni,
        matchDegreeTemplate?.hash_template
      )

      if (!successVerify) {
        showMessage('Xác minh chữ ký không thành công')
        return
      }

      showMessage('Xác minh chữ ký thành công')

      const blob = await issueDownloadDegreeZip(props.facultyId, selectDegreeTemplateId)

      await unzipAndSaveClient(blob, dirHandle)

      showNotification('success', 'Cấp bằng số thành công và tải thư mục thành công')
      setOpenSignDialog(false)
    } catch (err: any) {
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
            Cấp văn bằng số sẽ <strong>xác minh chữ ký</strong>,xác minh thành công sẽ tải xuống{' '}
            <strong>thư mục văn bằng.</strong>
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
            <Alert variant={'warning'}>
              <AlertCircleIcon />
              <AlertTitle>Cảnh báo</AlertTitle>
              <AlertDescription>
                Vui lòng chọn chuyên ngành trong <strong>phần tìm kiếm</strong> để tiến hành cấp bằng số.
              </AlertDescription>
            </Alert>
          )
        ) : (
          <Alert variant={'warning'}>
            <AlertCircleIcon />
            <AlertTitle>Cảnh báo</AlertTitle>
            <AlertDescription>
              Vui lòng cấu hình ký số cho <strong>link server ký số</strong>
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
            disabled={selectDegreeTemplateId === '' || signDegreeConfig.verifyService === '' || props.facultyId === ''}
          >
            Cấp bằng
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default IssueDegreeDialog
