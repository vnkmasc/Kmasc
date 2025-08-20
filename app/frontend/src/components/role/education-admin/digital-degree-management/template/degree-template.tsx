'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import { PlusIcon } from 'lucide-react'
import { useCallback, useState } from 'react'
import Filter from '../../../../common/filter'
import { formatFacultyOptionsByID, formatTemplateInterfaceOptions } from '@/lib/utils/format-api'
import { UseData } from '@/components/providers/data-provider'
import useSWR from 'swr'
import {
  createDegreeTemplate,
  getDegreeTemplateById,
  getTemplateInterfaces,
  searchDegreeTemplateByFaculty,
  updateDegreeTemplate
} from '@/lib/api/digital-degree'
import { showNotification } from '@/lib/utils/common'
import TableList from '../../../../common/table-list'
import useSWRMutation from 'swr/mutation'
import { OptionType } from '@/types/common'
import TableActionButton from './table-action-button'
import { Badge } from '@/components/ui/badge'
import { DEGREE_TEMPLATE_STATUS } from '@/constants/common'
import DetailDialog from '../../detail-dialog'
import { validateNoEmpty } from '@/lib/utils/validators'

const DegreeTemplate: React.FC = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)

  const [filter, setFilter] = useState<any>({
    faculty_id: ''
  })
  const handleCloseDetailDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])

  const queryDegreeTemplatesByFaculty = useSWR(
    'degree-templates' + filter.faculty,
    () => searchDegreeTemplateByFaculty(filter.faculty),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy danh sách mẫu bằng số')
      }
    }
  )
  const queryTemplateInterfaces = useSWR('example-templates', async () => {
    const res = await getTemplateInterfaces()

    return formatTemplateInterfaceOptions(res.data)
  })

  const queryDegreeTemplateDetail = useSWR(idDetail, () => getDegreeTemplateById(idDetail as string), {
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi lấy thông tin mẫu bằng số')
    }
  })

  const mutateCreateDegreeTemplate = useSWRMutation(
    'create-degree-template',
    (_key, { arg }: { arg: any }) => createDegreeTemplate(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Tạo mẫu bằng số thành công')
        setIdDetail(undefined)
        queryDegreeTemplatesByFaculty.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi tạo mẫu bằng số')
      }
    }
  )

  const mutateUpdateDegreeTemplate = useSWRMutation(
    'update-degree-template',
    (_key, { arg }: { arg: any }) => updateDegreeTemplate(idDetail as string, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Cập nhật mẫu bằng số thành công')
        setIdDetail(undefined)
        queryDegreeTemplatesByFaculty.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật mẫu bằng số')
      }
    }
  )

  const handleSubmit = (data: any) => {
    if (idDetail) {
      mutateUpdateDegreeTemplate.trigger(data)
    } else {
      mutateCreateDegreeTemplate.trigger(data)
    }
  }

  const handleRefetchQueryList = useCallback(() => {
    queryDegreeTemplatesByFaculty.mutate()
  }, [queryDegreeTemplatesByFaculty])

  return (
    <div>
      <PageHeader
        title='Quản lý mẫu bằng số'
        extra={[
          <Button key='create-degree-template' onClick={() => setIdDetail(null)}>
            <PlusIcon />
            <span className='hidden md:block'>Tạo mới</span>
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
                    options: facultyOptions
                  }
                ]
              }
            }
          }
        ]}
        handleSetFilter={setFilter}
      />
      <TableList
        data={queryDegreeTemplatesByFaculty.data?.data || []}
        items={[
          { header: 'Tên mẫu bằng số', value: 'name', className: 'min-w-[150px]' },
          { header: 'Mô tả', value: 'description', className: 'min-w-[150px] max-w-[350px]' },
          {
            header: 'Chuyên ngành',
            value: 'facultyId',
            className: 'min-w-[100px]',
            render: (item) => facultyOptions.find((faculty: OptionType) => faculty.value === item.facultyId)?.label
          },
          {
            header: 'Trạng thái ký',
            value: 'status',
            className: 'min-w-[100px]',
            render: (item) => (
              <Badge
                variant={
                  DEGREE_TEMPLATE_STATUS[item.status as keyof typeof DEGREE_TEMPLATE_STATUS].variant as
                    | 'outline'
                    | 'secondary'
                    | 'default'
                    | 'destructive'
                }
                title={`${DEGREE_TEMPLATE_STATUS[item.status as keyof typeof DEGREE_TEMPLATE_STATUS].label}`}
              >
                {DEGREE_TEMPLATE_STATUS[item.status as keyof typeof DEGREE_TEMPLATE_STATUS].label}
              </Badge>
            )
          },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <TableActionButton
                canSign={item.status === 'PENDING'}
                canEdit={!item.isLocked}
                handleSetIdDetail={setIdDetail}
                id={item.id}
                refetch={handleRefetchQueryList}
                hashTemplate={item.hash_template}
                templateSampleId={item.template_sample_id}
              />
            )
          }
        ]}
      />
      <DetailDialog
        mode={idDetail === undefined ? undefined : idDetail ? 'update' : 'create'}
        data={queryDegreeTemplateDetail.data || {}}
        handleClose={handleCloseDetailDialog}
        handleSubmit={handleSubmit}
        items={[
          {
            type: 'input',
            name: 'name',
            placeholder: 'Nhập tên mẫu bằng số',
            label: 'Tên mẫu bằng số',
            validator: validateNoEmpty('Tên mẫu bằng số')
          },
          { type: 'textarea', name: 'description', placeholder: 'Nhập mô tả', label: 'Mô tả' },
          {
            type: 'select',
            name: 'faculty_id',
            placeholder: 'Chọn chuyên ngành',
            label: 'Chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Chuyên ngành',
                    options: facultyOptions
                  }
                ]
              }
            },
            validator: validateNoEmpty('Chuyên ngành'),
            disabled: !!idDetail
          },
          {
            type: 'select',
            name: 'template_sample_id',
            placeholder: 'Chọn giao diện mẫu bằng',
            label: 'Mẫu bằng',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Mẫu bằng',
                    options: queryTemplateInterfaces.data || []
                  }
                ]
              }
            },
            validator: validateNoEmpty('Mẫu bằng')
          }
        ]}
      />
    </div>
  )
}

export default DegreeTemplate
