'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateInterface from '@/components/role/education-admin/digital-degree-management/template-interface/template-interface'
import DegreeTemplate from '@/components/role/education-admin/digital-degree-management/template/degree-template'
import { useRouter, useSearchParams } from 'next/navigation'
import DegreeManagement from '@/components/role/education-admin/digital-degree-management/degree/degree-management'
import { Suspense } from 'react'
import SuspendPage from '@/components/common/suspend-page'
import UseBreakpoint from '@/lib/hooks/use-breakpoint'

const DigitalDegreeManagement = () => {
  const { sm } = UseBreakpoint()
  const router = useRouter()
  const searchTabParams = useSearchParams()
  return (
    <Tabs
      defaultValue={searchTabParams.get('tab') || 'degree'}
      onValueChange={(value) => router.push(`/education-admin/digital-degree-management?tab=${value}`)}
    >
      <TabsList>
        <TabsTrigger value='degree'>{sm ? 'Quản lý văn bằng' : 'QLVB'} </TabsTrigger>
        <TabsTrigger value='template'>{sm ? 'Quản lý mẫu bằng' : 'QLMB'}</TabsTrigger>
        <TabsTrigger value='template-interface'>{sm ? 'Giao diện mẫu bằng' : 'GDMB'}</TabsTrigger>
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
