'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import { Key, KeyRound, PlusIcon } from 'lucide-react'
import { useCallback, useState } from 'react'
import Filter from '../../filter'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { UseData } from '@/components/providers/data-provider'
import useSWR from 'swr'
import {
  createDegreeTemplate,
  getDegreeTemplateById,
  searchDegreeTemplateByFaculty,
  signDegreeTemplateById,
  signDegreeTemplateFaculty,
  signDegreeTemplateUni,
  updateDegreeTemplate
} from '@/lib/api/degree'
import { showMessage, showNotification } from '@/lib/utils/common'
import TableList from '../../table-list'
import useSWRMutation from 'swr/mutation'
import DetailDialog from '../../detail-dialog'
import { validateNoEmpty } from '@/lib/utils/validators'
import { OptionType } from '@/types/common'
import TableActionButton from './table-action-button'
import { Badge } from '@/components/ui/badge'
import { DEGREE_TEMPLATE_STATUS } from '@/constants/common'
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
import CommonSelect from '../../common-select'
import { Label } from '@/components/ui/label'

const DegreeTemplate: React.FC = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)
  const facultyOptions = formatFacultyOptionsByID(UseData().facultyList)
  const [openSignDialog, setOpenSignDialog] = useState(false)
  const [selectFacultyId, setSelectFacultyId] = useState<string>('')
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

  const queryDegreeTemplateById = useSWR(idDetail, () => getDegreeTemplateById(idDetail as string), {
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

  const mutateSignDegreeTemplateFaculty = useSWRMutation(
    'sign-degree-template-faculty',
    (_key, { arg }: { arg: any }) =>
      async () => {
        if (selectFacultyId === '') {
          showMessage('Chuyên ngành không được để trống')
          return
        }
        return await signDegreeTemplateFaculty(arg)
      },
    {
      onSuccess: () => {
        if (selectFacultyId === '') return
        const matchFaculty = facultyOptions.find((faculty: OptionType) => faculty.value === selectFacultyId)
        showNotification('success', `Ký mẫu bằng số cho ${matchFaculty.label} thành công`)
        setOpenSignDialog(false)
        queryDegreeTemplatesByFaculty.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi ký mẫu bằng số cho chuyên ngành ')
      }
    }
  )

  const mutateSignDegreeTemplateUni = useSWRMutation('sign-degree-template-uni', () => signDegreeTemplateUni(), {
    onSuccess: () => {
      showNotification('success', 'Ký mẫu bằng số cho trường thành công')
      queryDegreeTemplatesByFaculty.mutate()
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi ký mẫu bằng số cho trường')
    }
  })

  const mutateSignDegreeTemplateById = useSWRMutation(
    'sign-degree-template-by-id',
    (_key, { arg }: { arg: any }) => signDegreeTemplateById(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Ký mẫu bằng số thành công')
        queryDegreeTemplatesByFaculty.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi ký mẫu bằng số')
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
          <Button
            key='sign-degree-template-all'
            variant={'secondary'}
            isLoading={mutateSignDegreeTemplateUni.isMutating}
            onClick={() => mutateSignDegreeTemplateUni.trigger()}
          >
            <Key />
            Ký trường
          </Button>,
          <Dialog key='sign-degree-template-faculty' open={openSignDialog} onOpenChange={setOpenSignDialog}>
            <DialogTrigger asChild>
              <Button variant={'outline'}>
                <KeyRound />
                Ký chuyên ngành
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Ký số cho chuyên ngành</DialogTitle>
                <DialogDescription>Ký số cho các mẫu văn bằng số của chuyên ngành</DialogDescription>
              </DialogHeader>
              <Label>Chọn chuyên ngành</Label>
              <CommonSelect
                value={selectFacultyId}
                options={facultyOptions}
                handleSelect={setSelectFacultyId}
                placeholder='Chọn chuyên ngành'
              />
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant='outline'>Hủy bỏ</Button>
                </DialogClose>
                <Button
                  type='submit'
                  isLoading={mutateSignDegreeTemplateFaculty.isMutating}
                  onClick={() => mutateSignDegreeTemplateFaculty.trigger(selectFacultyId)}
                >
                  Ký số
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>,
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
                    label: 'Chuyên ngành',
                    options: facultyOptions
                  }
                ]
              }
            }
          },
          {
            type: 'input',
            name: 'course',
            placeholder: 'Nhập khóa'
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
                {item.status}
              </Badge>
            )
          },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <TableActionButton
                canSign={item.status === 'PENDING'}
                handleSign={() => mutateSignDegreeTemplateById.trigger(item.id)}
                canEdit={!item.isLocked}
                handleSetIdDetail={setIdDetail}
                id={item.id}
              />
            )
          }
        ]}
      />
      <DetailDialog
        data={queryDegreeTemplateById.data?.data || {}}
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
                    label: 'Chuyên ngành',
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
