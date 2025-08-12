'use client'

import DegreeTemplate from '@/components/role/education-admin/digital-degree-management/template/degree-template'
import DigitalDegreeView from '@/components/role/education-admin/digital-degree-management/degree/digital-degree-view'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useRouter, useSearchParams } from 'next/navigation'
import ExampleTemplate from '@/components/role/education-admin/digital-degree-management/example-template/example-template'
import { Suspense } from 'react'

const DigitalDegreeManagementContent = () => {
  const router = useRouter()
  const searchParams = useSearchParams()

  const handleChangeTabs = (tabs: string) => {
    router.push(`/education-admin/digital-degree-management?tab=${tabs}`)
  }

  return (
    <Tabs value={searchParams.get('tab') ?? 'degree'} onValueChange={handleChangeTabs}>
      <TabsList>
        <TabsTrigger value='degree'>Văn bằng số</TabsTrigger>
        <TabsTrigger value='template'>Mẫu bằng số</TabsTrigger>
        <TabsTrigger value='example'>Mẫu ví dụ</TabsTrigger>
      </TabsList>
      <TabsContent value='degree'>
        <DigitalDegreeView />
      </TabsContent>
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
    <Suspense fallback={<div>Loading...</div>}>
      <DigitalDegreeManagementContent />
    </Suspense>
  )
}

export default DigitalDegreeManagementPage
