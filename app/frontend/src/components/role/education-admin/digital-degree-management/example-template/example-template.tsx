'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateView from '@/components/role/education-admin/digital-degree-management/template-view'
import PageHeader from '@/components/common/page-header'

const ExampleTemplate: React.FC = () => {
  return (
    <>
      <PageHeader title='Tạo mẫu bằng số' />
      <Tabs defaultValue='v1' className='w-full'>
        <TabsList>
          <TabsTrigger value='v1'>Mẫu 1</TabsTrigger>
          <TabsTrigger value='v2'>Mẫu 2</TabsTrigger>
          <TabsTrigger value='v3'>Mẫu 3</TabsTrigger>
        </TabsList>
        <TabsContent value='v1'>
          <TemplateView baseHtml='' />
        </TabsContent>
        <TabsContent value='v2'>
          <TemplateView baseHtml='' />
        </TabsContent>
        <TabsContent value='v3'>
          <TemplateView baseHtml='' />
        </TabsContent>
      </Tabs>
    </>
  )
}

export default ExampleTemplate
