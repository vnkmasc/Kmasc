import Image from 'next/image'
import logoKmasc from '../../../public/assets/images/logoKMA.png'
import { Facebook, Github, GraduationCap, MapPinned } from 'lucide-react'
import Link from 'next/link'

const defaultSocialLinks = [
  { icon: <GraduationCap className='size-5' />, href: 'https://actvn.edu.vn', label: 'Trang chủ học viện' },
  { icon: <Facebook className='size-5' />, href: 'https://www.facebook.com/hocvienkythuatmatma', label: 'Facebook' },
  { icon: <MapPinned className='size-5' />, href: 'https://maps.app.goo.gl/nH4ungjtTKWfox2c8', label: 'Địa chỉ' },
  { icon: <Github className='size-5' />, href: 'https://github.com/vnkmasc/Kmasc', label: 'Github' }
]

const Footer: React.FC = () => {
  return (
    <footer className='w-full border-gray-500 bg-gray-100 py-8 pt-8 dark:border-t dark:bg-background'>
      <div className='container space-y-4'>
        <div className='flex items-center gap-2 lg:justify-start'>
          <Link href='/'>
            <Image src={logoKmasc} alt='logo' title='Kmasc' width={32} height={32} className='h-8' />
          </Link>
          <h2 className='font-semibold text-main'>Kmasc</h2>
        </div>
        <p className='text-sm text-muted-foreground'>Giải pháp quản lý văn bằng chứng chỉ ứng dụng Blockchain.</p>
        <div className='flex flex-col justify-between gap-4 md:flex-row'>
          <p className='text-sm text-muted-foreground'>
            © 2025 Kmasc. Bản quyền thuộc về chuyên ngành CNTT Học Viện Kỹ Thuật Mật Mã.
          </p>
          <ul className='flex items-center space-x-6 text-muted-foreground'>
            {defaultSocialLinks.map((social, idx) => (
              <li key={idx} className='font-medium hover:text-primary'>
                <a href={social.href} aria-label={social.label} target='_blank'>
                  {social.icon}
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </footer>
  )
}

export default Footer
