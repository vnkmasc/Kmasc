'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateInterface from '@/components/role/education-admin/digital-degree-management/template-interface/template-interface'
import DegreeTemplate from '@/components/role/education-admin/digital-degree-management/template/degree-template'
import DigitalDegreeManagement from '@/components/role/education-admin/digital-degree-management/degree/degree-management'
import { useRouter } from 'next/navigation'

const DigitalDegreeManagementPage = () => {
  const router = useRouter()
  const searchTabParams = new URLSearchParams(window.location.search).get('tab') || 'degree'
  return (
    <Tabs
      defaultValue={searchTabParams}
      onValueChange={(value) => router.push(`/education-admin/digital-degree-management?tab=${value}`)}
    >
      <TabsList>
        <TabsTrigger value='degree'>Quản lý văn bằng</TabsTrigger>
        <TabsTrigger value='template'>Quản lý mẫu bằng</TabsTrigger>
        <TabsTrigger value='template-interface'>Giao diện mẫu bằng</TabsTrigger>
      </TabsList>
      <TabsContent value='degree'>
        <DigitalDegreeManagement />
      </TabsContent>
      <TabsContent value='template'>
        <DegreeTemplate />
      </TabsContent>
      <TabsContent value='template-interface'>
        <TemplateInterface />
      </TabsContent>
    </Tabs>
  )
}

export default DigitalDegreeManagementPage
