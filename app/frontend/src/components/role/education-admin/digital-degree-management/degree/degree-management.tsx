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
import {
  searchDigitalDegreeList,
  uploadDigitalDegreesBlockchain,
  verifyDigitalDegreeDataBlockchain
} from '@/lib/api/digital-degree'
import { Button } from '@/components/ui/button'
import { AlertCircleIcon, Blocks, CheckCircle2Icon, Eye, Grid2X2Check } from 'lucide-react'
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
import { encodeJSON } from '@/lib/utils/lz-string'

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
      const res = await uploadDigitalDegreesBlockchain(
        arg.faculty_id,
        arg.certificate_type,
        arg.course,
        Boolean(arg.issued)
      )
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

  const mutateVerifyDigitalDegreeDataBlockchain = useSWRMutation(
    'verify-digital-degree-data-blockchain',
    async (_key, { arg }: { arg: any }) =>
      verifyDigitalDegreeDataBlockchain(
        arg.university_id,
        arg.faculty_id,
        arg.certificate_type,
        arg.course,
        arg.ediploma_id
      ),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi xác minh dữ liệu trên blockchain')
      },
      onSuccess: (data) => {
        if (!data.verified) {
          showNotification('error', data.message || 'Dữ liệu không hợp lệ')
        } else {
          showNotification('success', 'Xác minh dữ liệu trên blockchain thành công')
        }
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
                      {filter.issued && <li>Trạng thái cấp: {filter.issued === 'true' ? 'Đã cấp' : 'Chưa cấp'}</li>}
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
            className: 'min-w-[100px]'
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
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button size={'icon'} title='Xác minh dữ liệu trên blockchain' disabled={!item.on_blockchain}>
                      <Grid2X2Check />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Xác minh dữ liệu trên blockchain</AlertDialogTitle>
                    </AlertDialogHeader>
                    {item.faculty_id ? (
                      <Alert variant={'success'}>
                        <CheckCircle2Icon />
                        <AlertTitle>Sẵn sàng</AlertTitle>
                        <AlertDescription>
                          <ul className='list-inside list-disc'>
                            <li>ID Trường: {item.university_id}</li>
                            <li>ID Chuyên ngành: {item.faculty_id}</li>
                            {item.certificate_type && <li>Loại bằng: {item.certificate_type}</li>}
                            {item.course && <li>Khóa học: {item.course}</li>}
                            <li>ID Văn bằng: {item.id}</li>
                          </ul>
                        </AlertDescription>
                      </Alert>
                    ) : (
                      <Alert variant={'warning'}>
                        <AlertCircleIcon />
                        <AlertTitle>Cảnh báo</AlertTitle>
                        <AlertDescription>Chuyên ngành của văn bằng số không hợp lệ</AlertDescription>
                      </Alert>
                    )}
                    <AlertDialogFooter>
                      <AlertDialogCancel>Hủy bỏ</AlertDialogCancel>
                      <AlertDialogAction
                        disabled={item.faculty_id === ''}
                        onClick={() => mutateVerifyDigitalDegreeDataBlockchain.trigger(item)}
                      >
                        Xác minh
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
                <CertificateQrCode
                  id={
                    encodeJSON({
                      university_id: item.university_id,
                      university_code: item.university_code,
                      faculty_id: item.faculty_id,
                      certificate_type: item.certificate_type,
                      course: item.course,
                      ediploma_id: item.id
                    }) ?? ''
                  }
                  isIcon={true}
                />
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
