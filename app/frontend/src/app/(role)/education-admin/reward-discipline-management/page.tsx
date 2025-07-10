'use client'

import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import DetailDialog from '@/components/role/education-admin/detail-dialog'
import Filter from '@/components/role/education-admin/filter'
import TableActionButton from '@/components/role/education-admin/table-action-button'
import TableList from '@/components/role/education-admin/table-list'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  LEVEL_DISCIPLINE,
  PAGE_SIZE,
  REWARD_DISCIPLINE_LEVEL_SETTING,
  REWARD_DISCIPLINE_TYPE_SETTING,
  STUDENT_CODE_SEARCH_SETTING
} from '@/constants/common'
import {
  createRewardDiscipline,
  deleteRewardDiscipline,
  getRewardDisciplineById,
  searchRewardDiscipline,
  updateRewardDiscipline
} from '@/lib/api/reward-discipline'
import { cn } from '@/lib/utils'
import { showNotification } from '@/lib/utils/common'
import { formatRewardDiscipline } from '@/lib/utils/format-api'
import { validateNoEmpty } from '@/lib/utils/validators'
import { PlusIcon } from 'lucide-react'

import { useCallback, useState } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'

const RewardDisciplineManagementPage: React.FC = () => {
  const [filter, setFilter] = useState<any>({})
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)

  const queryRewardDisciplineDetail = useSWR(idDetail, () => getRewardDisciplineById(idDetail as string), {
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi lấy thông tin khen thưởng & kỷ luật')
    }
  })

  const handleChangePage = useCallback(
    (page: number) => {
      setFilter({ ...filter, page })
    },
    [filter]
  )

  const queryRewardDiscipline = useSWR('reward-discipline' + JSON.stringify(filter), () =>
    searchRewardDiscipline({
      ...formatRewardDiscipline(filter, true),
      page: filter.page || 1,
      page_size: PAGE_SIZE
    })
  )

  const mutateCreateRewardDiscipline = useSWRMutation(
    'create-reward-discipline',
    (_key, { arg }: { arg: any }) => createRewardDiscipline(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Thêm khen thưởng & kỷ luật thành công')
        queryRewardDiscipline.mutate()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi thêm khen thưởng & kỷ luật')
      }
    }
  )

  const mutateUpdateRewardDiscipline = useSWRMutation(
    'update-reward-discipline',
    (_key, { arg }: { arg: any }) => updateRewardDiscipline(idDetail as string, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Cập nhật khen thưởng & kỷ luật thành công')
        queryRewardDiscipline.mutate()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật khen thưởng & kỷ luật')
      }
    }
  )

  const mutateDeleteRewardDiscipline = useSWRMutation(
    'delete-reward-discipline',
    (_key, { arg }: { arg: any }) => deleteRewardDiscipline(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Xóa khen thưởng & kỷ luật thành công')
        queryRewardDiscipline.mutate()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi xóa khen thưởng & kỷ luật')
      }
    }
  )

  const handleDelete = useCallback(
    (id: string) => {
      mutateDeleteRewardDiscipline.trigger(id)
    },
    [mutateDeleteRewardDiscipline]
  )

  const handleSubmit = useCallback(
    (data: any) => {
      if (idDetail) {
        mutateUpdateRewardDiscipline.trigger(data)
      } else {
        mutateCreateRewardDiscipline.trigger(data)
      }
    },
    [idDetail, mutateUpdateRewardDiscipline, mutateCreateRewardDiscipline]
  )

  const handleCloseDetailDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])

  return (
    <div>
      <PageHeader
        title='Khen thưởng & Kỷ luật'
        extra={[
          <Button key='add-reward-discipline' onClick={() => setIdDetail(null)}>
            <PlusIcon />
            Thêm KT & KL
          </Button>
        ]}
      />
      <Filter
        handleSetFilter={setFilter}
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn MSV',
            name: 'studentCode',
            setting: STUDENT_CODE_SEARCH_SETTING
          },
          {
            type: 'input',
            placeholder: 'Nhập mã QĐ',
            name: 'decisionNumber'
          },
          {
            type: 'select',
            placeholder: 'Chọn loại',
            name: 'isDiscipline',
            setting: REWARD_DISCIPLINE_TYPE_SETTING
          }
        ]}
      />
      <TableList
        data={queryRewardDiscipline.data?.data || []}
        items={[
          { header: 'Mã SV', value: 'studentCode', className: 'min-w-[80px] font-semibold text-blue-500' },
          { header: 'Họ và tên', value: 'studentName', className: 'min-w-[200px]' },
          { header: 'Mã QĐ', value: 'decisionNumber', className: 'min-w-[150px]' },
          { header: 'Tên QĐ', value: 'name', className: 'min-w-[200px] max-w-[250px]' },
          {
            header: 'Loại',
            value: 'isDiscipline',
            className: 'min-w-[100px]',
            render: (item) =>
              item.isDiscipline ? <Badge variant='destructive'>Kỷ luật</Badge> : <Badge>Khen thưởng</Badge>
          },
          {
            header: 'Mức độ kỷ luật',
            value: 'disciplineLevel',
            className: 'min-w-[150px]',
            render: (item) =>
              item.isDiscipline && (
                <Badge
                  className={cn(
                    item.disciplineLevel === '4' && 'bg-red-500 hover:bg-red-400',
                    item.disciplineLevel === '3' && 'bg-orange-500 hover:bg-orange-400',
                    item.disciplineLevel === '2' && 'bg-yellow-500 hover:bg-yellow-400',
                    item.disciplineLevel === '1' && 'bg-blue-500 hover:bg-blue-400'
                  )}
                >
                  {LEVEL_DISCIPLINE[Number(item.disciplineLevel) as keyof typeof LEVEL_DISCIPLINE]}
                </Badge>
              )
          },
          { header: 'Ngày tạo', value: 'createdAt', className: 'min-w-[150px]' },
          {
            header: 'Hành động',
            value: 'action',
            className: 'min-w-[90px]',
            render: (item) => (
              <TableActionButton handleDelete={handleDelete} handleSetIdDetail={setIdDetail} id={item.id} />
            )
          }
        ]}
        page={filter.page || 1}
      />
      <CommonPagination
        page={queryRewardDiscipline.data?.page || 1}
        totalPage={queryRewardDiscipline.data?.total_page || 1}
        handleChangePage={handleChangePage}
      />
      <DetailDialog
        mode={idDetail ? 'update' : idDetail === undefined ? undefined : 'create'}
        items={[
          {
            type: 'query_select',
            name: 'studentCode',
            placeholder: 'Nhập và chọn MSV',
            label: 'Mã sinh viên',
            setting: STUDENT_CODE_SEARCH_SETTING,
            validator: validateNoEmpty('Mã sinh viên')
          },
          {
            type: 'input',
            name: 'decisionNumber',
            placeholder: 'Nhập mã QĐ',
            label: 'Mã QĐ',
            validator: validateNoEmpty('Mã QĐ')
          },
          {
            type: 'input',
            name: 'name',
            placeholder: 'Nhập tên QĐ',
            label: 'Tên QĐ',
            validator: validateNoEmpty('Tên QĐ')
          },
          {
            type: 'select',
            name: 'disciplineLevel',
            label: 'Mức độ kỷ luật',
            placeholder: 'Chọn loại kỷ luật',
            setting: REWARD_DISCIPLINE_LEVEL_SETTING,
            description: 'Nếu là khen thưởng thì không cần chọn trường này'
          },
          {
            type: 'textarea',
            name: 'description',
            placeholder: 'Nhập mô tả',
            label: 'Mô tả'
          }
        ]}
        data={queryRewardDisciplineDetail.data || []}
        handleSubmit={handleSubmit}
        handleClose={handleCloseDetailDialog}
      />
    </div>
  )
}

export default RewardDisciplineManagementPage
