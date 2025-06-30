import Header from '@/components/common/header'
import logoKmasc from '../../public/assets/images/logoKMA.png'

import Image from 'next/image'
import { getSession } from '@/lib/auth/session'
import SearchVerifyCode from '@/components/common/search-verify-code'
import Footer from '@/components/common/footer'
const HomePage = async () => {
  const session = await getSession()

  return (
    <main className='flex h-screen flex-col'>
      <Header role={session?.role ? (session.role as 'admin' | 'student' | 'university_admin') : null} />
      <section className='container mt-16 flex flex-1 flex-col items-center pt-8'>
        <div className='flex items-center gap-2'>
          <Image src={logoKmasc} alt='logoKmasc' width={50} height={50} />
          <h1 className='text-2xl font-semibold text-main sm:text-4xl'>VnKmasc</h1>
        </div>
        <h1 className='mt-3 text-center text-xl font-semibold sm:text-3xl md:mt-6'>
          Giải pháp <span className='text-main'>quản lý văn bằng chứng chỉ </span> ứng dụng{' '}
          <span className='text-main'>Blockchain</span>
        </h1>
        <p className='mt-3 text-center text-sm text-muted-foreground sm:text-lg'>
          Dự án Web3 được xây dựng trên nền tảng Blockchain đảm bảo tính minh bạch cho văn bằng & chứng chỉ
        </p>
        <SearchVerifyCode />
      </section>
      <Footer />
    </main>
  )
}

export default HomePage
