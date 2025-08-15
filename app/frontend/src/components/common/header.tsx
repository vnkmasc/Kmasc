'use client'

import { Button } from '../ui/button'
import Image from 'next/image'
import { LogInIcon, MenuIcon, Settings } from 'lucide-react'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '../ui/sheet'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle
} from '../ui/navigation-menu'
import Link from 'next/link'
import logoKmasc from '../../../public/assets/images/logoKMA.png'
import UseBreakpoint from '@/lib/hooks/use-breakpoint'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '../ui/dropdown-menu'
import { useState } from 'react'
import SignoutDialog from './signout-dialog'
import ChangePassDialog from './change-pass-dialog'
import ThemeSwitch from './theme-switch'

interface Props {
  role: 'student' | 'university_admin' | 'admin' | null
}

const Header: React.FC<Props> = (props) => {
  const { md, lg } = UseBreakpoint()
  const educationAdminPages: { title: string; href: string }[] = [
    {
      title: md && !lg ? 'CN' : 'Chuyên ngành',
      href: '/education-admin/faculty-management'
    },
    {
      title: md && !lg ? 'SV' : 'Sinh viên',
      href: '/education-admin/student-management'
    },
    {
      title: md && !lg ? 'KT&KL' : 'Khen thưởng & Kỷ luật',
      href: '/education-admin/reward-discipline-management'
    },
    {
      title: md && !lg ? 'VB&CC' : 'Văn bằng & Chứng chỉ',
      href: '/education-admin/certificate-management'
    },
    {
      title: md && !lg ? 'VB số' : 'Văn bằng số',
      href: '/education-admin/digital-degree-management'
    },
    {
      title: md && !lg ? 'MB số' : 'Mẫu bằng số',
      href: '/education-admin/degree-template-management'
    }
  ]

  const studentPages: { title: string; href: string }[] = [
    {
      title: 'Thông tin cá nhân',
      href: '/student/information'
    },
    // {
    //   title: 'Kết quả học tập',
    //   href: '/student/score'
    // },
    {
      title: 'Văn bằng - chứng chỉ',
      href: '/student/certificate'
    }
  ]

  const adminPages: { title: string; href: string }[] = [
    {
      title: 'Tài khoản đào tạo',
      href: '/admin/education-management'
    },
    {
      title: 'Mẫu bằng số',
      href: '/admin/degree-template-management'
    }
  ]
  const [openSignoutDialog, setOpenSignoutDialog] = useState<boolean>(false)
  const [openChangePassDialog, setOpenChangePassDialog] = useState<boolean>(false)

  const navList =
    props.role === 'university_admin' ? educationAdminPages : props.role === 'admin' ? adminPages : studentPages
  return (
    <div className='fixed top-0 z-10 h-16 w-full shadow-lg'>
      <header className='container flex h-full items-center justify-between bg-white dark:bg-black'>
        {props.role !== null ? (
          <div className='flex gap-2 md:hidden'>
            <Sheet>
              <SheetTrigger>
                <div className='rounded-md border p-1 hover:bg-accent'>
                  <MenuIcon />
                </div>
              </SheetTrigger>

              <SheetContent side={'left'}>
                <SheetHeader className='mb-4'>
                  <SheetTitle className='text-start'>Chức năng</SheetTitle>
                </SheetHeader>
                {navList.map((item) => (
                  <Link href={item.href} key={item.href}>
                    <Button variant={'link'} className='flex-start mb-2 block pr-0'>
                      {item.title}
                    </Button>
                  </Link>
                ))}
              </SheetContent>
            </Sheet>
          </div>
        ) : null}

        <Link href='/'>
          <div className='flex items-center gap-1'>
            <Image src={logoKmasc} alt='logoKmasc' width={30} height={30} />
            <h1 className='text-lg font-semibold text-main sm:text-xl'>Kmasc</h1>
          </div>
        </Link>

        {props.role !== null ? (
          <NavigationMenu className='hidden md:flex md:gap-2'>
            <NavigationMenuList>
              {navList.map((item) => (
                <NavigationMenuItem key={item.href}>
                  <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
                    <Link href={item.href}>{item.title}</Link>
                  </NavigationMenuLink>
                </NavigationMenuItem>
              ))}
            </NavigationMenuList>
          </NavigationMenu>
        ) : null}

        {props.role !== null ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button size={'icon'}>
                <Settings />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align='end' className='w-40'>
              <DropdownMenuLabel>Cấu hình</DropdownMenuLabel>
              <DropdownMenuGroup>
                <DropdownMenuItem>Cấu hình ký số</DropdownMenuItem>
                <DropdownMenuItem>
                  <ThemeSwitch />
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuLabel>Tài khoản</DropdownMenuLabel>
              <DropdownMenuGroup>
                <DropdownMenuItem onClick={() => setOpenChangePassDialog(true)}>Đổi mật khẩu</DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setOpenSignoutDialog(true)}
                  className='text-destructive hover:!text-destructive'
                >
                  Đăng xuất
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : (
          <Link href='/auth/sign-in'>
            <Button>
              <LogInIcon /> <span className='hidden md:block'>Đăng nhập</span>
            </Button>
          </Link>
        )}

        <SignoutDialog open={openSignoutDialog} onOpenChange={setOpenSignoutDialog} />
        <ChangePassDialog open={openChangePassDialog} setOpen={setOpenChangePassDialog} />
      </header>
    </div>
  )
}

export default Header
