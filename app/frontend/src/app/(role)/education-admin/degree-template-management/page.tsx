'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import ExampleTemplate from '@/components/role/education-admin/degree-template-management/example-template/example-template'
import { Suspense } from 'react'
import DegreeTemplate from '@/components/role/education-admin/degree-template-management/template/degree-template'
import SuspendPage from '@/components/common/suspend-page'

const DegreeTemplateManagement = () => {
  return (
    <Tabs defaultValue='template'>
      <TabsList>
        <TabsTrigger value='template'>Mẫu bằng số</TabsTrigger>
        <TabsTrigger value='example'>Mẫu ví dụ</TabsTrigger>
      </TabsList>

      <TabsContent value='template'>
        <DegreeTemplate />
      </TabsContent>
      <TabsContent value='example'>
        <ExampleTemplate />
      </TabsContent>
    </Tabs>
  )
}

const DigitalDegreeManagementPage = () => {
  return (
    <Suspense fallback={<SuspendPage />}>
      <DegreeTemplateManagement />
    </Suspense>
  )
}

export default DigitalDegreeManagementPage
