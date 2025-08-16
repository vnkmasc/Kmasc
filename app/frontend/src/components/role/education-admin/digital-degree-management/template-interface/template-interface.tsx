'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateView from '@/components/role/education-admin/digital-degree-management/template-view'
import PageHeader from '@/components/common/page-header'
import useSWR from 'swr'
import { getTemplateInterfaceById, getTemplateInterfaces } from '@/lib/api/digital-degree'
import { useRouter, useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Edit2, Plus } from 'lucide-react'
import { useRef, useState } from 'react'
import TinyTextEdit, { TinyTextEditRef } from '../tiny-text-edit'

const TemplateInterface: React.FC = () => {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [mode, setMode] = useState<'view' | 'edit' | 'create'>('view')
  const tinyEditorRef = useRef<TinyTextEditRef>(null)
  const [tinyValue, setTinyValue] = useState<string | undefined>(undefined)

  const queryTemplateInterfaces = useSWR('example-templates', getTemplateInterfaces)
  const queryTemplateInterfaceById = useSWR(
    queryTemplateInterfaces.data ? 'example-templates-id' + searchParams.get('id') : undefined,
    async () => {
      const res = await getTemplateInterfaceById(searchParams.get('id') ?? queryTemplateInterfaces.data?.data[0].id)
      setTinyValue(res.data.HTMLContent)
      return res
    }
  )

  // useEffect(() => {
  //   if (queryTemplateInterfaces.data) {
  //     queryTemplateInterfaceById.mutate()
  //   }
  // })

  return (
    <>
      <PageHeader
        title='Giao diện mẫu bằng số'
        extra={[
          <Button key='edit' variant={'outline'} onClick={() => setMode('edit')}>
            <Edit2 /> <span className='hidden md:block'>Chỉnh sửa</span>
          </Button>,
          <Button
            key='create'
            onClick={() => {
              setMode('create')
              setTinyValue(undefined)
            }}
          >
            <Plus /> <span className='hidden md:block'>Tạo mới</span>
          </Button>
        ]}
      />
      <Tabs
        className='w-full'
        defaultValue={searchParams.get('id') ?? queryTemplateInterfaces.data?.data?.[0].id}
        onValueChange={(value) => {
          router.push(`/education-admin/digital-degree-management?tab=template-interface&id=${value}`)
          setMode('view')
        }}
      >
        <TabsList>
          {queryTemplateInterfaces.data?.data?.map((template: any) => (
            <TabsTrigger key={template.id} value={`${template.id}`}>
              {template.name}
            </TabsTrigger>
          ))}
        </TabsList>
        {queryTemplateInterfaces.data?.data?.map((template: any) => (
          <TabsContent key={template.id} value={`${template.id}`}>
            {mode !== 'view' ? (
              <div>
                <TinyTextEdit ref={tinyEditorRef} value={tinyValue} onChange={setTinyValue} />
              </div>
            ) : (
              <TemplateView
                baseHtml={queryTemplateInterfaceById.data?.data.HTMLContent}
                htmlLoading={queryTemplateInterfaceById.isLoading}
              />
            )}
          </TabsContent>
        ))}
      </Tabs>
    </>
  )
}

export default TemplateInterface
