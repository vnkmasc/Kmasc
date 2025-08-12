'use client'

import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import { UseData } from '@/components/providers/data-provider'
import DetailDialog from '@/components/role/education-admin/detail-dialog'
import Filter from '@/components/role/education-admin/filter'
import TableList from '@/components/role/education-admin/table-list'
import { Badge } from '@/components/ui/badge'
import { CERTIFICATE_TYPE_OPTIONS, PAGE_SIZE, STUDENT_CODE_SEARCH_SETTING } from '@/constants/common'

import { formatFacultyOptions } from '@/lib/utils/format-api'
import { validateNoEmpty } from '@/lib/utils/validators'
import { Fragment, useState } from 'react'

import { useCallback } from 'react'
import useSWR from 'swr'
import { searchDigitalDegreeList, uploadDegreeToMinio } from '@/lib/api/degree'
import SignDegreeDialog from './sign-degree-dialog'
import { formatDate } from 'date-fns'
import DownloadDialog from './download-dialog'
import { Button } from '@/components/ui/button'
import { Blocks, Hash } from 'lucide-react'
import useSWRMutation from 'swr/mutation'
import { showNotification } from '@/lib/utils/common'

const DigitalDegreeView = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)

  const [filter, setFilter] = useState<any>({})
  const handleCloseDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])

  const queryCertificates = useSWR('digital-degree-list' + JSON.stringify(filter), () =>
    searchDigitalDegreeList({
      ...filter,
      page: filter.page || 1,
      page_size: PAGE_SIZE
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
          <Fragment key='blockchain'>
            <Button variant={'secondary'}>
              <Blocks />
              Đẩy lên Blockchain
            </Button>
          </Fragment>,
          <Fragment key='hash'>
            <Button
              isLoading={mutateUploadDegreeToMinio.isMutating}
              onClick={() => mutateUploadDegreeToMinio.trigger()}
            >
              <Hash />
              Mã hóa & lưu Minio
            </Button>
          </Fragment>,
          <Fragment key='download-degree-faculty'>
            <DownloadDialog />
          </Fragment>,
          <Fragment key='sign-degree-faculty'>
            <SignDegreeDialog />
          </Fragment>
        ]}
      />

      <Filter
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'student_code',
            setting: STUDENT_CODE_SEARCH_SETTING
          },
          {
            type: 'select',
            name: 'faculty_code',
            placeholder: 'Chọn chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Chuyên ngành',
                    options: formatFacultyOptions(UseData().facultyList)
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
          }
        ]}
        handleSetFilter={setFilter}
      />
      <TableList
        items={[
          { header: 'Mã SV', value: 'student_code', className: 'min-w-[80px] font-semibold text-blue-500' },
          // { header: 'Họ và tên', value: 'studentName', className: 'min-w-[200px]' },
          { header: 'Chuyên ngành', value: 'faculty_name', className: 'min-w-[150px]' },
          { header: 'Tên văn bằng', value: 'full_name', className: 'min-w-[200px]' },
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
      <DetailDialog
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'studentCode',
            setting: STUDENT_CODE_SEARCH_SETTING,
            label: 'Mã sinh viên',
            validator: validateNoEmpty('Mã sinh viên')
          },
          {
            type: 'select',
            placeholder: 'Chọn loại bằng',
            name: 'certificateType',
            setting: {
              select: {
                groups: [
                  {
                    label: undefined,
                    options: CERTIFICATE_TYPE_OPTIONS
                  }
                ]
              }
            },
            label: 'Loại bằng',
            validator: validateNoEmpty('Loại bằng')
          },
          {
            type: 'input',
            placeholder: 'Nhập tên bằng',
            name: 'name',
            label: 'Tên bằng',
            validator: validateNoEmpty('Tên bằng')
          },
          {
            type: 'input',
            name: 'serialNumber',
            placeholder: 'Nhập số hiệu',
            label: 'Số hiệu',
            validator: validateNoEmpty('Số hiệu')
          },
          {
            type: 'input',
            name: 'regNo',
            placeholder: 'Nhập số vào sổ gốc cấp văn bằng',
            label: 'Số vào sổ gốc cấp văn bằng',
            validator: validateNoEmpty('Số vào sổ gốc cấp văn bằng')
          },
          {
            type: 'input',
            name: 'date',
            placeholder: 'Nhập ngày cấp',
            label: 'Ngày cấp',
            validator: validateNoEmpty('Ngày cấp'),
            setting: {
              input: {
                type: 'date'
              }
            }
          }
        ]}
        data={[]}
        mode={idDetail === null ? 'create' : undefined}
        handleSubmit={() => {}}
        handleClose={handleCloseDialog}
      />
    </>
  )
}

export default DigitalDegreeView
