'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import { Key, KeyRound, PlusIcon } from 'lucide-react'
import { useCallback, useState } from 'react'
import Filter from '../../filter'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { UseData } from '@/components/providers/data-provider'
import useSWR from 'swr'
import { searchDegreeTemplateByFaculty, signDegreeTemplateFaculty, signDegreeTemplateUni } from '@/lib/api/degree'
import { showMessage, showNotification } from '@/lib/utils/common'
import TableList from '../../table-list'
import useSWRMutation from 'swr/mutation'
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
import DegreeTemplateSheet from './degree-template-sheet'

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

  const handleRefetchQueryList = useCallback(() => {
    queryDegreeTemplatesByFaculty.mutate()
  }, [queryDegreeTemplatesByFaculty])

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
            <span className='hidden md:block'>Ký trường</span>
          </Button>,
          <Dialog key='sign-degree-template-faculty' open={openSignDialog} onOpenChange={setOpenSignDialog}>
            <DialogTrigger asChild>
              <Button variant={'outline'}>
                <KeyRound />
                <span className='hidden md:block'>Ký khoa</span>
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Ký số cho khoa</DialogTitle>
                <DialogDescription>Ký số cho các mẫu văn bằng số của khoa</DialogDescription>
              </DialogHeader>
              <Label>Chọn chuyên ngành</Label>
              <CommonSelect
                value={selectFacultyId}
                options={facultyOptions}
                handleSelect={setSelectFacultyId}
                placeholder='Chọn chuyên ngành'
                selectLabel='Chuyên ngành'
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
            <span className='hidden md:block'>Tạo mới</span>
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
                canEdit={!item.isLocked}
                handleSetIdDetail={setIdDetail}
                id={item.id}
                refetch={handleRefetchQueryList}
              />
            )
          }
        ]}
      />
      <DegreeTemplateSheet id={idDetail} onClose={handleCloseDetailDialog} handleRefetch={handleRefetchQueryList} />
    </div>
  )
}

export default DegreeTemplate
