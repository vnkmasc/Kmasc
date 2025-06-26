'use client'
import CustomFormItem from '@/components/common/ct-form-item'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Form } from '@/components/ui/form'

import { CustomZodFormItem } from '@/types/common'

import { zodResolver } from '@hookform/resolvers/zod'
import { CircleXIcon, SearchIcon } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

interface Props {
  items: CustomZodFormItem[]
  handleSetFilter: React.Dispatch<React.SetStateAction<any>>
}

const Filter: React.FC<Props> = (props) => {
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
    defaultValues: {
      ...props.items.reduce(
        (acc, obj) => {
          acc[obj.name] = obj.defaultValue || ''
          return acc
        },
        {} as Record<string, string>
      )
    }
  })

  const onSubmit = (data: z.infer<typeof formSchema>) => {
    props.handleSetFilter(data)
  }
  const handleReset = () => {
    form.reset()
    props.handleSetFilter(form.getValues())
  }
  return (
    <Card>
      <CardHeader>
        <CardTitle>
          <div className='flex items-center justify-between'>
            Tìm kiếm
            <div>
              <Button variant='destructive' className='mr-2' onClick={handleReset}>
                <CircleXIcon />
                <span className='hidden sm:block'>Xóa bộ lọc</span>
              </Button>
              <Button onClick={form.handleSubmit(onSubmit)}>
                <SearchIcon />
                <span className='hidden sm:block'>Tìm kiếm</span>
              </Button>
            </div>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className='grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3 md:gap-4 lg:grid-cols-4 xl:grid-cols-5'>
              {props.items.map((prop, index) => (
                <CustomFormItem {...prop} control={form.control} key={index} />
              ))}
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  )
}

export default Filter
