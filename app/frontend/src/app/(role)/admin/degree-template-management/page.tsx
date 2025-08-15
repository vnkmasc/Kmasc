'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  getUniversities,
  getFacultiesByUniversity,
  getTemplatesByUniAndFaculty,
  signTemplateByMinedu
} from '@/lib/api/admin-degree'
import { useEffect, useMemo, useState } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import TableList from '@/components/role/education-admin/table-list'
import { Badge } from '@/components/ui/badge'
import { showNotification } from '@/lib/utils/common'
import HtmlView from '@/components/role/education-admin/html-view'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'

const DegreeTemplateManagementPage = () => {
  const [selectedUniversityId, setSelectedUniversityId] = useState<string | undefined>()
  const [selectedFacultyId, setSelectedFacultyId] = useState<string | undefined>()
  const [previewHtml, setPreviewHtml] = useState<string | null>(null)

  const queryUniversities = useSWR('admin-universities', () => getUniversities())

  const queryFaculties = useSWR(selectedUniversityId ? `admin-faculties-${selectedUniversityId}` : null, () =>
    getFacultiesByUniversity(selectedUniversityId as string)
  )

  const queryTemplates = useSWR(
    selectedUniversityId && selectedFacultyId ? `admin-templates-${selectedUniversityId}-${selectedFacultyId}` : null,
    () => getTemplatesByUniAndFaculty(selectedUniversityId as string, selectedFacultyId as string)
  )

  useEffect(() => {
    if (queryUniversities.data && queryUniversities.data.length > 0 && !selectedUniversityId) {
      setSelectedUniversityId(queryUniversities.data[0].id)
    }
  }, [queryUniversities.data, selectedUniversityId])

  useEffect(() => {
    if (queryFaculties.data && queryFaculties.data.length > 0 && !selectedFacultyId) {
      setSelectedFacultyId(queryFaculties.data[0].id)
    }
  }, [queryFaculties.data, selectedFacultyId])

  const mutateSignTemplate = useSWRMutation(
    'admin-sign-template',
    (_, { arg }: { arg: string }) => signTemplateByMinedu(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Ký số mẫu thành công')
        queryTemplates.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Ký số mẫu thất bại')
      }
    }
  )

  const universityOptions = useMemo(() => {
    return (queryUniversities.data || []).map((u: any) => ({ label: u.university_name, value: u.id }))
  }, [queryUniversities.data])

  const facultyOptions = useMemo(() => {
    return (queryFaculties.data || []).map((f: any) => ({ label: f.faculty_name, value: f.id }))
  }, [queryFaculties.data])

  return (
    <>
      <PageHeader title='Quản lý văn bằng số' />

      <div className='flex flex-col gap-3 md:flex-row'>
        <div className='w-full md:w-1/2'>
          <label className='mb-1 block text-sm font-medium'>Trường</label>
          <Select
            value={selectedUniversityId}
            onValueChange={(val) => {
              setSelectedUniversityId(val)
              setSelectedFacultyId(undefined)
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder='Chọn trường' />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Danh sách trường</SelectLabel>
                {universityOptions.map((opt: { label: string; value: string }) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        <div className='w-full md:w-1/2'>
          <label className='mb-1 block text-sm font-medium'>Khoa</label>
          <Select value={selectedFacultyId} onValueChange={setSelectedFacultyId}>
            <SelectTrigger>
              <SelectValue placeholder='Chọn khoa' />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Danh sách khoa</SelectLabel>
                {facultyOptions.map((opt: { label: string; value: string }) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>

      <TableList
        items={[
          { header: 'Tên mẫu', value: 'name', className: 'font-semibold text-blue-500 min-w-[200px]' },
          { header: 'Mô tả', value: 'description' },
          {
            header: 'Trạng thái',
            value: 'status',
            render: (item) => (
              <Badge variant={item.status === 'SIGNED_BY_UNI' ? 'default' : 'outline'}>
                {item.status === 'SIGNED_BY_UNI'
                  ? 'Đã ký bởi Trường'
                  : item.status === 'SIGNED_BY_MINEDU'
                    ? 'Đã ký bởi Bộ'
                    : item.status}
              </Badge>
            )
          },
          {
            header: 'Hành động',
            value: 'action',
            render: (item) => (
              <div className='flex items-center gap-2'>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant={'secondary'} size={'sm'} onClick={() => setPreviewHtml(item.html_content)}>
                      Xem mẫu
                    </Button>
                  </DialogTrigger>
                  <DialogContent className='max-w-5xl'>
                    <DialogHeader>
                      <DialogTitle>Xem trước mẫu</DialogTitle>
                    </DialogHeader>
                    <div className='max-h-[80vh] overflow-y-auto'>
                      <HtmlView html={previewHtml || item.html_content} />
                    </div>
                  </DialogContent>
                </Dialog>
                <Button size={'sm'} onClick={() => mutateSignTemplate.trigger(item.id)}>
                  Ký số
                </Button>
              </div>
            )
          }
        ]}
        data={queryTemplates.data ?? []}
      />
    </>
  )
}

export default DegreeTemplateManagementPage
