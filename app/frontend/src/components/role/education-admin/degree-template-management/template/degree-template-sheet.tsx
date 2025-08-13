import CustomFormItem from '@/components/common/ct-form-item'
import { UseData } from '@/components/providers/data-provider'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { Sheet, SheetClose, SheetContent, SheetFooter, SheetTitle } from '@/components/ui/sheet'
import { createDegreeTemplate, getDegreeTemplateById, updateDegreeTemplate } from '@/lib/api/degree'
import { showNotification } from '@/lib/utils/common'
import { formatFacultyOptionsByID } from '@/lib/utils/format-api'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import useSWR from 'swr'
import useSWRMutation from 'swr/mutation'
import { z } from 'zod'
import HtmlEditView from './html-edit-view'

interface Props {
  id: string | null | undefined
  onClose: () => void
  handleRefetch: () => void
}

const formSchema = z.object({
  name: z.string({ message: 'Tên mẫu bằng số không được để trống' }),
  description: z.string({ message: 'Mô tả mẫu bằng số không được để trống' }),
  html_content: z.string({ message: 'Mẫu bằng số không được để trống' }),
  faculty_id: z.string({ message: 'Chuyên ngành không được để trống' })
})

const DegreeTemplateSheet: React.FC<Props> = (props) => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      description: '',
      html_content: '',
      faculty_id: undefined
    }
  })

  useSWR(
    props.id,
    async () => {
      const res = await getDegreeTemplateById(props.id as string)

      form.setValue('name', res.data?.name)
      form.setValue('description', res.data?.description)
      form.setValue('faculty_id', res.data?.facultyId)
      form.setValue('html_content', res.data?.html_content)
      return res.data
    },
    {
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi lấy thông tin mẫu bằng số')
      }
    }
  )

  const mutateCreateDegreeTemplate = useSWRMutation(
    'create-degree-template',
    (_key, { arg }: { arg: any }) => createDegreeTemplate(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Tạo mẫu bằng số thành công')
        props.onClose()
        props.handleRefetch()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi tạo mẫu bằng số')
      }
    }
  )

  const mutateUpdateDegreeTemplate = useSWRMutation(
    'update-degree-template',
    (_key, { arg }: { arg: any }) => updateDegreeTemplate(props.id as string, arg),
    {
      onSuccess: () => {
        showNotification('success', 'Cập nhật mẫu bằng số thành công')
        props.onClose()
        props.handleRefetch()
      },
      onError: (error) => {
        showNotification('error', error.message || 'Lỗi khi cập nhật mẫu bằng số')
      }
    }
  )

  const handleSubmit = (data: z.infer<typeof formSchema>) => {
    if (props.id) {
      mutateUpdateDegreeTemplate.trigger(data)
    } else {
      mutateCreateDegreeTemplate.trigger(data)
    }
  }

  return (
    <Sheet
      open={props.id !== undefined}
      onOpenChange={(open) => {
        if (!open) {
          props.onClose()
        }
        form.reset()
      }}
    >
      <SheetContent className='min-w-full max-w-full overflow-y-scroll md:min-w-[900px]'>
        <SheetTitle>{props.id ? 'Chỉnh sửa mẫu bằng số' : 'Thêm mẫu bằng số'}</SheetTitle>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className='mt-4'>
            <div className='space-y-4'>
              {' '}
              <CustomFormItem
                type='input'
                name='name'
                control={form.control}
                label='Tên mẫu bằng số'
                placeholder='Nhập tên mẫu bằng số'
              />
              <CustomFormItem
                type='textarea'
                name='description'
                control={form.control}
                label='Mô tả'
                placeholder='Nhập mô tả'
              />
              <CustomFormItem
                type='select'
                name='faculty_id'
                control={form.control}
                label='Chuyên ngành'
                placeholder='Chọn chuyên ngành'
                setting={{
                  select: {
                    groups: [
                      {
                        label: 'Chuyên ngành',
                        options: formatFacultyOptionsByID(UseData().facultyList)
                      }
                    ]
                  }
                }}
                disabled={!!props.id}
              />
              <HtmlEditView
                textarea={
                  <CustomFormItem
                    type='textarea'
                    name='html_content'
                    control={form.control}
                    placeholder='Nhập code mẫu bằng số'
                    label='Code mẫu bằng số'
                    setting={{
                      textarea: {
                        rows: 25
                      }
                    }}
                  />
                }
                html={form.watch('html_content')}
              />
            </div>
            <SheetFooter className='mt-4'>
              <SheetClose asChild>
                <Button variant='outline' type='button'>
                  Hủy bỏ
                </Button>
              </SheetClose>
              <Button
                type='submit'
                isLoading={mutateCreateDegreeTemplate.isMutating || mutateUpdateDegreeTemplate.isMutating}
              >
                {props.id ? 'Cập nhật' : 'Tạo mới'}
              </Button>
            </SheetFooter>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  )
}

export default DegreeTemplateSheet
