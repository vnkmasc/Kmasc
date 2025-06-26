'use client'

import CustomFormItem from '@/components/common/ct-form-item'
import { Button } from '@/components/ui/button'
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Form } from '@/components/ui/form'
import { type CustomZodFormItem } from '@/types/common'
import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

interface Props {
  items: CustomZodFormItem[]
  data: any
  mode: 'create' | 'update' | undefined
  // eslint-disable-next-line no-unused-vars
  handleSubmit: (data: any) => void
  handleClose: () => void
  title?: string
}

const DetailDialog: React.FC<Props> = (props) => {
  const [localMode, setLocalMode] = useState<'create' | 'update' | undefined>(props.mode)
  const [localData, setLocalData] = useState(props.data)

  const formSchema = z.object(
    props.items.reduce(
      (acc, obj) => {
        acc[obj.name] = obj.validator || z.any()
        return acc
      },
      {} as Record<string, z.ZodType>
    )
  )

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: props.items.reduce(
      (acc, obj) => {
        acc[obj.name] = obj.defaultValue || ''
        return acc
      },
      {} as Record<string, any>
    )
  })

  useEffect(() => {
    if (props.mode !== undefined) {
      setLocalMode(props.mode)
      setLocalData(props.data)
    }
  }, [props.mode, props.data])

  useEffect(() => {
    if (localMode === 'update' && localData) {
      for (const key in localData) {
        form.setValue(key, localData[key])
      }
    } else if (localMode === 'create') {
      form.reset()
    }
  }, [localMode, localData, form])

  const handleOpenChange = (open: boolean) => {
    if (open === false) {
      props.handleClose()

      setTimeout(() => {
        setLocalMode(undefined)
        setLocalData(null)
        form.reset()
      }, 150)
    }
  }

  return (
    <Dialog open={props.mode !== undefined} onOpenChange={handleOpenChange}>
      <DialogContent className='max-h-[80vh] overflow-y-scroll sm:max-w-[500px]'>
        <DialogHeader>
          <DialogTitle>{props.title || (localMode === 'create' ? 'Tạo mới dữ liệu' : 'Cập nhật dữ liệu')}</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(props.handleSubmit)} className='space-y-4'>
            {props.items.map((prop, index) => (
              <CustomFormItem {...prop} control={form.control} key={index} />
            ))}
            <DialogFooter>
              <DialogClose asChild>
                <Button variant='outline' type='button'>
                  Hủy bỏ
                </Button>
              </DialogClose>
              <Button type='submit'>{localMode === 'create' ? 'Tạo mới' : 'Cập nhật'}</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

export default DetailDialog
