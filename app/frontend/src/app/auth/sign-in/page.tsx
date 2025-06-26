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
import { signIn } from '@/lib/auth/auth'
import { showNotification } from '@/lib/utils/common'
import Link from 'next/link'
import { validateEmail, validatePassword } from '@/lib/utils/validators'
import { useRouter } from 'next/navigation'
import { LogInIcon, SchoolIcon, User } from 'lucide-react'
const formSchma = z.object({
  email: validateEmail,
  password: validatePassword
})

const SignInPage = () => {
  const router = useRouter()
  const form = useForm<z.infer<typeof formSchma>>({
    resolver: zodResolver(formSchma),
    defaultValues: {
      email: '',
      password: ''
    }
  })

  const handleSubmit = async (data: z.infer<typeof formSchma>) => {
    const res = await signIn(data)
    if (res === false) {
      showNotification('error', 'Email hoặc mật khẩu không chính xác')
    } else {
      showNotification('success', 'Đăng nhập thành công')
      router.refresh()
    }
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
              Đăng nhập
            </DialogTitle>
            <DialogDescription>Chào mừng bạn quay trở lại</DialogDescription>
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
                name='password'
                label='Mật khẩu'
                placeholder='Nhập mật khẩu'
                setting={{ input: { type: 'password' } }}
              />
              <Button type='submit' className='w-full' isLoading={form.formState.isSubmitting}>
                <LogInIcon /> Đăng nhập
              </Button>
            </form>
          </Form>

          <div className='relative'>
            <hr className='my-4' />
            <span className='absolute left-1/2 top-1 -translate-x-1/2 text-sm text-gray-500'>
              <div className='bg-white px-2 text-sm dark:bg-background'>hoặc</div>
            </span>
          </div>
          <p className='text-center text-sm'>
            Bạn chưa có tài khoản? <span className='underline'>Đăng ký với vai trò</span>
          </p>
          <div className='flex flex-col gap-2 md:flex-row'>
            <Link className='flex-1' href='/auth/sign-up'>
              <Button variant={'outline'} className='w-full'>
                {' '}
                <User /> Sinh viên
              </Button>
            </Link>
            <Link className='flex-1' href='/auth/education-sign-up'>
              <Button variant={'outline'} className='w-full'>
                {' '}
                <SchoolIcon />
                Quản lý đào tạo
              </Button>
            </Link>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default SignInPage
