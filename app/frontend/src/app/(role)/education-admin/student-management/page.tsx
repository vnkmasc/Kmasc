'use client'
import PageHeader from '@/components/common/page-header'
import CommonPagination from '@/components/common/pagination'
import { UseData } from '@/components/providers/data-provider'
import DetailDialog from '@/components/role/education-admin/detail-dialog'
import Filter from '@/components/common/filter'
import TableActionButton from '@/components/role/education-admin/table-action-button'
import TableList from '@/components/common/table-list'
import UploadButton from '@/components/role/education-admin/upload-button'
import { Button } from '@/components/ui/button'
import { GENDER_SELECT_SETTING, PAGE_SIZE, STUDENT_STATUS_OPTIONS } from '@/constants/common'

import {
  createStudent,
  deleteStudent,
  getStudentById,
  importExcel,
  searchStudent,
  updateStudent
} from '@/lib/api/student'
import { formatResponseImportExcel, showNotification } from '@/lib/utils/common'
import { formatDateISO, formatFacultyOptions, formatStudent } from '@/lib/utils/format-api'

import { validateAcademicEmail, validateCitizenId, validateNoEmpty } from '@/lib/utils/validators'
import { PlusIcon } from 'lucide-react'
import Link from 'next/link'
import { useCallback, useState } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import z from 'zod'
import { Badge } from '@/components/ui/badge'

const StudentManagementPage: React.FC = () => {
  const [idDetail, setIdDetail] = useState<string | null | undefined>(undefined)

  const [filter, setFilter] = useState<any>({})
  const handleCloseDetailDialog = useCallback(() => {
    setIdDetail(undefined)
  }, [])

  const handleChangePage = useCallback(
    (page: number) => {
      setFilter({ ...filter, page })
    },
    [filter]
  )

  const queryStudents = useSWR('students' + JSON.stringify(filter), () =>
    searchStudent({
      ...formatStudent(filter, true),
      page: filter.page || 1,
      page_size: PAGE_SIZE
    })
  )

  const queryStudentDetail = useSWR(idDetail, () => getStudentById(idDetail as string), {
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi lấy thông tin sinh viên')
    }
  })

  const mutateCreateStudent = useSWRMutation('create-student', (_key, { arg }: { arg: any }) => createStudent(arg), {
    onSuccess: () => {
      showNotification('success', 'Thêm sinh viên thành công')
      queryStudents.mutate()
      setIdDetail(undefined)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi thêm sinh viên')
    }
  })

  const mutateUpdateStudent = useSWRMutation(
    'update-student',
    (_key, { arg }: { arg: any }) => updateStudent(idDetail as string, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Cập nhật sinh viên thành công')
        queryStudents.mutate()
        setIdDetail(undefined)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật sinh viên')
      }
    }
  )

  const mutateDeleteStudent = useSWRMutation('delete-student', (_key, { arg }: { arg: any }) => deleteStudent(arg), {
    onSuccess: () => {
      showNotification('success', 'Xóa sinh viên thành công')
      queryStudents.mutate()
      setIdDetail(undefined)
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi xóa sinh viên')
    }
  })

  const mutateImportExcel = useSWRMutation('import-excel', (_key, { arg }: { arg: any }) => importExcel(arg), {
    onSuccess: (data) => {
      const formatData = formatResponseImportExcel(data)

      if (data.error_count === 0) {
        showNotification('success', `Thêm ${data.success_count} sinh viên thành công`)
        queryStudents.mutate()
        return
      }

      if (data.success_count === 0) {
        formatData.error.forEach((item) => {
          showNotification('error', `Sinh viên hàng thứ ${item.row.join(', ')} có lỗi: "${item.title}"`)
        })
        return
      }

      formatData.error.forEach((item) => {
        showNotification('error', `Sinh viên hàng thứ ${item.row.join(', ')} có lỗi: "${item.title}" `)
      })

      showNotification('success', `Sinh viên hàng thứ ${formatData.success.join(', ')} đã được thêm thành công`)
      queryStudents.mutate()
    },
    onError: (error) => {
      showNotification('error', error.message || 'Lỗi khi tải tệp lên')
    }
  })
  const handleDelete = useCallback(
    (id: string) => {
      mutateDeleteStudent.trigger(id)
    },
    [mutateDeleteStudent]
  )

  const handleSubmit = useCallback(
    (data: any) => {
      if (idDetail) {
        mutateUpdateStudent.trigger(data)
      } else {
        mutateCreateStudent.trigger(data)
      }
    },
    [idDetail, mutateUpdateStudent, mutateCreateStudent]
  )

  const handleUpload = useCallback(
    (file: FormData) => {
      mutateImportExcel.trigger(file)
    },
    [mutateImportExcel]
  )

  return (
    <div>
      <PageHeader
        title='Quản lý sinh viên'
        extra={[
          <UploadButton key='upload-excel' handleUpload={handleUpload} loading={mutateImportExcel.isMutating} />,
          <Button key='create-student' onClick={() => setIdDetail(null)}>
            <PlusIcon />
            <span className='hidden md:block'>Tạo mới</span>
          </Button>
        ]}
      />

      <Filter
        handleSetFilter={setFilter}
        items={[
          { type: 'input', name: 'code', placeholder: 'Nhập mã sinh viên' },
          {
            type: 'input',
            name: 'name',
            placeholder: 'Nhập họ và tên'
          },
          {
            type: 'input',
            name: 'citizenId',
            placeholder: 'Nhập CCCD'
          },
          {
            type: 'select',
            name: 'faculty',
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
            type: 'input',
            name: 'year',
            placeholder: 'Nhập khóa'
            // setting: {
            //   input: {
            //     type: 'number'
            //   }
            // }
          },

          {
            type: 'select',
            name: 'status',
            placeholder: 'Chọn trạng thái',
            setting: {
              select: {
                groups: [
                  {
                    label: undefined,
                    options: STUDENT_STATUS_OPTIONS
                  }
                ]
              }
            }
          }
        ]}
      />
      <TableList
        data={queryStudents.data?.data || []}
        items={[
          {
            header: 'Mã SV',
            value: 'code',
            className: 'min-w-[80px] font-semibold text-blue-500',
            render: (item) => <Link href={`/education-admin/student-management/${item.id}`}>{item.code}</Link>
          },
          { header: 'Họ và tên', value: 'name', className: 'min-w-[200px]' },
          // { header: 'Email', value: 'email', className: 'min-w-[200px]' },
          {
            header: 'Giới tính',
            value: 'gender',
            className: 'min-w-[150px]',
            render: (item) => (
              <Badge variant={item.gender === 'true' ? 'default' : 'secondary'}>
                {item.gender === 'true' ? 'Nam' : 'Nữ'}
              </Badge>
            )
          },
          {
            header: 'Ngày sinh',
            value: 'dateOfBirth',
            className: 'min-w-[150px]',
            render: (item) => formatDateISO(item.dateOfBirth)
          },
          { header: 'Chuyên ngành', value: 'facultyName', className: 'min-w-[200px]' },
          { header: 'Khóa', value: 'year', className: 'min-w-[150px]' },

          // {
          //   header: 'Trạng thái',
          //   value: 'status',
          //   className: 'min-w-[150px]',
          //   render: (item) => (
          //     <Badge variant={item.status === 'true' ? 'default' : 'secondary'}>
          //       {item.status === 'true' ? 'Đã tốt nghiệp' : 'Đang học'}
          //     </Badge>
          //   )
          // },
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
        page={queryStudents.data?.page || 1}
        totalPage={queryStudents.data?.total_page || 1}
        handleChangePage={handleChangePage}
      />
      <DetailDialog
        items={[
          {
            type: 'input',
            label: 'Mã sinh viên',
            name: 'code',
            validator: z.string().trim().nonempty({
              message: 'Mã sinh viên không được để trống'
            }),
            placeholder: 'VD: CT060111'
          },
          {
            type: 'input',
            label: 'Họ và tên',
            name: 'name',
            validator: z.string().trim().nonempty({
              message: 'Tên sinh viên không được để trống'
            }),
            placeholder: 'VD: Nguyễn Văn A'
          },
          {
            type: 'input',
            label: 'Ngày sinh',
            name: 'dateOfBirth',
            setting: {
              input: {
                type: 'date'
              }
            },
            validator: validateNoEmpty('Ngày sinh')
          },
          {
            type: 'select',
            label: 'Giới tính',
            name: 'gender',
            setting: GENDER_SELECT_SETTING,
            validator: validateNoEmpty('Giới tính')
          },
          {
            type: 'input',
            label: 'Căn cước công dân',
            name: 'citizenId',
            validator: validateCitizenId,
            placeholder: 'VD: 023456789012'
          },
          {
            type: 'input',
            label: 'Dân tộc',
            name: 'ethnicity',
            placeholder: 'VD: Kinh'
          },
          {
            type: 'input',
            label: 'Địa chỉ hiện tại',
            name: 'currentAddress',
            placeholder: 'VD: 123 Đường ABC, Quận XYZ, TP. HCM'
          },
          {
            type: 'input',
            label: 'Nơi sinh',
            name: 'birthAddress',
            placeholder: 'VD: 123 Đường ABC, Quận XYZ, TP. HCM'
          },
          {
            type: 'input',
            label: 'Ngày tham gia Đoàn',
            name: 'unionJoinDate',
            setting: {
              input: {
                type: 'date'
              }
            },

            placeholder: 'VD: 01/01/2025'
          },
          {
            type: 'input',
            label: 'Ngày tham gia Đảng',
            name: 'partyJoinDate',
            setting: {
              input: {
                type: 'date'
              }
            },

            placeholder: 'VD: 01/01/2025'
          },
          {
            type: 'input',
            label: 'Email',
            name: 'email',
            validator: validateAcademicEmail,
            placeholder: 'VD: CT060111@actvn.edu.vn'
          },
          {
            type: 'select',
            label: 'Chuyên ngành',
            placeholder: 'Chọn chuyên ngành',
            name: 'faculty',
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
            type: 'input',
            label: 'Khóa học',
            name: 'year'
            // setting: {
            //   input: {
            //     type: 'number'
            //   }
            // },
          },
          {
            type: 'textarea',
            name: 'description',
            placeholder: 'Nhập mô tả',
            label: 'Mô tả'
          },
          {
            type: 'select',
            label: 'Trạng thái',
            name: 'status',
            disabled: true,
            setting: {
              select: {
                groups: [
                  {
                    label: 'Trạng thái',
                    options: STUDENT_STATUS_OPTIONS
                  }
                ]
              }
            },
            placeholder: 'Không thể chỉnh sửa - Mặc định "Đang học"'
          }
        ]}
        data={queryStudentDetail.data || {}}
        handleSubmit={handleSubmit}
        mode={idDetail ? 'update' : idDetail === undefined ? undefined : 'create'}
        handleClose={handleCloseDetailDialog}
      />
    </div>
  )
}

export default StudentManagementPage
