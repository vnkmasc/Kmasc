'use client'
import PageHeader from '@/components/common/page-header'
import TableList from '@/components/role/education-admin/table-list'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { approveUniversity, getUniversityList, rejectUniversity } from '@/lib/api/university'
import { cn } from '@/lib/utils'

import { PackageCheckIcon, PackageXIcon } from 'lucide-react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import { showNotification } from '@/lib/utils/common'

const EducationManagementPage = () => {
  const queryUniversityList = useSWR('university-list', () => getUniversityList())
  const mutateApproveUniversity = useSWRMutation(
    'approve-university',
    (_, { arg }: { arg: string }) => approveUniversity(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Duyệt trường thành công')
        queryUniversityList.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Duyệt trường thất bại')
      }
    }
  )
  const mutateRejectUniversity = useSWRMutation(
    'reject-university',
    (_, { arg }: { arg: string }) => rejectUniversity(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Từ chối trường thành công')
        queryUniversityList.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Từ chối trường thất bại')
      }
    }
  )

  return (
    <>
      <PageHeader title='Quản lý tài khoản đào tạo' />
      <TableList
        items={[
          { header: 'Mã trường', value: 'university_code', className: 'font-semibold text-blue-500 min-w-[100px]' },
          { header: 'Tên trường', value: 'university_name' },
          { header: 'Địa chỉ', value: 'address' },
          { header: 'Email đào tạo', value: 'email_domain' },
          {
            header: 'Trạng thái',
            value: 'status',
            render: (item) => (
              <Badge
                variant={item.status === 'pending' ? 'outline' : item.status === 'approved' ? 'default' : 'destructive'}
              >
                {item.status === 'pending' ? 'Chờ duyệt' : item.status === 'approved' ? 'Đã duyệt' : 'Từ chối'}
              </Badge>
            )
          },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <div className={cn('flex items-center gap-2', item.status !== 'pending' && 'hidden')}>
                <Button size={'icon'} onClick={() => mutateApproveUniversity.trigger(item.id)}>
                  <PackageCheckIcon />
                </Button>
                <Button size={'icon'} variant={'destructive'} onClick={() => mutateRejectUniversity.trigger(item.id)}>
                  <PackageXIcon />
                </Button>
              </div>
            )
          }
        ]}
        data={queryUniversityList.data ?? []}
      />
    </>
  )
}

export default EducationManagementPage
