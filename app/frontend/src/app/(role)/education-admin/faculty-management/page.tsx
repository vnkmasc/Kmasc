'use client'
import PageHeader from '@/components/common/page-header'
import { UseData, UseRefetchFacultyList } from '@/components/providers/data-provider'
import DetailDialog from '@/components/role/education-admin/detail-dialog'
import TableActionButton from '@/components/role/education-admin/table-action-button'
import TableList from '@/components/role/education-admin/table-list'
import { Button } from '@/components/ui/button'

import { createFaculty, deleteFaculty, getFacultyById, updateFaculty } from '@/lib/api/faculty'
import { showNotification } from '@/lib/utils/common'
import { validateNoEmpty } from '@/lib/utils/validators'
import { PlusIcon } from 'lucide-react'
import { useCallback, useState } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'

const FacultyManagementPage = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)
  const data = UseData().facultyList

  const handleCloseDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])
  const queryFacultyDetail = useSWR(idDetail, () => getFacultyById(idDetail as string), {
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi lấy thông tin chuyên ngành')
    }
  })
  const mutateCreateFaculty = useSWRMutation('create-faculty', (_key, { arg }: { arg: any }) => createFaculty(arg), {
    onSuccess: () => {
      showNotification('success', 'Thêm chuyên ngành thành công')
      UseRefetchFacultyList()
      setIdDetail(undefined)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi thêm chuyên ngành')
    }
  })
  const mutateUpdateFaculty = useSWRMutation(
    'update-faculty',
    (_key, { arg }: { arg: any }) => updateFaculty(idDetail as string, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Cập nhật chuyên ngành thành công')
        UseRefetchFacultyList()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật chuyên ngành')
      }
    }
  )

  const mutateDeleteFaculty = useSWRMutation('delete-faculty', (_key, { arg }: { arg: any }) => deleteFaculty(arg), {
    onSuccess: () => {
      showNotification('success', 'Xóa chuyên ngành thành công')
      UseRefetchFacultyList()
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi xóa chuyên ngành')
    }
  })

  const handleSubmit = useCallback(
    (data: any) => {
      if (idDetail) {
        mutateUpdateFaculty.trigger(data)
      } else {
        mutateCreateFaculty.trigger(data)
      }
    },
    [idDetail, mutateUpdateFaculty, mutateCreateFaculty]
  )

  const handleDelete = useCallback(
    (id: string) => {
      mutateDeleteFaculty.trigger(id)
    },
    [mutateDeleteFaculty]
  )

  return (
    <>
      <PageHeader
        title='Quản lý chuyên ngành'
        extra={[
          <Button key='create-faculty' onClick={() => setIdDetail(null)}>
            <PlusIcon />
            <span className='hidden md:block'>Tạo mới</span>
          </Button>
        ]}
      />
      {/* <Filter
        items={[
          {
            type: 'input',
            placeholder: 'Nhập mã chuyên ngành',
            name: 'code'
          },
          {
            type: 'input',
            placeholder: 'Nhập chuyên ngành',
            name: 'name'
          }
        ]}
        handleSetFilter={setFilter}
      /> */}
      <TableList
        data={data}
        items={[
          { header: 'Mã chuyên ngành', value: 'code', className: 'text-blue-500 font-semibold min-w-[100px]' },
          { header: 'Chuyên ngành', value: 'name', className: 'min-w-[200px]' },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <TableActionButton id={item.id} handleDelete={handleDelete} handleSetIdDetail={setIdDetail} />
            )
          }
        ]}
      />
      <DetailDialog
        items={[
          {
            type: 'input',
            placeholder: 'Nhập mã chuyên ngành',
            name: 'code',
            label: 'Mã chuyên ngành',
            validator: validateNoEmpty('Mã chuyên ngành')
          },
          {
            type: 'input',
            placeholder: 'Nhập chuyên ngành',
            name: 'name',
            label: 'Chuyên ngành',
            validator: validateNoEmpty('Chuyên ngành')
          },
          {
            type: 'textarea',
            placeholder: 'Nhập mô tả',
            name: 'description',
            label: 'Mô tả'
          }
        ]}
        data={queryFacultyDetail.data || {}}
        mode={idDetail ? 'update' : idDetail === null ? 'create' : undefined}
        handleSubmit={handleSubmit}
        handleClose={handleCloseDialog}
      />
    </>
  )
}

export default FacultyManagementPage
