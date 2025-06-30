'use client'

import PageHeader from '@/components/common/page-header'
import Filter from '@/components/role/education-admin/filter'
import TableList from '@/components/role/education-admin/table-list'
import { PAGE_SIZE, REWARD_DISCIPLINE_TYPE_SETTING, STUDENT_CODE_SEARCH_SETTING } from '@/constants/common'
import { searchRewardDiscipline } from '@/lib/api/reward-discipline'
import { formatRewardDiscipline } from '@/lib/utils/format-api'

import { useState } from 'react'
import useSWR from 'swr'

const RewardDisciplineManagementPage: React.FC = () => {
  const [filter, setFilter] = useState<any>({})
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)
  console.log('🚀 ~ setIdDetail:', setIdDetail)
  console.log('🚀 ~ idDetail:', idDetail)

  // const handleCloseDetailDialog = useCallback(() => {
  //   setIdDetail(undefined)
  // }, [])

  // const handleChangePage = useCallback(
  //   (page: number) => {
  //     setFilter({ ...filter, page })
  //   },
  //   [filter]
  // )

  const queryRewardDiscipline = useSWR('reward-discipline' + JSON.stringify(filter), () =>
    searchRewardDiscipline({
      ...formatRewardDiscipline(filter, true),
      page: filter.page || 1,
      page_size: PAGE_SIZE
    })
  )
  return (
    <div>
      <PageHeader title='Khen thưởng & Kỷ luật' extra={[]} />
      <Filter
        handleSetFilter={setFilter}
        items={[
          {
            type: 'query_select',
            placeholder: 'Nhập và chọn mã sinh viên',
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
          { header: 'Họ và tên', value: 'name', className: 'min-w-[200px]' },
          { header: 'Mã QĐ', value: 'decisionNumber', className: 'min-w-[150px]' },
          { header: 'Tên QĐ', value: 'name', className: 'min-w-[200px]' },
          { header: 'Loại', value: 'isDiscipline', className: 'min-w-[150px]' },
          { header: 'Mức độ kỷ luật', value: 'level' },
          { header: 'Ngày tạo', value: 'date', className: 'min-w-[150px]' },
          { header: 'Hành động', value: 'action', className: 'min-w-[150px]' }
        ]}
        page={filter.page || 1}
      />
    </div>
  )
}

export default RewardDisciplineManagementPage
