'use client'
import { Button } from '@/components/ui/button'
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Checkbox } from '@/components/ui/checkbox'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import useSWRMutation from 'swr/mutation'
import { createVerifyCode } from '@/lib/api/certificate'

import CustomFormItem from '@/components/common/ct-form-item'
import { mutate } from 'swr'
import { showNotification } from '@/lib/utils/common'
import { Dispatch, SetStateAction } from 'react'

const permissionType = [
  {
    id: 'can_view_score',
    label: 'Điểm'
  },
  {
    id: 'can_view_data',
    label: 'Dữ liệu văn bằng'
  },
  {
    id: 'can_view_file',
    label: 'Tệp văn bằng'
  }
] as const

const formSchema = z.object({
  expiredAfter: z.number().min(1, { message: 'Thời gian tạo mã xác minh phải lớn hơn 0' }).max(100000, {
    message: 'Thời gian tạo mã xác minh phải nhỏ hơn 100000'
  }),
  permissionType: z.array(z.enum(['can_view_score', 'can_view_data', 'can_view_file'])).min(1, {
    message: 'Phải chọn ít nhất 1 quyền hạn'
  })
})

type FormSchema = z.infer<typeof formSchema>

interface Props {
  open: boolean
  handleSetOpen: Dispatch<SetStateAction<boolean>>
  swrKey: string
}

const CreateVerifyCodeDialog: React.FC<Props> = (props) => {
  const form = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      expiredAfter: 1,
      permissionType: ['can_view_score', 'can_view_data', 'can_view_file']
    }
  })

  const mutationCreateVerifyCode = useSWRMutation(
    'create-verifyCode',
    (_url, { arg }: { arg: any }) => createVerifyCode(arg),
    {
      onSuccess: () => {
        showNotification('success', 'Tạo mã xác minh thành công')
        form.reset()
        mutate(props.swrKey)
        props.handleSetOpen(false)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Tạo mã xác minh thất bại')
      }
    }
  )

  const onSubmit = (data: FormSchema) => {
    mutationCreateVerifyCode.trigger(data)
  }
  return (
    <Dialog open={props.open} onOpenChange={props.handleSetOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tạo mã xác minh</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <CustomFormItem
              type='input'
              control={form.control}
              name='expiredAfter'
              label='Thời gian hiệu lực (phút)'
              placeholder='Thời gian hiệu lực (phút)'
              setting={{
                input: {
                  type: 'number'
                }
              }}
            />
            <FormField
              control={form.control}
              name='permissionType'
              render={() => (
                <FormItem>
                  <FormLabel>Quyền hạn</FormLabel>

                  {permissionType.map((item) => (
                    <FormField
                      key={item.id}
                      control={form.control}
                      name='permissionType'
                      render={({ field }) => {
                        return (
                          <FormItem key={item.id} className='flex flex-row items-center gap-2'>
                            <FormControl>
                              <Checkbox
                                checked={field.value?.includes(item.id)}
                                onCheckedChange={(checked) => {
                                  return checked
                                    ? field.onChange([
                                        ...field.value,
                                        item.id as 'can_view_score' | 'can_view_data' | 'can_view_file'
                                      ])
                                    : field.onChange(field.value?.filter((value) => value !== item.id))
                                }}
                              />
                            </FormControl>
                            <FormLabel className='!mt-0 text-sm font-normal'>{item.label}</FormLabel>
                          </FormItem>
                        )
                      }}
                    />
                  ))}
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <DialogClose asChild>
                <Button variant={'outline'}>Hủy bỏ</Button>
              </DialogClose>
              <Button type='submit' isLoading={mutationCreateVerifyCode.isMutating}>
                Tạo mã
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

export default CreateVerifyCodeDialog
