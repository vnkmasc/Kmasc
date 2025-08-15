'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateView from '@/components/role/education-admin/degree-template-management/template-view'
import { useEffect, useState } from 'react'
import PageHeader from '@/components/common/page-header'

const ExampleTemplate: React.FC = () => {
  const [v1, setV1] = useState<string>('')
  const [v2, setV2] = useState<string>('')
  const [v3, setV3] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(true)

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true)
        const [r1, r2, r3] = await Promise.all([
          fetch('/api/templates/v1-degree').then((r) => r.text()),
          fetch('/api/templates/v2-degree').then((r) => r.text()),
          fetch('/api/templates/v3-degree').then((r) => r.text())
        ])
        setV1(r1)
        setV2(r2)
        setV3(r3)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  return (
    <>
      <PageHeader title='Danh sách mẫu ví dụ' />
      <Tabs defaultValue='v1' className='w-full'>
        <TabsList>
          <TabsTrigger value='v1'>Mẫu 1</TabsTrigger>
          <TabsTrigger value='v2'>Mẫu 2</TabsTrigger>
          <TabsTrigger value='v3'>Mẫu 3</TabsTrigger>
        </TabsList>
        <TabsContent value='v1'>
          <TemplateView baseHtml={v1} htmlLoading={loading} />
        </TabsContent>
        <TabsContent value='v2'>
          <TemplateView baseHtml={v2} htmlLoading={loading} />
        </TabsContent>
        <TabsContent value='v3'>
          <TemplateView baseHtml={v3} htmlLoading={loading} />
        </TabsContent>
      </Tabs>
    </>
  )
}

export default ExampleTemplate
