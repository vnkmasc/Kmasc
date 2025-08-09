'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import { PlusIcon, Signal } from 'lucide-react'
import { useCallback, useState } from 'react'
import Filter from '../filter'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { UseData } from '@/components/providers/data-provider'
import useSWR from 'swr'
import { createDegreeTemplate, searchDegreeTemplateByFaculty, updateDegreeTemplate } from '@/lib/api/degree'
import { showNotification } from '@/lib/utils/common'
import TableList from '../table-list'
import useSWRMutation from 'swr/mutation'
import DetailDialog from '../detail-dialog'
import { validateNoEmpty } from '@/lib/utils/validators'
import { OptionType } from '@/types/common'
import TableActionButton from './table-action-button'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { DEGREE_TEMPLATE_STATUS } from '@/constants/common'

const DegreeTemplate: React.FC = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)
  const [filter, setFilter] = useState<any>({
    faculty: ''
  })
  const handleCloseDetailDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])

  const queryDegreeTemplatesByFaculty = useSWR(
    filter.faculty === '' ? undefined : 'degree-templates' + filter.faculty,
    () => searchDegreeTemplateByFaculty(filter.faculty),
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy danh sách mẫu bằng số')
      }
    }
  )

  const mutateCreateDegreeTemplate = useSWRMutation(
    'create-degree-template',
    (_key, { arg }: { arg: any }) => createDegreeTemplate(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Tạo mẫu bằng số thành công')
        queryDegreeTemplatesByFaculty.mutate()
        setIdDetail(undefined)
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
        queryDegreeTemplatesByFaculty.mutate()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật mẫu bằng số')
      }
    }
  )

  const handleSubmit = useCallback(
    (data: any) => {
      if (idDetail) {
        mutateUpdateDegreeTemplate.trigger(data)
      } else {
        mutateCreateDegreeTemplate.trigger(data)
      }
    },
    [idDetail, mutateCreateDegreeTemplate, mutateUpdateDegreeTemplate]
  )

  return (
    <div>
      <PageHeader
        title='Quản lý mẫu bằng số'
        extra={[
          <Button key='sign-degree-template-all' variant={'secondary'}>
            <Signal />
            Ký toàn trường
          </Button>,
          <Button key='sign-degree-template' variant={'outline'}>
            <Signal />
            Ký khoa
          </Button>,
          <Button key='create-degree-template' onClick={() => setIdDetail(null)}>
            <PlusIcon />
            <span className='hidden sm:block'>Tạo mới</span>
          </Button>
        ]}
      />
      <Filter
        items={[
          {
            type: 'select',
            name: 'faculty',
            placeholder: 'Chọn chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Hệ đào tạo',
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
              <Tooltip>
                <TooltipTrigger>
                  <Badge
                    variant={
                      DEGREE_TEMPLATE_STATUS[item.status as keyof typeof DEGREE_TEMPLATE_STATUS].variant as
                        | 'outline'
                        | 'secondary'
                        | 'default'
                        | 'destructive'
                    }
                  >
                    {item.status}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  {DEGREE_TEMPLATE_STATUS[item.status as keyof typeof DEGREE_TEMPLATE_STATUS].label}
                </TooltipContent>
              </Tooltip>
            )
          },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <TableActionButton canEdit={!item.isLocked} handleSetIdDetail={setIdDetail} id={item.id} />
            )
          }
        ]}
      />
      <DetailDialog
        data={{}}
        mode={idDetail ? 'update' : idDetail === undefined ? undefined : 'create'}
        handleSubmit={handleSubmit}
        handleClose={handleCloseDetailDialog}
        items={[
          {
            type: 'input',
            placeholder: 'Nhập tên mẫu bằng số',
            name: 'name',
            label: 'Tên mẫu bằng số',
            validator: validateNoEmpty('Tên mẫu bằng số')
          },
          {
            type: 'textarea',
            placeholder: 'Nhập mô tả',
            name: 'description',
            label: 'Mô tả'
          },
          {
            type: 'select',
            name: 'facultyId',
            placeholder: 'Chọn chuyên ngành',
            label: 'Chuyên ngành',
            setting: {
              select: {
                groups: [
                  {
                    label: 'Hệ đào tạo',
                    options: facultyOptions
                  }
                ]
              }
            },
            disabled: idDetail !== null
          },
          {
            type: 'file',
            name: 'file',
            label: 'Mẫu văn bằng',
            setting: {
              file: {
                accept: '.html'
              }
            }
          }
        ]}
      />
    </div>
  )
}

export default DegreeTemplate
