'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TemplateView from '@/components/common/template-view'
import PageHeader from '@/components/common/page-header'
import useSWR from 'swr'
import {
  createTemplateInterface,
  getTemplateInterfaceById,
  getTemplateInterfaces,
  updateTemplateInterface
} from '@/lib/api/digital-degree'
import { useRouter, useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { ChevronsLeftRightEllipsis, CircleX, Edit2, Plus, Save } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import TinyTextEdit, { TinyTextEditRef } from '../tiny-text-edit'
import useSWRMutation from 'swr/mutation'
import { cn, showNotification } from '@/lib/utils/common'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import HtmlEditView from '../template/html-edit-view'
import { Textarea } from '@/components/ui/textarea'
import CommonSelect from '../../common-select'

const TemplateInterface: React.FC = () => {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [mode, setMode] = useState<'view' | 'edit' | 'create'>('view')
  const tinyEditorRef = useRef<TinyTextEditRef>(null)
  const [tinyValue, setTinyValue] = useState<string>('')
  const [templateName, setTemplateName] = useState<string>('')
  const [tinyMode, setTinyMode] = useState<boolean>(true)
  const [selectedTemplateId, setSelectedTemplateId] = useState<string>('')

  useEffect(() => {
    const fetchTemplateInterface = async () => {
      const res = await getTemplateInterfaceById(selectedTemplateId)
      if (res) {
        setTinyValue(res.data.HTMLContent)
      }
    }
    if (selectedTemplateId !== '') {
      fetchTemplateInterface()
    }
  }, [selectedTemplateId])

  const queryTemplateInterfaces = useSWR('sample-templates', getTemplateInterfaces)
  const queryTemplateInterfaceById = useSWR(
    queryTemplateInterfaces.data ? 'sample-templates-id' + searchParams.get('id') : undefined,
    async () => {
      const res = await getTemplateInterfaceById(searchParams.get('id') ?? queryTemplateInterfaces.data?.data[0].id)
      setTinyValue(res.data.HTMLContent)
      setTemplateName(res.data.Name)
      return res
    }
  )

  const mutateUpdateTemplateInterface = useSWRMutation(
    'update-template-interface',
    (_, { arg }: { arg: any }) =>
      updateTemplateInterface(searchParams.get('id') ?? queryTemplateInterfaces.data?.data?.[0].id, arg),
    {
      onSuccess: () => {
        queryTemplateInterfaceById.mutate()
        showNotification('success', 'Cập nhật giao diện mẫu thành công')
        setMode('view')
      },
      onError: (error: any) => {
        showNotification('error', error.message || 'Cập nhật giao diện mẫu thất bại')
      }
    }
  )

  const mutateCreateTemplateInterface = useSWRMutation(
    'create-template-interface',
    (_, { arg }: { arg: any }) => createTemplateInterface(arg),
    {
      onSuccess: () => {
        queryTemplateInterfaces.mutate()
        showNotification('success', 'Tạo giao diện mẫu thành công')
        setMode('view')
        setSelectedTemplateId('')
      },
      onError: (error: any) => {
        showNotification('error', error.message || 'Tạo giao diện mẫu thất bại')
      }
    }
  )

  const handleSubmit = () => {
    if (mode === 'edit') {
      mutateUpdateTemplateInterface.trigger({
        name: templateName,
        html_content: tinyValue
      })
    } else {
      mutateCreateTemplateInterface.trigger({
        name: templateName,
        html_content: tinyValue
      })
    }
  }

  useEffect(() => {
    if (mode === 'create') {
      setTinyValue('')
      setTemplateName('')
    }
    setTinyMode(true)
  }, [mode])

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
            }}
          >
            <Plus /> <span className='hidden md:block'>Tạo mới</span>
          </Button>
        ]}
      />
      {queryTemplateInterfaces.data && (
        <Tabs
          className='w-full'
          defaultValue={searchParams.get('id') ?? queryTemplateInterfaces.data?.data?.[0].id}
          onValueChange={(value) => {
            router.push(`/education-admin/digital-degree-management?tab=template-interface&id=${value}`)
            setMode('view')
          }}
        >
          <TabsList className={cn(mode === 'create' && 'hidden')}>
            {queryTemplateInterfaces.data?.data?.map((template: any) => (
              <TabsTrigger key={template.id} value={`${template.id}`}>
                {template.name}
              </TabsTrigger>
            ))}
          </TabsList>
          {queryTemplateInterfaces.data?.data?.map((template: any) => (
            <TabsContent key={template.id} value={`${template.id}`}>
              {mode !== 'view' ? (
                <Card>
                  <CardHeader className='flex flex-row items-center justify-between'>
                    <div>
                      {' '}
                      <CardTitle className='mb-1'>
                        {mode === 'edit' ? 'Chỉnh sửa giao diện mẫu' : 'Tạo giao diện mẫu'}
                      </CardTitle>
                      <CardDescription>
                        Chỉnh sửa mẫu giao diện bằng <strong>chế độ soạn thảo</strong> hoặc{' '}
                        <strong>nhập mã code</strong>
                      </CardDescription>
                    </div>
                    <div className='flex items-center gap-2'>
                      <Button variant={'destructive'} onClick={() => setMode('view')}>
                        <CircleX />
                        <span className='hidden md:block'>Hủy bỏ</span>
                      </Button>
                      <Button variant={'secondary'} onClick={handleSubmit}>
                        <Save />
                        <span className='hidden md:block'>{mode === 'edit' ? 'Cập nhật' : 'Tạo mới'}</span>
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className='flex flex-col gap-4 md:flex-row'>
                      {mode === 'create' && (
                        <div className='min-w-56'>
                          <Label htmlFor='template-html'>Chọn mẫu có sẵn</Label>
                          <CommonSelect
                            value={selectedTemplateId}
                            handleSelect={setSelectedTemplateId}
                            placeholder='Chọn mẫu'
                            options={queryTemplateInterfaces.data?.data?.map((item: any) => ({
                              label: item.name,
                              value: item.id
                            }))}
                          />
                        </div>
                      )}
                      <div>
                        {' '}
                        <Label htmlFor='template-name'>Tên giao diện mẫu</Label>
                        <Input
                          id='template-name'
                          placeholder='Nhập tên giao diện mẫu'
                          value={templateName}
                          onChange={(e) => setTemplateName(e.target.value)}
                          className='min-w-56'
                        />
                      </div>
                    </div>
                    <div className='mt-4 flex items-center gap-2'>
                      <Label>Mẫu bằng số</Label>
                      <ChevronsLeftRightEllipsis />
                      <span className='text-sm font-semibold'>Chế độ soạn thảo</span>
                      <Switch onCheckedChange={setTinyMode} checked={tinyMode} />
                    </div>
                    {tinyMode ? (
                      <TinyTextEdit ref={tinyEditorRef} value={tinyValue} onChange={setTinyValue} />
                    ) : (
                      <HtmlEditView
                        textarea={
                          <Textarea rows={25} value={tinyValue} onChange={(e) => setTinyValue(e.target.value)} />
                        }
                        html={tinyValue}
                      />
                    )}
                  </CardContent>
                </Card>
              ) : (
                <TemplateView
                  baseHtml={queryTemplateInterfaceById.data?.data.HTMLContent}
                  htmlLoading={queryTemplateInterfaceById.isLoading}
                />
              )}
            </TabsContent>
          ))}
        </Tabs>
      )}
    </>
  )
}

export default TemplateInterface
