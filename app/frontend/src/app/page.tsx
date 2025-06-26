import Header from '@/components/common/header'
import XrmSvg from '../../public/assets/svg/xrm.svg'

import Image from 'next/image'
import { getSession } from '@/lib/auth/session'
import SearchVerifyCode from '@/components/common/search-verify-code'
const HomePage = async () => {
  const session = await getSession()

  return (
    <main>
      <Header role={session?.role ? (session.role as 'admin' | 'student' | 'university_admin') : null} />
      <div className='container mt-16 flex flex-col items-center justify-center pt-8'>
        <Image src={XrmSvg} alt='xrm' width={200} height={100} />
        <h1 className='mt-3 text-center text-2xl font-semibold sm:text-4xl md:mt-6'>
          Giải pháp <span className='text-blue-500'>quản lý văn bằng chứng chỉ </span> ứng dụng{' '}
          <span className='text-blue-500'>Blockchain</span>
        </h1>
        <p className='mt-3 text-center text-sm text-muted-foreground sm:text-lg'>
          Dự án Web3 được xây dựng trên nền tảng Blockchain đảm bảo tính minh bạch cho văn bằng chứng chỉ
        </p>
        <SearchVerifyCode />
      </div>
    </main>
  )
}

export default HomePage
