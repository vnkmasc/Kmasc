'use client'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form } from '@/components/ui/form'
import CustomFormItem from '@/components/common/ct-form-item'
import { Button } from '@/components/ui/button'
import Image from 'next/image'
import XrmSvg from '../../../../public/assets/svg/xrm.svg'
import background from '../../../../public/assets/images/background.jpg'
import Link from 'next/link'
import { validateEmail, validateNoEmpty } from '@/lib/utils/validators'

import { requestEducationSignUp } from '@/lib/api/auth'
import { showNotification } from '@/lib/utils/common'
import useSWRMutation from 'swr/mutation'
const formSchma = z.object({
  email: validateEmail,
  name: validateNoEmpty('Tên của trường'),
  code: validateNoEmpty('Mã trường'),
  address: validateNoEmpty('Địa chỉ trường')
})

const EducationSignUpPage = () => {
  const form = useForm<z.infer<typeof formSchma>>({
    resolver: zodResolver(formSchma),
    defaultValues: {
      email: '',
      name: '',
      code: '',
      address: ''
    }
  })

  const mutateRequestEducationSignUp = useSWRMutation(
    '/universities',
    (_, { arg }: { arg: any }) => requestEducationSignUp(arg.email, arg.name, arg.code, arg.address),
    {
      onSuccess: (data) => {
        showNotification('success', data.message || 'Gửi yêu cầu đăng ký thành công')
      },
      onError: (error) => {
        showNotification('error', error.message || error.error || 'Gửi yêu cầu đăng ký thất bại')
      }
    }
  )

  const handleSubmit = async (data: z.infer<typeof formSchma>) => {
    mutateRequestEducationSignUp.trigger(data)
  }

  return (
    <div className='relative bottom-0 left-0 right-0 top-0 h-screen'>
      <Image src={background} width={1500} height={1500} className='h-full w-full object-cover' alt='no-image' />
      <Dialog open>
        <DialogContent className='rounded-lg sm:max-w-[450px] [&>button]:hidden'>
          <DialogHeader>
            <DialogTitle>
              <div>
                <Image src={XrmSvg} alt='xrm' width={150} height={150} className='mx-auto' />
              </div>
              Đăng ký
            </DialogTitle>
            <DialogDescription>Chào mừng quản lý đào tạo đến với hệ thống</DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className='space-y-4'>
              <CustomFormItem
                type='input'
                control={form.control}
                name='email'
                label='Email'
                placeholder='Nhập email'
                setting={{ input: { type: 'email' } }}
              />
              <CustomFormItem
                type='input'
                control={form.control}
                name='code'
                label='Mã trường'
                placeholder='Nhập mã trường'
              />
              <CustomFormItem
                type='input'
                control={form.control}
                name='name'
                label='Tên trường'
                placeholder='Nhập tên trường'
              />

              <CustomFormItem
                type='input'
                control={form.control}
                name='address'
                label='Địa chỉ trường'
                placeholder='Nhập địa chỉ trường'
              />
              <Button type='submit' className='w-full' isLoading={mutateRequestEducationSignUp.isMutating}>
                Gửi yêu cầu đăng ký
              </Button>
            </form>
          </Form>

          <div className='text-center text-sm'>
            Đã có tài khoản?{' '}
            <Link className='underline underline-offset-4' href='/auth/sign-in'>
              Đăng nhập
            </Link>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default EducationSignUpPage
