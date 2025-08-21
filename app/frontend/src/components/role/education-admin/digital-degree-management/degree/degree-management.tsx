'use client'

import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import { UseData } from '@/components/providers/data-provider'
import Filter from '@/components/common/filter'
import TableList from '@/components/common/table-list'
import { Badge } from '@/components/ui/badge'
import { CERTIFICATE_TYPE_OPTIONS, PAGE_SIZE } from '@/constants/common'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { useState } from 'react'
import useSWR from 'swr'
import { searchDigitalDegreeList, uploadDigitalDegreesBlockchain } from '@/lib/api/digital-degree'
import { formatDate } from 'date-fns'
import { Button } from '@/components/ui/button'
import { AlertCircleIcon, Blocks, CheckCircle2Icon, Eye } from 'lucide-react'
import IssueDegreeDialog from '@/components/role/education-admin/digital-degree-management/degree/issue-degree-dialog'
import SignDegreeButton from './sign-degree-button'
import { HashUploadButton } from './hash-upload-button'
import { findLabel, showNotification } from '@/lib/utils/common'
import useSWRMutation from 'swr/mutation'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger
} from '@/components/ui/alert-dialog'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import Link from 'next/link'
import CertificateQrCode from '@/components/common/certificate-qr-code'

const DegreeManagement = () => {
  const [filter, setFilter] = useState<any>({
    faculty_id: '',
    certificate_type: '',
    course: '',
    issued: 'true',
    page: 1
  })
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)
  const queryCertificates = useSWR('digital-degree-list' + JSON.stringify(filter), () =>
    searchDigitalDegreeList({
      ...filter,
      page: filter.page || 1,
      page_size: PAGE_SIZE,
      issued: filter.issued === 'true'
    })
  )

  const mutatePushDegreesBlockchain = useSWRMutation(
    'push-digital-degree-blockchain',
    async (_key, { arg }: { arg: any }) => {
      const formData = new FormData()
      formData.append('faculty_id', arg.faculty_id)
      if (filter.course !== '') formData.append('course', arg.course)
      if (filter.certificate_type !== '') formData.append('certificate_type', arg.certificate_type)

      const res = await uploadDigitalDegreesBlockchain(formData)
      queryCertificates.mutate()

      return res
    },
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi đẩy lên Blockchain')
      },
      onSuccess: () => {
        showNotification('success', 'Đẩy lên Blockchain thành công')
      }
    }
  )

  return (
    <>
      <PageHeader
        title='Quản lý văn bằng số'
        extra={[
          <IssueDegreeDialog
            key='sign-degree-faculty'
            facultyId={filter.faculty_id}
            certificateType={filter.certificate_type}
            course={filter.course}
          />,
          <SignDegreeButton key='sign-degree-button' />,
          <HashUploadButton key='hash-upload-button' />,
          <AlertDialog key='blockchain-alert'>
            <AlertDialogTrigger asChild>
              <Button title='Đẩy lên Blockchain' variant={'outline'}>
                <Blocks />
                <span className='hidden md:block'>Blockchain</span>
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Xác nhận đẩy lên Blockchain</AlertDialogTitle>
              </AlertDialogHeader>
              {filter.faculty_id ? (
                <Alert variant={'success'}>
                  <CheckCircle2Icon />
                  <AlertTitle>Sẵn sàng</AlertTitle>
                  <AlertDescription>
                    <ul className='list-inside list-disc'>
                      <li>Chuyên ngành: {findLabel(filter.faculty_id, facultyOptions)}</li>
                      {filter.certificate_type && <li>Loại bằng: {filter.certificate_type}</li>}
                      {filter.course && <li>Khóa học: {filter.course}</li>}
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
              )}
              <AlertDialogFooter>
                <AlertDialogCancel>Hủy bỏ</AlertDialogCancel>
                <AlertDialogAction
                  disabled={!filter.faculty_id}
                  onClick={() => mutatePushDegreesBlockchain.trigger(filter)}
                >
                  Xác nhận
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        ]}
      />

      <Filter
        items={[
          {
            type: 'select',
            name: 'faculty_id',
            placeholder: 'Chọn chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Chuyên ngành',
                    options: facultyOptions
                  }
                ]
              }
            }
          },
          {
            type: 'select',
            placeholder: 'Chọn loại bằng',
            name: 'certificate_type',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Bằng tốt nghiệp',
                    options: CERTIFICATE_TYPE_OPTIONS
                  }
                ]
              }
            }
          },
          {
            type: 'input',
            name: 'course',
            placeholder: 'Nhập khóa học'
            // setting: {
            //   input: {
            //     type: 'number'
            //   }
            // }
          },
          {
            type: 'select',
            name: 'issued',
            placeholder: 'Chọn trạng thái cấp',
            defaultValue: 'true',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Trạng thái',
                    options: [
                      { label: 'Đã cấp', value: 'true' },
                      { label: 'Chưa cấp', value: 'false' }
                    ]
                  }
                ]
              }
            }
          }
        ]}
        handleSetFilter={setFilter}
      />
      <TableList
        items={[
          { header: 'Mã SV', value: 'student_code', className: 'min-w-[80px] font-semibold text-blue-500' },
          { header: 'Họ và tên', value: 'student_name', className: 'min-w-[150px]' },
          { header: 'Chuyên ngành', value: 'faculty_name', className: 'min-w-[150px]' },
          { header: 'Tên văn bằng', value: 'full_name', className: 'min-w-[200px]' },
          {
            header: 'Phân loại',
            value: 'isDegree',
            render: (item) => (
              <Badge className='bg-blue-500 text-white hover:bg-blue-400'> {item.certificate_type}</Badge>
            )
          },
          { header: 'Mẫu bằng', value: 'template_name', className: 'min-w-[150px]' },
          { header: 'Khóa', value: 'course' },
          {
            header: 'Ngày cấp bằng',
            value: 'issue_date',
            className: 'min-w-[100px]',
            render: (item) => {
              return formatDate(item.issue_date, 'dd/MM/yyyy')
            }
          },
          {
            header: 'Trạng thái ký & mã',
            value: 'data_encrypted',
            className: 'min-w-[150px]',
            render: (item) => (
              <Badge variant={item.data_encrypted ? 'default' : 'outline'}>
                {item.data_encrypted ? 'Đã ký & mã' : 'Chưa ký & mã '}
              </Badge>
            )
          },
          {
            header: 'Blockchain',
            value: 'on_blockchain',

            render: (item) => (
              <Badge variant={item.on_blockchain ? 'default' : 'outline'}>
                {item.on_blockchain ? 'Đã đẩy' : 'Chưa đẩy'}
              </Badge>
            )
          },

          {
            header: 'Hành động',
            className: 'min-w-[100px]',
            value: 'action',
            render: (item) => (
              <div className='flex gap-2'>
                <Link href={`/education-admin/digital-degree-management/${item.id}`}>
                  <Button size={'icon'} variant={'outline'} title='Xem dữ liệu trên cơ sở dữ liệu'>
                    <Eye />
                  </Button>
                </Link>
                <Link href={`/education-admin/certificate-management/${item.id}/blockchain`}>
                  <Button size={'icon'} title='Xem dữ liệu trên blockchain'>
                    <Blocks />
                  </Button>
                </Link>
                <CertificateQrCode id={item.id} isIcon={true} />
              </div>
            )
          }
        ]}
        data={queryCertificates.data?.data || []}
        page={queryCertificates.data?.page || 1}
        pageSize={queryCertificates.data?.page_size || PAGE_SIZE}
      />
      <CommonPagination
        page={queryCertificates.data?.page || 1}
        totalPage={queryCertificates.data?.total_page || 1}
        handleChangePage={(page) => {
          setFilter({ ...filter, page })
        }}
      />
    </>
  )
}

export default DegreeManagement
