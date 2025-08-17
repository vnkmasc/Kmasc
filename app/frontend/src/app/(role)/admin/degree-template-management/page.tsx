'use client'

import PageHeader from '@/components/common/page-header'
import { Button } from '@/components/ui/button'
import {
  getUniversities,
  getFacultiesByUniversity,
  getTemplatesByUniAndFaculty,
  signTemplateByMinedu
} from '@/lib/api/admin-degree'
import { useState } from 'react'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import TableList from '@/components/common/table-list'
import { Badge } from '@/components/ui/badge'
import { showMessage, showNotification } from '@/lib/utils/common'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CircleXIcon, CodeXml, KeyRound } from 'lucide-react'
import CommonSelect from '@/components/role/education-admin/common-select'
import { DEGREE_TEMPLATE_STATUS } from '@/constants/common'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { getTemplateInterfaceById } from '@/lib/api/digital-degree'
import TemplateView from '@/components/common/template-view'
import { getSignDegreeConfig } from '@/lib/utils/handle-storage'
import { signDigitalSignature } from '@/lib/utils/handle-vgca'

const DegreeTemplateManagementPage = () => {
  const signDegreeConfig = getSignDegreeConfig()
  const [selectedUniversityId, setSelectedUniversityId] = useState<string>('')
  const [selectedFacultyId, setSelectedFacultyId] = useState<string>('')

  const handleReset = () => {
    setSelectedUniversityId('')
    setSelectedFacultyId('')
  }

  const queryUniversities = useSWR('admin-universities', async () => {
    const res = await getUniversities()

    return res.map((u: any) => ({ label: u.university_name, value: u.id }))
  })

  const queryFaculties = useSWR(selectedUniversityId ? `admin-faculties-${selectedUniversityId}` : null, async () => {
    const res = await getFacultiesByUniversity(selectedUniversityId)
    return res.map((u: any) => ({ label: u.faculty_name, value: u.id }))
  })

  const queryDigitalDegreeTemplate = useSWR(
    selectedUniversityId && selectedFacultyId ? `admin-templates-${selectedUniversityId}-${selectedFacultyId}` : null,
    () => getTemplatesByUniAndFaculty(selectedUniversityId as string, selectedFacultyId as string)
  )

  const mutateTemplateInterface = useSWRMutation('admin-html-template', (_, { arg }: { arg: string }) =>
    getTemplateInterfaceById(arg)
  )

  const mutateSignTemplate = useSWRMutation(
    'admin-sign-template',
    (_, { arg }: { arg: any }) => signTemplateByMinedu(arg.template_id, arg.signature),
    {
      onSuccess: () => {
        showNotification('success', 'Ký số mẫu thành công')
        queryDigitalDegreeTemplate.mutate()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Ký số mẫu thất bại')
      }
    }
  )

  return (
    <>
      <PageHeader title='Quản lý mẫu bằng số' />
      <Card>
        <CardHeader>
          <CardTitle>
            <div className='flex items-center justify-between'>
              Tìm kiếm
              <Button variant='destructive' onClick={handleReset}>
                <CircleXIcon />
                <span className='hidden md:block'>Xóa bộ lọc</span>
              </Button>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent className='grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3 md:gap-4 lg:grid-cols-4 xl:grid-cols-5'>
          <CommonSelect
            value={selectedUniversityId}
            handleSelect={setSelectedUniversityId}
            options={queryUniversities.data ?? []}
            selectLabel='Trường đại học'
            placeholder='Chọn trường'
          />
          <CommonSelect
            value={selectedFacultyId}
            handleSelect={setSelectedFacultyId}
            options={queryFaculties.data ?? []}
            selectLabel='Chuyên ngành'
            placeholder='Chọn chuyên ngành'
          />
        </CardContent>
      </Card>
      <TableList
        items={[
          { header: 'Tên mẫu', value: 'name', className: 'font-semibold text-blue-500 min-w-[200px]' },
          { header: 'Mô tả', value: 'description' },
          {
            header: 'Trạng thái',
            value: 'status',
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
              <div className='flex items-center gap-2'>
                <Sheet>
                  <SheetTrigger asChild>
                    <Button
                      variant={'outline'}
                      size={'icon'}
                      onClick={() => mutateTemplateInterface.trigger(item.template_sample_id)}
                    >
                      <CodeXml />
                    </Button>
                  </SheetTrigger>
                  <SheetContent className='min-w-full md:min-w-[1200px]'>
                    <SheetHeader>
                      <SheetTitle>Giao diện mẫu</SheetTitle>
                    </SheetHeader>
                    <TemplateView
                      baseHtml={mutateTemplateInterface.data?.data.HTMLContent}
                      htmlLoading={mutateTemplateInterface.isMutating}
                    />
                  </SheetContent>
                </Sheet>
                <Button
                  size='icon'
                  onClick={async () => {
                    if (signDegreeConfig?.signService === '') {
                      showMessage('Vui lòng cấu hình số cho link server ký số')
                      return
                    }
                    // *@*
                    const signature = await signDigitalSignature(item.hash_template)

                    if (!signature) {
                      showMessage('Ký số không thành công')
                      return
                    }

                    mutateSignTemplate.trigger({ template_id: item.template_sample_id, signature })
                  }}
                  disabled={item.status === 'SIGNED_BY_MINEDU'}
                  isLoading={mutateSignTemplate.isMutating}
                >
                  <KeyRound />
                </Button>
              </div>
            )
          }
        ]}
        data={queryDigitalDegreeTemplate.data ?? []}
      />
    </>
  )
}

export default DegreeTemplateManagementPage
