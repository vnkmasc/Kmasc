'use client'

import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import { UseData } from '@/components/providers/data-provider'
import Filter from '@/components/role/education-admin/filter'
import TableList from '@/components/role/education-admin/table-list'
import { Badge } from '@/components/ui/badge'
import { CERTIFICATE_TYPE_OPTIONS, PAGE_SIZE } from '@/constants/common'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { useState } from 'react'
import useSWR from 'swr'
import { searchDigitalDegreeList, uploadDegreeToMinio } from '@/lib/api/digital-degree'
import { formatDate } from 'date-fns'
import { Button } from '@/components/ui/button'
import { Blocks, FolderUp } from 'lucide-react'
import useSWRMutation from 'swr/mutation'
import { showNotification } from '@/lib/utils/common'
import IssueDegreeDialog from '@/components/role/education-admin/digital-degree-management/degree/issue-degree-dialog'
import SignDegreeDialog from './sign-degree-dialog'

const DigitalDegreeManagement = () => {
  const [filter, setFilter] = useState<any>({})

  const queryCertificates = useSWR('digital-degree-list' + JSON.stringify(filter), () =>
    searchDigitalDegreeList({
      ...filter,
      page: filter.page || 1,
      page_size: PAGE_SIZE,
      is_issued: filter.is_issued === 'true'
    })
  )
  const mutateUploadDegreeToMinio = useSWRMutation('upload-degree-to-minio', uploadDegreeToMinio, {
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi mã hóa và đẩy lên Minio')
    },
    onSuccess: () => {
      showNotification('success', 'Mã hóa và đẩy lên Minio thành công')
    }
  })

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
          <SignDegreeDialog
            key='sign-degree-dialog'
            facultyId={filter.faculty_id}
            certificateType={filter.certificate_type}
            course={filter.course}
          />,
          <Button
            key='minio'
            isLoading={mutateUploadDegreeToMinio.isMutating}
            onClick={() => mutateUploadDegreeToMinio.trigger()}
            title='Mã hóa & lưu lên Minio'
          >
            <FolderUp />
            <span className='hidden md:block'>Minio</span>
          </Button>,
          <Button key='blockchain' title='Đẩy lên Blockchain' variant={'outline'}>
            <Blocks />
            <span className='hidden md:block'>Blockchain</span>
          </Button>
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
                    options: formatFacultyOptionsByID(UseData().facultyList)
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
            name: 'is_issued',
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
          { header: 'Mẫu bằng', value: 'template_name', className: 'min-w-[150px]' },
          { header: 'Khóa', value: 'course' },
          {
            header: 'Ngày cấp',
            value: 'issue_date',
            className: 'min-w-[100px]',
            render: (item) => {
              return formatDate(item.issue_date, 'dd/MM/yyyy')
            }
          },
          {
            header: 'Trạng thái ký',
            value: 'signed',
            className: 'min-w-[100px]',
            render: (item) => (
              <Badge variant={item.signed ? 'default' : 'outline'}>{item.signed ? 'Đã ký' : 'Chưa ký'}</Badge>
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
            header: 'Trạng thái mã',
            value: 'data_encrypted',

            render: (item) => (
              <Badge variant={item.data_encrypted ? 'default' : 'outline'}>
                {item.data_encrypted ? 'Đã mã hóa' : 'Chưa mã hóa'}
              </Badge>
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

export default DigitalDegreeManagement
