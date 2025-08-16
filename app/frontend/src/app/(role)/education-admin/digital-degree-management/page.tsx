'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateInterface from '@/components/role/education-admin/digital-degree-management/template-interface/template-interface'
import DegreeTemplate from '@/components/role/education-admin/digital-degree-management/template/degree-template'
import { useRouter, useSearchParams } from 'next/navigation'
import DegreeManagement from '@/components/role/education-admin/digital-degree-management/degree/degree-management'
import { Suspense } from 'react'
import SuspendPage from '@/components/common/suspend-page'

const DigitalDegreeManagement = () => {
  const router = useRouter()
  const searchTabParams = useSearchParams()
  return (
    <Tabs
      defaultValue={searchTabParams.get('tab') || 'degree'}
      onValueChange={(value) => router.push(`/education-admin/digital-degree-management?tab=${value}`)}
    >
      <TabsList>
        <TabsTrigger value='degree'>Quản lý văn bằng</TabsTrigger>
        <TabsTrigger value='template'>Quản lý mẫu bằng</TabsTrigger>
        <TabsTrigger value='template-interface'>Giao diện mẫu bằng</TabsTrigger>
      </TabsList>
      <TabsContent value='degree'>
        <DegreeManagement />
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

const DigitalDegreeManagementPage = () => (
  <Suspense fallback={<SuspendPage />}>
    <DigitalDegreeManagement />
  </Suspense>
)

export default DigitalDegreeManagementPage
